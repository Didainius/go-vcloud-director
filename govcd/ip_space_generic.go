/*
 * Copyright 2024 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

type GenericIpSpace2 struct {
	IpSpace   *types.IpSpace
	vcdClient *VCDClient
}

// initialize is a hidden helper that helps to facilitate generic components
// It should fill all parent type (GenericIpSpace2) fields, except the "child" entity type that
func (g *GenericIpSpace2) initialize(child *types.IpSpace) *GenericIpSpace2 {
	g.IpSpace = child
	return g
}

func (vcdClient *VCDClient) GenCreateIpSpace(ipSpaceConfig *types.IpSpace) (*GenericIpSpace2, error) {
	c := genericCrudConfig{
		endpoint:   types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		entityName: "IP Space",
	}
	initializedParentType := &GenericIpSpace2{vcdClient: vcdClient}
	return genericInitializerCreateEntity[GenericIpSpace2, types.IpSpace](&vcdClient.Client, ipSpaceConfig, c, initializedParentType)
}

func (vcdClient *VCDClient) GenGetAllIpSpaceSummaries(queryParameters url.Values) ([]*GenericIpSpace2, error) {
	c := genericCrudConfig{
		endpoint:        types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaceSummaries,
		entityName:      "IP Space",
		queryParameters: queryParameters,
	}

	initializedParentType := &GenericIpSpace2{vcdClient: vcdClient}
	return genericGetAllEntities[GenericIpSpace2, types.IpSpace](&vcdClient.Client, c, initializedParentType)
}

func (vcdClient *VCDClient) GenGetIpSpaceById(id string) (*GenericIpSpace2, error) {
	if id == "" { // TODO - `genericCrudConfig` or `genericGetSingleBareEntity` could do such validation?
		return nil, fmt.Errorf("IP Space ID cannot be empty")
	}

	c := genericCrudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		endpointParams: []string{id},
		entityName:     "IP Space",
	}

	initializedParentType := &GenericIpSpace2{vcdClient: vcdClient}
	return genericGetSingleEntity[GenericIpSpace2, types.IpSpace](&vcdClient.Client, c, initializedParentType)
}

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

func (ipSpace *GenericIpSpace2) Update(ipSpaceConfig *types.IpSpace) (*GenericIpSpace2, error) {
	// if ipSpaceConfig.ID == "" { // TODO - `genericCrudConfig` or `genericGetSingleBareEntity` could do such validation?
	// 	return nil, fmt.Errorf("cannot update IP Space without ID")
	// }

	c := genericCrudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		endpointParams: []string{ipSpaceConfig.ID},
		entityName:     "IP Space",
	}
	initializedParentType := &GenericIpSpace2{vcdClient: ipSpace.vcdClient}
	return genericInitializerUpdateEntity[GenericIpSpace2, types.IpSpace](&ipSpace.vcdClient.Client, ipSpaceConfig, c, initializedParentType)
}

func (ipSpace *GenericIpSpace2) Delete() error {
	c := genericCrudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		endpointParams: []string{ipSpace.IpSpace.ID},
		entityName:     "IP Space",
	}
	return deleteById(&ipSpace.vcdClient.Client, c)
}
