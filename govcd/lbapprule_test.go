// +build lb lbAppRule functional ALL

/*
 * Copyright 2019 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_LBAppRule(check *C) {
	if vcd.config.VCD.EdgeGateway == "" {
		check.Skip("Skipping test because no edge gateway given")
	}
	edge, err := vcd.vdc.FindEdgeGateway(vcd.config.VCD.EdgeGateway)
	check.Assert(err, IsNil)
	check.Assert(edge.EdgeGateway.Name, Equals, vcd.config.VCD.EdgeGateway)

	// Used for creating
	lbAppRuleConfig := &types.LBAppRule{
		Name:   TestLBAppRule,
		Script: "acl vmware_page url_beg / vmware redirect location https://www.vmware.com/ ifvmware_page",
	}

	createdLbAppRule, err := edge.CreateLBAppRule(lbAppRuleConfig)
	check.Assert(err, IsNil)
	check.Assert(createdLbAppRule.ID, Not(IsNil))

	// // We created application rule successfully therefore let's add it to cleanup list
	parentEntity := vcd.org.Org.Name + "|" + vcd.vdc.Vdc.Name + "|" + vcd.config.VCD.EdgeGateway
	AddToCleanupList(TestLBAppRule, "lbAppRule", parentEntity, check.TestName())

	// // Lookup by both name and ID and compare that these are equal values
	lbAppRuleByID, err := edge.ReadLBAppRule(&types.LBAppRule{ID: createdLbAppRule.ID})
	check.Assert(err, IsNil)

	lbPoolByName, err := edge.ReadLBAppRule(&types.LBAppRule{Name: createdLbAppRule.Name})
	check.Assert(err, IsNil)
	check.Assert(createdLbAppRule.ID, Equals, lbPoolByName.ID)
	check.Assert(lbAppRuleByID.ID, Equals, lbPoolByName.ID)
	check.Assert(lbAppRuleByID.Name, Equals, lbPoolByName.Name)

	check.Assert(createdLbAppRule.Script, Equals, lbAppRuleConfig.Script)

	// Test updating fields
	// Update script to be multi-line
	lbAppRuleByID.Script = "acl other_page url_beg / other redirect location https://www.other.com/ ifother_page\n" +
		"acl other_page2 url_beg / other2 redirect location https://www.other2.com/ ifother_page2"
	updatedAppProfile, err := edge.UpdateLBAppRule(lbAppRuleByID)
	check.Assert(err, IsNil)
	check.Assert(updatedAppProfile.Script, Equals, lbAppRuleByID.Script)

	// Verify that updated pool and its configuration are identical
	check.Assert(updatedAppProfile, DeepEquals, lbAppRuleByID)

	// Try to set invalid script expect API to return error
	// invalid applicationRule script, invalid script line : invalid_script, error details :
	// Unknown keyword 'invalid_script'
	lbAppRuleByID.Script = "invalid_script"
	updatedAppProfile, err = edge.UpdateLBAppRule(lbAppRuleByID)
	check.Assert(err, ErrorMatches, ".*invalid applicationRule script.*")

	// Update should fail without name
	lbAppRuleByID.Name = ""
	_, err = edge.UpdateLBAppRule(lbAppRuleByID)
	check.Assert(err.Error(), Equals, "load balancer application rule Name cannot be empty")

	// Delete / cleanup
	err = edge.DeleteLBAppRule(&types.LBAppRule{ID: createdLbAppRule.ID})
	check.Assert(err, IsNil)
}