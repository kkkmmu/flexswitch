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

package server

import (
	"infra/platformd/objects"
)

func (svr *PlatformdServer) getSfpState(sfpId int32) (*objects.SfpState, error) {
	retObj, err := svr.pluginMgr.GetSfpState(sfpId)
	return retObj, err
}

func (svr *PlatformdServer) getBulkSfpState(fromIdx int, count int) (*objects.SfpStateGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkSfpState(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) getSfpConfig(sfpId int32) (*objects.SfpConfig, error) {
	retObj, err := svr.pluginMgr.GetSfpConfig(sfpId)
	return retObj, err
}

func (svr *PlatformdServer) getBulkSfpConfig(fromIdx int, count int) (*objects.SfpConfigGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkSfpConfig(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) updateSfpConfig(oldCfg *objects.SfpConfig, newCfg *objects.SfpConfig, attrset []bool) (bool, error) {
	ret, err := svr.pluginMgr.UpdateSfpConfig(oldCfg, newCfg, attrset)
	return ret, err
}
