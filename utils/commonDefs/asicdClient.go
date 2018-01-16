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

package commonDefs

import (
	"utils/logging"
)

type IPv4IntfState struct {
	IntfRef           string
	IfIndex           int32
	IpAddr            string
	OperState         string
	NumUpEvents       int32
	LastUpEventTime   string
	NumDownEvents     int32
	LastDownEventTime string
	L2IntfType        string
	L2IntfId          int32
}

type IPv4IntfStateGetInfo struct {
	StartIdx          int32
	EndIdx            int32
	Count             int32
	More              bool
	IPv4IntfStateList []IPv4IntfState
}

type Port struct {
	IntfRef     string
	IfIndex     int32
	Description string
	PhyIntfType string
	AdminState  string
	MacAddr     string
	Speed       int32
	Duplex      string
	Autoneg     string
	MediaType   string
	Mtu         int32
}

type PortGetInfo struct {
	StartIdx int32
	EndIdx   int32
	Count    int32
	More     bool
	PortList []Port
}

type PortState struct {
	IntfRef           string
	IfIndex           int32
	Name              string
	OperState         string
	NumUpEvents       int32
	LastUpEventTime   string
	NumDownEvents     int32
	LastDownEventTime string
	Pvid              int32
	IfInOctets        int64
	IfInUcastPkts     int64
	IfInDiscards      int64
	IfInErrors        int64
	IfInUnknownProtos int64
	IfOutOctets       int64
	IfOutUcastPkts    int64
	IfOutDiscards     int64
	IfOutErrors       int64
	ErrDisableReason  string
}

type PortStateGetInfo struct {
	StartIdx      int32
	EndIdx        int32
	Count         int32
	More          bool
	PortStateList []PortState
}

type Vlan struct {
	VlanId           int32
	IfIndexList      []int32
	UntagIfIndexList []int32
}

type VlanGetInfo struct {
	StartIdx int32
	EndIdx   int32
	Count    int32
	More     bool
	VlanList []Vlan
}

type VlanState struct {
	VlanId    int32
	VlanName  string
	OperState string
	IfIndex   int32
}

type VlanStateGetInfo struct {
	StartIdx      int32
	EndIdx        int32
	Count         int32
	More          bool
	VlanStateList []VlanState
}

const (
	//Notification msgs
	NOTIFY_L2INTF_STATE_CHANGE       = iota // 0
	NOTIFY_IPV4_L3INTF_STATE_CHANGE         // 1
	NOTIFY_IPV6_L3INTF_STATE_CHANGE         // 2
	NOTIFY_VLAN_CREATE                      // 3
	NOTIFY_VLAN_DELETE                      // 4
	NOTIFY_VLAN_UPDATE                      // 5
	NOTIFY_LOGICAL_INTF_CREATE              // 6
	NOTIFY_LOGICAL_INTF_DELETE              // 7
	NOTIFY_LOGICAL_INTF_UPDATE              // 8
	NOTIFY_IPV4INTF_CREATE                  // 9
	NOTIFY_IPV4INTF_DELETE                  // 10
	NOTIFY_IPV6INTF_CREATE                  // 11
	NOTIFY_IPV6INTF_DELETE                  // 12
	NOTIFY_LAG_CREATE                       // 13
	NOTIFY_LAG_DELETE                       // 14
	NOTIFY_LAG_UPDATE                       // 15
	NOTIFY_IPV4NBR_MAC_MOVE                 // 16
	NOTIFY_IPV6NBR_MAC_MOVE                 // 17
	NOTIFY_IPV4_ROUTE_CREATE_FAILURE        // 17
	NOTIFY_IPV4_ROUTE_DELETE_FAILURE        // 18
	NOTIFY_IPV6_ROUTE_CREATE_FAILURE        // 19
	NOTIFY_IPV6_ROUTE_DELETE_FAILURE        // 20
	NOTIFY_VTEP_CREATE                      // 21
	NOTIFY_VTEP_DELETE                      // 22
	NOTIFY_MPLSINTF_STATE_CHANGE            // 23
	NOTIFY_MPLSINTF_CREATE                  // 24
	NOTIFY_MPLSINTF_DELETE                  // 25
	NOTIFY_PORT_CONFIG_MODE_CHANGE          // 26
	NOTIFY_PORT_CONFIG_MTU_CHANGE           // 27
)

type AsicdNotification map[uint8]bool

type L2IntfStateNotifyMsg struct {
	MsgType uint8
	IfIndex int32
	IfState uint8
}

type IPv4L3IntfStateNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	IfState uint8
}

type IPv6L3IntfStateNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	IfState uint8
}

type VlanNotifyMsg struct {
	MsgType     uint8
	VlanId      uint16
	VlanIfIndex int32
	VlanName    string
	TagPorts    []int32
	UntagPorts  []int32
}

type LogicalIntfNotifyMsg struct {
	MsgType         uint8
	IfIndex         int32
	LogicalIntfName string
}

type LagNotifyMsg struct {
	MsgType     uint8
	LagName     string
	IfIndex     int32
	IfIndexList []int32
}

type IPv4IntfNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	IntfRef string
}

type IPv4NbrMacMoveNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	VlanId  int32
}

type IPv6NbrMacMoveNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	VlanId  int32
}

type IPv6IntfNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	IntfRef string
}

type PortConfigModeChgNotifyMsg struct {
	IfIndex int32
	OldMode string
	NewMode string
}

type PortConfigMtuChgNotifyMsg struct {
	IfIndex int32
	Mtu     int32
}

type AsicdNotificationHdl interface {
	ProcessNotification(msg AsicdNotifyMsg)
}

// Empty Interface
type AsicdNotifyMsg interface {
}

type AsicdClientStruct struct {
	Logger *logging.Writer
	NHdl   AsicdNotificationHdl
	NMap   AsicdNotification
}

type IPv6IntfState struct {
	IntfRef           string
	IfIndex           int32
	IpAddr            string
	OperState         string
	NumUpEvents       int32
	LastUpEventTime   string
	NumDownEvents     int32
	LastDownEventTime string
	L2IntfType        string
	L2IntfId          int32
}
