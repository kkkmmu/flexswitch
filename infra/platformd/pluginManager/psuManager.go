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
	"utils/logging"
)

type PsuManager struct {
	logger logging.LoggerIntf
	plugin PluginIntf
	psuCnt int
}

var PsuMgr PsuManager

func (psuMgr *PsuManager) Init(logger logging.LoggerIntf, plugin PluginIntf) {
	psuMgr.logger = logger
	psuMgr.plugin = plugin
	psuMgr.psuCnt = psuMgr.plugin.GetMaxNumOfPsu()
	psuMgr.logger.Info("PSU Manager Init()")
}

func (psuMgr *PsuManager) Deinit() {
	psuMgr.logger.Info("PSU Manager Deinit()")
}

func (psuMgr *PsuManager) GetPsuState(psuId int32) (*objects.PsuState, error) {
	var psuObj objects.PsuState

	if psuMgr.plugin == nil {
		return nil, errors.New("Invalid PSU platform plugin")
	}

	psuState, err := psuMgr.plugin.GetPsuState(psuId)
	if err != nil {
		return nil, err
	}

	psuObj.PsuId = psuId
	psuObj.AdminState = psuState.Status
	psuObj.ModelNum = psuState.Model
	psuObj.SerialNum = psuState.SerialNum
	psuObj.Vin = psuState.VoltIn
	psuObj.Vout = psuState.VoltOut
	psuObj.Iin = psuState.AmpIn
	psuObj.Iout = psuState.AmpOut
	psuObj.Pin = psuState.PwrIn
	psuObj.Pout = psuState.PwrOut

	return &psuObj, err
}

func (psuMgr *PsuManager) GetBulkPsuState(fromIdx int, cnt int) (*objects.PsuStateGetInfo, error) {
	var retObj objects.PsuStateGetInfo

	if psuMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}

	if fromIdx >= psuMgr.psuCnt {
		return nil, errors.New("Invalid range")
	}

	if fromIdx+cnt > psuMgr.psuCnt {
		retObj.EndIdx = psuMgr.psuCnt
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = psuMgr.psuCnt - retObj.EndIdx + 1
	}

	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		psuId := int32(idx)
		obj, err := psuMgr.GetPsuState(psuId)
		if err != nil {
			psuMgr.logger.Err(fmt.Sprintln("Error getting the PSU state for psuId:", psuId))
		}
		retObj.List = append(retObj.List, obj)
	}

	return &retObj, nil
}
