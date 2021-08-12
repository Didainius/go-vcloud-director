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

type NsxtAlbServiceEngineGroup struct {
	NsxtAlbServiceEngineGroup *types.NsxtAlbServiceEngineGroup
	client                    *Client
	// edgeGatewayId is stored here so that pointer receiver functions can embed edge gateway ID into path
	// edgeGatewayId string
}

// GetAllAlbImportableCloud returns ALB NSX-T
func (vcdClient *VCDClient) GetNsxtAlbServiceEngineGroup(parentAlbControllerUrn string, queryParameters url.Values) ([]*NsxtAlbServiceEngineGroup, error) {
	client := vcdClient.Client
	if parentAlbControllerUrn == "" {
		return nil, fmt.Errorf("parentAlbControllerUrn is required")
	}
	if !client.IsSysAdmin {
		return nil, errors.New("handling NSX-T ALB clouds require System user")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbServiceEngineGroups
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
	typeResponses := []*types.NsxtAlbServiceEngineGroup{{}}

	err = client.OpenApiGetAllItems(apiVersion, urlRef, queryParams, &typeResponses, nil)
	if err != nil {
		return nil, err
	}

	wrappedResponses := make([]*NsxtAlbServiceEngineGroup, len(typeResponses))
	for sliceIndex := range typeResponses {
		wrappedResponses[sliceIndex] = &NsxtAlbServiceEngineGroup{
			NsxtAlbServiceEngineGroup: typeResponses[sliceIndex],
			client:                    &client,
		}
	}

	return wrappedResponses, nil
}

func (vcdClient *VCDClient) CreateNsxtAlbServiceEngineGroup(albServiceEngineGroup *types.NsxtAlbServiceEngineGroup) (*NsxtAlbServiceEngineGroup, error) {
	client := vcdClient.Client
	if !client.IsSysAdmin {
		return nil, errors.New("handling NSX-T ALB clouds require System user")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbServiceEngineGroups
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	returnObject := &NsxtAlbServiceEngineGroup{
		NsxtAlbServiceEngineGroup: &types.NsxtAlbServiceEngineGroup{},
		client:                    &client,
	}

	err = client.OpenApiPostItem(minimumApiVersion, urlRef, nil, albServiceEngineGroup, returnObject.NsxtAlbServiceEngineGroup, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating NSX-T ALB Cloud: %s", err)
	}

	return returnObject, nil
}

func (nsxtAlbServiceEngineGroup *NsxtAlbServiceEngineGroup) Delete() error {
	client := nsxtAlbServiceEngineGroup.client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbServiceEngineGroups
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return err
	}

	if nsxtAlbServiceEngineGroup.NsxtAlbServiceEngineGroup.ID == "" {
		return fmt.Errorf("cannot delete NSX-T ALB Cloud without ID")
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, nsxtAlbServiceEngineGroup.NsxtAlbServiceEngineGroup.ID)
	if err != nil {
		return err
	}

	err = client.OpenApiDeleteItem(minimumApiVersion, urlRef, nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting NSX-T ALB Cloud: %s", err)
	}

	return nil
}
