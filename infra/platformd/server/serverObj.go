//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
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

type ServerOpId int

const (
	GET_FAN_STATE ServerOpId = iota
	GET_BULK_FAN_STATE
	GET_FAN_CONFIG
	GET_BULK_FAN_CONFIG
	UPDATE_FAN_CONFIG
	GET_SFP_STATE
	GET_BULK_SFP_STATE
	GET_SFP_CONFIG
	GET_BULK_SFP_CONFIG
	UPDATE_SFP_CONFIG
	GET_PSU_STATE
	GET_BULK_PSU_STATE
	GET_PSU_CONFIG
	GET_BULK_PSU_CONFIG
	UPDATE_PSU_CONFIG
	GET_PLATFORM_STATE
	GET_BULK_PLATFORM_STATE
	GET_THERMAL_STATE
	GET_BULK_THERMAL_STATE
	GET_FAN_SENSOR_STATE
	GET_BULK_FAN_SENSOR_STATE
	GET_BULK_FAN_SENSOR_CONFIG
	UPDATE_FAN_SENSOR_CONFIG
	GET_FAN_SENSOR_PM_STATE
	GET_TEMPERATURE_SENSOR_STATE
	GET_BULK_TEMPERATURE_SENSOR_STATE
	GET_BULK_TEMPERATURE_SENSOR_CONFIG
	UPDATE_TEMPERATURE_SENSOR_CONFIG
	GET_TEMPERATURE_SENSOR_PM_STATE
	GET_VOLTAGE_SENSOR_STATE
	GET_BULK_VOLTAGE_SENSOR_STATE
	GET_BULK_VOLTAGE_SENSOR_CONFIG
	UPDATE_VOLTAGE_SENSOR_CONFIG
	GET_VOLTAGE_SENSOR_PM_STATE
	GET_POWER_CONVERTER_SENSOR_STATE
	GET_BULK_POWER_CONVERTER_SENSOR_STATE
	GET_BULK_POWER_CONVERTER_SENSOR_CONFIG
	UPDATE_POWER_CONVERTER_SENSOR_CONFIG
	GET_POWER_CONVERTER_SENSOR_PM_STATE
	GET_QSFP_STATE
	GET_BULK_QSFP_STATE
	GET_BULK_QSFP_CONFIG
	UPDATE_QSFP_CONFIG
	GET_QSFP_PM_STATE
	GET_QSFP_CHANNEL_STATE
	GET_BULK_QSFP_CHANNEL_STATE
	GET_BULK_QSFP_CHANNEL_CONFIG
	UPDATE_QSFP_CHANNEL_CONFIG
	GET_QSFP_CHANNEL_PM_STATE
	GET_PLATFORM_MGMT_DEVICE_STATE
	GET_BULK_PLATFORM_MGMT_DEVICE_STATE
	GET_LED_STATE
	GET_BULK_LED_STATE
)

type ServerRequest struct {
	Op   ServerOpId
	Data interface{}
}

type GetBulkInArgs struct {
	FromIdx int
	Count   int
}

type GetFanStateInArgs struct {
	FanId int32
}

type GetFanStateOutArgs struct {
	Obj *objects.FanState
	Err error
}

type GetBulkFanStateOutArgs struct {
	BulkInfo *objects.FanStateGetInfo
	Err      error
}

type GetFanConfigInArgs struct {
	FanId int32
}

type GetFanConfigOutArgs struct {
	Obj *objects.FanConfig
	Err error
}

type GetBulkFanConfigOutArgs struct {
	BulkInfo *objects.FanConfigGetInfo
	Err      error
}

type UpdateFanConfigInArgs struct {
	FanOldCfg *objects.FanConfig
	FanNewCfg *objects.FanConfig
	AttrSet   []bool
}

type UpdateConfigOutArgs struct {
	RetVal bool
	Err    error
}

type GetSfpStateInArgs struct {
	SfpId int32
}

