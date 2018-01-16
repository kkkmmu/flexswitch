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

package openBMCVoyager

import (
	"errors"
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"time"
	"utils/logging"
)

const (
	SENSOR_POLLING_TIME = time.Duration(1) * time.Second
)

type openBMCVoyagerDriver struct {
	logger logging.LoggerIntf
	ipAddr string
	port   string
}

var driver openBMCVoyagerDriver

func NewOpenBMCVoyagerPlugin(params *pluginCommon.PluginInitParams) (*openBMCVoyagerDriver, error) {
	var err error
	driver.logger = params.Logger
	driver.ipAddr = params.IpAddr
	driver.port = params.Port
	return &driver, err
}

func (driver *openBMCVoyagerDriver) Init() error {
	driver.logger.Info("Initializing openBMC driver")
	return nil
}

func (driver *openBMCVoyagerDriver) DeInit() error {
	driver.logger.Info("DeInitializing openBMC voyager driver")
	return nil
}

func (driver *openBMCVoyagerDriver) GetFanState(fanId int32) (state pluginCommon.FanState, err error) {
	return state, errors.New("Not Supported")
}

func (driver *openBMCVoyagerDriver) GetFanConfig(fanId int32) (retObj *objects.FanConfig, err error) {
	return retObj, errors.New("Not Supported")
}

func (driver *openBMCVoyagerDriver) UpdateFanConfig(cfg *objects.FanConfig) (bool, error) {
	return false, errors.New("Not Supported")
}

func (driver *openBMCVoyagerDriver) GetMaxNumOfFans() int {
	return 0
}

func (driver *openBMCVoyagerDriver) GetAllFanState(states []pluginCommon.FanState, cnt int) error {
	return errors.New("Not Supported")
}

func (driver *openBMCVoyagerDriver) GetSfpState(sfpId int32) (obj pluginCommon.SfpState, err error) {
	return obj, errors.New("Not Supported")
}

func (driver *openBMCVoyagerDriver) GetSfpConfig(sfpId int32) (obj *objects.SfpConfig, err error) {
	return obj, errors.New("Not Supported")
}

func (driver *openBMCVoyagerDriver) UpdateSfpConfig(cfg *objects.SfpConfig) (bool, error) {
	return false, errors.New("Not Supported")
}

func (driver *openBMCVoyagerDriver) GetAllSfpState(states []pluginCommon.SfpState, cnt int) error {
	return errors.New("Not Supported")
}

func (driver *openBMCVoyagerDriver) GetSfpCnt() int {
	return 0
}

func (driver *openBMCVoyagerDriver) GetPlatformState() (retObj pluginCommon.PlatformState, err error) {
	return retObj, errors.New("Not Supported")
}

func (driver *openBMCVoyagerDriver) GetMaxNumOfThermal() int {
	return 0
}

func (driver *openBMCVoyagerDriver) GetThermalState(thermalId int32) (state pluginCommon.ThermalState, err error) {
	return state, errors.New("Not Supoorted")
}

func (driver *openBMCVoyagerDriver) GetAllThermalState(states []pluginCommon.ThermalState, cnt int) error {
	return nil
}

func (driver *openBMCVoyagerDriver) GetMaxNumOfPsu() int {
	return 0
}

func (driver *openBMCVoyagerDriver) GetPsuState(psuId int32) (pluginCommon.PsuState, error) {
	var retObj pluginCommon.PsuState
	return retObj, nil
}

func (driver *openBMCVoyagerDriver) GetAllPsuState(states []pluginCommon.PsuState, cnt int) error {
	return nil
}

func (driver *openBMCVoyagerDriver) GetMaxNumOfLed() int {
	return 0
}

func (driver *openBMCVoyagerDriver) GetLedState(ledId int32) (pluginCommon.LedState, error) {
	var retObj pluginCommon.LedState
	return retObj, nil
}

func (driver *openBMCVoyagerDriver) GetAllLedState(states []pluginCommon.LedState, cnt int) error {
	return nil
}
