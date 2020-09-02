/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// urn:vcloud:nsxtmanager:09722307-aee0-4623-af95-7f8e577c9ebc

type NsxtTier0Router struct {
	NsxtTier0Router *types.NsxtTier0Router
	client          *Client
}

// GetOpenApiRoleById retrieves NSX-T tier 0 router by given parent NSX-T manager ID and Tier 0 router ID
//
// Note. NSX-T manager ID format must be either UUID (09722307-aee0-4623-af95-7f8e577c9ebc) or complete
// URN (urn:vcloud:nsxtmanager:09722307-aee0-4623-af95-7f8e577c9ebc)
// func (adminOrg *AdminOrg) GetNsxtTier0RouterById(nsxtManagerId, id string) (*NsxtTier0Router, error) {
// 	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointImportableTier0Routers
// 	minimumApiVersion, err := adminOrg.client.checkOpenApiEndpointCompatibility(endpoint)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if nsxtManagerId == "" {
// 		return nil, fmt.Errorf("no NSX-T manager ID specified")
// 	}
//
// 	if id == "" {
// 		return nil, fmt.Errorf("empty Tier 0 router ID")
// 	}
//
// 	urlRef, err := adminOrg.client.OpenApiBuildEndpoint(endpoint, id)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	nsxtTier0Router := &NsxtTier0Router{
// 		NsxtTier0Router: &types.NsxtTier0Router{},
// 		client:          adminOrg.client,
// 	}
//
// 	// Get all Tier-0 routers that are accessible to an organization vDC. Routers that are already associated with an External Network are filtered out.
// 	// The “_context” filter key must be set with the id of the NSX-T manager for which we want to get the Tier-0 routers for.
// 	//
// 	// _context==urn:vcloud:nsxtmanager:09722307-aee0-4623-af95-7f8e577c9ebc
// 	queryParams := make(map[string][]string)
// 	queryParams["filter"] = []string{"_context==" + id}
//
// 	err = adminOrg.client.OpenApiGetItem(minimumApiVersion, urlRef, nil, nsxtTier0Router.NsxtTier0Router)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return nsxtTier0Router, nil
// }

// GetAllOpenApiRoles retrieves all roles using OpenAPI endpoint. Query parameters can be supplied to perform additional
// filtering.
//
// Note. IDs of Tier-0 routers do not look to have a standard and look as strings (often coinciding with DisplayName
// fields)
func (adminOrg *AdminOrg) GetAllNsxtTier0Routers(nsxtManagerId string) ([]*NsxtTier0Router, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointImportableTier0Routers
	minimumApiVersion, err := adminOrg.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := adminOrg.client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	// Get all Tier-0 routers that are accessible to an organization vDC. Routers that are already associated with an
	// External Network are filtered out. The “_context” filter key must be set with the id of the NSX-T manager for which
	// we want to get the Tier-0 routers for.
	//
	// _context==urn:vcloud:nsxtmanager:09722307-aee0-4623-af95-7f8e577c9ebc
	queryParams := make(map[string][]string)
	queryParams["filter"] = []string{"_context==" + nsxtManagerId}

	typeResponses := []*types.NsxtTier0Router{{}}

	err = adminOrg.client.OpenApiGetAllItems(minimumApiVersion, urlRef, queryParams, &typeResponses)
	if err != nil {
		return nil, err
	}

	returnObjects := make([]*NsxtTier0Router, len(typeResponses))
	for sliceIndex := range typeResponses {
		returnObjects[sliceIndex] = &NsxtTier0Router{
			NsxtTier0Router: typeResponses[sliceIndex],
			client:          adminOrg.client,
		}
	}

	return returnObjects, nil
}
