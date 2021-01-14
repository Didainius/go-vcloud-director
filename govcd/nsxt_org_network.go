/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// NsxtOrgVdcNetwork uses OpenAPI endpoint to operate NSX-T Edge Gateways
type NsxtOrgVdcNetwork struct {
	OrgVdcNetwork *types.OpenApiOrgVdcNetwork
	client        *Client
}

// GetNsxtOrgVdcNetworkById allows to retrieve NSX-T edge gateway by ID for Org
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

	return egw, nil
}

// GetNsxtOrgVdcNetworkByName allows to retrieve NSX-T edge gateway by Name in the VDC
func (vdc *Vdc) GetNsxtOrgVdcNetworkByName(name string) (*NsxtOrgVdcNetwork, error) {
	queryParameters := url.Values{}
	queryParameters.Add("filter", "name=="+name)

	allEdges, err := vdc.GetAllNsxtOrgVdcNetworks(queryParameters)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Org Vdc network by name '%s': %s", name, err)
	}

	// onlyNsxtEdges := filterOnlyNsxtEdges(allEdges)

	return returnSingleNsxtOrgVdcNetwork(name, allEdges)
}

// GetAllNsxtOrgVdcNetworks allows to retrieve all NSX-T Or gVdc networks in Vdc
//
// Note. If pageSize > 32 it will be limited to maximum of 32 in this function because API validation does not allow for
// higher number
func (vdc *Vdc) GetAllNsxtOrgVdcNetworks(queryParameters url.Values) ([]*NsxtOrgVdcNetwork, error) {
	filteredQueryParams := queryParameterFilterAnd("orgVdc.id=="+vdc.Vdc.ID, queryParameters)
	return getAllNsxtOrgVdcNetworks(vdc.client, filteredQueryParams)
}

// CreateNsxtOrgVdcNetwork allows to create NSX-T Org Vdc network
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

// Update allows to update Org Vdc network
func (orgVdcNet *NsxtOrgVdcNetwork) Update(OrgVdcNetworkConfig *types.OpenApiOrgVdcNetwork) (*NsxtOrgVdcNetwork, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointOrgVdcNetworks
	minimumApiVersion, err := orgVdcNet.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if OrgVdcNetworkConfig.ID == "" {
		return nil, fmt.Errorf("cannot update Org Vdc network without id")
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
		return nil, fmt.Errorf("error updating Org Vdc network: %s", err)
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
		return fmt.Errorf("cannot delete Org Vdc network without id")
	}

	urlRef, err := orgVdcNet.client.OpenApiBuildEndpoint(endpoint, orgVdcNet.OrgVdcNetwork.ID)
	if err != nil {
		return err
	}

	err = orgVdcNet.client.OpenApiDeleteItem(minimumApiVersion, urlRef, nil)

	if err != nil {
		return fmt.Errorf("error deleting Org Vdc: %s", err)
	}

	return nil
}

func (orgVdcNet *NsxtOrgVdcNetwork) GetType() string {

	return "nil"
}

func (orgVdcNet *NsxtOrgVdcNetwork) IsIsolated() bool {

	return false
}

func (orgVdcNet *NsxtOrgVdcNetwork) IsRouted() bool {

	return false
}

func (orgVdcNet *NsxtOrgVdcNetwork) IsImported() bool {

	return false
}

// getNsxtOrgVdcNetworkById is a private parent for wrapped functions:
// func (org *Org) GetNsxtOrgVdcNetworkById(id string) (*NsxtOrgVdcNetwork, error)
// func (vdc *Vdc) GetNsxtOrgVdcNetworkById(id string) (*NsxtOrgVdcNetwork, error)
func getNsxtOrgVdcNetworkById(client *Client, id string, queryParameters url.Values) (*NsxtOrgVdcNetwork, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointOrgVdcNetworks
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, fmt.Errorf("empty Org Vdc network id")
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

	return egw, nil
}

// returnSingleNsxtOrgVdcNetwork helps to reduce code duplication for `GetNsxtOrgVdcNetworkByName` functions with different
// receivers
func returnSingleNsxtOrgVdcNetwork(name string, allEdges []*NsxtOrgVdcNetwork) (*NsxtOrgVdcNetwork, error) {
	if len(allEdges) > 1 {
		return nil, fmt.Errorf("got more than 1 Org Vdc network by name '%s' %d", name, len(allEdges))
	}

	if len(allEdges) < 1 {
		return nil, fmt.Errorf("%s: got 0 Org Vdc network by name '%s'", ErrorEntityNotFound, name)
	}

	return allEdges[0], nil
}

// getAllNsxtOrgVdcNetworks is a private parent for wrapped functions:
// func (vdc *Vdc) GetAllNsxtOrgVdcNetworks(queryParameters url.Values) ([]*NsxtOrgVdcNetwork, error)
//
// Note. If pageSize > 32 it will be limited to maximum of 32 in this function because API validation does not allow for
// higher number
func getAllNsxtOrgVdcNetworks(client *Client, queryParameters url.Values) ([]*NsxtOrgVdcNetwork, error) {

	// Enforce maximum pageSize to be 32 as API endpoint throws error if it is > 32
	pageSizeString := queryParameters.Get("pageSize")

	switch pageSizeString {
	// If no pageSize is specified it must be set to 32 as by default low level API function OpenApiGetAllItems sets 128
	case "":
		queryParameters.Set("pageSize", "32")

		// If pageSize is specified ensure it is not >32
	default:
		pageSizeValue, err := strconv.Atoi(pageSizeString)
		if err != nil {
			return nil, fmt.Errorf("error parsing pageSize value: %s", err)
		}
		if pageSizeString != "" && pageSizeValue > 32 {
			queryParameters.Set("pageSize", "32")
		}
	}

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

	return wrappedResponses, nil
}
