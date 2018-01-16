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
	"infra/fMgrd/objects"
	//"models/events"
	"strings"
	"sync"
	"time"
	"utils/dbutils"
	"utils/eventUtils"
	"utils/logging"
	"utils/ringBuffer"
)

type PubIntf interface {
	Publish(string, interface{}, interface{})
	Connect() error
}

type FaultManager struct {
	logger                     logging.LoggerIntf
	dbHdl                      dbutils.DBIntf
	EventCh                    chan []byte
	PauseEventProcessCh        chan bool
	PauseEventProcessAckCh     chan bool
	FaultEventMap              map[EventKey]FaultDetail
	NonFaultEventMap           map[EventKey]NonFaultDetail
	OwnerEventNameMap          map[EventKeyStr]EventKey
	FMapRWMutex                sync.RWMutex
	FaultMap                   map[EventKey]FaultDataMap
	AMapRWMutex                sync.RWMutex
	AlarmMap                   map[EventKey]AlarmDataMap
	FRBRWMutex                 sync.RWMutex
	FaultRB                    *ringBuffer.RingBuffer
	ARBRWMutex                 sync.RWMutex
	AlarmRB                    *ringBuffer.RingBuffer
	DaemonList                 []string
	FaultSeqNumber             uint64
	AlarmSeqNumber             uint64
	FaultToAlarmTransitionTime time.Duration
	AlarmTransitionTime        time.Duration
	FaultPubHdl                PubIntf
	AlarmPubHdl                PubIntf
}

func NewFaultManager(logger logging.LoggerIntf) *FaultManager {
	fMgr := &FaultManager{}
	fMgr.logger = logger
	fMgr.EventCh = make(chan []byte, 1000)
	fMgr.PauseEventProcessCh = make(chan bool, 1)
	fMgr.PauseEventProcessAckCh = make(chan bool, 1)
	fMgr.FaultEventMap = make(map[EventKey]FaultDetail)
	fMgr.NonFaultEventMap = make(map[EventKey]NonFaultDetail)
	fMgr.OwnerEventNameMap = make(map[EventKeyStr]EventKey)
	fMgr.FaultMap = make(map[EventKey]FaultDataMap) //Existing Faults
	fMgr.AlarmMap = make(map[EventKey]AlarmDataMap) //Existing Alarm
	fMgr.FaultRB = new(ringBuffer.RingBuffer)
	fMgr.FaultRB.SetRingBufferCapacity(100000) // Max 100000 entries in fault database
	fMgr.AlarmRB = new(ringBuffer.RingBuffer)
	fMgr.AlarmRB.SetRingBufferCapacity(100000) // Max 100000 entries in alarm database
	fMgr.FaultSeqNumber = 0
	fMgr.AlarmSeqNumber = 0
	fMgr.FaultToAlarmTransitionTime = time.Duration(3) * time.Second
	fMgr.AlarmTransitionTime = time.Duration(3) * time.Second
	fMgr.dbHdl = dbutils.NewDBUtil(logger)
	fMgr.FaultPubHdl = dbutils.NewDBUtil(logger)
	fMgr.AlarmPubHdl = dbutils.NewDBUtil(logger)
	return fMgr
}

func (fMgr *FaultManager) InitFaultManager() error {
	err := fMgr.initFMgrDS()
	if err != nil {
		fMgr.logger.Err(fmt.Sprintln("Error Initializing Fault Manager DS:", err))
		return err
	}
	err = fMgr.dbHdl.Connect()
	if err != nil {
		fMgr.logger.Err(fmt.Sprintln("Error Initializing Fault Manager DB Handler:", err))
		return err
	}
	err = fMgr.FaultPubHdl.Connect()
	if err != nil {
		fMgr.logger.Err(fmt.Sprintln("Error Initializing Fault Publisher Handler:", err))
	}
	err = fMgr.AlarmPubHdl.Connect()
	if err != nil {
		fMgr.logger.Err(fmt.Sprintln("Error Initializing Alarm Publisher Handler:", err))
	}
	go fMgr.EventProcessor()
	return err
}

