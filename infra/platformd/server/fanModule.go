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

func (svr *PlatformdServer) getFanState(fanId int32) (*objects.FanState, error) {
	retObj, err := svr.pluginMgr.GetFanState(fanId)
	return retObj, err
}

func (svr *PlatformdServer) getBulkFanState(fromIdx int, count int) (*objects.FanStateGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkFanState(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) getFanConfig(fanId int32) (*objects.FanConfig, error) {
	retObj, err := svr.pluginMgr.GetFanConfig(fanId)
	return retObj, err
}

func (svr *PlatformdServer) getBulkFanConfig(fromIdx int, count int) (*objects.FanConfigGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkFanConfig(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) updateFanConfig(oldCfg *objects.FanConfig, newCfg *objects.FanConfig, attrset []bool) (bool, error) {
	ret, err := svr.pluginMgr.UpdateFanConfig(oldCfg, newCfg, attrset)
	return ret, err
}
