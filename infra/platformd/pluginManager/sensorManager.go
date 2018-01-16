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
	"infra/platformd/pluginManager/pluginCommon"
	"models/events"
	"sync"
	"time"
	"utils/dbutils"
	"utils/eventUtils"
	"utils/logging"
	"utils/ringBuffer"
)

type FanSensorConfig struct {
	AdminState             string
	HigherAlarmThreshold   int32
	HigherWarningThreshold int32
	LowerWarningThreshold  int32
	LowerAlarmThreshold    int32
	PMClassAAdminState     string
	PMClassBAdminState     string
	PMClassCAdminState     string
}

type TemperatureSensorConfig struct {
	AdminState             string
	HigherAlarmThreshold   float64
	HigherWarningThreshold float64
	LowerWarningThreshold  float64
	LowerAlarmThreshold    float64
	PMClassAAdminState     string
	PMClassBAdminState     string
	PMClassCAdminState     string
}

type VoltageSensorConfig struct {
	AdminState             string
	HigherAlarmThreshold   float64
	HigherWarningThreshold float64
	LowerWarningThreshold  float64
	LowerAlarmThreshold    float64
	PMClassAAdminState     string
	PMClassBAdminState     string
	PMClassCAdminState     string
}

type PowerConverterSensorConfig struct {
	AdminState             string
	HigherAlarmThreshold   float64
	HigherWarningThreshold float64
	LowerWarningThreshold  float64
	LowerAlarmThreshold    float64
	PMClassAAdminState     string
	PMClassBAdminState     string
	PMClassCAdminState     string
}

type FanSensorEventData struct {
	Value int32
}

type TempSensorEventData struct {
	Value float64
}

type VoltageSensorEventData struct {
	Value float64
}

type PowerConverterSensorEventData struct {
	Value float64
}

const (
	sensorClassAInterval time.Duration = time.Duration(10) * time.Second // Polling Interval 10 sec
	sensorClassABufSize  int           = 6 * 60 * 24                     //Storage for 24 hrs
	sensorClassBInterval time.Duration = time.Duration(15) * time.Minute // Polling Interval 15 mins
	sensorClassBBufSize  int           = 4 * 24                          // Storage for 24 hrs
	sensorClassCInterval time.Duration = time.Duration(24) * time.Hour   // Polling Interval 24 Hrs
	sensorClassCBufSize  int           = 365                             // Storage for 365 days
)

type EventStatus struct {
	SentHigherAlarm bool
	SentHigherWarn  bool
	SentLowerWarn   bool
	SentLowerAlarm  bool
}

type SensorManager struct {
	logger                       logging.LoggerIntf
	plugin                       PluginIntf
	eventDbHdl                   dbutils.DBIntf
	classAPMTimer                *time.Timer
	classBPMTimer                *time.Timer
	classCPMTimer                *time.Timer
	fanSensorList                []string
	fanConfigMutex               sync.RWMutex
	fanConfigDB                  map[string]FanSensorConfig
	fanMsgStatusMutex            sync.RWMutex
	fanMsgStatus                 map[string]EventStatus
	fanClassAPMMutex             sync.RWMutex
	fanSensorClassAPM            map[string]*ringBuffer.RingBuffer
	fanClassBPMMutex             sync.RWMutex
	fanSensorClassBPM            map[string]*ringBuffer.RingBuffer
	fanClassCPMMutex             sync.RWMutex
	fanSensorClassCPM            map[string]*ringBuffer.RingBuffer
	tempSensorList               []string
	tempConfigMutex              sync.RWMutex
	tempConfigDB                 map[string]TemperatureSensorConfig
	tempMsgStatusMutex           sync.RWMutex
	tempMsgStatus                map[string]EventStatus
	tempClassAPMMutex            sync.RWMutex
	tempSensorClassAPM           map[string]*ringBuffer.RingBuffer
	tempClassBPMMutex            sync.RWMutex
	tempSensorClassBPM           map[string]*ringBuffer.RingBuffer
	tempClassCPMMutex            sync.RWMutex
	tempSensorClassCPM           map[string]*ringBuffer.RingBuffer
	voltageSensorList            []string
	voltageConfigMutex           sync.RWMutex
	voltageConfigDB              map[string]VoltageSensorConfig
	voltageMsgStatusMutex        sync.RWMutex
	voltageMsgStatus             map[string]EventStatus
	voltageClassAPMMutex         sync.RWMutex
	voltageSensorClassAPM        map[string]*ringBuffer.RingBuffer
	voltageClassBPMMutex         sync.RWMutex
	voltageSensorClassBPM        map[string]*ringBuffer.RingBuffer
	voltageClassCPMMutex         sync.RWMutex
	voltageSensorClassCPM        map[string]*ringBuffer.RingBuffer
	powerConverterSensorList     []string
	powerConverterConfigMutex    sync.RWMutex
	powerConverterConfigDB       map[string]PowerConverterSensorConfig
	powerConverterMsgStatusMutex sync.RWMutex
	powerConverterMsgStatus      map[string]EventStatus
	powerConverterClassAPMMutex  sync.RWMutex
	powerConverterSensorClassAPM map[string]*ringBuffer.RingBuffer
	powerConverterClassBPMMutex  sync.RWMutex
	powerConverterSensorClassBPM map[string]*ringBuffer.RingBuffer
	powerConverterClassCPMMutex  sync.RWMutex
	powerConverterSensorClassCPM map[string]*ringBuffer.RingBuffer
	SensorStateMutex             sync.RWMutex
	SensorState                  *pluginCommon.SensorState
}

var SensorMgr SensorManager

func (sMgr *SensorManager) Init(logger logging.LoggerIntf, plugin PluginIntf, eventDbHdl dbutils.DBIntf) {
	var evtStatus EventStatus
	sMgr.logger = logger
	sMgr.eventDbHdl = eventDbHdl
	sMgr.logger.Info("sensor Manager Init( Start)")
	sMgr.plugin = plugin
	sMgr.SensorState = new(pluginCommon.SensorState)
	sMgr.SensorState.FanSensor = make(map[string]pluginCommon.FanSensorData)
	sMgr.SensorState.TemperatureSensor = make(map[string]pluginCommon.TemperatureSensorData)
	sMgr.SensorState.VoltageSensor = make(map[string]pluginCommon.VoltageSensorData)
	sMgr.SensorState.PowerConverterSensor = make(map[string]pluginCommon.PowerConverterSensorData)
	sMgr.fanSensorClassAPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.fanSensorClassBPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.fanSensorClassCPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.tempSensorClassAPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.tempSensorClassBPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.tempSensorClassCPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.voltageSensorClassAPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.voltageSensorClassBPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.voltageSensorClassCPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.powerConverterSensorClassAPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.powerConverterSensorClassBPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.powerConverterSensorClassCPM = make(map[string]*ringBuffer.RingBuffer)

	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.logger.Info("Sensor Manager Init() Failed")
		return
	}
	sMgr.logger.Info("Sensor State:", sMgr.SensorState)
	sMgr.fanConfigDB = make(map[string]FanSensorConfig)
	sMgr.fanMsgStatus = make(map[string]EventStatus)
	for name, _ := range sMgr.SensorState.FanSensor {
		sMgr.logger.Info("Fan Sensor:", name)
		sMgr.fanConfigMutex.Lock()
		fanCfgEnt, _ := sMgr.fanConfigDB[name]
		// TODO: Read Json
		fanCfgEnt.AdminState = "Enable"
		fanCfgEnt.HigherAlarmThreshold = 11000
		fanCfgEnt.HigherWarningThreshold = 11000
		fanCfgEnt.LowerAlarmThreshold = 1000
		fanCfgEnt.LowerWarningThreshold = 1000
		fanCfgEnt.PMClassAAdminState = "Enable"
		fanCfgEnt.PMClassBAdminState = "Enable"
		fanCfgEnt.PMClassCAdminState = "Enable"
		sMgr.fanConfigDB[name] = fanCfgEnt
		sMgr.fanMsgStatus[name] = evtStatus
		sMgr.fanConfigMutex.Unlock()
		sMgr.fanSensorList = append(sMgr.fanSensorList, name)
	}
	sMgr.tempConfigDB = make(map[string]TemperatureSensorConfig)
	sMgr.tempMsgStatus = make(map[string]EventStatus)
	for name, _ := range sMgr.SensorState.TemperatureSensor {
		sMgr.logger.Info("Temperature Sensor:", name)
		sMgr.tempConfigMutex.Lock()
		tempCfgEnt, _ := sMgr.tempConfigDB[name]
		// TODO: Read Json
		tempCfgEnt.AdminState = "Enable"
		tempCfgEnt.HigherAlarmThreshold = 11000.0
		tempCfgEnt.HigherWarningThreshold = 11000.0
		tempCfgEnt.LowerAlarmThreshold = -1000.0
		tempCfgEnt.LowerWarningThreshold = -1000.0
		tempCfgEnt.PMClassAAdminState = "Enable"
		tempCfgEnt.PMClassBAdminState = "Enable"
		tempCfgEnt.PMClassCAdminState = "Enable"
		sMgr.tempConfigDB[name] = tempCfgEnt
		sMgr.tempMsgStatus[name] = evtStatus
		sMgr.tempConfigMutex.Unlock()
		sMgr.tempSensorList = append(sMgr.tempSensorList, name)
	}
	sMgr.voltageConfigDB = make(map[string]VoltageSensorConfig)
	sMgr.voltageMsgStatus = make(map[string]EventStatus)
	for name, _ := range sMgr.SensorState.VoltageSensor {
		sMgr.logger.Info("Voltage Sensor:", name)
		sMgr.voltageConfigMutex.Lock()
		voltageCfgEnt, _ := sMgr.voltageConfigDB[name]
		// TODO: Read Json
		voltageCfgEnt.AdminState = "Enable"
		voltageCfgEnt.HigherAlarmThreshold = 11000
		voltageCfgEnt.HigherWarningThreshold = 11000
		voltageCfgEnt.LowerAlarmThreshold = 0
		voltageCfgEnt.LowerWarningThreshold = 0
		voltageCfgEnt.PMClassAAdminState = "Enable"
		voltageCfgEnt.PMClassBAdminState = "Enable"
		voltageCfgEnt.PMClassCAdminState = "Enable"
		sMgr.voltageConfigDB[name] = voltageCfgEnt
		sMgr.voltageMsgStatus[name] = evtStatus
		sMgr.voltageConfigMutex.Unlock()
		sMgr.voltageSensorList = append(sMgr.voltageSensorList, name)
	}
	sMgr.powerConverterConfigDB = make(map[string]PowerConverterSensorConfig)
	sMgr.powerConverterMsgStatus = make(map[string]EventStatus)
	for name, _ := range sMgr.SensorState.PowerConverterSensor {
		sMgr.logger.Info("Power Sensor:", name)
		sMgr.powerConverterConfigMutex.Lock()
		powerConverterCfgEnt, _ := sMgr.powerConverterConfigDB[name]
		// TODO: Read Json
		powerConverterCfgEnt.AdminState = "Enable"
		powerConverterCfgEnt.HigherAlarmThreshold = 11000
		powerConverterCfgEnt.HigherWarningThreshold = 11000
		powerConverterCfgEnt.LowerAlarmThreshold = 0
		powerConverterCfgEnt.LowerWarningThreshold = 0
		powerConverterCfgEnt.PMClassAAdminState = "Enable"
		powerConverterCfgEnt.PMClassBAdminState = "Enable"
		powerConverterCfgEnt.PMClassCAdminState = "Enable"
		sMgr.powerConverterConfigDB[name] = powerConverterCfgEnt
		sMgr.powerConverterMsgStatus[name] = evtStatus
		sMgr.powerConverterConfigMutex.Unlock()
		sMgr.powerConverterSensorList = append(sMgr.powerConverterSensorList, name)
	}
	sMgr.SensorStateMutex.Unlock()
	sMgr.StartSensorPM()
	sMgr.logger.Info("sensor Manager Init( Done)")
}

