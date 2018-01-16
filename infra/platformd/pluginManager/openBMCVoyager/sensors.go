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
	//	"bytes"
	"encoding/json"
	"fmt"
	"infra/platformd/pluginManager/pluginCommon"
	"io/ioutil"
	"net/http"
	//"net/url"
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

type FanSensor struct {
	Fan4Rear  string `json:"Fan 4 rear"`
	Fan3Front string `json:"Fan 3 front"`
	Fan2Front string `json:"Fan 2 front"`
	Fan5Front string `json:"Fan 5 front"`
	Name      string `json:"name"`
	Fan1Front string `json:"Fan 1 front"`
	Adapter   string `json:"Adapter"`
	Fan1Rear  string `json:"Fan 1 rear"`
	Fan5Rear  string `json:"Fan 5 rear"`
	Fan2Rear  string `json:"Fan 2 rear"`
	Fan3Rear  string `json:"Fan 3 rear"`
	Fan4Front string `json:"Fan 4 front"`
}

type IR3581Sensor struct {
	Adapter string `json:"Adapter"`
	Iout    string `json:"Iout"`
	Name    string `json:"name"`
	Vout    string `json:"Vout"`
}

type IR3584Sensor struct {
	Adapter string `json:"Adapter"`
	Iout    string `json:"Iout"`
	Name    string `json:"name"`
	Vout    string `json:"Vout"`
}

type InletTempSensor struct {
	Temp    string `json:"Inlet Temp"`
	Adapter string `json:"Adapter"`
	Name    string `json:"name"`
}

type MicroServerAmbientTempSensor struct {
	Temp    string `json:"Microserver Ambient Temp"`
	Adapter string `json:"Adapter"`
	Name    string `json:"name"`
}

type VoltageSensor struct {
	VMON5   string `json:"+3.3 VMON5 Voltage"`
	Name    string `json:"name"`
	VMON7   string `json:"+1.8 VMON7 Voltage"`
	Adapter string `json:"Adapter"`
	VMON9   string `json:"+1.2 VMON9 Voltage"`
	VMON6   string `json:"+3.3 VMON6 Voltage"`
	VMON2   string `json:"+1.0 VMON2 Voltage"`
	VMON1   string `json:"+1.25 VMON2 Voltage"`
	VMON8   string `json:"+3.3 VMON8 Voltage"`
	VMON4   string `json:"+1.8 VMON4 Voltage"`
	VMON3   string `json:"+1.0 VMON3 Voltage"`
	VMON10  string `json:"+12 VMON10 Voltage"`
}

type OutletTempSensor struct {
	Adapter string `json:"Adapter"`
	Temp    string `json:"Outlet Temp"`
	Name    string `json:"name"`
}

type SensorName struct {
	Name string `json:"name"`
}

func (driver *openBMCVoyagerDriver) GetAllSensorState(state *pluginCommon.SensorState) error {
	url := "http://" + driver.ipAddr + ":" + driver.port + "/api/sys/sensors"
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
	return extractSensorData(state, body)
}

func extractSensorData(state *pluginCommon.SensorState, body []byte) error {
	var info Sensor

	err := json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	for idx := 0; idx < 7; idx++ {
		var sensorName SensorName
		msg, _ := json.Marshal(info.Info[idx])
		err := json.Unmarshal(msg, &sensorName)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		switch sensorName.Name {
		case "IR3581-i2c-1-70":
			extractIR3581Data(state, info.Info[idx])
		case "IR3584-i2c-11-72":
			extractIR3584Data(state, info.Info[idx])
		case "pwr1014a-i2c-2-40":
			extractVoltageData(state, info.Info[idx])
		case "tmp75-i2c-3-48":
			extractInletTempData(state, info.Info[idx])
		case "tmp75-i2c-3-4a":
			extractMicroserverTempData(state, info.Info[idx])
		case "fancpld-i2c-8-33":
			extractFanData(state, info.Info[idx])
		case "tmp75-i2c-3-4b":
			extractOutletTempData(state, info.Info[idx])
		default:
			fmt.Println("Unrecognized sensor")
		}
	}

	return err
}

func getPowerValue(Iout string, Vout string) float64 {
	iout, _ := strconv.ParseFloat((strings.Split(Iout, " "))[0], 64)
	vout, _ := strconv.ParseFloat((strings.Split(Vout, " "))[0], 64)
	return iout * vout
}

func extractIR3581Data(state *pluginCommon.SensorState, info interface{}) error {
	var sensor IR3581Sensor
	msg, _ := json.Marshal(info)
	err := json.Unmarshal(msg, &sensor)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	ent, _ := state.PowerConverterSensor[sensor.Name]
	ent.Value = getPowerValue(sensor.Iout, sensor.Vout)
	state.PowerConverterSensor[sensor.Name] = ent
	return nil
}

func extractIR3584Data(state *pluginCommon.SensorState, info interface{}) error {
	var sensor IR3584Sensor
	msg, _ := json.Marshal(info)
	err := json.Unmarshal(msg, &sensor)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	ent, _ := state.PowerConverterSensor[sensor.Name]
	ent.Value = getPowerValue(sensor.Iout, sensor.Vout)
	state.PowerConverterSensor[sensor.Name] = ent
	return nil
}

func getTempValue(temp string) float64 {
	val, _ := strconv.ParseFloat((strings.Split(temp, " "))[0], 64)
	return val
}

func extractInletTempData(state *pluginCommon.SensorState, info interface{}) error {
	var sensor InletTempSensor
	msg, _ := json.Marshal(info)
	err := json.Unmarshal(msg, &sensor)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	ent, _ := state.TemperatureSensor["Inlet"]
	ent.Value = getTempValue(sensor.Temp)
	state.TemperatureSensor["Inlet"] = ent
	return nil

}

