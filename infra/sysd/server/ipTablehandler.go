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

package server

import (
	"models/objects"
	"sysd"
	"utils/dbutils"
)

func (server *SYSDServer) ReadIpAclConfigFromDB(dbHdl *dbutils.DBUtil) error {
	server.logger.Info("Reading Ip Acl Config From Db")
	if dbHdl != nil {
		var dbObj objects.IpTableAcl
		objList, err := dbHdl.GetAllObjFromDb(dbObj)
		if err != nil {
			server.logger.Err("DB query failed for IpTableAcl config")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := sysd.NewIpTableAcl()
			dbObject := objList[idx].(objects.IpTableAcl)
			objects.ConvertsysdIpTableAclObjToThrift(&dbObject, obj)
			server.AddIpTableRule(obj, true /*restart*/)
		}
	}
	server.logger.Info("reading ip acl config done")
	return nil
}

func (server *SYSDServer) AddIpTableRule(ipaclConfig *sysd.IpTableAcl,
	restart bool) (bool, error) {
	return (server.sysdIpTableMgr.AddIpRule(ipaclConfig, restart))
}

func (server *SYSDServer) DelIpTableRule(ipaclConfig *sysd.IpTableAcl) (bool, error) {
	return (server.sysdIpTableMgr.DelIpRule(ipaclConfig))
}