func (sMgr *SensorManager) Deinit() {
	sMgr.logger.Info("Sensor Manager Deinit()")
}

func (sMgr *SensorManager) GetFanSensorState(Name string) (*objects.FanSensorState, error) {
	var fanSensorObj objects.FanSensorState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.fanConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid Fan Sensor Name")
	}

	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.SensorStateMutex.Unlock()
		return nil, err
	}
	fanSensorState, _ := sMgr.SensorState.FanSensor[Name]
	fanSensorObj.Name = Name
	fanSensorObj.CurrentSpeed = fanSensorState.Value
	sMgr.SensorStateMutex.Unlock()
	return &fanSensorObj, err
}

func (sMgr *SensorManager) GetBulkFanSensorState(fromIdx int, cnt int) (*objects.FanSensorStateGetInfo, error) {
	var retObj objects.FanSensorStateGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.fanSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.fanSensorList) {
		retObj.EndIdx = len(sMgr.fanSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.fanSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		fanSensorName := sMgr.fanSensorList[idx]
		obj, err := sMgr.GetFanSensorState(fanSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the fan state for fan Sensor:", fanSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (sMgr *SensorManager) GetFanSensorConfig(Name string) (*objects.FanSensorConfig, error) {
	var fanSensorObj objects.FanSensorConfig
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.fanConfigMutex.RLock()
	fanSensorCfgEnt, exist := sMgr.fanConfigDB[Name]
	if !exist {
		sMgr.fanConfigMutex.RUnlock()
		return nil, errors.New("Invalid Fan Sensor Name")
	}
	fanSensorObj.Name = Name
	fanSensorObj.AdminState = fanSensorCfgEnt.AdminState
	fanSensorObj.HigherAlarmThreshold = fanSensorCfgEnt.HigherAlarmThreshold
	fanSensorObj.HigherWarningThreshold = fanSensorCfgEnt.HigherWarningThreshold
	fanSensorObj.LowerAlarmThreshold = fanSensorCfgEnt.LowerAlarmThreshold
	fanSensorObj.LowerWarningThreshold = fanSensorCfgEnt.LowerWarningThreshold
	fanSensorObj.PMClassAAdminState = fanSensorCfgEnt.PMClassAAdminState
	fanSensorObj.PMClassBAdminState = fanSensorCfgEnt.PMClassBAdminState
	fanSensorObj.PMClassCAdminState = fanSensorCfgEnt.PMClassCAdminState
	sMgr.fanConfigMutex.RUnlock()
	return &fanSensorObj, nil
}

func (sMgr *SensorManager) GetBulkFanSensorConfig(fromIdx int, cnt int) (*objects.FanSensorConfigGetInfo, error) {
	var retObj objects.FanSensorConfigGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.fanSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.fanSensorList) {
		retObj.EndIdx = len(sMgr.fanSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.fanSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		fanSensorName := sMgr.fanSensorList[idx]
		obj, err := sMgr.GetFanSensorConfig(fanSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the fan state for fan sensor:", fanSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func genFanSensorUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0

	if attrset == nil {
		mask = objects.FAN_SENSOR_UPDATE_ADMIN_STATE |
			objects.FAN_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD |
			objects.FAN_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD |
			objects.FAN_SENSOR_UPDATE_LOWER_WARN_THRESHOLD |
			objects.FAN_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD |
			objects.FAN_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE |
			objects.FAN_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE |
			objects.FAN_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE
	} else {
		for idx, val := range attrset {
			if true == val {
				switch idx {
				case 0:
					//ObjKey Fan Name
				case 1:
					mask |= objects.FAN_SENSOR_UPDATE_ADMIN_STATE
				case 2:
					mask |= objects.FAN_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD
				case 3:
					mask |= objects.FAN_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD
				case 4:
					mask |= objects.FAN_SENSOR_UPDATE_LOWER_WARN_THRESHOLD
				case 5:
					mask |= objects.FAN_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD
				case 6:
					mask |= objects.FAN_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE
				case 7:
					mask |= objects.FAN_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE
				case 8:
					mask |= objects.FAN_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE
				}
			}
		}
	}
	return mask
}

func (sMgr *SensorManager) UpdateFanSensorConfig(oldCfg *objects.FanSensorConfig, newCfg *objects.FanSensorConfig, attrset []bool) (bool, error) {
	if sMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	sMgr.fanConfigMutex.Lock()
	fanSensorCfgEnt, exist := sMgr.fanConfigDB[newCfg.Name]
	if !exist {
		sMgr.fanConfigMutex.Unlock()
		return false, errors.New("Invalid FanSensor Name")
	}
	mask := genFanSensorUpdateMask(attrset)
	if mask&objects.FAN_SENSOR_UPDATE_ADMIN_STATE == objects.FAN_SENSOR_UPDATE_ADMIN_STATE {
		if newCfg.AdminState != "Enable" && newCfg.AdminState != "Disable" {
			sMgr.fanConfigMutex.Unlock()
			return false, errors.New("Invalid AdminState Value")
		}
		fanSensorCfgEnt.AdminState = newCfg.AdminState
	}
	if mask&objects.FAN_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD == objects.FAN_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD {
		fanSensorCfgEnt.HigherAlarmThreshold = newCfg.HigherAlarmThreshold
	}
	if mask&objects.FAN_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD == objects.FAN_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD {
		fanSensorCfgEnt.HigherWarningThreshold = newCfg.HigherWarningThreshold
	}
	if mask&objects.FAN_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD == objects.FAN_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD {
		fanSensorCfgEnt.LowerAlarmThreshold = newCfg.LowerAlarmThreshold
	}
	if mask&objects.FAN_SENSOR_UPDATE_LOWER_WARN_THRESHOLD == objects.FAN_SENSOR_UPDATE_LOWER_WARN_THRESHOLD {
		fanSensorCfgEnt.LowerWarningThreshold = newCfg.LowerWarningThreshold
	}
	if mask&objects.FAN_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE == objects.FAN_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE {
		if newCfg.PMClassAAdminState != "Enable" && newCfg.PMClassAAdminState != "Disable" {
			sMgr.fanConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassAAdminState Value")
		}
		fanSensorCfgEnt.PMClassAAdminState = newCfg.PMClassAAdminState
	}
	if mask&objects.FAN_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE == objects.FAN_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE {
		if newCfg.PMClassBAdminState != "Enable" && newCfg.PMClassBAdminState != "Disable" {
			sMgr.fanConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassBAdminState Value")
		}
		fanSensorCfgEnt.PMClassBAdminState = newCfg.PMClassBAdminState
	}
	if mask&objects.FAN_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE == objects.FAN_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE {
		if newCfg.PMClassCAdminState != "Enable" && newCfg.PMClassCAdminState != "Disable" {
			sMgr.fanConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassCAdminState Value")
		}
		fanSensorCfgEnt.PMClassCAdminState = newCfg.PMClassCAdminState
	}

	if !(fanSensorCfgEnt.HigherAlarmThreshold >= fanSensorCfgEnt.HigherWarningThreshold &&
		fanSensorCfgEnt.HigherWarningThreshold > fanSensorCfgEnt.LowerWarningThreshold &&
		fanSensorCfgEnt.LowerWarningThreshold >= fanSensorCfgEnt.LowerAlarmThreshold) {
		sMgr.fanConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, Please verify the thresholds")
	}
	if sMgr.fanConfigDB[newCfg.Name].AdminState != fanSensorCfgEnt.AdminState {
		if fanSensorCfgEnt.AdminState == "Disable" {
			//Clear all alarms
			sMgr.clearExistingFanSensorFaults(newCfg.Name)
			sMgr.logger.Info("Clear all the existing faults")
		}
	}
	if sMgr.fanConfigDB[newCfg.Name].PMClassAAdminState != fanSensorCfgEnt.PMClassAAdminState {
		if fanSensorCfgEnt.PMClassAAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.fanClassAPMMutex.Lock()
			sMgr.fanSensorClassAPM[newCfg.Name].FlushRingBuffer()
			sMgr.fanClassAPMMutex.Unlock()
			sMgr.logger.Info("Flush Class A PM Ring Buffer")
		}
	}
	if sMgr.fanConfigDB[newCfg.Name].PMClassBAdminState != fanSensorCfgEnt.PMClassBAdminState {
		if fanSensorCfgEnt.PMClassBAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.fanClassBPMMutex.Lock()
			sMgr.fanSensorClassBPM[newCfg.Name].FlushRingBuffer()
			sMgr.fanClassBPMMutex.Unlock()
			sMgr.logger.Info("Flush Class B PM Ring Buffer")
		}
	}
	if sMgr.fanConfigDB[newCfg.Name].PMClassCAdminState != fanSensorCfgEnt.PMClassCAdminState {
		if fanSensorCfgEnt.PMClassCAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.fanClassCPMMutex.Lock()
			sMgr.fanSensorClassCPM[newCfg.Name].FlushRingBuffer()
			sMgr.fanClassCPMMutex.Unlock()
			sMgr.logger.Info("Flush Class C PM Ring Buffer")
		}
	}
	sMgr.fanConfigDB[newCfg.Name] = fanSensorCfgEnt
	sMgr.fanConfigMutex.Unlock()
	return true, nil
}

func (sMgr *SensorManager) GetTemperatureSensorState(Name string) (*objects.TemperatureSensorState, error) {
	var tempSensorObj objects.TemperatureSensorState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.tempConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid Temperature Sensor Name")
	}

	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.SensorStateMutex.Unlock()
		return nil, err
	}
	tempSensorState, _ := sMgr.SensorState.TemperatureSensor[Name]
	tempSensorObj.Name = Name
	tempSensorObj.CurrentTemperature = tempSensorState.Value
	sMgr.SensorStateMutex.Unlock()
	return &tempSensorObj, err
}

