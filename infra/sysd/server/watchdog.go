//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"infra/sysd/sysdCommonDefs"
	"io/ioutil"
	"models/objects"
	"os"
	"os/exec"
	"strings"
	"sysd"
	"time"
)

const (
	KA_TIMEOUT_COUNT_MIN = 0
	KA_TIMEOUT_COUNT     = 5 // After 5 KA missed from a daemon, sysd will assume the daemon as non-responsive. Restart it.
	WD_MAX_NUM_RESTARTS  = 5 // After 5 restarts, if this daemon is still not responsive then stop it.
	KA_RECV_CH_SIZE      = 1024
)

const (
	REASON_NONE           = "None"
	REASON_KA_FAIL        = "Failed to receive keepalive messages"
	REASON_USER_RESTART   = "Restarted by user"
	REASON_DAEMON_STOPPED = "Stopped by user"
	REASON_COMING_UP      = "Started by user"
)

type DaemonJson struct {
	Name string `json:"Name"`
}

type DaemonEnableJson struct {
	Name   string `json:"Name"`
	Enable bool   `json:"Enable"`
}

type DaemonsEnableList struct {
	DaemonsEnable []DaemonEnableJson `json:"Daemons"`
}

type DaemonInfo struct {
	Enable        bool
	State         sysdCommonDefs.SRDaemonStatus
	Reason        string
	StartTime     string
	RecvedKACount int32
	NumRestarts   int32
	RestartTime   string
	RestartReason string
	WatchDog      bool
}

func (daemonInfo *DaemonInfo) Initialize() error {
	daemonInfo.Enable = true
	daemonInfo.State = sysdCommonDefs.STARTING
	daemonInfo.Reason = REASON_COMING_UP
	daemonInfo.StartTime = time.Now().String()
	daemonInfo.RecvedKACount = 0
	daemonInfo.NumRestarts = 0
	daemonInfo.RestartTime = ""
	daemonInfo.RestartReason = ""
	daemonInfo.WatchDog = false
	return nil
}

func (server *SYSDServer) StartWDRoutine() {
	server.DaemonMap = make(map[string]*DaemonInfo)
	server.KaRecvCh = make(chan string, KA_RECV_CH_SIZE)
	server.DaemonStateDBCh = make(chan string, KA_RECV_CH_SIZE)
	server.DaemonConfigCh = make(chan DaemonConfig, sysdCommonDefs.SYSD_TOTAL_KA_DAEMONS)
	server.DaemonRestartCh = make(chan string, sysdCommonDefs.SYSD_TOTAL_KA_DAEMONS)
	server.PopulateDaemonList()
	go server.KeepAliveReceiver()
	go server.DaemonStateDBWriter()
	for {
		select {
		case daemonConfig := <-server.DaemonConfigCh:
			server.logger.Info("Received daemon config for:", daemonConfig.Name, "Enable ", daemonConfig.Enable)
			daemon := daemonConfig.Name
			enable := daemonConfig.Enable
			watchDog := daemonConfig.WatchDog
			daemonInfo, exist := server.DaemonMap[daemon]
			daemonUpdated := false
			if exist {
				if enable {
					if daemonInfo.State == sysdCommonDefs.STOPPED {
						server.ToggleFlexswitchDaemon(daemon, true)
						daemonInfo.State = sysdCommonDefs.STARTING
						daemonInfo.Reason = REASON_COMING_UP
						daemonInfo.Enable = true
						daemonUpdated = true
					} else {
						server.logger.Info("Daemon", daemonConfig.Name, "is not in stopped state, ignoring enable command")
					}
				} else {
					if daemonInfo.State == sysdCommonDefs.UP {
						server.ToggleFlexswitchDaemon(daemon, false)
						daemonInfo.State = sysdCommonDefs.STOPPED
						daemonInfo.Reason = REASON_DAEMON_STOPPED
						daemonInfo.RecvedKACount = 0
						server.PublishDaemonKANotification(daemon, daemonInfo.State)
						daemonInfo.Enable = false
						daemonUpdated = true
					} else {
						server.logger.Info("Daemon", daemonConfig.Name, "is not in up state, ignoring disable command")
					}
				}
				if daemonInfo != nil {
					daemonInfo.WatchDog = watchDog
				}
				if daemonUpdated {
					server.DaemonStateDBCh <- daemon
				}
			} else {
				server.logger.Info("Unknown daemon:", daemonConfig.Name)
			}
		case daemon := <-server.DaemonRestartCh:
			daemonInfo, exist := server.DaemonMap[daemon]
			if exist && daemonInfo.State == sysdCommonDefs.UP {
				daemonInfo.State = sysdCommonDefs.RESTARTING
				daemonInfo.NumRestarts++
				daemonInfo.RestartTime = time.Now().String()
				daemonInfo.RestartReason = REASON_KA_FAIL
				go server.RestartFlexswitchDaemon(daemon)
				server.DaemonStateDBCh <- daemon
			}
		}
	}
}

