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

// NsxtAppPortProfile uses OpenAPI endpoint to operate NSX-T Application Port Profiles
type NsxtAppPortProfile struct {
	NsxtAppPortProfile *types.NsxtAppPortProfile
	client            *Client
}

// CreateNsxtAppPortProfile allows users to create NSX-T Firewall Group
func (vdc *Vdc) CreateNsxtAppPortProfile(firewallGroupConfig *types.NsxtAppPortProfile) (*NsxtAppPortProfile, error) {
	return createNsxtAppPortProfile(vdc.client, firewallGroupConfig)
}

// CreateNsxtAppPortProfile allows users to create NSX-T Firewall Group
func (egw *NsxtEdgeGateway) CreateNsxtAppPortProfile(firewallGroupConfig *types.NsxtAppPortProfile) (*NsxtAppPortProfile, error) {
	return createNsxtAppPortProfile(egw.client, firewallGroupConfig)
}

// GetAllNsxtAppPortProfiles allows users to retrieve all Firewall Groups for Org
// firewallGroupType can be one of the following:
// * types.FirewallGroupTypeSecurityGroup - for NSX-T Security Groups
// * types.FirewallGroupTypeIpSet - for NSX-T IP Sets
// * "" (empty) - search will not be limited and will get both - IP Sets and Security Groups
//
// It is possible to add additional filtering by using queryParameters of type 'url.Values'.
// One special filter is `_context==` filtering. Value can be one of the following:
//
// * Org Vdc Network ID (_context==networkId) - Returns all the firewall groups which the specified
// network is a member of.
//
// * Edge Gateway ID (_context==edgeGatewayId) - Returns all the firewall groups which are available
// to the specific edge gateway. Or use a shorthand NsxtEdgeGateway.GetAllNsxtAppPortProfiles() which
// automatically injects this filter.
//
// * Network Provider ID (_context==networkProviderId) - Returns all the firewall groups which are
// available under a specific network provider. This context requires system admin privilege.
// 'networkProviderId' is NSX-T manager ID
func (org *Org) GetAllNsxtAppPortProfiles(queryParameters url.Values, firewallGroupType string) ([]*NsxtAppPortProfile, error) {
	queryParams := copyOrNewUrlValues(queryParameters)
	if firewallGroupType != "" {
		queryParams = queryParameterFilterAnd("type=="+firewallGroupType, queryParams)
	}

	return getAllNsxtAppPortProfiles(org.client, queryParams)
}

// GetAllNsxtAppPortProfiles allows users to retrieve all NSX-T Firewall Groups
func (vdc *Vdc) GetAllNsxtAppPortProfiles(queryParameters url.Values, firewallGroupType string) ([]*NsxtAppPortProfile, error) {
	if vdc.IsNsxv() {
		return nil, errors.New("only NSX-T VDCs support Firewall Groups")
	}
	return getAllNsxtAppPortProfiles(vdc.client, queryParameters)
}

// GetAllNsxtAppPortProfiles allows users to retrieve all NSX-T Firewall Groups in a particular Edge Gateway
// firewallGroupType can be one of the following:
// * types.FirewallGroupTypeSecurityGroup - for NSX-T Security Groups
// * types.FirewallGroupTypeIpSet - for NSX-T IP Sets
// * "" (empty) - search will not be limited and will get both - IP Sets and Security Groups
func (egw *NsxtEdgeGateway) GetAllNsxtAppPortProfiles(queryParameters url.Values, firewallGroupType string) ([]*NsxtAppPortProfile, error) {
	queryParams := copyOrNewUrlValues(queryParameters)

	if firewallGroupType != "" {
		queryParams = queryParameterFilterAnd("type=="+firewallGroupType, queryParams)
	}

	// Automatically inject Edge Gateway filter because this is an Edge Gateway scoped query
	queryParams = queryParameterFilterAnd("_context=="+egw.EdgeGateway.ID, queryParams)

	return getAllNsxtAppPortProfiles(egw.client, queryParams)
}

