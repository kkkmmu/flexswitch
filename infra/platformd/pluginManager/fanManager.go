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

type FanId int32

type FanState struct {
	OperMode      string
	OperSpeed     int32
	OperDirection string
	Status        string
	Model         string
	SerialNum     string
	LedId         int32
}

type FanConfig struct {
	AdminSpeed int32
	AdminState string
}

type FanManager struct {
	logger    logging.LoggerIntf
	plugin    PluginIntf
	fanIdList []FanId
	stateDB   map[FanId]FanState
	configDB  map[FanId]FanConfig
}

var FanMgr FanManager

func (fMgr *FanManager) Init(logger logging.LoggerIntf, plugin PluginIntf) {
	fMgr.logger = logger
	fMgr.plugin = plugin
	fMgr.stateDB = make(map[FanId]FanState)
	fMgr.configDB = make(map[FanId]FanConfig)
	numOfFans := fMgr.plugin.GetMaxNumOfFans()
	fanState := make([]pluginCommon.FanState, numOfFans)
	fMgr.plugin.GetAllFanState(fanState, numOfFans)
	for _, fan := range fanState {
		if fan.Valid == false {
			continue
		}
		fanEnt, _ := fMgr.stateDB[FanId(fan.FanId)]
		fanEnt.OperMode = fan.OperMode
		fanEnt.OperSpeed = fan.OperSpeed
		fanEnt.OperDirection = fan.OperDirection
		fanEnt.Status = fan.Status
		fanEnt.Model = fan.Model
		fanEnt.SerialNum = fan.SerialNum
		fMgr.stateDB[FanId(fan.FanId)] = fanEnt
		fanCfgEnt, _ := fMgr.configDB[FanId(fan.FanId)]
		fanCfgEnt.AdminState = fan.OperMode
		fanCfgEnt.AdminSpeed = fan.OperSpeed
		fMgr.configDB[FanId(fan.FanId)] = fanCfgEnt
		fMgr.fanIdList = append(fMgr.fanIdList, FanId(fan.FanId))
	}
	fMgr.logger.Info("Fan Manager Init()")
}

func (fMgr *FanManager) Deinit() {
	fMgr.logger.Info("Fan Manager Deinit()")
}

func (fMgr *FanManager) GetFanState(fanId int32) (*objects.FanState, error) {
	var fanObj objects.FanState
	if fMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	fanStateEnt, exist := fMgr.stateDB[FanId(fanId)]
	if !exist {
		return nil, errors.New("Invalid FanId")
	}

	fanState, err := fMgr.plugin.GetFanState(fanId)
	if err != nil {
		return nil, err
	}
	fanStateEnt.OperMode = fanState.OperMode
	fanStateEnt.OperSpeed = fanState.OperSpeed
	fanStateEnt.OperDirection = fanState.OperDirection
	fanStateEnt.Status = fanState.Status
	fanStateEnt.Model = fanState.Model
	fanStateEnt.SerialNum = fanState.SerialNum
	fMgr.stateDB[FanId(fanId)] = fanStateEnt
	fanObj.FanId = fanId
	fanObj.OperMode = fanState.OperMode
	fanObj.OperSpeed = fanState.OperSpeed
	fanObj.OperDirection = fanState.OperDirection
	fanObj.Status = fanState.Status
	fanObj.Model = fanState.Model
	fanObj.SerialNum = fanState.SerialNum
	fanObj.LedId = fanState.LedId
	return &fanObj, err
}

func (fMgr *FanManager) GetBulkFanState(fromIdx int, cnt int) (*objects.FanStateGetInfo, error) {
	var retObj objects.FanStateGetInfo
	if fMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(fMgr.fanIdList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(fMgr.fanIdList) {
		retObj.EndIdx = len(fMgr.fanIdList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(fMgr.fanIdList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		fanId := int32(fMgr.fanIdList[idx])
		obj, err := fMgr.GetFanState(fanId)
		if err != nil {
			fMgr.logger.Err(fmt.Sprintln("Error getting the fan state for fanId:", fanId))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (fMgr *FanManager) GetFanConfig(fanId int32) (*objects.FanConfig, error) {
	var fanObj objects.FanConfig
	if fMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	fanCfgEnt, exist := fMgr.configDB[FanId(fanId)]
	if !exist {
		return nil, errors.New("Invalid FanId")
	}
	fanObj.FanId = fanId
	fanObj.AdminSpeed = fanCfgEnt.AdminSpeed
	fanObj.AdminState = fanCfgEnt.AdminState
	return &fanObj, nil
}

func (fMgr *FanManager) GetBulkFanConfig(fromIdx int, cnt int) (*objects.FanConfigGetInfo, error) {
	var retObj objects.FanConfigGetInfo
	if fMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(fMgr.fanIdList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(fMgr.fanIdList) {
		retObj.EndIdx = len(fMgr.fanIdList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(fMgr.fanIdList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		fanId := int32(fMgr.fanIdList[idx])
		obj, err := fMgr.GetFanConfig(fanId)
		if err != nil {
			fMgr.logger.Err(fmt.Sprintln("Error getting the fan state for fanId:", fanId))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (fMgr *FanManager) UpdateFanConfig(oldCfg *objects.FanConfig, newCfg *objects.FanConfig, attrset []bool) (bool, error) {
	if fMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	ret, err := fMgr.plugin.UpdateFanConfig(newCfg)
	return ret, err
}