func extractMicroserverTempData(state *pluginCommon.SensorState, info interface{}) error {
	var sensor MicroServerAmbientTempSensor
	msg, _ := json.Marshal(info)
	err := json.Unmarshal(msg, &sensor)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	ent, _ := state.TemperatureSensor["Microserver"]
	ent.Value = getTempValue(sensor.Temp)
	state.TemperatureSensor["Microserver"] = ent
	return nil
}

func extractOutletTempData(state *pluginCommon.SensorState, info interface{}) error {
	var sensor OutletTempSensor
	msg, _ := json.Marshal(info)
	err := json.Unmarshal(msg, &sensor)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	ent, _ := state.TemperatureSensor["Outlet"]
	ent.Value = getTempValue(sensor.Temp)
	state.TemperatureSensor["Outlet"] = ent
	return nil

}

func getVoltageValue(volt string) float64 {
	val, _ := strconv.ParseFloat((strings.Split(volt, " "))[0], 64)
	return val
}

func extractVoltageData(state *pluginCommon.SensorState, info interface{}) error {
	var sensor VoltageSensor
	msg, _ := json.Marshal(info)
	err := json.Unmarshal(msg, &sensor)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	ent, _ := state.VoltageSensor["+3.3 VMON5"]
	ent.Value = getVoltageValue(sensor.VMON5)
	state.VoltageSensor["+3.3 VMON5"] = ent

	ent, _ = state.VoltageSensor["+1.8 VMON7"]
	ent.Value = getVoltageValue(sensor.VMON7)
	state.VoltageSensor["+1.8 VMON7"] = ent

	ent, _ = state.VoltageSensor["+1.2 VMON9"]
	ent.Value = getVoltageValue(sensor.VMON9)
	state.VoltageSensor["+1.2 VMON9"] = ent

	ent, _ = state.VoltageSensor["+3.3 VMON6"]
	ent.Value = getVoltageValue(sensor.VMON6)
	state.VoltageSensor["+3.3 VMON6"] = ent

	ent, _ = state.VoltageSensor["+1.0 VMON2"]
	ent.Value = getVoltageValue(sensor.VMON2)
	state.VoltageSensor["+1.0 VMON2"] = ent

	ent, _ = state.VoltageSensor["+1.25 VMON1"]
	ent.Value = getVoltageValue(sensor.VMON1)
	state.VoltageSensor["+1.25 VMON1"] = ent

	ent, _ = state.VoltageSensor["+3.3 VMON8"]
	ent.Value = getVoltageValue(sensor.VMON8)
	state.VoltageSensor["+3.3 VMON8"] = ent

	ent, _ = state.VoltageSensor["+1.8 VMON4"]
	ent.Value = getVoltageValue(sensor.VMON4)
	state.VoltageSensor["+1.8 VMON4"] = ent

	ent, _ = state.VoltageSensor["+1.0 VMON3"]
	ent.Value = getVoltageValue(sensor.VMON3)
	state.VoltageSensor["+1.0 VMON3"] = ent

	ent, _ = state.VoltageSensor["+12 VMON10"]
	ent.Value = getVoltageValue(sensor.VMON10)
	state.VoltageSensor["+12 VMON10"] = ent
	return nil
}

func getFanValue(rpm string) int32 {
	val, _ := strconv.Atoi((strings.Split(rpm, " "))[0])
	return int32(val)
}

func extractFanData(state *pluginCommon.SensorState, info interface{}) error {
	var sensor FanSensor
	msg, _ := json.Marshal(info)
	err := json.Unmarshal(msg, &sensor)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	ent, _ := state.FanSensor["Fan 1 front"]
	ent.Value = getFanValue(sensor.Fan1Front)
	state.FanSensor["Fan 1 front"] = ent

	ent, _ = state.FanSensor["Fan 1 rear"]
	ent.Value = getFanValue(sensor.Fan1Rear)
	state.FanSensor["Fan 1 rear"] = ent

	ent, _ = state.FanSensor["Fan 2 front"]
	ent.Value = getFanValue(sensor.Fan2Front)
	state.FanSensor["Fan 2 front"] = ent

	ent, _ = state.FanSensor["Fan 2 rear"]
	ent.Value = getFanValue(sensor.Fan2Rear)
	state.FanSensor["Fan 2 rear"] = ent

	ent, _ = state.FanSensor["Fan 3 front"]
	ent.Value = getFanValue(sensor.Fan3Front)
	state.FanSensor["Fan 3 front"] = ent

	ent, _ = state.FanSensor["Fan 3 rear"]
	ent.Value = getFanValue(sensor.Fan3Rear)
	state.FanSensor["Fan 3 rear"] = ent

	ent, _ = state.FanSensor["Fan 4 front"]
	ent.Value = getFanValue(sensor.Fan4Front)
	state.FanSensor["Fan 4 front"] = ent

	ent, _ = state.FanSensor["Fan 4 rear"]
	ent.Value = getFanValue(sensor.Fan4Rear)
	state.FanSensor["Fan 4 rear"] = ent

	ent, _ = state.FanSensor["Fan 5 front"]
	ent.Value = getFanValue(sensor.Fan5Front)
	state.FanSensor["Fan 5 front"] = ent

	ent, _ = state.FanSensor["Fan 5 rear"]
	ent.Value = getFanValue(sensor.Fan5Rear)
	state.FanSensor["Fan 5 rear"] = ent

	return nil
}
