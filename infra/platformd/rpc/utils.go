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

package rpc

import (
	"infra/platformd/objects"
	"platformd"
)

func convertToRPCFmtPlatformState(obj *objects.PlatformState) *platformd.PlatformState {
	return &platformd.PlatformState{
		ObjName:      "Platform",
		ProductName:  obj.ProductName,
		SerialNum:    obj.SerialNum,
		Manufacturer: obj.Manufacturer,
		Vendor:       obj.Vendor,
		Release:      obj.Release,
		PlatformName: obj.PlatformName,
		Version:      obj.Version,
	}
}

func convertToRPCFmtFanState(obj *objects.FanState) *platformd.FanState {
	return &platformd.FanState{
		FanId:         obj.FanId,
		OperMode:      obj.OperMode,
		OperSpeed:     obj.OperSpeed,
		OperDirection: obj.OperDirection,
		Status:        obj.Status,
		Model:         obj.Model,
		SerialNum:     obj.SerialNum,
		LedId:         obj.LedId,
	}
}

func convertToRPCFmtFanConfig(obj *objects.FanConfig) *platformd.Fan {
	return &platformd.Fan{
		FanId:      obj.FanId,
		AdminSpeed: obj.AdminSpeed,
		AdminState: obj.AdminState,
	}
}

func convertRPCToObjFmtFanConfig(rpcObj *platformd.Fan) *objects.FanConfig {
	return &objects.FanConfig{
		FanId:      rpcObj.FanId,
		AdminSpeed: rpcObj.AdminSpeed,
		AdminState: rpcObj.AdminState,
	}
}

func convertToRPCFmtSfpConfig(obj *objects.SfpConfig) *platformd.Sfp {
	return &platformd.Sfp{
		SfpId: obj.SfpId,
	}
}

func convertToRPCFmtThermalState(obj *objects.ThermalState) *platformd.ThermalState {
	return &platformd.ThermalState{
		ThermalId:                 obj.ThermalId,
		Location:                  obj.Location,
		Temperature:               obj.Temperature,
		LowerWatermarkTemperature: obj.LowerWatermarkTemperature,
		UpperWatermarkTemperature: obj.UpperWatermarkTemperature,
		ShutdownTemperature:       obj.ShutdownTemperature,
	}
}

func convertToRPCFmtSfpState(obj *objects.SfpState) *platformd.SfpState {
	return &platformd.SfpState{
		SfpId:      obj.SfpId,
		SfpSpeed:   obj.SfpSpeed,
		SfpLOS:     obj.SfpLos,
		SfpPresent: obj.SfpPresent,
		SfpType:    obj.SfpType,
		SerialNum:  obj.SerialNum,
		EEPROM:     obj.EEPROM,
	}
}

func convertToRPCFmtFanSensorConfig(obj *objects.FanSensorConfig) *platformd.FanSensor {
	return &platformd.FanSensor{
		Name:                   obj.Name,
		AdminState:             obj.AdminState,
		HigherAlarmThreshold:   obj.HigherAlarmThreshold,
		HigherWarningThreshold: obj.HigherWarningThreshold,
		LowerAlarmThreshold:    obj.LowerAlarmThreshold,
		LowerWarningThreshold:  obj.LowerWarningThreshold,
		PMClassAAdminState:     obj.PMClassAAdminState,
		PMClassBAdminState:     obj.PMClassBAdminState,
		PMClassCAdminState:     obj.PMClassCAdminState,
	}
}

func convertRPCToObjFmtFanSensorConfig(rpcObj *platformd.FanSensor) *objects.FanSensorConfig {
	return &objects.FanSensorConfig{
		Name:                   rpcObj.Name,
		AdminState:             rpcObj.AdminState,
		HigherAlarmThreshold:   rpcObj.HigherAlarmThreshold,
		HigherWarningThreshold: rpcObj.HigherWarningThreshold,
		LowerAlarmThreshold:    rpcObj.LowerAlarmThreshold,
		LowerWarningThreshold:  rpcObj.LowerWarningThreshold,
		PMClassAAdminState:     rpcObj.PMClassAAdminState,
		PMClassBAdminState:     rpcObj.PMClassBAdminState,
		PMClassCAdminState:     rpcObj.PMClassCAdminState,
	}
}

func convertToRPCFmtFanSensorState(obj *objects.FanSensorState) *platformd.FanSensorState {
	return &platformd.FanSensorState{
		Name:         obj.Name,
		CurrentSpeed: obj.CurrentSpeed,
	}
}

