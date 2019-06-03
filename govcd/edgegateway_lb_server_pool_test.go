// +build lb lbsp functional ALL

/*
 * Copyright 2019 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (

	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_LBServerPool(check *C) {
	if vcd.config.VCD.EdgeGateway == "" {
		check.Skip("Skipping test because no edge gateway given")
	}
	edge, err := vcd.vdc.FindEdgeGateway(vcd.config.VCD.EdgeGateway)
	check.Assert(err, IsNil)
	check.Assert(edge.EdgeGateway.Name, Equals, vcd.config.VCD.EdgeGateway)

	// Used for creating
	lbPoolConfig := &types.LBPool{
		Name:       Test_LBServerPool,
		Algorithm:       "round-robin",
	}

	createdLbPool, err := edge.CreateLBServerPool(lbPoolConfig)
	check.Assert(err, IsNil)
	check.Assert(createdLbPool.ID, Not(IsNil))

	// // We created monitor successfully therefore let's add it to cleanup list
	parentEntity := vcd.org.Org.Name + "|" + vcd.vdc.Vdc.Name + "|" + vcd.config.VCD.EdgeGateway
	AddToCleanupList(Test_LBServerPool, "lbServerPool", parentEntity, check.TestName())

	// // Lookup by both name and ID and compare that these are equal values
	lbPoolByID, err := edge.ReadLBServerPool(&types.LBPool{ID: createdLbPool.ID})
	check.Assert(err, IsNil)

	lbPoolByName, err := edge.ReadLBServerPool(&types.LBPool{Name: createdLbPool.Name})
	check.Assert(err, IsNil)
	check.Assert(createdLbPool.ID, Equals, lbPoolByName.ID)
	check.Assert(lbPoolByID.ID, Equals, lbPoolByName.ID)
	check.Assert(lbPoolByID.Name, Equals, lbPoolByName.Name)

	check.Assert(createdLbPool.Algorithm, Equals, lbPoolConfig.Algorithm)
	// check.Assert(lbMonitor.Timeout, Equals, lbMonitorByID.Timeout)
	// check.Assert(lbMonitor.MaxRetries, Equals, lbMonitorByID.MaxRetries)

	// Test updating fields
	// Update algorith
	lbPoolByID.Algorithm = "ip-hash"
	updatedLBPool, err := edge.UpdateLBServerPool(lbPoolByID)
	check.Assert(err, IsNil)
	check.Assert(updatedLBPool.Algorithm, Equals, lbPoolByID.Algorithm)
	
	// Try to set invalid algorithm hash and excpect API to return error
	// Invalid algorithm hash. Valid algorithms are: IP-HASH|ROUND-ROBIN|URI|LEASTCONN|URL|HTTP-HEADER.
	lbPoolByID.Algorithm = "invalid_algorithm"
	updatedLBPool, err = edge.UpdateLBServerPool(lbPoolByID)
	check.Assert(err, ErrorMatches, ".*Invalid algorithm.*Valid algorithms are:.*")

	// // Update should fail without name
	lbPoolByID.Name = ""
	_, err = edge.UpdateLBServerPool(lbPoolByID)
	check.Assert(err.Error(), Equals, "load balancer server pool Name cannot be empty")

	// Delete / cleanup
	err = edge.DeleteLBServerPool(&types.LBPool{ID: createdLbPool.ID})
	check.Assert(err, IsNil)
}
