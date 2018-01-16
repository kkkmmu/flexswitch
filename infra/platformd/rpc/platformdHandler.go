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

package rpc

import (
	"infra/platformd/api"
	"platformd"
)

func (rpcHdl *rpcServiceHandler) GetPlatformState(objName string) (*platformd.PlatformState, error) {
	var rpcObj *platformd.PlatformState

	obj, err := api.GetPlatformState(objName)
	if err == nil {
		rpcObj = convertToRPCFmtPlatformState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkPlatformState(fromIdx, count platformd.Int) (*platformd.PlatformStateGetInfo, error) {
	var getBulkObj platformd.PlatformStateGetInfo

	info, err := api.GetBulkPlatformState(int(fromIdx), int(count))
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.PlatformStateList = append(getBulkObj.PlatformStateList, convertToRPCFmtPlatformState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetFanState(fanId int32) (*platformd.FanState, error) {
	var rpcObj *platformd.FanState
	var err error

	obj, err := api.GetFanState(fanId)
	if err == nil {
		rpcObj = convertToRPCFmtFanState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkFanState(fromIdx, count platformd.Int) (*platformd.FanStateGetInfo, error) {
	var getBulkObj platformd.FanStateGetInfo
	var err error

	info, err := api.GetBulkFanState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.FanStateList = append(getBulkObj.FanStateList, convertToRPCFmtFanState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) CreateFan(config *platformd.Fan) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) DeleteFan(config *platformd.Fan) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) UpdateFan(oldConfig *platformd.Fan, newConfig *platformd.Fan, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtFanConfig(oldConfig)
	newCfg := convertRPCToObjFmtFanConfig(newConfig)
	rv, err := api.UpdateFan(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetFan(fanId int32) (*platformd.Fan, error) {
	var rpcObj *platformd.Fan
	var err error

	obj, err := api.GetFanConfig(fanId)
	if err == nil {
		rpcObj = convertToRPCFmtFanConfig(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkFan(fromIdx, count platformd.Int) (*platformd.FanGetInfo, error) {
	var getBulkObj platformd.FanGetInfo
	var err error

	info, err := api.GetBulkFanConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.FanList = append(getBulkObj.FanList, convertToRPCFmtFanConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkSfp(fromIdx, count platformd.Int) (*platformd.SfpGetInfo, error) {
	var getBulkObj platformd.SfpGetInfo
	var err error

	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) CreateSfp(config *platformd.Sfp) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) DeleteSfp(config *platformd.Sfp) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) UpdateSfp(oldConfig *platformd.Sfp, newConfig *platformd.Sfp, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	return false, nil
}

func (rpcHdl *rpcServiceHandler) GetSfp(sfpID int32) (*platformd.Sfp, error) {
	var rpcObj platformd.Sfp

	return &rpcObj, nil
}

func (rpcHdl *rpcServiceHandler) GetSfpState(sfpId int32) (*platformd.SfpState, error) {
	var rpcObj *platformd.SfpState

	obj, err := api.GetSfpState(sfpId)
	if err == nil {
		rpcObj = convertToRPCFmtSfpState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkSfpState(fromIdx, count platformd.Int) (*platformd.SfpStateGetInfo, error) {
	var getBulkObj platformd.SfpStateGetInfo
	var err error

	info, err := api.GetBulkSfpState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.SfpStateList = append(getBulkObj.SfpStateList, convertToRPCFmtSfpState(info.List[idx]))
	}
	return &getBulkObj, err
}

// TODO
func (rpcHdl *rpcServiceHandler) CreateLed(config *platformd.Led) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) DeleteLed(config *platformd.Led) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) UpdateLed(oldConfig *platformd.Led, newConfig *platformd.Led, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	return false, nil
}

func (rpcHdl *rpcServiceHandler) GetBulkLed(fromIdx, count platformd.Int) (*platformd.LedGetInfo, error) {
	var getBulkObj platformd.LedGetInfo
	var err error

	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetLed(ledID int32) (*platformd.Led, error) {
	var rpcObj platformd.Led

	return &rpcObj, nil
}

func (rpcHdl *rpcServiceHandler) GetBulkLedState(fromIdx, count platformd.Int) (*platformd.LedStateGetInfo, error) {
	var getBulkObj platformd.LedStateGetInfo
	var err error

	info, err := api.GetBulkLedState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.LedStateList = append(getBulkObj.LedStateList, convertToRPCFmtLedState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetLedState(ledId int32) (*platformd.LedState, error) {
	var rpcObj *platformd.LedState
	var err error

	obj, err := api.GetLedState(ledId)
	if err == nil {
		rpcObj = convertToRPCFmtLedState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkThermalState(fromIdx, count platformd.Int) (*platformd.ThermalStateGetInfo, error) {
	var getBulkObj platformd.ThermalStateGetInfo
	var err error

	info, err := api.GetBulkThermalState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.ThermalStateList = append(getBulkObj.ThermalStateList, convertToRPCFmtThermalState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetThermalState(thermalId int32) (*platformd.ThermalState, error) {
	var rpcObj *platformd.ThermalState
	var err error

	obj, err := api.GetThermalState(thermalId)
	if err == nil {
		rpcObj = convertToRPCFmtThermalState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) CreatePsu(config *platformd.Psu) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) DeletePsu(config *platformd.Psu) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) UpdatePsu(oldConfig *platformd.Psu, newConfig *platformd.Psu, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	return false, nil
}

func (rpcHdl *rpcServiceHandler) GetBulkPsu(fromIdx, count platformd.Int) (*platformd.PsuGetInfo, error) {
	var getBulkObj platformd.PsuGetInfo
	var err error

	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetPsu(psuID int32) (*platformd.Psu, error) {
	var rpcObj platformd.Psu

	return &rpcObj, nil
}

func (rpcHdl *rpcServiceHandler) GetBulkPsuState(fromIdx, count platformd.Int) (*platformd.PsuStateGetInfo, error) {
	var getBulkObj platformd.PsuStateGetInfo
	var err error

	info, err := api.GetBulkPsuState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.PsuStateList = append(getBulkObj.PsuStateList, convertToRPCFmtPsuState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetPsuState(psuId int32) (*platformd.PsuState, error) {
	var rpcObj *platformd.PsuState
	var err error

	obj, err := api.GetPsuState(psuId)
	if err == nil {
		rpcObj = convertToRPCFmtPsuState(obj)
	}
	return rpcObj, err
}
