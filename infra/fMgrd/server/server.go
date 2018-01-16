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
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"infra/fMgrd/faultMgr"
	"time"
	"utils/logging"
)

type FMGRServer struct {
	Logger    logging.LoggerIntf
	dbHdl     redis.Conn
	subHdl    redis.PubSubConn
	fMgr      *faultMgr.FaultManager
	InitDone  chan bool
	ReqChan   chan *ServerRequest
	ReplyChan chan interface{}
}

func NewFMGRServer(logger *logging.Writer) *FMGRServer {
	fMgrServer := &FMGRServer{}
	fMgrServer.Logger = logger
	fMgrServer.InitDone = make(chan bool)
	fMgrServer.ReqChan = make(chan *ServerRequest)
	fMgrServer.ReplyChan = make(chan interface{})
	return fMgrServer
}

func (server *FMGRServer) dial() (redis.Conn, error) {
	retryCount := 0
	ticker := time.NewTicker(2 * time.Second)
	for _ = range ticker.C {
		retryCount += 1
		dbHdl, err := redis.Dial("tcp", ":6379")
		if err != nil {
			if retryCount%100 == 0 {
				server.Logger.Err(fmt.Sprintln("Failed to dail out to Redis server. Retrying connection. Num of retries = ", retryCount))
			}
		} else {
			return dbHdl, nil
		}
	}
	err := errors.New("Error opening db handler")
	return nil, err
}

func (server *FMGRServer) Subscriber() {
	for {
		switch n := server.subHdl.Receive().(type) {
		case redis.Message:
			server.fMgr.EventCh <- n.Data
		case redis.Subscription:
			if n.Count == 0 {
				server.Logger.Err("Empty data Received")
			}
		case error:
			server.Logger.Err(fmt.Sprintf("error: %v\n", n))
			return
		}
	}
	server.subHdl.Unsubscribe()
	//server.subHdl.PUnsubscribe()
}

func (server *FMGRServer) InitSubscriber() error {
	var errMsg string
	for _, daemon := range server.fMgr.DaemonList {
		err := server.subHdl.Subscribe(daemon)
		if err != nil {
			errMsg = fmt.Sprintf("%s : %s", errMsg, err)
		}
	}

	if errMsg == "" {
		return nil
	}

	return errors.New(fmt.Sprintln("Error Initializing Subscriber:", errMsg))
}

func (server *FMGRServer) InitServer() error {
	server.fMgr = faultMgr.NewFaultManager(server.Logger)
	err := server.fMgr.InitFaultManager()
	if err != nil {
		server.Logger.Err(fmt.Sprintln(err))
		return err
	}
	server.dbHdl, err = server.dial()
	if err != nil {
		server.Logger.Err(fmt.Sprintln(err))
		return err
	}
	//defer server.dbHdl.Close()
	server.subHdl = redis.PubSubConn{Conn: server.dbHdl}

	err = server.InitSubscriber()
	if err != nil {
		server.Logger.Err(fmt.Sprintln("Error in Initializing Subscriber", err))
		return err
	}
	go server.Subscriber()
	return err
}

func (server *FMGRServer) handleRPCRequest(req *ServerRequest) {
	server.Logger.Debug(fmt.Sprintln("Calling handle RPC Request for:", *req))
	switch req.Op {
	case GET_BULK_FAULT_STATE:
		var retObj GetBulkFaultStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = server.getBulkFaultState(val.FromIdx, val.Count)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_BULK_ALARM_STATE:
		var retObj GetBulkAlarmStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = server.getBulkAlarmState(val.FromIdx, val.Count)
		}
		server.ReplyChan <- interface{}(&retObj)
	case FAULT_ENABLE_ACTION:
		var retObj FaultEnableActionOutArgs
		if val, ok := req.Data.(*FaultEnableActionInArgs); ok {
			retObj.RetVal, retObj.Err = server.faultEnableAction(val.Config)
		}
		server.ReplyChan <- interface{}(&retObj)
	case FAULT_CLEAR_ACTION:
		var retObj FaultClearActionOutArgs
		if val, ok := req.Data.(*FaultClearActionInArgs); ok {
			retObj.RetVal, retObj.Err = server.faultClearAction(val.Config)
		}
		server.ReplyChan <- interface{}(&retObj)
	default:
		server.Logger.Err(fmt.Sprintln("Error: Server received unrecognized request - ", req.Op))
	}
}

func (server *FMGRServer) StartServer() {
	server.InitServer()
	server.InitDone <- true
	for {
		select {
		case req := <-server.ReqChan:
			server.Logger.Info(fmt.Sprintln("Server request received - ", *req))
			server.handleRPCRequest(req)
		}
	}
}