func convertToRPCFmtTemperatureSensorConfig(obj *objects.TemperatureSensorConfig) *platformd.TemperatureSensor {
	return &platformd.TemperatureSensor{
		Name:                   obj.Name,
		AdminState:             obj.AdminState,
		HigherAlarmThreshold:   obj.HigherAlarmThreshold,
		HigherWarningThreshold: obj.HigherWarningThreshold,
		LowerAlarmThreshold:    obj.LowerAlarmThreshold,
		LowerWarningThreshold:  obj.LowerWarningThreshold,
		PMClassAAdminState:     obj.PMClassAAdminState,
		PMClassBAdminState:     obj.PMClassBAdminState,
		PMClassCAdminState:     obj.PMClassCAdminState,
	}
}

func convertRPCToObjFmtTemperatureSensorConfig(rpcObj *platformd.TemperatureSensor) *objects.TemperatureSensorConfig {
	return &objects.TemperatureSensorConfig{
		Name:                   rpcObj.Name,
		AdminState:             rpcObj.AdminState,
		HigherAlarmThreshold:   rpcObj.HigherAlarmThreshold,
		HigherWarningThreshold: rpcObj.HigherWarningThreshold,
		LowerAlarmThreshold:    rpcObj.LowerAlarmThreshold,
		LowerWarningThreshold:  rpcObj.LowerWarningThreshold,
		PMClassAAdminState:     rpcObj.PMClassAAdminState,
		PMClassBAdminState:     rpcObj.PMClassBAdminState,
		PMClassCAdminState:     rpcObj.PMClassCAdminState,
	}
}

func convertToRPCFmtTemperatureSensorState(obj *objects.TemperatureSensorState) *platformd.TemperatureSensorState {
	return &platformd.TemperatureSensorState{
		Name:               obj.Name,
		CurrentTemperature: obj.CurrentTemperature,
	}
}

func convertToRPCFmtVoltageSensorConfig(obj *objects.VoltageSensorConfig) *platformd.VoltageSensor {
	return &platformd.VoltageSensor{
		Name:                   obj.Name,
		AdminState:             obj.AdminState,
		HigherAlarmThreshold:   obj.HigherAlarmThreshold,
		HigherWarningThreshold: obj.HigherWarningThreshold,
		LowerAlarmThreshold:    obj.LowerAlarmThreshold,
		LowerWarningThreshold:  obj.LowerWarningThreshold,
		PMClassAAdminState:     obj.PMClassAAdminState,
		PMClassBAdminState:     obj.PMClassBAdminState,
		PMClassCAdminState:     obj.PMClassCAdminState,
	}
}

func convertRPCToObjFmtVoltageSensorConfig(rpcObj *platformd.VoltageSensor) *objects.VoltageSensorConfig {
	return &objects.VoltageSensorConfig{
		Name:                   rpcObj.Name,
		AdminState:             rpcObj.AdminState,
		HigherAlarmThreshold:   rpcObj.HigherAlarmThreshold,
		HigherWarningThreshold: rpcObj.HigherWarningThreshold,
		LowerAlarmThreshold:    rpcObj.LowerAlarmThreshold,
		LowerWarningThreshold:  rpcObj.LowerWarningThreshold,
		PMClassAAdminState:     rpcObj.PMClassAAdminState,
		PMClassBAdminState:     rpcObj.PMClassBAdminState,
		PMClassCAdminState:     rpcObj.PMClassCAdminState,
	}
}

func convertToRPCFmtVoltageSensorState(obj *objects.VoltageSensorState) *platformd.VoltageSensorState {
	return &platformd.VoltageSensorState{
		Name:           obj.Name,
		CurrentVoltage: obj.CurrentVoltage,
	}
}

func convertToRPCFmtPowerConverterSensorConfig(obj *objects.PowerConverterSensorConfig) *platformd.PowerConverterSensor {
	return &platformd.PowerConverterSensor{
		Name:                   obj.Name,
		AdminState:             obj.AdminState,
		HigherAlarmThreshold:   obj.HigherAlarmThreshold,
		HigherWarningThreshold: obj.HigherWarningThreshold,
		LowerAlarmThreshold:    obj.LowerAlarmThreshold,
		LowerWarningThreshold:  obj.LowerWarningThreshold,
		PMClassAAdminState:     obj.PMClassAAdminState,
		PMClassBAdminState:     obj.PMClassBAdminState,
		PMClassCAdminState:     obj.PMClassCAdminState,
	}
}

