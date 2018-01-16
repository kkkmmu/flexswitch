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

package faultMgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"models/events"
	"strings"
)

func (fMgr *FaultManager) generateFaultObjKey(ownerName, srcObjName string, srcObjKey interface{}) (FaultObjKey, string, string, error) {
	objKey, dbObjKey, err := getEventObjKey(ownerName, srcObjName, srcObjKey)
	if err != nil {
		fMgr.logger.Err("Unable to find the Obj Key", srcObjName, srcObjKey, err)
		return "", "", "", errors.New(fmt.Sprintln("Unable to find the ObjKey of", srcObjName, srcObjKey, err))
	}

	srcObjUUID, err := fMgr.getUUID(srcObjName, dbObjKey)
	if err != nil {
		fMgr.logger.Err("Unable to find the UUID of", srcObjName, srcObjKey, err)
		return "", "", "", errors.New(fmt.Sprintln("Unable to find the UUID of", srcObjName, srcObjKey, err))
	}
	return FaultObjKey(fmt.Sprintf("%s#%s#%s", srcObjName, objKey, srcObjUUID)), srcObjUUID, objKey, err
}

func (fMgr *FaultManager) getUUID(srcObjName, dbObjKey string) (uuid string, err error) {
	return fMgr.dbHdl.GetUUIDFromObjKey(dbObjKey)
}

func getResolutionReason(reason Reason) string {
	switch reason {
	case AUTOCLEARED:
		return "Automatically Cleared"
	case FAULTDISABLED:
		return "Cleared because of FaultEnable(Enable=false) Action"
	case FAULTCLEARED:
		return "Cleared because of FaultClear Action"
	}
	return "Unknown"
}

func getEventObjKey(ownerName, srcObjName string, srcObjKey interface{}) (objKey string, dbObjKey string, err error) {
	objKeyMap, _ := events.EventKeyMap[strings.ToUpper(ownerName)]
	obj, _ := objKeyMap[srcObjName]
	bytes, _ := json.Marshal(srcObjKey)

	return obj.GetObjDBKey(bytes)
}
