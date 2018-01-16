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

package ipTable

import (
	"errors"
	_ "fmt"
	"strconv"
	"strings"
	"sysd"
	"utils/logging"
)

/*
#cgo CFLAGS: -I../../../../../external/src/github.com/netfilter/libiptables/include -I../../../../../external/src/github.com/netfilter/iptables/include
#cgo LDFLAGS: -L../../../../../external/src/github.com/netfilter/libiptables/lib -lip4tc
#include "ipTable.h"
*/
import "C"

func SysdNewSysdIpTableHandler(logger *logging.Writer) *SysdIpTableHandler {
	ipTableHdl := &SysdIpTableHandler{}
	ipTableHdl.logger = logger
	ipTableHdl.ruleInfo = make(map[string]C.ipt_config_t, 100)
	return ipTableHdl
}

func (hdl *SysdIpTableHandler) AddIpRule(config *sysd.IpTableAcl,
	restart bool) (bool, error) {
	port, err := strconv.Atoi(config.Port)
	var iptEntry C.ipt_config_t
	rv := -1
	if err != nil {
		if strings.Compare(config.Port, ALL_RULE_STR) == 0 {
			hdl.logger.Info("Rule to be applied on all ports")
			port = 0
		}
	}
	ip := config.IpAddr
	pl := 0 // prefix length
	if strings.Contains(config.IpAddr, "/") {
		splitStr := strings.Split(config.IpAddr, "/")
		ip = splitStr[0]
		pl, _ = strconv.Atoi(splitStr[1])
	}
	entry := &C.rule_entry_t{
		Name:         C.CString(config.Name),
		PhysicalPort: C.CString(config.PhysicalPort),
		Action:       C.CString(config.Action),
		IpAddr:       C.CString(ip),
		PrefixLength: C.int(pl),
		Protocol:     C.CString(config.Protocol),
		Port:         C.uint16_t(port),
		Restart:      C.bool(restart),
	}
	switch config.Protocol {
	case "udp":
		rv = int(C.add_iptable_udp_rule(entry, &iptEntry))

	case "tcp":
		rv = int(C.add_iptable_tcp_rule(entry, &iptEntry))

	case "icmp":
		rv = int(C.add_iptable_icmp_rule(entry, &iptEntry))
	default:
		hdl.logger.Err("Rule adding for " + config.Protocol +
			" is not supported")
		return true, nil
	}
	// If rv = -2 or -3 then new entry insert failed....
	// If rv = -1 then duplicated entry (rule)....do not update this into sysd
	if rv <= 0 {
		var errString C.err_t
		C.get_iptc_error_string(&errString, C.int(iptEntry.err_num))
		return false, errors.New(INSERTING_RULE_ERROR +
			C.GoString(&errString.err_string[0]))
	} else {
		hdl.ruleInfo[config.Name] = iptEntry
		return true, nil
	}
}

func (hdl *SysdIpTableHandler) DelIpRule(config *sysd.IpTableAcl) (bool, error) {
	entry, entryFound := hdl.ruleInfo[config.Name]
	if !entryFound {
		hdl.logger.Err("No rule found for " + config.Name +
			" in sysd runtime db.. This means that either the rule is " +
			"not created or it was duplicate rule")
		return true, nil
	}

	rv := int(C.del_iptable_rule(&entry))
	if rv <= 0 {
		hdl.logger.Err("Delete rule failed for " + config.Name)
		var errString C.err_t
		C.get_iptc_error_string(&errString, entry.err_num)
		return false, errors.New(DELETING_RULE_ERROR +
			C.GoString(&errString.err_string[0]))

	}
	delete(hdl.ruleInfo, config.Name)
	return true, nil
}
