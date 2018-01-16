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
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"models/events"
	"sync"
	"time"
	"utils/eventUtils"
	"utils/logging"
	"utils/ringBuffer"
)

type QsfpConfig struct {
	AdminState               string
	HigherAlarmTemperature   float64
	HigherAlarmVoltage       float64
	HigherWarningTemperature float64
	HigherWarningVoltage     float64
	LowerAlarmTemperature    float64
	LowerAlarmVoltage        float64
	LowerWarningTemperature  float64
	LowerWarningVoltage      float64
	PMClassAAdminState       string
	PMClassBAdminState       string
	PMClassCAdminState       string
}

type QsfpChannelState struct {
	Present bool
	RXPower float64
	TXPower float64
	TXBias  float64
}

type QsfpChannelConfig struct {
	AdminState           string
	HigherAlarmRXPower   float64
	HigherAlarmTXPower   float64
	HigherAlarmTXBias    float64
	HigherWarningRXPower float64
	HigherWarningTXPower float64
	HigherWarningTXBias  float64
	LowerAlarmRXPower    float64
	LowerAlarmTXPower    float64
	LowerAlarmTXBias     float64
	LowerWarningRXPower  float64
	LowerWarningTXPower  float64
	LowerWarningTXBias   float64
	PMClassAAdminState   string
	PMClassBAdminState   string
	PMClassCAdminState   string
}

type QsfpResource struct {
	QsfpId int32
	ResId  uint8
}

type QsfpChannel struct {
	QsfpId     int32
	ChannelNum uint8
}

type QsfpChannelResource struct {
	QsfpChannel QsfpChannel
	ResId       uint8
}

const (
	qsfpClassAInterval time.Duration = time.Duration(10) * time.Second // Polling Interval 10 sec
	qsfpClassABufSize  int           = 6 * 60 * 24                     //Storage for 24 hrs
	qsfpClassBInterval time.Duration = time.Duration(15) * time.Minute // Polling Interval 15 mins
	qsfpClassBBufSize  int           = 4 * 24                          // Storage for 24 hrs
	qsfpClassCInterval time.Duration = time.Duration(24) * time.Hour   // Polling Interval 24 Hrs
	qsfpClassCBufSize  int           = 365                             // Storage for 365 days
)

const (
	TemperatureRes uint8 = 0
	VoltageRes     uint8 = 1
	MaxNumQsfpRes  uint8 = 2 // Should be always last
)
const (
	RXPowerRes           uint8 = 0
	TXPowerRes           uint8 = 1
	TXBiasRes            uint8 = 2
	MaxNumQsfpChannelRes uint8 = 3 // Should be always last
)

const (
	MaxNumOfQsfpChannel uint8 = 4
)

type QsfpEventStatus [MaxNumQsfpRes]EventStatus
type QsfpChannelEventStatus [MaxNumQsfpChannelRes]EventStatus

type QsfpEventData struct {
	Value float64
}

type QsfpChannelEventData struct {
	Value float64
}

func getQsfpResourcId(res string) (uint8, error) {
	switch res {
	case "Temperature":
		return 0, nil
	case "Voltage":
		return 1, nil
	default:
		return 0, errors.New("Invalid Resource Name")
	}
	return 0, errors.New("Invalid Resource Name")
}

func getQsfpChannelResourcId(res string) (uint8, error) {
	switch res {
	case "RXPower":
		return 0, nil
	case "TXPower":
		return 1, nil
	case "TXBias":
		return 2, nil
	default:
		return 0, errors.New("Invalid Resource Name")
	}
	return 0, errors.New("Invalid Resource Name")
}

type QsfpManager struct {
	logger                         logging.LoggerIntf
	plugin                         PluginIntf
	stateMutex                     sync.RWMutex
	numOfQsfps                     int32
	qsfpConfigMutex                sync.RWMutex
	qsfpConfigDB                   map[int32]QsfpConfig
	classAPMTimer                  *time.Timer
	qsfpClassAMutex                sync.RWMutex
	qsfpClassAPM                   map[QsfpResource]*ringBuffer.RingBuffer
	classBPMTimer                  *time.Timer
	qsfpClassBMutex                sync.RWMutex
	qsfpClassBPM                   map[QsfpResource]*ringBuffer.RingBuffer
	classCPMTimer                  *time.Timer
	qsfpClassCMutex                sync.RWMutex
	qsfpClassCPM                   map[QsfpResource]*ringBuffer.RingBuffer
	qsfpChannelConfigMutex         sync.RWMutex
	qsfpChannelConfigDB            map[QsfpChannel]QsfpChannelConfig
	qsfpChannelClassAMutex         sync.RWMutex
	qsfpChannelClassAPM            map[QsfpChannelResource]*ringBuffer.RingBuffer
	qsfpChannelClassBMutex         sync.RWMutex
	qsfpChannelClassBPM            map[QsfpChannelResource]*ringBuffer.RingBuffer
	qsfpChannelClassCMutex         sync.RWMutex
	qsfpChannelClassCPM            map[QsfpChannelResource]*ringBuffer.RingBuffer
	qsfpEventMsgStatusMutex        sync.RWMutex
	qsfpEventMsgStatus             map[QsfpResource]EventStatus
	qsfpChannelEventMsgStatusMutex sync.RWMutex
	qsfpChannelEventMsgStatus      map[QsfpChannelResource]EventStatus
	//	qsfpStatus                []bool
}

var QsfpMgr QsfpManager

func (qMgr *QsfpManager) Init(logger logging.LoggerIntf, plugin PluginIntf) {
	qMgr.logger = logger
	qMgr.plugin = plugin
	numOfQsfps := qMgr.plugin.GetMaxNumOfQsfp()
	qMgr.numOfQsfps = int32(numOfQsfps)

	qMgr.qsfpClassAPM = make(map[QsfpResource]*ringBuffer.RingBuffer, numOfQsfps*int(MaxNumQsfpRes))
	qMgr.qsfpClassBPM = make(map[QsfpResource]*ringBuffer.RingBuffer, numOfQsfps*int(MaxNumQsfpRes))
	qMgr.qsfpClassCPM = make(map[QsfpResource]*ringBuffer.RingBuffer, numOfQsfps*int(MaxNumQsfpRes))
	qMgr.qsfpEventMsgStatus = make(map[QsfpResource]EventStatus, numOfQsfps*int(MaxNumQsfpRes))
	qMgr.qsfpChannelClassAPM = make(map[QsfpChannelResource]*ringBuffer.RingBuffer, numOfQsfps*int(MaxNumOfQsfpChannel)*int(MaxNumQsfpChannelRes))
	qMgr.qsfpChannelClassBPM = make(map[QsfpChannelResource]*ringBuffer.RingBuffer, numOfQsfps*int(MaxNumOfQsfpChannel)*int(MaxNumQsfpChannelRes))
	qMgr.qsfpChannelClassCPM = make(map[QsfpChannelResource]*ringBuffer.RingBuffer, numOfQsfps*int(MaxNumOfQsfpChannel)*int(MaxNumQsfpChannelRes))
	qMgr.qsfpChannelEventMsgStatus = make(map[QsfpChannelResource]EventStatus, numOfQsfps*int(MaxNumOfQsfpChannel)*int(MaxNumQsfpChannelRes))
	//	qMgr.qsfpStatus = make([]bool, numOfQsfps)

	qMgr.qsfpConfigMutex.Lock()
	qMgr.qsfpChannelConfigMutex.Lock()
	qMgr.qsfpConfigDB = make(map[int32]QsfpConfig, numOfQsfps)
	qMgr.qsfpChannelConfigDB = make(map[QsfpChannel]QsfpChannelConfig, numOfQsfps*int(MaxNumOfQsfpChannel))
	for id := 1; id <= numOfQsfps; id++ {
		qsfpId := int32(id)
		qsfpCfgEnt, _ := qMgr.qsfpConfigDB[qsfpId]
		qsfpCfgEnt.AdminState = "Disable"
		qsfpCfgEnt.HigherAlarmTemperature = 100.0
		qsfpCfgEnt.HigherAlarmVoltage = 10.0
		qsfpCfgEnt.HigherWarningTemperature = 100.0
		qsfpCfgEnt.HigherWarningVoltage = 10.0
		qsfpCfgEnt.LowerAlarmTemperature = -100.0
		qsfpCfgEnt.LowerAlarmVoltage = -10.0
		qsfpCfgEnt.LowerWarningTemperature = -100.0
		qsfpCfgEnt.LowerWarningVoltage = -10.0
		qsfpCfgEnt.PMClassAAdminState = "Disable"
		qsfpCfgEnt.PMClassBAdminState = "Disable"
		qsfpCfgEnt.PMClassCAdminState = "Disable"
		qMgr.qsfpConfigDB[qsfpId] = qsfpCfgEnt
		for ch := 1; ch <= int(MaxNumOfQsfpChannel); ch++ {
			qsfpChannel := QsfpChannel{
				QsfpId:     qsfpId,
				ChannelNum: uint8(ch),
			}
			qsfpChCfgEnt, _ := qMgr.qsfpChannelConfigDB[qsfpChannel]
			qsfpChCfgEnt.AdminState = "Disable"
			qsfpChCfgEnt.HigherAlarmRXPower = 100.0
			qsfpChCfgEnt.HigherAlarmTXPower = 100.0
			qsfpChCfgEnt.HigherAlarmTXBias = 100.0
			qsfpChCfgEnt.HigherWarningRXPower = 100.0
			qsfpChCfgEnt.HigherWarningTXPower = 100.0
			qsfpChCfgEnt.HigherWarningTXBias = 100.0
			qsfpChCfgEnt.LowerAlarmRXPower = -100.0
			qsfpChCfgEnt.LowerAlarmTXPower = -100.0
			qsfpChCfgEnt.LowerAlarmTXBias = -100.0
			qsfpChCfgEnt.LowerWarningRXPower = -100.0
			qsfpChCfgEnt.LowerWarningTXPower = -100.0
			qsfpChCfgEnt.LowerWarningTXBias = -100.0
			qsfpChCfgEnt.PMClassAAdminState = "Disable"
			qsfpChCfgEnt.PMClassBAdminState = "Disable"
			qsfpChCfgEnt.PMClassCAdminState = "Disable"
			qMgr.qsfpChannelConfigDB[qsfpChannel] = qsfpChCfgEnt
		}
	}
	qMgr.qsfpChannelConfigMutex.Unlock()
	qMgr.qsfpConfigMutex.Unlock()
	qMgr.StartQsfpPM()
	qMgr.logger.Info("Qsfp Manager Init()")
}

func (qMgr *QsfpManager) Deinit() {
	qMgr.logger.Info("Fan Manager Deinit()")
}

