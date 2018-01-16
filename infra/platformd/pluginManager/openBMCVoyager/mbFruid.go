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

/*
import (
	"bytes"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	//"strconv"
	//"strings"
)

type MBFruidResource []string
type MBFruidActions []string

type MBFruid struct {
	Info      MBFruidInfo     `json:"Information"`
	Resources MBFruidResource `json:"Resources"`
	Actions   MBFruidActions  `json:"Actions"`
}

type MBFruidInfo struct {
	LocOnFabric string `json:"Location on Fabric"`
	//ProSubVer         string `json:"Product Sub-Version"`
	ProSubVer         int32  `json:"Product Sub-Version"`
	FbPCBPartNum      string `json:"Facebook PCB Part Number"`
	CRC8              string `json:"CRC8"`
	SysAssemPartNum   string `json:"System Assembly Part Number"`
	ProSerialNum      string `json:"Product Serial Number"`
	SysManuDate       string `json:"System Manufacturing Date"`
	LocalMAC          string `json:"Local MAC"`
	AssemAt           string `json:"Assembled At"`
	ODMPCBASerialNum  string `json:"ODM PCBA Serial Number"`
	ProductAssetTag   string `json:"Product Asset Tag"`
	ProductName       string `json:"Product Name"`
	ODMPCBAPartNumber string `json:"ODM PCBA Part Number"`
	//ProductProductionState string `json:"Product Production State"`
	ProductProductionState int32  `json:"Product Production State"`
	ProductPartNumber      string `json:"Product Part Number"`
	PCBManufacturer        string `json:"PCB Manufacturer"`
	SystemManufacturer     string `json:"System Manufacturer"`
	//ExtendedMACAddressSize string `json:"Extended MAC Address Size"`
	ExtendedMACAddressSize int32  `json:"Extended MAC Address Size"`
	FacebookPCBAPartNumber string `json:"Facebook PCBA Part Number"`
	//Version                string `json:"Version"`
	Version         int32  `json:"Version"`
	ExtendedMACBase string `json:"Extended MAC Base"`
	//ProductVer             string `json:"Product Version"`
	ProductVer int32 `json:"Product Version"`
}

func (driver *openBMCDriver) GetMBFruidInfo() (info MBFruidInfo, err error) {
	var jsonStr = []byte(nil)
	url := "http://" + driver.ipAddr + ":" + driver.port + "/api/sys/mb/fruid"
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
	return extractMBFruidData(body)
}

func extractMBFruidData(body []byte) (data MBFruidInfo, err error) {
	var mbFruid MBFruid

	err = json.Unmarshal(body, &mbFruid)
	if err != nil {
		fmt.Println("Error:", err)
		return data, err
	}

	return mbFruid.Info, err
}
*/
