//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package main

import (
	"encoding/json"
	"io/ioutil"
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type initCfgParams struct {
	thriftServerPort int
}

func parseConfigFile(paramsDir string) *initCfgParams {
	var clientsList []ClientJson
	var initCfg initCfgParams

	//Set default thrift port to use, override based on config file below
	initCfg.thriftServerPort = 4000
	bytes, err := ioutil.ReadFile(paramsDir + "clients.json")
	if err != nil {
		logger.Err("Error retrieving thrift server port number using default port 4000")
	} else {
		err := json.Unmarshal(bytes, &clientsList)
		if err != nil {
			logger.Err("Error retrieving thrift server port number using default port 4000")
		} else {
			for _, client := range clientsList {
				if client.Name == "asicd" {
					initCfg.thriftServerPort = client.Port
				}
			}
		}
	}
	return &initCfg
}
