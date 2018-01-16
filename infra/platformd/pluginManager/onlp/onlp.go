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

package onlp

import (
	"errors"
	"fmt"
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"strconv"
	"utils/logging"
)

/*
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include "onlp.h"
#include "pluginCommon.h"
*/
import "C"

type onlpDriver struct {
	logger logging.LoggerIntf
}

var driver onlpDriver

func NewONLPPlugin(params *pluginCommon.PluginInitParams) (*onlpDriver, error) {
	driver.logger = params.Logger
	return &driver, nil
}

func (driver *onlpDriver) Init() error {
	driver.logger.Info("Initializing onlp driver")
	rv := int(C.Init())
	if rv < 0 {
		return errors.New("Error initializing Onlp Driver")
	}
	return nil
}

func (driver *onlpDriver) DeInit() error {
	driver.logger.Info("DeInitializing onlp driver")
	C.DeInit()
	return nil
}

func (driver *onlpDriver) GetFanState(fanId int32) (pluginCommon.FanState, error) {
	var retObj pluginCommon.FanState
	var fanInfo C.fan_info_t

	retVal := int(C.GetFanState(&fanInfo, C.int(fanId)))
	if retVal < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch Fan State of", fanId))
	}
	retObj.FanId = int32(fanInfo.FanId)
	switch int(fanInfo.Mode) {
	case pluginCommon.FAN_MODE_OFF:
		retObj.OperMode = pluginCommon.FAN_MODE_OFF_STR
	case pluginCommon.FAN_MODE_ON:
		retObj.OperMode = pluginCommon.FAN_MODE_ON_STR
	}
	retObj.OperSpeed = int32(fanInfo.Speed)
	switch int(fanInfo.Direction) {
	case pluginCommon.FAN_DIR_B2F:
		retObj.OperDirection = pluginCommon.FAN_DIR_B2F_STR
	case pluginCommon.FAN_DIR_F2B:
		retObj.OperDirection = pluginCommon.FAN_DIR_F2B_STR
	case pluginCommon.FAN_DIR_INVALID:
		retObj.OperDirection = pluginCommon.FAN_DIR_INVALID_STR
	}
	switch int(fanInfo.Status) {
	case pluginCommon.FAN_STATUS_PRESENT:
		retObj.Status = pluginCommon.FAN_STATUS_PRESENT_STR
	case pluginCommon.FAN_STATUS_MISSING:
		retObj.Status = pluginCommon.FAN_STATUS_MISSING_STR
	case pluginCommon.FAN_STATUS_FAILED:
		retObj.Status = pluginCommon.FAN_STATUS_FAILED_STR
	case pluginCommon.FAN_STATUS_NORMAL:
		retObj.Status = pluginCommon.FAN_STATUS_NORMAL_STR
	}
	retObj.Model = C.GoString(&fanInfo.Model[0])
	//states[idx].Model = ""
	retObj.SerialNum = C.GoString(&fanInfo.SerialNum[0])
	return retObj, nil
}

func (driver *onlpDriver) GetFanConfig(fanId int32) (retObj *objects.FanConfig, err error) {
	return retObj, err
}

func (driver *onlpDriver) UpdateFanConfig(cfg *objects.FanConfig) (bool, error) {
	driver.logger.Info("Updating Onlp Fan Config")
	return true, nil
}

func (driver *onlpDriver) GetMaxNumOfFans() int {
	return int(C.GetMaxNumOfFans())
}

