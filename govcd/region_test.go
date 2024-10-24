//go:build tm || functional || ALL

/*
 * Copyright 2024 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v3/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_TmRegion(check *C) {
	skipNonTm(vcd, check)
	sysadminOnly(vcd, check)

	// existingRegion, err := vcd.client.GetRegionByName(check.TestName())
	// check.Assert(err, IsNil)

	// err = existingRegion.Delete()
	// check.Assert(err, IsNil)

	// check.Assert(true, Equals, false)

	superVisorName := "vcfcons-mgmt-vc03-wcp"
	nsxtManagerId := "urn:vcloud:nsxtmanager:1c9d361e-e69b-449e-a8d0-476d59806a06"
	storagePolicyName := "vSAN Default Storage Policy"

	supervisor, err := vcd.client.GetSupervisorByName(superVisorName)
	check.Assert(err, IsNil)

	r := &types.Region{
		Name: check.TestName(),
		NsxManager: &types.OpenApiReference{
			ID: nsxtManagerId,
		},
		Supervisors: []types.OpenApiReference{
			{
				ID:   supervisor.Supervisor.SupervisorID,
				Name: supervisor.Supervisor.Name,
			},
		},
		StoragePolicies: []string{storagePolicyName},
		IsEnabled:       true,
	}

	createdRegion, err := vcd.client.CreateRegion(r)
	check.Assert(err, IsNil)
	check.Assert(createdRegion.Region, NotNil)
	fmt.Println(createdRegion.Region.ID)
	AddToCleanupListOpenApi(createdRegion.Region.ID, check.TestName(), types.OpenApiPathVcf+types.OpenApiEndpointRegions+createdRegion.Region.ID)

	check.Assert(createdRegion.Region.Status, Equals, "READY") // Region is operational

	// Get By Name
	byName, err := vcd.client.GetRegionByName(r.Name)
	check.Assert(err, IsNil)
	check.Assert(byName, NotNil)

	// Get By ID
	byId, err := vcd.client.GetRegionById(createdRegion.Region.ID)
	check.Assert(err, IsNil)
	check.Assert(byId, NotNil)

	// Get All
	allRegions, err := vcd.client.GetAllRegions(nil)
	check.Assert(err, IsNil)
	check.Assert(allRegions, NotNil)
	check.Assert(len(allRegions) > 0, Equals, true)

	// TODO: TM: No Update so far
	// Update
	// createdRegion.Region.IsEnabled = false
	// updated, err := createdRegion.Update(createdRegion.Region)
	// check.Assert(err, IsNil)
	// check.Assert(updated, NotNil)

	// Delete
	err = createdRegion.Delete()
	check.Assert(err, IsNil)

	notFoundByName, err := vcd.client.GetRegionByName(createdRegion.Region.Name)
	check.Assert(ContainsNotFound(err), Equals, true)
	check.Assert(notFoundByName, IsNil)

	// region, err := vcd.client.GetRegionByName(vcd.config.Tm.Region)
	// check.Assert(err, IsNil)
	// check.Assert(region, NotNil)

	// // Get by ID
	// regionById, err := vcd.client.GetRegionById(region.Region.ID)
	// check.Assert(err, IsNil)
	// check.Assert(regionById, NotNil)

	// check.Assert(region.Region, DeepEquals, regionById.Region)

	// allTmVdc, err := vcd.client.GetAllRegions(nil)
	// check.Assert(err, IsNil)
	// check.Assert(len(allTmVdc) > 0, Equals, true)
}
