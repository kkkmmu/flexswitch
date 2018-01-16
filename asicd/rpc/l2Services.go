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

// This file defines all interfaces provided for L2 service
package rpc

import (
	"asicdInt"
	"asicdServices"
)

/* Method to create a lag
 * hashType  - Hash algorithm to use for this lag
 * ports     - List of ports to be added as members of this lag
 */
func (svcHdlr AsicDaemonServiceHandler) CreateLag(ifName string, hashType int32, ports string) (rval int32, err error) {
	return rval, nil
}

/* Method to delete a lag
 * ifIndex - ID of lag to be deleted
 */
func (svcHdlr AsicDaemonServiceHandler) DeleteLag(ifIndex int32) (rval int32, err error) {
	return rval, nil
}

/* Method to update a lag
 * ifIndex   - ID of lag to be deleted
 * hashType  - Hash algorithm to use for this lag
 * ports     - List of new ports to be added as members of this lag
 */
func (svcHdlr AsicDaemonServiceHandler) UpdateLag(ifIndex, hashType int32, ports string) (rval int32, err error) {
	return rval, nil
}

/* Method to retrieve all LAG objects */
func (svcHdlr AsicDaemonServiceHandler) GetBulkLag(currMarker, count asicdInt.Int) (*asicdInt.LagGetInfo, error) {
	bulkObj := asicdInt.NewLagGetInfo()
	return bulkObj, nil
}

/* Method to create an stg
 * vlanList - List of vlan id's that are members of this stg
 */
func (svcHdlr AsicDaemonServiceHandler) CreateStg(vlanList []int32) (stgId int32, err error) {
	return stgId, err
}

/* Method to delete an stg
 * stgId - Id of spanning tree group to be deleted
 */
func (svcHdlr AsicDaemonServiceHandler) DeleteStg(stgId int32) (rv bool, err error) {
	return rv, err
}

/* Method to set stp state of a port
 * stgId - Id of spanning tree group
 * port  - Id of port memeber in the stg
 * stpState - Spanning tree state of member port to set
 */
func (svcHdlr AsicDaemonServiceHandler) SetPortStpState(stgId, port, stpState int32) (rv bool, err error) {
	return rv, err
}

/* Method to get stp state of a port
 * stgId - Id of spanning tree group
 * port  - Id of port memeber in the stg
 */
func (svcHdlr AsicDaemonServiceHandler) GetPortStpState(stgId, port int32) (stpState int32, err error) {
	return stpState, err
}

/* Method to update an stg
 * stgId    - Id of stg to update
 * VlanList - List of vlan id's that are to be added as members of this stg
 */
func (svcHdlr AsicDaemonServiceHandler) UpdateStgVlanList(stgId int32, vlanList []int32) (rv bool, err error) {
	return rv, err
}

/* Method to flush FDB table per vlan
 * stgId - Id of stg for flush operation
 */
func (svcHdlr AsicDaemonServiceHandler) FlushFdbStgGroup(stgId, port int32) error {
	return nil
}

/* Method to retrieve MAC table information for specific mac addr */
func (svcHdlr AsicDaemonServiceHandler) GetMacTableEntryState(macAddr string) (*asicdServices.MacTableEntryState, error) {
	return nil, nil
}

/* Method to retrieve all MAC table objects */
func (svcHdlr AsicDaemonServiceHandler) GetBulkMacTableEntryState(currMarker, count asicdServices.Int) (*asicdServices.MacTableEntryStateGetInfo, error) {
	return nil, nil
}

/* Method to enable packet reception for specific protocol mac */
func (svcHdlr AsicDaemonServiceHandler) EnablePacketReception(macObj *asicdInt.RsvdProtocolMacConfig) (rv bool, err error) {
	return rv, err
}

/* Method to disable packet reception for specific protocol mac */
func (svcHdlr AsicDaemonServiceHandler) DisablePacketReception(macObj *asicdInt.RsvdProtocolMacConfig) (rv bool, err error) {
	return rv, err
}