func (sMgr *SensorManager) GetBulkTemperatureSensorState(fromIdx int, cnt int) (*objects.TemperatureSensorStateGetInfo, error) {
	var retObj objects.TemperatureSensorStateGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.tempSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.tempSensorList) {
		retObj.EndIdx = len(sMgr.tempSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.tempSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		tempSensorName := sMgr.tempSensorList[idx]
		obj, err := sMgr.GetTemperatureSensorState(tempSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the temp state for temp Sensor:", tempSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (sMgr *SensorManager) GetTemperatureSensorConfig(Name string) (*objects.TemperatureSensorConfig, error) {
	var tempSensorObj objects.TemperatureSensorConfig
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.tempConfigMutex.RLock()
	tempSensorCfgEnt, exist := sMgr.tempConfigDB[Name]
	if !exist {
		sMgr.tempConfigMutex.RUnlock()
		return nil, errors.New("Invalid Temperature Sensor Name")
	}
	tempSensorObj.Name = Name
	tempSensorObj.AdminState = tempSensorCfgEnt.AdminState
	tempSensorObj.HigherAlarmThreshold = tempSensorCfgEnt.HigherAlarmThreshold
	tempSensorObj.HigherWarningThreshold = tempSensorCfgEnt.HigherWarningThreshold
	tempSensorObj.LowerAlarmThreshold = tempSensorCfgEnt.LowerAlarmThreshold
	tempSensorObj.LowerWarningThreshold = tempSensorCfgEnt.LowerWarningThreshold
	tempSensorObj.PMClassAAdminState = tempSensorCfgEnt.PMClassAAdminState
	tempSensorObj.PMClassBAdminState = tempSensorCfgEnt.PMClassBAdminState
	tempSensorObj.PMClassCAdminState = tempSensorCfgEnt.PMClassCAdminState
	sMgr.tempConfigMutex.RUnlock()
	return &tempSensorObj, nil
}

func (sMgr *SensorManager) GetBulkTemperatureSensorConfig(fromIdx int, cnt int) (*objects.TemperatureSensorConfigGetInfo, error) {
	var retObj objects.TemperatureSensorConfigGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.tempSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.tempSensorList) {
		retObj.EndIdx = len(sMgr.tempSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.tempSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		tempSensorName := sMgr.tempSensorList[idx]
		obj, err := sMgr.GetTemperatureSensorConfig(tempSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the temp state for temp sensor:", tempSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func genTempSensorUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0
	if attrset == nil {
		mask = objects.TEMP_SENSOR_UPDATE_ADMIN_STATE |
			objects.TEMP_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD |
			objects.TEMP_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD |
			objects.TEMP_SENSOR_UPDATE_LOWER_WARN_THRESHOLD |
			objects.TEMP_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD |
			objects.TEMP_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE |
			objects.TEMP_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE |
			objects.TEMP_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE
	} else {
		for idx, val := range attrset {
			if true == val {
				switch idx {
				case 0:
					//ObjKey Temp Sensor Name
				case 1:
					mask |= objects.TEMP_SENSOR_UPDATE_ADMIN_STATE
				case 2:
					mask |= objects.TEMP_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD
				case 3:
					mask |= objects.TEMP_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD
				case 4:
					mask |= objects.TEMP_SENSOR_UPDATE_LOWER_WARN_THRESHOLD
				case 5:
					mask |= objects.TEMP_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD
				case 6:
					mask |= objects.TEMP_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE
				case 7:
					mask |= objects.TEMP_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE
				case 8:
					mask |= objects.TEMP_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE
				}
			}
		}
	}
	return mask
}

func (sMgr *SensorManager) UpdateTemperatureSensorConfig(oldCfg *objects.TemperatureSensorConfig, newCfg *objects.TemperatureSensorConfig, attrset []bool) (bool, error) {
	if sMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	sMgr.tempConfigMutex.Lock()
	tempSensorCfgEnt, exist := sMgr.tempConfigDB[newCfg.Name]
	if !exist {
		sMgr.tempConfigMutex.Unlock()
		return false, errors.New("Invalid TemperatureSensor Name")
	}
	mask := genTempSensorUpdateMask(attrset)
	if mask&objects.TEMP_SENSOR_UPDATE_ADMIN_STATE == objects.TEMP_SENSOR_UPDATE_ADMIN_STATE {
		if newCfg.AdminState != "Enable" && newCfg.AdminState != "Disable" {
			sMgr.tempConfigMutex.Unlock()
			return false, errors.New("Invalid AdminState Value")
		}
		tempSensorCfgEnt.AdminState = newCfg.AdminState
	}
	if mask&objects.TEMP_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD == objects.TEMP_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD {
		tempSensorCfgEnt.HigherAlarmThreshold = newCfg.HigherAlarmThreshold
	}
	if mask&objects.TEMP_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD == objects.TEMP_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD {
		tempSensorCfgEnt.HigherWarningThreshold = newCfg.HigherWarningThreshold
	}
	if mask&objects.TEMP_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD == objects.TEMP_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD {
		tempSensorCfgEnt.LowerAlarmThreshold = newCfg.LowerAlarmThreshold
	}
	if mask&objects.TEMP_SENSOR_UPDATE_LOWER_WARN_THRESHOLD == objects.TEMP_SENSOR_UPDATE_LOWER_WARN_THRESHOLD {
		tempSensorCfgEnt.LowerWarningThreshold = newCfg.LowerWarningThreshold
	}
	if mask&objects.TEMP_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE == objects.TEMP_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE {
		if newCfg.PMClassAAdminState != "Enable" && newCfg.PMClassAAdminState != "Disable" {
			sMgr.tempConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassAAdminState Value")
		}
		tempSensorCfgEnt.PMClassAAdminState = newCfg.PMClassAAdminState
	}
	if mask&objects.TEMP_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE == objects.TEMP_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE {
		if newCfg.PMClassBAdminState != "Enable" && newCfg.PMClassBAdminState != "Disable" {
			sMgr.tempConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassBAdminState Value")
		}
		tempSensorCfgEnt.PMClassBAdminState = newCfg.PMClassBAdminState
	}
	if mask&objects.TEMP_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE == objects.TEMP_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE {
		if newCfg.PMClassCAdminState != "Enable" && newCfg.PMClassCAdminState != "Disable" {
			sMgr.tempConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassCAdminState Value")
		}
		tempSensorCfgEnt.PMClassCAdminState = newCfg.PMClassCAdminState
	}
	if !(tempSensorCfgEnt.HigherAlarmThreshold >= tempSensorCfgEnt.HigherWarningThreshold &&
		tempSensorCfgEnt.HigherWarningThreshold > tempSensorCfgEnt.LowerWarningThreshold &&
		tempSensorCfgEnt.LowerWarningThreshold >= tempSensorCfgEnt.LowerAlarmThreshold) {
		sMgr.tempConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, Please verify the thresholds")
	}
	if sMgr.tempConfigDB[newCfg.Name].AdminState != tempSensorCfgEnt.AdminState {
		if tempSensorCfgEnt.AdminState == "Disable" {
			//Clear all alarms
			sMgr.clearExistingTempSensorFaults(newCfg.Name)
			sMgr.logger.Info("Clear all the existing faults")
		}
	}
	if sMgr.tempConfigDB[newCfg.Name].PMClassAAdminState != tempSensorCfgEnt.PMClassAAdminState {
		if tempSensorCfgEnt.PMClassAAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.tempClassAPMMutex.Lock()
			sMgr.tempSensorClassAPM[newCfg.Name].FlushRingBuffer()
			sMgr.tempClassAPMMutex.Unlock()
			sMgr.logger.Info("Flush Class A PM Ring Buffer")
		}
	}
	if sMgr.tempConfigDB[newCfg.Name].PMClassBAdminState != tempSensorCfgEnt.PMClassBAdminState {
		if tempSensorCfgEnt.PMClassBAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.tempClassBPMMutex.Lock()
			sMgr.tempSensorClassBPM[newCfg.Name].FlushRingBuffer()
			sMgr.tempClassBPMMutex.Unlock()
			sMgr.logger.Info("Flush Class B PM Ring Buffer")
		}
	}
	if sMgr.tempConfigDB[newCfg.Name].PMClassCAdminState != tempSensorCfgEnt.PMClassCAdminState {
		if tempSensorCfgEnt.PMClassCAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.tempClassCPMMutex.Lock()
			sMgr.tempSensorClassCPM[newCfg.Name].FlushRingBuffer()
			sMgr.tempClassCPMMutex.Unlock()
			sMgr.logger.Info("Flush Class C PM Ring Buffer")
		}
	}
	sMgr.tempConfigDB[newCfg.Name] = tempSensorCfgEnt
	sMgr.tempConfigMutex.Unlock()
	return true, nil
}

func (sMgr *SensorManager) GetVoltageSensorState(Name string) (*objects.VoltageSensorState, error) {
	var voltageSensorObj objects.VoltageSensorState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.voltageConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid Voltage Sensor Name")
	}

	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.SensorStateMutex.Unlock()
		return nil, err
	}
	voltageSensorState, _ := sMgr.SensorState.VoltageSensor[Name]
	voltageSensorObj.Name = Name
	voltageSensorObj.CurrentVoltage = voltageSensorState.Value
	sMgr.SensorStateMutex.Unlock()
	return &voltageSensorObj, err
}

