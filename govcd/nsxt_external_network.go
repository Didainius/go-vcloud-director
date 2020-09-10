/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

type ExternalNetworkV2 struct {
	ExternalNetwork *types.ExternalNetworkV2
	client          *Client
}

// GetExternalNetworkById retrieves external network by given ID
func (adminOrg *AdminOrg) GetExternalNetworkById(id string) (*ExternalNetworkV2, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointExternalNetworks
	minimumApiVersion, err := adminOrg.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, fmt.Errorf("empty external network id")
	}

	urlRef, err := adminOrg.client.OpenApiBuildEndpoint(endpoint, id)
	if err != nil {
		return nil, err
	}

	extNet := &ExternalNetworkV2{
		ExternalNetwork: &types.ExternalNetworkV2{},
		client:          adminOrg.client,
	}

	err = adminOrg.client.OpenApiGetItem(minimumApiVersion, urlRef, nil, extNet.ExternalNetwork)
	if err != nil {
		return nil, err
	}

	return extNet, nil
}

// GetAllExternalNetworks retrieves all roles using OpenAPI endpoint. Query parameters can be supplied to perform
// additional filtering
func (adminOrg *AdminOrg) GetAllExternalNetworks(queryParameters url.Values) ([]*ExternalNetworkV2, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointExternalNetworks
	minimumApiVersion, err := adminOrg.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := adminOrg.client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	typeResponses := []*types.ExternalNetworkV2{{}}
	err = adminOrg.client.OpenApiGetAllItems(minimumApiVersion, urlRef, queryParameters, &typeResponses)
	if err != nil {
		return nil, err
	}

	// Wrap all typeResponses into external network types with client
	returnExtNetworks := make([]*ExternalNetworkV2, len(typeResponses))
	for sliceIndex := range typeResponses {
		returnExtNetworks[sliceIndex] = &ExternalNetworkV2{
			ExternalNetwork: typeResponses[sliceIndex],
			client:          adminOrg.client,
		}
	}

	return returnExtNetworks, nil
}

func CreateExternalNetworkV2(vcdClient *VCDClient, newExtNet *types.ExternalNetworkV2) (*ExternalNetworkV2, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointExternalNetworks
	minimumApiVersion, err := vcdClient.Client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := vcdClient.Client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	returnRole := &ExternalNetworkV2{
		ExternalNetwork: &types.ExternalNetworkV2{},
		client:          &vcdClient.Client,
	}

	err = vcdClient.Client.OpenApiPostItem(minimumApiVersion, urlRef, nil, newExtNet, returnRole.ExternalNetwork)
	if err != nil {
		return nil, fmt.Errorf("error creating external network: %s", err)
	}

	return returnRole, nil
}

// CreateExternalNetwork creates a new external network using OpenAPI endpoint
func (adminOrg *AdminOrg) CreateExternalNetwork(newExtNet *types.ExternalNetworkV2) (*ExternalNetworkV2, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointExternalNetworks
	minimumApiVersion, err := adminOrg.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := adminOrg.client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	returnRole := &ExternalNetworkV2{
		ExternalNetwork: &types.ExternalNetworkV2{},
		client:          adminOrg.client,
	}

	err = adminOrg.client.OpenApiPostItem(minimumApiVersion, urlRef, nil, newExtNet, returnRole.ExternalNetwork)
	if err != nil {
		return nil, fmt.Errorf("error creating external network: %s", err)
	}

	return returnRole, nil
}

// Update updates existing OpenAPI external network
func (role *ExternalNetworkV2) Update() (*ExternalNetworkV2, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointExternalNetworks
	minimumApiVersion, err := role.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if role.ExternalNetwork.ID == "" {
		return nil, fmt.Errorf("cannot update external network without id")
	}

	urlRef, err := role.client.OpenApiBuildEndpoint(endpoint, role.ExternalNetwork.ID)
	if err != nil {
		return nil, err
	}

	returnExtNet := &ExternalNetworkV2{
		ExternalNetwork: &types.ExternalNetworkV2{},
		client:          role.client,
	}

	err = role.client.OpenApiPutItem(minimumApiVersion, urlRef, nil, role.ExternalNetwork, returnExtNet.ExternalNetwork)
	if err != nil {
		return nil, fmt.Errorf("error updating external network: %s", err)
	}

	return returnExtNet, nil
}

// Delete deletes OpenAPI external network
func (role *ExternalNetworkV2) Delete() error {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointExternalNetworks
	minimumApiVersion, err := role.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return err
	}

	if role.ExternalNetwork.ID == "" {
		return fmt.Errorf("cannot delete external network without id")
	}

	urlRef, err := role.client.OpenApiBuildEndpoint(endpoint, role.ExternalNetwork.ID)
	if err != nil {
		return err
	}

	err = role.client.OpenApiDeleteItem(minimumApiVersion, urlRef, nil)

	if err != nil {
		return fmt.Errorf("error deleting role: %s", err)
	}

	return nil
}
