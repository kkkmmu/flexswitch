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

// This file defines all interfaces provided for the AsicGlobal service

package rpc

import (
	"asicdServices"
)

func (svcHdlr AsicDaemonServiceHandler) GetAsicGlobalState(moduleId int8) (*asicdServices.AsicGlobalState, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkAsicGlobalState(currMarker, count asicdServices.Int) (*asicdServices.AsicGlobalStateGetInfo, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateAsicGlobalPM(obj *asicdServices.AsicGlobalPM) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteAsicGlobalPM(obj *asicdServices.AsicGlobalPM) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateAsicGlobalPM(oldObj, newObj *asicdServices.AsicGlobalPM, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetAsicGlobalPM(moduleId int8, resource string) (*asicdServices.AsicGlobalPM, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkAsicGlobalPM(currMarker, count asicdServices.Int) (*asicdServices.AsicGlobalPMGetInfo, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetAsicGlobalPMState(moduleId int8, resource string) (*asicdServices.AsicGlobalPMState, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkAsicGlobalPMState(currMarker, count asicdServices.Int) (*asicdServices.AsicGlobalPMStateGetInfo, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetAsicSummaryState(moduleId int8) (*asicdServices.AsicSummaryState, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkAsicSummaryState(currMarker, count asicdServices.Int) (*asicdServices.AsicSummaryStateGetInfo, error) {
	return nil, nil
}
