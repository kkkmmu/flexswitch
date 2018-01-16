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
	"fmt"
	//"infra/platformd/objects"
	"infra/platformd/pluginManager"
	"infra/platformd/pluginManager/pluginCommon"
	"strings"
	"utils/dbutils"
	"utils/eventUtils"
	"utils/logging"
)

type PlatformdServer struct {
	dmnName        string
	paramsDir      string
	pluginMgr      *pluginManager.PluginManager
	eventDbHdl     dbutils.DBIntf
	Logger         logging.LoggerIntf
	InitCompleteCh chan bool
	ReqChan        chan *ServerRequest
	ReplyChan      chan interface{}
}

type InitParams struct {
	DmnName     string
	ParamsDir   string
	CfgFileName string
	EventDbHdl  dbutils.DBIntf
	Logger      logging.LoggerIntf
}

func NewPlatformdServer(initParams *InitParams) (*PlatformdServer, error) {
	var svr PlatformdServer

	svr.dmnName = initParams.DmnName
	svr.paramsDir = initParams.ParamsDir
	svr.eventDbHdl = initParams.EventDbHdl
	svr.Logger = initParams.Logger
	svr.InitCompleteCh = make(chan bool)
	svr.ReqChan = make(chan *ServerRequest)
	svr.ReplyChan = make(chan interface{})

	CfgFileInfo, err := parseCfgFile(initParams.CfgFileName)
	if err != nil {
		svr.Logger.Err("Failed to parse platformd config file, using default values for all attributes")
	}
	pluginInitParams := &pluginCommon.PluginInitParams{
		Logger:     svr.Logger,
		PluginName: CfgFileInfo.PluginName,
		IpAddr:     CfgFileInfo.IpAddr,
		Port:       CfgFileInfo.Port,
		EventDbHdl: svr.eventDbHdl,
	}
	svr.pluginMgr, err = pluginManager.NewPluginMgr(pluginInitParams)
	if err != nil {
		return nil, err
	}
	return &svr, err
}

func (svr *PlatformdServer) initServer() error {
	//Initialize plugin layer first
	err := eventUtils.InitEvents(strings.ToUpper(svr.dmnName), svr.eventDbHdl, svr.eventDbHdl, svr.Logger, 1000)
	if err != nil {
		return err
	}
	err = svr.pluginMgr.Init()
	if err != nil {
		return err
	}

	return err
}

