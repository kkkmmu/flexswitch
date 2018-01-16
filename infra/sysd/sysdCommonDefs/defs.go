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

package sysdCommonDefs

import ()

const (
	PUB_SOCKET_ADDR = "ipc:///tmp/sysd.ipc"
)

const (
	G_LOG       uint8 = 1 // Global logging configuration
	C_LOG       uint8 = 2 // Component level logging configuration
	KA_DAEMON   uint8 = 3 // Daemon keepalive notification
	SYSTEM_Info uint8 = 4 // System Information notification
)

type Notification struct {
	Type    uint8
	Payload []byte
}

//Logging levels
type SRDebugLevel uint8

const (
	OFF    SRDebugLevel = 0
	CRIT   SRDebugLevel = 1
	ERR    SRDebugLevel = 2
	WARN   SRDebugLevel = 3
	ALERT  SRDebugLevel = 4
	EMERG  SRDebugLevel = 5
	NOTICE SRDebugLevel = 6
	INFO   SRDebugLevel = 7
	DEBUG  SRDebugLevel = 8
	TRACE  SRDebugLevel = 9
)

type GlobalLogging struct {
	Level SRDebugLevel
}

type ComponentLogging struct {
	Name  string
	Level SRDebugLevel
}

const (
	SYSD_TOTAL_KA_DAEMONS = 32
)

type SRDaemonStatus uint8

const (
	UP         SRDaemonStatus = 0
	STARTING   SRDaemonStatus = 1
	RESTARTING SRDaemonStatus = 2
	STOPPED    SRDaemonStatus = 3
)

func ConvertDaemonStateCodeToString(status SRDaemonStatus) string {
	var statusStr string
	switch status {
	case UP:
		statusStr = "up"
	case STARTING:
		statusStr = "starting"
	case RESTARTING:
		statusStr = "restarting"
	case STOPPED:
		statusStr = "stopped"
	default:
		statusStr = "unknown"
	}
	return statusStr
}

type DaemonStatus struct {
	Name   string
	Status SRDaemonStatus
}
