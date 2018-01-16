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
	"encoding/json"
	"errors"
	"fmt"
	nanomsg "github.com/op/go-nanomsg"
	"io/ioutil"
	"ribd"
	"strconv"
	"time"
	"utils/ipcutils"
)

type RIBdClient struct {
	DmnClientBase
	ClientHdl *ribd.RIBDServicesClient
}

type L3Daemon struct {
	FSDaemon
	Ribdclnt           RIBdClient
	RibdSubSocket      *nanomsg.SubSocket
	RibdSubSocketCh    chan []byte
	RibdSubSocketErrCh chan error
}

func (dmn *L3Daemon) ConnectToRIBd() error {
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
		if client.Name == "ribd" {
			dmn.Logger.Info(fmt.Sprintln("found  ribd at port ", client.Port))
			dmn.Ribdclnt.Address = "localhost:" + strconv.Itoa(client.Port)
			dmn.Ribdclnt.Transport, dmn.Ribdclnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(dmn.Ribdclnt.Address)
			if dmn.Ribdclnt.Transport != nil && dmn.Ribdclnt.PtrProtocolFactory != nil {
				dmn.Ribdclnt.ClientHdl = ribd.NewRIBDServicesClientFactory(dmn.Ribdclnt.Transport, dmn.Ribdclnt.PtrProtocolFactory)
				dmn.Ribdclnt.IsConnected = true
			} else {
				dmn.Logger.Info(fmt.Sprintf("Failed to connect to Asicd, retrying until connection is successful"))
				count := 0
				ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
				for _ = range ticker.C {
					dmn.Ribdclnt.Transport, dmn.Ribdclnt.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(dmn.Ribdclnt.Address)
					if err == nil {
						ticker.Stop()
						break
					}
					count++
					if (count % 10) == 0 {
						dmn.Logger.Info("Still can't connect to Ribd, retrying...")
					}
				}
			}
		}
	}
	return err
}

func (dmn *L3Daemon) CreateRIBdSubscriber(sub string) error {
	dmn.Logger.Info("Listen for RIBd updates")
	err := dmn.ListenForRIBdUpdates(sub)
	if err != nil {
		dmn.Logger.Err("Error initialzing RIBd subscriber")
		return err
	}
	for {
		dmn.Logger.Info("Read on RIBd subscriber socket...")
		rxBuf, err := dmn.RibdSubSocket.Recv(0)
		if err != nil {
			dmn.Logger.Err(fmt.Sprintln("Recv on RIBd subscriber socket failed with error:", err))
			dmn.RibdSubSocketErrCh <- err
			continue
		}
		dmn.Logger.Info(fmt.Sprintln("RIB subscriber recv returned:", rxBuf))
		dmn.RibdSubSocketCh <- rxBuf
	}
	return nil
}

func (dmn *L3Daemon) ListenForRIBdUpdates(address string) error {
	var err error
	if dmn.RibdSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		dmn.Logger.Err(fmt.Sprintln("Failed to create RIBd subscribe socket, error:", err))
		return err
	}

	if _, err = dmn.RibdSubSocket.Connect(address); err != nil {
		dmn.Logger.Err(fmt.Sprintln("Failed to connect to RIBd publisher socket, address:", address, "error:", err))
		return err
	}

	if err = dmn.RibdSubSocket.Subscribe(""); err != nil {
		dmn.Logger.Err(fmt.Sprintln("Failed to subscribe to \"\" on RIBd subscribe socket, error:", err))
		return err
	}

	dmn.Logger.Info(fmt.Sprintln("Connected to RIBd publisher at address:", address))
	if err = dmn.RibdSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		dmn.Logger.Err(fmt.Sprintln("Failed to set the buffer size for RIBd publisher socket, error:", err))
		return err
	}
	return nil
}

func (dmn *L3Daemon) InitSubscribers(ribdsubscriberList []string) (err error) {
	dmn.Logger.Info("L3 Dmn InitSubscribers")
	dmn.FSDaemon.InitSubscribers(nil)
	for _, sub := range ribdsubscriberList {
		go dmn.CreateRIBdSubscriber(sub)
	}
	return err
}
func (dmn *L3Daemon) Init(dmnName string, logPrefix string) bool {
	if !dmn.FSDaemon.Init(dmnName, logPrefix) {
		dmn.Logger.Err("Init failed")
		return false
	}
	return true
}

func (dmn *L3Daemon) ConnectToServers() error {
	err := dmn.FSDaemon.ConnectToServers()
	if err != nil {
		dmn.Logger.Err("Failed to connect to servers")
		return err
	}
	err = dmn.ConnectToRIBd()
	if err != nil {
		dmn.Logger.Err("Failed to connect to RIBd")
		return err
	}
	return nil
}
func (dmn *L3Daemon) NewServer() {
	dmn.FSDaemon.NewServer()
	dmn.RibdSubSocketCh = make(chan []byte)
	dmn.RibdSubSocketErrCh = make(chan error)
}
