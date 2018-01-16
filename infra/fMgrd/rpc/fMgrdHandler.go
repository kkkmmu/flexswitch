//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
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
	"fMgrd"
	"fmt"
	"infra/fMgrd/api"
	//"utils/logging"
)

func (h *rpcServiceHandler) CreateFMgrGlobal(conf *fMgrd.FMgrGlobal) (bool, error) {
	h.logger.Info(fmt.Sprintln("Received CreateFMgrGlobal call"))
	return true, nil
}

func (h *rpcServiceHandler) DeleteFMgrGlobal(conf *fMgrd.FMgrGlobal) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete FMgr config attrs:", conf))
	return true, nil
}

func (h *rpcServiceHandler) UpdateFMgrGlobal(origConf *fMgrd.FMgrGlobal, newConf *fMgrd.FMgrGlobal, attrset []bool, op []*fMgrd.PatchOpInfo) (bool, error) {
	return true, nil
}

func (h *rpcServiceHandler) GetBulkFaultState(fromIndex fMgrd.Int, count fMgrd.Int) (*fMgrd.FaultStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("Get bulk call for Faults"))
	var getBulkObj fMgrd.FaultStateGetInfo
	info, err := api.GetBulkFault(int(fromIndex), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fMgrd.Int(fromIndex)
	getBulkObj.EndIdx = fMgrd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = fMgrd.Int(info.Count)
	for idx := 0; idx < info.Count; idx++ {
		getBulkObj.FaultStateList = append(getBulkObj.FaultStateList, convertToRPCFmtFaultState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (h *rpcServiceHandler) GetFaultState(ownerId int32, eventId int32, ownerName string, eventName string, srcObjName string) (*fMgrd.FaultState, error) {
	return nil, nil
}

func (h *rpcServiceHandler) GetBulkAlarmState(fromIndex fMgrd.Int, count fMgrd.Int) (*fMgrd.AlarmStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("Get bulk call for Alarm"))
	var getBulkObj fMgrd.AlarmStateGetInfo
	info, err := api.GetBulkAlarm(int(fromIndex), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fMgrd.Int(fromIndex)
	getBulkObj.EndIdx = fMgrd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = fMgrd.Int(info.Count)
	for idx := 0; idx < info.Count; idx++ {
		getBulkObj.AlarmStateList = append(getBulkObj.AlarmStateList, convertToRPCFmtAlarmState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (h *rpcServiceHandler) GetAlarmState(ownerId int32, eventId int32, ownerName string, eventName string, srcObjName string) (*fMgrd.AlarmState, error) {
	return nil, nil
}

func (h *rpcServiceHandler) ExecuteActionFaultEnable(config *fMgrd.FaultEnable) (bool, error) {
	h.logger.Info(fmt.Sprintln("ExecuteActionFaultEnable ", config))

	return api.FaultEnableAction(convertToObjFmtFaultEnable(config))
}

func (h *rpcServiceHandler) ExecuteActionFaultClear(config *fMgrd.FaultClear) (bool, error) {
	h.logger.Info(fmt.Sprintln("ExecuteActionFaultClear ", config))

	return api.FaultClearAction(convertToObjFmtFaultClear(config))
}
