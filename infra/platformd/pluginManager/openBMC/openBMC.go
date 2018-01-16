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

package openBMC

import (
	"errors"
	"fmt"
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"sync"
	"time"
	"utils/logging"
)

const (
	NUM_OF_FAN          int = 10
	NUM_OF_THERMAL      int = 7
	SENSOR_POLLING_TIME     = time.Duration(1) * time.Second
)

/*
 * FanId Mapping
 * Fan1Front = FanId(0)
 * Fan1Rear  = FanId(1)
 * Fan2Front = FanId(2)
 * Fan2Rear  = FanId(3)
 * Fan3Front = FanId(4)
 * Fan3Rear  = FanId(5)
 * Fan4Front = FanId(6)
 * Fan4Rear  = FanId(7)
 * Fan5Front = FanId(8)
 * Fan5Rear  = FanId(9)
 */

/*
 * ThermalId Mapping
 * Switch 		= ThermalId(0)
 * Inlet Left 		= ThermalId(1)
 * Inlet Middle 	= ThermalId(2)
 * Inlet Right 		= ThermalId(3)
 * Outlet Left 		= ThermalId(4)
 * Outlet Middle 	= ThermalId(5)
 * Outlet Right 	= ThermalId(6)
 */

type openBMCDriver struct {
	logger      logging.LoggerIntf
	ipAddr      string
	port        string
	sensorMutex sync.RWMutex
	sensorData  SensorData
	mbFruidInfo MBFruidInfo
	bmcInfo     BMCInfo
}

var driver openBMCDriver

func NewOpenBMCPlugin(params *pluginCommon.PluginInitParams) (*openBMCDriver, error) {
	var err error
	driver.logger = params.Logger
	driver.ipAddr = params.IpAddr
	driver.port = params.Port
	return &driver, err
}

func (driver *openBMCDriver) processSensorData() (err error) {
	driver.sensorMutex.Lock()
	driver.sensorData, err = driver.GetSensorState()
	if err != nil {
		driver.logger.Err(fmt.Sprintln("Error getting OpenBMC Senssor Data", err))
		driver.sensorMutex.Unlock()
		return err
	}
	driver.sensorMutex.Unlock()
	return err
}

func (driver *openBMCDriver) processMBFruidInfo() (err error) {
	driver.mbFruidInfo, err = driver.GetMBFruidInfo()
	if err != nil {
		driver.logger.Err(fmt.Sprintln("Error getting OpenBMC MB Fruid Info", err))
		return err
	}
	return err
}

func (driver *openBMCDriver) processBMCInfo() (err error) {
	driver.bmcInfo, err = driver.GetBMCInfo()
	if err != nil {
		driver.logger.Err(fmt.Sprintln("Error getting OpenBMC BMC Info", err))
		return err
	}
	return err
}

func (driver *openBMCDriver) Init() error {
	driver.logger.Info("Initializing openBMC driver")
	err := driver.processSensorData()
	if err != nil {
		return err
	}
	err = driver.processMBFruidInfo()
	if err != nil {
		return err
	}
	err = driver.processBMCInfo()
	if err != nil {
		return err
	}
	go driver.collectSensorData()
	return err
}

func (driver *openBMCDriver) collectSensorData() {
	var err error
	for {
		time.Sleep(SENSOR_POLLING_TIME)
		driver.sensorMutex.Lock()
		driver.sensorData, err = driver.GetSensorState()
		if err != nil {
			driver.logger.Err(fmt.Sprintln("Error getting OpenBMC Senssor Data", err))
		}
		driver.sensorMutex.Unlock()
	}
}

func (driver *openBMCDriver) DeInit() error {
	driver.logger.Info("DeInitializing openBMC driver")
	return nil
}