func convertRPCToObjFmtPowerConverterSensorConfig(rpcObj *platformd.PowerConverterSensor) *objects.PowerConverterSensorConfig {
	return &objects.PowerConverterSensorConfig{
		Name:                   rpcObj.Name,
		AdminState:             rpcObj.AdminState,
		HigherAlarmThreshold:   rpcObj.HigherAlarmThreshold,
		HigherWarningThreshold: rpcObj.HigherWarningThreshold,
		LowerAlarmThreshold:    rpcObj.LowerAlarmThreshold,
		LowerWarningThreshold:  rpcObj.LowerWarningThreshold,
		PMClassAAdminState:     rpcObj.PMClassAAdminState,
		PMClassBAdminState:     rpcObj.PMClassBAdminState,
		PMClassCAdminState:     rpcObj.PMClassCAdminState,
	}
}

func convertToRPCFmtPowerConverterSensorState(obj *objects.PowerConverterSensorState) *platformd.PowerConverterSensorState {
	return &platformd.PowerConverterSensorState{
		Name:         obj.Name,
		CurrentPower: obj.CurrentPower,
	}
}

func convertToRPCFmtQsfpConfig(obj *objects.QsfpConfig) *platformd.Qsfp {
	return &platformd.Qsfp{
		QsfpId:                   obj.QsfpId,
		AdminState:               obj.AdminState,
		HigherAlarmTemperature:   obj.HigherAlarmTemperature,
		HigherAlarmVoltage:       obj.HigherAlarmVoltage,
		HigherWarningTemperature: obj.HigherWarningTemperature,
		HigherWarningVoltage:     obj.HigherWarningVoltage,
		LowerAlarmTemperature:    obj.LowerAlarmTemperature,
		LowerAlarmVoltage:        obj.LowerAlarmVoltage,
		LowerWarningTemperature:  obj.LowerWarningTemperature,
		LowerWarningVoltage:      obj.LowerWarningVoltage,
		PMClassAAdminState:       obj.PMClassAAdminState,
		PMClassBAdminState:       obj.PMClassBAdminState,
		PMClassCAdminState:       obj.PMClassCAdminState,
	}
}

func convertRPCToObjFmtQsfpConfig(rpcObj *platformd.Qsfp) *objects.QsfpConfig {
	return &objects.QsfpConfig{
		QsfpId:                   rpcObj.QsfpId,
		AdminState:               rpcObj.AdminState,
		HigherAlarmTemperature:   rpcObj.HigherAlarmTemperature,
		HigherAlarmVoltage:       rpcObj.HigherAlarmVoltage,
		HigherWarningTemperature: rpcObj.HigherWarningTemperature,
		HigherWarningVoltage:     rpcObj.HigherWarningVoltage,
		LowerAlarmTemperature:    rpcObj.LowerAlarmTemperature,
		LowerAlarmVoltage:        rpcObj.LowerAlarmVoltage,
		LowerWarningTemperature:  rpcObj.LowerWarningTemperature,
		LowerWarningVoltage:      rpcObj.LowerWarningVoltage,
		PMClassAAdminState:       rpcObj.PMClassAAdminState,
		PMClassBAdminState:       rpcObj.PMClassBAdminState,
		PMClassCAdminState:       rpcObj.PMClassCAdminState,
	}
}

func convertToRPCFmtQsfpState(obj *objects.QsfpState) *platformd.QsfpState {
	return &platformd.QsfpState{
		QsfpId:             obj.QsfpId,
		Present:            obj.Present,
		VendorName:         obj.VendorName,
		VendorOUI:          obj.VendorOUI,
		VendorPartNumber:   obj.VendorPartNumber,
		VendorRevision:     obj.VendorRevision,
		VendorSerialNumber: obj.VendorSerialNumber,
		DataCode:           obj.DataCode,
		Temperature:        obj.Temperature,
		Voltage:            obj.Voltage,
		CurrBER:            obj.CurrBER,
		AccBER:             obj.AccBER,
		MinBER:             obj.MinBER,
		MaxBER:             obj.MaxBER,
		UDF0:               obj.UDF0,
		UDF1:               obj.UDF1,
		UDF2:               obj.UDF2,
		UDF3:               obj.UDF3,
	}
}

