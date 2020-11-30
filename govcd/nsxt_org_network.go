/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// NsxtOrgVdcNetwork uses OpenAPI endpoint to operate NSX-T Edge Gateways
type NsxtOrgVdcNetwork struct {
	OrgVdcNetwork *types.OpenApiOrgVdcNetwork
	client        *Client
}

// GetNsxtOrgVdcNetworkById allows to retrieve NSX-T edge gateway by ID for Org admins
// func (adminOrg *AdminOrg) GetNsxtOrgVdcNetworkById(id string) (*NsxtOrgVdcNetwork, error) {
// 	return getNsxtOrgVdcNetworkById(adminOrg.client, id, nil)
// }

// GetNsxtOrgVdcNetworkById allows to retrieve NSX-T edge gateway by ID for Org users
func (org *Org) GetNsxtOrgVdcNetworkById(id string) (*NsxtOrgVdcNetwork, error) {
	// Inject Org ID filter to perform filtering on server side
	params := url.Values{}
	filterParams := queryParameterFilterAnd("orgRef.id=="+org.Org.ID, params)
	return getNsxtOrgVdcNetworkById(org.client, id, filterParams)
}

// GetNsxtOrgVdcNetworkById allows to retrieve NSX-T edge gateway by ID for specific Vdc
func (vdc *Vdc) GetNsxtOrgVdcNetworkById(id string) (*NsxtOrgVdcNetwork, error) {
	// Inject Vdc ID filter to perform filtering on server side
	params := url.Values{}
	filterParams := queryParameterFilterAnd("orgVdc.id=="+vdc.Vdc.ID, params)
	egw, err := getNsxtOrgVdcNetworkById(vdc.client, id, filterParams)
	if err != nil {
		return nil, err
	}
	//
	// if egw.OrgVdcNetwork.OrgVdc.ID != vdc.Vdc.ID {
	// 	return nil, fmt.Errorf("%s: no NSX-T edge gateway with ID '%s' found in VDC '%s'",
	// 		ErrorEntityNotFound, id, vdc.Vdc.ID)
	// }

	return egw, nil
}

// GetNsxtOrgVdcNetworkByName allows to retrieve NSX-T edge gateway by Name for Org admins
// func (adminOrg *AdminOrg) GetNsxtOrgVdcNetworkByName(name string) (*NsxtOrgVdcNetwork, error) {
// 	queryParameters := url.Values{}
// 	queryParameters.Add("filter", "name=="+name)
//
// 	allEdges, err := adminOrg.GetAllNsxtOrgVdcNetworks(queryParameters)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve edge gateway by name '%s': %s", name, err)
// 	}
//
// 	onlyNsxtEdges := filterOnlyNsxtEdges(allEdges)
//
// 	return returnSingleNsxtOrgVdcNetwork(name, onlyNsxtEdges)
// }

// GetNsxtOrgVdcNetworkByName allows to retrieve NSX-T edge gateway by Name for Org admins
func (org *Org) GetNsxtOrgVdcNetworkByName(name string) (*NsxtOrgVdcNetwork, error) {
	queryParameters := url.Values{}
	queryParameters.Add("filter", "name=="+name)
	queryParameters = queryParameterFilterAnd("orgRef.id=="+org.Org.ID, queryParameters)

	allEdges, err := org.GetAllNsxtOrgVdcNetworks(queryParameters)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve edge gateway by name '%s': %s", name, err)
	}

	// onlyNsxtEdges := filterOnlyNsxtEdges(allEdges)

	return returnSingleNsxtOrgVdcNetwork(name, allEdges)
}

// GetNsxtOrgVdcNetworkByName allows to retrieve NSX-T edge gateway by Name for specifi Vdc
func (vdc *Vdc) GetNsxtOrgVdcNetworkByName(name string) (*NsxtOrgVdcNetwork, error) {
	queryParameters := url.Values{}
	queryParameters.Add("filter", "name=="+name)

	allEdges, err := vdc.GetAllNsxtOrgVdcNetworks(queryParameters)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve edge gateway by name '%s': %s", name, err)
	}

	// onlyNsxtEdges := filterOnlyNsxtEdges(allEdges)

	return returnSingleNsxtOrgVdcNetwork(name, allEdges)
}

// GetAllNsxtOrgVdcNetworks allows to retrieve all NSX-T edge gateways for Org Admins
func (adminOrg *AdminOrg) GetAllNsxtOrgVdcNetworks(queryParameters url.Values) ([]*NsxtOrgVdcNetwork, error) {
	return getAllNsxtOrgVdcNetworks(adminOrg.client, queryParameters)
}

