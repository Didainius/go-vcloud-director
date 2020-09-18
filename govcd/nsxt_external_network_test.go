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
	manId, err := buildUrnWithUuid("urn:vcloud:nsxtmanager:", extractUuid(man[0].HREF))
	check.Assert(err, IsNil)

	tier0Router, err := vcd.client.GetImportableNsxtTier0RouterByName(vcd.config.Nsxt.Tier0router, manId)
	check.Assert(err, IsNil)

	neT := testExternalNetworkV2(types.ExternalNetworkBackingTypeNsxtTier0Router, tier0Router.NsxtTier0Router.ID, manId)

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
	case types.ExternalNetworkBackingDvPortgroup:
		pgs, err = QueryDistributedPortGroup(vcd.client, vcd.config.VCD.ExternalNetworkPortGroup)
	case types.ExternalNetworkBackingTypeNetwork:
		pgs, err = QueryNetworkPortGroup(vcd.client, vcd.config.VCD.ExternalNetworkPortGroup)
	default:
		check.Errorf("unrecognized external network portgroup type: %s", vcd.config.VCD.ExternalNetworkPortGroupType)
	}
	check.Assert(err, IsNil)
	check.Assert(len(pgs), Equals, 1)

	// Query
	vcHref, err := getVcenterHref(vcd.client, vcd.config.VCD.VimServer)
	check.Assert(err, IsNil)
	vcuuid := extractUuid(vcHref)

	vcUrn, err := buildUrnWithUuid("urn:vcloud:vimserver:", vcuuid)
	check.Assert(err, IsNil)

	neT := testExternalNetworkV2(vcd.config.VCD.ExternalNetworkPortGroupType, pgs[0].MoRef, vcUrn)

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
		Subnets: types.ExternalNetworkV2Subnets{[]types.ExternalNetworkV2Subnet{
			{
				Gateway:      "1.1.1.1",
				PrefixLength: 24,
				DNSSuffix:    "",
				DNSServer1:   "",
				DNSServer2:   "",
				IPRanges: types.ExternalNetworkV2IPRanges{[]types.ExternalNetworkV2IPRange{
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
		NetworkBackings: types.ExternalNetworkV2Backings{[]types.ExternalNetworkV2Backing{
			{
				BackingID: backingId,
				// Name:        tier0Router.NsxtTier0Router.DisplayName,
				BackingType: backingType,
				NetworkProvider: types.NetworkProviderProvider{
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

// buildUrnWithUuid helps to build valid URNs where APIs require URN format, but other API responds with UUID (or
// extracted from HREF)
func buildUrnWithUuid(urnPrefix, uuid string) (string, error) {
	if !IsUuid(uuid) {
		return "", fmt.Errorf("supplied uuid '%s' is not valid UUID", uuid)
	}

	urn := urnPrefix + uuid
	if !isUrn(urn) {
		return "", fmt.Errorf("failed building valid URN '%s'", urn)
	}

	return urn, nil
}
