// +build network nsxt functional openapi ALL

package govcd

import (
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_NsxtFirewallGroup(check *C) {
	skipNoNsxtConfiguration(vcd, check)
	skipOpenApiEndpointTest(vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointFirewallGroups)

	org, err := vcd.client.GetOrgByName(vcd.config.VCD.Org)
	check.Assert(err, IsNil)

	nsxtVdc, err := org.GetVDCByName(vcd.config.VCD.Nsxt.Vdc, false)
	check.Assert(err, IsNil)

	egw, err := nsxtVdc.GetNsxtEdgeGatewayByName(vcd.config.VCD.Nsxt.EdgeGateway)
	check.Assert(err, IsNil)

	fwGroupDefinition := &types.NsxtFirewallGroup{
		Name:        check.TestName(),
		Description: check.TestName() + "-Description",
		Type:        types.FirewallGroupTypeSecurityGroup,
		OwnerRef:    &types.OpenApiReference{ID: egw.EdgeGateway.ID},
	}

	createdFwGroup, err := org.CreateNsxtFirewallGroup(fwGroupDefinition)
	check.Assert(err, IsNil)
	check.Assert(createdFwGroup.NsxtFirewallGroup.ID, Not(Equals), "")
	check.Assert(createdFwGroup.NsxtFirewallGroup.EdgeGatewayRef.Name, Equals, vcd.config.VCD.Nsxt.EdgeGateway)

	openApiEndpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointFirewallGroups + createdFwGroup.NsxtFirewallGroup.ID
	AddToCleanupListOpenApi(createdFwGroup.NsxtFirewallGroup.Name, check.TestName(), openApiEndpoint)

	fwGroupDefinition.ID = createdFwGroup.NsxtFirewallGroup.ID
	// fwGroupDefinition.Ownerref.Name = createdFwGroup.NsxtFirewallGroup.Ownerref.Name

	// On creation one sets OwnerRef field, but in GET Edge Gateway is returned in EdgeGatewayRef
	// field
	check.Assert(createdFwGroup.NsxtFirewallGroup.EdgeGatewayRef.ID, Equals, fwGroupDefinition.OwnerRef.ID)
	check.Assert(createdFwGroup.NsxtFirewallGroup.Description, Equals, fwGroupDefinition.Description)
	check.Assert(createdFwGroup.NsxtFirewallGroup.Name, Equals, fwGroupDefinition.Name)
	check.Assert(createdFwGroup.NsxtFirewallGroup.Type, Equals, fwGroupDefinition.Type)

	// check.Assert(createdFwGroup.NsxtFirewallGroup, DeepEquals, fwGroupDefinition)

	// Update
	createdFwGroup.NsxtFirewallGroup.Description = "updated-description"
	createdFwGroup.NsxtFirewallGroup.Name = check.TestName() + "-updated"

	updatedFwGroup, err := createdFwGroup.Update(createdFwGroup.NsxtFirewallGroup)
	check.Assert(err, IsNil)
	check.Assert(updatedFwGroup.NsxtFirewallGroup, DeepEquals, createdFwGroup.NsxtFirewallGroup)

	check.Assert(updatedFwGroup, DeepEquals, createdFwGroup)

	// Get all Firewall Groups and check if the created one is there
	allFwGroups, err := org.GetAllNsxtFirewallGroups(nil)
	check.Assert(err, IsNil)
	fwGroupFound := false
	for i := range allFwGroups {
		if allFwGroups[i].NsxtFirewallGroup.ID == updatedFwGroup.NsxtFirewallGroup.ID {
			fwGroupFound = true
			break
		}
	}
	check.Assert(fwGroupFound, Equals, true)

	// Get firewall group by name
	fwGroupByName, err := org.GetNsxtFirewallGroupByName(updatedFwGroup.NsxtFirewallGroup.Name)
	check.Assert(err, IsNil)

	fwGroupById, err := org.GetNsxtFirewallGroupById(updatedFwGroup.NsxtFirewallGroup.ID)
	check.Assert(err, IsNil)

	check.Assert(fwGroupById, DeepEquals, fwGroupByName)

	// Get Firewall Group using Edge Gateway
	egwFirewallGroup, err := egw.GetNsxtFirewallGroupByName(updatedFwGroup.NsxtFirewallGroup.Name)

	check.Assert(egwFirewallGroup.NsxtFirewallGroup, DeepEquals, fwGroupByName.NsxtFirewallGroup)

	// Remove
	err = createdFwGroup.Delete()
	check.Assert(err, IsNil)
}
