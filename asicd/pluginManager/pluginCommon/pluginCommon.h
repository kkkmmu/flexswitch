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

#ifndef PLUGIN_MGR_H
#define PLUGIN_MGR_H
#define MAC_ADDR_LEN 6
#define DEFAULT_VLAN_ID 4095
#define BOOT_MODE_COLDBOOT 0
#define BOOT_MODE_WARMBOOT 1
#define MAX_VLAN_ID 4096
#define NEIGHBOR_TYPE_COPY_TO_CPU 0x1
#define NEIGHBOR_TYPE_BLACKHOLE 0x2
#define NEIGHBOR_TYPE_FULL_SPEC_NEXTHOP 0x4
#define NEIGHBOR_L2_ACCESS_TYPE_PORT 0x8
#define NEIGHBOR_L2_ACCESS_TYPE_LAG 0x10
#define NEXTHOP_TYPE_COPY_TO_CPU 0x1
#define NEXTHOP_TYPE_BLACKHOLE 0x2
#define NEXTHOP_TYPE_FULL_SPEC_NEXTHOP 0x4
#define NEXTHOP_L2_ACCESS_TYPE_PORT 0x8
#define NEXTHOP_L2_ACCESS_TYPE_LAG 0x10
#define INTF_STATE_DOWN 0
#define INTF_STATE_UP 1
#define INVALID_OBJECT_ID 0xFFFFFFFFFFFFFFFF
#define MAC_ENTRY_LEARNED 0x1
#define MAC_ENTRY_AGED 0x2
#define ROUTE_TYPE_CONNECTED 0x00000001
#define ROUTE_TYPE_SINGLEPATH 0x00000002
#define ROUTE_TYPE_MULTIPATH 0x00000004
#define ROUTE_OPERATION_TYPE_UPDATE 0x00000008
#define ROUTE_TYPE_V6 0x000000010
#define MAX_NEXTHOPS_PER_GROUP 32
#define INTF_TYPE_MASK 0x7f000000
#define INTF_TYPE_SHIFT 24
#define INTF_ID_MASK 0xffffff
#define INTF_ID_SHIFT 0

#define	PORT_PROTOCOL_ARP 		0x1
#define	PORT_PROTOCOL_DHCP 		0x2
#define	PORT_PROTOCOL_DHCP_RELAY 	0x4
#define	PORT_PROTOCOL_BGP 		0x8
#define	PORT_PROTOCOL_OSPF 		0x10
#define	PORT_PROTOCOL_VXLAN 		0x20
#define	PORT_PROTOCOL_MPLS		0x40
#define	PORT_PROTOCOL_BFD 		0x80

//Port attribute update flags
#define PORT_ATTR_PHY_INTF_TYPE 0x00000001
#define PORT_ATTR_ADMIN_STATE   0x00000002
#define PORT_ATTR_MAC_ADDR      0x00000004
#define PORT_ATTR_SPEED         0x00000080
#define PORT_ATTR_DUPLEX        0x00000010
#define PORT_ATTR_AUTONEG       0x00000020
#define PORT_ATTR_MEDIA_TYPE    0x00000040
#define PORT_ATTR_MTU           0x00000080
#define PORT_ATTR_BREAKOUT_MODE 0x00000100
//Port breakout modes
#define PORT_BREAKOUT_MODE_UNSUPPORTED 0
#define PORT_BREAKOUT_MODE_1x40  0x00000001
#define PORT_BREAKOUT_MODE_4x10  0x00000002
#define PORT_BREAKOUT_MODE_1x100 0x00000004
// IP values copied from golang syscall package
#define IP_TYPE_IPV4 0x2
#define IP_TYPE_IPV6 0xa
/* STP STATE definitions */
enum stpPortStates {
    StpPortStateBlocking = 0,
    StpPortStateLearning,
    StpPortStateForwarding,
    StpPortStateCount
};

enum hashTypes {
    HASHTYPE_SRCMAC_DSTMAC = 0,
    HASHTYPE_SRCIP_DSTIP = 6,
    HASHTYPE_END
};
enum mediaType {
    MediaTypeCount
};
enum duplexType {
    HalfDuplex = 0,
    FullDuplex,
    DuplexCount
};
enum portIfType {
    PortIfTypeMII,
    PortIfTypeGMII,
    PortIfTypeSGMII,
    PortIfTypeQSGMII,
    PortIfTypeSFI,
    PortIfTypeXFI,
    PortIfTypeXAUI,
    PortIfTypeXLAUI,
    PortIfTypeRXAUI,
    PortIfTypeCR,
    PortIfTypeCR2,
    PortIfTypeCR4,
    PortIfTypeKR,
    PortIfTypeKR2,
    PortIfTypeKR4,
    PortIfTypeSR,
    PortIfTypeSR2,
    PortIfTypeSR4,
    PortIfTypeSR10,
    PortIfTypeLR,
    PortIfTypeLR4,
    PortIfTypeCount
};
enum portStatTypes {
    IfInOctets = 0,
    IfInUcastPkts,
    IfInDiscards,
    IfInErrors,
    IfInUnknownProtos,
    IfOutOctets,
    IfOutUcastPkts,
    IfOutDiscards,
    IfOutErrors,
    portStatTypesMax
};
#endif //PLUGIN_MGR_H
