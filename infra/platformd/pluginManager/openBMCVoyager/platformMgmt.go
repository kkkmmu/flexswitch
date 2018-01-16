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

package openBMCVoyager

import (
	"encoding/json"
	"fmt"
	"infra/platformd/pluginManager/pluginCommon"
	"io/ioutil"
	"net/http"
)

type PlatformMgmtResource []string
type PlatformMgmtActions []string

type PlatformMgmt struct {
	Info PlatformMgmtInfo `json:"Information"`
	Reso SensorResource   `json:"Resources"`
	Act  SensorActions    `json:"Actions"`
}

type PlatformMgmtInfo struct {
	Uptime      string `json:"Uptime"`
	Description string `json:"Description"`
	ResetReason string `json:"Reset Reason"`
	MemoryUsage string `json:"Memory Usage"`
	Version     string `json:"OpenBMC Version"`
	CPUUsage    string `json:"CPU Usage"`
}

func (driver *openBMCVoyagerDriver) GetPlatformMgmtDeviceState(state *pluginCommon.PlatformMgmtDeviceState) error {
	url := "http://" + driver.ipAddr + ":" + driver.port + "/api/sys/bmc"
	//fmt.Println("URL:>", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return extractPlatformData(state, body)
}

func extractPlatformData(state *pluginCommon.PlatformMgmtDeviceState, body []byte) error {
	var plat PlatformMgmt

	err := json.Unmarshal(body, &plat)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	state.Uptime = plat.Info.Uptime
	state.DeviceName = "BMC"
	state.Description = plat.Info.Description
	state.ResetReason = plat.Info.ResetReason
	state.MemoryUsage = plat.Info.MemoryUsage
	state.Version = plat.Info.Version
	state.CPUUsage = plat.Info.CPUUsage
	return nil
}