/*
 * Read clients.json and populate DaemonMap. Keep the Daemon's watchdog as false so that sysd won't
 * try to restart it until first keepalive is received.
 * Also, read systemProfile.json and disable (as stopped) the daemons that have Enabled == false.
 */
func (server *SYSDServer) PopulateDaemonList() error {
	var daemonsList []DaemonJson
	var daemonsEnableList DaemonsEnableList
	var daemonName string

	daemonsFile := server.paramsDir + "clients.json"
	bytes, err := ioutil.ReadFile(daemonsFile)
	if err != nil {
		server.logger.Info("Failed to read daemons list file")
		return err
	}
	err = json.Unmarshal(bytes, &daemonsList)
	if err != nil {
		server.logger.Info("Failed to unmarshal daemons list json")
		return err
	}
	for _, daemon := range daemonsList {
		if daemon.Name == "sysd" {
			continue
		}
		if daemon.Name == "local" {
			daemonName = "confd"
		} else {
			daemonName = daemon.Name
		}
		daemonInfo := &DaemonInfo{}
		daemonInfo.Initialize()
		server.DaemonMap[daemonName] = daemonInfo
		server.logger.Info("Added daemon", daemonName)
	}

	daemonsEnableFile := server.paramsDir + "systemProfile.json"
	bytes, err = ioutil.ReadFile(daemonsEnableFile)
	if err != nil {
		server.logger.Info("Failed to read daemons enablefile")
		return err
	}
	err = json.Unmarshal(bytes, &daemonsEnableList)
	if err != nil {
		server.logger.Info("Failed to unmarshal daemons enable list json")
		return err
	}
	for _, daemon := range daemonsEnableList.DaemonsEnable {
		daemonInfo, exist := server.DaemonMap[daemon.Name]
		if exist {
			if daemon.Enable == false {
				daemonInfo.Enable = false
				daemonInfo.State = sysdCommonDefs.STOPPED
				daemonInfo.Reason = REASON_DAEMON_STOPPED
				server.logger.Info("Daemon", daemon.Name, "is stopped")
			}
		}
	}
	// Write all daemons' information into DB
	for daemonName, _ = range server.DaemonMap {
		server.UpdateDaemonStateInDb(daemonName)
	}
	return nil
}

func (server *SYSDServer) PublishDaemonKANotification(name string, status sysdCommonDefs.SRDaemonStatus) error {
	msg := sysdCommonDefs.DaemonStatus{
		Name:   name,
		Status: status,
	}
	msgBuf, err := json.Marshal(msg)
	if err != nil {
		server.logger.Err("Failed to marshal daemon status")
		return err
	}
	notification := sysdCommonDefs.Notification{
		Type:    uint8(sysdCommonDefs.KA_DAEMON),
		Payload: msgBuf,
	}
	notificationBuf, err := json.Marshal(notification)
	if err != nil {
		server.logger.Err("Failed to marshal daemon status message")
		return err
	}
	server.notificationCh <- notificationBuf
	return nil
}

func (server *SYSDServer) ToggleFlexswitchDaemon(daemon string, up bool) error {
	var (
		cmdOut []byte
		err    error
		op     string
	)
	cmdDir := strings.TrimSuffix(server.paramsDir, "params/")
	cmdName := cmdDir + "flexswitch"
	if up {
		op = "start"
	} else {
		op = "stop"
	}
	//cmdArgs := []string{"-n", daemon, "-o", op, "-d", cmdDir}
	cmdArgs := []string{"-n", daemon, "-o", op}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		server.logger.Info(os.Stderr, "There was an error to", op, "flexswitch daemon", daemon, " : ", err)
		return err
	}
	out := string(cmdOut)
	server.logger.Info("Flexswitch daemon", daemon, op, "returned", out)

	return nil
}

func (server *SYSDServer) RestartFlexswitchDaemon(daemon string) {
	server.PublishDaemonKANotification(daemon, sysdCommonDefs.RESTARTING)
	server.ToggleFlexswitchDaemon(daemon, false)
	server.ToggleFlexswitchDaemon(daemon, true)
}

func (server *SYSDServer) DaemonStateDBWriter() {
	for {
		select {
		case daemon := <-server.DaemonStateDBCh:
			server.UpdateDaemonStateInDb(daemon)
		}
	}
}