func (svr *PlatformdServer) handleRPCRequest(req *ServerRequest) {
	svr.Logger.Info(fmt.Sprintln("Calling handle RPC Request for:", *req))
	switch req.Op {
	case GET_FAN_STATE:
		var retObj GetFanStateOutArgs
		if val, ok := req.Data.(*GetFanStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getFanState(val.FanId)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_FAN_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_FAN_STATE:
		var retObj GetBulkFanStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkFanState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_FAN_CONFIG:
		var retObj GetFanConfigOutArgs
		if val, ok := req.Data.(*GetFanConfigInArgs); ok {
			retObj.Obj, retObj.Err = svr.getFanConfig(val.FanId)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_FAN_CONFIG request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_FAN_CONFIG:
		var retObj GetBulkFanConfigOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkFanConfig(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case UPDATE_FAN_CONFIG:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdateFanConfigInArgs); ok {
			retObj.RetVal, retObj.Err = svr.updateFanConfig(val.FanOldCfg, val.FanNewCfg, val.AttrSet)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_SFP_STATE:
		var retObj GetSfpStateOutArgs
		if val, ok := req.Data.(*GetSfpStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getSfpState(val.SfpId)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_SFP_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_SFP_STATE:
		var retObj GetBulkSfpStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkSfpState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_PLATFORM_STATE:
		var retObj GetPlatformStateOutArgs
		if val, ok := req.Data.(*GetPlatformStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getPlatformState(val.ObjName)
		}
		svr.Logger.Info(fmt.Sprintln("Server GET_PLATFORM_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_PLATFORM_STATE:
		var retObj GetBulkPlatformStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkPlatformState(val.FromIdx, val.Count)
		}
		svr.Logger.Info(fmt.Sprintln("Server GET BULK GET_PLATFORM_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_THERMAL_STATE:
		var retObj GetThermalStateOutArgs
		if val, ok := req.Data.(*GetThermalStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getThermalState(val.ThermalId)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_THERMAL_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_THERMAL_STATE:
		var retObj GetBulkThermalStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkThermalState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_FAN_SENSOR_STATE:
		var retObj GetFanSensorStateOutArgs
		if val, ok := req.Data.(*GetFanSensorStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getFanSensorState(val.Name)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_FAN_SENSOR_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_FAN_SENSOR_STATE:
		var retObj GetBulkFanSensorStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkFanSensorState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_FAN_SENSOR_CONFIG:
		var retObj GetBulkFanSensorConfigOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkFanSensorConfig(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case UPDATE_FAN_SENSOR_CONFIG:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdateFanSensorConfigInArgs); ok {
			retObj.RetVal, retObj.Err = svr.updateFanSensorConfig(val.FanSensorOldCfg, val.FanSensorNewCfg, val.AttrSet)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_TEMPERATURE_SENSOR_STATE:
		var retObj GetTemperatureSensorStateOutArgs
		if val, ok := req.Data.(*GetTemperatureSensorStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getTemperatureSensorState(val.Name)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_TEMPERATURE_SENSOR_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_TEMPERATURE_SENSOR_STATE:
		var retObj GetBulkTemperatureSensorStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkTemperatureSensorState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_TEMPERATURE_SENSOR_CONFIG:
		var retObj GetBulkTemperatureSensorConfigOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkTemperatureSensorConfig(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case UPDATE_TEMPERATURE_SENSOR_CONFIG:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdateTemperatureSensorConfigInArgs); ok {
			retObj.RetVal, retObj.Err = svr.updateTemperatureSensorConfig(val.TemperatureSensorOldCfg, val.TemperatureSensorNewCfg, val.AttrSet)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_VOLTAGE_SENSOR_STATE:
		var retObj GetVoltageSensorStateOutArgs
		if val, ok := req.Data.(*GetVoltageSensorStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getVoltageSensorState(val.Name)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_VOLTAGE_SENSOR_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_VOLTAGE_SENSOR_STATE:
		var retObj GetBulkVoltageSensorStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkVoltageSensorState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_VOLTAGE_SENSOR_CONFIG:
		var retObj GetBulkVoltageSensorConfigOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkVoltageSensorConfig(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case UPDATE_VOLTAGE_SENSOR_CONFIG:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdateVoltageSensorConfigInArgs); ok {
			retObj.RetVal, retObj.Err = svr.updateVoltageSensorConfig(val.VoltageSensorOldCfg, val.VoltageSensorNewCfg, val.AttrSet)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_POWER_CONVERTER_SENSOR_STATE:
		var retObj GetPowerConverterSensorStateOutArgs
		if val, ok := req.Data.(*GetPowerConverterSensorStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getPowerConverterSensorState(val.Name)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_POWER_CONVERTER_SENSOR_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_POWER_CONVERTER_SENSOR_STATE:
		var retObj GetBulkPowerConverterSensorStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkPowerConverterSensorState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_POWER_CONVERTER_SENSOR_CONFIG:
		var retObj GetBulkPowerConverterSensorConfigOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkPowerConverterSensorConfig(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case UPDATE_POWER_CONVERTER_SENSOR_CONFIG:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdatePowerConverterSensorConfigInArgs); ok {
			retObj.RetVal, retObj.Err = svr.updatePowerConverterSensorConfig(val.PowerConverterSensorOldCfg, val.PowerConverterSensorNewCfg, val.AttrSet)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_QSFP_STATE:
		var retObj GetQsfpStateOutArgs
		if val, ok := req.Data.(*GetQsfpStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getQsfpState(val.QsfpId)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_QSFP_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_QSFP_STATE:
		var retObj GetBulkQsfpStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkQsfpState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_QSFP_CONFIG:
		var retObj GetBulkQsfpConfigOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkQsfpConfig(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case UPDATE_QSFP_CONFIG:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdateQsfpConfigInArgs); ok {
			retObj.RetVal, retObj.Err = svr.updateQsfpConfig(val.QsfpOldCfg, val.QsfpNewCfg, val.AttrSet)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_QSFP_CHANNEL_STATE:
		var retObj GetQsfpChannelStateOutArgs
		if val, ok := req.Data.(*GetQsfpChannelStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getQsfpChannelState(val.QsfpId, val.ChannelNum)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_QSFP_CHANNEL_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_QSFP_CHANNEL_STATE:
		var retObj GetBulkQsfpChannelStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkQsfpChannelState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_QSFP_CHANNEL_CONFIG:
		var retObj GetBulkQsfpChannelConfigOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkQsfpChannelConfig(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case UPDATE_QSFP_CHANNEL_CONFIG:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdateQsfpChannelConfigInArgs); ok {
			retObj.RetVal, retObj.Err = svr.updateQsfpChannelConfig(val.QsfpChannelOldCfg, val.QsfpChannelNewCfg, val.AttrSet)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_PLATFORM_MGMT_DEVICE_STATE:
		var retObj GetPlatformMgmtDeviceStateOutArgs
		if val, ok := req.Data.(*GetPlatformMgmtDeviceStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getPlatformMgmtDeviceState(val.DeviceName)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_PLATFORM_MGMT_DEVICE_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_PLATFORM_MGMT_DEVICE_STATE:
		var retObj GetBulkPlatformMgmtDeviceStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkPlatformMgmtDeviceState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_FAN_SENSOR_PM_STATE:
		var retObj GetFanSensorPMStateOutArgs
		if val, ok := req.Data.(*GetFanSensorPMStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getFanSensorPMState(val.Name, val.Class)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_FAN_SENSOR_PM_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_TEMPERATURE_SENSOR_PM_STATE:
		var retObj GetTempSensorPMStateOutArgs
		if val, ok := req.Data.(*GetTempSensorPMStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getTempSensorPMState(val.Name, val.Class)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_TEMPERATURE_SENSOR_PM_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_VOLTAGE_SENSOR_PM_STATE:
		var retObj GetVoltageSensorPMStateOutArgs
		if val, ok := req.Data.(*GetVoltageSensorPMStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getVoltageSensorPMState(val.Name, val.Class)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_VOLTAGE_SENSOR_PM_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_POWER_CONVERTER_SENSOR_PM_STATE:
		var retObj GetPowerConverterSensorPMStateOutArgs
		if val, ok := req.Data.(*GetPowerConverterSensorPMStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getPowerConverterSensorPMState(val.Name, val.Class)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_QSFP_PM_STATE:
		var retObj GetQsfpPMStateOutArgs
		if val, ok := req.Data.(*GetQsfpPMStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getQsfpPMState(val.QsfpId, val.Resource, val.Class)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_QSFP_PM_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_QSFP_CHANNEL_PM_STATE:
		var retObj GetQsfpChannelPMStateOutArgs
		if val, ok := req.Data.(*GetQsfpChannelPMStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getQsfpChannelPMState(val.QsfpId, val.ChannelNum, val.Resource, val.Class)
		}
		//svr.Logger.Info(fmt.Sprintln("Server GET_QSFP_CHANNEL_PM_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_PSU_STATE:
		var retObj GetPsuStateOutArgs
		if val, ok := req.Data.(*GetPsuStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getPsuState(val.PsuId)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_PSU_STATE:
		var retObj GetBulkPsuStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkPsuState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_LED_STATE:
		var retObj GetLedStateOutArgs
		if val, ok := req.Data.(*GetLedStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getLedState(val.LedId)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_LED_STATE:
		var retObj GetBulkLedStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkLedState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	default:
		svr.Logger.Err(fmt.Sprintln("Error : Server recevied unrecognized request - ", req.Op))
	}
}

func (svr *PlatformdServer) Serve() {
	svr.Logger.Info("Server initialization started")
	err := svr.initServer()
	if err != nil {
		panic(err)
	}
	svr.InitCompleteCh <- true
	svr.Logger.Info("Server initialization complete, starting cfg/state listerner")
	for {
		select {
		case req := <-svr.ReqChan:
			svr.Logger.Info(fmt.Sprintln("Server request received - ", *req))
			svr.handleRPCRequest(req)

		}
	}
}
