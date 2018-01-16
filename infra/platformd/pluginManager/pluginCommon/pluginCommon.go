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

package pluginCommon

type FanStateStruct struct {
	OperMode      string
	OperSpeed     int32
	OperDirection string
	Status        string
	Model         [100]byte
	SerialNum     [100]byte
}

const (
	FAN_MODE_OFF int = 0x0
	FAN_MODE_ON  int = 0x1
)

const (
	FAN_MODE_OFF_STR string = "OFF"
	FAN_MODE_ON_STR         = "ON"
)

const (
	FAN_DIR_B2F     int = 0x0
	FAN_DIR_F2B     int = 0x1
	FAN_DIR_INVALID     = 0x2
)

const (
	FAN_DIR_B2F_STR     string = "Back2Front"
	FAN_DIR_F2B_STR            = "Front2Back"
	FAN_DIR_INVALID_STR        = "InvalidDir"
)

const (
	FAN_STATUS_PRESENT int = 0x0
	FAN_STATUS_MISSING     = 0x2
	FAN_STATUS_FAILED      = 0x3
	FAN_STATUS_NORMAL      = 0x4
)

const (
	FAN_STATUS_PRESENT_STR string = "PRESENT"
	FAN_STATUS_MISSING_STR        = "MISSING"
	FAN_STATUS_FAILED_STR         = "FAILED"
	FAN_STATUS_NORMAL_STR         = "NORMAL"
)

type FanSensorData struct {
	Value int32
}

type TemperatureSensorData struct {
	Value float64
}

type VoltageSensorData struct {
	Value float64
}

type PowerConverterSensorData struct {
	Value float64
}

type SensorState struct {
	FanSensor            map[string]FanSensorData
	TemperatureSensor    map[string]TemperatureSensorData
	VoltageSensor        map[string]VoltageSensorData
	PowerConverterSensor map[string]PowerConverterSensorData
}

const (
	QsfpNumChannel int32 = 4
)

type QsfpState struct {
	VendorName         string
	VendorOUI          string
	VendorPartNumber   string
	VendorRevision     string
	VendorSerialNumber string
	DataCode           string
	Temperature        float64
	Voltage            float64
	CurrBER            float64
	AccBER             float64
	MinBER             float64
	MaxBER             float64
	UDF0               float64
	UDF1               float64
	UDF2               float64
	UDF3               float64
	RXPower            [QsfpNumChannel]float64
	TXPower            [QsfpNumChannel]float64
	TXBias             [QsfpNumChannel]float64
}

type PlatformMgmtDeviceState struct {
	DeviceName  string
	Uptime      string
	Description string
	ResetReason string
	MemoryUsage string
	Version     string
	CPUUsage    string
}

type QsfpPMData struct {
	Temperature float64
	Voltage     float64
	RXPower     [QsfpNumChannel]float64
	TXPower     [QsfpNumChannel]float64
	TXBias      [QsfpNumChannel]float64
}
