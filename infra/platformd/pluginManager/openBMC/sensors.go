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

package openBMC

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type SensorResource []string
type SensorActions []string

type Sensor struct {
	Info []interface{}  `json:"Information"`
	Reso SensorResource `json:"Resources"`
	Act  SensorActions  `json:"Actions"`
}

type CPUSensors struct {
	Voltage5     string `json:"+5V Voltage"`
	CoreVoltage  string `json:"CPU Vcore"`
	Name         string `json:"name"`
	Adapter      string `json:"Adapter"`
	Voltage12    string `json:"+12V Voltage"`
	Voltage3     string `json:"+3V Voltage"`
	MemTemp      string `json:"Memory Temp"`
	CPUTemp      string `json:"CPU Temp"`
	VDIMMVoltage string `json:"VDIMM Voltage"`
}

type FanSensors struct {
	Fan4Rear  string `json:"Fan 4 rear"`
	Fan3Front string `json:"Fan 3 front"`
	Fan2Front string `json:"Fan 2 front"`
	Fan5Front string `json:"Fan 5 front"`
	Name      string `json:"name"`
	Fan1Front string `json:"Fan 1 Front"`
	Adapter   string `json:"Adapter"`
	Fan1Rear  string `json:"Fan 1 rear"`
	Fan5Rear  string `json:"Fan 5 rear"`
	Fan2Rear  string `json:"Fan 2 rear"`
	Fan3Rear  string `json:"Fan 3 rear"`
	Fan4Front string `json:"Fan 4 front"`
}

type CurrentVoltageSensors struct {
	Current   string `json:"Current"`
	Adapter   string `json:"Adapter"`
	Voltage12 string `json:"+12 Voltage"`
	Name      string `json:"name"`
}

type OutletMidTempSensors struct {
	Adapter string `json:"Adapter"`
	Name    string `json:"name"`
	Temp    string `json:"Outlet Middle Temp"`
}

type InletMidTempSensors struct {
	Adapter string `json:"Adapter"`
	Temp    string `json:"Inlet Middle Temp"`
	Name    string `json:"name"`
}

type InletLeftTempSensors struct {
	Temp    string `json:"Inlet Left Temp"`
	Adapter string `json:"Adapter"`
	Name    string `json:"name"`
}

type SwitchTempSensors struct {
	Adapter string `json:"Adapter"`
	Name    string `json:"name"`
	Temp    string `json:"Switch Temp"`
}

type InletRightTempSensors struct {
	Adapter string `json:"Adapter"`
	Temp    string `json:"Inlet Right Temp"`
	Name    string `json:"name"`
}

type OutletRightTempSensors struct {
	Adapter string `json:"Adapter"`
	Temp    string `json:"Outlet Right Temp"`
	Name    string `json:"name"`
}

type OutletLeftTempSensors struct {
	Adapter string `json:"Adapter"`
	Temp    string `json:"Outlet Left Temp"`
	Name    string `json:"name"`
}

type TempSensor struct {
	OutMidTemp   OutletMidTempSensors
	OutLeftTemp  OutletLeftTempSensors
	OutRightTemp OutletRightTempSensors
	InMidTemp    InletMidTempSensors
	InLeftTemp   InletLeftTempSensors
	InRightTemp  InletRightTempSensors
	SwitchTemp   SwitchTempSensors
}

type SensorData struct {
	CpuSensor      CPUSensors
	FanSensor      FanSensors
	CurrVoltSensor CurrentVoltageSensors
	TempSensor     TempSensor
}

func (driver *openBMCDriver) GetSensorState() (data SensorData, err error) {
	var jsonStr = []byte(nil)
	url := "http://" + driver.ipAddr + ":" + driver.port + "/api/sys/sensors"
	//fmt.Println("URL:>", url)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return data, err
	}
	req.Header.Set("Accept", "application/json")
	body, err := SendHttpCmd(req)
	if err != nil {
		return data, err
	}
	return extractSensorData(body)
}

func SendHttpCmd(req *http.Request) (body []byte, err error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	//fmt.Println("response Body:", string(body))
	return body, err
}

func extractSensorData(body []byte) (data SensorData, err error) {
	var info Sensor

	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println("Error:", err)
		return data, err
	}
	//data.CpuSensor, _ = extractCPUData(info.Info[0])
	data.FanSensor, _ = extractFanData(info.Info[0])
	//data.CurrVoltSensor, _ = extractCurrentVoltageData(info.Info[2])
	data.TempSensor.OutMidTemp, _ = extractOutletMidTempData(info.Info[2])
	data.TempSensor.InMidTemp, _ = extractInletMidTempData(info.Info[3])
	data.TempSensor.InLeftTemp, _ = extractInletLeftTempData(info.Info[4])
	data.TempSensor.SwitchTemp, _ = extractSwitchTempData(info.Info[5])
	data.TempSensor.InRightTemp, _ = extractInletRightTempData(info.Info[6])
	data.TempSensor.OutRightTemp, _ = extractOutletRightTempData(info.Info[7])
	data.TempSensor.OutLeftTemp, _ = extractOutletLeftTempData(info.Info[8])
	return data, err
}

