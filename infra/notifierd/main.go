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
	"infra/notifierd/api"
	"infra/notifierd/rpc"
	"infra/notifierd/server"
	"strconv"
	"utils/dmnBase"
)

const (
	DMN_NAME = "notifierd"
)

type nMgrDaemon struct {
	*dmnBase.FSBaseDmn
	server    *server.NMGRServer
	rpcServer *rpc.RPCServer
}

var dmn nMgrDaemon

func main() {
	// Get base daemon handle and initialize
	dmn.FSBaseDmn = dmnBase.NewBaseDmn(DMN_NAME, DMN_NAME)
	ok := dmn.Init()
	if ok == false {
		panic("Notifier Daemon: Base Daemon Initialization failed")
	}

	// Get server handle and start server
	svrInitParams := &server.ServerInitParams{
		ParamsDir: dmn.ParamsDir,
		Logger:    dmn.FSBaseDmn.Logger,
	}
	dmn.server = server.NewNMGRServer(svrInitParams)
	go dmn.server.StartServer()

	//Initialize API layer
	api.InitApiLayer(dmn.server)

	//Get RPC server handle
	var rpcServerAddr string
	for _, value := range dmn.FSBaseDmn.ClientsList {
		if value.Name == "notifierd" {
			rpcServerAddr = "localhost:" + strconv.Itoa(value.Port)
			break
		}
	}

	if rpcServerAddr == "" {
		panic("Notifier Daemon is not part of system profile")
	}

	dmn.rpcServer = rpc.NewRPCServer(rpcServerAddr, dmn.FSBaseDmn.Logger)

	// Start Keep Alive for watchdog
	dmn.StartKeepAlive()

	_ = <-dmn.server.InitDone

	//Start RPC server
	dmn.FSBaseDmn.Logger.Info("Notifier Daemon server started")
	dmn.rpcServer.Serve()
	panic("Notifier Daemon RPC server terminated")
}
