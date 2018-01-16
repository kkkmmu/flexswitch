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

package openBMC

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	//"strconv"
	//"strings"
)

type BMCResource []string
type BMCActions []string

type BMC struct {
	Info      BMCInfo     `json:"Information"`
	Resources BMCResource `json:"Resources"`
	Actions   BMCActions  `json:"Actions"`
}

type BMCInfo struct {
	Uptime         string `json:"Uptime"`
	Description    string `json:"Description"`
	ResetReson     string `json:"Reset Reason"`
	MemoryUsage    string `json:"Memory Usage"`
	OpenBMCVersion string `json:"OpenBMC Version"`
	CPUUsage       string `json:"CPU Usage"`
}

func (driver *openBMCDriver) GetBMCInfo() (info BMCInfo, err error) {
	var jsonStr = []byte(nil)
	url := "http://" + driver.ipAddr + ":" + driver.port + "/api/sys/bmc"
	//fmt.Println("URL:>", url)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return info, err
	}
	req.Header.Set("Accept", "application/json")
	body, err := SendHttpCmd(req)
	if err != nil {
		return info, err
	}
	return extractBMCData(body)
}

func extractBMCData(body []byte) (data BMCInfo, err error) {
	var bmc BMC

	err = json.Unmarshal(body, &bmc)
	if err != nil {
		fmt.Println("Error:", err)
		return data, err
	}

	return bmc.Info, err
}
