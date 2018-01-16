//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package pluginManager

import (
	"errors"
	"fmt"
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"utils/logging"
)

type PlatformManager struct {
	logger                  logging.LoggerIntf
	plugin                  PluginIntf
	PlatformMgmtDeviceState *pluginCommon.PlatformMgmtDeviceState
}

var PlatformMgr PlatformManager

func (pMgr *PlatformManager) Init(logger logging.LoggerIntf, plugin PluginIntf) {
	pMgr.logger = logger
	pMgr.logger.Info("PlarformManager Init Start")
	pMgr.plugin = plugin
	pMgr.PlatformMgmtDeviceState = new(pluginCommon.PlatformMgmtDeviceState)
}

func (pMgr *PlatformManager) Deinit() {
	pMgr.logger.Info("PlarformManager deinit")
}

func (pMgr *PlatformManager) GetPlatformMgmtDeviceState(Name string) (*objects.PlatformMgmtDeviceState, error) {
	var platformMgmtDeviceObj objects.PlatformMgmtDeviceState
	if pMgr.plugin == nil {
		return &platformMgmtDeviceObj, errors.New("Invalid platform plugin")
	}
	err := pMgr.plugin.GetPlatformMgmtDeviceState(pMgr.PlatformMgmtDeviceState)
	if err != nil {
		return &platformMgmtDeviceObj, errors.New("Unable to get the PlatformMgmtDeviceState")
	}
	platformMgmtDeviceObj.CPUUsage = pMgr.PlatformMgmtDeviceState.CPUUsage
	platformMgmtDeviceObj.Version = pMgr.PlatformMgmtDeviceState.Version
	platformMgmtDeviceObj.Description = pMgr.PlatformMgmtDeviceState.Description
	platformMgmtDeviceObj.DeviceName = pMgr.PlatformMgmtDeviceState.DeviceName
	platformMgmtDeviceObj.MemoryUsage = pMgr.PlatformMgmtDeviceState.MemoryUsage
	platformMgmtDeviceObj.ResetReason = pMgr.PlatformMgmtDeviceState.ResetReason
	platformMgmtDeviceObj.Uptime = pMgr.PlatformMgmtDeviceState.Uptime
	return &platformMgmtDeviceObj, err
}

func (pMgr *PlatformManager) GetBulkPlatformMgmtDeviceState(fromIdx int, count int) (*objects.PlatformMgmtDeviceStateGetInfo, error) {
	var retObj objects.PlatformMgmtDeviceStateGetInfo
	if pMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx > 0 {
		return nil, errors.New("Invalid range, There is only one object")
	}
	retObj.EndIdx = 1
	retObj.More = false
	retObj.Count = 1
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		obj, err := pMgr.GetPlatformMgmtDeviceState("BMC")
		if err != nil {
			pMgr.logger.Err(fmt.Sprintln("Error getting the platform management state for :BMC"))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}
