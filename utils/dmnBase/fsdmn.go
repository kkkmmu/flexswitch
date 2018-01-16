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
	"asicd/asicdCommonDefs"
	"asicdServices"
	"encoding/json"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	nanomsg "github.com/op/go-nanomsg"
	"io/ioutil"
	"strconv"
	"time"
	"utils/asicdClient"
	"utils/commonDefs"
	"utils/dbutils"
	"utils/ipcutils"
	"utils/keepalive"
	"utils/logging"
)

const (
	CLIENTS_FILE_NAME = "clients.json"
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type DmnClientBase struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
	IsConnected        bool
}

type AsicdClient struct {
	DmnClientBase
	ClientHdl *asicdServices.ASICDServicesClient
}

type FSBaseDmn struct {
	DmnName     string
	ParamsDir   string
	LogPrefix   string
	Logger      *logging.Writer
	DbHdl       *dbutils.DBUtil
	ClientsList []ClientJson
}

// @TODO: need to remove this struct, it duplicate and introducing bugs
type FSDaemon struct {
	*FSBaseDmn
	Asicdclnt      AsicdClient
	AsicdSubSocket *nanomsg.SubSocket
	// @ALERT ANY FUTURE DEVELOPER PLEASE DO NOT USE THIS, REFER NDP AND SEE HOW TO USE FSBaseDmn
	AsicdSubSocketCh    chan []byte
	AsicdSubSocketErrCh chan error
}

func (dmn *FSBaseDmn) InitLogger() (err error) {
	fmt.Println(dmn.LogPrefix, " Starting ", dmn.DmnName, "logger")
	dmnLogger, err := logging.NewLogger(dmn.DmnName, dmn.LogPrefix, true)
	if err != nil {
		fmt.Println("Failed to start the logger. Nothing will be logged...")
		return err
	}
	dmn.Logger = dmnLogger
	return err
}

func (dmn *FSBaseDmn) InitDBHdl() (err error) {
	dbHdl := dbutils.NewDBUtil(dmn.Logger)
	err = dbHdl.Connect()
	if err != nil {
		dmn.Logger.Err("Failed to dial out to Redis server")
		return err
	}
	dmn.DbHdl = dbHdl
	return err
}

func (dmn *FSBaseDmn) Init() bool {
	err := dmn.InitLogger()
	if err != nil {
		return false
	}
	err = dmn.InitDBHdl()
	if err != nil {
		return false
	}
	configFile := dmn.ParamsDir + "clients.json"
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		dmn.Logger.Info(fmt.Sprintln("Error in reading configuration file ", configFile))
		return false
	}
	err = json.Unmarshal(bytes, &dmn.ClientsList)
	if err != nil {
		dmn.Logger.Info("Error in Unmarshalling Json")
		return false
	}
	dmn.Logger.Info("Base daemon init completed")
	return true
}

func (dmn *FSBaseDmn) GetParams() string {
	paramsDir := flag.String("params", "./params", "Params directory")
	flag.Parse()
	dirName := *paramsDir
	if dirName[len(dirName)-1] != '/' {
		dirName = dirName + "/"
	}
	return dirName
}

func (dmn *FSBaseDmn) StartKeepAlive() {
	go keepalive.InitKeepAlive(dmn.DmnName, dmn.ParamsDir)
}

func (dmn *FSDaemon) StartKeepAlive() {
	dmn.FSBaseDmn.StartKeepAlive()
}

func NewBaseDmn(dmnName, logPrefix string) *FSBaseDmn {
	var dmn = new(FSBaseDmn)
	dmn.DmnName = dmnName
	dmn.LogPrefix = logPrefix
	dmn.ParamsDir = dmn.GetParams()
	return dmn
}

