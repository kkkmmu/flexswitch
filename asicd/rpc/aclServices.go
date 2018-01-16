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

// This file defines all interfaces provided for the LAG service
package rpc

import (
	"asicdServices"
)

func (svcHdlr AsicDaemonServiceHandler) CreateAcl(aclObj *asicdServices.Acl) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateAclInternal(aclObj *asicdServices.Acl, clientInt string) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteAcl(aclObj *asicdServices.Acl) (bool, error) {
	return true, nil

}

func (svcHdlr AsicDaemonServiceHandler) UpdateAcl(oldAclObj, newAclObj *asicdServices.Acl, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateAclRule(aclRuleObj *asicdServices.AclRule) (bool, error) {
	return true, nil

}

func (svcHdlr AsicDaemonServiceHandler) CreateAclRuleInternal(aclRuleObj *asicdServices.AclRule, clientInt string) (bool, error) {
	return true, nil

}

func (svcHdlr AsicDaemonServiceHandler) DeleteAclRule(aclRuleObj *asicdServices.AclRule) (bool, error) {
	return true, nil

}

func (svcHdlr AsicDaemonServiceHandler) UpdateAclRule(oldAclRuleObj, newAclRuleObj *asicdServices.AclRule, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil

}

func (svcHdlr AsicDaemonServiceHandler) GetAclRuleState(ruleName string) (*asicdServices.AclRuleState, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetAclState(aclName string, name string) (*asicdServices.AclState, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkAclState(currMarker, count asicdServices.Int) (*asicdServices.AclStateGetInfo, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkAclRuleState(currMarker, count asicdServices.Int) (*asicdServices.AclRuleStateGetInfo, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkCoppStatState(currMarker, count asicdServices.Int) (*asicdServices.CoppStatStateGetInfo, error) {
	return nil, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetCoppStatState(proto string) (*asicdServices.CoppStatState, error) {
	return nil, nil
}
