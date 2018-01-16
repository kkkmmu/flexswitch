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

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"utils/eventUtils"
)

type SysProfile struct {
	Port int `json:"Notifier_Port"`
}

func (server *NMGRServer) getNotifierPort() (string, error) {
	var sysProfile SysProfile

	sysProfileFile := server.paramsDir + "systemProfile.json"
	bytes, err := ioutil.ReadFile(sysProfileFile)
	if err != nil {
		return "", errors.New(fmt.Sprintln("Error reading the sysProfile file", sysProfileFile))
	}
	err = json.Unmarshal(bytes, &sysProfile)
	if err != nil {
		return "", errors.New(fmt.Sprintln("Error unmarshalling sysProfile file", sysProfileFile))
	}
	port := strconv.Itoa(sysProfile.Port)
	return port, nil
}

func (server *NMGRServer) getDmnList() ([]string, error) {
	evtJson, err := eventUtils.ParseEventsJson()
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Error parsing events.json file", err))
	}
	var dmnList []string
	for _, daemon := range evtJson.DaemonEvents {
		dmnList = append(dmnList, daemon.DaemonName)
	}
	return dmnList, nil
}
