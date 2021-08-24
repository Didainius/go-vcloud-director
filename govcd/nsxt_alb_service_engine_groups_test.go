// +build alb functional ALL

package govcd

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

// Tests VDC storage profile update
func (vcd *TestVCD) Test_GetAllAlbServiceEngineGroups(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}
	controllers, err := vcd.client.GetAllAlbControllers(nil)
	check.Assert(err, IsNil)
	check.Assert(len(controllers), Equals, 1)

	importableClouds, err := controllers[0].GetAllAlbImportableClouds(nil)
	check.Assert(err, IsNil)
	check.Assert(len(importableClouds) > 0, Equals, true)

	albCloudConfig := &types.NsxtAlbCloud{
		Name:        "test-1",
		Description: "",
		LoadBalancerCloudBacking: types.NsxtAlbCloudBacking{
			BackingId: importableClouds[0].NsxtAlbImportableCloud.ID,
			LoadBalancerControllerRef: types.OpenApiReference{
				ID: controllers[0].NsxtAlbController.ID,
			},
		},
		NetworkPoolRef: types.OpenApiReference{
			ID: importableClouds[0].NsxtAlbImportableCloud.NetworkPoolRef.ID,
		},
	}

	createdAlbCloud, err := vcd.client.CreateAlbCloud(albCloudConfig)
	check.Assert(err, IsNil)

	importableSeGroups, err := vcd.client.GetAllAlbImportableServiceEngineGroups(createdAlbCloud.NsxtAlbCloud.ID, nil)
	check.Assert(err, IsNil)
	check.Assert(len(importableSeGroups) > 0, Equals, true)

	// Create SE group

	albSeGroup := &types.NsxtAlbServiceEngineGroup{
		//Status:                     "",
		//ID:                         "",
		Name:            "one-se",
		Description:     "one-descr",
		ReservationType: "DEDICATED",
		ServiceEngineGroupBacking: types.ServiceEngineGroupBacking{
			BackingId: importableSeGroups[0].NsxtAlbImportableServiceEngineGroups.ID,
			//BackingType:               "",
			LoadBalancerCloudRef: types.OpenApiReference{
				//Name: "",
				ID: createdAlbCloud.NsxtAlbCloud.ID,
			},
		},
		//	HaMode:                     "",
		//	ReservationType:            "",
		//	MaxVirtualServices:         0,
		//	NumDeployedVirtualServices: 0,
		//	ReservedVirtualServices:    0,
		//	OverAllocated:              false,
	}

	//{
	//	"description": "",
	//	"name": "asdasd",
	//	"reservationType": "DEDICATED",
	//	"serviceEngineGroupBacking": {
	//	"backingId": "serviceenginegroup-b494f4df-1aee-464b-bfb8-e4709decb2ed",
	//		"loadBalancerCloudRef": {
	//		"id": "urn:vcloud:loadBalancerCloud:575b56cb-abba-4128-8d1d-7261212a0ffb",
	//			"name": "test-1"
	//	}
	//}
	//}
	//
	createdSeGroup, err := vcd.client.CreateNsxtAlbServiceEngineGroup(albSeGroup)
	check.Assert(err, IsNil)

	err = createdSeGroup.Delete()
	check.Assert(err, IsNil)

	//

	// Cleanup
	err = createdAlbCloud.Delete()
	check.Assert(err, IsNil)
}