// GetNsxtAppPortProfileByName allows users to retrieve Firewall Group by Name
// firewallGroupType can be one of the following:
// * types.FirewallGroupTypeSecurityGroup - for NSX-T Security Groups
// * types.FirewallGroupTypeIpSet - for NSX-T IP Sets
// * "" (empty) - search will not be limited and will get both - IP Sets and Security Groups
//
// Note. One might get an error if IP Set and Security Group exist with the same name (two objects
// of the same type cannot exist) and firewallGroupType is left empty.
func (org *Org) GetNsxtAppPortProfileByName(name, firewallGroupType string) (*NsxtAppPortProfile, error) {
	queryParameters := url.Values{}
	if firewallGroupType != "" {
		queryParameters = queryParameterFilterAnd("type=="+firewallGroupType, queryParameters)
	}

	return getNsxtAppPortProfileByName(org.client, name, queryParameters)
}

// GetNsxtAppPortProfileByName allows users to retrieve Firewall Group by Name
// firewallGroupType can be one of the following:
// * types.FirewallGroupTypeSecurityGroup - for NSX-T Security Groups
// * types.FirewallGroupTypeIpSet - for NSX-T IP Sets
// * "" (empty) - search will not be limited and will get both - IP Sets and Security Groups
//
// Note. One might get an error if IP Set and Security Group exist with the same name (two objects
// of the same type cannot exist) and firewallGroupType is left empty.
func (vdc *Vdc) GetNsxtAppPortProfileByName(name, firewallGroupType string) (*NsxtAppPortProfile, error) {

	queryParameters := url.Values{}
	if firewallGroupType != "" {
		queryParameters = queryParameterFilterAnd("type=="+firewallGroupType, queryParameters)
	}
	return getNsxtAppPortProfileByName(vdc.client, name, queryParameters)
}

// GetNsxtAppPortProfileByName allows users to retrieve Firewall Group by Name in a particular Edge Gateway
// firewallGroupType can be one of the following:
// * types.FirewallGroupTypeSecurityGroup - for NSX-T Security Groups
// * types.FirewallGroupTypeIpSet - for NSX-T IP Sets
// * "" (empty) - search will not be limited and will get both - IP Sets and Security Groups
//
// Note. One might get an error if IP Set and Security Group exist with the same name (two objects
// of the same type cannot exist) and firewallGroupType is left empty.
func (egw *NsxtEdgeGateway) GetNsxtAppPortProfileByName(name string, firewallGroupType string) (*NsxtAppPortProfile, error) {
	queryParameters := url.Values{}

	if firewallGroupType != "" {
		queryParameters = queryParameterFilterAnd("type=="+firewallGroupType, queryParameters)
	}

	// Automatically inject Edge Gateway filter because this is an Edge Gateway scoped query
	queryParameters = queryParameterFilterAnd("_context=="+egw.EdgeGateway.ID, queryParameters)

	return getNsxtAppPortProfileByName(egw.client, name, queryParameters)
}

// GetNsxtAppPortProfileById retrieves NSX-T Firewall Group by ID
func (org *Org) GetNsxtAppPortProfileById(id string) (*NsxtAppPortProfile, error) {
	return getNsxtAppPortProfileById(org.client, id)
}

// GetNsxtAppPortProfileById retrieves NSX-T Firewall Group by ID
func (vdc *Vdc) GetNsxtAppPortProfileById(id string) (*NsxtAppPortProfile, error) {
	return getNsxtAppPortProfileById(vdc.client, id)
}

// GetNsxtAppPortProfileById retrieves NSX-T Firewall Group by ID
func (egw *NsxtEdgeGateway) GetNsxtAppPortProfileById(id string) (*NsxtAppPortProfile, error) {
	return getNsxtAppPortProfileById(egw.client, id)
}

// Update allows users to update NSX-T Firewall Group
func (firewallGroup *NsxtAppPortProfile) Update(firewallGroupConfig *types.NsxtAppPortProfile) (*NsxtAppPortProfile, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAppPortProfiles
	minimumApiVersion, err := firewallGroup.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if firewallGroupConfig.ID == "" {
		return nil, fmt.Errorf("cannot update NSX-T Firewall Group without ID")
	}

	urlRef, err := firewallGroup.client.OpenApiBuildEndpoint(endpoint, firewallGroupConfig.ID)
	if err != nil {
		return nil, err
	}

	returnObject := &NsxtAppPortProfile{
		NsxtAppPortProfile: &types.NsxtAppPortProfile{},
		client:            firewallGroup.client,
	}

	err = firewallGroup.client.OpenApiPutItem(minimumApiVersion, urlRef, nil, firewallGroupConfig, returnObject.NsxtAppPortProfile)
	if err != nil {
		return nil, fmt.Errorf("error updating NSX-T firewall group: %s", err)
	}

	return returnObject, nil
}