func (fMgr *FaultManager) EventProcessor() {
	for {
		select {
		case msg := <-fMgr.EventCh:
			var evt eventUtils.Event
			err := json.Unmarshal(msg, &evt)
			if err != nil {
				fMgr.logger.Err(fmt.Sprintln("Unable to Unmarshal the byte stream", err))
				continue
			}
			fMgr.logger.Debug(fmt.Sprintln("OwnerId:", evt.OwnerId))
			fMgr.logger.Debug(fmt.Sprintln("OwnerName:", evt.OwnerName))
			fMgr.logger.Debug(fmt.Sprintln("EvtId:", evt.EvtId))
			fMgr.logger.Debug(fmt.Sprintln("EventName:", evt.EventName))
			fMgr.logger.Debug(fmt.Sprintln("Timestamp:", evt.TimeStamp))
			fMgr.logger.Debug(fmt.Sprintln("Description:", evt.Description))
			fMgr.logger.Debug(fmt.Sprintln("SrcObjName:", evt.SrcObjName))

			fMgr.processEvents(evt)
		case _ = <-fMgr.PauseEventProcessCh:
			fMgr.PauseEventProcessAckCh <- true
			<-fMgr.PauseEventProcessCh
		}
	}
}

func (fMgr *FaultManager) initFMgrDS() error {
	evtMap := make(map[EventKey]EvtDetail)
	evtJson, err := eventUtils.ParseEventsJson()
	if err != nil {
		fMgr.logger.Err(fmt.Sprintln("Error Parsing the events.json", err))
		return err
	}
	for _, daemon := range evtJson.DaemonEvents {
		fMgr.logger.Debug(fmt.Sprintln("daemon.DaemonName:", daemon.DaemonName))
		fMgr.DaemonList = append(fMgr.DaemonList, daemon.DaemonName)
		for _, evt := range daemon.EventList {
			fId := EventKey{
				DaemonId: int(daemon.DaemonId),
				EventId:  int(evt.EventId),
			}
			fName := EventKeyStr{
				OwnerName: daemon.DaemonName,
				EventName: evt.EventName,
			}
			fMgr.OwnerEventNameMap[fName] = fId
			evtEnt, exist := evtMap[fId]
			if exist {
				fMgr.logger.Err(fmt.Sprintln("Duplicate entry found"))
				continue
			}
			if evt.IsFault == true {
				evtEnt.IsFault = true
				evtEnt.IsClearingEvent = false
				evtEnt.RaiseFault = evt.Fault.RaiseFault
				evtEnt.ClearingEventId = evt.Fault.ClearingEventId
				evtEnt.ClearingDaemonId = evt.Fault.ClearingDaemonId
				evtEnt.AlarmSeverity = evt.Fault.AlarmSeverity
				evtEnt.OwnerName = daemon.DaemonName
				evtEnt.EventName = evt.EventName
				evtEnt.SrcObjName = evt.SrcObjName
			} else {
				evtEnt.IsFault = false
				evtEnt.IsClearingEvent = false
				evtEnt.RaiseFault = false
				evtEnt.ClearingEventId = -1
				evtEnt.ClearingDaemonId = -1
				evtEnt.AlarmSeverity = ""
				evtEnt.OwnerName = daemon.DaemonName
				evtEnt.EventName = evt.EventName
				evtEnt.SrcObjName = evt.SrcObjName
			}
			evtMap[fId] = evtEnt
		}
	}

	for fId, evt := range evtMap {
		if evt.IsFault == true {
			cFId := EventKey{
				DaemonId: evt.ClearingDaemonId,
				EventId:  evt.ClearingEventId,
			}
			cEvt, exist := evtMap[cFId]
			if !exist {
				fMgr.logger.Err(fmt.Sprintln("No clearing event found for fault:", fId))
				continue
			}

			cEvt.IsClearingEvent = true
			evtMap[cFId] = cEvt
		}
	}

	for fId, evt := range evtMap {
		if evt.IsFault == true {
			evtEnt, _ := fMgr.FaultEventMap[fId]
			evtEnt.RaiseFault = evt.RaiseFault
			evtEnt.ClearingEventId = evt.ClearingEventId
			evtEnt.ClearingDaemonId = evt.ClearingDaemonId
			evtEnt.FaultOwnerName = evt.OwnerName
			evtEnt.FaultEventName = evt.EventName
			evtEnt.FaultSrcObjName = evt.SrcObjName
			evtEnt.AlarmSeverity = evt.AlarmSeverity
			fMgr.FaultEventMap[fId] = evtEnt
			cFId := EventKey{
				DaemonId: evtEnt.ClearingDaemonId,
				EventId:  evtEnt.ClearingEventId,
			}
			cEvtEnt, _ := fMgr.NonFaultEventMap[cFId]
			cEvtEnt.FaultOwnerId = fId.DaemonId
			cEvtEnt.FaultEventId = fId.EventId
			fMgr.NonFaultEventMap[cFId] = cEvtEnt
		} else {
			evtEnt, _ := fMgr.NonFaultEventMap[fId]
			evtEnt.IsClearingEvent = evt.IsClearingEvent
			fMgr.NonFaultEventMap[fId] = evtEnt
		}
	}
	return nil
}

