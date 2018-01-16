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

package pluginManager

import (
	"errors"
	"fmt"
	"infra/platformd/objects"
	"utils/logging"
)

type LedManager struct {
	logger logging.LoggerIntf
	plugin PluginIntf
	ledCnt int
}

var LedMgr LedManager

func (lMgr *LedManager) Init(logger logging.LoggerIntf, plugin PluginIntf) {
	lMgr.logger = logger
	lMgr.plugin = plugin
	lMgr.ledCnt = lMgr.plugin.GetMaxNumOfLed()
	lMgr.logger.Info("Led Manager Init()")
}

func (lMgr *LedManager) Deinit() {
	lMgr.logger.Info("Led Manager Deinit()")
}

func (lMgr *LedManager) GetLedState(ledId int32) (*objects.LedState, error) {
	var ledObj objects.LedState

	if lMgr.plugin == nil {
		return nil, errors.New("Invalid Led platform plugin")
	}

	ledState, err := lMgr.plugin.GetLedState(ledId)
	if err != nil {
		return nil, err
	}

	ledObj.LedId = ledId
	ledObj.LedIdentify = ledState.LedIdentify
	ledObj.LedColor = ledState.LedColor
	ledObj.LedState = ledState.LedState

	return &ledObj, err
}

func (lMgr *LedManager) GetBulkLedState(fromIdx int, cnt int) (*objects.LedStateGetInfo, error) {
	var retObj objects.LedStateGetInfo

	if lMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}

	if fromIdx >= lMgr.ledCnt {
		return nil, errors.New("Invalid range")
	}

	if fromIdx+cnt > lMgr.ledCnt {
		retObj.EndIdx = lMgr.ledCnt
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = lMgr.ledCnt - retObj.EndIdx + 1
	}

	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		ledId := int32(idx)
		obj, err := lMgr.GetLedState(ledId)
		if err != nil {
			lMgr.logger.Err(fmt.Sprintln("Error getting the Led state for ledId:", ledId))
		}
		retObj.List = append(retObj.List, obj)
	}

	return &retObj, nil
}