// GetAllNsxtOrgVdcNetworks  allows to retrieve all NSX-T edge gateways for Org users
func (org *Org) GetAllNsxtOrgVdcNetworks(queryParameters url.Values) ([]*NsxtOrgVdcNetwork, error) {
	return getAllNsxtOrgVdcNetworks(org.client, queryParameters)
}

// GetAllNsxtOrgVdcNetworks allows to retrieve all NSX-T edge gateways for Vdc users.
func (vdc *Vdc) GetAllNsxtOrgVdcNetworks(queryParameters url.Values) ([]*NsxtOrgVdcNetwork, error) {
	filteredQueryParams := queryParameterFilterAnd("orgVdc.id=="+vdc.Vdc.ID, queryParameters)
	return getAllNsxtOrgVdcNetworks(vdc.client, filteredQueryParams)
}

// CreateNsxtOrgVdcNetwork
// func (org *Org) CreateNsxtOrgVdcNetwork(OrgVdcNetworkConfig *types.OpenApiOrgVdcNetwork) (*NsxtOrgVdcNetwork, error) {
// 	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointOrgVdcNetworks
// 	minimumApiVersion, err := org.client.checkOpenApiEndpointCompatibility(endpoint)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	urlRef, err := org.client.OpenApiBuildEndpoint(endpoint)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	returnEgw := &NsxtOrgVdcNetwork{
// 		OrgVdcNetwork: &types.OpenApiOrgVdcNetwork{},
// 		client:        org.client,
// 	}
//
// 	err = org.client.OpenApiPostItem(minimumApiVersion, urlRef, nil, OrgVdcNetworkConfig, returnEgw.OrgVdcNetwork)
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating Org Vdc network: %s", err)
// 	}
//
// 	return returnEgw, nil
// }

// CreateNsxtOrgVdcNetwork
func (vdc *Vdc) CreateNsxtOrgVdcNetwork(OrgVdcNetworkConfig *types.OpenApiOrgVdcNetwork) (*NsxtOrgVdcNetwork, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointOrgVdcNetworks
	minimumApiVersion, err := vdc.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := vdc.client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	returnEgw := &NsxtOrgVdcNetwork{
		OrgVdcNetwork: &types.OpenApiOrgVdcNetwork{},
		client:        vdc.client,
	}

	err = vdc.client.OpenApiPostItem(minimumApiVersion, urlRef, nil, OrgVdcNetworkConfig, returnEgw.OrgVdcNetwork)
	if err != nil {
		return nil, fmt.Errorf("error creating Org Vdc network: %s", err)
	}

	return returnEgw, nil
}

// Update allows to update
func (orgVdcNet *NsxtOrgVdcNetwork) Update(OrgVdcNetworkConfig *types.OpenApiOrgVdcNetwork) (*NsxtOrgVdcNetwork, error) {
	if !orgVdcNet.client.IsSysAdmin {
		return nil, fmt.Errorf("only System Administrator can update Edge Gateway")
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointOrgVdcNetworks
	minimumApiVersion, err := orgVdcNet.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if OrgVdcNetworkConfig.ID == "" {
		return nil, fmt.Errorf("cannot update Edge Gateway without id")
	}

	urlRef, err := orgVdcNet.client.OpenApiBuildEndpoint(endpoint, OrgVdcNetworkConfig.ID)
	if err != nil {
		return nil, err
	}

	returnEgw := &NsxtOrgVdcNetwork{
		OrgVdcNetwork: &types.OpenApiOrgVdcNetwork{},
		client:        orgVdcNet.client,
	}

	err = orgVdcNet.client.OpenApiPutItem(minimumApiVersion, urlRef, nil, OrgVdcNetworkConfig, returnEgw.OrgVdcNetwork)
	if err != nil {
		return nil, fmt.Errorf("error updating Edge Gateway: %s", err)
	}

	return returnEgw, nil
}

// Delete allows to delete
func (orgVdcNet *NsxtOrgVdcNetwork) Delete() error {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointOrgVdcNetworks
	minimumApiVersion, err := orgVdcNet.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return err
	}

	if orgVdcNet.OrgVdcNetwork.ID == "" {
		return fmt.Errorf("cannot delete Edge Gateway without id")
	}

	urlRef, err := orgVdcNet.client.OpenApiBuildEndpoint(endpoint, orgVdcNet.OrgVdcNetwork.ID)
	if err != nil {
		return err
	}

	err = orgVdcNet.client.OpenApiDeleteItem(minimumApiVersion, urlRef, nil)

	if err != nil {
		return fmt.Errorf("error deleting Edge Gateway: %s", err)
	}

	return nil
}

