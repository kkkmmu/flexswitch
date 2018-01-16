//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package server

import (
	"infra/fMgrd/objects"
)

type ServerOpId int

const (
	GET_BULK_FAULT_STATE ServerOpId = iota
	GET_BULK_ALARM_STATE
	FAULT_ENABLE_ACTION
	FAULT_CLEAR_ACTION
)

type ServerRequest struct {
	Op   ServerOpId
	Data interface{}
}

type GetBulkInArgs struct {
	FromIdx int
	Count   int
}

type GetBulkFaultStateOutArgs struct {
	BulkInfo *objects.FaultStateGetInfo
	Err      error
}

type GetBulkAlarmStateOutArgs struct {
	BulkInfo *objects.AlarmStateGetInfo
	Err      error
}

type FaultEnableActionInArgs struct {
	Config *objects.FaultEnable
}

type FaultEnableActionOutArgs struct {
	RetVal bool
	Err    error
}

type FaultClearActionInArgs struct {
	Config *objects.FaultClear
}

type FaultClearActionOutArgs struct {
	RetVal bool
	Err    error
}