type GetSfpStateOutArgs struct {
	Obj *objects.SfpState
	Err error
}

type GetBulkSfpStateOutArgs struct {
	BulkInfo *objects.SfpStateGetInfo
	Err      error
}

type GetSfpConfigInArgs struct {
	SfpId int32
}

type GetSfpConfigOutArgs struct {
	Obj *objects.SfpConfig
	Err error
}

type GetBulkSfpConfigOutArgs struct {
	BulkInfo *objects.SfpConfigGetInfo
	Err      error
}

type UpdateSfpConfigInArgs struct {
	SfpOldCfg *objects.SfpConfig
	SfpNewCfg *objects.SfpConfig
	AttrSet   []bool
}

type GetPlatformStateInArgs struct {
	ObjName string
}

type GetPlatformStateOutArgs struct {
	Obj *objects.PlatformState
	Err error
}

type GetBulkPlatformStateOutArgs struct {
	BulkInfo *objects.PlatformStateGetInfo
	Err      error
}

type GetThermalStateInArgs struct {
	ThermalId int32
}

type GetThermalStateOutArgs struct {
	Obj *objects.ThermalState
	Err error
}

type GetBulkThermalStateOutArgs struct {
	BulkInfo *objects.ThermalStateGetInfo
	Err      error
}

type GetFanSensorStateInArgs struct {
	Name string
}

type GetFanSensorStateOutArgs struct {
	Obj *objects.FanSensorState
	Err error
}

type GetBulkFanSensorStateOutArgs struct {
	BulkInfo *objects.FanSensorStateGetInfo
	Err      error
}

type GetBulkFanSensorConfigOutArgs struct {
	BulkInfo *objects.FanSensorConfigGetInfo
	Err      error
}

type UpdateFanSensorConfigInArgs struct {
	FanSensorOldCfg *objects.FanSensorConfig
	FanSensorNewCfg *objects.FanSensorConfig
	AttrSet         []bool
}

type GetTemperatureSensorStateInArgs struct {
	Name string
}

type GetTemperatureSensorStateOutArgs struct {
	Obj *objects.TemperatureSensorState
	Err error
}

type GetBulkTemperatureSensorStateOutArgs struct {
	BulkInfo *objects.TemperatureSensorStateGetInfo
	Err      error
}

type GetBulkTemperatureSensorConfigOutArgs struct {
	BulkInfo *objects.TemperatureSensorConfigGetInfo
	Err      error
}

type UpdateTemperatureSensorConfigInArgs struct {
	TemperatureSensorOldCfg *objects.TemperatureSensorConfig
	TemperatureSensorNewCfg *objects.TemperatureSensorConfig
	AttrSet                 []bool
}

type GetVoltageSensorStateInArgs struct {
	Name string
}

type GetVoltageSensorStateOutArgs struct {
	Obj *objects.VoltageSensorState
	Err error
}

type GetBulkVoltageSensorStateOutArgs struct {
	BulkInfo *objects.VoltageSensorStateGetInfo
	Err      error
}

type GetBulkVoltageSensorConfigOutArgs struct {
	BulkInfo *objects.VoltageSensorConfigGetInfo
	Err      error
}

type UpdateVoltageSensorConfigInArgs struct {
	VoltageSensorOldCfg *objects.VoltageSensorConfig
	VoltageSensorNewCfg *objects.VoltageSensorConfig
	AttrSet             []bool
}

type GetPowerConverterSensorStateInArgs struct {
	Name string
}

type GetPowerConverterSensorStateOutArgs struct {
	Obj *objects.PowerConverterSensorState
	Err error
}

type GetBulkPowerConverterSensorStateOutArgs struct {
	BulkInfo *objects.PowerConverterSensorStateGetInfo
	Err      error
}

type GetBulkPowerConverterSensorConfigOutArgs struct {
	BulkInfo *objects.PowerConverterSensorConfigGetInfo
	Err      error
}

