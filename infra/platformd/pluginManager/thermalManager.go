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
	"errors"
	"fmt"
	"infra/platformd/objects"
	"utils/logging"
)

type ThermalId int32

type ThermalState struct {
	Location                  string
	Temperature               string
	LowerWatermarkTemperature string
	UpperWatermarkTemperature string
	ShutdownTemperature       string
}

type ThermalManager struct {
	logger      logging.LoggerIntf
	plugin      PluginIntf
	numOfSensor int
}

var ThermalMgr ThermalManager

func (tMgr *ThermalManager) Init(logger logging.LoggerIntf, plugin PluginIntf) {
	tMgr.logger = logger
	tMgr.plugin = plugin
	tMgr.numOfSensor = tMgr.plugin.GetMaxNumOfThermal()
	tMgr.logger.Info("Thermal Manager Init()")
}

func (tMgr *ThermalManager) Deinit() {
	tMgr.logger.Info("Thermal Manager Deinit()")
}

func (tMgr *ThermalManager) GetThermalState(thermalId int32) (*objects.ThermalState, error) {
	var thermalObj objects.ThermalState
	if tMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}

	thermalState, err := tMgr.plugin.GetThermalState(thermalId)
	if err != nil {
		return nil, err
	}
	thermalObj.ThermalId = thermalId
	thermalObj.Location = thermalState.Location
	thermalObj.Temperature = thermalState.Temperature
	thermalObj.LowerWatermarkTemperature = thermalState.LowerWatermarkTemperature
	thermalObj.UpperWatermarkTemperature = thermalState.UpperWatermarkTemperature
	thermalObj.ShutdownTemperature = thermalState.ShutdownTemperature
	return &thermalObj, err
}

func (tMgr *ThermalManager) GetBulkThermalState(fromIdx int, cnt int) (*objects.ThermalStateGetInfo, error) {
	var retObj objects.ThermalStateGetInfo
	if tMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= tMgr.numOfSensor {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > tMgr.numOfSensor {
		retObj.EndIdx = tMgr.numOfSensor
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = tMgr.numOfSensor - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		thermalId := int32(idx)
		obj, err := tMgr.GetThermalState(thermalId)
		if err != nil {
			tMgr.logger.Err(fmt.Sprintln("Error getting the thermal state for thermalId:", thermalId))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}