func (sMgr *SensorManager) GetBulkVoltageSensorState(fromIdx int, cnt int) (*objects.VoltageSensorStateGetInfo, error) {
	var retObj objects.VoltageSensorStateGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.voltageSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.voltageSensorList) {
		retObj.EndIdx = len(sMgr.voltageSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.voltageSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		voltageSensorName := sMgr.voltageSensorList[idx]
		obj, err := sMgr.GetVoltageSensorState(voltageSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the voltage state for voltage Sensor:", voltageSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (sMgr *SensorManager) GetVoltageSensorConfig(Name string) (*objects.VoltageSensorConfig, error) {
	var voltageSensorObj objects.VoltageSensorConfig
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.voltageConfigMutex.RLock()
	voltageSensorCfgEnt, exist := sMgr.voltageConfigDB[Name]
	if !exist {
		sMgr.voltageConfigMutex.RUnlock()
		return nil, errors.New("Invalid Voltage Sensor Name")
	}
	voltageSensorObj.Name = Name
	voltageSensorObj.AdminState = voltageSensorCfgEnt.AdminState
	voltageSensorObj.HigherAlarmThreshold = voltageSensorCfgEnt.HigherAlarmThreshold
	voltageSensorObj.HigherWarningThreshold = voltageSensorCfgEnt.HigherWarningThreshold
	voltageSensorObj.LowerAlarmThreshold = voltageSensorCfgEnt.LowerAlarmThreshold
	voltageSensorObj.LowerWarningThreshold = voltageSensorCfgEnt.LowerWarningThreshold
	voltageSensorObj.PMClassAAdminState = voltageSensorCfgEnt.PMClassAAdminState
	voltageSensorObj.PMClassBAdminState = voltageSensorCfgEnt.PMClassBAdminState
	voltageSensorObj.PMClassCAdminState = voltageSensorCfgEnt.PMClassCAdminState
	sMgr.voltageConfigMutex.RUnlock()
	return &voltageSensorObj, nil
}

func (sMgr *SensorManager) GetBulkVoltageSensorConfig(fromIdx int, cnt int) (*objects.VoltageSensorConfigGetInfo, error) {
	var retObj objects.VoltageSensorConfigGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.voltageSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.voltageSensorList) {
		retObj.EndIdx = len(sMgr.voltageSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.voltageSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		voltageSensorName := sMgr.voltageSensorList[idx]
		obj, err := sMgr.GetVoltageSensorConfig(voltageSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the voltage state for voltage sensor:", voltageSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func genVoltageSensorUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0
	if attrset == nil {
		mask = objects.VOLTAGE_SENSOR_UPDATE_ADMIN_STATE |
			objects.VOLTAGE_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD |
			objects.VOLTAGE_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD |
			objects.VOLTAGE_SENSOR_UPDATE_LOWER_WARN_THRESHOLD |
			objects.VOLTAGE_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD |
			objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE |
			objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE |
			objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE
	} else {
		for idx, val := range attrset {
			if true == val {
				switch idx {
				case 0:
					//ObjKey Voltage Sensor Name
				case 1:
					mask |= objects.VOLTAGE_SENSOR_UPDATE_ADMIN_STATE
				case 2:
					mask |= objects.VOLTAGE_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD
				case 3:
					mask |= objects.VOLTAGE_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD
				case 4:
					mask |= objects.VOLTAGE_SENSOR_UPDATE_LOWER_WARN_THRESHOLD
				case 5:
					mask |= objects.VOLTAGE_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD
				case 6:
					mask |= objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE
				case 7:
					mask |= objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE
				case 8:
					mask |= objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE
				}
			}
		}
	}
	return mask
}

func (sMgr *SensorManager) UpdateVoltageSensorConfig(oldCfg *objects.VoltageSensorConfig, newCfg *objects.VoltageSensorConfig, attrset []bool) (bool, error) {
	if sMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	sMgr.voltageConfigMutex.Lock()
	voltageSensorCfgEnt, exist := sMgr.voltageConfigDB[newCfg.Name]
	if !exist {
		sMgr.voltageConfigMutex.Unlock()
		return false, errors.New("Invalid VoltageSensor Name")
	}
	mask := genVoltageSensorUpdateMask(attrset)
	if mask&objects.VOLTAGE_SENSOR_UPDATE_ADMIN_STATE == objects.VOLTAGE_SENSOR_UPDATE_ADMIN_STATE {
		if newCfg.AdminState != "Enable" && newCfg.AdminState != "Disable" {
			sMgr.voltageConfigMutex.Unlock()
			return false, errors.New("Invalid AdminState Value")
		}
		voltageSensorCfgEnt.AdminState = newCfg.AdminState
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD == objects.VOLTAGE_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD {
		voltageSensorCfgEnt.HigherAlarmThreshold = newCfg.HigherAlarmThreshold
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD == objects.VOLTAGE_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD {
		voltageSensorCfgEnt.HigherWarningThreshold = newCfg.HigherWarningThreshold
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD == objects.VOLTAGE_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD {
		voltageSensorCfgEnt.LowerAlarmThreshold = newCfg.LowerAlarmThreshold
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_LOWER_WARN_THRESHOLD == objects.VOLTAGE_SENSOR_UPDATE_LOWER_WARN_THRESHOLD {
		voltageSensorCfgEnt.LowerWarningThreshold = newCfg.LowerWarningThreshold
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE == objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE {
		if newCfg.PMClassAAdminState != "Enable" && newCfg.PMClassAAdminState != "Disable" {
			sMgr.voltageConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassAAdminState Value")
		}
		voltageSensorCfgEnt.PMClassAAdminState = newCfg.PMClassAAdminState
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE == objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE {
		if newCfg.PMClassBAdminState != "Enable" && newCfg.PMClassBAdminState != "Disable" {
			sMgr.voltageConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassBAdminState Value")
		}
		voltageSensorCfgEnt.PMClassBAdminState = newCfg.PMClassBAdminState
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE == objects.VOLTAGE_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE {
		if newCfg.PMClassCAdminState != "Enable" && newCfg.PMClassCAdminState != "Disable" {
			sMgr.voltageConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassCAdminState Value")
		}
		voltageSensorCfgEnt.PMClassCAdminState = newCfg.PMClassCAdminState
	}
	if !(voltageSensorCfgEnt.HigherAlarmThreshold >= voltageSensorCfgEnt.HigherWarningThreshold &&
		voltageSensorCfgEnt.HigherWarningThreshold > voltageSensorCfgEnt.LowerWarningThreshold &&
		voltageSensorCfgEnt.LowerWarningThreshold >= voltageSensorCfgEnt.LowerAlarmThreshold) {
		sMgr.voltageConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, Please verify the thresholds")
	}
	if sMgr.voltageConfigDB[newCfg.Name].AdminState != voltageSensorCfgEnt.AdminState {
		if voltageSensorCfgEnt.AdminState == "Disable" {
			//Clear all alarms
			sMgr.clearExistingVoltageSensorFaults(newCfg.Name)
			sMgr.logger.Info("Clear all the existing faults")
		}
	}
	if sMgr.voltageConfigDB[newCfg.Name].PMClassAAdminState != voltageSensorCfgEnt.PMClassAAdminState {
		if voltageSensorCfgEnt.PMClassAAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.voltageClassAPMMutex.Lock()
			sMgr.voltageSensorClassAPM[newCfg.Name].FlushRingBuffer()
			sMgr.voltageClassAPMMutex.Unlock()
			sMgr.logger.Info("Flush Class A PM Ring Buffer")
		}
	}
	if sMgr.voltageConfigDB[newCfg.Name].PMClassBAdminState != voltageSensorCfgEnt.PMClassBAdminState {
		if voltageSensorCfgEnt.PMClassBAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.voltageClassBPMMutex.Lock()
			sMgr.voltageSensorClassBPM[newCfg.Name].FlushRingBuffer()
			sMgr.voltageClassBPMMutex.Unlock()
			sMgr.logger.Info("Flush Class B PM Ring Buffer")
		}
	}
	if sMgr.voltageConfigDB[newCfg.Name].PMClassCAdminState != voltageSensorCfgEnt.PMClassCAdminState {
		if voltageSensorCfgEnt.PMClassCAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.voltageClassCPMMutex.Lock()
			sMgr.voltageSensorClassCPM[newCfg.Name].FlushRingBuffer()
			sMgr.voltageClassCPMMutex.Unlock()
			sMgr.logger.Info("Flush Class C PM Ring Buffer")
		}
	}
	sMgr.voltageConfigDB[newCfg.Name] = voltageSensorCfgEnt
	sMgr.voltageConfigMutex.Unlock()

	return true, nil
}

func (sMgr *SensorManager) GetPowerConverterSensorState(Name string) (*objects.PowerConverterSensorState, error) {
	var powerConverterSensorObj objects.PowerConverterSensorState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.powerConverterConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid PowerConverter Sensor Name")
	}

	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.SensorStateMutex.Unlock()
		return nil, err
	}
	powerConverterSensorState, _ := sMgr.SensorState.PowerConverterSensor[Name]
	powerConverterSensorObj.Name = Name
	powerConverterSensorObj.CurrentPower = powerConverterSensorState.Value
	sMgr.SensorStateMutex.Unlock()
	return &powerConverterSensorObj, err
}

func (sMgr *SensorManager) GetBulkPowerConverterSensorState(fromIdx int, cnt int) (*objects.PowerConverterSensorStateGetInfo, error) {
	var retObj objects.PowerConverterSensorStateGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.powerConverterSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.powerConverterSensorList) {
		retObj.EndIdx = len(sMgr.powerConverterSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.powerConverterSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		powerConverterSensorName := sMgr.powerConverterSensorList[idx]
		obj, err := sMgr.GetPowerConverterSensorState(powerConverterSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the powerConverter state for powerConverter Sensor:", powerConverterSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (sMgr *SensorManager) GetPowerConverterSensorConfig(Name string) (*objects.PowerConverterSensorConfig, error) {
	var powerConverterSensorObj objects.PowerConverterSensorConfig
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.powerConverterConfigMutex.RLock()
	powerConverterSensorCfgEnt, exist := sMgr.powerConverterConfigDB[Name]
	if !exist {
		sMgr.powerConverterConfigMutex.RUnlock()
		return nil, errors.New("Invalid PowerConverter Sensor Name")
	}
	powerConverterSensorObj.Name = Name
	powerConverterSensorObj.AdminState = powerConverterSensorCfgEnt.AdminState
	powerConverterSensorObj.HigherAlarmThreshold = powerConverterSensorCfgEnt.HigherAlarmThreshold
	powerConverterSensorObj.HigherWarningThreshold = powerConverterSensorCfgEnt.HigherWarningThreshold
	powerConverterSensorObj.LowerAlarmThreshold = powerConverterSensorCfgEnt.LowerAlarmThreshold
	powerConverterSensorObj.LowerWarningThreshold = powerConverterSensorCfgEnt.LowerWarningThreshold
	powerConverterSensorObj.PMClassAAdminState = powerConverterSensorCfgEnt.PMClassAAdminState
	powerConverterSensorObj.PMClassBAdminState = powerConverterSensorCfgEnt.PMClassBAdminState
	powerConverterSensorObj.PMClassCAdminState = powerConverterSensorCfgEnt.PMClassCAdminState
	sMgr.powerConverterConfigMutex.RUnlock()
	return &powerConverterSensorObj, nil
}

func (sMgr *SensorManager) GetBulkPowerConverterSensorConfig(fromIdx int, cnt int) (*objects.PowerConverterSensorConfigGetInfo, error) {
	var retObj objects.PowerConverterSensorConfigGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.powerConverterSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.powerConverterSensorList) {
		retObj.EndIdx = len(sMgr.powerConverterSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.powerConverterSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		powerConverterSensorName := sMgr.powerConverterSensorList[idx]
		obj, err := sMgr.GetPowerConverterSensorConfig(powerConverterSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the powerConverter state for powerConverter sensor:", powerConverterSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func genPowerConverterSensorUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0

	if attrset == nil {
		mask = objects.POWER_CONVERTER_SENSOR_UPDATE_ADMIN_STATE |
			objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD |
			objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD |
			objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_WARN_THRESHOLD |
			objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD |
			objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE |
			objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE |
			objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE
	} else {
		for idx, val := range attrset {
			if true == val {
				switch idx {
				case 0:
					//ObjKey PowerConverter Sensor Name
				case 1:
					mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_ADMIN_STATE
				case 2:
					mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD
				case 3:
					mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD
				case 4:
					mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_WARN_THRESHOLD
				case 5:
					mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD
				case 6:
					mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE
				case 7:
					mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE
				case 8:
					mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE
				}
			}
		}
	}
	return mask
}

func (sMgr *SensorManager) UpdatePowerConverterSensorConfig(oldCfg *objects.PowerConverterSensorConfig, newCfg *objects.PowerConverterSensorConfig, attrset []bool) (bool, error) {
	if sMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	sMgr.powerConverterConfigMutex.Lock()
	powerConverterSensorCfgEnt, exist := sMgr.powerConverterConfigDB[newCfg.Name]
	if !exist {
		sMgr.powerConverterConfigMutex.Unlock()
		return false, errors.New("Invalid PowerConverterSensor Name")
	}
	mask := genPowerConverterSensorUpdateMask(attrset)
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_ADMIN_STATE == objects.POWER_CONVERTER_SENSOR_UPDATE_ADMIN_STATE {
		if newCfg.AdminState != "Enable" && newCfg.AdminState != "Disable" {
			sMgr.powerConverterConfigMutex.Unlock()
			return false, errors.New("Invalid AdminState Value")
		}
		powerConverterSensorCfgEnt.AdminState = newCfg.AdminState
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD == objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD {
		powerConverterSensorCfgEnt.HigherAlarmThreshold = newCfg.HigherAlarmThreshold
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD == objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD {
		powerConverterSensorCfgEnt.HigherWarningThreshold = newCfg.HigherWarningThreshold
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD == objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD {
		powerConverterSensorCfgEnt.LowerAlarmThreshold = newCfg.LowerAlarmThreshold
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_WARN_THRESHOLD == objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_WARN_THRESHOLD {
		powerConverterSensorCfgEnt.LowerWarningThreshold = newCfg.LowerWarningThreshold
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE == objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE {
		if newCfg.PMClassAAdminState != "Enable" && newCfg.PMClassAAdminState != "Disable" {
			sMgr.powerConverterConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassAAdminState Value")
		}
		powerConverterSensorCfgEnt.PMClassAAdminState = newCfg.PMClassAAdminState
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE == objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE {
		if newCfg.PMClassBAdminState != "Enable" && newCfg.PMClassBAdminState != "Disable" {
			sMgr.powerConverterConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassBAdminState Value")
		}
		powerConverterSensorCfgEnt.PMClassBAdminState = newCfg.PMClassBAdminState
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE == objects.POWER_CONVERTER_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE {
		if newCfg.PMClassCAdminState != "Enable" && newCfg.PMClassCAdminState != "Disable" {
			sMgr.powerConverterConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassCAdminState Value")
		}
		powerConverterSensorCfgEnt.PMClassCAdminState = newCfg.PMClassCAdminState
	}
	if !(powerConverterSensorCfgEnt.HigherAlarmThreshold >= powerConverterSensorCfgEnt.HigherWarningThreshold &&
		powerConverterSensorCfgEnt.HigherWarningThreshold > powerConverterSensorCfgEnt.LowerWarningThreshold &&
		powerConverterSensorCfgEnt.LowerWarningThreshold >= powerConverterSensorCfgEnt.LowerAlarmThreshold) {
		sMgr.powerConverterConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, Please verify the thresholds")
	}
	if sMgr.powerConverterConfigDB[newCfg.Name].AdminState != powerConverterSensorCfgEnt.AdminState {
		if powerConverterSensorCfgEnt.AdminState == "Disable" {
			//Clear all alarms
			sMgr.clearExistingPowerConverterSensorFaults(newCfg.Name)
			sMgr.logger.Info("Clear all the existing faults")
		}
	}
	if sMgr.powerConverterConfigDB[newCfg.Name].PMClassAAdminState != powerConverterSensorCfgEnt.PMClassAAdminState {
		if powerConverterSensorCfgEnt.PMClassAAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.powerConverterClassAPMMutex.Lock()
			sMgr.powerConverterSensorClassAPM[newCfg.Name].FlushRingBuffer()
			sMgr.powerConverterClassAPMMutex.Unlock()
			sMgr.logger.Info("Flush Class A PM Ring Buffer")
		}
	}
	if sMgr.powerConverterConfigDB[newCfg.Name].PMClassBAdminState != powerConverterSensorCfgEnt.PMClassBAdminState {
		if powerConverterSensorCfgEnt.PMClassBAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.powerConverterClassBPMMutex.Lock()
			sMgr.powerConverterSensorClassBPM[newCfg.Name].FlushRingBuffer()
			sMgr.powerConverterClassBPMMutex.Unlock()
			sMgr.logger.Info("Flush Class B PM Ring Buffer")
		}
	}
	if sMgr.powerConverterConfigDB[newCfg.Name].PMClassCAdminState != powerConverterSensorCfgEnt.PMClassCAdminState {
		if powerConverterSensorCfgEnt.PMClassCAdminState == "Disable" {
			//Flush PM Ringbuffer
			sMgr.powerConverterClassCPMMutex.Lock()
			sMgr.powerConverterSensorClassCPM[newCfg.Name].FlushRingBuffer()
			sMgr.powerConverterClassCPMMutex.Unlock()
			sMgr.logger.Info("Flush Class C PM Ring Buffer")
		}
	}
	sMgr.powerConverterConfigDB[newCfg.Name] = powerConverterSensorCfgEnt
	sMgr.powerConverterConfigMutex.Unlock()

	return true, nil
}

func (sMgr *SensorManager) InitSensorPM() {
	for fanSensorName, _ := range sMgr.fanConfigDB {
		sMgr.fanSensorClassAPM[fanSensorName] = new(ringBuffer.RingBuffer)
		sMgr.fanSensorClassAPM[fanSensorName].SetRingBufferCapacity(sensorClassABufSize)
		sMgr.fanSensorClassBPM[fanSensorName] = new(ringBuffer.RingBuffer)
		sMgr.fanSensorClassBPM[fanSensorName].SetRingBufferCapacity(sensorClassBBufSize)
		sMgr.fanSensorClassCPM[fanSensorName] = new(ringBuffer.RingBuffer)
		sMgr.fanSensorClassCPM[fanSensorName].SetRingBufferCapacity(sensorClassCBufSize)
	}
	for tempSensorName, _ := range sMgr.tempConfigDB {
		sMgr.tempSensorClassAPM[tempSensorName] = new(ringBuffer.RingBuffer)
		sMgr.tempSensorClassAPM[tempSensorName].SetRingBufferCapacity(sensorClassABufSize)
		sMgr.tempSensorClassBPM[tempSensorName] = new(ringBuffer.RingBuffer)
		sMgr.tempSensorClassBPM[tempSensorName].SetRingBufferCapacity(sensorClassBBufSize)
		sMgr.tempSensorClassCPM[tempSensorName] = new(ringBuffer.RingBuffer)
		sMgr.tempSensorClassCPM[tempSensorName].SetRingBufferCapacity(sensorClassCBufSize)
	}

	for voltageSensorName, _ := range sMgr.voltageConfigDB {
		sMgr.voltageSensorClassAPM[voltageSensorName] = new(ringBuffer.RingBuffer)
		sMgr.voltageSensorClassAPM[voltageSensorName].SetRingBufferCapacity(sensorClassABufSize)
		sMgr.voltageSensorClassBPM[voltageSensorName] = new(ringBuffer.RingBuffer)
		sMgr.voltageSensorClassBPM[voltageSensorName].SetRingBufferCapacity(sensorClassBBufSize)
		sMgr.voltageSensorClassCPM[voltageSensorName] = new(ringBuffer.RingBuffer)
		sMgr.voltageSensorClassCPM[voltageSensorName].SetRingBufferCapacity(sensorClassCBufSize)
	}
	for powerConverterSensorName, _ := range sMgr.powerConverterConfigDB {
		sMgr.powerConverterSensorClassAPM[powerConverterSensorName] = new(ringBuffer.RingBuffer)
		sMgr.powerConverterSensorClassAPM[powerConverterSensorName].SetRingBufferCapacity(sensorClassABufSize)
		sMgr.powerConverterSensorClassBPM[powerConverterSensorName] = new(ringBuffer.RingBuffer)
		sMgr.powerConverterSensorClassBPM[powerConverterSensorName].SetRingBufferCapacity(sensorClassBBufSize)
		sMgr.powerConverterSensorClassCPM[powerConverterSensorName] = new(ringBuffer.RingBuffer)
		sMgr.powerConverterSensorClassCPM[powerConverterSensorName].SetRingBufferCapacity(sensorClassCBufSize)
	}
}

func (sMgr *SensorManager) StartSensorPM() {
	sMgr.InitSensorPM()
	sMgr.StartSensorPMClass("Class-A")
	sMgr.StartSensorPMClass("Class-B")
	sMgr.StartSensorPMClass("Class-C")
}

func (sMgr *SensorManager) clearExistingFanSensorFaults(name string) {
	var evts []events.EventId
	eventKey := events.FanSensorKey{
		Name: name,
	}
	txEvent := eventUtils.TxEvent{
		Key:            eventKey,
		AdditionalInfo: "Fault cleared because of Admin Disable",
		AdditionalData: nil,
	}
	sMgr.fanMsgStatusMutex.RLock()
	prevEventStatus := sMgr.fanMsgStatus[name]
	sMgr.fanMsgStatusMutex.RUnlock()
	if prevEventStatus.SentHigherAlarm == true {
		prevEventStatus.SentHigherAlarm = false
		evts = append(evts, events.FanHigherTCAAlarmClear)
	}
	if prevEventStatus.SentHigherWarn == true {
		prevEventStatus.SentLowerWarn = false
		evts = append(evts, events.FanHigherTCAWarnClear)
	}
	if prevEventStatus.SentLowerWarn == true {
		prevEventStatus.SentLowerWarn = false
		evts = append(evts, events.FanLowerTCAWarnClear)
	}
	if prevEventStatus.SentLowerAlarm == true {
		prevEventStatus.SentLowerAlarm = false
		evts = append(evts, events.FanLowerTCAAlarmClear)
	}
	sMgr.fanMsgStatusMutex.Lock()
	sMgr.fanMsgStatus[name] = prevEventStatus
	sMgr.fanMsgStatusMutex.Unlock()
	for _, evt := range evts {
		txEvent.EventId = evt
		txEvt := txEvent
		err := eventUtils.PublishEvents(&txEvt)
		if err != nil {
			sMgr.logger.Err("Error publish events")
		}
	}
}

func (sMgr *SensorManager) ProcessFanSensorPM(sensorState *pluginCommon.SensorState, class string) {
	sMgr.fanConfigMutex.RLock()
	for fanSensorName, fanSensorCfgEnt := range sMgr.fanConfigDB {
		fanSensorState, _ := sMgr.SensorState.FanSensor[fanSensorName]
		fanSensorPMData := objects.FanSensorPMData{
			TimeStamp: time.Now().String(),
			Value:     fanSensorState.Value,
		}
		if fanSensorCfgEnt.AdminState == "Enable" {
			eventKey := events.FanSensorKey{
				Name: fanSensorName,
			}
			eventData := FanSensorEventData{
				Value: fanSensorState.Value,
			}
			txEvent := eventUtils.TxEvent{
				Key:            eventKey,
				AdditionalInfo: "",
				AdditionalData: eventData,
			}
			var evts []events.EventId
			var curEvents EventStatus
			sMgr.fanMsgStatusMutex.RLock()
			prevEventStatus := sMgr.fanMsgStatus[fanSensorName]
			sMgr.fanMsgStatusMutex.RUnlock()
			if fanSensorState.Value >= fanSensorCfgEnt.HigherAlarmThreshold {
				curEvents.SentHigherAlarm = true
			}
			if fanSensorState.Value >= fanSensorCfgEnt.HigherWarningThreshold {
				curEvents.SentHigherWarn = true
			}
			if fanSensorState.Value <= fanSensorCfgEnt.LowerWarningThreshold {
				curEvents.SentLowerWarn = true
			}
			if fanSensorState.Value <= fanSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerAlarm = true
			}
			if prevEventStatus.SentHigherAlarm != curEvents.SentHigherAlarm {
				if curEvents.SentHigherAlarm == true {
					evts = append(evts, events.FanHigherTCAAlarm)
				} else {
					evts = append(evts, events.FanHigherTCAAlarmClear)
				}
			}
			if prevEventStatus.SentHigherWarn != curEvents.SentHigherWarn {
				if curEvents.SentHigherWarn == true {
					evts = append(evts, events.FanHigherTCAWarn)
				} else {
					evts = append(evts, events.FanHigherTCAWarnClear)
				}
			}
			if prevEventStatus.SentLowerAlarm != curEvents.SentLowerAlarm {
				if curEvents.SentLowerAlarm == true {
					evts = append(evts, events.FanLowerTCAAlarm)
				} else {
					evts = append(evts, events.FanLowerTCAAlarmClear)
				}
			}
			if prevEventStatus.SentLowerWarn != curEvents.SentLowerWarn {
				if curEvents.SentLowerWarn == true {
					evts = append(evts, events.FanLowerTCAWarn)
				} else {
					evts = append(evts, events.FanLowerTCAWarnClear)
				}
			}
			if prevEventStatus != curEvents {
				sMgr.fanMsgStatusMutex.Lock()
				sMgr.fanMsgStatus[fanSensorName] = curEvents
				sMgr.fanMsgStatusMutex.Unlock()
			}
			for _, evt := range evts {
				txEvent.EventId = evt
				txEvt := txEvent
				err := eventUtils.PublishEvents(&txEvt)
				if err != nil {
					sMgr.logger.Err("Error publish events")
				}
			}
		}
		switch class {
		case "Class-A":
			if fanSensorCfgEnt.PMClassAAdminState == "Enable" {
				sMgr.fanClassAPMMutex.Lock()
				sMgr.fanSensorClassAPM[fanSensorName].InsertIntoRingBuffer(fanSensorPMData)
				sMgr.fanClassAPMMutex.Unlock()
			}
		case "Class-B":
			if fanSensorCfgEnt.PMClassBAdminState == "Enable" {
				sMgr.fanClassBPMMutex.Lock()
				sMgr.fanSensorClassBPM[fanSensorName].InsertIntoRingBuffer(fanSensorPMData)
				sMgr.fanClassBPMMutex.Unlock()
			}
		case "Class-C":
			if fanSensorCfgEnt.PMClassCAdminState == "Enable" {
				sMgr.fanClassCPMMutex.Lock()
				sMgr.fanSensorClassCPM[fanSensorName].InsertIntoRingBuffer(fanSensorPMData)
				sMgr.fanClassCPMMutex.Unlock()
			}
		}
	}
	sMgr.fanConfigMutex.RUnlock()
}

func (sMgr *SensorManager) clearExistingTempSensorFaults(name string) {
	var evts []events.EventId
	eventKey := events.TemperatureSensorKey{
		Name: name,
	}
	txEvent := eventUtils.TxEvent{
		Key:            eventKey,
		AdditionalInfo: "Fault cleared because of Admin Disable",
		AdditionalData: nil,
	}
	sMgr.tempMsgStatusMutex.RLock()
	prevEventStatus := sMgr.tempMsgStatus[name]
	sMgr.tempMsgStatusMutex.RUnlock()
	if prevEventStatus.SentHigherAlarm == true {
		prevEventStatus.SentHigherAlarm = false
		evts = append(evts, events.TemperatureHigherTCAAlarmClear)
	}
	if prevEventStatus.SentHigherWarn == true {
		prevEventStatus.SentLowerWarn = false
		evts = append(evts, events.TemperatureHigherTCAWarnClear)
	}
	if prevEventStatus.SentLowerWarn == true {
		prevEventStatus.SentLowerWarn = false
		evts = append(evts, events.TemperatureLowerTCAWarnClear)
	}
	if prevEventStatus.SentLowerAlarm == true {
		prevEventStatus.SentLowerAlarm = false
		evts = append(evts, events.TemperatureLowerTCAAlarmClear)
	}
	sMgr.tempMsgStatusMutex.Lock()
	sMgr.tempMsgStatus[name] = prevEventStatus
	sMgr.tempMsgStatusMutex.Unlock()
	for _, evt := range evts {
		txEvent.EventId = evt
		txEvt := txEvent
		err := eventUtils.PublishEvents(&txEvt)
		if err != nil {
			sMgr.logger.Err("Error publish events")
		}
	}
}

func (sMgr *SensorManager) ProcessTempSensorPM(sensorState *pluginCommon.SensorState, class string) {
	sMgr.tempConfigMutex.RLock()
	for tempSensorName, tempSensorCfgEnt := range sMgr.tempConfigDB {
		tempSensorState, _ := sMgr.SensorState.TemperatureSensor[tempSensorName]
		tempSensorPMData := objects.TemperatureSensorPMData{
			TimeStamp: time.Now().String(),
			Value:     tempSensorState.Value,
		}
		if tempSensorCfgEnt.AdminState == "Enable" {
			eventKey := events.TemperatureSensorKey{
				Name: tempSensorName,
			}
			eventData := TempSensorEventData{
				Value: tempSensorState.Value,
			}
			txEvent := eventUtils.TxEvent{
				Key:            eventKey,
				AdditionalInfo: "",
				AdditionalData: eventData,
			}
			var evts []events.EventId
			var curEvents EventStatus
			sMgr.tempMsgStatusMutex.RLock()
			prevEventStatus := sMgr.tempMsgStatus[tempSensorName]
			sMgr.tempMsgStatusMutex.RUnlock()
			if tempSensorState.Value >= tempSensorCfgEnt.HigherAlarmThreshold {
				curEvents.SentHigherAlarm = true
			}
			if tempSensorState.Value >= tempSensorCfgEnt.HigherWarningThreshold {
				curEvents.SentHigherWarn = true
			}
			if tempSensorState.Value <= tempSensorCfgEnt.LowerWarningThreshold {
				curEvents.SentLowerWarn = true
			}
			if tempSensorState.Value <= tempSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerAlarm = true
			}
			if prevEventStatus.SentHigherAlarm != curEvents.SentHigherAlarm {
				if curEvents.SentHigherAlarm == true {
					evts = append(evts, events.TemperatureHigherTCAAlarm)
				} else {
					evts = append(evts, events.TemperatureHigherTCAAlarmClear)
				}
			}
			if prevEventStatus.SentHigherWarn != curEvents.SentHigherWarn {
				if curEvents.SentHigherWarn == true {
					evts = append(evts, events.TemperatureHigherTCAWarn)
				} else {
					evts = append(evts, events.TemperatureHigherTCAWarnClear)
				}
			}
			if prevEventStatus.SentLowerAlarm != curEvents.SentLowerAlarm {
				if curEvents.SentLowerAlarm == true {
					evts = append(evts, events.TemperatureLowerTCAAlarm)
				} else {
					evts = append(evts, events.TemperatureLowerTCAAlarmClear)
				}
			}
			if prevEventStatus.SentLowerWarn != curEvents.SentLowerWarn {
				if curEvents.SentLowerWarn == true {
					evts = append(evts, events.TemperatureLowerTCAWarn)
				} else {
					evts = append(evts, events.TemperatureLowerTCAWarnClear)
				}
			}
			if prevEventStatus != curEvents {
				sMgr.tempMsgStatusMutex.Lock()
				sMgr.tempMsgStatus[tempSensorName] = curEvents
				sMgr.tempMsgStatusMutex.Unlock()
			}
			for _, evt := range evts {
				txEvent.EventId = evt
				txEvt := txEvent
				err := eventUtils.PublishEvents(&txEvt)
				if err != nil {
					sMgr.logger.Err("Error publish events")
				}
			}

		}

		switch class {
		case "Class-A":
			if tempSensorCfgEnt.PMClassAAdminState == "Enable" {
				sMgr.tempClassAPMMutex.Lock()
				sMgr.tempSensorClassAPM[tempSensorName].InsertIntoRingBuffer(tempSensorPMData)
				sMgr.tempClassAPMMutex.Unlock()
			}
		case "Class-B":
			if tempSensorCfgEnt.PMClassBAdminState == "Enable" {
				sMgr.tempClassBPMMutex.Lock()
				sMgr.tempSensorClassBPM[tempSensorName].InsertIntoRingBuffer(tempSensorPMData)
				sMgr.tempClassBPMMutex.Unlock()
			}
		case "Class-C":
			if tempSensorCfgEnt.PMClassCAdminState == "Enable" {
				sMgr.tempClassCPMMutex.Lock()
				sMgr.tempSensorClassCPM[tempSensorName].InsertIntoRingBuffer(tempSensorPMData)
				sMgr.tempClassCPMMutex.Unlock()
			}
		}
	}
	sMgr.tempConfigMutex.RUnlock()
}

func (sMgr *SensorManager) clearExistingVoltageSensorFaults(name string) {
	var evts []events.EventId
	eventKey := events.VoltageSensorKey{
		Name: name,
	}
	txEvent := eventUtils.TxEvent{
		Key:            eventKey,
		AdditionalInfo: "Fault cleared because of Admin Disable",
		AdditionalData: nil,
	}
	sMgr.voltageMsgStatusMutex.RLock()
	prevEventStatus := sMgr.voltageMsgStatus[name]
	sMgr.voltageMsgStatusMutex.RUnlock()
	if prevEventStatus.SentHigherAlarm == true {
		prevEventStatus.SentHigherAlarm = false
		evts = append(evts, events.VoltageHigherTCAAlarmClear)
	}
	if prevEventStatus.SentHigherWarn == true {
		prevEventStatus.SentLowerWarn = false
		evts = append(evts, events.VoltageHigherTCAWarnClear)
	}
	if prevEventStatus.SentLowerWarn == true {
		prevEventStatus.SentLowerWarn = false
		evts = append(evts, events.VoltageLowerTCAWarnClear)
	}
	if prevEventStatus.SentLowerAlarm == true {
		prevEventStatus.SentLowerAlarm = false
		evts = append(evts, events.VoltageLowerTCAAlarmClear)
	}
	sMgr.voltageMsgStatusMutex.Lock()
	sMgr.voltageMsgStatus[name] = prevEventStatus
	sMgr.voltageMsgStatusMutex.Unlock()
	for _, evt := range evts {
		txEvent.EventId = evt
		txEvt := txEvent
		err := eventUtils.PublishEvents(&txEvt)
		if err != nil {
			sMgr.logger.Err("Error publish events")
		}
	}
}

func (sMgr *SensorManager) ProcessVoltageSensorPM(sensorState *pluginCommon.SensorState, class string) {
	sMgr.voltageConfigMutex.RLock()
	for voltageSensorName, voltageSensorCfgEnt := range sMgr.voltageConfigDB {
		voltageSensorState, _ := sMgr.SensorState.VoltageSensor[voltageSensorName]
		voltageSensorPMData := objects.VoltageSensorPMData{
			TimeStamp: time.Now().String(),
			Value:     voltageSensorState.Value,
		}
		if voltageSensorCfgEnt.AdminState == "Enable" {
			eventKey := events.VoltageSensorKey{
				Name: voltageSensorName,
			}
			eventData := VoltageSensorEventData{
				Value: voltageSensorState.Value,
			}
			txEvent := eventUtils.TxEvent{
				Key:            eventKey,
				AdditionalInfo: "",
				AdditionalData: eventData,
			}
			var evts []events.EventId
			var curEvents EventStatus
			sMgr.voltageMsgStatusMutex.RLock()
			prevEventStatus := sMgr.voltageMsgStatus[voltageSensorName]
			sMgr.voltageMsgStatusMutex.RUnlock()
			if voltageSensorState.Value >= voltageSensorCfgEnt.HigherAlarmThreshold {
				curEvents.SentHigherAlarm = true
			}
			if voltageSensorState.Value >= voltageSensorCfgEnt.HigherWarningThreshold {
				curEvents.SentHigherWarn = true
			}
			if voltageSensorState.Value <= voltageSensorCfgEnt.LowerWarningThreshold {
				curEvents.SentLowerWarn = true
			}
			if voltageSensorState.Value <= voltageSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerAlarm = true
			}
			if prevEventStatus.SentHigherAlarm != curEvents.SentHigherAlarm {
				if curEvents.SentHigherAlarm == true {
					evts = append(evts, events.VoltageHigherTCAAlarm)
				} else {
					evts = append(evts, events.VoltageHigherTCAAlarmClear)
				}
			}
			if prevEventStatus.SentHigherWarn != curEvents.SentHigherWarn {
				if curEvents.SentHigherWarn == true {
					evts = append(evts, events.VoltageHigherTCAWarn)
				} else {
					evts = append(evts, events.VoltageHigherTCAWarnClear)
				}
			}
			if prevEventStatus.SentLowerAlarm != curEvents.SentLowerAlarm {
				if curEvents.SentLowerAlarm == true {
					evts = append(evts, events.VoltageLowerTCAAlarm)
				} else {
					evts = append(evts, events.VoltageLowerTCAAlarmClear)
				}
			}
			if prevEventStatus.SentLowerWarn != curEvents.SentLowerWarn {
				if curEvents.SentLowerWarn == true {
					evts = append(evts, events.VoltageLowerTCAWarn)
				} else {
					evts = append(evts, events.VoltageLowerTCAWarnClear)
				}
			}
			if prevEventStatus != curEvents {
				sMgr.voltageMsgStatusMutex.Lock()
				sMgr.voltageMsgStatus[voltageSensorName] = curEvents
				sMgr.voltageMsgStatusMutex.Unlock()
			}
			for _, evt := range evts {
				txEvent.EventId = evt
				txEvt := txEvent
				err := eventUtils.PublishEvents(&txEvt)
				if err != nil {
					sMgr.logger.Err("Error publish events")
				}
			}

		}

		switch class {
		case "Class-A":
			if voltageSensorCfgEnt.PMClassAAdminState == "Enable" {
				sMgr.voltageClassAPMMutex.Lock()
				sMgr.voltageSensorClassAPM[voltageSensorName].InsertIntoRingBuffer(voltageSensorPMData)
				sMgr.voltageClassAPMMutex.Unlock()
			}
		case "Class-B":
			if voltageSensorCfgEnt.PMClassBAdminState == "Enable" {
				sMgr.voltageClassBPMMutex.Lock()
				sMgr.voltageSensorClassBPM[voltageSensorName].InsertIntoRingBuffer(voltageSensorPMData)
				sMgr.voltageClassBPMMutex.Unlock()
			}
		case "Class-C":
			if voltageSensorCfgEnt.PMClassCAdminState == "Enable" {
				sMgr.voltageClassCPMMutex.Lock()
				sMgr.voltageSensorClassCPM[voltageSensorName].InsertIntoRingBuffer(voltageSensorPMData)
				sMgr.voltageClassCPMMutex.Unlock()
			}
		}
	}
	sMgr.voltageConfigMutex.RUnlock()
}

func (sMgr *SensorManager) clearExistingPowerConverterSensorFaults(name string) {
	var evts []events.EventId
	eventKey := events.PowerConverterSensorKey{
		Name: name,
	}
	txEvent := eventUtils.TxEvent{
		Key:            eventKey,
		AdditionalInfo: "Fault cleared because of Admin Disable",
		AdditionalData: nil,
	}
	sMgr.powerConverterMsgStatusMutex.RLock()
	prevEventStatus := sMgr.powerConverterMsgStatus[name]
	sMgr.powerConverterMsgStatusMutex.RUnlock()
	if prevEventStatus.SentHigherAlarm == true {
		prevEventStatus.SentHigherAlarm = false
		evts = append(evts, events.PowerConverterHigherTCAAlarmClear)
	}
	if prevEventStatus.SentHigherWarn == true {
		prevEventStatus.SentLowerWarn = false
		evts = append(evts, events.PowerConverterHigherTCAWarnClear)
	}
	if prevEventStatus.SentLowerWarn == true {
		prevEventStatus.SentLowerWarn = false
		evts = append(evts, events.PowerConverterLowerTCAWarnClear)
	}
	if prevEventStatus.SentLowerAlarm == true {
		prevEventStatus.SentLowerAlarm = false
		evts = append(evts, events.PowerConverterLowerTCAAlarmClear)
	}
	sMgr.powerConverterMsgStatusMutex.Lock()
	sMgr.powerConverterMsgStatus[name] = prevEventStatus
	sMgr.powerConverterMsgStatusMutex.Unlock()
	for _, evt := range evts {
		txEvent.EventId = evt
		txEvt := txEvent
		err := eventUtils.PublishEvents(&txEvt)
		if err != nil {
			sMgr.logger.Err("Error publish events")
		}
	}
}

func (sMgr *SensorManager) ProcessPowerConverterSensorPM(sensorState *pluginCommon.SensorState, class string) {
	sMgr.powerConverterConfigMutex.RLock()
	for powerConverterSensorName, powerConverterSensorCfgEnt := range sMgr.powerConverterConfigDB {
		powerConverterSensorState, _ := sMgr.SensorState.PowerConverterSensor[powerConverterSensorName]
		powerConverterSensorPMData := objects.PowerConverterSensorPMData{
			TimeStamp: time.Now().String(),
			Value:     powerConverterSensorState.Value,
		}
		if powerConverterSensorCfgEnt.AdminState == "Enable" {
			eventKey := events.PowerConverterSensorKey{
				Name: powerConverterSensorName,
			}
			eventData := PowerConverterSensorEventData{
				Value: powerConverterSensorState.Value,
			}
			txEvent := eventUtils.TxEvent{
				Key:            eventKey,
				AdditionalInfo: "",
				AdditionalData: eventData,
			}
			var evts []events.EventId
			var curEvents EventStatus
			sMgr.powerConverterMsgStatusMutex.RLock()
			prevEventStatus := sMgr.powerConverterMsgStatus[powerConverterSensorName]
			sMgr.powerConverterMsgStatusMutex.RUnlock()
			if powerConverterSensorState.Value >= powerConverterSensorCfgEnt.HigherAlarmThreshold {
				curEvents.SentHigherAlarm = true
			}
			if powerConverterSensorState.Value >= powerConverterSensorCfgEnt.HigherWarningThreshold {
				curEvents.SentHigherWarn = true
			}
			if powerConverterSensorState.Value <= powerConverterSensorCfgEnt.LowerWarningThreshold {
				curEvents.SentLowerWarn = true
			}
			if powerConverterSensorState.Value <= powerConverterSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerAlarm = true
			}
			if prevEventStatus.SentHigherAlarm != curEvents.SentHigherAlarm {
				if curEvents.SentHigherAlarm == true {
					evts = append(evts, events.PowerConverterHigherTCAAlarm)
				} else {
					evts = append(evts, events.PowerConverterHigherTCAAlarmClear)
				}
			}
			if prevEventStatus.SentHigherWarn != curEvents.SentHigherWarn {
				if curEvents.SentHigherWarn == true {
					evts = append(evts, events.PowerConverterHigherTCAWarn)
				} else {
					evts = append(evts, events.PowerConverterHigherTCAWarnClear)
				}
			}
			if prevEventStatus.SentLowerAlarm != curEvents.SentLowerAlarm {
				if curEvents.SentLowerAlarm == true {
					evts = append(evts, events.PowerConverterLowerTCAAlarm)
				} else {
					evts = append(evts, events.PowerConverterLowerTCAAlarmClear)
				}
			}
			if prevEventStatus.SentLowerWarn != curEvents.SentLowerWarn {
				if curEvents.SentLowerWarn == true {
					evts = append(evts, events.PowerConverterLowerTCAWarn)
				} else {
					evts = append(evts, events.PowerConverterLowerTCAWarnClear)
				}
			}
			if prevEventStatus != curEvents {
				sMgr.powerConverterMsgStatusMutex.Lock()
				sMgr.powerConverterMsgStatus[powerConverterSensorName] = curEvents
				sMgr.powerConverterMsgStatusMutex.Unlock()
			}
			for _, evt := range evts {
				txEvent.EventId = evt
				txEvt := txEvent
				err := eventUtils.PublishEvents(&txEvt)
				if err != nil {
					sMgr.logger.Err("Error publish events")
				}
			}

		}

		switch class {
		case "Class-A":
			if powerConverterSensorCfgEnt.PMClassAAdminState == "Enable" {
				sMgr.powerConverterClassAPMMutex.Lock()
				sMgr.powerConverterSensorClassAPM[powerConverterSensorName].InsertIntoRingBuffer(powerConverterSensorPMData)
				sMgr.powerConverterClassAPMMutex.Unlock()
			}
		case "Class-B":
			if powerConverterSensorCfgEnt.PMClassBAdminState == "Enable" {
				sMgr.powerConverterClassBPMMutex.Lock()
				sMgr.powerConverterSensorClassBPM[powerConverterSensorName].InsertIntoRingBuffer(powerConverterSensorPMData)
				sMgr.powerConverterClassBPMMutex.Unlock()
			}
		case "Class-C":
			if powerConverterSensorCfgEnt.PMClassBAdminState == "Enable" {
				sMgr.powerConverterClassCPMMutex.Lock()
				sMgr.powerConverterSensorClassCPM[powerConverterSensorName].InsertIntoRingBuffer(powerConverterSensorPMData)
				sMgr.powerConverterClassCPMMutex.Unlock()
			}
		}
	}
	sMgr.powerConverterConfigMutex.RUnlock()
}

func (sMgr *SensorManager) StartSensorPMClass(class string) {
	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.logger.Err("Error getting sensor data during start of PM")
		return
	}
	sMgr.ProcessFanSensorPM(sMgr.SensorState, class)
	sMgr.ProcessTempSensorPM(sMgr.SensorState, class)
	sMgr.ProcessVoltageSensorPM(sMgr.SensorState, class)
	sMgr.ProcessPowerConverterSensorPM(sMgr.SensorState, class)
	sMgr.SensorStateMutex.Unlock()

	switch class {
	case "Class-A":
		classAPMFunc := func() {
			sMgr.SensorStateMutex.Lock()
			err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
			if err != nil {
				sMgr.logger.Err("Error getting sensor data during PM processing")
				sMgr.SensorStateMutex.Unlock()
				return
			}
			sMgr.ProcessFanSensorPM(sMgr.SensorState, class)
			sMgr.ProcessTempSensorPM(sMgr.SensorState, class)
			sMgr.ProcessVoltageSensorPM(sMgr.SensorState, class)
			sMgr.ProcessPowerConverterSensorPM(sMgr.SensorState, class)
			sMgr.SensorStateMutex.Unlock()
			sMgr.classAPMTimer.Reset(sensorClassAInterval)
		}
		sMgr.classAPMTimer = time.AfterFunc(sensorClassAInterval, classAPMFunc)
	case "Class-B":
		classBPMFunc := func() {
			sMgr.SensorStateMutex.Lock()
			err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
			if err != nil {
				sMgr.logger.Err("Error getting sensor data during PM processing")
				sMgr.SensorStateMutex.Unlock()
				return
			}
			sMgr.ProcessFanSensorPM(sMgr.SensorState, class)
			sMgr.ProcessTempSensorPM(sMgr.SensorState, class)
			sMgr.ProcessVoltageSensorPM(sMgr.SensorState, class)
			sMgr.ProcessPowerConverterSensorPM(sMgr.SensorState, class)
			sMgr.SensorStateMutex.Unlock()
			sMgr.classBPMTimer.Reset(sensorClassBInterval)
		}
		sMgr.classBPMTimer = time.AfterFunc(sensorClassBInterval, classBPMFunc)
	case "Class-C":
		classCPMFunc := func() {
			sMgr.SensorStateMutex.Lock()
			err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
			if err != nil {
				sMgr.logger.Err("Error getting sensor data during PM processing")
				sMgr.SensorStateMutex.Unlock()
				return
			}
			sMgr.ProcessFanSensorPM(sMgr.SensorState, class)
			sMgr.ProcessTempSensorPM(sMgr.SensorState, class)
			sMgr.ProcessVoltageSensorPM(sMgr.SensorState, class)
			sMgr.ProcessPowerConverterSensorPM(sMgr.SensorState, class)
			sMgr.SensorStateMutex.Unlock()
			sMgr.classCPMTimer.Reset(sensorClassCInterval)
		}
		sMgr.classCPMTimer = time.AfterFunc(sensorClassCInterval, classCPMFunc)
	}

}

func (sMgr *SensorManager) GetFanSensorPMState(Name string, Class string) (*objects.FanSensorPMState, error) {
	var fanSensorPMObj objects.FanSensorPMState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.fanConfigMutex.RLock()
	fanCfgEnt, exist := sMgr.fanConfigDB[Name]
	if !exist {
		sMgr.fanConfigMutex.RUnlock()
		return nil, errors.New("Invalid Fan Sensor Name")
	}

	fanSensorPMObj.Name = Name
	fanSensorPMObj.Class = Class
	switch Class {
	case "Class-A":
		if fanCfgEnt.PMClassAAdminState == "Enable" {
			sMgr.fanClassAPMMutex.RLock()
			fanSensorPMObj.Data = sMgr.fanSensorClassAPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.fanClassAPMMutex.RUnlock()
		} else {
			sMgr.fanConfigMutex.RUnlock()
			return nil, errors.New("PM Class A is Disabled")
		}
	case "Class-B":
		if fanCfgEnt.PMClassBAdminState == "Enable" {
			sMgr.fanClassBPMMutex.RLock()
			fanSensorPMObj.Data = sMgr.fanSensorClassBPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.fanClassBPMMutex.RUnlock()
		} else {
			sMgr.fanConfigMutex.RUnlock()
			return nil, errors.New("PM Class B is Disabled")
		}
	case "Class-C":
		if fanCfgEnt.PMClassCAdminState == "Enable" {
			sMgr.fanClassCPMMutex.RLock()
			fanSensorPMObj.Data = sMgr.fanSensorClassCPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.fanClassCPMMutex.RUnlock()
		} else {
			sMgr.fanConfigMutex.RUnlock()
			return nil, errors.New("PM Class C is Disabled")
		}
	default:
		sMgr.fanConfigMutex.RUnlock()
		return nil, errors.New("Invalid Class")
	}
	sMgr.fanConfigMutex.RUnlock()
	return &fanSensorPMObj, nil
}

func (sMgr *SensorManager) GetTempSensorPMState(Name string, Class string) (*objects.TemperatureSensorPMState, error) {
	var tempSensorPMObj objects.TemperatureSensorPMState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.tempConfigMutex.RLock()
	tempCfgEnt, exist := sMgr.tempConfigDB[Name]
	if !exist {
		sMgr.tempConfigMutex.RUnlock()
		return nil, errors.New("Invalid Temp Sensor Name")
	}

	tempSensorPMObj.Name = Name
	tempSensorPMObj.Class = Class
	switch Class {
	case "Class-A":
		if tempCfgEnt.PMClassAAdminState == "Enable" {
			sMgr.tempClassAPMMutex.RLock()
			tempSensorPMObj.Data = sMgr.tempSensorClassAPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.tempClassAPMMutex.RUnlock()
		} else {
			sMgr.tempConfigMutex.RUnlock()
			return nil, errors.New("PM Class A is Disabled")
		}
	case "Class-B":
		if tempCfgEnt.PMClassBAdminState == "Enable" {
			sMgr.tempClassBPMMutex.RLock()
			tempSensorPMObj.Data = sMgr.tempSensorClassBPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.tempClassBPMMutex.RUnlock()
		} else {
			sMgr.tempConfigMutex.RUnlock()
			return nil, errors.New("PM Class B is Disabled")
		}
	case "Class-C":
		if tempCfgEnt.PMClassCAdminState == "Enable" {
			sMgr.tempClassCPMMutex.RLock()
			tempSensorPMObj.Data = sMgr.tempSensorClassCPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.tempClassCPMMutex.RUnlock()
		} else {
			sMgr.tempConfigMutex.RUnlock()
			return nil, errors.New("PM Class C is Disabled")
		}
	default:
		sMgr.tempConfigMutex.RUnlock()
		return nil, errors.New("Invalid Class")
	}
	sMgr.tempConfigMutex.RUnlock()
	return &tempSensorPMObj, nil
}

func (sMgr *SensorManager) GetVoltageSensorPMState(Name string, Class string) (*objects.VoltageSensorPMState, error) {
	var voltageSensorPMObj objects.VoltageSensorPMState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.voltageConfigMutex.RLock()
	voltageCfgEnt, exist := sMgr.voltageConfigDB[Name]
	if !exist {
		sMgr.voltageConfigMutex.RUnlock()
		return nil, errors.New("Invalid Voltage Sensor Name")
	}

	voltageSensorPMObj.Name = Name
	voltageSensorPMObj.Class = Class
	switch Class {
	case "Class-A":
		if voltageCfgEnt.PMClassAAdminState == "Enable" {
			sMgr.voltageClassAPMMutex.RLock()
			voltageSensorPMObj.Data = sMgr.voltageSensorClassAPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.voltageClassAPMMutex.RUnlock()
		} else {
			sMgr.voltageConfigMutex.RUnlock()
			return nil, errors.New("PM Class A is Disabled")
		}
	case "Class-B":
		if voltageCfgEnt.PMClassBAdminState == "Enable" {
			sMgr.voltageClassBPMMutex.RLock()
			voltageSensorPMObj.Data = sMgr.voltageSensorClassBPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.voltageClassBPMMutex.RUnlock()
		} else {
			sMgr.voltageConfigMutex.RUnlock()
			return nil, errors.New("PM Class B is Disabled")
		}
	case "Class-C":
		if voltageCfgEnt.PMClassCAdminState == "Enable" {
			sMgr.voltageClassCPMMutex.RLock()
			voltageSensorPMObj.Data = sMgr.voltageSensorClassCPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.voltageClassCPMMutex.RUnlock()
		} else {
			sMgr.voltageConfigMutex.RUnlock()
			return nil, errors.New("PM Class C is Disabled")
		}
	default:
		sMgr.voltageConfigMutex.RUnlock()
		return nil, errors.New("Invalid Class")
	}
	sMgr.voltageConfigMutex.RUnlock()
	return &voltageSensorPMObj, nil
}

func (sMgr *SensorManager) GetPowerConverterSensorPMState(Name string, Class string) (*objects.PowerConverterSensorPMState, error) {
	var powerConverterSensorPMObj objects.PowerConverterSensorPMState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.powerConverterConfigMutex.RLock()
	powerConverterCfgEnt, exist := sMgr.powerConverterConfigDB[Name]
	if !exist {
		sMgr.powerConverterConfigMutex.RUnlock()
		return nil, errors.New("Invalid PowerConverter Sensor Name")
	}

	powerConverterSensorPMObj.Name = Name
	powerConverterSensorPMObj.Class = Class
	switch Class {
	case "Class-A":
		if powerConverterCfgEnt.PMClassAAdminState == "Enable" {
			sMgr.powerConverterClassAPMMutex.RLock()
			powerConverterSensorPMObj.Data = sMgr.powerConverterSensorClassAPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.powerConverterClassAPMMutex.RUnlock()
		} else {
			sMgr.powerConverterConfigMutex.RUnlock()
			return nil, errors.New("PM Class A is Disabled")
		}
	case "Class-B":
		if powerConverterCfgEnt.PMClassBAdminState == "Enable" {
			sMgr.powerConverterClassBPMMutex.RLock()
			powerConverterSensorPMObj.Data = sMgr.powerConverterSensorClassBPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.powerConverterClassBPMMutex.RUnlock()
		} else {
			sMgr.powerConverterConfigMutex.RUnlock()
			return nil, errors.New("PM Class B is Disabled")
		}
	case "Class-C":
		if powerConverterCfgEnt.PMClassCAdminState == "Enable" {
			sMgr.powerConverterClassCPMMutex.RLock()
			powerConverterSensorPMObj.Data = sMgr.powerConverterSensorClassCPM[Name].GetListOfEntriesFromRingBuffer()
			sMgr.powerConverterClassCPMMutex.RUnlock()
		} else {
			sMgr.powerConverterConfigMutex.RUnlock()
			return nil, errors.New("PM Class C is Disabled")
		}
	default:
		sMgr.powerConverterConfigMutex.RUnlock()
		return nil, errors.New("Invalid Class")
	}
	sMgr.powerConverterConfigMutex.RUnlock()
	return &powerConverterSensorPMObj, nil
}
