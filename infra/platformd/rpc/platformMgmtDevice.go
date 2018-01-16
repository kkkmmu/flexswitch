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
	"infra/platformd/api"
	"platformd"
)

func (rpcHdl *rpcServiceHandler) GetPlatformMgmtDeviceState(DeviceName string) (*platformd.PlatformMgmtDeviceState, error) {
	var rpcObj *platformd.PlatformMgmtDeviceState
	var err error

	obj, err := api.GetPlatformMgmtDeviceState(DeviceName)
	if err == nil {
		rpcObj = convertToRPCFmtPlatformMgmtDeviceState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkPlatformMgmtDeviceState(fromIdx, count platformd.Int) (*platformd.PlatformMgmtDeviceStateGetInfo, error) {
	var getBulkObj platformd.PlatformMgmtDeviceStateGetInfo
	var err error

	info, err := api.GetBulkPlatformMgmtDeviceState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.PlatformMgmtDeviceStateList = append(getBulkObj.PlatformMgmtDeviceStateList, convertToRPCFmtPlatformMgmtDeviceState(info.List[idx]))
	}
	return &getBulkObj, err
}
