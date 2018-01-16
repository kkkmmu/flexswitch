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

package objects

type TemperatureSensorState struct {
	Name               string
	CurrentTemperature float64
}

type TemperatureSensorStateGetInfo struct {
	EndIdx int
	Count  int
	More   bool
	List   []*TemperatureSensorState
}

type TemperatureSensorConfig struct {
	Name                   string
	AdminState             string
	HigherAlarmThreshold   float64
	HigherWarningThreshold float64
	LowerWarningThreshold  float64
	LowerAlarmThreshold    float64
	PMClassAAdminState     string
	PMClassBAdminState     string
	PMClassCAdminState     string
}

type TemperatureSensorConfigGetInfo struct {
	EndIdx int
	Count  int
	More   bool
	List   []*TemperatureSensorConfig
}

const (
	TEMP_SENSOR_UPDATE_ADMIN_STATE            = 0x1
	TEMP_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD = 0x2
	TEMP_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD  = 0x4
	TEMP_SENSOR_UPDATE_LOWER_WARN_THRESHOLD   = 0x8
	TEMP_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD  = 0x10
	TEMP_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE = 0x20
	TEMP_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE = 0x40
	TEMP_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE = 0x80
)

type TemperatureSensorPMData struct {
	TimeStamp string
	Value     float64
}

type TemperatureSensorPMState struct {
	Name  string
	Class string
	Data  []interface{}
}
