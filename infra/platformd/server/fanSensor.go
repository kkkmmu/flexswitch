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

func (svr *PlatformdServer) getFanSensorState(Name string) (*objects.FanSensorState, error) {
	retObj, err := svr.pluginMgr.GetFanSensorState(Name)
	return retObj, err
}

func (svr *PlatformdServer) getBulkFanSensorState(fromIdx int, count int) (*objects.FanSensorStateGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkFanSensorState(fromIdx, count)
	svr.Logger.Info("Fan State:", retObj)
	return retObj, err
}

func (svr *PlatformdServer) getBulkFanSensorConfig(fromIdx int, count int) (*objects.FanSensorConfigGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkFanSensorConfig(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) updateFanSensorConfig(oldCfg *objects.FanSensorConfig, newCfg *objects.FanSensorConfig, attrset []bool) (bool, error) {
	ret, err := svr.pluginMgr.UpdateFanSensorConfig(oldCfg, newCfg, attrset)
	return ret, err
}

func (svr *PlatformdServer) getFanSensorPMState(Name string, Class string) (*objects.FanSensorPMState, error) {
	retObj, err := svr.pluginMgr.GetFanSensorPMState(Name, Class)
	return retObj, err
}
