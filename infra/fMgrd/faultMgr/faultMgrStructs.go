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
	"time"
)

type EventKey struct {
	DaemonId int
	EventId  int
}

type EventKeyStr struct {
	OwnerName string
	EventName string
}

type NonFaultDetail struct {
	IsClearingEvent bool
	FaultEventId    int
	FaultOwnerId    int
}

type FaultDetail struct {
	RaiseFault       bool
	ClearingEventId  int
	ClearingDaemonId int
	AlarmSeverity    string
	FaultOwnerName   string
	FaultEventName   string
	FaultSrcObjName  string
}

type EvtDetail struct {
	IsClearingEvent  bool
	IsFault          bool
	RaiseFault       bool
	ClearingEventId  int
	ClearingDaemonId int
	AlarmSeverity    string
	OwnerName        string
	EventName        string
	SrcObjName       string
}

type Reason uint8

const (
	AUTOCLEARED   Reason = 0
	FAULTDISABLED Reason = 1
	FAULTCLEARED  Reason = 2
)

type FaultRBEntry struct {
	OwnerId          int
	EventId          int
	ResolutionTime   time.Time
	OccuranceTime    time.Time
	SrcObjKey        string
	FaultSeqNumber   uint64
	Description      string
	Resolved         bool
	ResolutionReason Reason
	SrcObjUUID       string
}

type AlarmRBEntry struct {
	OwnerId          int
	EventId          int
	ResolutionTime   time.Time
	OccuranceTime    time.Time
	SrcObjKey        string
	AlarmSeqNumber   uint64
	Description      string
	Resolved         bool
	ResolutionReason Reason
	SrcObjUUID       string
}

type FaultData struct {
	FaultListIdx int
	//AlarmListIdx     int
	CreateAlarmTimer *time.Timer
	FaultSeqNumber   uint64
}

type AlarmData struct {
	AlarmListIdx     int
	AlarmSeqNumber   uint64
	RemoveAlarmTimer *time.Timer
}

type FaultObjKey string
type FaultDataMap map[FaultObjKey]FaultData
type AlarmDataMap map[FaultObjKey]AlarmData
