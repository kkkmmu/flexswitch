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

package pluginManager

import (
	"fmt"
	"infra/platformd/objects"
	"infra/platformd/pluginManager/dummy"
	"infra/platformd/pluginManager/onlp"
	"infra/platformd/pluginManager/openBMC"
	"infra/platformd/pluginManager/openBMCVoyager"
	"infra/platformd/pluginManager/pluginCommon"
	"strings"
	"utils/dbutils"
	"utils/logging"
)

type PluginIntf interface {
	Init() error
	DeInit() error
	GetPlatformState() (pluginCommon.PlatformState, error)

	GetSfpState(sfpId int32) (pluginCommon.SfpState, error)
	GetSfpConfig(sfpId int32) (*objects.SfpConfig, error)
	UpdateSfpConfig(cfg *objects.SfpConfig) (bool, error)
	GetAllSfpState(state []pluginCommon.SfpState, count int) error
	GetSfpCnt() int

	GetThermalState(thermalId int32) (pluginCommon.ThermalState, error)
	GetAllThermalState(state []pluginCommon.ThermalState, count int) error
	GetMaxNumOfThermal() int

	GetFanState(fanId int32) (pluginCommon.FanState, error)
	GetFanConfig(fanId int32) (*objects.FanConfig, error)
	UpdateFanConfig(cfg *objects.FanConfig) (bool, error)
	GetMaxNumOfFans() int
	GetAllFanState(state []pluginCommon.FanState, count int) error

	GetPsuState(psuId int32) (pluginCommon.PsuState, error)
	GetAllPsuState(state []pluginCommon.PsuState, count int) error
	GetMaxNumOfPsu() int

	GetLedState(ledId int32) (pluginCommon.LedState, error)
	GetAllLedState(state []pluginCommon.LedState, count int) error
	GetMaxNumOfLed() int

	GetAllSensorState(state *pluginCommon.SensorState) error
	GetMaxNumOfQsfp() int
	GetQsfpState(id int32) (pluginCommon.QsfpState, error)
	GetQsfpPMData(id int32) (pluginCommon.QsfpPMData, error)
	GetPlatformMgmtDeviceState(state *pluginCommon.PlatformMgmtDeviceState) error
}

type ResourceManagers struct {
	*SysManager
	*FanManager
	*PsuManager
	*SfpManager
	*ThermalManager
	*LedManager
	*SensorManager
	*QsfpManager
	*PlatformManager
}

type PluginManager struct {
	*ResourceManagers
	logger     logging.LoggerIntf
	plugin     PluginIntf
	EventDbHdl dbutils.DBIntf
}

func NewPluginMgr(initParams *pluginCommon.PluginInitParams) (*PluginManager, error) {
	var plugin PluginIntf
	var err error
	pluginMgr := new(PluginManager)
	pluginMgr.ResourceManagers = new(ResourceManagers)
	pluginMgr.logger = initParams.Logger
	pluginMgr.EventDbHdl = initParams.EventDbHdl
	pluginName := strings.ToLower(initParams.PluginName)
	switch pluginName {
	case pluginCommon.ONLP_PLUGIN:
		fmt.Println("===== ONLP_PLUGIN =====")
		plugin, err = onlp.NewONLPPlugin(initParams)
		if err != nil {
			return nil, err
		}
		pluginMgr.plugin = plugin
	case pluginCommon.OpenBMC_PLUGIN:
		fmt.Println("===== OPENBMC_PLUGIN =====")
		plugin, err = openBMC.NewOpenBMCPlugin(initParams)
		if err != nil {
			return nil, err
		}
		pluginMgr.plugin = plugin
	case pluginCommon.OpenBMCVoyager_PLUGIN:
		fmt.Println("===== OPENBMCVOYAGER_PLUGIN =====")
		plugin, err = openBMCVoyager.NewOpenBMCVoyagerPlugin(initParams)
		if err != nil {
			return nil, err
		}
		pluginMgr.plugin = plugin
	case pluginCommon.Dummy_PLUGIN:
		fmt.Println("===== Dummy_PLUGIN =====")
		plugin, err = dummy.NewDummyPlugin(initParams)
		if err != nil {
			return nil, err
		}
		pluginMgr.plugin = plugin
	default:
	}
	pluginMgr.SysManager = &SysMgr
	pluginMgr.FanManager = &FanMgr
	pluginMgr.PsuManager = &PsuMgr
	pluginMgr.SfpManager = &SfpMgr
	pluginMgr.ThermalManager = &ThermalMgr
	pluginMgr.LedManager = &LedMgr
	pluginMgr.SensorManager = &SensorMgr
	pluginMgr.QsfpManager = &QsfpMgr
	pluginMgr.PlatformManager = &PlatformMgr
	return pluginMgr, nil
}

func (pMgr *PluginManager) Init() error {
	pMgr.plugin.Init()
	pMgr.SysManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.FanManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.PsuManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.SfpManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.ThermalManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.LedManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.SensorManager.Init(pMgr.logger, pMgr.plugin, pMgr.EventDbHdl)
	pMgr.QsfpManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.PlatformManager.Init(pMgr.logger, pMgr.plugin)
	return nil
}

func (pMgr *PluginManager) Deinit() {
	pMgr.SysManager.Deinit()
	pMgr.FanManager.Deinit()
	pMgr.PsuManager.Deinit()
	pMgr.SfpManager.Deinit()
	pMgr.ThermalManager.Deinit()
	pMgr.LedManager.Deinit()
	pMgr.SensorManager.Deinit()
	pMgr.QsfpManager.Deinit()
	pMgr.PlatformManager.Deinit()
}