func (qMgr *QsfpManager) GetQsfpState(id int32) (*objects.QsfpState, error) {
	var qsfpObj objects.QsfpState
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}

	_, exist := qMgr.qsfpConfigDB[id]
	if !exist {
		return nil, errors.New("Invalid QsfpId")
	}
	qMgr.stateMutex.Lock()
	qsfpState, err := qMgr.plugin.GetQsfpState(id)
	qMgr.stateMutex.Unlock()
	if err != nil {
		qsfpObj.QsfpId = id
		qsfpObj.Present = false
	} else {
		qsfpObj.QsfpId = id
		qsfpObj.Present = true
		qsfpObj.VendorName = qsfpState.VendorName
		qsfpObj.VendorOUI = qsfpState.VendorOUI
		qsfpObj.VendorPartNumber = qsfpState.VendorPartNumber
		qsfpObj.VendorRevision = qsfpState.VendorRevision
		qsfpObj.VendorSerialNumber = qsfpState.VendorSerialNumber
		qsfpObj.DataCode = qsfpState.DataCode
		qsfpObj.Temperature = qsfpState.Temperature
		qsfpObj.Voltage = qsfpState.Voltage
	}
	return &qsfpObj, nil
}

func (qMgr *QsfpManager) GetQsfpChannelState(id int32, channelNum int32) (*objects.QsfpChannelState, error) {
	var qsfpChannelObj objects.QsfpChannelState
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if channelNum > int32(MaxNumOfQsfpChannel) {
		return nil, errors.New("Invalid ChannelNum")
	}
	qsfpChannel := QsfpChannel{
		QsfpId:     id,
		ChannelNum: uint8(channelNum),
	}
	_, exist := qMgr.qsfpChannelConfigDB[qsfpChannel]
	if !exist {
		return nil, errors.New("Invalid QsfpId")
	}

	qMgr.stateMutex.Lock()
	qsfpState, err := qMgr.plugin.GetQsfpState(id)
	qMgr.stateMutex.Unlock()
	if err != nil {
		qsfpChannelObj.QsfpId = id
		qsfpChannelObj.ChannelNum = channelNum
		qsfpChannelObj.Present = false
	} else {
		qsfpChannelObj.QsfpId = id
		qsfpChannelObj.ChannelNum = channelNum
		qsfpChannelObj.Present = true
		qsfpChannelObj.RXPower = qsfpState.RXPower[channelNum-1]
		qsfpChannelObj.TXPower = qsfpState.TXPower[channelNum-1]
		qsfpChannelObj.TXBias = qsfpState.TXBias[channelNum-1]
	}

	return &qsfpChannelObj, nil
}