// Delete allows users to delete NSX-T Firewall Group
func (firewallGroup *NsxtAppPortProfile) Delete() error {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAppPortProfiles
	minimumApiVersion, err := firewallGroup.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return err
	}

	if firewallGroup.NsxtAppPortProfile.ID == "" {
		return fmt.Errorf("cannot delete NSX-T Firewall Group without ID")
	}

	urlRef, err := firewallGroup.client.OpenApiBuildEndpoint(endpoint, firewallGroup.NsxtAppPortProfile.ID)
	if err != nil {
		return err
	}

	err = firewallGroup.client.OpenApiDeleteItem(minimumApiVersion, urlRef, nil)

	if err != nil {
		return fmt.Errorf("error deleting NSX-T Firewall Group: %s", err)
	}

	return nil
}

func getNsxtAppPortProfileByName(client *Client, name string, queryParameters url.Values) (*NsxtAppPortProfile, error) {
	queryParams := copyOrNewUrlValues(queryParameters)
	queryParams = queryParameterFilterAnd("name=="+name, queryParams)

	allGroups, err := getAllNsxtAppPortProfiles(client, queryParams)
	if err != nil {
		return nil, fmt.Errorf("could not find NSX-T Firewall Group with name '%s': %s", name, err)
	}

	if len(allGroups) == 0 {
		return nil, fmt.Errorf("%s: expected exactly one NSX-T Firewall Group with name '%s'. Got %d", ErrorEntityNotFound, name, len(allGroups))
	}

	if len(allGroups) > 1 {
		return nil, fmt.Errorf("expected exactly one NSX-T Firewall Group with name '%s'. Got %d", name, len(allGroups))
	}

	// TODO API V36.0 - maybe it is fixed
	// There is a bug that not all data is present (e.g. missing IpAddresses field for IP_SET) when
	// using "getAll" endpoint therefore after finding the object by name we must retrieve it once
	// again using its direct endpoint.
	//
	// return allGroups[0], nil

	return getNsxtAppPortProfileById(client, allGroups[0].NsxtAppPortProfile.ID)
}

func getNsxtAppPortProfileById(client *Client, id string) (*NsxtAppPortProfile, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAppPortProfiles
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, fmt.Errorf("empty NSX-T Firewall Group ID specified")
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, id)
	if err != nil {
		return nil, err
	}

	fwGroup := &NsxtAppPortProfile{
		NsxtAppPortProfile: &types.NsxtAppPortProfile{},
		client:            client,
	}

	err = client.OpenApiGetItem(minimumApiVersion, urlRef, nil, fwGroup.NsxtAppPortProfile)
	if err != nil {
		return nil, err
	}

	return fwGroup, nil
}

func getAllNsxtAppPortProfiles(client *Client, queryParameters url.Values) ([]*NsxtAppPortProfile, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAppPortProfiles
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	// This Object does not follow regular REST scheme and for get the endpoint must be
	// 1.0.0/firewallGroups/summaries therefore bellow "summaries" is appended to the path
	urlRef, err := client.OpenApiBuildEndpoint(endpoint, "summaries")
	if err != nil {
		return nil, err
	}

	typeResponses := []*types.NsxtAppPortProfile{{}}
	err = client.OpenApiGetAllItems(minimumApiVersion, urlRef, queryParameters, &typeResponses)
	if err != nil {
		return nil, err
	}

	// Wrap all typeResponses into NsxtEdgeGateway types with client
	wrappedResponses := make([]*NsxtAppPortProfile, len(typeResponses))
	for sliceIndex := range typeResponses {
		wrappedResponses[sliceIndex] = &NsxtAppPortProfile{
			NsxtAppPortProfile: typeResponses[sliceIndex],
			client:            client,
		}
	}

	return wrappedResponses, nil
}

func createNsxtAppPortProfile(client *Client, firewallGroupConfig *types.NsxtAppPortProfile) (*NsxtAppPortProfile, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAppPortProfiles
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	returnObject := &NsxtAppPortProfile{
		NsxtAppPortProfile: &types.NsxtAppPortProfile{},
		client:            client,
	}

	err = client.OpenApiPostItem(minimumApiVersion, urlRef, nil, firewallGroupConfig, returnObject.NsxtAppPortProfile)
	if err != nil {
		return nil, fmt.Errorf("error creating NSX-T Firewall Group: %s", err)
	}

	return returnObject, nil
}
