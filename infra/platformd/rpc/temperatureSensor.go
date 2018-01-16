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

func (rpcHdl *rpcServiceHandler) CreateTemperatureSensor(config *platformd.TemperatureSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeleteTemperatureSensor(config *platformd.TemperatureSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdateTemperatureSensor(oldConfig *platformd.TemperatureSensor, newConfig *platformd.TemperatureSensor, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtTemperatureSensorConfig(oldConfig)
	newCfg := convertRPCToObjFmtTemperatureSensorConfig(newConfig)
	rv, err := api.UpdateTemperatureSensor(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetTemperatureSensor(Name string) (*platformd.TemperatureSensor, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkTemperatureSensor(fromIdx, count platformd.Int) (*platformd.TemperatureSensorGetInfo, error) {
	var getBulkObj platformd.TemperatureSensorGetInfo
	var err error

	info, err := api.GetBulkTemperatureSensorConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.TemperatureSensorList = append(getBulkObj.TemperatureSensorList, convertToRPCFmtTemperatureSensorConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetTemperatureSensorState(Name string) (*platformd.TemperatureSensorState, error) {
	var rpcObj *platformd.TemperatureSensorState
	var err error

	obj, err := api.GetTemperatureSensorState(Name)
	if err == nil {
		rpcObj = convertToRPCFmtTemperatureSensorState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkTemperatureSensorState(fromIdx, count platformd.Int) (*platformd.TemperatureSensorStateGetInfo, error) {
	var getBulkObj platformd.TemperatureSensorStateGetInfo
	var err error

	info, err := api.GetBulkTemperatureSensorState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.TemperatureSensorStateList = append(getBulkObj.TemperatureSensorStateList, convertToRPCFmtTemperatureSensorState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetTemperatureSensorPMDataState(Name string, Class string) (*platformd.TemperatureSensorPMDataState, error) {
	var rpcObj *platformd.TemperatureSensorPMDataState
	var err error

	obj, err := api.GetTempSensorPMDataState(Name, Class)
	if err == nil {
		rpcObj = convertToRPCFmtTempSensorPMState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkTemperatureSensorPMDataState(fromIdx, count platformd.Int) (*platformd.TemperatureSensorPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) restoreTemperatureSensorConfigFromDB() (bool, error) {
	var temperatureSensorCfg objects.TemperatureSensor
	temperatureSensorCfgList, err := rpcHdl.dbHdl.GetAllObjFromDb(temperatureSensorCfg)
	if err != nil {
		return false, errors.New("Failed to retrieve TemperatureSensor config object from DB")
	}
	for idx := 0; idx < len(temperatureSensorCfgList); idx++ {
		dbObj := temperatureSensorCfgList[idx].(objects.TemperatureSensor)
		obj := new(platformd.TemperatureSensor)
		objects.ConvertplatformdTemperatureSensorObjToThrift(&dbObj, obj)
		convNewCfg := convertRPCToObjFmtTemperatureSensorConfig(obj)
		ok, err := api.UpdateTemperatureSensor(nil, convNewCfg, nil)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}