type UpdatePowerConverterSensorConfigInArgs struct {
	PowerConverterSensorOldCfg *objects.PowerConverterSensorConfig
	PowerConverterSensorNewCfg *objects.PowerConverterSensorConfig
	AttrSet                    []bool
}

type GetQsfpStateInArgs struct {
	QsfpId int32
}

type GetQsfpStateOutArgs struct {
	Obj *objects.QsfpState
	Err error
}

type GetBulkQsfpStateOutArgs struct {
	BulkInfo *objects.QsfpStateGetInfo
	Err      error
}

type GetBulkQsfpConfigOutArgs struct {
	BulkInfo *objects.QsfpConfigGetInfo
	Err      error
}

type UpdateQsfpConfigInArgs struct {
	QsfpOldCfg *objects.QsfpConfig
	QsfpNewCfg *objects.QsfpConfig
	AttrSet    []bool
}

type GetQsfpPMStateInArgs struct {
	QsfpId   int32
	Resource string
	Class    string
}

type GetQsfpPMStateOutArgs struct {
	Obj *objects.QsfpPMState
	Err error
}

type GetQsfpChannelStateInArgs struct {
	QsfpId     int32
	ChannelNum int32
}

type GetQsfpChannelStateOutArgs struct {
	Obj *objects.QsfpChannelState
	Err error
}

type GetBulkQsfpChannelStateOutArgs struct {
	BulkInfo *objects.QsfpChannelStateGetInfo
	Err      error
}

type GetBulkQsfpChannelConfigOutArgs struct {
	BulkInfo *objects.QsfpChannelConfigGetInfo
	Err      error
}

type UpdateQsfpChannelConfigInArgs struct {
	QsfpChannelOldCfg *objects.QsfpChannelConfig
	QsfpChannelNewCfg *objects.QsfpChannelConfig
	AttrSet           []bool
}

type GetQsfpChannelPMStateInArgs struct {
	QsfpId     int32
	ChannelNum int32
	Resource   string
	Class      string
}

type GetQsfpChannelPMStateOutArgs struct {
	Obj *objects.QsfpChannelPMState
	Err error
}

type GetPlatformMgmtDeviceStateInArgs struct {
	DeviceName string
}

type GetPlatformMgmtDeviceStateOutArgs struct {
	Obj *objects.PlatformMgmtDeviceState
	Err error
}

type GetBulkPlatformMgmtDeviceStateOutArgs struct {
	BulkInfo *objects.PlatformMgmtDeviceStateGetInfo
	Err      error
}

type GetFanSensorPMStateInArgs struct {
	Name  string
	Class string
}

type GetFanSensorPMStateOutArgs struct {
	Obj *objects.FanSensorPMState
	Err error
}

type GetTempSensorPMStateInArgs struct {
	Name  string
	Class string
}

type GetTempSensorPMStateOutArgs struct {
	Obj *objects.TemperatureSensorPMState
	Err error
}

type GetVoltageSensorPMStateInArgs struct {
	Name  string
	Class string
}

type GetVoltageSensorPMStateOutArgs struct {
	Obj *objects.VoltageSensorPMState
	Err error
}

type GetPowerConverterSensorPMStateInArgs struct {
	Name  string
	Class string
}

type GetPowerConverterSensorPMStateOutArgs struct {
	Obj *objects.PowerConverterSensorPMState
	Err error
}

type GetPsuStateInArgs struct {
	PsuId int32
}

type GetPsuStateOutArgs struct {
	Obj *objects.PsuState
	Err error
}

type GetBulkPsuStateOutArgs struct {
	BulkInfo *objects.PsuStateGetInfo
	Err      error
}

type GetLedStateInArgs struct {
	LedId int32
}

type GetLedStateOutArgs struct {
	Obj *objects.LedState
	Err error
}

type GetBulkLedStateOutArgs struct {
	BulkInfo *objects.LedStateGetInfo
	Err      error
}
