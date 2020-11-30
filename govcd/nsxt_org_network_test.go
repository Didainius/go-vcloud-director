// +build network nsxt functional openapi ALL

package govcd

import (
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_NsxtOrgVdcNetwork(check *C) {
	skipOpenApiEndpointTest(vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointOrgVdcNetworks)
	skipNoNsxtConfiguration(vcd, check)

	orgVdcNetwork := &types.OpenApiOrgVdcNetwork{
		Name:        check.TestName() + "-routed",
		Description: check.TestName() + "-description",
		OwnerRef:    &types.OpenApiReference{ID: vcd.vdc.Vdc.ID},
		NetworkType: types.OrgVdcNetworkTypeRouted,
		Connection:  types.Connection{},
		// BackingNetworkId:   "",
		// BackingNetworkType: "",
		// GuestVlanTaggingAllowed: false,
		Subnets: types.OrgVdcNetworkSubnets{
			Values: []types.OrgVdcNetworkSubnetValues{
				{
					Gateway:      "1.1.1.1",
					PrefixLength: 24,
					DNSServer1:   "8.8.8.8",
					DNSServer2:   "8.8.4.4",
					DNSSuffix:    "foo.bar",
					IPRanges: types.OrgVdcNetworkSubnetIPRanges{
						Values: []types.ExternalNetworkV2IPRange{
							{
								StartAddress: "1.1.1.20",
								EndAddress:   "1.1.1.30",
							},
							{
								StartAddress: "1.1.1.40",
								EndAddress:   "1.1.1.50",
							},
						}},
				},
			},
		},
		// RouteAdvertised:         nil,
		// SecurityGroups:
	}

	orgVdcNetwork = nil

}
