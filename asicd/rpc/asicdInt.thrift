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
// _______   __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __  
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  | 
// |  |__   |  |     |  |__   \  V  /     |   (----  \   \/    \/   /  |  |  ---|  |---- |  ,---- |  |__|  | 
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   | 
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  | 
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__| 
//                                                                                                           

namespace go asicdInt
typedef i32 int
typedef i16 uint16

// Struct for Configuring Reserved Mac Addr on Chip
// so that packets are punted to CPU when enabled
struct RsvdProtocolMacConfig {
        1: string MacAddr
        2: string MacAddrMask
        3: i32 VlanId
}
struct Lag {
    1:i32 LagIfIndex
    2:i32 HashType
    3:list<i32> IfIndexList
}
struct LagGetInfo {
    1: int StartIdx
    2: int EndIdx
    3: int Count
    4: bool More
    5: list<Lag> LagList
}

struct Vtep {
	1 : i32 IfIndex
	2 : string IfName
	3 : i32 Vni
	4 : i32 SrcIfIndex
	5 : string SrcIfName
	6 : string SrcMac
	7 : string DstIp
	8 : string SrcIp
	9 : i16 VlanId
	10 : i16 UDP
	11 : i16 TTL
	12 : i32 NextHopIfIndex
	13 : string NextHopIfName
	14 : string NextHopIp 
	15 : bool Learning
}
struct Vxlan {
	1 : i32 Vni
	2 : string McDestIp
	3 : i16 VlanId
	4 : i32 Mtu
}

struct IPv4NextHop {
    1: string NextHopIp
    2: i32 Weight
    3: i32 NextHopIfType
}
struct IPv6NextHop {
    1: string NextHopIp
    2: i32 Weight
    3: i32 NextHopIfType
}
struct IPv4Route {
    1: string destinationNw
    2: string networkMask
    3: list<IPv4NextHop> NextHopList
}
struct IPv6Route {
    1: string destinationNw
    2: string networkMask
    3: list<IPv6NextHop> NextHopList
}
struct Vlan {
	1 : i32 VlanId
	2 : list<i32> IfIndexList
	3 : list<i32> UntagIfIndexList
}
struct VlanGetInfo {
	1: int StartIdx
	2: int EndIdx
	3: int Count
	4: bool More
	5: list<Vlan> VlanList
}
struct Intf {
    1: string IfName
    2: i32 IfIndex
}
struct IntfGetInfo {
	1: int StartIdx
	2: int EndIdx
	3: int Count
	4: bool More
	5: list<Intf> IntfList
}
service ASICDINTServices {
    // All services listed here are utilities provided to other daemons. These are hand-written and not auto genereated.
    //Vlan
    VlanGetInfo GetBulkVlan(1: int fromIndex, 2: int count);

    //Intf
    IntfGetInfo GetBulkIntf(1: int fromIndex, 2: int count);

    //STP
    i32 CreateStg(1:list<i32> vlanList);
    bool DeleteStg(1:i32 stgId);
    bool SetPortStpState(1:i32 stgId, 2:i32 port, 3:i32 stpState);
    i32 GetPortStpState(1:i32 stgId, 2:i32 port);
    bool UpdateStgVlanList(1:i32 stgId, 2:list<i32> vlanList);
    oneway void FlushFdbStgGroup(1:i32 stgId, 2:i32 port);

    //LAG
    i32 CreateLag(1:string ifName, 2:i32 hashType, 3:string ifIndexList);
    i32 DeleteLag(1:i32 lagId);
    i32 UpdateLag(1:i32 lagId, 2:i32 hashType, 3:string ifIndexList);
    LagGetInfo GetBulkLag(1:int fromIndex, 2:int count);

    //IPv4 neighbors
    i32 CreateIPv4Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);
    i32 UpdateIPv4Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);
    i32 DeleteIPv4Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);

    //IPv4 routes
    oneway void OnewayCreateIPv4Route(1:list<IPv4Route> ipv4RouteList);
    oneway void OnewayDeleteIPv4Route(1:list<IPv4Route> ipv4RouteList);

    //IPv6 neighbors
    i32 CreateIPv6Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);
    i32 UpdateIPv6Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);
    i32 DeleteIPv6Neighbor(1:string ipAddr, 2:string macAddr, 3:i32 vlanId, 4:i32 ifIndex);
    
    //IPv6 routes
    oneway void OnewayCreateIPv6Route(1:list<IPv6Route> ipv6RouteList);
    oneway void OnewayDeleteIPv6Route(1:list<IPv6Route> ipv6RouteList);

    //Protocol Mac Addr
    bool EnablePacketReception(1:RsvdProtocolMacConfig config);
    bool DisablePacketReception(1:RsvdProtocolMacConfig config);
	
    //Err-disable	
    bool ErrorDisablePort(1: i32 ifIndex, 2:string AdminState, 3:string ErrDisableReason)
	
    i32 CreateVxlanVtep(1: Vtep config);
    bool DeleteVxlanVtep(1: Vtep config);
    i32 CreateVxlan(1: Vxlan config);
    bool DeleteVxlan(1: Vxlan config);
    bool LearnFdbVtep(1:string mac, 2:string vtep, 3:i32 ifindex);
}
