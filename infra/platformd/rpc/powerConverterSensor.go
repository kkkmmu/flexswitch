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
	"errors"
	"infra/platformd/api"
	"models/objects"
	"platformd"
)

func (rpcHdl *rpcServiceHandler) CreatePowerConverterSensor(config *platformd.PowerConverterSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeletePowerConverterSensor(config *platformd.PowerConverterSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdatePowerConverterSensor(oldConfig *platformd.PowerConverterSensor, newConfig *platformd.PowerConverterSensor, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtPowerConverterSensorConfig(oldConfig)
	newCfg := convertRPCToObjFmtPowerConverterSensorConfig(newConfig)
	rv, err := api.UpdatePowerConverterSensor(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetPowerConverterSensor(Name string) (*platformd.PowerConverterSensor, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkPowerConverterSensor(fromIdx, count platformd.Int) (*platformd.PowerConverterSensorGetInfo, error) {
	var getBulkObj platformd.PowerConverterSensorGetInfo
	var err error

	info, err := api.GetBulkPowerConverterSensorConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.PowerConverterSensorList = append(getBulkObj.PowerConverterSensorList, convertToRPCFmtPowerConverterSensorConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetPowerConverterSensorState(Name string) (*platformd.PowerConverterSensorState, error) {
	var rpcObj *platformd.PowerConverterSensorState
	var err error

	obj, err := api.GetPowerConverterSensorState(Name)
	if err == nil {
		rpcObj = convertToRPCFmtPowerConverterSensorState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkPowerConverterSensorState(fromIdx, count platformd.Int) (*platformd.PowerConverterSensorStateGetInfo, error) {
	var getBulkObj platformd.PowerConverterSensorStateGetInfo
	var err error

	info, err := api.GetBulkPowerConverterSensorState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.PowerConverterSensorStateList = append(getBulkObj.PowerConverterSensorStateList, convertToRPCFmtPowerConverterSensorState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetPowerConverterSensorPMDataState(Name string, Class string) (*platformd.PowerConverterSensorPMDataState, error) {
	var rpcObj *platformd.PowerConverterSensorPMDataState
	var err error

	obj, err := api.GetPowerConverterSensorPMDataState(Name, Class)
	if err == nil {
		rpcObj = convertToRPCFmtPowerConverterSensorPMState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkPowerConverterSensorPMDataState(fromIdx, count platformd.Int) (*platformd.PowerConverterSensorPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) restorePowerConverterSensorConfigFromDB() (bool, error) {
	var powerConverterSensorCfg objects.PowerConverterSensor
	powerConverterSensorCfgList, err := rpcHdl.dbHdl.GetAllObjFromDb(powerConverterSensorCfg)
	if err != nil {
		return false, errors.New("Failed to retrieve PowerConverterSensor config object from DB")
	}
	for idx := 0; idx < len(powerConverterSensorCfgList); idx++ {
		dbObj := powerConverterSensorCfgList[idx].(objects.PowerConverterSensor)
		obj := new(platformd.PowerConverterSensor)
		objects.ConvertplatformdPowerConverterSensorObjToThrift(&dbObj, obj)
		convNewCfg := convertRPCToObjFmtPowerConverterSensorConfig(obj)
		ok, err := api.UpdatePowerConverterSensor(nil, convNewCfg, nil)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}
