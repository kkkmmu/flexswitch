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

func (svr *PlatformdServer) getQsfpState(QsfpId int32) (*objects.QsfpState, error) {
	retObj, err := svr.pluginMgr.GetQsfpState(QsfpId)
	return retObj, err
}

func (svr *PlatformdServer) getBulkQsfpState(fromIdx int, count int) (*objects.QsfpStateGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkQsfpState(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) getBulkQsfpConfig(fromIdx int, count int) (*objects.QsfpConfigGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkQsfpConfig(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) updateQsfpConfig(oldCfg *objects.QsfpConfig, newCfg *objects.QsfpConfig, attrset []bool) (bool, error) {
	ret, err := svr.pluginMgr.UpdateQsfpConfig(oldCfg, newCfg, attrset)
	return ret, err
}

func (svr *PlatformdServer) getQsfpPMState(QsfpId int32, Resource string, Class string) (*objects.QsfpPMState, error) {
	retObj, err := svr.pluginMgr.GetQsfpPMState(QsfpId, Resource, Class)
	return retObj, err
}

func (svr *PlatformdServer) getQsfpChannelState(QsfpId int32, ChannelNum int32) (*objects.QsfpChannelState, error) {
	retObj, err := svr.pluginMgr.GetQsfpChannelState(QsfpId, ChannelNum)
	return retObj, err
}

func (svr *PlatformdServer) getBulkQsfpChannelState(fromIdx int, count int) (*objects.QsfpChannelStateGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkQsfpChannelState(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) getBulkQsfpChannelConfig(fromIdx int, count int) (*objects.QsfpChannelConfigGetInfo, error) {
	retObj, err := svr.pluginMgr.GetBulkQsfpChannelConfig(fromIdx, count)
	return retObj, err
}

func (svr *PlatformdServer) updateQsfpChannelConfig(oldCfg *objects.QsfpChannelConfig, newCfg *objects.QsfpChannelConfig, attrset []bool) (bool, error) {
	ret, err := svr.pluginMgr.UpdateQsfpChannelConfig(oldCfg, newCfg, attrset)
	return ret, err
}

func (svr *PlatformdServer) getQsfpChannelPMState(QsfpId int32, ChannelNum int32, Resource string, Class string) (*objects.QsfpChannelPMState, error) {
	retObj, err := svr.pluginMgr.GetQsfpChannelPMState(QsfpId, ChannelNum, Resource, Class)
	return retObj, err
}
