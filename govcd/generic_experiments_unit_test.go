//go:build unit || ALL

/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

func Test_generics(t *testing.T) {

	vcdClient := &VCDClient{
		Client: Client{
			APIVersion: "2000",
		},
	}

	internalField := types.NsxtAlbController{ID: "testing-id"}

	asd := genericConstructor[NsxtAlbController, types.NsxtAlbController](internalField, vcdClient, nil)

	spew.Dump(asd)

}

// func Test_gebn2(t *testing.T) {
// 	gg := genConstructorExp2[NsxtAlbController, types.NsxtAlbController]()
// }

func Test_gebn3(t *testing.T) {
	internalField := &types.NsxtAlbController{ID: "testing-id"}
	vcdClient := &VCDClient{
		Client: Client{
			APIVersion: "2000",
		},
	}

	client := &Client{APIVersion: "3000"}

	aaaa := &NsxtAlbControllerExp2{}

	initializedType := genericNew[NsxtAlbControllerExp2, *types.NsxtAlbController, *VCDClient, *Client](aaaa, internalField, vcdClient, client)

	spew.Dump(initializedType)
}

func Test_gebn33(t *testing.T) {
	internalField := &types.NsxtAlbController{ID: "testing-id"}
	vcdClient := &VCDClient{
		Client: Client{
			APIVersion: "2000",
		},
	}

	client := &Client{APIVersion: "3000"}

	// aaaa := &NsxtAlbControllerExp2{}

	initializedType := genericNew22[NsxtAlbControllerExp2, *types.NsxtAlbController, *VCDClient, *Client](internalField, vcdClient, client)

	spew.Dump(initializedType)
}

func Test_gebn33wrappedResponses(t *testing.T) {

	// typeResponses := make([]*types.NsxtAlbController, 0)

	wrappedResp := wrappedResponse[NsxtAlbControllerExp3, types.NsxtAlbController](&types.NsxtAlbController{ID: "single"})
	spew.Dump(wrappedResp)

	typeResponses := make([]*types.NsxtAlbController, 2)
	typeResponses[0] = &types.NsxtAlbController{ID: "one"}
	typeResponses[1] = &types.NsxtAlbController{ID: "two"}

	allwrap := wrappedResponses[NsxtAlbControllerExp3, types.NsxtAlbController](typeResponses)
	spew.Dump(allwrap)

	// allwrap2 := wrappedResponses[NsxtAlbController, types.NsxtAlbController](typeResponses)

}
