/*
 * Copyright 2021 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// NsxtFirewallGroup uses OpenAPI endpoint to operate NSX-T Security Groups and IP Sets which use
// the same Firewall Group API endpoint
type NsxtFirewallGroup struct {
	NsxtFirewallGroup *types.NsxtFirewallGroup
	client            *Client
}

// GetNsxtFirewallGroupByName retrieves NSX-T Firewall Group by name
//
// Note. Name uniqueness is enforced in the API so there can only be one result
func (adminOrg *AdminOrg) GetNsxtFirewallGroupByName(name, firewallGroupType string) (*NsxtFirewallGroup, error) {
	validateFirewallGroupType(firewallGroupType)

	queryParameters := url.Values{}
	queryParameters = queryParameterFilterAnd("type=="+firewallGroupType, queryParameters)
	return getNsxtFirewallGroupByName(adminOrg.client, name, queryParameters)
}

func (org *Org) GetNsxtFirewallGroupByName(name, firewallGroupType string) (*NsxtFirewallGroup, error) {
	validateFirewallGroupType(firewallGroupType)

	queryParameters := url.Values{}
	queryParameters = queryParameterFilterAnd("type=="+firewallGroupType, queryParameters)
	return getNsxtFirewallGroupByName(org.client, name, queryParameters)
}

func (vdc *Vdc) GetNsxtFirewallGroupByName(name, firewallGroupType string) (*NsxtFirewallGroup, error) {
	validateFirewallGroupType(firewallGroupType)

	queryParameters := url.Values{}
	queryParameters = queryParameterFilterAnd("type=="+firewallGroupType, queryParameters)
	return getNsxtFirewallGroupByName(vdc.client, name, queryParameters)
}

func validateFirewallGroupType(firewallGroupType string) error {
	if firewallGroupType != "IP_SET" && firewallGroupType != "SECURITY_GROUP" {
		return fmt.Errorf("NSX-T Firewall Group type can be 'IP_SET' or 'SECURITY_GROUP', not '%s'",
			firewallGroupType)
	}

	return nil
}

// GetNsxtFirewallGroupByName will limit scope of Firewall Groups
func (egw *NsxtEdgeGateway) GetNsxtFirewallGroupByName(name string, firewallGroupType string) (*NsxtFirewallGroup, error) {
	queryParameters := url.Values{}

	queryParameters = queryParameterFilterAnd("type=="+firewallGroupType, queryParameters)
	queryParameters = queryParameterFilterAnd("_context=="+egw.EdgeGateway.ID, queryParameters)

	return getNsxtFirewallGroupByName(egw.client, name, queryParameters)
}

func getNsxtFirewallGroupByName(client *Client, name string, queryParameters url.Values) (*NsxtFirewallGroup, error) {
	queryParams := copyOrNewUrlValues(queryParameters)
	queryParams = queryParameterFilterAnd("name=="+name, queryParams)

	allGroups, err := getAllNsxtFirewallGroups(client, queryParams)
	if err != nil {
		return nil, fmt.Errorf("could not find NSX-T Firewall Group with name '%s': %s", name, err)
	}

	if len(allGroups) == 0 {
		return nil, fmt.Errorf("%s: expected exactly one NSX-T Firewall Group with name '%s'. Got %d", ErrorEntityNotFound, name, len(allGroups))
	}

	if len(allGroups) > 1 {
		return nil, fmt.Errorf("expected exactly one NSX-T Firewall Group with name '%s'. Got %d", name, len(allGroups))
	}

	return allGroups[0], nil
}

// GetNsxtFirewallGroupById retrieves NSX-T Firewall Group by id
func (adminOrg *AdminOrg) GetNsxtFirewallGroupById(id string) (*NsxtFirewallGroup, error) {
	return getNsxtFirewallGroupById(adminOrg.client, id)
}

// GetNsxtFirewallGroupById retrieves NSX-T Firewall Group by id
func (org *Org) GetNsxtFirewallGroupById(id string) (*NsxtFirewallGroup, error) {
	return getNsxtFirewallGroupById(org.client, id)
}

func (vdc *Vdc) GetNsxtFirewallGroupById(id string) (*NsxtFirewallGroup, error) {
	return getNsxtFirewallGroupById(vdc.client, id)
}

func getNsxtFirewallGroupById(client *Client, id string) (*NsxtFirewallGroup, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointFirewallGroups
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

	fwGroup := &NsxtFirewallGroup{
		NsxtFirewallGroup: &types.NsxtFirewallGroup{},
		client:            client,
	}

	err = client.OpenApiGetItem(minimumApiVersion, urlRef, nil, fwGroup.NsxtFirewallGroup)
	if err != nil {
		return nil, err
	}

	return fwGroup, nil
}

// GetAllNsxtFirewallGroups allows to retrieve all NSX-T edge gateways for Org users
//
// It is possible to add additional filtering by using `_context==` filtering. Value can be one of
// the following:
// * Org Vdc Network ID (_context==networkId) - Returns all the firewall groups which the specified
// network is a member of.
// * Edge Gateway ID (_context==edgeGatewayId) - Returns all the firewall
// groups which are available to the specific edge gateway.
// * Network Provider ID (_context==networkProviderId) - Returns all the firewall groups which are
// available under a specific network provider. This context requires system admin privilege.
// 'networkProviderId' is NSX-T manager ID
func (org *Org) GetAllNsxtFirewallGroups(queryParameters url.Values) ([]*NsxtFirewallGroup, error) {
	return getAllNsxtFirewallGroups(org.client, queryParameters)
}

func (adminOrg *AdminOrg) GetAllNsxtFirewallGroups(queryParameters url.Values) ([]*NsxtFirewallGroup, error) {
	return getAllNsxtFirewallGroups(adminOrg.client, queryParameters)
}

func (vdc *Vdc) GetAllNsxtFirewallGroups(queryParameters url.Values) ([]*NsxtFirewallGroup, error) {
	return getAllNsxtFirewallGroups(vdc.client, queryParameters)
}

func getAllNsxtFirewallGroups(client *Client, queryParameters url.Values) ([]*NsxtFirewallGroup, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointFirewallGroups
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

	typeResponses := []*types.NsxtFirewallGroup{{}}
	err = client.OpenApiGetAllItems(minimumApiVersion, urlRef, queryParameters, &typeResponses)
	if err != nil {
		return nil, err
	}

	// Wrap all typeResponses into NsxtEdgeGateway types with client
	wrappedResponses := make([]*NsxtFirewallGroup, len(typeResponses))
	for sliceIndex := range typeResponses {
		wrappedResponses[sliceIndex] = &NsxtFirewallGroup{
			NsxtFirewallGroup: typeResponses[sliceIndex],
			client:            client,
		}
	}

	return wrappedResponses, nil
}

// CreateNsxtFirewallGroup allows to create NSX-T Firewall Group
func (adminOrg *AdminOrg) CreateNsxtFirewallGroup(firewallGroupConfig *types.NsxtFirewallGroup) (*NsxtFirewallGroup, error) {
	return createNsxtFirewallGroup(adminOrg.client, firewallGroupConfig)
}

// CreateNsxtFirewallGroup allows to create NSX-T Firewall Group
func (org *Org) CreateNsxtFirewallGroup(firewallGroupConfig *types.NsxtFirewallGroup) (*NsxtFirewallGroup, error) {
	return createNsxtFirewallGroup(org.client, firewallGroupConfig)
}

// CreateNsxtFirewallGroup allows to create NSX-T Firewall Group
func (vdc *Vdc) CreateNsxtFirewallGroup(firewallGroupConfig *types.NsxtFirewallGroup) (*NsxtFirewallGroup, error) {
	return createNsxtFirewallGroup(vdc.client, firewallGroupConfig)
}

// CreateNsxtFirewallGroup allows to create NSX-T Firewall Group
func (egw *NsxtEdgeGateway) CreateNsxtFirewallGroup(firewallGroupConfig *types.NsxtFirewallGroup) (*NsxtFirewallGroup, error) {
	return createNsxtFirewallGroup(egw.client, firewallGroupConfig)
}

func createNsxtFirewallGroup(client *Client, firewallGroupConfig *types.NsxtFirewallGroup) (*NsxtFirewallGroup, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointFirewallGroups
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	returnObject := &NsxtFirewallGroup{
		NsxtFirewallGroup: &types.NsxtFirewallGroup{},
		client:            client,
	}

	err = client.OpenApiPostItem(minimumApiVersion, urlRef, nil, firewallGroupConfig, returnObject.NsxtFirewallGroup)
	if err != nil {
		return nil, fmt.Errorf("error creating NSX-T Firewall Group: %s", err)
	}

	return returnObject, nil
}

// Update allows to update NSX-T Firewall Group
func (firewallGroup *NsxtFirewallGroup) Update(firewallGroupConfig *types.NsxtFirewallGroup) (*NsxtFirewallGroup, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointFirewallGroups
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

	returnObject := &NsxtFirewallGroup{
		NsxtFirewallGroup: &types.NsxtFirewallGroup{},
		client:            firewallGroup.client,
	}

	err = firewallGroup.client.OpenApiPutItem(minimumApiVersion, urlRef, nil, firewallGroupConfig, returnObject.NsxtFirewallGroup)
	if err != nil {
		return nil, fmt.Errorf("error updating NSX-T firewall group: %s", err)
	}

	return returnObject, nil
}

// Delete allows to delete NSX-T Firewall Group
func (firewallGroup *NsxtFirewallGroup) Delete() error {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointFirewallGroups
	minimumApiVersion, err := firewallGroup.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return err
	}

	if firewallGroup.NsxtFirewallGroup.ID == "" {
		return fmt.Errorf("cannot delete NSX-T Firewall Group without ID")
	}

	urlRef, err := firewallGroup.client.OpenApiBuildEndpoint(endpoint, firewallGroup.NsxtFirewallGroup.ID)
	if err != nil {
		return err
	}

	err = firewallGroup.client.OpenApiDeleteItem(minimumApiVersion, urlRef, nil)

	if err != nil {
		return fmt.Errorf("error deleting NSX-T Firewall Group: %s", err)
	}

	return nil
}

// GetAssociatedVms allows to retrieve a list of references to child VMs (with vApps if exist)
func (firewallGroup *NsxtFirewallGroup) GetAssociatedVms() ([]*types.NsxtFirewallGroupMemberVms, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointFirewallGroups
	minimumApiVersion, err := firewallGroup.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if firewallGroup.NsxtFirewallGroup.ID == "" {
		return nil, fmt.Errorf("cannot retrieve associated VMs for  NSX-T Firewall Group without ID")
	}

	if !firewallGroup.IsSecurityGroup() {
		return nil, fmt.Errorf("only Security Groups have associated VMs. This Firewall Group has type '%s'",
			firewallGroup.NsxtFirewallGroup.Type)
	}

	urlRef, err := firewallGroup.client.OpenApiBuildEndpoint(endpoint, firewallGroup.NsxtFirewallGroup.ID, "/associatedVMs")
	if err != nil {
		return nil, err
	}

	associatedVms := []*types.NsxtFirewallGroupMemberVms{{}}

	err = firewallGroup.client.OpenApiGetAllItems(minimumApiVersion, urlRef, nil, &associatedVms)

	if err != nil {
		return nil, fmt.Errorf("error retrieving associated VMs: %s", err)
	}

	return associatedVms, nil
}

// IsSecurityGroup allows to check if Firewall Group is a Security Group
func (firewallGroup *NsxtFirewallGroup) IsSecurityGroup() bool {
	return firewallGroup.NsxtFirewallGroup.Type == types.FirewallGroupTypeSecurityGroup
}

// IsIpSet allows to check if Firewall Group is an IP Set
func (firewallGroup *NsxtFirewallGroup) IsIpSet() bool {
	return firewallGroup.NsxtFirewallGroup.Type == types.FirewallGroupTypeIpSet
}