func (qMgr *QsfpManager) GetBulkQsfpState(fromIdx int, cnt int) (*objects.QsfpStateGetInfo, error) {
	var retObj objects.QsfpStateGetInfo
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= int(qMgr.numOfQsfps) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > int(qMgr.numOfQsfps) {
		retObj.EndIdx = int(qMgr.numOfQsfps)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = int(qMgr.numOfQsfps) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		obj, err := qMgr.GetQsfpState(int32(idx + 1))
		if err != nil {
			qMgr.logger.Err("Error getting the qsfp state for QsfpId:", idx+1)
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (qMgr *QsfpManager) GetBulkQsfpChannelState(fromIdx int, cnt int) (*objects.QsfpChannelStateGetInfo, error) {
	var retObj objects.QsfpChannelStateGetInfo
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= int(qMgr.numOfQsfps)*int(MaxNumOfQsfpChannel) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > int(qMgr.numOfQsfps)*int(MaxNumOfQsfpChannel) {
		retObj.EndIdx = int(qMgr.numOfQsfps) * int(MaxNumOfQsfpChannel)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = (int(qMgr.numOfQsfps) * int(MaxNumOfQsfpChannel)) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		qsfpId := int32((idx / int(MaxNumOfQsfpChannel)) + 1)
		chNum := int32((idx % int(MaxNumOfQsfpChannel)) + 1)
		obj, err := qMgr.GetQsfpChannelState(qsfpId, chNum)
		if err != nil {
			qMgr.logger.Err("Error getting the qsfp channel state for QsfpId:", qsfpId, "Channel Number:", chNum)
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (qMgr *QsfpManager) GetQsfpConfig(id int32) (*objects.QsfpConfig, error) {
	var obj objects.QsfpConfig
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	qMgr.qsfpConfigMutex.RLock()
	qsfpCfgEnt, exist := qMgr.qsfpConfigDB[id]
	if !exist {
		qMgr.qsfpConfigMutex.RUnlock()
		return nil, errors.New("Invalid QsfpId")
	}

	obj.QsfpId = id
	obj.AdminState = qsfpCfgEnt.AdminState
	obj.HigherAlarmTemperature = qsfpCfgEnt.HigherAlarmTemperature
	obj.HigherAlarmVoltage = qsfpCfgEnt.HigherAlarmVoltage
	obj.HigherWarningTemperature = qsfpCfgEnt.HigherWarningTemperature
	obj.HigherWarningVoltage = qsfpCfgEnt.HigherWarningVoltage
	obj.LowerAlarmTemperature = qsfpCfgEnt.LowerAlarmTemperature
	obj.LowerAlarmVoltage = qsfpCfgEnt.LowerAlarmVoltage
	obj.LowerWarningTemperature = qsfpCfgEnt.LowerWarningTemperature
	obj.LowerWarningVoltage = qsfpCfgEnt.LowerWarningVoltage
	obj.PMClassAAdminState = qsfpCfgEnt.PMClassAAdminState
	obj.PMClassBAdminState = qsfpCfgEnt.PMClassBAdminState
	obj.PMClassCAdminState = qsfpCfgEnt.PMClassCAdminState
	qMgr.qsfpConfigMutex.RUnlock()

	return &obj, nil
}

func (qMgr *QsfpManager) GetQsfpChannelConfig(id int32, channelNum int32) (*objects.QsfpChannelConfig, error) {
	var obj objects.QsfpChannelConfig
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if channelNum > int32(MaxNumOfQsfpChannel) {
		return nil, errors.New("Invalid ChannelNum")
	}
	qsfpChannel := QsfpChannel{
		QsfpId:     id,
		ChannelNum: uint8(channelNum),
	}
	qMgr.qsfpChannelConfigMutex.RLock()
	qsfpChannelCfgEnt, exist := qMgr.qsfpChannelConfigDB[qsfpChannel]
	if !exist {
		qMgr.qsfpChannelConfigMutex.RUnlock()
		return nil, errors.New("Invalid QsfpId")
	}
	obj.QsfpId = id
	obj.ChannelNum = channelNum
	obj.AdminState = qsfpChannelCfgEnt.AdminState
	obj.HigherAlarmRXPower = qsfpChannelCfgEnt.HigherAlarmRXPower
	obj.HigherAlarmTXPower = qsfpChannelCfgEnt.HigherAlarmTXPower
	obj.HigherAlarmTXBias = qsfpChannelCfgEnt.HigherAlarmTXBias
	obj.HigherWarningRXPower = qsfpChannelCfgEnt.HigherWarningRXPower
	obj.HigherWarningTXPower = qsfpChannelCfgEnt.HigherWarningTXPower
	obj.HigherWarningTXBias = qsfpChannelCfgEnt.HigherWarningTXBias
	obj.LowerAlarmRXPower = qsfpChannelCfgEnt.LowerAlarmRXPower
	obj.LowerAlarmTXPower = qsfpChannelCfgEnt.LowerAlarmTXPower
	obj.LowerAlarmTXBias = qsfpChannelCfgEnt.LowerAlarmTXBias
	obj.LowerWarningRXPower = qsfpChannelCfgEnt.LowerWarningRXPower
	obj.LowerWarningTXPower = qsfpChannelCfgEnt.LowerWarningTXPower
	obj.LowerWarningTXBias = qsfpChannelCfgEnt.LowerWarningTXBias
	obj.PMClassAAdminState = qsfpChannelCfgEnt.PMClassAAdminState
	obj.PMClassBAdminState = qsfpChannelCfgEnt.PMClassBAdminState
	obj.PMClassCAdminState = qsfpChannelCfgEnt.PMClassCAdminState
	qMgr.qsfpChannelConfigMutex.RUnlock()
	return &obj, nil
}

func (qMgr *QsfpManager) GetBulkQsfpConfig(fromIdx int, cnt int) (*objects.QsfpConfigGetInfo, error) {
	var retObj objects.QsfpConfigGetInfo
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= int(qMgr.numOfQsfps) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > int(qMgr.numOfQsfps) {
		retObj.EndIdx = int(qMgr.numOfQsfps)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = int(qMgr.numOfQsfps) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		obj, err := qMgr.GetQsfpConfig(int32(idx + 1))
		if err != nil {
			qMgr.logger.Err("Error getting the Qsfp Config for QsfpId:", idx+1)
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (qMgr *QsfpManager) GetBulkQsfpChannelConfig(fromIdx int, cnt int) (*objects.QsfpChannelConfigGetInfo, error) {
	var retObj objects.QsfpChannelConfigGetInfo
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= int(qMgr.numOfQsfps)*int(MaxNumOfQsfpChannel) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > int(qMgr.numOfQsfps)*int(MaxNumOfQsfpChannel) {
		retObj.EndIdx = int(qMgr.numOfQsfps) * int(MaxNumOfQsfpChannel)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = (int(qMgr.numOfQsfps) * int(MaxNumOfQsfpChannel)) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		qsfpId := int32((idx / int(MaxNumOfQsfpChannel)) + 1)
		chNum := int32((idx % int(MaxNumOfQsfpChannel)) + 1)
		obj, err := qMgr.GetQsfpChannelConfig(qsfpId, chNum)
		if err != nil {
			qMgr.logger.Err("Error getting the qsfp channel config for QsfpId:", qsfpId, "Channel Number:", chNum)
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func genQsfpUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0

	if attrset == nil {
		mask = objects.QSFP_UPDATE_ADMIN_STATE |
			objects.QSFP_UPDATE_HIGHER_ALARM_TEMPERATURE |
			objects.QSFP_UPDATE_HIGHER_ALARM_VOLTAGE |
			objects.QSFP_UPDATE_HIGHER_WARN_TEMPERATURE |
			objects.QSFP_UPDATE_HIGHER_WARN_VOLTAGE |
			objects.QSFP_UPDATE_LOWER_ALARM_TEMPERATURE |
			objects.QSFP_UPDATE_LOWER_ALARM_VOLTAGE |
			objects.QSFP_UPDATE_LOWER_WARN_TEMPERATURE |
			objects.QSFP_UPDATE_LOWER_WARN_VOLTAGE |
			objects.QSFP_UPDATE_PM_CLASS_A_ADMIN_STATE |
			objects.QSFP_UPDATE_PM_CLASS_B_ADMIN_STATE |
			objects.QSFP_UPDATE_PM_CLASS_C_ADMIN_STATE
	} else {
		for idx, val := range attrset {
			if true == val {
				switch idx {
				case 0:
					//QSFP Id
				case 1:
					mask |= objects.QSFP_UPDATE_ADMIN_STATE
				case 2:
					mask |= objects.QSFP_UPDATE_HIGHER_ALARM_TEMPERATURE
				case 3:
					mask |= objects.QSFP_UPDATE_HIGHER_ALARM_VOLTAGE
				case 4:
					mask |= objects.QSFP_UPDATE_HIGHER_WARN_TEMPERATURE
				case 5:
					mask |= objects.QSFP_UPDATE_HIGHER_WARN_VOLTAGE
				case 6:
					mask |= objects.QSFP_UPDATE_LOWER_ALARM_TEMPERATURE
				case 7:
					mask |= objects.QSFP_UPDATE_LOWER_ALARM_VOLTAGE
				case 8:
					mask |= objects.QSFP_UPDATE_LOWER_WARN_TEMPERATURE
				case 9:
					mask |= objects.QSFP_UPDATE_LOWER_WARN_VOLTAGE
				case 10:
					mask |= objects.QSFP_UPDATE_PM_CLASS_A_ADMIN_STATE
				case 11:
					mask |= objects.QSFP_UPDATE_PM_CLASS_B_ADMIN_STATE
				case 12:
					mask |= objects.QSFP_UPDATE_PM_CLASS_C_ADMIN_STATE
				}
			}
		}
	}
	return mask
}

func (qMgr *QsfpManager) UpdateQsfpConfig(oldCfg *objects.QsfpConfig, newCfg *objects.QsfpConfig, attrset []bool) (bool, error) {
	if qMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	qMgr.qsfpConfigMutex.Lock()
	qsfpCfgEnt, exist := qMgr.qsfpConfigDB[newCfg.QsfpId]
	if !exist {
		qMgr.qsfpConfigMutex.Unlock()
		return false, errors.New("Invalid QsfpId")
	}
	var cfgEnt QsfpConfig
	mask := genQsfpUpdateMask(attrset)
	if mask&objects.QSFP_UPDATE_ADMIN_STATE == objects.QSFP_UPDATE_ADMIN_STATE {
		if newCfg.AdminState != "Enable" && newCfg.AdminState != "Disable" {
			qMgr.qsfpConfigMutex.Unlock()
			return false, errors.New("Invalid AdminState Value")
		}

		cfgEnt.AdminState = newCfg.AdminState
	} else {
		cfgEnt.AdminState = qsfpCfgEnt.AdminState
	}
	if mask&objects.QSFP_UPDATE_HIGHER_ALARM_TEMPERATURE == objects.QSFP_UPDATE_HIGHER_ALARM_TEMPERATURE {
		cfgEnt.HigherAlarmTemperature = newCfg.HigherAlarmTemperature
	} else {
		cfgEnt.HigherAlarmTemperature = qsfpCfgEnt.HigherAlarmTemperature
	}
	if mask&objects.QSFP_UPDATE_HIGHER_ALARM_VOLTAGE == objects.QSFP_UPDATE_HIGHER_ALARM_VOLTAGE {
		cfgEnt.HigherAlarmVoltage = newCfg.HigherAlarmVoltage
	} else {
		cfgEnt.HigherAlarmVoltage = qsfpCfgEnt.HigherAlarmVoltage
	}
	if mask&objects.QSFP_UPDATE_HIGHER_WARN_TEMPERATURE == objects.QSFP_UPDATE_HIGHER_WARN_TEMPERATURE {
		cfgEnt.HigherWarningTemperature = newCfg.HigherWarningTemperature
	} else {
		cfgEnt.HigherWarningTemperature = qsfpCfgEnt.HigherWarningTemperature
	}
	if mask&objects.QSFP_UPDATE_HIGHER_WARN_VOLTAGE == objects.QSFP_UPDATE_HIGHER_WARN_VOLTAGE {
		cfgEnt.HigherWarningVoltage = newCfg.HigherWarningVoltage
	} else {
		cfgEnt.HigherWarningVoltage = qsfpCfgEnt.HigherWarningVoltage
	}
	if mask&objects.QSFP_UPDATE_LOWER_ALARM_TEMPERATURE == objects.QSFP_UPDATE_LOWER_ALARM_TEMPERATURE {
		cfgEnt.LowerAlarmTemperature = newCfg.LowerAlarmTemperature
	} else {
		cfgEnt.LowerAlarmTemperature = qsfpCfgEnt.LowerAlarmTemperature
	}
	if mask&objects.QSFP_UPDATE_LOWER_ALARM_VOLTAGE == objects.QSFP_UPDATE_LOWER_ALARM_VOLTAGE {
		cfgEnt.LowerAlarmVoltage = newCfg.LowerAlarmVoltage
	} else {
		cfgEnt.LowerAlarmVoltage = qsfpCfgEnt.LowerAlarmVoltage
	}
	if mask&objects.QSFP_UPDATE_LOWER_WARN_TEMPERATURE == objects.QSFP_UPDATE_LOWER_WARN_TEMPERATURE {
		cfgEnt.LowerWarningTemperature = newCfg.LowerWarningTemperature
	} else {
		cfgEnt.LowerWarningTemperature = qsfpCfgEnt.LowerWarningTemperature
	}
	if mask&objects.QSFP_UPDATE_LOWER_WARN_VOLTAGE == objects.QSFP_UPDATE_LOWER_WARN_VOLTAGE {
		cfgEnt.LowerWarningVoltage = newCfg.LowerWarningVoltage
	} else {
		cfgEnt.LowerWarningVoltage = qsfpCfgEnt.LowerWarningVoltage
	}
	if mask&objects.QSFP_UPDATE_PM_CLASS_A_ADMIN_STATE == objects.QSFP_UPDATE_PM_CLASS_A_ADMIN_STATE {
		if newCfg.PMClassAAdminState != "Enable" && newCfg.PMClassAAdminState != "Disable" {
			qMgr.qsfpConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassAAdminState Value")
		}
		cfgEnt.PMClassAAdminState = newCfg.PMClassAAdminState
	} else {
		cfgEnt.PMClassAAdminState = qsfpCfgEnt.PMClassAAdminState
	}
	if mask&objects.QSFP_UPDATE_PM_CLASS_B_ADMIN_STATE == objects.QSFP_UPDATE_PM_CLASS_B_ADMIN_STATE {
		if newCfg.PMClassBAdminState != "Enable" && newCfg.PMClassBAdminState != "Disable" {
			qMgr.qsfpConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassBAdminState Value")
		}
		cfgEnt.PMClassBAdminState = newCfg.PMClassBAdminState
	} else {
		cfgEnt.PMClassBAdminState = qsfpCfgEnt.PMClassBAdminState
	}
	if mask&objects.QSFP_UPDATE_PM_CLASS_C_ADMIN_STATE == objects.QSFP_UPDATE_PM_CLASS_C_ADMIN_STATE {
		if newCfg.PMClassCAdminState != "Enable" && newCfg.PMClassCAdminState != "Disable" {
			qMgr.qsfpConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassCAdminState Value")
		}
		cfgEnt.PMClassCAdminState = newCfg.PMClassCAdminState
	} else {
		cfgEnt.PMClassCAdminState = qsfpCfgEnt.PMClassCAdminState
	}

	if !(cfgEnt.HigherAlarmTemperature >= cfgEnt.HigherWarningTemperature &&
		cfgEnt.HigherWarningTemperature > cfgEnt.LowerWarningTemperature &&
		cfgEnt.LowerWarningTemperature >= cfgEnt.LowerAlarmTemperature) {
		qMgr.qsfpConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, please verify the thresholds")
	}

	if !(cfgEnt.HigherAlarmVoltage >= cfgEnt.HigherWarningVoltage &&
		cfgEnt.HigherWarningVoltage > cfgEnt.LowerWarningVoltage &&
		cfgEnt.LowerWarningVoltage >= cfgEnt.LowerAlarmVoltage) {
		qMgr.qsfpConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, please verify the thresholds")
	}

	if qsfpCfgEnt.AdminState != cfgEnt.AdminState {
		if cfgEnt.AdminState == "Disable" {
			qMgr.clearExistingQsfpFaults(newCfg.QsfpId)
			qMgr.logger.Info("Clear all the existing Faults")
		}
	}
	if qsfpCfgEnt.PMClassAAdminState != cfgEnt.PMClassAAdminState {
		if cfgEnt.PMClassAAdminState == "Disable" {
			// Flush PM RingBuffer
			qMgr.flushHistoricQsfpPM(newCfg.QsfpId, "Class-A")
			qMgr.logger.Info("Flush Class A PM Ring buffer")
		}
	}
	if qsfpCfgEnt.PMClassBAdminState != cfgEnt.PMClassBAdminState {
		if cfgEnt.PMClassBAdminState == "Disable" {
			// Flush PM RingBuffer
			qMgr.flushHistoricQsfpPM(newCfg.QsfpId, "Class-B")
			qMgr.logger.Info("Flush Class B PM Ring buffer")
		}
	}
	if qsfpCfgEnt.PMClassCAdminState != cfgEnt.PMClassCAdminState {
		if cfgEnt.PMClassCAdminState == "Disable" {
			// Flush PM RingBuffer
			qMgr.flushHistoricQsfpPM(newCfg.QsfpId, "Class-C")
			qMgr.logger.Info("Flush Class C PM Ring buffer")
		}
	}
	qMgr.qsfpConfigDB[newCfg.QsfpId] = cfgEnt
	qMgr.qsfpConfigMutex.Unlock()

	return true, nil
}

func genQsfpChannelUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0

	if attrset == nil {
		mask = objects.QSFP_CHANNEL_UPDATE_ADMIN_STATE |
			objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_RX_POWER |
			objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_TX_POWER |
			objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_TX_BIAS |
			objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_RX_POWER |
			objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_TX_POWER |
			objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_TX_BIAS |
			objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_RX_POWER |
			objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_TX_POWER |
			objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_TX_BIAS |
			objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_RX_POWER |
			objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_TX_POWER |
			objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_TX_BIAS |
			objects.QSFP_CHANNEL_UPDATE_PM_CLASS_A_ADMIN_STATE |
			objects.QSFP_CHANNEL_UPDATE_PM_CLASS_B_ADMIN_STATE |
			objects.QSFP_CHANNEL_UPDATE_PM_CLASS_C_ADMIN_STATE
	} else {
		for idx, val := range attrset {
			if true == val {
				switch idx {
				case 0:
					//QSFP Id
				case 1:
					// Channel Number
				case 2:
					mask |= objects.QSFP_CHANNEL_UPDATE_ADMIN_STATE
				case 3:
					mask |= objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_RX_POWER
				case 4:
					mask |= objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_TX_POWER
				case 5:
					mask |= objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_TX_BIAS
				case 6:
					mask |= objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_RX_POWER
				case 7:
					mask |= objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_TX_POWER
				case 8:
					mask |= objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_TX_BIAS
				case 9:
					mask |= objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_RX_POWER
				case 10:
					mask |= objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_TX_POWER
				case 11:
					mask |= objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_TX_BIAS
				case 12:
					mask |= objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_RX_POWER
				case 13:
					mask |= objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_TX_POWER
				case 14:
					mask |= objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_TX_BIAS
				case 15:
					mask |= objects.QSFP_CHANNEL_UPDATE_PM_CLASS_A_ADMIN_STATE
				case 16:
					mask |= objects.QSFP_CHANNEL_UPDATE_PM_CLASS_B_ADMIN_STATE
				case 17:
					mask |= objects.QSFP_CHANNEL_UPDATE_PM_CLASS_C_ADMIN_STATE
				}
			}
		}
	}
	return mask
}

func (qMgr *QsfpManager) UpdateQsfpChannelConfig(oldCfg *objects.QsfpChannelConfig, newCfg *objects.QsfpChannelConfig, attrset []bool) (bool, error) {
	if qMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	if newCfg.ChannelNum > int32(MaxNumOfQsfpChannel) {
		return false, errors.New("Invalid ChannelNum")
	}
	qsfpChannel := QsfpChannel{
		QsfpId:     newCfg.QsfpId,
		ChannelNum: uint8(newCfg.ChannelNum),
	}
	qMgr.qsfpChannelConfigMutex.Lock()
	qsfpChannelCfgEnt, exist := qMgr.qsfpChannelConfigDB[qsfpChannel]
	if !exist {
		qMgr.qsfpChannelConfigMutex.Unlock()
		return false, errors.New("Invalid QsfpId")
	}
	var cfgEnt QsfpChannelConfig
	mask := genQsfpChannelUpdateMask(attrset)
	if mask&objects.QSFP_CHANNEL_UPDATE_ADMIN_STATE == objects.QSFP_CHANNEL_UPDATE_ADMIN_STATE {
		if newCfg.AdminState != "Enable" && newCfg.AdminState != "Disable" {
			qMgr.qsfpChannelConfigMutex.Unlock()
			return false, errors.New("Invalid AdminState Value")
		}

		cfgEnt.AdminState = newCfg.AdminState
	} else {
		cfgEnt.AdminState = qsfpChannelCfgEnt.AdminState
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_RX_POWER == objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_RX_POWER {
		cfgEnt.HigherAlarmRXPower = newCfg.HigherAlarmRXPower
	} else {
		cfgEnt.HigherAlarmRXPower = qsfpChannelCfgEnt.HigherAlarmRXPower
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_TX_POWER == objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_TX_POWER {
		cfgEnt.HigherAlarmTXPower = newCfg.HigherAlarmTXPower
	} else {
		cfgEnt.HigherAlarmTXPower = qsfpChannelCfgEnt.HigherAlarmTXPower
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_TX_BIAS == objects.QSFP_CHANNEL_UPDATE_HIGHER_ALARM_TX_BIAS {
		cfgEnt.HigherAlarmTXBias = newCfg.HigherAlarmTXBias
	} else {
		cfgEnt.HigherAlarmTXBias = qsfpChannelCfgEnt.HigherAlarmTXBias
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_RX_POWER == objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_RX_POWER {
		cfgEnt.HigherWarningRXPower = newCfg.HigherWarningRXPower
	} else {
		cfgEnt.HigherWarningRXPower = qsfpChannelCfgEnt.HigherWarningRXPower
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_TX_POWER == objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_TX_POWER {
		cfgEnt.HigherWarningTXPower = newCfg.HigherWarningTXPower
	} else {
		cfgEnt.HigherWarningTXPower = qsfpChannelCfgEnt.HigherWarningTXPower
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_TX_BIAS == objects.QSFP_CHANNEL_UPDATE_HIGHER_WARN_TX_BIAS {
		cfgEnt.HigherWarningTXBias = newCfg.HigherWarningTXBias
	} else {
		cfgEnt.HigherWarningTXBias = qsfpChannelCfgEnt.HigherWarningTXBias
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_RX_POWER == objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_RX_POWER {
		cfgEnt.LowerAlarmRXPower = newCfg.LowerAlarmRXPower
	} else {
		cfgEnt.LowerAlarmRXPower = qsfpChannelCfgEnt.LowerAlarmRXPower
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_TX_POWER == objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_TX_POWER {
		cfgEnt.LowerAlarmTXPower = newCfg.LowerAlarmTXPower
	} else {
		cfgEnt.LowerAlarmTXPower = qsfpChannelCfgEnt.LowerAlarmTXPower
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_TX_BIAS == objects.QSFP_CHANNEL_UPDATE_LOWER_ALARM_TX_BIAS {
		cfgEnt.LowerAlarmTXBias = newCfg.LowerAlarmTXBias
	} else {
		cfgEnt.LowerAlarmTXBias = qsfpChannelCfgEnt.LowerAlarmTXBias
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_RX_POWER == objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_RX_POWER {
		cfgEnt.LowerWarningRXPower = newCfg.LowerWarningRXPower
	} else {
		cfgEnt.LowerWarningRXPower = qsfpChannelCfgEnt.LowerWarningRXPower
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_TX_POWER == objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_TX_POWER {
		cfgEnt.LowerWarningTXPower = newCfg.LowerWarningTXPower
	} else {
		cfgEnt.LowerWarningTXPower = qsfpChannelCfgEnt.LowerWarningTXPower
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_TX_BIAS == objects.QSFP_CHANNEL_UPDATE_LOWER_WARN_TX_BIAS {
		cfgEnt.LowerWarningTXBias = newCfg.LowerWarningTXBias
	} else {
		cfgEnt.LowerWarningTXBias = qsfpChannelCfgEnt.LowerWarningTXBias
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_PM_CLASS_A_ADMIN_STATE == objects.QSFP_CHANNEL_UPDATE_PM_CLASS_A_ADMIN_STATE {
		if newCfg.PMClassAAdminState != "Enable" && newCfg.PMClassAAdminState != "Disable" {
			qMgr.qsfpChannelConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassAAdminState Value")
		}
		cfgEnt.PMClassAAdminState = newCfg.PMClassAAdminState
	} else {
		cfgEnt.PMClassAAdminState = qsfpChannelCfgEnt.PMClassAAdminState
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_PM_CLASS_B_ADMIN_STATE == objects.QSFP_CHANNEL_UPDATE_PM_CLASS_B_ADMIN_STATE {
		if newCfg.PMClassBAdminState != "Enable" && newCfg.PMClassBAdminState != "Disable" {
			qMgr.qsfpChannelConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassBAdminState Value")
		}
		cfgEnt.PMClassBAdminState = newCfg.PMClassBAdminState
	} else {
		cfgEnt.PMClassBAdminState = qsfpChannelCfgEnt.PMClassBAdminState
	}
	if mask&objects.QSFP_CHANNEL_UPDATE_PM_CLASS_C_ADMIN_STATE == objects.QSFP_CHANNEL_UPDATE_PM_CLASS_C_ADMIN_STATE {
		if newCfg.PMClassCAdminState != "Enable" && newCfg.PMClassCAdminState != "Disable" {
			qMgr.qsfpChannelConfigMutex.Unlock()
			return false, errors.New("Invalid PMClassCAdminState Value")
		}
		cfgEnt.PMClassCAdminState = newCfg.PMClassCAdminState
	} else {
		cfgEnt.PMClassCAdminState = qsfpChannelCfgEnt.PMClassCAdminState
	}

	if !(cfgEnt.HigherAlarmRXPower >= cfgEnt.HigherWarningRXPower &&
		cfgEnt.HigherWarningRXPower > cfgEnt.LowerWarningRXPower &&
		cfgEnt.LowerWarningRXPower >= cfgEnt.LowerAlarmRXPower) {
		qMgr.qsfpChannelConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, please verify the thresholds")
	}

	if !(cfgEnt.HigherAlarmTXPower >= cfgEnt.HigherWarningTXPower &&
		cfgEnt.HigherWarningTXPower > cfgEnt.LowerWarningTXPower &&
		cfgEnt.LowerWarningTXPower >= cfgEnt.LowerAlarmTXPower) {
		qMgr.qsfpChannelConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, please verify the thresholds")
	}

	if !(cfgEnt.HigherAlarmTXBias >= cfgEnt.HigherWarningTXBias &&
		cfgEnt.HigherWarningTXBias > cfgEnt.LowerWarningTXBias &&
		cfgEnt.LowerWarningTXBias >= cfgEnt.LowerAlarmTXBias) {
		qMgr.qsfpChannelConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, please verify the thresholds")
	}

	if qsfpChannelCfgEnt.AdminState != cfgEnt.AdminState {
		if cfgEnt.AdminState == "Disable" {
			qMgr.clearExistingQsfpChannelFaults(newCfg.QsfpId, newCfg.ChannelNum)
			qMgr.logger.Info("Clear all the existing Faults")
		}
	}
	if qsfpChannelCfgEnt.PMClassAAdminState != cfgEnt.PMClassAAdminState {
		if cfgEnt.PMClassAAdminState == "Disable" {
			// Flush PM RingBuffer
			qMgr.flushHistoricQsfpChannelPM(newCfg.QsfpId, newCfg.ChannelNum, "Class-A")
			qMgr.logger.Info("Flush Class A PM Ring buffer")
		}
	}
	if qsfpChannelCfgEnt.PMClassBAdminState != cfgEnt.PMClassBAdminState {
		if cfgEnt.PMClassBAdminState == "Disable" {
			// Flush PM RingBuffer
			qMgr.flushHistoricQsfpChannelPM(newCfg.QsfpId, newCfg.ChannelNum, "Class-B")
			qMgr.logger.Info("Flush Class B PM Ring buffer")
		}
	}
	if qsfpChannelCfgEnt.PMClassCAdminState != cfgEnt.PMClassCAdminState {
		if cfgEnt.PMClassCAdminState == "Disable" {
			// Flush PM RingBuffer
			qMgr.flushHistoricQsfpChannelPM(newCfg.QsfpId, newCfg.ChannelNum, "Class-C")
			qMgr.logger.Info("Flush Class C PM Ring buffer")
		}
	}
	qMgr.qsfpChannelConfigDB[qsfpChannel] = cfgEnt

	qMgr.qsfpChannelConfigMutex.Unlock()
	return true, nil
}

func (qMgr *QsfpManager) flushHistoricQsfpPM(QsfpId int32, Class string) {
	switch Class {
	case "Class-A":
		qMgr.qsfpClassAMutex.Lock()
		for idx := 0; idx < int(MaxNumQsfpRes); idx++ {
			resId := uint8(idx)
			qsfpResource := QsfpResource{
				QsfpId: QsfpId,
				ResId:  resId,
			}
			qMgr.qsfpClassAPM[qsfpResource].FlushRingBuffer()
		}
		qMgr.qsfpClassAMutex.Unlock()
	case "Class-B":
		qMgr.qsfpClassBMutex.Lock()
		for idx := 0; idx < int(MaxNumQsfpRes); idx++ {
			resId := uint8(idx)
			qsfpResource := QsfpResource{
				QsfpId: QsfpId,
				ResId:  resId,
			}
			qMgr.qsfpClassBPM[qsfpResource].FlushRingBuffer()
		}
		qMgr.qsfpClassBMutex.Unlock()
	case "Class-C":
		qMgr.qsfpClassCMutex.Lock()
		for idx := 0; idx < int(MaxNumQsfpRes); idx++ {
			resId := uint8(idx)
			qsfpResource := QsfpResource{
				QsfpId: QsfpId,
				ResId:  resId,
			}
			qMgr.qsfpClassCPM[qsfpResource].FlushRingBuffer()
		}
		qMgr.qsfpClassCMutex.Unlock()
	}
}

func (qMgr *QsfpManager) flushHistoricQsfpChannelPM(QsfpId int32, ChannelNum int32, Class string) {
	qsfpChannel := QsfpChannel{
		QsfpId:     QsfpId,
		ChannelNum: uint8(ChannelNum),
	}
	switch Class {
	case "Class-A":
		qMgr.qsfpChannelClassAMutex.Lock()
		for idx := 0; idx < int(MaxNumQsfpChannelRes); idx++ {
			resId := uint8(idx)
			qsfpChannelRes := QsfpChannelResource{
				QsfpChannel: qsfpChannel,
				ResId:       resId,
			}
			qMgr.qsfpChannelClassAPM[qsfpChannelRes].FlushRingBuffer()
		}
		qMgr.qsfpChannelClassAMutex.Unlock()
	case "Class-B":
		qMgr.qsfpChannelClassBMutex.Lock()
		for idx := 0; idx < int(MaxNumQsfpChannelRes); idx++ {
			resId := uint8(idx)
			qsfpChannelRes := QsfpChannelResource{
				QsfpChannel: qsfpChannel,
				ResId:       resId,
			}
			qMgr.qsfpChannelClassBPM[qsfpChannelRes].FlushRingBuffer()
		}
		qMgr.qsfpChannelClassBMutex.Unlock()
	case "Class-C":
		qMgr.qsfpChannelClassCMutex.Lock()
		for idx := 0; idx < int(MaxNumQsfpChannelRes); idx++ {
			resId := uint8(idx)
			qsfpChannelRes := QsfpChannelResource{
				QsfpChannel: qsfpChannel,
				ResId:       resId,
			}
			qMgr.qsfpChannelClassCPM[qsfpChannelRes].FlushRingBuffer()
		}
		qMgr.qsfpChannelClassCMutex.Unlock()
	}
}

func (qMgr *QsfpManager) GetQsfpPMState(QsfpId int32, Resource string, Class string) (*objects.QsfpPMState, error) {
	var qsfpPMObj objects.QsfpPMState
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	qMgr.qsfpConfigMutex.RLock()
	qsfpCfgEnt, exist := qMgr.qsfpConfigDB[QsfpId]
	if !exist {
		qMgr.qsfpConfigMutex.RUnlock()
		return nil, errors.New("Invalid QsfpId")
	}
	resId, err := getQsfpResourcId(Resource)
	if err != nil {
		qMgr.qsfpConfigMutex.RUnlock()
		return nil, errors.New("Invalid Resource Name")
	}
	qsfpResource := QsfpResource{
		QsfpId: QsfpId,
		ResId:  resId,
	}
	switch Class {
	case "Class-A":
		if qsfpCfgEnt.PMClassAAdminState == "Enable" {
			qMgr.qsfpClassAMutex.RLock()
			qsfpPMObj.Data = qMgr.qsfpClassAPM[qsfpResource].GetListOfEntriesFromRingBuffer()
			qMgr.qsfpClassAMutex.RUnlock()
		}
	case "Class-B":
		if qsfpCfgEnt.PMClassBAdminState == "Enable" {
			qMgr.qsfpClassBMutex.RLock()
			qsfpPMObj.Data = qMgr.qsfpClassBPM[qsfpResource].GetListOfEntriesFromRingBuffer()
			qMgr.qsfpClassBMutex.RUnlock()
		}
	case "Class-C":
		if qsfpCfgEnt.PMClassCAdminState == "Enable" {
			qMgr.qsfpClassCMutex.RLock()
			qsfpPMObj.Data = qMgr.qsfpClassCPM[qsfpResource].GetListOfEntriesFromRingBuffer()
			qMgr.qsfpClassCMutex.RUnlock()
		}
	default:
		qMgr.qsfpConfigMutex.RUnlock()
		return nil, errors.New("Invalid Class")
	}
	qMgr.qsfpConfigMutex.RUnlock()
	qsfpPMObj.QsfpId = QsfpId
	qsfpPMObj.Resource = Resource
	qsfpPMObj.Class = Class
	return &qsfpPMObj, nil
}

func (qMgr *QsfpManager) GetQsfpChannelPMState(QsfpId int32, ChannelNum int32, Resource string, Class string) (*objects.QsfpChannelPMState, error) {
	var qsfpChannelPMObj objects.QsfpChannelPMState
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if ChannelNum > int32(MaxNumOfQsfpChannel) {
		return nil, errors.New("Invalid ChannelNum")
	}
	qsfpChannel := QsfpChannel{
		QsfpId:     QsfpId,
		ChannelNum: uint8(ChannelNum),
	}
	qMgr.qsfpChannelConfigMutex.RLock()
	qsfpChannelCfgEnt, exist := qMgr.qsfpChannelConfigDB[qsfpChannel]
	if !exist {
		qMgr.qsfpChannelConfigMutex.RUnlock()
		return nil, errors.New("Invalid QsfpId")
	}
	resId, err := getQsfpChannelResourcId(Resource)
	if err != nil {
		qMgr.qsfpChannelConfigMutex.RUnlock()
		return nil, errors.New("Invalid Resource Name")
	}
	qsfpChannelResource := QsfpChannelResource{
		QsfpChannel: qsfpChannel,
		ResId:       resId,
	}
	switch Class {
	case "Class-A":
		if qsfpChannelCfgEnt.PMClassAAdminState == "Enable" {
			qMgr.qsfpChannelClassAMutex.RLock()
			qsfpChannelPMObj.Data = qMgr.qsfpChannelClassAPM[qsfpChannelResource].GetListOfEntriesFromRingBuffer()
			qMgr.qsfpChannelClassAMutex.RUnlock()
		}
	case "Class-B":
		if qsfpChannelCfgEnt.PMClassBAdminState == "Enable" {
			qMgr.qsfpChannelClassBMutex.RLock()
			qsfpChannelPMObj.Data = qMgr.qsfpChannelClassBPM[qsfpChannelResource].GetListOfEntriesFromRingBuffer()
			qMgr.qsfpChannelClassBMutex.RUnlock()
		}
	case "Class-C":
		if qsfpChannelCfgEnt.PMClassCAdminState == "Enable" {
			qMgr.qsfpChannelClassCMutex.RLock()
			qsfpChannelPMObj.Data = qMgr.qsfpChannelClassCPM[qsfpChannelResource].GetListOfEntriesFromRingBuffer()
			qMgr.qsfpChannelClassCMutex.RUnlock()
		}
	default:
		qMgr.qsfpChannelConfigMutex.RUnlock()
		return nil, errors.New("Invalid Class")
	}
	qMgr.qsfpChannelConfigMutex.RUnlock()
	qsfpChannelPMObj.QsfpId = QsfpId
	qsfpChannelPMObj.ChannelNum = ChannelNum
	qsfpChannelPMObj.Resource = Resource
	qsfpChannelPMObj.Class = Class
	return &qsfpChannelPMObj, nil
}

func getQsfpPMData(data *pluginCommon.QsfpPMData, resId uint8) (objects.QsfpPMData, error) {
	qsfpPMData := objects.QsfpPMData{
		TimeStamp: time.Now().String(),
	}
	switch resId {
	case TemperatureRes:
		qsfpPMData.Value = data.Temperature
	case VoltageRes:
		qsfpPMData.Value = data.Voltage
	default:
		return qsfpPMData, errors.New("Invalid resource Id")
	}
	return qsfpPMData, nil
}

func getQsfpChannelPMData(data *pluginCommon.QsfpPMData, ChannelNum uint8, resId uint8) (objects.QsfpPMData, error) {
	qsfpPMData := objects.QsfpPMData{
		TimeStamp: time.Now().String(),
	}
	switch resId {
	case RXPowerRes:
		qsfpPMData.Value = data.RXPower[ChannelNum-1]
	case TXPowerRes:
		qsfpPMData.Value = data.TXPower[ChannelNum-1]
	case TXBiasRes:
		qsfpPMData.Value = data.TXBias[ChannelNum-1]
	default:
		return qsfpPMData, errors.New("Invalid resource Id")
	}
	return qsfpPMData, nil
}

func (qMgr *QsfpManager) getCurrentQsfpEventStatus(pmData *pluginCommon.QsfpPMData, qsfpId int32, resId uint8) EventStatus {
	var curEventStatus EventStatus
	var data float64
	var highAlarmData float64
	var highWarnData float64
	var lowWarnData float64
	var lowAlarmData float64
	qsfpCfgEnt, _ := qMgr.qsfpConfigDB[qsfpId]
	switch resId {
	case TemperatureRes:
		data = pmData.Temperature
		highAlarmData = qsfpCfgEnt.HigherAlarmTemperature
		highWarnData = qsfpCfgEnt.HigherWarningTemperature
		lowWarnData = qsfpCfgEnt.LowerWarningTemperature
		lowAlarmData = qsfpCfgEnt.LowerAlarmTemperature
	case VoltageRes:
		data = pmData.Voltage
		highAlarmData = qsfpCfgEnt.HigherAlarmVoltage
		highWarnData = qsfpCfgEnt.HigherWarningVoltage
		lowWarnData = qsfpCfgEnt.LowerWarningVoltage
		lowAlarmData = qsfpCfgEnt.LowerAlarmVoltage
	}
	if data >= highAlarmData {
		curEventStatus.SentHigherAlarm = true
	}
	if data >= highWarnData {
		curEventStatus.SentHigherWarn = true
	}
	if data <= lowAlarmData {
		curEventStatus.SentLowerAlarm = true
	}
	if data <= lowWarnData {
		curEventStatus.SentLowerWarn = true
	}

	return curEventStatus
}

func (qMgr *QsfpManager) getCurrentQsfpChannelEventStatus(pmData *pluginCommon.QsfpPMData, qsfpId int32, channelNum uint8, resId uint8) EventStatus {
	var curEventStatus EventStatus
	var data float64
	var highAlarmData float64
	var highWarnData float64
	var lowWarnData float64
	var lowAlarmData float64
	qsfpChannel := QsfpChannel{
		QsfpId:     qsfpId,
		ChannelNum: channelNum,
	}
	qsfpChannelCfgEnt, _ := qMgr.qsfpChannelConfigDB[qsfpChannel]
	switch resId {
	case RXPowerRes:
		data = pmData.RXPower[channelNum-1]
		highAlarmData = qsfpChannelCfgEnt.HigherAlarmRXPower
		highWarnData = qsfpChannelCfgEnt.HigherWarningRXPower
		lowWarnData = qsfpChannelCfgEnt.LowerWarningRXPower
		lowAlarmData = qsfpChannelCfgEnt.LowerAlarmRXPower
	case TXPowerRes:
		data = pmData.TXPower[channelNum-1]
		highAlarmData = qsfpChannelCfgEnt.HigherAlarmTXPower
		highWarnData = qsfpChannelCfgEnt.HigherWarningTXPower
		lowWarnData = qsfpChannelCfgEnt.LowerWarningTXPower
		lowAlarmData = qsfpChannelCfgEnt.LowerAlarmTXPower
	case TXBiasRes:
		data = pmData.TXBias[channelNum-1]
		highAlarmData = qsfpChannelCfgEnt.HigherAlarmTXBias
		highWarnData = qsfpChannelCfgEnt.HigherWarningTXBias
		lowWarnData = qsfpChannelCfgEnt.LowerWarningTXBias
		lowAlarmData = qsfpChannelCfgEnt.LowerAlarmTXBias
	}
	if data >= highAlarmData {
		curEventStatus.SentHigherAlarm = true
	}
	if data >= highWarnData {
		curEventStatus.SentHigherWarn = true
	}
	if data <= lowAlarmData {
		curEventStatus.SentLowerAlarm = true
	}
	if data <= lowWarnData {
		curEventStatus.SentLowerWarn = true
	}
	return curEventStatus
}

func (qMgr *QsfpManager) publishQsfpEvents(qsfpId int32, resId uint8, data *pluginCommon.QsfpPMData, evts []events.EventId) {
	eventKey := events.QsfpKey{
		QsfpId: qsfpId,
	}
	txEvent := eventUtils.TxEvent{
		Key: eventKey,
	}
	if data == nil {
		txEvent.AdditionalInfo = "Clearing because of AdminState Disable"
	}

	for _, evt := range evts {
		txEvent.EventId = evt
		if data != nil {
			var eventData QsfpEventData
			switch evt {
			case events.QsfpTemperatureHigherTCAAlarm,
				events.QsfpTemperatureHigherTCAWarn,
				events.QsfpTemperatureLowerTCAAlarm,
				events.QsfpTemperatureLowerTCAWarn:
				eventData.Value = data.Temperature
			case events.QsfpVoltageHigherTCAAlarm,
				events.QsfpVoltageHigherTCAWarn,
				events.QsfpVoltageLowerTCAAlarm,
				events.QsfpVoltageLowerTCAWarn:
				eventData.Value = data.Voltage
			}
			txEvent.AdditionalData = eventData
		}
		txEvt := txEvent
		err := eventUtils.PublishEvents(&txEvt)
		if err != nil {
			qMgr.logger.Err("Error publishing event:", txEvt)
		}
	}
}

func (qMgr *QsfpManager) publishQsfpChannelEvents(qsfpId int32, channelNum uint8, resId uint8, data *pluginCommon.QsfpPMData, evts []events.EventId) {
	eventKey := events.QsfpChannelKey{
		QsfpId:     qsfpId,
		ChannelNum: int32(channelNum),
	}
	txEvent := eventUtils.TxEvent{
		Key: eventKey,
	}
	if data == nil {
		txEvent.AdditionalInfo = "Clearing because of AdminState Disable"
	}

	for _, evt := range evts {
		txEvent.EventId = evt
		if data != nil {
			var eventData QsfpChannelEventData
			switch evt {
			case events.QsfpRXPowerHigherTCAAlarm,
				events.QsfpRXPowerHigherTCAWarn,
				events.QsfpRXPowerLowerTCAAlarm,
				events.QsfpRXPowerLowerTCAWarn:
				eventData.Value = data.RXPower[channelNum-1]
			case events.QsfpTXPowerHigherTCAAlarm,
				events.QsfpTXPowerHigherTCAWarn,
				events.QsfpTXPowerLowerTCAAlarm,
				events.QsfpTXPowerLowerTCAWarn:
				eventData.Value = data.TXPower[channelNum-1]
			case events.QsfpTXBiasHigherTCAAlarm,
				events.QsfpTXBiasHigherTCAWarn,
				events.QsfpTXBiasLowerTCAAlarm,
				events.QsfpTXBiasLowerTCAWarn:
				eventData.Value = data.TXBias[channelNum-1]
			}
			txEvent.AdditionalData = eventData
		}
		txEvt := txEvent
		err := eventUtils.PublishEvents(&txEvt)
		if err != nil {
			qMgr.logger.Err("Error publishing event:", txEvt)
		}
	}
}

func getListofQsfpEvent(resId uint8, prevEvt EventStatus, curEvt EventStatus) []events.EventId {
	var evts []events.EventId
	switch resId {
	case TemperatureRes:
		if prevEvt.SentHigherAlarm != curEvt.SentHigherAlarm {
			if curEvt.SentHigherAlarm == true {
				evts = append(evts, events.QsfpTemperatureHigherTCAAlarm)
			} else {
				evts = append(evts, events.QsfpTemperatureHigherTCAAlarmClear)
			}
		}
		if prevEvt.SentHigherWarn != curEvt.SentHigherWarn {
			if curEvt.SentHigherWarn == true {
				evts = append(evts, events.QsfpTemperatureHigherTCAWarn)
			} else {
				evts = append(evts, events.QsfpTemperatureHigherTCAWarnClear)
			}
		}
		if prevEvt.SentLowerAlarm != curEvt.SentLowerAlarm {
			if curEvt.SentLowerAlarm == true {
				evts = append(evts, events.QsfpTemperatureLowerTCAAlarm)
			} else {
				evts = append(evts, events.QsfpTemperatureLowerTCAAlarmClear)
			}
		}
		if prevEvt.SentLowerWarn != curEvt.SentLowerWarn {
			if curEvt.SentLowerWarn == true {
				evts = append(evts, events.QsfpTemperatureLowerTCAWarn)
			} else {
				evts = append(evts, events.QsfpTemperatureLowerTCAWarnClear)
			}
		}
	case VoltageRes:
		if prevEvt.SentHigherAlarm != curEvt.SentHigherAlarm {
			if curEvt.SentHigherAlarm == true {
				evts = append(evts, events.QsfpVoltageHigherTCAAlarm)
			} else {
				evts = append(evts, events.QsfpVoltageHigherTCAAlarmClear)
			}
		}
		if prevEvt.SentHigherWarn != curEvt.SentHigherWarn {
			if curEvt.SentHigherWarn == true {
				evts = append(evts, events.QsfpVoltageHigherTCAWarn)
			} else {
				evts = append(evts, events.QsfpVoltageHigherTCAWarnClear)
			}
		}
		if prevEvt.SentLowerAlarm != curEvt.SentLowerAlarm {
			if curEvt.SentLowerAlarm == true {
				evts = append(evts, events.QsfpVoltageLowerTCAAlarm)
			} else {
				evts = append(evts, events.QsfpVoltageLowerTCAAlarmClear)
			}
		}
		if prevEvt.SentLowerWarn != curEvt.SentLowerWarn {
			if curEvt.SentLowerWarn == true {
				evts = append(evts, events.QsfpVoltageLowerTCAWarn)
			} else {
				evts = append(evts, events.QsfpVoltageLowerTCAWarnClear)
			}
		}
	}
	return evts
}

func getListofQsfpChannelEvent(resId uint8, prevEvt EventStatus, curEvt EventStatus) []events.EventId {
	var evts []events.EventId
	switch resId {
	case RXPowerRes:
		if prevEvt.SentHigherAlarm != curEvt.SentHigherAlarm {
			if curEvt.SentHigherAlarm == true {
				evts = append(evts, events.QsfpRXPowerHigherTCAAlarm)
			} else {
				evts = append(evts, events.QsfpRXPowerHigherTCAAlarmClear)
			}
		}
		if prevEvt.SentHigherWarn != curEvt.SentHigherWarn {
			if curEvt.SentHigherWarn == true {
				evts = append(evts, events.QsfpRXPowerHigherTCAWarn)
			} else {
				evts = append(evts, events.QsfpRXPowerHigherTCAWarnClear)
			}
		}
		if prevEvt.SentLowerAlarm != curEvt.SentLowerAlarm {
			if curEvt.SentLowerAlarm == true {
				evts = append(evts, events.QsfpRXPowerLowerTCAAlarm)
			} else {
				evts = append(evts, events.QsfpRXPowerLowerTCAAlarmClear)
			}
		}
		if prevEvt.SentLowerWarn != curEvt.SentLowerWarn {
			if curEvt.SentLowerWarn == true {
				evts = append(evts, events.QsfpRXPowerLowerTCAWarn)
			} else {
				evts = append(evts, events.QsfpRXPowerLowerTCAWarnClear)
			}
		}
	case TXPowerRes:
		if prevEvt.SentHigherAlarm != curEvt.SentHigherAlarm {
			if curEvt.SentHigherAlarm == true {
				evts = append(evts, events.QsfpTXPowerHigherTCAAlarm)
			} else {
				evts = append(evts, events.QsfpTXPowerHigherTCAAlarmClear)
			}
		}
		if prevEvt.SentHigherWarn != curEvt.SentHigherWarn {
			if curEvt.SentHigherWarn == true {
				evts = append(evts, events.QsfpTXPowerHigherTCAWarn)
			} else {
				evts = append(evts, events.QsfpTXPowerHigherTCAWarnClear)
			}
		}
		if prevEvt.SentLowerAlarm != curEvt.SentLowerAlarm {
			if curEvt.SentLowerAlarm == true {
				evts = append(evts, events.QsfpTXPowerLowerTCAAlarm)
			} else {
				evts = append(evts, events.QsfpTXPowerLowerTCAAlarmClear)
			}
		}
		if prevEvt.SentLowerWarn != curEvt.SentLowerWarn {
			if curEvt.SentLowerWarn == true {
				evts = append(evts, events.QsfpTXPowerLowerTCAWarn)
			} else {
				evts = append(evts, events.QsfpTXPowerLowerTCAWarnClear)
			}
		}
	case TXBiasRes:
		if prevEvt.SentHigherAlarm != curEvt.SentHigherAlarm {
			if curEvt.SentHigherAlarm == true {
				evts = append(evts, events.QsfpTXBiasHigherTCAAlarm)
			} else {
				evts = append(evts, events.QsfpTXBiasHigherTCAAlarmClear)
			}
		}
		if prevEvt.SentHigherWarn != curEvt.SentHigherWarn {
			if curEvt.SentHigherWarn == true {
				evts = append(evts, events.QsfpTXBiasHigherTCAWarn)
			} else {
				evts = append(evts, events.QsfpTXBiasHigherTCAWarnClear)
			}
		}
		if prevEvt.SentLowerAlarm != curEvt.SentLowerAlarm {
			if curEvt.SentLowerAlarm == true {
				evts = append(evts, events.QsfpTXBiasLowerTCAAlarm)
			} else {
				evts = append(evts, events.QsfpTXBiasLowerTCAAlarmClear)
			}
		}
		if prevEvt.SentLowerWarn != curEvt.SentLowerWarn {
			if curEvt.SentLowerWarn == true {
				evts = append(evts, events.QsfpTXBiasLowerTCAWarn)
			} else {
				evts = append(evts, events.QsfpTXBiasLowerTCAWarnClear)
			}
		}
	}
	return evts
}

func (qMgr *QsfpManager) processQsfpChannelEvents(data *pluginCommon.QsfpPMData, qsfpId int32, channelNum uint8) {
	qsfpChannel := QsfpChannel{
		QsfpId:     qsfpId,
		ChannelNum: channelNum,
	}

	for id := 0; id < int(MaxNumQsfpChannelRes); id++ {
		var evts []events.EventId
		resId := uint8(id)
		qsfpChannelResource := QsfpChannelResource{
			QsfpChannel: qsfpChannel,
			ResId:       resId,
		}
		qMgr.qsfpChannelEventMsgStatusMutex.RLock()
		prevEventStatus, _ := qMgr.qsfpChannelEventMsgStatus[qsfpChannelResource]
		qMgr.qsfpChannelEventMsgStatusMutex.RUnlock()

		curEventStatus := qMgr.getCurrentQsfpChannelEventStatus(data, qsfpId, channelNum, resId)
		evts = append(evts, getListofQsfpChannelEvent(resId, prevEventStatus, curEventStatus)...)
		if prevEventStatus != curEventStatus {
			qMgr.qsfpChannelEventMsgStatusMutex.Lock()
			qMgr.qsfpChannelEventMsgStatus[qsfpChannelResource] = curEventStatus
			qMgr.qsfpChannelEventMsgStatusMutex.Unlock()
		}
		qMgr.publishQsfpChannelEvents(qsfpId, channelNum, resId, data, evts)
	}
}

func (qMgr *QsfpManager) processQsfpEvents(data *pluginCommon.QsfpPMData, qsfpId int32) {
	for id := 0; id < int(MaxNumQsfpRes); id++ {
		var evts []events.EventId
		resId := uint8(id)
		qsfpResource := QsfpResource{
			QsfpId: qsfpId,
			ResId:  resId,
		}
		qMgr.qsfpEventMsgStatusMutex.RLock()
		prevEventStatus, _ := qMgr.qsfpEventMsgStatus[qsfpResource]
		qMgr.qsfpEventMsgStatusMutex.RUnlock()

		curEventStatus := qMgr.getCurrentQsfpEventStatus(data, qsfpId, resId)
		evts = append(evts, getListofQsfpEvent(resId, prevEventStatus, curEventStatus)...)

		if prevEventStatus != curEventStatus {
			qMgr.qsfpEventMsgStatusMutex.Lock()
			qMgr.qsfpEventMsgStatus[qsfpResource] = curEventStatus
			qMgr.qsfpEventMsgStatusMutex.Unlock()
		}
		qMgr.publishQsfpEvents(qsfpId, resId, data, evts)
	}
}

func (qMgr *QsfpManager) clearExistingQsfpFaults(qsfpId int32) {
	for idx := 0; idx < int(MaxNumQsfpRes); idx++ {
		var evts []events.EventId
		resId := uint8(idx)
		qsfpResource := QsfpResource{
			QsfpId: qsfpId,
			ResId:  resId,
		}
		qMgr.qsfpEventMsgStatusMutex.RLock()
		eventStatus, _ := qMgr.qsfpEventMsgStatus[qsfpResource]
		qMgr.qsfpEventMsgStatusMutex.RUnlock()
		switch resId {
		case TemperatureRes:
			if eventStatus.SentHigherAlarm == true {
				evts = append(evts, events.QsfpTemperatureHigherTCAAlarmClear)
				eventStatus.SentHigherAlarm = false
			}
			if eventStatus.SentHigherWarn == true {
				evts = append(evts, events.QsfpTemperatureHigherTCAWarnClear)
				eventStatus.SentHigherWarn = false
			}
			if eventStatus.SentLowerAlarm == true {
				evts = append(evts, events.QsfpTemperatureLowerTCAAlarmClear)
				eventStatus.SentLowerAlarm = false
			}
			if eventStatus.SentLowerWarn == true {
				evts = append(evts, events.QsfpTemperatureLowerTCAWarnClear)
				eventStatus.SentLowerWarn = false
			}
		case VoltageRes:
			if eventStatus.SentHigherAlarm == true {
				evts = append(evts, events.QsfpVoltageHigherTCAAlarmClear)
				eventStatus.SentHigherAlarm = false
			}
			if eventStatus.SentHigherWarn == true {
				evts = append(evts, events.QsfpVoltageHigherTCAWarnClear)
				eventStatus.SentHigherWarn = false
			}
			if eventStatus.SentLowerAlarm == true {
				evts = append(evts, events.QsfpVoltageLowerTCAAlarmClear)
				eventStatus.SentLowerAlarm = false
			}
			if eventStatus.SentLowerWarn == true {
				evts = append(evts, events.QsfpVoltageLowerTCAWarnClear)
				eventStatus.SentLowerWarn = false
			}
		}
		qMgr.qsfpEventMsgStatusMutex.Lock()
		qMgr.qsfpEventMsgStatus[qsfpResource] = eventStatus
		qMgr.qsfpEventMsgStatusMutex.Unlock()
		qMgr.publishQsfpEvents(qsfpId, resId, nil, evts)
	}
}

func (qMgr *QsfpManager) clearExistingQsfpChannelFaults(qsfpId int32, channelNum int32) {
	qsfpChannel := QsfpChannel{
		QsfpId:     qsfpId,
		ChannelNum: uint8(channelNum),
	}
	for idx := 0; idx < int(MaxNumQsfpChannelRes); idx++ {
		var evts []events.EventId
		resId := uint8(idx)
		qsfpChannelResource := QsfpChannelResource{
			QsfpChannel: qsfpChannel,
			ResId:       resId,
		}
		qMgr.qsfpChannelEventMsgStatusMutex.RLock()
		eventStatus, _ := qMgr.qsfpChannelEventMsgStatus[qsfpChannelResource]
		qMgr.qsfpChannelEventMsgStatusMutex.RUnlock()
		switch resId {
		case RXPowerRes:
			if eventStatus.SentHigherAlarm == true {
				evts = append(evts, events.QsfpRXPowerHigherTCAAlarmClear)
				eventStatus.SentHigherAlarm = false
			}
			if eventStatus.SentHigherWarn == true {
				evts = append(evts, events.QsfpRXPowerHigherTCAWarnClear)
				eventStatus.SentHigherWarn = false
			}
			if eventStatus.SentLowerAlarm == true {
				evts = append(evts, events.QsfpRXPowerLowerTCAAlarmClear)
				eventStatus.SentLowerAlarm = false
			}
			if eventStatus.SentLowerWarn == true {
				evts = append(evts, events.QsfpRXPowerLowerTCAWarnClear)
				eventStatus.SentLowerWarn = false
			}
		case TXPowerRes:
			if eventStatus.SentHigherAlarm == true {
				evts = append(evts, events.QsfpTXPowerHigherTCAAlarmClear)
				eventStatus.SentHigherAlarm = false
			}
			if eventStatus.SentHigherWarn == true {
				evts = append(evts, events.QsfpTXPowerHigherTCAWarnClear)
				eventStatus.SentHigherWarn = false
			}
			if eventStatus.SentLowerAlarm == true {
				evts = append(evts, events.QsfpTXPowerLowerTCAAlarmClear)
				eventStatus.SentLowerAlarm = false
			}
			if eventStatus.SentLowerWarn == true {
				evts = append(evts, events.QsfpTXPowerLowerTCAWarnClear)
				eventStatus.SentLowerWarn = false
			}
		case TXBiasRes:
			if eventStatus.SentHigherAlarm == true {
				evts = append(evts, events.QsfpTXBiasHigherTCAAlarmClear)
				eventStatus.SentHigherAlarm = false
			}
			if eventStatus.SentHigherWarn == true {
				evts = append(evts, events.QsfpTXBiasHigherTCAWarnClear)
				eventStatus.SentHigherWarn = false
			}
			if eventStatus.SentLowerAlarm == true {
				evts = append(evts, events.QsfpTXBiasLowerTCAAlarmClear)
				eventStatus.SentLowerAlarm = false
			}
			if eventStatus.SentLowerWarn == true {
				evts = append(evts, events.QsfpTXBiasLowerTCAWarnClear)
				eventStatus.SentLowerWarn = false
			}
		}
		qMgr.qsfpChannelEventMsgStatusMutex.Lock()
		qMgr.qsfpChannelEventMsgStatus[qsfpChannelResource] = eventStatus
		qMgr.qsfpChannelEventMsgStatusMutex.Unlock()
		qMgr.publishQsfpChannelEvents(qsfpId, uint8(channelNum), resId, nil, evts)
	}
}

func (qMgr *QsfpManager) processQsfpChannelPMData(data *pluginCommon.QsfpPMData, QsfpId int32, ChannelNum uint8, Class string) {
	qsfpChannel := QsfpChannel{
		QsfpId:     QsfpId,
		ChannelNum: ChannelNum,
	}
	qMgr.qsfpChannelConfigMutex.RLock()
	qsfpChannelCfgEnt, _ := qMgr.qsfpChannelConfigDB[qsfpChannel]
	switch Class {
	case "Class-A":
		if qsfpChannelCfgEnt.PMClassAAdminState == "Enable" {
			for id := 0; id < int(MaxNumQsfpChannelRes); id++ {
				resId := uint8(id)
				pmData, err := getQsfpChannelPMData(data, ChannelNum, resId)
				if err != nil {
					qMgr.logger.Err("Wrong resource Id:", resId)
					continue
				}
				qsfpChannelResource := QsfpChannelResource{
					QsfpChannel: qsfpChannel,
					ResId:       resId,
				}
				qMgr.qsfpChannelClassAMutex.Lock()
				qMgr.qsfpChannelClassAPM[qsfpChannelResource].InsertIntoRingBuffer(pmData)
				qMgr.qsfpChannelClassAMutex.Unlock()
			}
		}
	case "Class-B":
		if qsfpChannelCfgEnt.PMClassBAdminState == "Enable" {
			for id := 0; id < int(MaxNumQsfpChannelRes); id++ {
				resId := uint8(id)
				pmData, err := getQsfpChannelPMData(data, ChannelNum, resId)
				if err != nil {
					qMgr.logger.Err("Wrong resource Id:", resId)
					continue
				}
				qsfpChannelResource := QsfpChannelResource{
					QsfpChannel: qsfpChannel,
					ResId:       resId,
				}
				qMgr.qsfpChannelClassBMutex.Lock()
				qMgr.qsfpChannelClassBPM[qsfpChannelResource].InsertIntoRingBuffer(pmData)
				qMgr.qsfpChannelClassBMutex.Unlock()
			}
		}
	case "Class-C":
		if qsfpChannelCfgEnt.PMClassCAdminState == "Enable" {
			for id := 0; id < int(MaxNumQsfpChannelRes); id++ {
				resId := uint8(id)
				pmData, err := getQsfpChannelPMData(data, ChannelNum, resId)
				if err != nil {
					qMgr.logger.Err("Wrong resource Id:", resId)
					continue
				}
				qsfpChannelResource := QsfpChannelResource{
					QsfpChannel: qsfpChannel,
					ResId:       resId,
				}
				qMgr.qsfpChannelClassCMutex.Lock()
				qMgr.qsfpChannelClassCPM[qsfpChannelResource].InsertIntoRingBuffer(pmData)
				qMgr.qsfpChannelClassCMutex.Unlock()
			}
		}
	}
	qMgr.qsfpChannelConfigMutex.RUnlock()
}

func (qMgr *QsfpManager) processQsfpPMData(data *pluginCommon.QsfpPMData, QsfpId int32, Class string) {
	qMgr.qsfpConfigMutex.RLock()
	qsfpCfgEnt, _ := qMgr.qsfpConfigDB[QsfpId]
	switch Class {
	case "Class-A":
		if qsfpCfgEnt.PMClassAAdminState == "Enable" {
			for id := 0; id < int(MaxNumQsfpRes); id++ {
				resId := uint8(id)
				pmData, err := getQsfpPMData(data, resId)
				if err != nil {
					qMgr.logger.Err("Wrong resource Id:", resId)
					continue
				}
				qsfpResource := QsfpResource{
					QsfpId: QsfpId,
					ResId:  resId,
				}
				qMgr.qsfpClassAMutex.Lock()
				qMgr.qsfpClassAPM[qsfpResource].InsertIntoRingBuffer(pmData)
				qMgr.qsfpClassAMutex.Unlock()
			}
		}
	case "Class-B":
		if qsfpCfgEnt.PMClassBAdminState == "Enable" {
			for id := 0; id < int(MaxNumQsfpRes); id++ {
				resId := uint8(id)
				pmData, err := getQsfpPMData(data, resId)
				if err != nil {
					qMgr.logger.Err("Wrong resource Id:", resId)
					continue
				}
				qsfpResource := QsfpResource{
					QsfpId: QsfpId,
					ResId:  resId,
				}
				qMgr.qsfpClassBMutex.Lock()
				qMgr.qsfpClassBPM[qsfpResource].InsertIntoRingBuffer(pmData)
				qMgr.qsfpClassBMutex.Unlock()
			}
		}
	case "Class-C":
		if qsfpCfgEnt.PMClassCAdminState == "Enable" {
			for id := 0; id < int(MaxNumQsfpRes); id++ {
				resId := uint8(id)
				pmData, err := getQsfpPMData(data, resId)
				if err != nil {
					qMgr.logger.Err("Wrong resource Id:", resId)
					continue
				}
				qsfpResource := QsfpResource{
					QsfpId: QsfpId,
					ResId:  resId,
				}
				qMgr.qsfpClassCMutex.Lock()
				qMgr.qsfpClassCPM[qsfpResource].InsertIntoRingBuffer(pmData)
				qMgr.qsfpClassCMutex.Unlock()
			}
		}
	}
	qMgr.qsfpConfigMutex.RUnlock()
}

func (qMgr *QsfpManager) ProcessQsfpPMData(data *pluginCommon.QsfpPMData, qsfpId int32, class string) {
	qMgr.qsfpConfigMutex.RLock()
	qsfpCfgEnt, _ := qMgr.qsfpConfigDB[qsfpId]
	if qsfpCfgEnt.AdminState == "Enable" {
		qMgr.processQsfpEvents(data, qsfpId)
	}
	qMgr.qsfpConfigMutex.RUnlock()
	for ch := 1; ch <= int(MaxNumOfQsfpChannel); ch++ {
		qsfpChannel := QsfpChannel{
			QsfpId:     qsfpId,
			ChannelNum: uint8(ch),
		}
		qMgr.qsfpChannelConfigMutex.RLock()
		qsfpChannelCfgEnt, _ := qMgr.qsfpChannelConfigDB[qsfpChannel]
		if qsfpChannelCfgEnt.AdminState == "Enable" {
			qMgr.processQsfpChannelEvents(data, qsfpId, uint8(ch))
		}
		qMgr.qsfpChannelConfigMutex.RUnlock()
	}
	qMgr.processQsfpPMData(data, qsfpId, class)
	for ch := 1; ch <= int(MaxNumOfQsfpChannel); ch++ {
		qMgr.processQsfpChannelPMData(data, qsfpId, uint8(ch), class)
	}
}

func (qMgr *QsfpManager) StartQsfpPMClass(class string) {
	for id := 1; id <= int(qMgr.numOfQsfps); id++ {
		qsfpId := int32(id)
		qMgr.stateMutex.Lock()
		qsfpPMData, err := qMgr.plugin.GetQsfpPMData(qsfpId)
		qMgr.stateMutex.Unlock()
		if err == nil {
			qMgr.ProcessQsfpPMData(&qsfpPMData, qsfpId, class)
		}
	}
	switch class {
	case "Class-A":
		classAPMFunc := func() {
			for id := 1; id <= int(qMgr.numOfQsfps); id++ {
				qsfpId := int32(id)
				qMgr.stateMutex.Lock()
				qsfpPMData, err := qMgr.plugin.GetQsfpPMData(qsfpId)
				qMgr.stateMutex.Unlock()
				if err == nil {
					qMgr.ProcessQsfpPMData(&qsfpPMData, qsfpId, class)
				}
			}
			qMgr.classAPMTimer.Reset(qsfpClassAInterval)
		}
		qMgr.classAPMTimer = time.AfterFunc(qsfpClassAInterval, classAPMFunc)
	case "Class-B":
		classBPMFunc := func() {
			for id := 1; id <= int(qMgr.numOfQsfps); id++ {
				qsfpId := int32(id)
				qMgr.stateMutex.Lock()
				qsfpPMData, err := qMgr.plugin.GetQsfpPMData(qsfpId)
				qMgr.stateMutex.Unlock()
				if err == nil {
					qMgr.ProcessQsfpPMData(&qsfpPMData, qsfpId, class)
				}
			}
			qMgr.classBPMTimer.Reset(qsfpClassBInterval)
		}
		qMgr.classBPMTimer = time.AfterFunc(qsfpClassBInterval, classBPMFunc)
	case "Class-C":
		classCPMFunc := func() {
			for id := 1; id <= int(qMgr.numOfQsfps); id++ {
				qsfpId := int32(id)
				qMgr.stateMutex.Lock()
				qsfpPMData, err := qMgr.plugin.GetQsfpPMData(qsfpId)
				qMgr.stateMutex.Unlock()
				if err == nil {
					qMgr.ProcessQsfpPMData(&qsfpPMData, qsfpId, class)
				}
			}
			qMgr.classCPMTimer.Reset(qsfpClassCInterval)
		}
		qMgr.classCPMTimer = time.AfterFunc(qsfpClassCInterval, classCPMFunc)
	}
}

func (qMgr *QsfpManager) InitQsfpPM() {
	for idx := 1; idx <= int(qMgr.numOfQsfps); idx++ {
		for qsfpRes := 0; qsfpRes < int(MaxNumQsfpRes); qsfpRes++ {
			qsfpResource := QsfpResource{
				QsfpId: int32(idx),
				ResId:  uint8(qsfpRes),
			}
			qMgr.qsfpClassAPM[qsfpResource] = new(ringBuffer.RingBuffer)
			qMgr.qsfpClassAPM[qsfpResource].SetRingBufferCapacity(qsfpClassABufSize)
			qMgr.qsfpClassBPM[qsfpResource] = new(ringBuffer.RingBuffer)
			qMgr.qsfpClassBPM[qsfpResource].SetRingBufferCapacity(qsfpClassBBufSize)
			qMgr.qsfpClassCPM[qsfpResource] = new(ringBuffer.RingBuffer)
			qMgr.qsfpClassCPM[qsfpResource].SetRingBufferCapacity(qsfpClassCBufSize)
		}
		for ch := 1; ch <= int(MaxNumOfQsfpChannel); ch++ {
			qsfpChannel := QsfpChannel{
				QsfpId:     int32(idx),
				ChannelNum: uint8(ch),
			}
			for qsfpChannelRes := 0; qsfpChannelRes < int(MaxNumQsfpChannelRes); qsfpChannelRes++ {
				qsfpChannelResource := QsfpChannelResource{
					QsfpChannel: qsfpChannel,
					ResId:       uint8(qsfpChannelRes),
				}
				qMgr.qsfpChannelClassAPM[qsfpChannelResource] = new(ringBuffer.RingBuffer)
				qMgr.qsfpChannelClassAPM[qsfpChannelResource].SetRingBufferCapacity(qsfpClassABufSize)
				qMgr.qsfpChannelClassBPM[qsfpChannelResource] = new(ringBuffer.RingBuffer)
				qMgr.qsfpChannelClassBPM[qsfpChannelResource].SetRingBufferCapacity(qsfpClassBBufSize)
				qMgr.qsfpChannelClassCPM[qsfpChannelResource] = new(ringBuffer.RingBuffer)
				qMgr.qsfpChannelClassCPM[qsfpChannelResource].SetRingBufferCapacity(qsfpClassCBufSize)
			}
		}
	}
}
func (qMgr *QsfpManager) StartQsfpPM() {
	qMgr.InitQsfpPM()
	qMgr.StartQsfpPMClass("Class-A")
	qMgr.StartQsfpPMClass("Class-B")
	qMgr.StartQsfpPMClass("Class-C")
}
