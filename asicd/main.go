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
	"asicd/rpc"
	"flag"
	"fmt"
	"strconv"
	"utils/keepalive"
	"utils/logging"
)

var logger *logging.Writer
var asicdServer *rpc.AsicDaemonServerInfo

func main() {
	var err error
	fmt.Println("Starting asicd daemon")
	paramsDirStr := flag.String("params", "", "Directory Location for config file")
	flag.Parse()
	paramsDir := *paramsDirStr
	if paramsDir[len(paramsDir)-1] != '/' {
		paramsDir = paramsDir + "/"
	}

	//Initialize logger
	logger, err = logging.NewLogger("asicd", "ASICD :", true)
	if err != nil {
		fmt.Println("Failed to start the logger. Nothing will be logged...")
	}

	//Parse cfg file
	cfgFileInfo := parseConfigFile(paramsDir)

	// Start keepalive routine
	go keepalive.InitKeepAlive("asicd", paramsDir)

	//Start rpc server
	asicdServer = rpc.NewAsicdServer("localhost:"+strconv.Itoa(cfgFileInfo.thriftServerPort), logger)
	logger.Info("ASICD: server started")
	asicdServer.Server.Serve()
}