func convertToRPCFmtQsfpChannelConfig(obj *objects.QsfpChannelConfig) *platformd.QsfpChannel {
	return &platformd.QsfpChannel{
		QsfpId:               obj.QsfpId,
		ChannelNum:           obj.ChannelNum,
		AdminState:           obj.AdminState,
		HigherAlarmRXPower:   obj.HigherAlarmRXPower,
		HigherAlarmTXPower:   obj.HigherAlarmTXPower,
		HigherAlarmTXBias:    obj.HigherAlarmTXBias,
		HigherWarningRXPower: obj.HigherWarningRXPower,
		HigherWarningTXPower: obj.HigherWarningTXPower,
		HigherWarningTXBias:  obj.HigherWarningTXBias,
		LowerAlarmRXPower:    obj.LowerAlarmRXPower,
		LowerAlarmTXPower:    obj.LowerAlarmTXPower,
		LowerAlarmTXBias:     obj.LowerAlarmTXBias,
		LowerWarningRXPower:  obj.LowerWarningRXPower,
		LowerWarningTXPower:  obj.LowerWarningTXPower,
		LowerWarningTXBias:   obj.LowerWarningTXBias,
		PMClassAAdminState:   obj.PMClassAAdminState,
		PMClassBAdminState:   obj.PMClassBAdminState,
		PMClassCAdminState:   obj.PMClassCAdminState,
	}
}

func convertRPCToObjFmtQsfpChannelConfig(rpcObj *platformd.QsfpChannel) *objects.QsfpChannelConfig {
	return &objects.QsfpChannelConfig{
		QsfpId:               rpcObj.QsfpId,
		ChannelNum:           rpcObj.ChannelNum,
		AdminState:           rpcObj.AdminState,
		HigherAlarmRXPower:   rpcObj.HigherAlarmRXPower,
		HigherAlarmTXPower:   rpcObj.HigherAlarmTXPower,
		HigherAlarmTXBias:    rpcObj.HigherAlarmTXBias,
		HigherWarningRXPower: rpcObj.HigherWarningRXPower,
		HigherWarningTXPower: rpcObj.HigherWarningTXPower,
		HigherWarningTXBias:  rpcObj.HigherWarningTXBias,
		LowerAlarmRXPower:    rpcObj.LowerAlarmRXPower,
		LowerAlarmTXPower:    rpcObj.LowerAlarmTXPower,
		LowerAlarmTXBias:     rpcObj.LowerAlarmTXBias,
		LowerWarningRXPower:  rpcObj.LowerWarningRXPower,
		LowerWarningTXPower:  rpcObj.LowerWarningTXPower,
		LowerWarningTXBias:   rpcObj.LowerWarningTXBias,
		PMClassAAdminState:   rpcObj.PMClassAAdminState,
		PMClassBAdminState:   rpcObj.PMClassBAdminState,
		PMClassCAdminState:   rpcObj.PMClassCAdminState,
	}
}

func convertToRPCFmtQsfpChannelState(obj *objects.QsfpChannelState) *platformd.QsfpChannelState {
	return &platformd.QsfpChannelState{
		QsfpId:     obj.QsfpId,
		ChannelNum: obj.ChannelNum,
		Present:    obj.Present,
		RXPower:    obj.RXPower,
		TXPower:    obj.TXPower,
		TXBias:     obj.TXBias,
	}
}

func convertToRPCFmtPlatformMgmtDeviceState(obj *objects.PlatformMgmtDeviceState) *platformd.PlatformMgmtDeviceState {
	return &platformd.PlatformMgmtDeviceState{
		DeviceName:  obj.DeviceName,
		Uptime:      obj.Uptime,
		Description: obj.Description,
		ResetReason: obj.ResetReason,
		MemoryUsage: obj.MemoryUsage,
		Version:     obj.Version,
		CPUUsage:    obj.CPUUsage,
	}
}

func convertToRPCFmtFanSensorPMState(obj *objects.FanSensorPMState) *platformd.FanSensorPMDataState {
	length := len(obj.Data)
	data := make([]*platformd.FanSensorPMData, length)
	// Revisit
	for idx := 0; idx < length; idx++ {
		objData := obj.Data[length-1-idx].(objects.FanSensorPMData)
		data[idx] = new(platformd.FanSensorPMData)
		data[idx].TimeStamp = objData.TimeStamp
		data[idx].Value = objData.Value
	}
	return &platformd.FanSensorPMDataState{
		Name:  obj.Name,
		Class: obj.Class,
		Data:  data,
	}

}

func convertToRPCFmtTempSensorPMState(obj *objects.TemperatureSensorPMState) *platformd.TemperatureSensorPMDataState {
	length := len(obj.Data)
	data := make([]*platformd.TemperatureSensorPMData, length)
	// Revisit
	for idx := 0; idx < length; idx++ {
		objData := obj.Data[length-1-idx].(objects.TemperatureSensorPMData)
		data[idx] = new(platformd.TemperatureSensorPMData)
		data[idx].TimeStamp = objData.TimeStamp
		data[idx].Value = objData.Value
	}
	return &platformd.TemperatureSensorPMDataState{
		Name:  obj.Name,
		Class: obj.Class,
		Data:  data,
	}

}