func (driver *onlpDriver) GetAllFanState(states []pluginCommon.FanState, cnt int) error {
	var fanInfo []C.fan_info_t

	fanInfo = make([]C.fan_info_t, cnt)
	retVal := int(C.GetAllFanState(&fanInfo[0], C.int(cnt)))
	if retVal < 0 {
		return errors.New(fmt.Sprintln("Unable to fetch the fan State:"))
	}
	for idx := 0; idx < cnt; idx++ {
		if int(fanInfo[idx].valid) == 0 {
			states[idx].Valid = false
			continue
		}
		states[idx].Valid = true
		states[idx].FanId = int32(fanInfo[idx].FanId)
		switch int(fanInfo[idx].Mode) {
		case pluginCommon.FAN_MODE_OFF:
			states[idx].OperMode = pluginCommon.FAN_MODE_OFF_STR
		case pluginCommon.FAN_MODE_ON:
			states[idx].OperMode = pluginCommon.FAN_MODE_ON_STR
		}
		states[idx].OperSpeed = int32(fanInfo[idx].Speed)

		switch int(fanInfo[idx].Direction) {
		case pluginCommon.FAN_DIR_B2F:
			states[idx].OperDirection = pluginCommon.FAN_DIR_B2F_STR
		case pluginCommon.FAN_DIR_F2B:
			states[idx].OperDirection = pluginCommon.FAN_DIR_F2B_STR
		case pluginCommon.FAN_DIR_INVALID:
			states[idx].OperDirection = pluginCommon.FAN_DIR_INVALID_STR
		}

		switch int(fanInfo[idx].Status) {
		case pluginCommon.FAN_STATUS_PRESENT:
			states[idx].Status = pluginCommon.FAN_STATUS_PRESENT_STR
		case pluginCommon.FAN_STATUS_MISSING:
			states[idx].Status = pluginCommon.FAN_STATUS_MISSING_STR
		case pluginCommon.FAN_STATUS_FAILED:
			states[idx].Status = pluginCommon.FAN_STATUS_FAILED_STR
		case pluginCommon.FAN_STATUS_NORMAL:
			states[idx].Status = pluginCommon.FAN_STATUS_NORMAL_STR
		}
		states[idx].Model = C.GoString(&fanInfo[idx].Model[0])
		states[idx].SerialNum = C.GoString(&fanInfo[idx].SerialNum[0])
	}
	return nil
}

func (driver *onlpDriver) GetSfpCnt() int {
	return int(C.GetSfpCnt())
}

func (driver *onlpDriver) GetSfpState(sfpId int32) (pluginCommon.SfpState, error) {
	var retObj pluginCommon.SfpState
	var sfpInfo C.sfp_info_t
	var rt int

	rt = int(C.GetSfpState(&sfpInfo, C.int(sfpId)))
	if rt < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch SFP info for ", sfpId))
	}

	if int(rt) > 0 {
		retObj.SfpId = sfpId
		retObj.SfpPresent = "SfpNotPresent"
		return retObj, nil
	}

	retObj.SfpPresent = "SfpPresent"
	if int(sfpInfo.sfp_los) > 0 {
		retObj.SfpLos = "LaserUp"
	} else {
		retObj.SfpLos = "LaserDown"
	}

	retObj.SerialNum = C.GoString(&sfpInfo.serial_number[0])
	q := strconv.Quote(C.GoStringN(&sfpInfo.eeprom[0], 256))
	retObj.EEPROM = q

	return retObj, nil
}

func (driver *onlpDriver) GetAllSfpState(states []pluginCommon.SfpState, cnt int) error {

	if cnt > driver.GetSfpCnt() {
		return errors.New("Error GetAllSfpState Invalid Count")
	}

	for idx := 0; idx < cnt; idx++ {
		states[idx], _ = driver.GetSfpState(int32(idx))
	}
	return nil
}

func (driver *onlpDriver) GetSfpConfig(sfpId int32) (*objects.SfpConfig, error) {
	var retObj objects.SfpConfig

	retObj.SfpId = sfpId
	return &retObj, nil
}

func (driver *onlpDriver) UpdateSfpConfig(cfg *objects.SfpConfig) (bool, error) {
	driver.logger.Info("Updating Onlp SFP Config")
	return true, nil
}

func (driver *onlpDriver) GetPlatformState() (pluginCommon.PlatformState, error) {
	var retObj pluginCommon.PlatformState
	var sysInfo C.sys_info_t

	rt := int(C.GetPlatformState(&sysInfo))

	if rt < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch System info"))
	}

	retObj.ProductName = C.GoString(&sysInfo.product_name[0])
	retObj.Vendor = C.GoString(&sysInfo.vendor[0])
	retObj.SerialNum = C.GoString(&sysInfo.serial_number[0])
	retObj.Manufacturer = C.GoString(&sysInfo.manufacturer[0])
	retObj.Release = C.GoString(&sysInfo.label_revision[0])
	retObj.PlatformName = C.GoString(&sysInfo.platform_name[0])
	retObj.Version = C.GoString(&sysInfo.onie_version[0])

	return retObj, nil
}

func (driver *onlpDriver) GetMaxNumOfThermal() int {
	return 8
}

