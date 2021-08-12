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

type NsxtAlbImportableCloud struct {
	NsxtAlbImportableCloud *types.NsxtAlbImportableCloud
	client                 *Client
	// edgeGatewayId is stored here so that pointer receiver functions can embed edge gateway ID into path
	// edgeGatewayId string
}

// GetAllAlbImportableCloud returns ALB NSX-T
func (vcdClient *VCDClient) GetAllAlbImportableCloud(parentAlbControllerUrn string, queryParameters url.Values) ([]*NsxtAlbImportableCloud, error) {
	client := vcdClient.Client
	if parentAlbControllerUrn == "" {
		return nil, fmt.Errorf("parentAlbControllerUrn is required")
	}
	if !client.IsSysAdmin {
		return nil, errors.New("handling NSX-T ALB clouds require System user")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbImportableClouds
	apiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	queryParams := copyOrNewUrlValues(queryParameters)
	queryParams = queryParameterFilterAnd(fmt.Sprintf("_context==%s", parentAlbControllerUrn), queryParams)
	typeResponses := []*types.NsxtAlbImportableCloud{{}}
	err = client.OpenApiGetAllItems(apiVersion, urlRef, queryParams, &typeResponses, nil)
	if err != nil {
		return nil, err
	}

	wrappedResponses := make([]*NsxtAlbImportableCloud, len(typeResponses))
	for sliceIndex := range typeResponses {
		wrappedResponses[sliceIndex] = &NsxtAlbImportableCloud{
			NsxtAlbImportableCloud: typeResponses[sliceIndex],
			client:                 &client,
		}
	}

	return wrappedResponses, nil
}

func (nsxtAlbController *NsxtAlbController) GetAllAlbImportableCloud(queryParameters url.Values) ([]*NsxtAlbImportableCloud, error) {
	return nsxtAlbController.vcdClient.GetAllAlbImportableCloud(nsxtAlbController.NsxtAlbController.ID, queryParameters)
}