func convertToRPCFmtVoltageSensorPMState(obj *objects.VoltageSensorPMState) *platformd.VoltageSensorPMDataState {
	length := len(obj.Data)
	data := make([]*platformd.VoltageSensorPMData, length)
	// Revisit
	for idx := 0; idx < length; idx++ {
		objData := obj.Data[length-1-idx].(objects.VoltageSensorPMData)
		data[idx] = new(platformd.VoltageSensorPMData)
		data[idx].TimeStamp = objData.TimeStamp
		data[idx].Value = objData.Value
	}
	return &platformd.VoltageSensorPMDataState{
		Name:  obj.Name,
		Class: obj.Class,
		Data:  data,
	}

}

func convertToRPCFmtPowerConverterSensorPMState(obj *objects.PowerConverterSensorPMState) *platformd.PowerConverterSensorPMDataState {
	length := len(obj.Data)
	data := make([]*platformd.PowerConverterSensorPMData, length)
	// Revisit
	for idx := 0; idx < length; idx++ {
		objData := obj.Data[length-1-idx].(objects.PowerConverterSensorPMData)
		data[idx] = new(platformd.PowerConverterSensorPMData)
		data[idx].TimeStamp = objData.TimeStamp
		data[idx].Value = objData.Value
	}
	return &platformd.PowerConverterSensorPMDataState{
		Name:  obj.Name,
		Class: obj.Class,
		Data:  data,
	}

}

func convertToRPCFmtQsfpPMState(obj *objects.QsfpPMState) *platformd.QsfpPMDataState {
	length := len(obj.Data)
	data := make([]*platformd.QsfpPMData, length)
	// Revisit
	for idx := 0; idx < length; idx++ {
		objData := obj.Data[length-1-idx].(objects.QsfpPMData)
		data[idx] = new(platformd.QsfpPMData)
		data[idx].TimeStamp = objData.TimeStamp
		data[idx].Value = objData.Value
	}
	return &platformd.QsfpPMDataState{
		QsfpId:   obj.QsfpId,
		Resource: obj.Resource,
		Class:    obj.Class,
		Data:     data,
	}
}

func convertToRPCFmtQsfpChannelPMState(obj *objects.QsfpChannelPMState) *platformd.QsfpChannelPMDataState {
	length := len(obj.Data)
	data := make([]*platformd.QsfpChannelPMData, length)
	// Revisit
	for idx := 0; idx < length; idx++ {
		objData := obj.Data[length-1-idx].(objects.QsfpPMData)
		data[idx] = new(platformd.QsfpChannelPMData)
		data[idx].TimeStamp = objData.TimeStamp
		data[idx].Value = objData.Value
	}
	return &platformd.QsfpChannelPMDataState{
		QsfpId:     obj.QsfpId,
		ChannelNum: obj.ChannelNum,
		Resource:   obj.Resource,
		Class:      obj.Class,
		Data:       data,
	}
}

func convertToRPCFmtPsuState(obj *objects.PsuState) *platformd.PsuState {
	return &platformd.PsuState{
		PsuId:      obj.PsuId,
		AdminState: obj.AdminState,
		ModelNum:   obj.ModelNum,
		SerialNum:  obj.SerialNum,
		Vin:        obj.Vin,
		Vout:       obj.Vout,
		Iin:        obj.Iin,
		Iout:       obj.Iout,
		Pin:        obj.Pin,
		Pout:       obj.Pout,
	}
}

//type LedState struct {
//	baseObj
//	LedId       int32  `SNAPROUTE: "KEY", ACCESS:"r", MULTIPLICITY:"*", DESCRIPTION: "LED id", DEFAULT:0`
//	LedIdentify string `DESCRIPTION: "LED represents FAN/PSU/RESET etc"`
//	LedState    string `DESCRIPTION: "LED State ON/OFF"`
//	LedColor    string `DESCRIPTION: "LED Color"`
//}

//type LedState struct {
//    LedId       int32
//    LedIdentify string
//    LedState    string
//    LedColor    string
//}

func convertToRPCFmtLedState(obj *objects.LedState) *platformd.LedState {
	return &platformd.LedState{
		LedId:       obj.LedId,
		LedIdentify: obj.LedIdentify,
		LedState:    obj.LedState,
		LedColor:    obj.LedColor,
	}
}
