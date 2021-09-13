//go:build network || nsxt || functional || ALL
// +build network nsxt functional ALL

/*
 * Copyright 2021 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_GetAllNsxtImportableSwitches(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	if vcd.client.Client.APIVCDMaxVersionIs("< 34") {
		check.Skip("At least VCD 10.1 is required")
	}
	skipNoNsxtConfiguration(vcd, check)

	nsxtVdc, err := vcd.org.GetVDCByNameOrId(vcd.config.VCD.Nsxt.Vdc, true)
	check.Assert(err, IsNil)

	allSwitches, err := nsxtVdc.GetAllNsxtImportableSwitches()
	check.Assert(err, IsNil)
	check.Assert(len(allSwitches) > 0, Equals, true)
}

func (vcd *TestVCD) Test_GetNsxtImportableSwitchByName(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	if vcd.client.Client.APIVCDMaxVersionIs("< 34") {
		check.Skip("At least VCD 10.1 is required")
	}
	skipNoNsxtConfiguration(vcd, check)

	nsxtVdc, err := vcd.org.GetVDCByNameOrId(vcd.config.VCD.Nsxt.Vdc, true)
	check.Assert(err, IsNil)

	logicalSwitch, err := nsxtVdc.GetNsxtImportableSwitchByName(vcd.config.VCD.Nsxt.NsxtImportSegment)
	check.Assert(err, IsNil)
	check.Assert(logicalSwitch.NsxtImportableSwitch.Name, Equals, vcd.config.VCD.Nsxt.NsxtImportSegment)
}

func (vcd *TestVCD) Test_GetAllNsxtImportableSwitchesByManager(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	if vcd.client.Client.APIVCDMaxVersionIs("< 34") {
		check.Skip("At least VCD 10.1 is required")
	}
	skipNoNsxtConfiguration(vcd, check)

	nsxtManagerId := testGetNsxtManagerId(vcd, check)

	importableSwitches, err := vcd.client.GetAllNsxtImportableSwitchesByManager(nsxtManagerId)
	check.Assert(err, IsNil)
	check.Assert(len(importableSwitches) > 0, Equals, true)

	var foundSegment bool
	for i := range importableSwitches {
		if importableSwitches[i].NsxtImportableSwitch.Name == vcd.config.VCD.Nsxt.NsxtImportSegment {
			foundSegment = true
		}
	}

	check.Assert(foundSegment, Equals, true)
}

func (vcd *TestVCD) GetNsxtImportableSwitchByName(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	if vcd.client.Client.APIVCDMaxVersionIs("< 34") {
		check.Skip("At least VCD 10.1 is required")
	}
	skipNoNsxtConfiguration(vcd, check)

	nsxtManagerId := testGetNsxtManagerId(vcd, check)

	importableSwitch, err := vcd.client.GetNsxtImportableSwitchByName(nsxtManagerId, vcd.config.VCD.Nsxt.NsxtImportSegment)
	check.Assert(err, IsNil)
	check.Assert(importableSwitch, NotNil)
	check.Assert(importableSwitch.NsxtImportableSwitch.Name, Equals, vcd.config.VCD.Nsxt.NsxtImportSegment)

}

// testGetNsxtManagerId helps to get NSX-T Manager ID. NSX-T Manager endpoint does not return URN, but some VCD
// endpoints require it.
func testGetNsxtManagerId(vcd *TestVCD, check *C) string {
	nsxtManagers, err := vcd.client.QueryNsxtManagerByName(vcd.config.VCD.Nsxt.Manager)
	check.Assert(err, IsNil)

	id := extractUuid(nsxtManagers[0].HREF)
	urn, err := BuildUrnWithUuid("urn:vcloud:nsxtmanager:", id)
	check.Assert(err, IsNil)

	return urn
}
