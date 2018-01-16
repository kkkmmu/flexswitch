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

func (svr *PlatformdServer) getVoltageSensorState(Name string) (*objects.VoltageSensorState, error) {
	retObj, err := svr.pluginMgr.GetVoltageSensorState(Name)
	return retObj, err
}

func (svr *PlatformdServer) getBulkVoltageSensorState(fromIdx int, count int) (*objects.VoltageSensorStateGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkVoltageSensorState(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) getBulkVoltageSensorConfig(fromIdx int, count int) (*objects.VoltageSensorConfigGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkVoltageSensorConfig(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) updateVoltageSensorConfig(oldCfg *objects.VoltageSensorConfig, newCfg *objects.VoltageSensorConfig, attrset []bool) (bool, error) {
	ret, err := svr.pluginMgr.UpdateVoltageSensorConfig(oldCfg, newCfg, attrset)
	return ret, err
}

func (svr *PlatformdServer) getVoltageSensorPMState(Name string, Class string) (*objects.VoltageSensorPMState, error) {
	retObj, err := svr.pluginMgr.GetVoltageSensorPMState(Name, Class)
	return retObj, err
}
