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
	"flag"
	"fmt"
	"infra/sysd/rpc"
	"infra/sysd/server"
	"utils/dbutils"
	"utils/logging"
)

/*
#cgo LDFLAGS: -L../../../../external/src/github.com/netfilter/libiptables/lib
*/
import "C"

func main() {
	fmt.Println("Starting system daemon")
	paramsDir := flag.String("params", "./params", "Params directory")
	flag.Parse()
	fileName := *paramsDir
	if fileName[len(fileName)-1] != '/' {
		fileName = fileName + "/"
	}

	logger, err := logging.NewLogger("sysd", "SYSTEM", false)
	if err != nil {
		fmt.Println("Failed to start the logger. Nothing will be logged...")
	}
	logger.Info("Started the logger successfully.")

	dbHdl := dbutils.NewDBUtil(logger)
	if err := dbHdl.Connect(); err != nil {
		return
	}

	clientsFileName := fileName + "clients.json"

	logger.Info("Starting Sysd Server...")
	sysdServer := server.NewSYSDServer(logger, dbHdl, fileName)
	// Initialize sysd server
	sysdServer.InitServer()
	// Start signal handler first
	go sysdServer.SigHandler(dbHdl)

	// Start sysd server
	go sysdServer.StartServer()
	<-sysdServer.ServerStartedCh

	// Read IpTableAclConfig during restart case
	sysdServer.ReadIpAclConfigFromDB(dbHdl)

	logger.Info("Starting Sysd Config listener...")
	confIface := rpc.NewSYSDHandler(logger, sysdServer)
	rpc.StartServer(logger, confIface, clientsFileName)
}
