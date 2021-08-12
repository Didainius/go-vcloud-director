/*
 * Copyright 2021 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"errors"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

type NsxtAlbController struct {
	NsxtAlbController *types.NsxtAlbController
	client            *Client
	vcdClient         *VCDClient
}

// GetAllAlbControllers returns all configured NSX-T ALB controllers
func (vcdClient *VCDClient) GetAllAlbControllers(queryParameters url.Values) ([]*NsxtAlbController, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, errors.New("reading NSX-T ALB controllers require System user")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbController
	apiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	typeResponses := []*types.NsxtAlbController{{}}
	err = client.OpenApiGetAllItems(apiVersion, urlRef, queryParameters, &typeResponses, nil)
	if err != nil {
		return nil, err
	}

	// Wrap all typeResponses into NsxtAlbController types with client
	wrappedResponses := make([]*NsxtAlbController, len(typeResponses))
	for sliceIndex := range typeResponses {
		wrappedResponses[sliceIndex] = &NsxtAlbController{
			NsxtAlbController: typeResponses[sliceIndex],
			client:            &client,
			vcdClient:         vcdClient,
		}
	}

	return wrappedResponses, nil
}
