/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_QueryNsxtManagerByName(check *C) {
	nsxtManagers, err := QueryNsxtManagerByName(vcd.client, "nsx*")
	check.Assert(err, IsNil)
	check.Assert(len(nsxtManagers), Equals, 1)
}

func (vcd *TestVCD) Test_GetAllNsxtTier0Routers(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	adminOrg, err := vcd.client.GetAdminOrgByName(vcd.config.VCD.Org)
	check.Assert(err, IsNil)
	check.Assert(adminOrg, NotNil)

	nsxtManagers, err := QueryNsxtManagerByName(vcd.client, "nsx*")
	check.Assert(err, IsNil)
	check.Assert(len(nsxtManagers), Equals, 1)

	idddd, err := GetUuidFromHref(nsxtManagers[0].HREF, true)
	check.Assert(err, IsNil)
	fullId := "urn:vcloud:nsxtmanager:" + idddd

	tier0Routers, err := adminOrg.GetAllNsxtTier0Routers(fullId)
	check.Assert(err, IsNil)

	spew.Dump(tier0Routers[0].NsxtTier0Router)

}

// GetAllNsxtTier0Routers