func extractInletRightTempData(info interface{}) (tempSensor InletRightTempSensors, err error) {
	msg, _ := json.Marshal(info)
	err = json.Unmarshal(msg, &tempSensor)
	if err != nil {
		fmt.Println("Error:", err)
		return tempSensor, err
	}

	//fmt.Println("Inlet Right Temp Sensors: ", tempSensor, "Name:", tempSensor.Name)
	return tempSensor, err
}

func extractOutletRightTempData(info interface{}) (tempSensor OutletRightTempSensors, err error) {
	msg, _ := json.Marshal(info)
	err = json.Unmarshal(msg, &tempSensor)
	if err != nil {
		fmt.Println("Error:", err)
		return tempSensor, err
	}

	//fmt.Println("Outlet Right Temp Sensors: ", tempSensor, "Name:", tempSensor.Name)
	return tempSensor, err
}

func extractOutletLeftTempData(info interface{}) (tempSensor OutletLeftTempSensors, err error) {
	msg, _ := json.Marshal(info)
	err = json.Unmarshal(msg, &tempSensor)
	if err != nil {
		fmt.Println("Error:", err)
		return tempSensor, err
	}

	//fmt.Println("Outlet Left Temp Sensors: ", tempSensor, "Name:", tempSensor.Name)
	return tempSensor, err
}

func extractSwitchTempData(info interface{}) (tempSensor SwitchTempSensors, err error) {
	msg, _ := json.Marshal(info)
	err = json.Unmarshal(msg, &tempSensor)
	if err != nil {
		fmt.Println("Error:", err)
		return tempSensor, err
	}

	//fmt.Println("Switch Temp Sensors: ", tempSensor, "Name:", tempSensor.Name)
	return tempSensor, err
}

func extractInletMidTempData(info interface{}) (tempSensor InletMidTempSensors, err error) {
	msg, _ := json.Marshal(info)
	err = json.Unmarshal(msg, &tempSensor)
	if err != nil {
		fmt.Println("Error:", err)
		return tempSensor, err
	}

	//fmt.Println("Inlet Mid Temp Sensors: ", tempSensor, "Name:", tempSensor.Name)
	return tempSensor, err
}

func extractInletLeftTempData(info interface{}) (tempSensor InletLeftTempSensors, err error) {
	msg, _ := json.Marshal(info)
	err = json.Unmarshal(msg, &tempSensor)
	if err != nil {
		fmt.Println("Error:", err)
		return tempSensor, err
	}

	//fmt.Println("Inlet Left Temp Sensors: ", tempSensor, "Name:", tempSensor.Name)
	return tempSensor, err
}

func extractOutletMidTempData(info interface{}) (tempSensor OutletMidTempSensors, err error) {
	msg, _ := json.Marshal(info)
	err = json.Unmarshal(msg, &tempSensor)
	if err != nil {
		fmt.Println("Error:", err)
		return tempSensor, err
	}

	//fmt.Println("Outlet Mid Temp Sensors: ", tempSensor, "Name:", tempSensor.Name)
	return tempSensor, err
}

func extractCPUData(info interface{}) (cpuSensor CPUSensors, err error) {
	msg, _ := json.Marshal(info)
	err = json.Unmarshal(msg, &cpuSensor)
	if err != nil {
		fmt.Println("Error:", err)
		return cpuSensor, err
	}

	//fmt.Println("CPU Sensors: ", cpuSensor, "Name:", cpuSensor.Name)
	return cpuSensor, err
}

func extractFanData(info interface{}) (fanSensor FanSensors, err error) {
	msg, _ := json.Marshal(info)
	err = json.Unmarshal(msg, &fanSensor)
	if err != nil {
		fmt.Println("Error:", err)
		return fanSensor, err
	}

	//fmt.Println("Fan Sensors: ", fanSensor, "Name:", fanSensor.Name)
	return fanSensor, err
}

func extractCurrentVoltageData(info interface{}) (currVoltSensor CurrentVoltageSensors, err error) {
	msg, _ := json.Marshal(info)
	err = json.Unmarshal(msg, &currVoltSensor)
	if err != nil {
		fmt.Println("Error:", err)
		return currVoltSensor, err
	}

	//fmt.Println("Current Voltage Sensors: ", currVoltSensor, "Name:", currVoltSensor.Name)
	return currVoltSensor, err
}

func convertFanSpeedStringToInt32(speedStr string) (speed int32) {
	if speedStr != "" {
		spd, err := strconv.Atoi((strings.Split(speedStr, " "))[0])
		if err != nil {
			speed = int32(0)
		} else {
			speed = int32(spd)
		}
	}
	return speed
}
