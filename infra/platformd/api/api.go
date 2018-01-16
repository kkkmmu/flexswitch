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

package api

import (
	"errors"
	"infra/platformd/objects"
	"infra/platformd/server"
)

var svr *server.PlatformdServer

//var logger *logging.LogginIntf

//Initialize server handle
func InitApiLayer(server *server.PlatformdServer) {
	svr = server
	svr.Logger.Info("Initializing API layer")
}

func GetPlatformState(objName string) (*objects.PlatformState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_PLATFORM_STATE,
		Data: interface{}(&server.GetPlatformStateInArgs{
			ObjName: objName,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetPlatformStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetPlatformState")
	}
}

func GetBulkPlatformState(fromIdx, count int) (*objects.PlatformStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_PLATFORM_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkPlatformStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkPlatformState")
	}
}

func GetFanState(fanId int32) (*objects.FanState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_FAN_STATE,
		Data: interface{}(&server.GetFanStateInArgs{
			FanId: fanId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetFanStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetFanState")
	}
}

func GetBulkFanState(fromIdx, count int) (*objects.FanStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_FAN_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkFanStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkFanState")
	}
}

func UpdateFan(oldCfg *objects.FanConfig, newCfg *objects.FanConfig, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_FAN_CONFIG,
		Data: interface{}(&server.UpdateFanConfigInArgs{
			FanOldCfg: oldCfg,
			FanNewCfg: newCfg,
			AttrSet:   attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateFan")
}

func GetFanConfig(fanId int32) (*objects.FanConfig, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_FAN_CONFIG,
		Data: interface{}(&server.GetFanConfigInArgs{
			FanId: fanId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetFanConfigOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetFanConfig")
	}
}

func GetBulkFanConfig(fromIdx, count int) (*objects.FanConfigGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_FAN_CONFIG,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkFanConfigOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkFanConfig")
	}
}

func UpdateSfp(oldCfg *objects.SfpConfig, newCfg *objects.SfpConfig, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_SFP_CONFIG,
		Data: interface{}(&server.UpdateSfpConfigInArgs{
			SfpOldCfg: oldCfg,
			SfpNewCfg: newCfg,
			AttrSet:   attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateFan")
}

func GetSfpConfig(sfpId int32) (*objects.SfpConfig, error) {
	var obj objects.SfpConfig

	return &obj, nil
}

func GetBulkSfpConfig(fromIdx, count int) (*objects.SfpConfigGetInfo, error) {
	var obj objects.SfpConfigGetInfo

	return &obj, nil
}

func GetSfpState(sfpId int32) (*objects.SfpState, error) {
	var obj objects.SfpState

	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_SFP_STATE,
		Data: interface{}(&server.GetSfpStateInArgs{
			SfpId: sfpId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetSfpStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during SfpStateGet")
	}

	return &obj, nil
}

func GetBulkSfpState(fromIdx, count int) (*objects.SfpStateGetInfo, error) {
	var obj objects.SfpStateGetInfo

	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_SFP_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkSfpStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkSfpState")
	}

	return &obj, nil
}

func GetThermalState(thermalId int32) (*objects.ThermalState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_THERMAL_STATE,
		Data: interface{}(&server.GetThermalStateInArgs{
			ThermalId: thermalId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetThermalStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetThermalState")
	}
}

func GetBulkThermalState(fromIdx, count int) (*objects.ThermalStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_THERMAL_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkThermalStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkThermalState")
	}
}

func GetFanSensorState(name string) (*objects.FanSensorState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_FAN_SENSOR_STATE,
		Data: interface{}(&server.GetFanSensorStateInArgs{
			Name: name,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetFanSensorStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetFanSensorState")
	}
}

func GetBulkFanSensorState(fromIdx, count int) (*objects.FanSensorStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_FAN_SENSOR_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkFanSensorStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkFanSensorState")
	}
}

func UpdateFanSensor(oldCfg *objects.FanSensorConfig, newCfg *objects.FanSensorConfig, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_FAN_SENSOR_CONFIG,
		Data: interface{}(&server.UpdateFanSensorConfigInArgs{
			FanSensorOldCfg: oldCfg,
			FanSensorNewCfg: newCfg,
			AttrSet:         attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateFanSensor")
}

func GetBulkFanSensorConfig(fromIdx, count int) (*objects.FanSensorConfigGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_FAN_SENSOR_CONFIG,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkFanSensorConfigOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkFanSensorConfig")
	}
}

func GetTemperatureSensorState(name string) (*objects.TemperatureSensorState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_TEMPERATURE_SENSOR_STATE,
		Data: interface{}(&server.GetTemperatureSensorStateInArgs{
			Name: name,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetTemperatureSensorStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetTemperatureSensorState")
	}
}

func GetBulkTemperatureSensorState(fromIdx, count int) (*objects.TemperatureSensorStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_TEMPERATURE_SENSOR_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkTemperatureSensorStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkTemperatureSensorState")
	}
}

func UpdateTemperatureSensor(oldCfg *objects.TemperatureSensorConfig, newCfg *objects.TemperatureSensorConfig, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_TEMPERATURE_SENSOR_CONFIG,
		Data: interface{}(&server.UpdateTemperatureSensorConfigInArgs{
			TemperatureSensorOldCfg: oldCfg,
			TemperatureSensorNewCfg: newCfg,
			AttrSet:                 attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateTemperatureSensor")
}

func GetBulkTemperatureSensorConfig(fromIdx, count int) (*objects.TemperatureSensorConfigGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_TEMPERATURE_SENSOR_CONFIG,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkTemperatureSensorConfigOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkTemperatureSensorConfig")
	}
}

func GetVoltageSensorState(name string) (*objects.VoltageSensorState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_VOLTAGE_SENSOR_STATE,
		Data: interface{}(&server.GetVoltageSensorStateInArgs{
			Name: name,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetVoltageSensorStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetVoltageSensorState")
	}
}

func GetBulkVoltageSensorState(fromIdx, count int) (*objects.VoltageSensorStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_VOLTAGE_SENSOR_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkVoltageSensorStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkVoltageSensorState")
	}
}

func UpdateVoltageSensor(oldCfg *objects.VoltageSensorConfig, newCfg *objects.VoltageSensorConfig, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_VOLTAGE_SENSOR_CONFIG,
		Data: interface{}(&server.UpdateVoltageSensorConfigInArgs{
			VoltageSensorOldCfg: oldCfg,
			VoltageSensorNewCfg: newCfg,
			AttrSet:             attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateVoltageSensor")
}

func GetBulkVoltageSensorConfig(fromIdx, count int) (*objects.VoltageSensorConfigGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_VOLTAGE_SENSOR_CONFIG,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkVoltageSensorConfigOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkVoltageSensorConfig")
	}
}

func GetPowerConverterSensorState(name string) (*objects.PowerConverterSensorState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_POWER_CONVERTER_SENSOR_STATE,
		Data: interface{}(&server.GetPowerConverterSensorStateInArgs{
			Name: name,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetPowerConverterSensorStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetPowerConverterSensorState")
	}
}

func GetBulkPowerConverterSensorState(fromIdx, count int) (*objects.PowerConverterSensorStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_POWER_CONVERTER_SENSOR_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkPowerConverterSensorStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkPowerConverterSensorState")
	}
}

func UpdatePowerConverterSensor(oldCfg *objects.PowerConverterSensorConfig, newCfg *objects.PowerConverterSensorConfig, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_POWER_CONVERTER_SENSOR_CONFIG,
		Data: interface{}(&server.UpdatePowerConverterSensorConfigInArgs{
			PowerConverterSensorOldCfg: oldCfg,
			PowerConverterSensorNewCfg: newCfg,
			AttrSet:                    attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdatePowerConverterSensor")
}

func GetBulkPowerConverterSensorConfig(fromIdx, count int) (*objects.PowerConverterSensorConfigGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_POWER_CONVERTER_SENSOR_CONFIG,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkPowerConverterSensorConfigOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkPowerConverterSensorConfig")
	}
}

func GetQsfpState(QsfpId int32) (*objects.QsfpState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_QSFP_STATE,
		Data: interface{}(&server.GetQsfpStateInArgs{
			QsfpId: QsfpId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetQsfpStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetQsfpState")
	}
}

func GetBulkQsfpState(fromIdx, count int) (*objects.QsfpStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_QSFP_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkQsfpStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkQsfpState")
	}
}

func UpdateQsfp(oldCfg *objects.QsfpConfig, newCfg *objects.QsfpConfig, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_QSFP_CONFIG,
		Data: interface{}(&server.UpdateQsfpConfigInArgs{
			QsfpOldCfg: oldCfg,
			QsfpNewCfg: newCfg,
			AttrSet:    attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateQsfp")
}

func GetBulkQsfpConfig(fromIdx, count int) (*objects.QsfpConfigGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_QSFP_CONFIG,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkQsfpConfigOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkQsfpConfig")
	}
}

func GetQsfpChannelState(QsfpId int32, ChannelNum int32) (*objects.QsfpChannelState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_QSFP_CHANNEL_STATE,
		Data: interface{}(&server.GetQsfpChannelStateInArgs{
			QsfpId:     QsfpId,
			ChannelNum: ChannelNum,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetQsfpChannelStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetQsfpChannelState")
	}
}

func GetBulkQsfpChannelState(fromIdx, count int) (*objects.QsfpChannelStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_QSFP_CHANNEL_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkQsfpChannelStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkQsfpChannelState")
	}
}

func UpdateQsfpChannel(oldCfg *objects.QsfpChannelConfig, newCfg *objects.QsfpChannelConfig, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_QSFP_CHANNEL_CONFIG,
		Data: interface{}(&server.UpdateQsfpChannelConfigInArgs{
			QsfpChannelOldCfg: oldCfg,
			QsfpChannelNewCfg: newCfg,
			AttrSet:           attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateQsfpChannel")
}

func GetBulkQsfpChannelConfig(fromIdx, count int) (*objects.QsfpChannelConfigGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_QSFP_CHANNEL_CONFIG,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkQsfpChannelConfigOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkQsfpChannelConfig")
	}
}

func GetPlatformMgmtDeviceState(deviceName string) (*objects.PlatformMgmtDeviceState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_PLATFORM_MGMT_DEVICE_STATE,
		Data: interface{}(&server.GetPlatformMgmtDeviceStateInArgs{
			DeviceName: deviceName,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetPlatformMgmtDeviceStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetPlatformMgmtDeviceState")
	}
}

func GetBulkPlatformMgmtDeviceState(fromIdx, count int) (*objects.PlatformMgmtDeviceStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_PLATFORM_MGMT_DEVICE_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkPlatformMgmtDeviceStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkPlatformMgmtDeviceState")
	}
}

func GetFanSensorPMDataState(name string, class string) (*objects.FanSensorPMState, error) {
	if class != "Class-A" &&
		class != "Class-B" &&
		class != "Class-C" {
		return nil, errors.New("Invalid Class")
	}
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_FAN_SENSOR_PM_STATE,
		Data: interface{}(&server.GetFanSensorPMStateInArgs{
			Name:  name,
			Class: class,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetFanSensorPMStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetFanSensorPMDataState")

	}
}

func GetTempSensorPMDataState(name string, class string) (*objects.TemperatureSensorPMState, error) {
	if class != "Class-A" &&
		class != "Class-B" &&
		class != "Class-C" {
		return nil, errors.New("Invalid Class")
	}
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_TEMPERATURE_SENSOR_PM_STATE,
		Data: interface{}(&server.GetTempSensorPMStateInArgs{
			Name:  name,
			Class: class,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetTempSensorPMStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetTempSensorPMDataState")
	}
}

func GetVoltageSensorPMDataState(name string, class string) (*objects.VoltageSensorPMState, error) {
	if class != "Class-A" &&
		class != "Class-B" &&
		class != "Class-C" {
		return nil, errors.New("Invalid Class")
	}
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_VOLTAGE_SENSOR_PM_STATE,
		Data: interface{}(&server.GetVoltageSensorPMStateInArgs{
			Name:  name,
			Class: class,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetVoltageSensorPMStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetVoltageSensorPMDataState")
	}
}

func GetPowerConverterSensorPMDataState(name string, class string) (*objects.PowerConverterSensorPMState, error) {
	if class != "Class-A" &&
		class != "Class-B" &&
		class != "Class-C" {
		return nil, errors.New("Invalid Class")
	}
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_POWER_CONVERTER_SENSOR_PM_STATE,
		Data: interface{}(&server.GetPowerConverterSensorPMStateInArgs{
			Name:  name,
			Class: class,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetPowerConverterSensorPMStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetPowerConverterSensorPMDataState")
	}
}

func GetQsfpPMDataState(qsfpId int32, resource string, class string) (*objects.QsfpPMState, error) {
	if class != "Class-A" &&
		class != "Class-B" &&
		class != "Class-C" {
		return nil, errors.New("Invalid Class")
	}
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_QSFP_PM_STATE,
		Data: interface{}(&server.GetQsfpPMStateInArgs{
			QsfpId:   qsfpId,
			Resource: resource,
			Class:    class,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetQsfpPMStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetQsfpPMDataState")
	}
}
func GetQsfpChannelPMDataState(qsfpId int32, channelNum int32, resource string, class string) (*objects.QsfpChannelPMState, error) {
	if class != "Class-A" &&
		class != "Class-B" &&
		class != "Class-C" {
		return nil, errors.New("Invalid Class")
	}
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_QSFP_CHANNEL_PM_STATE,
		Data: interface{}(&server.GetQsfpChannelPMStateInArgs{
			QsfpId:     qsfpId,
			ChannelNum: channelNum,
			Resource:   resource,
			Class:      class,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetQsfpChannelPMStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetQsfpChannelPMDataState")
	}
}

func GetPsuState(psuId int32) (*objects.PsuState, error) {
	var obj objects.PsuState

	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_PSU_STATE,
		Data: interface{}(&server.GetPsuStateInArgs{
			PsuId: psuId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetPsuStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during PsuStateGet")
	}

	return &obj, nil
}

func GetBulkPsuState(fromIdx, count int) (*objects.PsuStateGetInfo, error) {
	var obj objects.PsuStateGetInfo

	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_PSU_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkPsuStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkPsuState")
	}

	return &obj, nil
}

func GetLedState(ledId int32) (*objects.LedState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_LED_STATE,
		Data: interface{}(&server.GetLedStateInArgs{
			LedId: ledId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetLedStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetLedState")
	}
}

func GetBulkLedState(fromIdx, count int) (*objects.LedStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_LED_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkLedStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkLedState")
	}
}
