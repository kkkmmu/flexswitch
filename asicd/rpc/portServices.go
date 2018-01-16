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

// This file defines all interfaces provided for the LAG service
package rpc

import (
	"asicdServices"
	"errors"
)

func (svcHdlr AsicDaemonServiceHandler) CreatePort(portObj *asicdServices.Port) (bool, error) {
	return false, errors.New("Create operation not supported on PortConfig object")
}

func (svcHdlr AsicDaemonServiceHandler) DeletePort(portObj *asicdServices.Port) (bool, error) {
	return false, errors.New("Delete operation not supported on PortConfig object")
}

/* Method to update attributes on a port object
 * attrset - list of boolean values indicating what object attributes are modified
 */
func (svcHdlr AsicDaemonServiceHandler) UpdatePort(oldPortObj, newPortObj *asicdServices.Port, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}

/* Method to retrieve all port config objects */
func (svcHdlr AsicDaemonServiceHandler) GetBulkPort(currMarker, count asicdServices.Int) (*asicdServices.PortGetInfo, error) {
	bulkObj := asicdServices.NewPortGetInfo()
	return bulkObj, nil
}

/* Method to retrieve a specific port config object */
func (svcHdlr AsicDaemonServiceHandler) GetPort(intfRef string) (*asicdServices.Port, error) {
	return nil, nil
}

/* Method to retrieve all port state objects */
func (svcHdlr AsicDaemonServiceHandler) GetBulkPortState(currMarker, count asicdServices.Int) (*asicdServices.PortStateGetInfo, error) {
	bulkObj := asicdServices.NewPortStateGetInfo()
	return bulkObj, nil
}

/* Method to retrieve a specific port state object */
func (svcHdlr AsicDaemonServiceHandler) GetPortState(intfRef string) (*asicdServices.PortState, error) {
	return nil, nil
}

/* Method to error disable a port */
func (svcHdlr AsicDaemonServiceHandler) ErrorDisablePort(ifIndex int32, adminState string, errDisableReason string) (bool, error) {
	return true, nil
}
func (svcHdlr AsicDaemonServiceHandler) GetBulkBufferPortStatState(currMarker, count asicdServices.Int) (*asicdServices.BufferPortStatStateGetInfo, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBufferPortStatState(intfRef string) (*asicdServices.BufferPortStatState, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkBufferGlobalStatState(currMarker, count asicdServices.Int) (*asicdServices.BufferGlobalStatStateGetInfo, error) {
	return nil, nil

}

func (svcHdlr AsicDaemonServiceHandler) GetBufferGlobalStatState(deviceId int32) (*asicdServices.BufferGlobalStatState, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateEthernetPM(obj *asicdServices.EthernetPM) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteEthernetPM(obj *asicdServices.EthernetPM) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateEthernetPM(oldObj, newObj *asicdServices.EthernetPM, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetEthernetPM(intfRef, resource string) (*asicdServices.EthernetPM, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkEthernetPM(currMarker, count asicdServices.Int) (*asicdServices.EthernetPMGetInfo, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetEthernetPMState(intfRef, resource string) (*asicdServices.EthernetPMState, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkEthernetPMState(currMarker, count asicdServices.Int) (*asicdServices.EthernetPMStateGetInfo, error) {
	return nil, nil
}
