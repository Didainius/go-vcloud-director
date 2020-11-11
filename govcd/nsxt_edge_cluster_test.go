// +build network functional openapi ALL

/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_GetAllNsxtEdgeClusters(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	nsxtVdc, err := vcd.org.GetVDCByNameOrId("nsxt-vdc-dainius", true)
	check.Assert(err, IsNil)

	tier0Router, err := vcd.client.GetAllNsxtEdgeClusters(nsxtVdc.Vdc.ID, nil)
	check.Assert(err, IsNil)
	check.Assert(tier0Router, NotNil)
}

func (vcd *TestVCD) Test_GetNsxtEdgeClusterByName(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	nsxtVdc, err := vcd.org.GetVDCByNameOrId("nsxt-vdc-dainius", true)
	check.Assert(err, IsNil)

	tier0Router, err := vcd.client.GetAllNsxtEdgeClusters(nsxtVdc.Vdc.ID, nil)
	check.Assert(err, IsNil)
	check.Assert(tier0Router, NotNil)

	ecl, err := vcd.client.GetNsxtEdgeClusterByName(tier0Router[0].NsxtEdgeCluster.Name, nsxtVdc.Vdc.ID)
	check.Assert(err, IsNil)
	check.Assert(tier0Router, NotNil)
	check.Assert(ecl, DeepEquals, tier0Router[0])

}
