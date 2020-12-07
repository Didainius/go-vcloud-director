// +build network nsxt functional openapi ALL

package govcd

import (
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_NsxtOrgVdcNetworkIsolated(check *C) {
	skipOpenApiEndpointTest(vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointOrgVdcNetworks)
	skipNoNsxtConfiguration(vcd, check)

	orgVdcNetworkConfig := &types.OpenApiOrgVdcNetwork{
		Name:        check.TestName(),
		Description: check.TestName() + "-description",
		OrgVdc:      &types.OpenApiReference{ID: vcd.nsxtVdc.Vdc.ID},

		NetworkType: types.OrgVdcNetworkTypeIsolated,
		Subnets: types.OrgVdcNetworkSubnets{
			Values: []types.OrgVdcNetworkSubnetValues{
				{
					Gateway:      "4.1.1.1",
					PrefixLength: 25,
					DNSServer1:   "8.8.8.8",
					DNSServer2:   "8.8.4.4",
					DNSSuffix:    "bar.foo",
					IPRanges: types.OrgVdcNetworkSubnetIPRanges{
						Values: []types.OrgVdcNetworkSubnetIPRangeValues{
							{
								StartAddress: "4.1.1.20",
								EndAddress:   "4.1.1.30",
							},
							{
								StartAddress: "4.1.1.40",
								EndAddress:   "4.1.1.50",
							},
							{
								StartAddress: "4.1.1.88",
								EndAddress:   "4.1.1.92",
							},
						}},
				},
			},
		},
	}

	runOpenApiOrgVdcNetworkTest(vcd, check, orgVdcNetworkConfig)
}

func (vcd *TestVCD) Test_NsxtOrgVdcNetworkRouted(check *C) {
	skipOpenApiEndpointTest(vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointOrgVdcNetworks)
	skipNoNsxtConfiguration(vcd, check)

	egw, err := vcd.org.GetNsxtEdgeGatewayByName(vcd.config.VCD.Nsxt.EdgeGateway)
	check.Assert(err, IsNil)

	orgVdcNetworkConfig := &types.OpenApiOrgVdcNetwork{
		Name:        check.TestName(),
		Description: check.TestName() + "-description",
		OrgVdc:      &types.OpenApiReference{ID: vcd.nsxtVdc.Vdc.ID},

		NetworkType: types.OrgVdcNetworkTypeRouted,

		// Connection is used for "routed" network
		Connection: &types.Connection{
			RouterRef: types.OpenApiReference{
				ID: egw.EdgeGateway.ID,
			},
			ConnectionType: "INTERNAL",
		},
		Subnets: types.OrgVdcNetworkSubnets{
			Values: []types.OrgVdcNetworkSubnetValues{
				{
					Gateway:      "2.1.1.1",
					PrefixLength: 24,
					DNSServer1:   "8.8.8.8",
					DNSServer2:   "8.8.4.4",
					DNSSuffix:    "foo.bar",
					IPRanges: types.OrgVdcNetworkSubnetIPRanges{
						Values: []types.OrgVdcNetworkSubnetIPRangeValues{
							{
								StartAddress: "2.1.1.20",
								EndAddress:   "2.1.1.30",
							},
							{
								StartAddress: "2.1.1.40",
								EndAddress:   "2.1.1.50",
							},
							{
								StartAddress: "2.1.1.60",
								EndAddress:   "2.1.1.62",
							}, {
								StartAddress: "2.1.1.72",
								EndAddress:   "2.1.1.74",
							}, {
								StartAddress: "2.1.1.84",
								EndAddress:   "2.1.1.85",
							},
						}},
				},
			},
		},
	}

	runOpenApiOrgVdcNetworkTest(vcd, check, orgVdcNetworkConfig)
}

func (vcd *TestVCD) Test_NsxtOrgVdcNetworkImported(check *C) {
	skipOpenApiEndpointTest(vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointOrgVdcNetworks)
	skipNoNsxtConfiguration(vcd, check)

	logicalSwitch, err := vcd.nsxtVdc.GetNsxtLogicalSwitchByName(vcd.config.VCD.Nsxt.LogicalSwitch)
	check.Assert(err, IsNil)

	orgVdcNetworkConfig := &types.OpenApiOrgVdcNetwork{
		Name:        check.TestName(),
		Description: check.TestName() + "-description",
		OrgVdc:      &types.OpenApiReference{ID: vcd.nsxtVdc.Vdc.ID},

		NetworkType: types.OrgVdcNetworkTypeOpaque,
		// BackingNetworkId contains NSX-T logical switch ID for Imported networks
		BackingNetworkId: logicalSwitch.NsxtLogicalSwitch.ID,

		Subnets: types.OrgVdcNetworkSubnets{
			Values: []types.OrgVdcNetworkSubnetValues{
				{
					Gateway:      "3.1.1.1",
					PrefixLength: 24,
					DNSServer1:   "8.8.8.8",
					DNSServer2:   "8.8.4.4",
					DNSSuffix:    "foo.bar",
					IPRanges: types.OrgVdcNetworkSubnetIPRanges{
						Values: []types.OrgVdcNetworkSubnetIPRangeValues{
							{
								StartAddress: "3.1.1.20",
								EndAddress:   "3.1.1.30",
							},
							{
								StartAddress: "3.1.1.40",
								EndAddress:   "3.1.1.50",
							},
						}},
				},
			},
		},
	}

	runOpenApiOrgVdcNetworkTest(vcd, check, orgVdcNetworkConfig)

}

func runOpenApiOrgVdcNetworkTest(vcd *TestVCD, check *C, orgVdcNetworkConfig *types.OpenApiOrgVdcNetwork) {
	orgVdcNet, err := vcd.vdc.CreateNsxtOrgVdcNetwork(orgVdcNetworkConfig)
	check.Assert(err, IsNil)

	// Use generic "OpenApiEntity" resource cleanup type
	openApiEndpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointOrgVdcNetworks + orgVdcNet.OrgVdcNetwork.ID
	AddToCleanupList(orgVdcNet.OrgVdcNetwork.Name, "OpenApiEntity", "", check.TestName(), openApiEndpoint)

	// Check it can be found
	orgVdcNetById, err := vcd.nsxtVdc.GetNsxtOrgVdcNetworkById(orgVdcNet.OrgVdcNetwork.ID)
	check.Assert(err, IsNil)
	orgVdcNetByName, err := vcd.nsxtVdc.GetNsxtOrgVdcNetworkByName(orgVdcNet.OrgVdcNetwork.Name)
	check.Assert(err, IsNil)

	check.Assert(orgVdcNetById.OrgVdcNetwork.ID, Equals, orgVdcNet.OrgVdcNetwork.ID)
	check.Assert(orgVdcNetByName.OrgVdcNetwork.ID, Equals, orgVdcNet.OrgVdcNetwork.ID)

	// Retrieve all networks in VDC and expect newly created network to be there
	var foundNetInVdc bool
	allOrgVdcNets, err := vcd.nsxtVdc.GetAllNsxtOrgVdcNetworks(nil)
	check.Assert(err, IsNil)
	for _, net := range allOrgVdcNets {
		if net.OrgVdcNetwork.ID == orgVdcNet.OrgVdcNetwork.ID {
			foundNetInVdc = true
		}
	}
	check.Assert(foundNetInVdc, Equals, true)

	// Update
	orgVdcNet.OrgVdcNetwork.Description = check.TestName() + "updated description"
	updatedOrgVdcNet, err := orgVdcNet.Update(orgVdcNet.OrgVdcNetwork)
	check.Assert(err, IsNil)

	check.Assert(updatedOrgVdcNet.OrgVdcNetwork.Name, Equals, orgVdcNet.OrgVdcNetwork.Name)
	check.Assert(updatedOrgVdcNet.OrgVdcNetwork.ID, Equals, orgVdcNet.OrgVdcNetwork.ID)
	check.Assert(updatedOrgVdcNet.OrgVdcNetwork.Description, Equals, orgVdcNet.OrgVdcNetwork.Description)

	// Delete
	err = orgVdcNet.Delete()
	check.Assert(err, IsNil)

	// Test again if it was deleted and expect it to contain ErrorEntityNotFound
	_, err = vcd.nsxtVdc.GetNsxtOrgVdcNetworkByName(orgVdcNet.OrgVdcNetwork.Name)
	check.Assert(ContainsNotFound(err), Equals, true)

	_, err = vcd.nsxtVdc.GetNsxtOrgVdcNetworkById(orgVdcNet.OrgVdcNetwork.ID)
	check.Assert(ContainsNotFound(err), Equals, true)
}