func (server *SYSDServer) KeepAliveReceiver() {
	wdTimerCh := time.NewTicker(time.Second * KA_TIMEOUT_COUNT).C
	for {
		select {
		case kaDaemon := <-server.KaRecvCh:
			daemonInfo, exist := server.DaemonMap[kaDaemon]
			if exist {
				if daemonInfo.Enable {
					if daemonInfo.WatchDog == false {
						daemonInfo.WatchDog = true
					}
					daemonInfo.RecvedKACount++
					if daemonInfo.State != sysdCommonDefs.UP {
						daemonInfo.State = sysdCommonDefs.UP
						daemonInfo.Reason = REASON_NONE
						server.PublishDaemonKANotification(kaDaemon, daemonInfo.State)
					}
					server.DaemonStateDBCh <- kaDaemon
				}
			} else {
				server.logger.Info("Received keepalive from unknown daemon", kaDaemon)
			}
		case <-wdTimerCh:
			for daemon, daemonInfo := range server.DaemonMap {
				if daemonInfo.WatchDog && daemonInfo.State == sysdCommonDefs.UP {
					if daemonInfo.RecvedKACount < KA_TIMEOUT_COUNT && daemonInfo.RecvedKACount > KA_TIMEOUT_COUNT_MIN {
						server.logger.Info("Daemon", daemon, "is slowing down. Monitoring it.")
					}
					if daemonInfo.RecvedKACount <= KA_TIMEOUT_COUNT_MIN {
						server.logger.Info("Daemon", daemon, "is not responsive. Restarting it.")
						server.DaemonRestartCh <- daemon
					}
				}
				daemonInfo.RecvedKACount = 0
			}
		}
	}
}

func (server *SYSDServer) ConvertDaemonStateToThrift(ent DaemonState) *sysd.DaemonState {
	dState := sysd.NewDaemonState()
	dState.Name = string(ent.Name)
	dState.Enable = ent.Enable
	dState.State = string(sysdCommonDefs.ConvertDaemonStateCodeToString(ent.State))
	dState.Reason = string(ent.Reason)
	dState.StartTime = string(ent.StartTime)
	kaStr := fmt.Sprintf("Received %d keepalives", ent.RecvedKACount)
	dState.KeepAlive = string(kaStr)
	dState.RestartCount = int32(ent.NumRestarts)
	dState.RestartTime = string(ent.RestartTime)
	dState.RestartReason = string(ent.RestartReason)
	return dState
}

func (server *SYSDServer) ConvertDaemonStateToObj(ent DaemonState) objects.DaemonState {
	kaStr := fmt.Sprintf("Received %d keepalives", ent.RecvedKACount)
	dState := objects.DaemonState{
		Name:          ent.Name,
		Enable:        ent.Enable,
		State:         sysdCommonDefs.ConvertDaemonStateCodeToString(ent.State),
		Reason:        ent.Reason,
		StartTime:     ent.StartTime,
		KeepAlive:     kaStr,
		RestartCount:  ent.NumRestarts,
		RestartTime:   ent.RestartTime,
		RestartReason: ent.RestartReason,
	}
	return dState
}

func (server *SYSDServer) GetDaemonState(name string) *DaemonState {
	daemonState := new(DaemonState)
	daemonInfo, found := server.DaemonMap[name]
	if found {
		daemonState.Name = name
		daemonState.Enable = daemonInfo.Enable
		daemonState.State = daemonInfo.State
		daemonState.Reason = daemonInfo.Reason
		daemonState.StartTime = daemonInfo.StartTime
		daemonState.RecvedKACount = daemonInfo.RecvedKACount
		daemonState.NumRestarts = daemonInfo.NumRestarts
		daemonState.RestartTime = daemonInfo.RestartTime
		daemonState.RestartReason = daemonInfo.RestartReason
	}
	return daemonState
}

func (server *SYSDServer) GetBulkDaemonStates(idx int, cnt int) (int, int, []DaemonState) {
	var nextIdx int
	var count int
	server.logger.Info("GetBulk DaemonStates")
	result := make([]DaemonState, cnt)
	i := 0
	for daemon, daemonInfo := range server.DaemonMap {
		result[i].Name = daemon
		result[i].Enable = daemonInfo.Enable
		result[i].State = daemonInfo.State
		result[i].Reason = daemonInfo.Reason
		result[i].StartTime = daemonInfo.StartTime
		result[i].RecvedKACount = daemonInfo.RecvedKACount
		result[i].NumRestarts = daemonInfo.NumRestarts
		result[i].RestartTime = daemonInfo.RestartTime
		result[i].RestartReason = daemonInfo.RestartReason
		i++
	}
	count = i
	nextIdx = 0
	return nextIdx, count, result
}

func (server *SYSDServer) UpdateDaemonStateInDb(name string) error {
	var err error
	daemonState := server.GetDaemonState(name)
	if daemonState != nil {
		obj := server.ConvertDaemonStateToObj(*daemonState)
		server.dbHdl.StoreObjectInDb(obj)
	} else {
		errStr := "Failed to get daemon " + name
		server.logger.Info(errStr)
		err = errors.New(errStr)
	}
	return err
}
