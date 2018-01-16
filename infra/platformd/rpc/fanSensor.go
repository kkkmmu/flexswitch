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

func (rpcHdl *rpcServiceHandler) CreateFanSensor(config *platformd.FanSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeleteFanSensor(config *platformd.FanSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdateFanSensor(oldConfig *platformd.FanSensor, newConfig *platformd.FanSensor, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtFanSensorConfig(oldConfig)
	newCfg := convertRPCToObjFmtFanSensorConfig(newConfig)
	rv, err := api.UpdateFanSensor(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetFanSensor(Name string) (*platformd.FanSensor, error) {
	return nil, errors.New("Not Supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkFanSensor(fromIdx, count platformd.Int) (*platformd.FanSensorGetInfo, error) {
	var getBulkObj platformd.FanSensorGetInfo
	var err error

	info, err := api.GetBulkFanSensorConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.FanSensorList = append(getBulkObj.FanSensorList, convertToRPCFmtFanSensorConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetFanSensorState(Name string) (*platformd.FanSensorState, error) {
	var rpcObj *platformd.FanSensorState
	var err error

	obj, err := api.GetFanSensorState(Name)
	if err == nil {
		rpcObj = convertToRPCFmtFanSensorState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkFanSensorState(fromIdx, count platformd.Int) (*platformd.FanSensorStateGetInfo, error) {
	var getBulkObj platformd.FanSensorStateGetInfo
	var err error

	info, err := api.GetBulkFanSensorState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.FanSensorStateList = append(getBulkObj.FanSensorStateList, convertToRPCFmtFanSensorState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetFanSensorPMDataState(Name string, Class string) (*platformd.FanSensorPMDataState, error) {
	var rpcObj *platformd.FanSensorPMDataState
	var err error

	obj, err := api.GetFanSensorPMDataState(Name, Class)
	if err == nil {
		rpcObj = convertToRPCFmtFanSensorPMState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkFanSensorPMDataState(fromIdx, count platformd.Int) (*platformd.FanSensorPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) restoreFanSensorConfigFromDB() (bool, error) {
	var fanSensorCfg objects.FanSensor
	fanSensorCfgList, err := rpcHdl.dbHdl.GetAllObjFromDb(fanSensorCfg)
	if err != nil {
		return false, errors.New("Failed to retrieve FanSensor config object from DB")
	}
	for idx := 0; idx < len(fanSensorCfgList); idx++ {
		dbObj := fanSensorCfgList[idx].(objects.FanSensor)
		obj := new(platformd.FanSensor)
		objects.ConvertplatformdFanSensorObjToThrift(&dbObj, obj)
		convNewCfg := convertRPCToObjFmtFanSensorConfig(obj)
		ok, err := api.UpdateFanSensor(nil, convNewCfg, nil)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}
