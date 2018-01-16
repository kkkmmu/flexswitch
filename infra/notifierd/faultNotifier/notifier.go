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

package faultNotifier

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"infra/notifierd/objects"
	"net/http"
	"utils/logging"
)

type Notifier struct {
	logger   logging.LoggerIntf
	DmnList  []string
	upgrader websocket.Upgrader
}

func NewNotifier(param *objects.NotifierParam) *Notifier {
	notifier := &Notifier{}
	notifier.logger = param.Logger
	notifier.DmnList = append(notifier.DmnList, param.DmnList...)
	notifier.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	return notifier
}

func (notifier *Notifier) InitNotifier(subHdl *redis.PubSubConn) error {
	dbHdl, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	subHdl.Conn = dbHdl

	for _, daemon := range notifier.DmnList {
		channel := daemon + "Faults"
		err := subHdl.Subscribe(channel)
		if err != nil {
			notifier.logger.Err("Error Initializing event subscriber for", daemon)
		}
	}

	return nil
}

func (notifier *Notifier) ProcessNotification(w http.ResponseWriter, r *http.Request) {
	conn, err := notifier.upgrader.Upgrade(w, r, nil)
	if err != nil {
		notifier.logger.Err("Error setting up websocket connection faults", err)
		return
	}
	defer conn.Close()
	var subHdl redis.PubSubConn
	err = notifier.InitNotifier(&subHdl)
	if err != nil {
		notifier.logger.Err("Error setting up fault notifier", err)
		return
	}

	for {
		switch n := subHdl.Receive().(type) {
		case redis.Message:
			err = conn.WriteMessage(websocket.TextMessage, n.Data)
			if err != nil {
				notifier.logger.Err("Error sending faults. Hence Closing connection", err)
				return
			}
		case redis.Subscription:
			if n.Count == 0 {
				notifier.logger.Err("Invalid data recevied while Subscription")
			}
		case error:
			notifier.logger.Err("Error while Subscription", n)
		}
	}
}
