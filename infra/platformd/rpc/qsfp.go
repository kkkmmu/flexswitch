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

func (rpcHdl *rpcServiceHandler) CreateQsfp(config *platformd.Qsfp) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeleteQsfp(config *platformd.Qsfp) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdateQsfp(oldConfig *platformd.Qsfp, newConfig *platformd.Qsfp, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtQsfpConfig(oldConfig)
	newCfg := convertRPCToObjFmtQsfpConfig(newConfig)
	rv, err := api.UpdateQsfp(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetQsfp(QsfpId int32) (*platformd.Qsfp, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkQsfp(fromIdx, count platformd.Int) (*platformd.QsfpGetInfo, error) {
	var getBulkObj platformd.QsfpGetInfo
	var err error

	info, err := api.GetBulkQsfpConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.QsfpList = append(getBulkObj.QsfpList, convertToRPCFmtQsfpConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetQsfpState(QsfpId int32) (*platformd.QsfpState, error) {
	var rpcObj *platformd.QsfpState
	var err error

	obj, err := api.GetQsfpState(QsfpId)
	if err == nil {
		rpcObj = convertToRPCFmtQsfpState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkQsfpState(fromIdx, count platformd.Int) (*platformd.QsfpStateGetInfo, error) {
	var getBulkObj platformd.QsfpStateGetInfo
	var err error

	info, err := api.GetBulkQsfpState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.QsfpStateList = append(getBulkObj.QsfpStateList, convertToRPCFmtQsfpState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetQsfpPMDataState(QsfpId int32, Resource string, Class string) (*platformd.QsfpPMDataState, error) {
	var rpcObj *platformd.QsfpPMDataState
	var err error

	obj, err := api.GetQsfpPMDataState(QsfpId, Resource, Class)
	if err == nil {
		rpcObj = convertToRPCFmtQsfpPMState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkQsfpPMDataState(fromIdx, count platformd.Int) (*platformd.QsfpPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) restoreQsfpConfigFromDB() (bool, error) {
	var qsfpCfg objects.Qsfp
	qsfpCfgList, err := rpcHdl.dbHdl.GetAllObjFromDb(qsfpCfg)
	if err != nil {
		return false, errors.New("Failed to retrive Qsfp config object from DB")
	}
	for idx := 0; idx < len(qsfpCfgList); idx++ {
		dbObj := qsfpCfgList[idx].(objects.Qsfp)
		obj := new(platformd.Qsfp)
		objects.ConvertplatformdQsfpObjToThrift(&dbObj, obj)
		convNewCfg := convertRPCToObjFmtQsfpConfig(obj)
		ok, err := api.UpdateQsfp(nil, convNewCfg, nil)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (rpcHdl *rpcServiceHandler) CreateQsfpChannel(config *platformd.QsfpChannel) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeleteQsfpChannel(config *platformd.QsfpChannel) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdateQsfpChannel(oldConfig *platformd.QsfpChannel, newConfig *platformd.QsfpChannel, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtQsfpChannelConfig(oldConfig)
	newCfg := convertRPCToObjFmtQsfpChannelConfig(newConfig)
	rv, err := api.UpdateQsfpChannel(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetQsfpChannel(QsfpId int32, ChannelNum int32) (*platformd.QsfpChannel, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkQsfpChannel(fromIdx, count platformd.Int) (*platformd.QsfpChannelGetInfo, error) {
	var getBulkObj platformd.QsfpChannelGetInfo
	var err error

	info, err := api.GetBulkQsfpChannelConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.QsfpChannelList = append(getBulkObj.QsfpChannelList, convertToRPCFmtQsfpChannelConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetQsfpChannelState(QsfpId int32, ChannelNum int32) (*platformd.QsfpChannelState, error) {
	var rpcObj *platformd.QsfpChannelState
	var err error

	obj, err := api.GetQsfpChannelState(QsfpId, ChannelNum)
	if err == nil {
		rpcObj = convertToRPCFmtQsfpChannelState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkQsfpChannelState(fromIdx, count platformd.Int) (*platformd.QsfpChannelStateGetInfo, error) {
	var getBulkObj platformd.QsfpChannelStateGetInfo
	var err error

	info, err := api.GetBulkQsfpChannelState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.QsfpChannelStateList = append(getBulkObj.QsfpChannelStateList, convertToRPCFmtQsfpChannelState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetQsfpChannelPMDataState(QsfpId int32, ChannelNum int32, Resource string, Class string) (*platformd.QsfpChannelPMDataState, error) {
	var rpcObj *platformd.QsfpChannelPMDataState
	var err error

	obj, err := api.GetQsfpChannelPMDataState(QsfpId, ChannelNum, Resource, Class)
	if err == nil {
		rpcObj = convertToRPCFmtQsfpChannelPMState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkQsfpChannelPMDataState(fromIdx, count platformd.Int) (*platformd.QsfpChannelPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) restoreQsfpChannelConfigFromDB() (bool, error) {
	var qsfpCfg objects.QsfpChannel
	qsfpCfgList, err := rpcHdl.dbHdl.GetAllObjFromDb(qsfpCfg)
	if err != nil {
		return false, errors.New("Failed to retrive Qsfp Channel config object from DB")
	}
	for idx := 0; idx < len(qsfpCfgList); idx++ {
		dbObj := qsfpCfgList[idx].(objects.QsfpChannel)
		obj := new(platformd.QsfpChannel)
		objects.ConvertplatformdQsfpChannelObjToThrift(&dbObj, obj)
		convNewCfg := convertRPCToObjFmtQsfpChannelConfig(obj)
		ok, err := api.UpdateQsfpChannel(nil, convNewCfg, nil)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}
