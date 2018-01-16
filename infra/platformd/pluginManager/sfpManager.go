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
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"utils/logging"
)

//type SfpId int32

type SfpManager struct {
	logger    logging.LoggerIntf
	plugin    PluginIntf
	sfpIdList []int32
	stateDB   map[int32]SfpState
	configDB  map[int32]SfpConfig
}

type SfpConfig struct {
	SfpId      int32
	AdminState string
}

type SfpState struct {
	SfpId      int32
	SfpSpeed   string
	SfpLos     string
	SfpPresent string
	SfpType    string
	SerialNum  string
	EEPROM     string
}

var SfpMgr SfpManager

func (sfpMgr *SfpManager) Init(logger logging.LoggerIntf, plugin PluginIntf) {
	sfpMgr.logger = logger
	sfpMgr.plugin = plugin

	sfpMgr.stateDB = make(map[int32]SfpState)
	sfpMgr.configDB = make(map[int32]SfpConfig)

	sfpCnt := sfpMgr.plugin.GetSfpCnt()
	sfpList := make([]pluginCommon.SfpState, sfpCnt)
	sfpMgr.plugin.GetAllSfpState(sfpList, sfpCnt)

	for _, sfp := range sfpList {
		sfpEnt, _ := sfpMgr.stateDB[int32(sfp.SfpId)]

		sfpEnt.SfpId = sfp.SfpId
		sfpEnt.SfpSpeed = sfp.SfpSpeed
		sfpEnt.SfpLos = sfp.SfpLos
		sfpEnt.SfpPresent = sfp.SfpPresent
		sfpEnt.SfpType = sfp.SfpType
		sfpEnt.SerialNum = sfp.SerialNum
		sfpEnt.EEPROM = sfp.EEPROM

		sfpMgr.sfpIdList = append(sfpMgr.sfpIdList, int32(sfp.SfpId))
	}
	sfpMgr.logger.Info("SFP Manager Init()")
}

func (sfpMgr *SfpManager) Deinit() {
	sfpMgr.logger.Info("SFP Manager Deinit()")
}

func (sfpMgr *SfpManager) GetSfpState(sfpId int32) (*objects.SfpState, error) {
	var obj objects.SfpState

	if sfpMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}

	sfpState, err := sfpMgr.plugin.GetSfpState(sfpId)
	if err != nil {
		return &obj, err
	}
	obj.SfpId = sfpState.SfpId
	obj.SfpSpeed = sfpState.SfpSpeed
	obj.SfpLos = sfpState.SfpLos
	obj.SfpPresent = sfpState.SfpPresent
	obj.SfpType = sfpState.SfpType
	obj.SerialNum = sfpState.SerialNum
	obj.EEPROM = sfpState.EEPROM

	return &obj, nil
}

func (sfpMgr *SfpManager) GetBulkSfpState(fromIdx, cnt int) (*objects.SfpStateGetInfo, error) {
	var retObj objects.SfpStateGetInfo

	if sfpMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sfpMgr.sfpIdList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sfpMgr.sfpIdList) {
		retObj.EndIdx = len(sfpMgr.sfpIdList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sfpMgr.sfpIdList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		obj, err := sfpMgr.GetSfpState(int32(idx))
		if err != nil {
			sfpMgr.logger.Err("Error getting the SFP state for sfpId:", idx)
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (sfpMgr *SfpManager) GetSfpConfig(spfID int32) (*objects.SfpConfig, error) {
	var obj objects.SfpConfig
	return &obj, nil
}

func (sfpMgr *SfpManager) GetBulkSfpConfig(fromIdx, count int) (*objects.SfpConfigGetInfo, error) {
	var obj objects.SfpConfigGetInfo
	return &obj, nil
}

func (sfpMgr *SfpManager) UpdateSfpConfig(oldCfg *objects.SfpConfig, newCfg *objects.SfpConfig, attrset []bool) (bool, error) {
	return false, nil
}
