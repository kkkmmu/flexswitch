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
	"infra/fMgrd/objects"
	"time"
	"utils/eventUtils"
)

func (fMgr *FaultManager) GetAlarmStateObject(alarm *AlarmRBEntry) (aObj objects.AlarmState, err error) {
	aObj.OwnerId = int32(alarm.OwnerId)
	aObj.EventId = int32(alarm.EventId)
	evtKey := EventKey{
		DaemonId: alarm.OwnerId,
		EventId:  alarm.EventId,
	}
	fEnt, exist := fMgr.FaultEventMap[evtKey]
	if !exist {
		return aObj, errors.New("Error finding the entry in AlarmRB")
	}
	aObj.OwnerName = fEnt.FaultOwnerName
	aObj.EventName = fEnt.FaultEventName
	aObj.SrcObjName = fEnt.FaultSrcObjName
	aObj.Severity = fEnt.AlarmSeverity
	aObj.Description = alarm.Description
	aObj.OccuranceTime = alarm.OccuranceTime.String()
	aObj.SrcObjKey = alarm.SrcObjKey
	aObj.SrcObjUUID = alarm.SrcObjUUID

	if alarm.Resolved == true {
		aObj.ResolutionTime = alarm.ResolutionTime.String()
		aObj.ResolutionReason = getResolutionReason(alarm.ResolutionReason)
	} else {
		aObj.ResolutionTime = "N/A"
		aObj.ResolutionReason = "N/A"
	}
	return aObj, nil
}

func (fMgr *FaultManager) PublishAlarms(idx int) {
	fMgr.ARBRWMutex.RLock()
	aIntf := fMgr.AlarmRB.GetEntryFromRingBuffer(idx)
	fMgr.ARBRWMutex.RUnlock()
	alarm := aIntf.(AlarmRBEntry)
	aObj, err := fMgr.GetAlarmStateObject(&alarm)
	if err != nil {
		fMgr.logger.Err("Error Fetching the fault state object", err)
		return
	}
	msg, _ := json.Marshal(aObj)
	channel := aObj.OwnerName + "Alarms"
	fMgr.AlarmPubHdl.Publish("PUBLISH", channel, msg)
}

func (fMgr *FaultManager) GetBulkAlarmState(fromIdx int, count int) (*objects.AlarmStateGetInfo, error) {
	var retObj objects.AlarmStateGetInfo

	fMgr.ARBRWMutex.RLock()
	alarms := fMgr.AlarmRB.GetListOfEntriesFromRingBuffer()
	fMgr.ARBRWMutex.RUnlock()
	length := len(alarms)
	aState := make([]objects.AlarmState, count)

	var i int
	var j int

	for i, j = 0, fromIdx; i < count && j < length; j++ {
		aIntf := alarms[length-j-1]
		alarm := aIntf.(AlarmRBEntry)
		aObj, err := fMgr.GetAlarmStateObject(&alarm)
		if err != nil {
			continue
		}
		aState[i] = aObj
		i++
	}
	retObj.EndIdx = j
	retObj.Count = i
	if j != length {
		retObj.More = true
	}
	retObj.List = aState

	return &retObj, nil
}

func (fMgr *FaultManager) StartAlarmTimer(evt eventUtils.Event) *time.Timer {
	evtKey := EventKey{
		DaemonId: int(evt.OwnerId),
		EventId:  int(evt.EvtId),
	}

	alarmFunc := func() {
		fMgr.AMapRWMutex.Lock()
		if fMgr.AlarmMap[evtKey] == nil {
			fMgr.logger.Debug("Alarm Database does not exist, hence creating one")
			fMgr.AlarmMap[evtKey] = make(map[FaultObjKey]AlarmData)
		}

		aDataMapEnt, _ := fMgr.AlarmMap[evtKey]
		fObjKey, fObjKeyUUId, objKey, err := fMgr.generateFaultObjKey(evt.OwnerName, evt.SrcObjName, evt.SrcObjKey)
		if err != nil {
			fMgr.logger.Err("Fault Obj key, hence skipping alarm generation")
			fMgr.AMapRWMutex.Unlock()
			return
		}

		aDataEnt, exist := aDataMapEnt[fObjKey]
		if exist {
			fMgr.logger.Err("Alarm Data entry already exist, hence skipping this")
			fMgr.AMapRWMutex.Unlock()
			return
		}
		aDataEnt.AlarmListIdx = fMgr.AddAlarmEntryInRB(evt, objKey, fObjKeyUUId)
		fMgr.PublishAlarms(aDataEnt.AlarmListIdx)
		aDataEnt.AlarmSeqNumber = fMgr.AlarmSeqNumber
		fMgr.AlarmSeqNumber++
		aDataMapEnt[fObjKey] = aDataEnt
		fMgr.AlarmMap[evtKey] = aDataMapEnt
		fMgr.AMapRWMutex.Unlock()
	}

	return time.AfterFunc(fMgr.FaultToAlarmTransitionTime, alarmFunc)
}

