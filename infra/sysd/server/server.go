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
	nanomsg "github.com/op/go-nanomsg"
	"infra/sysd/iptables"
	"infra/sysd/sysdCommonDefs"
	"models/objects"
	"os"
	"os/signal"
	"syscall"
	"sysd"
	"utils/dbutils"
	"utils/logging"
)

type GlobalLoggingConfig struct {
	Level sysdCommonDefs.SRDebugLevel
}

type ComponentLoggingConfig struct {
	Component string
	Level     sysdCommonDefs.SRDebugLevel
}

type DaemonConfig struct {
	Name     string
	Enable   bool
	WatchDog bool
}

type DaemonState struct {
	Name          string
	Enable        bool
	State         sysdCommonDefs.SRDaemonStatus
	Reason        string
	StartTime     string
	RecvedKACount int32
	NumRestarts   int32
	RestartTime   string
	RestartReason string
}

type SYSDServer struct {
	logger                   *logging.Writer
	dbHdl                    *dbutils.DBUtil
	ServerStartedCh          chan bool
	paramsDir                string
	GlobalLoggingConfigCh    chan GlobalLoggingConfig
	ComponentLoggingConfigCh chan ComponentLoggingConfig
	sysdPubSocket            *nanomsg.PubSocket
	sysdIpTableMgr           *ipTable.SysdIpTableHandler
	notificationCh           chan []byte
	IptableAddCh             chan *sysd.IpTableAcl
	IptableDelCh             chan *sysd.IpTableAcl
	SystemParamConfig        chan objects.SystemParam
	KaRecvCh                 chan string
	DaemonMap                map[string]*DaemonInfo
	DaemonConfigCh           chan DaemonConfig
	DaemonRestartCh          chan string
	SysInfo                  *objects.SystemParam
	SysUpdCh                 chan *SystemParamUpdate
	DaemonStateDBCh          chan string
}

func NewSYSDServer(logger *logging.Writer, dbHdl *dbutils.DBUtil, paramsDir string) *SYSDServer {
	sysdServer := &SYSDServer{}
	sysdServer.sysdIpTableMgr = ipTable.SysdNewSysdIpTableHandler(logger)
	sysdServer.logger = logger
	sysdServer.dbHdl = dbHdl
	sysdServer.paramsDir = paramsDir
	sysdServer.ServerStartedCh = make(chan bool)
	sysdServer.GlobalLoggingConfigCh = make(chan GlobalLoggingConfig)
	sysdServer.ComponentLoggingConfigCh = make(chan ComponentLoggingConfig)
	sysdServer.notificationCh = make(chan []byte)
	sysdServer.IptableAddCh = make(chan *sysd.IpTableAcl)
	sysdServer.IptableDelCh = make(chan *sysd.IpTableAcl)
	sysdServer.SystemParamConfig = make(chan objects.SystemParam)
	sysdServer.SysUpdCh = make(chan *SystemParamUpdate)
	return sysdServer
}

func (server *SYSDServer) SigHandler(dbHdl *dbutils.DBUtil) {
	server.logger.Info("Starting SigHandler")
	sigChan := make(chan os.Signal, 1)
	signalList := []os.Signal{syscall.SIGHUP}
	signal.Notify(sigChan, signalList...)

	for {
		select {
		case signal := <-sigChan:
			switch signal {
			case syscall.SIGHUP:
				server.logger.Info("Received SIGHUP signal. Exiting")
				dbHdl.Disconnect()
				os.Exit(0)
			default:
				server.logger.Info("Unhandled signal : ", signal)
			}
		}
	}
}

func (server *SYSDServer) InitServer() {
	server.logger.Info("Initializing Sysd Server")
}

func (server *SYSDServer) InitPublisher(pub_str string) (pub *nanomsg.PubSocket) {
	server.logger.Info("Setting up ", pub_str, "publisher")
	pub, err := nanomsg.NewPubSocket()
	if err != nil {
		server.logger.Info("Failed to open pub socket")
		return nil
	}
	ep, err := pub.Bind(pub_str)
	if err != nil {
		server.logger.Info("Failed to bind pub socket - ", ep)
		return nil
	}
	err = pub.SetSendBuffer(1024)
	if err != nil {
		server.logger.Info("Failed to set send buffer size")
		return nil
	}
	return pub
}

func (server *SYSDServer) PublishSysdNotifications() {
	server.sysdPubSocket = server.InitPublisher(sysdCommonDefs.PUB_SOCKET_ADDR)
	for {
		select {
		case event := <-server.notificationCh:
			_, err := server.sysdPubSocket.Send(event, nanomsg.DontWait)
			if err == syscall.EAGAIN {
				server.logger.Err("Failed to publish event")
			}
		}
	}
}

func (server *SYSDServer) ProcessGlobalLoggingConfig(gLogConf GlobalLoggingConfig) error {
	server.logger.SetLevel(gLogConf.Level)
	server.logger.UpdateComponentLoggingInDb()
	msg := sysdCommonDefs.GlobalLogging{
		Level: gLogConf.Level,
	}
	msgBuf, err := json.Marshal(msg)
	if err != nil {
		server.logger.Err("Failed to marshal Global logging message")
		return err
	}
	notification := sysdCommonDefs.Notification{
		Type:    uint8(sysdCommonDefs.G_LOG),
		Payload: msgBuf,
	}
	notificationBuf, err := json.Marshal(notification)
	if err != nil {
		server.logger.Err("Failed to marshal Global logging message")
		return err
	}
	server.notificationCh <- notificationBuf
	return nil
}

func (server *SYSDServer) ProcessComponentLoggingConfig(cLogConf ComponentLoggingConfig) error {
	if cLogConf.Component == server.logger.MyComponentName {
		server.logger.SetLevel(cLogConf.Level)
	} else {
		msg := sysdCommonDefs.ComponentLogging{
			Name:  cLogConf.Component,
			Level: cLogConf.Level,
		}
		msgBuf, err := json.Marshal(msg)
		if err != nil {
			server.logger.Err("Failed to marshal Global logging message")
			return err
		}
		notification := sysdCommonDefs.Notification{
			Type:    uint8(sysdCommonDefs.C_LOG),
			Payload: msgBuf,
		}
		notificationBuf, err := json.Marshal(notification)
		if err != nil {
			server.logger.Err("Failed to marshal Global logging message")
			return err
		}
		server.notificationCh <- notificationBuf
	}
	return nil
}

func (server *SYSDServer) StartServer() {
	// Start notification publish thread
	go server.PublishSysdNotifications()
	// Start watchdog routine
	go server.StartWDRoutine()
	server.ReadSystemInfoFromDB()
	server.ServerStartedCh <- true
	// Now, wait on below channels to process
	for {
		select {
		case gLogConf := <-server.GlobalLoggingConfigCh:
			server.logger.Info("Received call for performing Global logging Configuration", gLogConf)
			server.ProcessGlobalLoggingConfig(gLogConf)
		case compLogConf := <-server.ComponentLoggingConfigCh:
			server.logger.Info("Received call for performing Component logging Configuration", compLogConf)
			server.ProcessComponentLoggingConfig(compLogConf)
		case addConfig := <-server.IptableAddCh:
			server.sysdIpTableMgr.AddIpRule(addConfig, false /*non-restart*/)
		case delConfig := <-server.IptableDelCh:
			server.sysdIpTableMgr.DelIpRule(delConfig)
		case sysConfig := <-server.SystemParamConfig:
			server.InitSystemInfo(sysConfig)
		case updateInfo := <-server.SysUpdCh:
			server.UpdateSystemInfo(updateInfo)
		}
	}
}
