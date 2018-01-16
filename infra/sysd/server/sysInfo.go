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
	"bufio"
	"encoding/json"
	"infra/sysd/sysdCommonDefs"
	"models/objects"
	"os"
	"os/exec"
	"strings"
)

type SystemParamUpdate struct {
	EntriesUpdated []string
	NewCfg         *objects.SystemParam
}

func (svr *SYSDServer) ReadSystemInfoFromDB() error {
	svr.logger.Info("Reading System Information From Db")
	dbHdl := svr.dbHdl
	if dbHdl != nil {
		var dbObj objects.SystemParam
		objList, err := dbHdl.GetAllObjFromDb(dbObj)
		if err != nil {
			svr.logger.Err("DB query failed for System Info")
			return err
		}
		svr.logger.Info("Total System Entries are", len(objList))
		for idx := 0; idx < len(objList); idx++ {
			if svr.SysInfo == nil {
				svr.SysInfo = &objects.SystemParam{}
			}
			dbObject := objList[idx].(objects.SystemParam)
			svr.SysInfo.Vrf = dbObject.Vrf
			svr.SysInfo.SwitchMac = dbObject.SwitchMac
			svr.SysInfo.MgmtIp = dbObject.MgmtIp
			svr.SysInfo.SwVersion = dbObject.SwVersion
			svr.SysInfo.Description = dbObject.Description
			svr.SysInfo.Hostname = dbObject.Hostname
			svr.logger.Info("System Information:", *svr.SysInfo)
			break
		}
	}
	svr.logger.Info("reading system info from db done")
	return nil
}

// Func to send update nanomsg update notification to all the dameons on the system
func (svr *SYSDServer) SendSystemUpdate() {
	svr.logger.Debug("Sysd Sending System Update Notification with type:", sysdCommonDefs.SYSTEM_Info)
	msgBuf, err := json.Marshal(svr.SysInfo)
	if err != nil {
		svr.logger.Err("Failed to marshal system information, err:", err)
		return
	}
	svr.logger.Debug("Sysd new system information being send out is", *svr.SysInfo)
	notification := sysdCommonDefs.Notification{
		Type:    uint8(sysdCommonDefs.SYSTEM_Info),
		Payload: msgBuf,
	}
	notificationBuf, err := json.Marshal(notification)
	if err != nil {
		svr.logger.Err("Failed to marshal system notification, err:", err)
		return
	}
	svr.notificationCh <- notificationBuf
	svr.logger.Debug("Sysd system update notification send out")
}

func getDistro() string {
	inFile, _ := os.Open("/etc/os-release")
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	idStr := ""
	idLikeStr := ""
	versionStr := ""
	for scanner.Scan() {
		str := scanner.Text()
		strs := strings.Split(str, "=")
		switch strs[0] {
		case "ID":
			idStr = strs[1]
		case "ID_LIKE":
			idLikeStr = strs[1]
		case "VERSION":
			versionStr = strs[1]
		}
	}
	distroStr := idStr + " " + idLikeStr + " " + versionStr
	return distroStr
}

func getKernel() string {
	cmd := exec.Command("uname", "-a")
	cmdOut, _ := cmd.Output()
	cmdOuts := strings.Split(string(cmdOut), " ")
	return cmdOuts[2]
}

func (svr *SYSDServer) copyAndSendSystemParam(cfg objects.SystemParam) {
	sysInfo := svr.SysInfo
	sysInfo.Vrf = cfg.Vrf
	sysInfo.SwitchMac = cfg.SwitchMac
	sysInfo.MgmtIp = cfg.MgmtIp
	sysInfo.SwVersion = cfg.SwVersion
	sysInfo.Description = cfg.Description
	sysInfo.Hostname = cfg.Hostname
	svr.SendSystemUpdate()
}

// Initialize system information using json file...or whatever other means are
func (svr *SYSDServer) InitSystemInfo(cfg objects.SystemParam) {
	if svr.SysInfo == nil {
		svr.SysInfo = &objects.SystemParam{}
	}
	svr.copyAndSendSystemParam(cfg)
}

// Helper function for NB listener to determine whether a global object is created or not
func (svr *SYSDServer) SystemInfoCreated() bool {
	return (svr.SysInfo != nil)
}

// During Get calls we will use below api to read from run-time information
func (svr *SYSDServer) GetSystemParam(name string) *objects.SystemParamState {
	if svr.SysInfo == nil || svr.SysInfo.Vrf != name {
		return nil
	}
	sysParamsInfo := new(objects.SystemParamState)

	sysParamsInfo.Vrf = svr.SysInfo.Vrf
	sysParamsInfo.SwitchMac = svr.SysInfo.SwitchMac
	sysParamsInfo.MgmtIp = svr.SysInfo.MgmtIp
	sysParamsInfo.SwVersion = svr.SysInfo.SwVersion
	sysParamsInfo.Description = svr.SysInfo.Description
	sysParamsInfo.Hostname = svr.SysInfo.Hostname
	sysParamsInfo.Distro = getDistro()
	sysParamsInfo.Kernel = getKernel()
	return sysParamsInfo
}

// Update runtime system param info and send a notification
func (svr *SYSDServer) UpdateSystemInfo(updateInfo *SystemParamUpdate) {
	svr.logger.Info("Updating System Information with newconfig:", *updateInfo.NewCfg)
	for _, entry := range updateInfo.EntriesUpdated {
		switch entry {
		case "Hostname":
			binary, lookErr := exec.LookPath("hostname")
			if lookErr != nil {
				svr.logger.Err("Error searching path for hostname", lookErr)
				continue
			}
			cmd := exec.Command(binary, updateInfo.NewCfg.Hostname)
			err := cmd.Run()
			if err != nil {
				svr.logger.Err("Updating hostname in linux failed", err)
			}
			svr.logger.Info("Updated System Host Name", updateInfo.NewCfg.Hostname)
		case "Description":
			svr.SysInfo.Description = updateInfo.NewCfg.Description
			svr.logger.Info("Updated System Description", updateInfo.NewCfg.Description)
		}
	}
	svr.copyAndSendSystemParam(*updateInfo.NewCfg)
}