func (driver *openBMCDriver) GetFanState(fanId int32) (pluginCommon.FanState, error) {
	var state pluginCommon.FanState
	state.Valid = true
	state.FanId = fanId
	//driver.logger.Info(fmt.Sprintln("Sensor Data:", sensorData))
	driver.sensorMutex.Lock()
	switch fanId {
	case 0:
		state.OperSpeed = convertFanSpeedStringToInt32(driver.sensorData.FanSensor.Fan1Front)
	case 1:
		state.OperSpeed = convertFanSpeedStringToInt32(driver.sensorData.FanSensor.Fan1Rear)
	case 2:
		state.OperSpeed = convertFanSpeedStringToInt32(driver.sensorData.FanSensor.Fan2Front)
	case 3:
		state.OperSpeed = convertFanSpeedStringToInt32(driver.sensorData.FanSensor.Fan2Rear)
	case 4:
		state.OperSpeed = convertFanSpeedStringToInt32(driver.sensorData.FanSensor.Fan3Front)
	case 5:
		state.OperSpeed = convertFanSpeedStringToInt32(driver.sensorData.FanSensor.Fan3Rear)
	case 6:
		state.OperSpeed = convertFanSpeedStringToInt32(driver.sensorData.FanSensor.Fan4Front)
	case 7:
		state.OperSpeed = convertFanSpeedStringToInt32(driver.sensorData.FanSensor.Fan4Rear)
	case 8:
		state.OperSpeed = convertFanSpeedStringToInt32(driver.sensorData.FanSensor.Fan5Front)
	case 9:
		state.OperSpeed = convertFanSpeedStringToInt32(driver.sensorData.FanSensor.Fan5Rear)
	}
	driver.sensorMutex.Unlock()
	if state.OperSpeed == 0 {
		state.OperMode = pluginCommon.FAN_MODE_OFF_STR
	} else {
		state.OperMode = pluginCommon.FAN_MODE_ON_STR
	}
	state.OperDirection = "Not Supported"
	state.Status = "Not Supported"
	state.Model = "Not Supported"
	state.SerialNum = "Not Supported"
	state.LedId = -1
	return state, nil
}

func (driver *openBMCDriver) GetFanConfig(fanId int32) (retObj *objects.FanConfig, err error) {
	return retObj, nil
}

func (driver *openBMCDriver) UpdateFanConfig(cfg *objects.FanConfig) (bool, error) {
	driver.logger.Info("Updating OpenBMC Fan Config")
	return true, nil
}

func (driver *openBMCDriver) GetMaxNumOfFans() int {
	driver.logger.Info("Inside OpenBMC: GetMaxNumOfFans()")
	return NUM_OF_FAN
}

func (driver *openBMCDriver) GetAllFanState(states []pluginCommon.FanState, cnt int) error {
	for idx := 0; idx < cnt; idx++ {
		states[idx], _ = driver.GetFanState(int32(idx))
	}
	return nil
}

func (driver *openBMCDriver) GetSfpState(sfpId int32) (pluginCommon.SfpState, error) {
	var retObj pluginCommon.SfpState

	// TODO
	retObj.SfpId = sfpId
	return retObj, nil
}

func (driver *openBMCDriver) GetSfpConfig(sfpId int32) (*objects.SfpConfig, error) {
	var retObj objects.SfpConfig

	// TODO
	retObj.SfpId = sfpId
	return &retObj, nil
}

func (driver *openBMCDriver) UpdateSfpConfig(cfg *objects.SfpConfig) (bool, error) {
	driver.logger.Info("Updating Onlp SFP Config")
	return true, nil
}

func (driver *openBMCDriver) GetAllSfpState(states []pluginCommon.SfpState, cnt int) error {
	driver.logger.Info("GetAllSfpState")
	return nil
}

func (driver *openBMCDriver) GetSfpCnt() int {
	driver.logger.Info("GetSfpCnt")
	return 0
}