// getNsxtOrgVdcNetworkById is a private parent for wrapped functions:
// func (adminOrg *AdminOrg) GetNsxtOrgVdcNetworkByName(id string) (*NsxtOrgVdcNetwork, error)
// func (org *Org) GetNsxtOrgVdcNetworkByName(id string) (*NsxtOrgVdcNetwork, error)
func getNsxtOrgVdcNetworkById(client *Client, id string, queryParameters url.Values) (*NsxtOrgVdcNetwork, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointOrgVdcNetworks
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, fmt.Errorf("empty Edge Gateway id")
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, id)
	if err != nil {
		return nil, err
	}

	egw := &NsxtOrgVdcNetwork{
		OrgVdcNetwork: &types.OpenApiOrgVdcNetwork{},
		client:        client,
	}

	err = client.OpenApiGetItem(minimumApiVersion, urlRef, queryParameters, egw.OrgVdcNetwork)
	if err != nil {
		return nil, err
	}

	// if egw.OrgVdcNetwork.GatewayBacking.GatewayType != "NSXT_BACKED" {
	// 	return nil, fmt.Errorf("%s: this is not NSX-T edge gateway (%s)",
	// 		ErrorEntityNotFound, egw.OrgVdcNetwork.GatewayBacking.GatewayType)
	// }

	return egw, nil
}

// returnSingleNsxtOrgVdcNetwork helps to reduce code duplication for `GetNsxtOrgVdcNetworkByName` functions with different
// receivers
func returnSingleNsxtOrgVdcNetwork(name string, allEdges []*NsxtOrgVdcNetwork) (*NsxtOrgVdcNetwork, error) {
	if len(allEdges) > 1 {
		return nil, fmt.Errorf("got more than 1 edge gateway by name '%s' %d", name, len(allEdges))
	}

	if len(allEdges) < 1 {
		return nil, fmt.Errorf("%s: got 0 edge gateways by name '%s'", ErrorEntityNotFound, name)
	}

	return allEdges[0], nil
}

// getAllNsxtOrgVdcNetworks is a private parent for wrapped functions:
// func (adminOrg *AdminOrg) GetAllNsxtOrgVdcNetworks(queryParameters url.Values) ([]*NsxtOrgVdcNetwork, error)
// func (org *Org) GetAllNsxtOrgVdcNetworks(queryParameters url.Values) ([]*NsxtOrgVdcNetwork, error)
func getAllNsxtOrgVdcNetworks(client *Client, queryParameters url.Values) ([]*NsxtOrgVdcNetwork, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointOrgVdcNetworks
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	typeResponses := []*types.OpenApiOrgVdcNetwork{{}}
	err = client.OpenApiGetAllItems(minimumApiVersion, urlRef, queryParameters, &typeResponses)
	if err != nil {
		return nil, err
	}

	// Wrap all typeResponses into NsxtOrgVdcNetwork types with client
	wrappedResponses := make([]*NsxtOrgVdcNetwork, len(typeResponses))
	for sliceIndex := range typeResponses {
		wrappedResponses[sliceIndex] = &NsxtOrgVdcNetwork{
			OrgVdcNetwork: typeResponses[sliceIndex],
			client:        client,
		}
	}

	// onlyNsxtEdges := filterOnlyNsxtEdges(wrappedResponses)

	return wrappedResponses, nil
}

// filterOnlyNsxtEdges filters our list of edge gateways only for NSXT_BACKED ones because original endpoint can
// return NSX-V and NSX-T backed edge gateways.
// func filterOnlyNsxtEdges(allEdges []*NsxtOrgVdcNetwork) []*NsxtOrgVdcNetwork {
// 	filteredEdges := make([]*NsxtOrgVdcNetwork, 0)
//
// 	for index := range allEdges {
// 		if allEdges[index] != nil && allEdges[index].OrgVdcNetwork != nil &&
// 			allEdges[index].OrgVdcNetwork.GatewayBacking != nil &&
// 			allEdges[index].OrgVdcNetwork.GatewayBacking.GatewayType == "NSXT_BACKED" {
// 			filteredEdges = append(filteredEdges, allEdges[index])
// 		}
// 	}
//
// 	return filteredEdges
// }