func (driver *onlpDriver) GetThermalState(thermalId int32) (pluginCommon.ThermalState, error) {
	var retObj pluginCommon.ThermalState
	var tInfo C.thermal_info_t

	rt := int(C.GetThermalState(&tInfo, C.int(thermalId)))
	if rt < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch sensor state of", thermalId))
	}

	retObj.ThermalId = int32(tInfo.sensor_id)
	retObj.Location = C.GoString(&tInfo.description[0])
	retObj.Temperature = strconv.Itoa(int(tInfo.temp))
	retObj.LowerWatermarkTemperature = strconv.Itoa(int(tInfo.threshold_warning))
	retObj.UpperWatermarkTemperature = strconv.Itoa(int(tInfo.threshold_error))
	retObj.ShutdownTemperature = strconv.Itoa(int(tInfo.threshold_shutdown))

	return retObj, nil
}

func (driver *onlpDriver) GetAllThermalState(states []pluginCommon.ThermalState, cnt int) error {

	if cnt > driver.GetMaxNumOfThermal() {
		return errors.New("Error GetAllThermalState Invalid Count")
	}

	for idx := 1; idx <= cnt; idx++ {
		states[idx], _ = driver.GetThermalState(int32(idx))
	}
	return nil
}

func (driver *onlpDriver) GetMaxNumOfPsu() int {
	return 2
}

func (driver *onlpDriver) GetPsuState(psuId int32) (pluginCommon.PsuState, error) {
	var retObj pluginCommon.PsuState
	var pInfo C.psu_info_t

	rt := int(C.GetPsuState(&pInfo, C.int(psuId)))
	if rt < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch PSU state of", psuId))
	}

	retObj.PsuId = int32(pInfo.psu_id)

	if pInfo.status != 0 {
		retObj.Status = "PSU PRESENT"
		retObj.VoltIn = int32(pInfo.mvin)
		retObj.VoltOut = int32(pInfo.mvout)
		retObj.AmpIn = int32(pInfo.miin)
		retObj.AmpOut = int32(pInfo.miout)
		retObj.PwrIn = int32(pInfo.mpin)
		retObj.PwrOut = int32(pInfo.mpout)
	} else {
		retObj.Status = "PSU UNPLUGGED"
	}

	return retObj, nil
}

func (driver *onlpDriver) GetAllPsuState(states []pluginCommon.PsuState, cnt int) error {

	if cnt > driver.GetMaxNumOfPsu() {
		return errors.New("Error GetAllPsuState Invalid Count")
	}

	for idx := 0; idx < cnt; idx++ {
		states[idx], _ = driver.GetPsuState(int32(idx))
	}
	return nil
}

func (driver *onlpDriver) GetAllSensorState(state *pluginCommon.SensorState) error {
	return errors.New("Not supported")
}

func (driver *onlpDriver) GetQsfpState(Id int32) (retObj pluginCommon.QsfpState, err error) {
	return retObj, nil
}

func (driver *onlpDriver) GetQsfpPMData(Id int32) (retObj pluginCommon.QsfpPMData, err error) {
	return retObj, nil
}

func (driver *onlpDriver) GetMaxNumOfQsfp() int {
	driver.logger.Info("Inside Dummy: GetMaxNumOfQsfps()")
	return 0
}

func (driver *onlpDriver) GetPlatformMgmtDeviceState(state *pluginCommon.PlatformMgmtDeviceState) error {
	return errors.New("Not supported")
}

func (driver *onlpDriver) GetMaxNumOfLed() int {
	return 12
}

func (driver *onlpDriver) GetLedState(ledId int32) (pluginCommon.LedState, error) {
	var retObj pluginCommon.LedState
	var ledInfo C.led_info_t

	rt := int(C.GetLedState(&ledInfo, C.int(ledId)))
	if rt < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch Led state of", ledId))
	}

	retObj.LedId = int32(ledInfo.led_id)

	if ledInfo.status != 0 {
		retObj.LedState = "LED PRESENT"
		retObj.LedIdentify = C.GoString(&ledInfo.description[0])
		retObj.LedColor = C.GoString(&ledInfo.color[0])
	} else {
		retObj.LedState = "LED NOT PRESENT"
	}

	return retObj, nil
}

func (driver *onlpDriver) GetAllLedState(states []pluginCommon.LedState, cnt int) error {

	if cnt > driver.GetMaxNumOfLed() {
		return errors.New("Error GetAllLedState Invalid Count")
	}

	for idx := 0; idx < cnt; idx++ {
		states[idx], _ = driver.GetLedState(int32(idx))
	}
	return nil
}
