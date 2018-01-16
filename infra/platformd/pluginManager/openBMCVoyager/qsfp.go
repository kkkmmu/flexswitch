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

package openBMCVoyager

import (
	"errors"
	"fmt"
	"infra/platformd/pluginManager/pluginCommon"
)

/*
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include "openBMCVoyager.h"
#include "pluginCommon.h"
*/
import "C"

const (
	MAX_NUM_OF_QSFP int = 12
)

func (driver *openBMCVoyagerDriver) GetQsfpState(Id int32) (retObj pluginCommon.QsfpState, err error) {
	var qsfpInfo C.qsfp_info_t

	retval := int(C.GetQsfpState(&qsfpInfo, C.int(Id)))
	if retval < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch qsft state of", Id))
	}
	retObj.VendorName = C.GoString(&qsfpInfo.VendorName[0])
	retObj.VendorOUI = C.GoString(&qsfpInfo.VendorOUI[0])
	retObj.VendorPartNumber = C.GoString(&qsfpInfo.VendorPN[0])
	retObj.VendorRevision = C.GoString(&qsfpInfo.VendorRev[0])
	retObj.VendorSerialNumber = C.GoString(&qsfpInfo.VendorSN[0])
	retObj.DataCode = C.GoString(&qsfpInfo.DataCode[0])
	retObj.Temperature = float64(qsfpInfo.Temperature)
	retObj.Voltage = float64(qsfpInfo.SupplyVoltage)
	for idx := 0; idx < int(pluginCommon.QsfpNumChannel); idx++ {
		retObj.RXPower[idx] = float64(qsfpInfo.RXPower[idx])
		retObj.TXPower[idx] = float64(qsfpInfo.TXPower[idx])
		retObj.TXBias[idx] = float64(qsfpInfo.TXBias[idx])
	}

	retObj.CurrBER = float64(qsfpInfo.CurrBER)
	retObj.AccBER = float64(qsfpInfo.AccBER)
	retObj.MinBER = float64(qsfpInfo.MinBER)
	retObj.MaxBER = float64(qsfpInfo.MaxBER)
	retObj.UDF0 = float64(qsfpInfo.UDF0)
	retObj.UDF1 = float64(qsfpInfo.UDF1)
	retObj.UDF2 = float64(qsfpInfo.UDF2)
	retObj.UDF3 = float64(qsfpInfo.UDF3)
	return retObj, nil
}

func (driver *openBMCVoyagerDriver) GetMaxNumOfQsfp() int {
	driver.logger.Info("Inside OpenBMC Voyager: GetMaxNumOfQsfps()")
	return MAX_NUM_OF_QSFP
}

func (driver *openBMCVoyagerDriver) GetQsfpPMData(Id int32) (retObj pluginCommon.QsfpPMData, err error) {
	var qsfpPMInfo C.qsfp_pm_info_t

	retval := int(C.GetQsfpPMData(&qsfpPMInfo, C.int(Id)))
	if retval < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch qsft pm data of", Id))
	}
	retObj.Temperature = float64(qsfpPMInfo.Temperature)
	retObj.Voltage = float64(qsfpPMInfo.SupplyVoltage)
	for idx := 0; idx < int(pluginCommon.QsfpNumChannel); idx++ {
		retObj.RXPower[idx] = float64(qsfpPMInfo.RXPower[idx])
		retObj.TXPower[idx] = float64(qsfpPMInfo.TXPower[idx])
		retObj.TXBias[idx] = float64(qsfpPMInfo.TXBias[idx])
	}
	return retObj, nil
}
