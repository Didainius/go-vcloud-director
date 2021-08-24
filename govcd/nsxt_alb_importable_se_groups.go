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

type NsxtAlbImportableServiceEngineGroups struct {
	NsxtAlbImportableServiceEngineGroups *types.NsxtAlbImportableServiceEngineGroups
	client                               *Client
	// edgeGatewayId is stored here so that pointer receiver functions can embed edge gateway ID into path
	// edgeGatewayId string
}

func (vcdClient *VCDClient) GetAllAlbImportableServiceEngineGroups(parentAlbCloudUrn string, queryParameters url.Values) ([]*NsxtAlbImportableServiceEngineGroups, error) {
	client := vcdClient.Client
	if parentAlbCloudUrn == "" {
		return nil, fmt.Errorf("parentAlbCloudUrn is required")
	}
	if !client.IsSysAdmin {
		return nil, errors.New("handling NSX-T ALB clouds require System user")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbImportableServiceEngineGroups
	apiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	queryParams := copyOrNewUrlValues(queryParameters)
	queryParams = queryParameterFilterAnd(fmt.Sprintf("_context==%s", parentAlbCloudUrn), queryParams)
	typeResponses := []*types.NsxtAlbImportableServiceEngineGroups{{}}
	err = client.OpenApiGetAllItems(apiVersion, urlRef, queryParams, &typeResponses, nil)
	if err != nil {
		return nil, err
	}

	wrappedResponses := make([]*NsxtAlbImportableServiceEngineGroups, len(typeResponses))
	for sliceIndex := range typeResponses {
		wrappedResponses[sliceIndex] = &NsxtAlbImportableServiceEngineGroups{
			NsxtAlbImportableServiceEngineGroups: typeResponses[sliceIndex],
			client:                               &client,
		}
	}

	return wrappedResponses, nil
}

//
//func (nsxtAlbController *NsxtAlbController) GetAllAlbImportableClouds(queryParameters url.Values) ([]*NsxtAlbImportableCloud, error) {
//	return nsxtAlbController.vcdClient.GetAllAlbImportableClouds(nsxtAlbController.NsxtAlbController.ID, queryParameters)
//}
