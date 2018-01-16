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

func (rpcHdl *rpcServiceHandler) CreateVoltageSensor(config *platformd.VoltageSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeleteVoltageSensor(config *platformd.VoltageSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdateVoltageSensor(oldConfig *platformd.VoltageSensor, newConfig *platformd.VoltageSensor, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtVoltageSensorConfig(oldConfig)
	newCfg := convertRPCToObjFmtVoltageSensorConfig(newConfig)
	rv, err := api.UpdateVoltageSensor(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetVoltageSensor(Name string) (*platformd.VoltageSensor, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkVoltageSensor(fromIdx, count platformd.Int) (*platformd.VoltageSensorGetInfo, error) {
	var getBulkObj platformd.VoltageSensorGetInfo
	var err error

	info, err := api.GetBulkVoltageSensorConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.VoltageSensorList = append(getBulkObj.VoltageSensorList, convertToRPCFmtVoltageSensorConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetVoltageSensorState(Name string) (*platformd.VoltageSensorState, error) {
	var rpcObj *platformd.VoltageSensorState
	var err error

	obj, err := api.GetVoltageSensorState(Name)
	if err == nil {
		rpcObj = convertToRPCFmtVoltageSensorState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkVoltageSensorState(fromIdx, count platformd.Int) (*platformd.VoltageSensorStateGetInfo, error) {
	var getBulkObj platformd.VoltageSensorStateGetInfo
	var err error

	info, err := api.GetBulkVoltageSensorState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.VoltageSensorStateList = append(getBulkObj.VoltageSensorStateList, convertToRPCFmtVoltageSensorState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetVoltageSensorPMDataState(Name string, Class string) (*platformd.VoltageSensorPMDataState, error) {
	var rpcObj *platformd.VoltageSensorPMDataState
	var err error

	obj, err := api.GetVoltageSensorPMDataState(Name, Class)
	if err == nil {
		rpcObj = convertToRPCFmtVoltageSensorPMState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkVoltageSensorPMDataState(fromIdx, count platformd.Int) (*platformd.VoltageSensorPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) restoreVoltageSensorConfigFromDB() (bool, error) {
	var voltageSensorCfg objects.VoltageSensor
	voltageSensorCfgList, err := rpcHdl.dbHdl.GetAllObjFromDb(voltageSensorCfg)
	if err != nil {
		return false, errors.New("Failed to retrieve VoltageSensor config object from DB")
	}
	for idx := 0; idx < len(voltageSensorCfgList); idx++ {
		dbObj := voltageSensorCfgList[idx].(objects.VoltageSensor)
		obj := new(platformd.VoltageSensor)
		objects.ConvertplatformdVoltageSensorObjToThrift(&dbObj, obj)
		convNewCfg := convertRPCToObjFmtVoltageSensorConfig(obj)
		ok, err := api.UpdateVoltageSensor(nil, convNewCfg, nil)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}
