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

package api

import (
	"errors"
	"infra/fMgrd/objects"
	"infra/fMgrd/server"
)

var svr *server.FMGRServer

func InitApiLayer(server *server.FMGRServer) {
	svr = server
	svr.Logger.Info("Initializing API Layer")
}

func GetBulkFault(fromIdx, count int) (*objects.FaultStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_FAULT_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkFaultStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response recevied from server during GetBulkFaultState")
	}
}

func GetBulkAlarm(fromIdx, count int) (*objects.AlarmStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_ALARM_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkAlarmStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response recevied from server during GetBulkFaultState")
	}
}

func FaultEnableAction(cfg *objects.FaultEnable) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.FAULT_ENABLE_ACTION,
		Data: interface{}(&server.FaultEnableActionInArgs{
			Config: cfg,
		}),
	}

	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.FaultEnableActionOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response recevied from server during Executing Fault Enable Action")

}

func FaultClearAction(cfg *objects.FaultClear) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.FAULT_CLEAR_ACTION,
		Data: interface{}(&server.FaultClearActionInArgs{
			Config: cfg,
		}),
	}

	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.FaultClearActionOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response recevied from server during Executing Fault Clear Action")
}
