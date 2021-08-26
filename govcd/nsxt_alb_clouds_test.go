// +build alb functional ALL

package govcd

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"

	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_AlbClouds(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}
	skipNoNsxtAlbConfiguration(vcd, check)

	albController, err := vcd.client.GetAlbControllerByUrl(vcd.config.VCD.Nsxt.NsxtAlbControllerUrl)
	check.Assert(err, IsNil)
	check.Assert(albController, NotNil)

	importableClouds, err := albController.GetAllAlbImportableClouds(nil)
	check.Assert(err, IsNil)
	check.Assert(len(importableClouds) > 0, Equals, true)

	albCloudConfig := &types.NsxtAlbCloud{
		Name:        "test-1",
		Description: "alb-cloud-description",
		LoadBalancerCloudBacking: types.NsxtAlbCloudBacking{
			BackingId: importableClouds[0].NsxtAlbImportableCloud.ID,
			LoadBalancerControllerRef: types.OpenApiReference{
				ID: albController.NsxtAlbController.ID,
			},
		},
		NetworkPoolRef: types.OpenApiReference{
			ID: importableClouds[0].NsxtAlbImportableCloud.NetworkPoolRef.ID,
		},
	}

	createdAlbCloud, err := vcd.client.CreateAlbCloud(albCloudConfig)
	check.Assert(err, IsNil)
	openApiEndpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbCloud + createdAlbCloud.NsxtAlbCloud.ID
	AddToCleanupListOpenApi(createdAlbCloud.NsxtAlbCloud.Name, check.TestName(), openApiEndpoint)

	// Get all clouds and ensure the needed on is found
	allClouds, err := vcd.client.GetAllAlbClouds(nil)
	check.Assert(err, IsNil)
	var foundCreatedCloud bool
	for cloudIndex := range allClouds {
		if allClouds[cloudIndex].NsxtAlbCloud.ID == createdAlbCloud.NsxtAlbCloud.ID {
			foundCreatedCloud = true
			break
		}
	}
	check.Assert(foundCreatedCloud, Equals, true)

	// Filter lookup by name
	filter := url.Values{}
	filter.Add("filter", "name=="+createdAlbCloud.NsxtAlbCloud.Name)
	allCloudsFiltered, err := vcd.client.GetAllAlbClouds(filter)
	check.Assert(err, IsNil)
	check.Assert(len(allCloudsFiltered), Equals, 1)
	check.Assert(allCloudsFiltered[0].NsxtAlbCloud.ID, Equals, createdAlbCloud.NsxtAlbCloud.ID)

	// Get by Name
	albCloudByName, err := vcd.client.GetAlbCloudByName(createdAlbCloud.NsxtAlbCloud.Name)
	check.Assert(err, IsNil)
	check.Assert(albCloudByName.NsxtAlbCloud.Name, Equals, createdAlbCloud.NsxtAlbCloud.Name)

	err = createdAlbCloud.Delete()
	check.Assert(err, IsNil)
}
