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

package main

import (
	"infra/platformd/api"
	"infra/platformd/rpc"
	"infra/platformd/server"
	"strconv"
	"strings"
	"utils/dmnBase"
)

const (
	DMN_NAME = "platformd"
	CFG_FILE = "platformd.conf"
	//CFG_DIR  = "/etc/flexswitch/"
)

type platformDaemon struct {
	*dmnBase.FSBaseDmn
	server    *server.PlatformdServer
	rpcServer *rpc.RPCServer
}

var dmn platformDaemon

func main() {
	var err error
	// Get base daemon handle and initialize
	dmn.FSBaseDmn = dmnBase.NewBaseDmn(DMN_NAME, DMN_NAME)
	ok := dmn.Init()
	if ok == false {
		panic("PlatformD Base daemon initialization failed")
	}

	//Get server handle and start server
	cfgFileName := dmn.ParamsDir + "/" + CFG_FILE
	InitParams := &server.InitParams{
		DmnName:     DMN_NAME,
		ParamsDir:   dmn.ParamsDir,
		CfgFileName: cfgFileName,
		EventDbHdl:  dmn.DbHdl,
		Logger:      dmn.FSBaseDmn.Logger,
	}
	dmn.server, err = server.NewPlatformdServer(InitParams)
	if err != nil {
		panic("Unable to Initialize Platform Daemon Plugin")
		return
	}
	go dmn.server.Serve()

	// Initialize api layer
	api.InitApiLayer(dmn.server)

	//Start keepalive for watchdog
	dmn.StartKeepAlive()

	//Wait for server started msg
	_ = <-dmn.server.InitCompleteCh

	//Get RPC server handle
	var rpcServerAddr string
	for _, value := range dmn.FSBaseDmn.ClientsList {
		if value.Name == strings.ToLower(DMN_NAME) {
			rpcServerAddr = "localhost:" + strconv.Itoa(value.Port)
			break
		}
	}
	if rpcServerAddr == "" {
		panic("Platform Daemon is not part of the system profile")
	}
	dmn.rpcServer = rpc.NewRPCServer(rpcServerAddr, dmn.FSBaseDmn.Logger, dmn.DbHdl)

	//Start RPC server
	dmn.FSBaseDmn.Logger.Info("Platform Daemon server started")
	dmn.rpcServer.Serve()
	panic("Platform Daemon RPC Server terminated")
}
