/*
 * Copyright 2024 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// IpSpace provides structured approach to allocating public and private IP addresses by preventing
// the use of overlapping IP addresses across organizations and organization VDCs.
//
// An IP space consists of a set of defined non-overlapping IP ranges and small CIDR blocks that are
// reserved and used during the consumption aspect of the IP space life cycle. An IP space can be
// either IPv4 or IPv6, but not both.
//
// Every IP space has an internal scope and an external scope. The internal scope of an IP space is
// a list of CIDR notations that defines the exact span of IP addresses in which all ranges and
// blocks must be contained in. The external scope defines the total span of IP addresses to which
// the IP space has access, for example the internet or a WAN.
type GenericIpSpace2 struct {
	IpSpace   *types.IpSpace
	vcdClient *VCDClient
}

// wrap is a hidden helper that helps to facilitate usage of generic CRUD function
func (g GenericIpSpace2) wrap(inner *types.IpSpace) *GenericIpSpace2 {
	// TODO TODO TODO note - a copy of the structure happens because it is value receiver
	g.IpSpace = inner
	return &g
}

// CreateIpSpace creates IP Space with desired configuration
func (vcdClient *VCDClient) GenCreateIpSpace(ipSpaceConfig *types.IpSpace) (*GenericIpSpace2, error) {
	c := crudConfig{
		endpoint:   types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		entityName: "IP Space",
	}
	outerType := GenericIpSpace2{vcdClient: vcdClient}
	return genericCreateEntity(&vcdClient.Client, outerType, c, ipSpaceConfig)
}

// GetAllIpSpaceSummaries retrieve summaries of all IP Spaces with an optional filter
// Note. There is no API endpoint to get multiple IP Spaces with their full definitions. Only
// "summaries" endpoint exists, but it does not include all fields. To retrieve complete structure
// one can use `GetIpSpaceById` or `GetIpSpaceByName`
func (vcdClient *VCDClient) GenGetAllIpSpaceSummaries(queryParameters url.Values) ([]*GenericIpSpace2, error) {
	c := crudConfig{
		endpoint:        types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaceSummaries,
		entityName:      "IP Space",
		queryParameters: queryParameters,
	}

	initializedParentType := GenericIpSpace2{vcdClient: vcdClient}
	return genericGetAllEntities[GenericIpSpace2, types.IpSpace](&vcdClient.Client, initializedParentType, c)
}

// GetIpSpaceById retrieves IP Space with a given ID
func (vcdClient *VCDClient) GenGetIpSpaceById(id string) (*GenericIpSpace2, error) {
	if id == "" { // TODO - `genericCrudConfig` or `genericGetSingleBareEntity` could do such validation?
		return nil, fmt.Errorf("IP Space ID cannot be empty")
	}

	c := crudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		endpointParams: []string{id},
		entityName:     "IP Space",
	}

	outerType := GenericIpSpace2{vcdClient: vcdClient}
	return genericGetSingleEntity[GenericIpSpace2, types.IpSpace](&vcdClient.Client, outerType, c)
}

// GetIpSpaceByName retrieves IP Space with a given name
// Note. It will return an error if multiple IP Spaces exist with the same name
func (vcdClient *VCDClient) GenGetIpSpaceByName(name string) (*GenericIpSpace2, error) {
	if name == "" {
		return nil, fmt.Errorf("IP Space lookup requires name")
	}

	queryParams := url.Values{}
	queryParams.Add("filter", "name=="+name)

	filteredEntities, err := vcdClient.GenGetAllIpSpaceSummaries(queryParams)
	if err != nil {
		return nil, err
	}

	singleIpSpace, err := oneOrError("name", name, filteredEntities)
	if err != nil {
		return nil, err
	}

	return vcdClient.GenGetIpSpaceById(singleIpSpace.IpSpace.ID)
}

// GetIpSpaceByNameAndOrgId retrieves IP Space with a given name in a particular Org
// Note. Only PRIVATE IP spaces belong to Orgs
func (vcdClient *VCDClient) GenGetIpSpaceByNameAndOrgId(name, orgId string) (*GenericIpSpace2, error) {
	if name == "" || orgId == "" {
		return nil, fmt.Errorf("IP Space lookup requires name and Org ID")
	}

	queryParams := url.Values{}
	queryParams.Add("filter", "name=="+name)
	queryParams = queryParameterFilterAnd("orgRef.id=="+orgId, queryParams)

	filteredEntities, err := vcdClient.GenGetAllIpSpaceSummaries(queryParams)
	if err != nil {
		return nil, err
	}

	singleIpSpace, err := oneOrError("name", name, filteredEntities)
	if err != nil {
		return nil, err
	}

	return vcdClient.GenGetIpSpaceById(singleIpSpace.IpSpace.ID)
}

// Update updates IP Space with new config
func (ipSpace *GenericIpSpace2) Update(ipSpaceConfig *types.IpSpace) (*GenericIpSpace2, error) {
	if ipSpaceConfig.ID == "" { // TODO - `genericCrudConfig` or `genericGetSingleBareEntity` could do such validation?
		return nil, fmt.Errorf("cannot update IP Space without ID")
	}

	c := crudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		endpointParams: []string{ipSpaceConfig.ID},
		entityName:     "IP Space",
	}
	outerType := GenericIpSpace2{vcdClient: ipSpace.vcdClient}
	return genericUpdateEntity(&ipSpace.vcdClient.Client, outerType, c, ipSpaceConfig)
}

// Delete deletes IP Space
func (ipSpace *GenericIpSpace2) Delete() error {
	c := crudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		endpointParams: []string{ipSpace.IpSpace.ID},
		entityName:     "IP Space",
	}
	return deleteById(&ipSpace.vcdClient.Client, c)
}