func (driver *openBMCDriver) GetPlatformState() (pluginCommon.PlatformState, error) {
	var retObj pluginCommon.PlatformState
	retObj.ProductName = driver.mbFruidInfo.ProductName
	retObj.SerialNum = driver.mbFruidInfo.ProSerialNum
	retObj.Manufacturer = driver.mbFruidInfo.SystemManufacturer
	retObj.Vendor = driver.mbFruidInfo.AssemAt
	retObj.Release = fmt.Sprintf("%d.%d", driver.mbFruidInfo.ProductVer, driver.mbFruidInfo.ProSubVer)
	retObj.PlatformName = driver.bmcInfo.Description
	retObj.Version = driver.bmcInfo.OpenBMCVersion
	return retObj, nil
}

func (driver *openBMCDriver) GetMaxNumOfThermal() int {
	driver.logger.Info("Inside OpenBMC: GetMaxNumOfThermal()")
	return NUM_OF_THERMAL
}

func (driver *openBMCDriver) GetThermalState(thermalId int32) (pluginCommon.ThermalState, error) {
	var state pluginCommon.ThermalState
	state.Valid = true
	state.ThermalId = thermalId
	//driver.logger.Info(fmt.Sprintln("Sensor Data:", sensorData))
	driver.sensorMutex.Lock()
	switch thermalId {
	case 0:
		state.Location = "Switch"
		state.Temperature = driver.sensorData.TempSensor.SwitchTemp.Temp
	case 1:
		state.Location = "Inlet Left"
		state.Temperature = driver.sensorData.TempSensor.InLeftTemp.Temp
	case 2:
		state.Location = "Inlet Middle"
		state.Temperature = driver.sensorData.TempSensor.InMidTemp.Temp
	case 3:
		state.Location = "Inlet Right"
		state.Temperature = driver.sensorData.TempSensor.InRightTemp.Temp
	case 4:
		state.Location = "Outlet Left"
		state.Temperature = driver.sensorData.TempSensor.OutLeftTemp.Temp
	case 5:
		state.Location = "Outlet Middle"
		state.Temperature = driver.sensorData.TempSensor.OutMidTemp.Temp
	case 6:
		state.Location = "Outlet Right"
		state.Temperature = driver.sensorData.TempSensor.OutRightTemp.Temp
	}
	driver.sensorMutex.Unlock()
	state.LowerWatermarkTemperature = "Not Supported"
	state.UpperWatermarkTemperature = "Not Supported"
	state.ShutdownTemperature = "Not Supported"
	return state, nil
}

func (driver *openBMCDriver) GetAllThermalState(states []pluginCommon.ThermalState, cnt int) error {
	return nil
}

func (driver *openBMCDriver) GetAllSensorState(state *pluginCommon.SensorState) error {
	return errors.New("Not supported")
}

func (driver *openBMCDriver) GetQsfpState(Id int32) (retObj pluginCommon.QsfpState, err error) {
	return retObj, nil
}

func (driver *openBMCDriver) GetQsfpPMData(Id int32) (retObj pluginCommon.QsfpPMData, err error) {
	return retObj, nil
}

func (driver *openBMCDriver) GetMaxNumOfQsfp() int {
	driver.logger.Info("Inside Dummy: GetMaxNumOfQsfps()")
	return 0
}

func (driver *openBMCDriver) GetPlatformMgmtDeviceState(state *pluginCommon.PlatformMgmtDeviceState) error {
	return errors.New("Not supported")
}

func (driver *openBMCDriver) GetMaxNumOfPsu() int {
	return 0
}

func (driver *openBMCDriver) GetPsuState(psuId int32) (pluginCommon.PsuState, error) {
	var retObj pluginCommon.PsuState
	return retObj, nil
}

func (driver *openBMCDriver) GetAllPsuState(states []pluginCommon.PsuState, cnt int) error {
	return nil
}

func (driver *openBMCDriver) GetMaxNumOfLed() int {
	return 0
}

func (driver *openBMCDriver) GetLedState(ledId int32) (pluginCommon.LedState, error) {
	var retObj pluginCommon.LedState
	return retObj, nil
}

func (driver *openBMCDriver) GetAllLedState(states []pluginCommon.LedState, cnt int) error {
	return nil
}