func (fMgr *FaultManager) AddAlarmEntryInRB(evt eventUtils.Event, objKey, uuid string) int {
	aRBEnt := AlarmRBEntry{
		OwnerId:        int(evt.OwnerId),
		EventId:        int(evt.EvtId),
		OccuranceTime:  time.Now(),
		SrcObjKey:      objKey,
		SrcObjUUID:     uuid,
		AlarmSeqNumber: fMgr.AlarmSeqNumber,
		Description:    evt.Description,
	}

	fMgr.ARBRWMutex.Lock()
	idx, _ := fMgr.AlarmRB.InsertIntoRingBuffer(aRBEnt)
	fMgr.ARBRWMutex.Unlock()
	return idx
}

func (fMgr *FaultManager) StartAlarmRemoveTimer(evt eventUtils.Event, reason Reason) *time.Timer {
	evtKey := EventKey{
		DaemonId: int(evt.OwnerId),
		EventId:  int(evt.EvtId),
	}

	cFEnt, exist := fMgr.NonFaultEventMap[evtKey]
	if !exist {
		fMgr.logger.Err("Error finding the fault for fault clearing event")
		return nil
	}
	fEvtKey := EventKey{
		DaemonId: cFEnt.FaultOwnerId,
		EventId:  cFEnt.FaultEventId,
	}

	fObjKey, _, _, err := fMgr.generateFaultObjKey(evt.OwnerName, evt.SrcObjName, evt.SrcObjKey)
	if err != nil {
		fMgr.logger.Err("Error generating fault object key")
		return nil
	}

	alarmFunc := func() {
		fMgr.AMapRWMutex.Lock()
		aDataMapEnt, exist := fMgr.AlarmMap[fEvtKey]
		if !exist {
			fMgr.logger.Err("Alarm Database does not exist, hence skipping removal of Alarm")
			fMgr.AMapRWMutex.Unlock()
			return
		}
		aDataEnt, exist := aDataMapEnt[fObjKey]
		if !exist {
			fMgr.logger.Err("Alarm Data entry doesnot exist, hence skipping this")
			fMgr.AMapRWMutex.Unlock()
			return
		}
		aIntf := fMgr.AlarmRB.GetEntryFromRingBuffer(aDataEnt.AlarmListIdx)
		aRBData := aIntf.(AlarmRBEntry)
		if aRBData.AlarmSeqNumber == aDataEnt.AlarmSeqNumber {
			aRBData.ResolutionTime = time.Now()
			aRBData.ResolutionReason = reason
			aRBData.Resolved = true
			fMgr.ARBRWMutex.Lock()
			fMgr.AlarmRB.UpdateEntryInRingBuffer(aRBData, aDataEnt.AlarmListIdx)
			fMgr.ARBRWMutex.Unlock()
			fMgr.PublishAlarms(aDataEnt.AlarmListIdx)
			delete(aDataMapEnt, fObjKey)
			fMgr.AlarmMap[fEvtKey] = aDataMapEnt
		}
		fMgr.AMapRWMutex.Unlock()
	}

	return time.AfterFunc(fMgr.AlarmTransitionTime, alarmFunc)
}

func (fMgr *FaultManager) ClearExistingAlarms(evtKey EventKey, uuid string, reason Reason) {
	fMgr.AMapRWMutex.Lock()
	aDataMapEnt, exist := fMgr.AlarmMap[evtKey]
	if !exist {
		fMgr.AMapRWMutex.Unlock()
		return
	}
	for aDataKey, aDataEnt := range aDataMapEnt {
		fMgr.ARBRWMutex.Lock()
		aIntf := fMgr.AlarmRB.GetEntryFromRingBuffer(aDataEnt.AlarmListIdx)
		aRBData := aIntf.(AlarmRBEntry)
		if aRBData.AlarmSeqNumber == aDataEnt.AlarmSeqNumber {
			if uuid == "" || uuid == aRBData.SrcObjUUID {
				aRBData.ResolutionTime = time.Now()
				aRBData.ResolutionReason = reason
				aRBData.Resolved = true
				fMgr.AlarmRB.UpdateEntryInRingBuffer(aRBData, aDataEnt.AlarmListIdx)
				if aDataEnt.RemoveAlarmTimer != nil {
					aDataEnt.RemoveAlarmTimer.Stop()
				}
				delete(aDataMapEnt, aDataKey)
			}
		}
		fMgr.ARBRWMutex.Unlock()
		if aRBData.AlarmSeqNumber == aDataEnt.AlarmSeqNumber {
			if uuid == "" || uuid == aRBData.SrcObjUUID {
				fMgr.PublishAlarms(aDataEnt.AlarmListIdx)
			}
		}
	}
	if len(aDataMapEnt) == 0 {
		delete(fMgr.AlarmMap, evtKey)
	}
	fMgr.AMapRWMutex.Unlock()
}