func (fMgr *FaultManager) faultEnable(evtKey EventKey, enable bool) (retVal bool, err error) {
	_, exist := fMgr.FaultEventMap[evtKey]
	if !exist {
		err = errors.New("Unable to find the corresponding fault event")
	} else {
		if enable == false {
			err = fMgr.DisableFaults(evtKey)
			if err == nil {
				fMgr.ClearExistingFaults(evtKey, "", FAULTDISABLED)
				fMgr.ClearExistingAlarms(evtKey, "", FAULTDISABLED)
				retVal = true
			}
		} else {
			err = fMgr.EnableFaults(evtKey)
			if err == nil {
				retVal = true
			}
		}
	}
	return retVal, err
}

func (fMgr *FaultManager) FaultEnableAction(config *objects.FaultEnable) (retVal bool, err error) {
	fMgr.PauseEventProcessCh <- true
	<-fMgr.PauseEventProcessAckCh
	if strings.ToLower(config.EventName) == objects.ALL_EVENTS {
		ownerName := strings.ToLower(config.OwnerName)
		for evtKeyStr, evtKey := range fMgr.OwnerEventNameMap {
			if strings.ToLower(evtKeyStr.OwnerName) == ownerName {
				_, exist := fMgr.FaultEventMap[evtKey]
				if exist {
					retVal, err = fMgr.faultEnable(evtKey, config.Enable)
				}
			} else {
				continue
			}
		}
	} else {
		evtKeyStr := EventKeyStr{
			OwnerName: config.OwnerName,
			EventName: config.EventName,
		}
		evtKey, exist := fMgr.OwnerEventNameMap[evtKeyStr]
		if !exist {
			err = errors.New("Unable to find the corresponding event")
		} else {
			retVal, err = fMgr.faultEnable(evtKey, config.Enable)
		}
	}
	fMgr.PauseEventProcessCh <- true
	return retVal, err
}

func (fMgr *FaultManager) DisableFaults(evtKey EventKey) error {
	fEnt, _ := fMgr.FaultEventMap[evtKey]
	if fEnt.RaiseFault == false {
		return errors.New("Fault is already disabled")
	}
	fEnt.RaiseFault = false
	fMgr.FaultEventMap[evtKey] = fEnt
	return nil
}

func (fMgr *FaultManager) EnableFaults(evtKey EventKey) error {
	fEnt, _ := fMgr.FaultEventMap[evtKey]
	if fEnt.RaiseFault == true {
		return errors.New("Fault is already enabled")
	}
	fEnt.RaiseFault = true
	fMgr.FaultEventMap[evtKey] = fEnt
	return nil
}

func (fMgr *FaultManager) FaultClearAction(config *objects.FaultClear) (retVal bool, err error) {
	fMgr.PauseEventProcessCh <- true
	<-fMgr.PauseEventProcessAckCh
	evtKeyStr := EventKeyStr{
		OwnerName: config.OwnerName,
		EventName: config.EventName,
	}
	evtKey, exist := fMgr.OwnerEventNameMap[evtKeyStr]
	if !exist {
		err = errors.New("Unable to find the corresponding event")
	} else {
		fEnt, exist := fMgr.FaultEventMap[evtKey]
		if !exist {
			err = errors.New("Unable to find the corresponding faulty event")
		} else {
			if fEnt.RaiseFault == true {
				fMgr.ClearExistingFaults(evtKey, config.SrcObjUUID, FAULTCLEARED)
				fMgr.ClearExistingAlarms(evtKey, config.SrcObjUUID, FAULTCLEARED)
				retVal = true
			} else {
				err = errors.New("Fault for this Event is already disabled, nothing to be cleared")
			}
		}
	}
	fMgr.PauseEventProcessCh <- true
	return retVal, err
}
