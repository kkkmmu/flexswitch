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

func (svr *PlatformdServer) getPowerConverterSensorState(Name string) (*objects.PowerConverterSensorState, error) {
	retObj, err := svr.pluginMgr.GetPowerConverterSensorState(Name)
	return retObj, err
}

func (svr *PlatformdServer) getBulkPowerConverterSensorState(fromIdx int, count int) (*objects.PowerConverterSensorStateGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkPowerConverterSensorState(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) getBulkPowerConverterSensorConfig(fromIdx int, count int) (*objects.PowerConverterSensorConfigGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkPowerConverterSensorConfig(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) updatePowerConverterSensorConfig(oldCfg *objects.PowerConverterSensorConfig, newCfg *objects.PowerConverterSensorConfig, attrset []bool) (bool, error) {
	ret, err := svr.pluginMgr.UpdatePowerConverterSensorConfig(oldCfg, newCfg, attrset)
	return ret, err
}

func (svr *PlatformdServer) getPowerConverterSensorPMState(Name string, Class string) (*objects.PowerConverterSensorPMState, error) {
	retObj, err := svr.pluginMgr.GetPowerConverterSensorPMState(Name, Class)
	return retObj, err
}
