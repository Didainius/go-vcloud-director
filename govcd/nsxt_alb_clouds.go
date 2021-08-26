/*
 * Copyright 2021 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// NsxtAlbCloud is a service provider-level construct that consists of an NSX-T Manager and an NSX-T Data Center
// transport zone. An NSX-T Data Center transport zone dictates which hosts and virtual machines can participate in the
// use of a particular network. An NSX-T Cloud has a one-to-one relationship with a network pool backed by an NSX-T Data
// Center transport zone.
type NsxtAlbCloud struct {
	NsxtAlbCloud *types.NsxtAlbCloud
	client       *Client
}

// GetAllAlbClouds returns all configured NSX-T ALB Clouds
func (vcdClient *VCDClient) GetAllAlbClouds(queryParameters url.Values) ([]*NsxtAlbCloud, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, errors.New("handling NSX-T ALB Clouds require System user")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbCloud
	apiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	typeResponses := []*types.NsxtAlbCloud{{}}
	err = client.OpenApiGetAllItems(apiVersion, urlRef, queryParameters, &typeResponses, nil)
	if err != nil {
		return nil, err
	}

	// Wrap all typeResponses into NsxtAlbController types with client
	wrappedResponses := make([]*NsxtAlbCloud, len(typeResponses))
	for sliceIndex := range typeResponses {
		wrappedResponses[sliceIndex] = &NsxtAlbCloud{
			NsxtAlbCloud: typeResponses[sliceIndex],
			client:       &client,
		}
	}

	return wrappedResponses, nil
}

// GetAlbCloudByName returns ALB cloud by name
func (vcdClient *VCDClient) GetAlbCloudByName(name string) (*NsxtAlbCloud, error) {

	queryParameters := copyOrNewUrlValues(nil)
	queryParameters.Add("filter", "name=="+name)

	controllers, err := vcdClient.GetAllAlbClouds(queryParameters)
	if err != nil {
		return nil, fmt.Errorf("error reading ALB Cloud with Name '%s': %s", name, err)
	}

	if len(controllers) == 0 {
		return nil, fmt.Errorf("%s could not find ALB Cloud by Name '%s'", ErrorEntityNotFound, name)
	}

	if len(controllers) > 1 {
		return nil, fmt.Errorf("found more than 1 ALB Cloud by Name '%s'", name)
	}

	return controllers[0], nil
}

func (vcdClient *VCDClient) CreateAlbCloud(albCloudConfig *types.NsxtAlbCloud) (*NsxtAlbCloud, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, errors.New("handling NSX-T ALB Clouds require System user")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbCloud
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	returnObject := &NsxtAlbCloud{
		NsxtAlbCloud: &types.NsxtAlbCloud{},
		client:       &client,
	}

	err = client.OpenApiPostItem(minimumApiVersion, urlRef, nil, albCloudConfig, returnObject.NsxtAlbCloud, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating NSX-T ALB Cloud: %s", err)
	}

	return returnObject, nil
}

// Update is not supported up to at least VCD 10.3 therefore this function remains commented
//
// Update updates existing NSX-T ALB Cloud with new supplied albCloudConfig configuration
//func (nsxtAlbCloud *NsxtAlbCloud) Update(albCloudConfig *types.NsxtAlbCloud) (*NsxtAlbCloud, error) {
//	client := nsxtAlbCloud.client
//	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbCloud
//	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
//	if err != nil {
//		return nil, err
//	}
//
//	if nsxtAlbCloud.NsxtAlbCloud.ID == "" {
//		return nil, fmt.Errorf("cannot update NSX-T ALB Cloud without ID")
//	}
//
//	urlRef, err := client.OpenApiBuildEndpoint(endpoint, nsxtAlbCloud.NsxtAlbCloud.ID)
//	if err != nil {
//		return nil, err
//	}
//
//	responseAlbCloud := &NsxtAlbCloud{
//		NsxtAlbCloud: &types.NsxtAlbCloud{},
//		client:       nsxtAlbCloud.client,
//	}
//
//	err = client.OpenApiPutItem(minimumApiVersion, urlRef, nil, albCloudConfig, responseAlbCloud.NsxtAlbCloud, nil)
//	if err != nil {
//		return nil, fmt.Errorf("error updating NSX-T ALB Cloud: %s", err)
//	}
//
//	return responseAlbCloud, nil
//}

// Delete removes NSX-T ALB Cloud configuration
func (nsxtAlbCloud *NsxtAlbCloud) Delete() error {
	client := nsxtAlbCloud.client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbCloud
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return err
	}

	if nsxtAlbCloud.NsxtAlbCloud.ID == "" {
		return fmt.Errorf("cannot delete NSX-T ALB Cloud without ID")
	}

	urlRef, err := nsxtAlbCloud.client.OpenApiBuildEndpoint(endpoint, nsxtAlbCloud.NsxtAlbCloud.ID)
	if err != nil {
		return err
	}

	err = nsxtAlbCloud.client.OpenApiDeleteItem(minimumApiVersion, urlRef, nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting NSX-T ALB Cloud: %s", err)
	}

	return nil
}