func (dmn *FSDaemon) ConnectToAsicd() error {
	var err error
	for _, client := range dmn.FSBaseDmn.ClientsList {
		if client.Name == "asicd" {
			dmn.Logger.Info(fmt.Sprintln("found  asicd at port ", client.Port))
			dmn.Asicdclnt.Address = "localhost:" + strconv.Itoa(client.Port)
			dmn.Asicdclnt.Transport, dmn.Asicdclnt.PtrProtocolFactory, err =
				ipcutils.CreateIPCHandles(dmn.Asicdclnt.Address)
			if dmn.Asicdclnt.Transport != nil && dmn.Asicdclnt.PtrProtocolFactory != nil {
				dmn.Asicdclnt.ClientHdl =
					asicdServices.NewASICDServicesClientFactory(dmn.Asicdclnt.Transport,
						dmn.Asicdclnt.PtrProtocolFactory)
				dmn.Asicdclnt.IsConnected = true
			} else {
				dmn.Logger.Info(fmt.Sprintf("Failed to connect to Asicd, retrying until connection is successful"))
				count := 0
				ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
				for _ = range ticker.C {
					dmn.Asicdclnt.Transport, dmn.Asicdclnt.PtrProtocolFactory, err =
						ipcutils.CreateIPCHandles(dmn.Asicdclnt.Address)
					if err == nil {
						ticker.Stop()
						break
					}
					count++
					if (count % 10) == 0 {
						dmn.Logger.Info("Still can't connect to Asicd, retrying...")
					}
				}
			}
		}
	}
	return err
}

func (dmn *FSDaemon) CreateASICdSubscriber() error {
	dmn.Logger.Info("Listen for ASICd updates")
	err := dmn.ListenForASICdUpdates(asicdCommonDefs.PUB_SOCKET_ADDR)
	if err != nil {
		dmn.Logger.Err("Error initializing ASICD subscriber")
		return err
	}
	for {
		dmn.Logger.Info("Read on ASICd subscriber socket...")
		asicdrxBuf, err := dmn.AsicdSubSocket.Recv(0)
		if err != nil {
			dmn.Logger.Err(fmt.Sprintln("Recv on ASICd subscriber socket failed with error:", err))
			dmn.AsicdSubSocketErrCh <- err
			continue
		}
		dmn.Logger.Info(fmt.Sprintln("ASIC subscriber recv returned:", asicdrxBuf))
		dmn.AsicdSubSocketCh <- asicdrxBuf
	}
	return nil
}

func (dmn *FSDaemon) ListenForASICdUpdates(address string) error {
	var err error
	if dmn.AsicdSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		dmn.Logger.Err(fmt.Sprintln("Failed to create ASICd subscribe socket, error:", err))
		return err
	}

	if _, err = dmn.AsicdSubSocket.Connect(address); err != nil {
		dmn.Logger.Err(fmt.Sprintln("Failed to connect to ASICd publisher socket, address:", address, "error:", err))
		return err
	}

	if err = dmn.AsicdSubSocket.Subscribe(""); err != nil {
		dmn.Logger.Err(fmt.Sprintln("Failed to subscribe to \"\" on ASICd subscribe socket, error:", err))
		return err
	}

	dmn.Logger.Info(fmt.Sprintln("Connected to ASICd publisher at address:", address))
	if err = dmn.AsicdSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		dmn.Logger.Err(fmt.Sprintln("Failed to set the buffer size for ASICd publisher socket, error:", err))
		return err
	}
	return nil
}

func (dmn *FSDaemon) InitSubscribers([]string) (err error) {
	go dmn.CreateASICdSubscriber()
	return err
}

func (dmn *FSDaemon) ConnectToServers() error {
	err := dmn.ConnectToAsicd()
	if err != nil {
		return err
	}
	return nil
}

// @TODO: remove this when l2 & l3 daemons have moved to Plugin Model
func (dmn *FSDaemon) Init(dmnName, logPrefix string) bool {
	dmn.FSBaseDmn = NewBaseDmn(dmnName, logPrefix)
	return dmn.FSBaseDmn.Init()
}

func (dmn *FSDaemon) NewServer() {
	dmn.AsicdSubSocketCh = make(chan []byte)
	dmn.AsicdSubSocketErrCh = make(chan error)
}

func (dmn *FSBaseDmn) GetLogger() *logging.Writer {
	return dmn.Logger
}

func (dmn *FSBaseDmn) GetDbHdl() *dbutils.DBUtil {
	return dmn.DbHdl
}

func (dmn *FSBaseDmn) InitSwitch(plugin, dmnName, logPrefix string, switchHdl commonDefs.AsicdClientStruct) asicdClient.AsicdClientIntf {
	// @TODO: need to change second argument
	return asicdClient.NewAsicdClientInit(plugin, dmn.ParamsDir+"clients.json", switchHdl)

}
