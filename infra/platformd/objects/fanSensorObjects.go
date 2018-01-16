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

type FanSensorState struct {
	Name         string
	CurrentSpeed int32
}

type FanSensorStateGetInfo struct {
	EndIdx int
	Count  int
	More   bool
	List   []*FanSensorState
}

type FanSensorConfig struct {
	Name                   string
	AdminState             string
	HigherAlarmThreshold   int32
	HigherWarningThreshold int32
	LowerWarningThreshold  int32
	LowerAlarmThreshold    int32
	PMClassAAdminState     string
	PMClassBAdminState     string
	PMClassCAdminState     string
}

type FanSensorConfigGetInfo struct {
	EndIdx int
	Count  int
	More   bool
	List   []*FanSensorConfig
}

const (
	FAN_SENSOR_UPDATE_ADMIN_STATE            = 0x1
	FAN_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD = 0x2
	FAN_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD  = 0x4
	FAN_SENSOR_UPDATE_LOWER_WARN_THRESHOLD   = 0x8
	FAN_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD  = 0x10
	FAN_SENSOR_UPDATE_PM_CLASS_A_ADMIN_STATE = 0x20
	FAN_SENSOR_UPDATE_PM_CLASS_B_ADMIN_STATE = 0x40
	FAN_SENSOR_UPDATE_PM_CLASS_C_ADMIN_STATE = 0x80
)

type FanSensorPMData struct {
	TimeStamp string
	Value     int32
}

type FanSensorPMState struct {
	Name  string
	Class string
	Data  []interface{}
}
