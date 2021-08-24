// +build alb functional ALL

package govcd

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/vmware/go-vcloud-director/v2/types/v56"

	. "gopkg.in/check.v1"
)

// Tests VDC storage profile update
func (vcd *TestVCD) Test_GetAllAlbImportableServiceEngineGroups(check *C) {
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

	spew.Dump(importableSeGroups[0].NsxtAlbImportableServiceEngineGroups)

	// Cleanup
	err = createdAlbCloud.Delete()
	check.Assert(err, IsNil)
}
