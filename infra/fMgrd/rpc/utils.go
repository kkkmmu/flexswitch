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
	"infra/fMgrd/objects"
)

func convertToRPCFmtFaultState(obj objects.FaultState) *fMgrd.FaultState {
	return &fMgrd.FaultState{
		OwnerId:          obj.OwnerId,
		EventId:          obj.EventId,
		OwnerName:        obj.OwnerName,
		EventName:        obj.EventName,
		SrcObjName:       obj.SrcObjName,
		Description:      obj.Description,
		OccuranceTime:    obj.OccuranceTime,
		SrcObjKey:        obj.SrcObjKey,
		SrcObjUUID:       obj.SrcObjUUID,
		ResolutionTime:   obj.ResolutionTime,
		ResolutionReason: obj.ResolutionReason,
	}
}

func convertToRPCFmtAlarmState(obj objects.AlarmState) *fMgrd.AlarmState {
	return &fMgrd.AlarmState{
		OwnerId:          obj.OwnerId,
		EventId:          obj.EventId,
		OwnerName:        obj.OwnerName,
		EventName:        obj.EventName,
		SrcObjName:       obj.SrcObjName,
		Description:      obj.Description,
		OccuranceTime:    obj.OccuranceTime,
		SrcObjKey:        obj.SrcObjKey,
		SrcObjUUID:       obj.SrcObjUUID,
		ResolutionTime:   obj.ResolutionTime,
		Severity:         obj.Severity,
		ResolutionReason: obj.ResolutionReason,
	}
}

func convertToObjFmtFaultEnable(config *fMgrd.FaultEnable) *objects.FaultEnable {
	return &objects.FaultEnable{
		OwnerName: config.OwnerName,
		EventName: config.EventName,
		Enable:    config.Enable,
	}
}

func convertToObjFmtFaultClear(config *fMgrd.FaultClear) *objects.FaultClear {
	return &objects.FaultClear{
		OwnerName:  config.OwnerName,
		EventName:  config.EventName,
		SrcObjUUID: config.SrcObjUUID,
	}
}
