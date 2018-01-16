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

import (
	"utils/dbutils"
	"utils/logging"
)

//Plugin name constants
const (
	ONLP_PLUGIN           = "onlp"
	OpenBMC_PLUGIN        = "openbmc"
	OpenBMCVoyager_PLUGIN = "openbmcvoyager"
	Dummy_PLUGIN          = "dummy"
)

type PluginInitParams struct {
	Logger     logging.LoggerIntf
	PluginName string
	IpAddr     string
	Port       string
	EventDbHdl dbutils.DBIntf
}

type FanState struct {
	FanId         int32
	OperMode      string
	OperSpeed     int32
	OperDirection string
	Status        string
	Model         string
	SerialNum     string
	LedId         int32
	Valid         bool
}

type SfpState struct {
	SfpId      int32
	SfpSpeed   string
	SfpLos     string
	SfpPresent string
	SfpType    string
	SerialNum  string
	EEPROM     string
}

type PlatformState struct {
	ObjName      string
	ProductName  string
	SerialNum    string
	Manufacturer string
	Vendor       string
	Release      string
	PlatformName string
	Version      string
}

type ThermalState struct {
	ThermalId                 int32
	Location                  string
	Temperature               string
	LowerWatermarkTemperature string
	UpperWatermarkTemperature string
	ShutdownTemperature       string
	Valid                     bool
}

type PsuState struct {
	PsuId     int32
	Model     string
	SerialNum string
	Status    string
	VoltIn    int32
	VoltOut   int32
	AmpIn     int32
	AmpOut    int32
	PwrIn     int32
	PwrOut    int32
}

type LedState struct {
	LedId       int32
	LedIdentify string
	LedState    string
	LedColor    string
}
