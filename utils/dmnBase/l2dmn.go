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

package dmnBase

import (
	"arpd"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
	"utils/ipcutils"
)

type ArpdClient struct {
	DmnClientBase
	ClientHdl *arpd.ARPDServicesClient
}

type L2Daemon struct {
	FSDaemon
	Arpdclnt ArpdClient
}

func (dmn *L2Daemon) ConnectToArpd() error {
	configFile := dmn.ParamsDir + "clients.json"
	var clientsList []ClientJson

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		dmn.Logger.Info(fmt.Sprintln("Error in reading configuration file ", configFile))
		return errors.New(fmt.Sprintln("Error in reading config file: ", configFile))
	}
	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		dmn.Logger.Info("Error in Unmarshalling Json")
		return errors.New("Error unmarshaling")
	}

	for _, client := range clientsList {
		if client.Name == "arpd" {
			dmn.Logger.Info(fmt.Sprintln("found  arpd at port ", client.Port))
			dmn.Arpdclnt.Address = "localhost:" + strconv.Itoa(client.Port)
			dmn.Arpdclnt.Transport, dmn.Arpdclnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(dmn.Arpdclnt.Address)
			if dmn.Arpdclnt.Transport != nil && dmn.Arpdclnt.PtrProtocolFactory != nil {
				dmn.Logger.Info(fmt.Sprintln("connecting to arpd,asicdclnt.IsConnected:", dmn.Arpdclnt.IsConnected))
				dmn.Arpdclnt.ClientHdl =
					arpd.NewARPDServicesClientFactory(dmn.Arpdclnt.Transport, dmn.Arpdclnt.PtrProtocolFactory)
				dmn.Arpdclnt.IsConnected = true
			} else {
				dmn.Logger.Info(fmt.Sprintf("Failed to connect to Asicd, retrying until connection is successful"))
				count := 0
				ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
				for _ = range ticker.C {
					dmn.Arpdclnt.Transport, dmn.Arpdclnt.PtrProtocolFactory, err =
						ipcutils.CreateIPCHandles(dmn.Arpdclnt.Address)
					if err == nil {
						ticker.Stop()
						break
					}
					count++
					if (count % 10) == 0 {
						dmn.Logger.Info("Still can't connect to Arpd, retrying...")
					}
				}
			}
		}
	}
	return err
}
func (dmn *L2Daemon) InitSubscribers([]string) (err error) {
	dmn.Logger.Info("L2 Dmn InitSubscribers")
	dmn.FSDaemon.InitSubscribers(nil)
	return err
}

func (dmn *L2Daemon) Init(dmnName string, logPrefix string) bool {
	if !dmn.FSDaemon.Init(dmnName, logPrefix) {
		dmn.Logger.Err("Init failed")
		return false
	}
	return true
}
func (dmn *L2Daemon) ConnectToServers() error {
	err := dmn.FSDaemon.ConnectToServers()
	if err != nil {
		dmn.Logger.Err("Failed to connect to server")
		return err
	}
	err = dmn.ConnectToArpd()
	if err != nil {
		dmn.Logger.Err("Failed to connect to ARPD")
		return err
	}
	return nil
}
func (dmn *L2Daemon) NewServer() {
	dmn.FSDaemon.NewServer()
}
