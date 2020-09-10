/*
 * Copyright 2019 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_CreateExternalNetworkV2NsxT(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	// NSX-T details
	man, err := vcd.client.QueryNsxtManagerByName(vcd.config.Nsxt.Manager)
	check.Assert(err, IsNil)
	manId := "urn:vcloud:nsxtmanager:" + extractUuid(man[0].HREF)

	tier0Router, err := vcd.client.GetNsxtTier0RouterByName(vcd.config.Nsxt.Tier0router, manId)
	check.Assert(err, IsNil)

	neT := testExternalNetworkV2("NSXT_TIER0", tier0Router.NsxtTier0Router.ID, manId)

	r, err := CreateExternalNetworkV2(vcd.client, neT)
	check.Assert(err, IsNil)

	r.ExternalNetwork.Name = "changed_name"
	_, err = r.Update()
	check.Assert(err, IsNil)

	err = r.Delete()
	check.Assert(err, IsNil)
}

func (vcd *TestVCD) Test_CreateExternalNetworkV2PortGroup(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	var err error
	var pgs []*types.PortGroupRecordType

	switch vcd.config.VCD.ExternalNetworkPortGroupType {
	case "DV_PORTGROUP":
		pgs, err = QueryDistributedPortGroup(vcd.client, vcd.config.VCD.ExternalNetworkPortGroup)
	case "NETWORK":
		pgs, err = QueryNetworkPortGroup(vcd.client, vcd.config.VCD.ExternalNetworkPortGroup)
	default:

	}
	check.Assert(err, IsNil)
	check.Assert(len(pgs), Equals, 1)

	// Query
	vcHref, err := getVcenterHref(vcd.client, vcd.config.VCD.VimServer)
	check.Assert(err, IsNil)
	vcuuid := extractUuid(vcHref)

	neT := testExternalNetworkV2(vcd.config.VCD.ExternalNetworkPortGroupType, pgs[0].MoRef, "urn:vcloud:vimserver:"+vcuuid)

	r, err := CreateExternalNetworkV2(vcd.client, neT)
	check.Assert(err, IsNil)

	r.ExternalNetwork.Name = "changed_name"
	_, err = r.Update()
	check.Assert(err, IsNil)

	err = r.Delete()
	check.Assert(err, IsNil)
}

func testExternalNetworkV2(backingType, backingId, NetworkProviderId string) *types.ExternalNetworkV2 {
	neT := &types.ExternalNetworkV2{
		ID:          "",
		Name:        "testNet",
		Description: "",
		Subnets: types.Subnets{[]types.Subnet{
			{
				Gateway:      "1.1.1.1",
				PrefixLength: 24,
				DNSSuffix:    "",
				DNSServer1:   "",
				DNSServer2:   "",
				IPRanges: types.IPRanges2{[]types.IPRange2{
					{
						StartAddress: "1.1.1.3",
						EndAddress:   "1.1.1.50",
					},
				}},
				Enabled:      true,
				UsedIPCount:  0,
				TotalIPCount: 0,
			},
		}},
		NetworkBackings: types.NetworkBackings{[]types.NetworkBacking{
			{
				BackingID: backingId,
				// Name:        tier0Router.NsxtTier0Router.DisplayName,
				BackingType: backingType,
				NetworkProvider: types.NetworkProvider{
					// Name: vcd.config.Nsxt.Manager,
					ID: NetworkProviderId,
				},
			},
		}},
	}

	return neT
}

func getVcenterHref(vcdClient *VCDClient, name string) (string, error) {
	virtualCenters, err := QueryVirtualCenters(vcdClient, fmt.Sprintf("(name==%s)", name))
	if err != nil {
		return "", err
	}
	if len(virtualCenters) == 0 || len(virtualCenters) > 1 {
		return "", fmt.Errorf("vSphere server found %d instances with name '%s' while expected one", len(virtualCenters), name)
	}
	return virtualCenters[0].HREF, nil
}
