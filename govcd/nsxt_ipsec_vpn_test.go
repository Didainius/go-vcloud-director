// +build network nsxt functional openapi ALL

package govcd

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_NsxtIpSecVpn(check *C) {
	skipNoNsxtConfiguration(vcd, check)
	skipOpenApiEndpointTest(vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointFirewallGroups)

	org, err := vcd.client.GetOrgByName(vcd.config.VCD.Org)
	check.Assert(err, IsNil)

	nsxtVdc, err := org.GetVDCByName(vcd.config.VCD.Nsxt.Vdc, false)
	check.Assert(err, IsNil)

	edge, err := nsxtVdc.GetNsxtEdgeGatewayByName(vcd.config.VCD.Nsxt.EdgeGateway)
	check.Assert(err, IsNil)

	ipSecDef := &types.NsxtIpSecVpn{
		Name:        check.TestName(),
		Description: check.TestName() + "-description",
		Enabled:     true,
		LocalEndpoint: types.NsxtIpSecVpnLocalEndpoint{
			LocalAddress:  edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
			LocalNetworks: []string{"10.10.10.0/24"},
		},
		RemoteEndpoint: types.NsxtIpSecVpnRemoteEndpoint{
			RemoteId:       "192.168.140.1",
			RemoteAddress:  "192.168.140.1",
			RemoteNetworks: []string{"20.20.20.0/24"},
		},
		PreSharedKey: "PSK-Sec",
		SecurityType: "DEFAULT",
		Logging:      false,
	}

	createdIpSecVpn, err := edge.CreateIpSecVpn(ipSecDef)
	check.Assert(err, IsNil)
	openApiEndpoint := types.OpenApiPathVersion1_0_0 + fmt.Sprintf(types.OpenApiEndpointIpSecVpn, createdIpSecVpn.edgeGatewayId) + createdIpSecVpn.NsxtIpSecVpn.ID
	AddToCleanupListOpenApi(createdIpSecVpn.NsxtIpSecVpn.Name, check.TestName(), openApiEndpoint)

	foundIpSecVpnById, err := edge.GetIpSecVpnById(createdIpSecVpn.NsxtIpSecVpn.ID)
	check.Assert(err, IsNil)
	check.Assert(foundIpSecVpnById.NsxtIpSecVpn, DeepEquals, createdIpSecVpn.NsxtIpSecVpn)

	foundIpSecVpnByName, err := edge.GetIpSecVpnByName(createdIpSecVpn.NsxtIpSecVpn.Name)
	check.Assert(err, IsNil)
	check.Assert(foundIpSecVpnByName.NsxtIpSecVpn, DeepEquals, createdIpSecVpn.NsxtIpSecVpn)
	check.Assert(foundIpSecVpnByName.NsxtIpSecVpn, DeepEquals, foundIpSecVpnById.NsxtIpSecVpn)

	check.Assert(createdIpSecVpn.NsxtIpSecVpn.ID, Not(Equals), "")

	ipSecDef.Name = check.TestName() + "-updated"
	ipSecDef.RemoteEndpoint.RemoteAddress = "192.168.40.1"
	ipSecDef.ID = createdIpSecVpn.NsxtIpSecVpn.ID

	updatedIpSecVpn, err := createdIpSecVpn.Update(ipSecDef)
	check.Assert(updatedIpSecVpn.NsxtIpSecVpn.Name, Equals, ipSecDef.Name)
	check.Assert(updatedIpSecVpn.NsxtIpSecVpn.ID, Equals, ipSecDef.ID)
	check.Assert(updatedIpSecVpn.NsxtIpSecVpn.RemoteEndpoint.RemoteAddress, Equals, ipSecDef.RemoteEndpoint.RemoteAddress)

	err = createdIpSecVpn.Delete()
	check.Assert(err, IsNil)

	// Ensure rule does not exist in the list
	allVpnConfigs, err := edge.GetAllIpSecVpns(nil)
	check.Assert(err, IsNil)
	for _, vpnConfig := range allVpnConfigs {
		check.Assert(vpnConfig.IsEqualTo(updatedIpSecVpn.NsxtIpSecVpn), Equals, false)
	}

}