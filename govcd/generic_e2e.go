/*
 * Copyright 2024 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// type GenericNsxtFirewall struct {
// 	NsxtFirewallRuleContainer *types.NsxtFirewallRuleContainer
// 	client                    *Client
// 	// edgeGatewayId is stored for usage in NsxtFirewall receiver functions
// 	edgeGatewayId string
// }

type GenericIpSpace struct {
	IpSpace   *types.IpSpace
	vcdClient *VCDClient
}

func (g *GenericIpSpace) initialize(child *types.IpSpace) *GenericIpSpace {
	g.IpSpace = child
	return g
}

func genericInitializerSquared(vcdClient *VCDClient) *GenericIpSpace {
	return &GenericIpSpace{vcdClient: vcdClient}
}

func genericIpSpaceInitializer222(vcdClient *VCDClient) genericInitializerType[GenericIpSpace, types.IpSpace] {
	// genericInitializerType[P any, C any] func(child *C) *P
	retFunc := func(c *types.IpSpace) *GenericIpSpace {
		return &GenericIpSpace{
			IpSpace:   c,
			vcdClient: vcdClient,
		}
	}

	return retFunc
}

func genericIpSpaceInitializer(vcdClient *VCDClient) genericInitializerType[GenericIpSpace, types.IpSpace] {
	// genericInitializerType[P any, C any] func(child *C) *P
	retFunc := func(c *types.IpSpace) *GenericIpSpace {
		return &GenericIpSpace{
			IpSpace:   c,
			vcdClient: vcdClient,
		}
	}

	return retFunc
}

func (t GenericIpSpace) New(ct *types.IpSpace, vcdClient *VCDClient, client *Client) *GenericIpSpace {
	t.IpSpace = ct
	t.vcdClient = vcdClient

	return &t
}

func (t GenericIpSpace) New2(ct *types.IpSpace) *GenericIpSpace {
	t.IpSpace = ct
	// t.vcdClient = vcdClient

	return &t
}

/* func genericInitializer[P any]() *P {
	p := new(P)
	return p
} */

func NewGenericIpSpace(vcdClient *VCDClient) *GenericIpSpace {
	internalTypeField := &types.IpSpace{ID: "one"}

	return genericNew22[GenericIpSpace, *types.IpSpace, *VCDClient, *Client](internalTypeField, vcdClient, &vcdClient.Client)
}

// CreateIpSpace creates IP Space with desired configuration
func (vcdClient *VCDClient) GenericCreateIpSpace(ipSpaceConfig *types.IpSpace) (*GenericIpSpace, error) {
	c := genericCrudConfig{
		endpoint:   types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		entityName: "IP Space",
	}

	createdEntity, err := genericCreateBareEntity(&vcdClient.Client, ipSpaceConfig, c)
	if err != nil {
		return nil, err
	}

	// Response wrapper
	wrappedEntry := genericWrappedResponse[GenericIpSpace, *types.IpSpace, *VCDClient, *Client](createdEntity, vcdClient, &vcdClient.Client)
	return wrappedEntry, nil
}

// CreateIpSpace creates IP Space with desired configuration
func (vcdClient *VCDClient) GenericCreateIpSpace2(ipSpaceConfig *types.IpSpace) (*GenericIpSpace, error) {
	c := genericCrudConfig{
		endpoint:   types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		entityName: "IP Space",
	}

	return genericCreateEntityVcdClient[GenericIpSpace, *types.IpSpace](vcdClient, &ipSpaceConfig, c)
}

func (vcdClient *VCDClient) GenericCreateIpSpace3(ipSpaceConfig *types.IpSpace) (*GenericIpSpace, error) {
	c := genericCrudConfig{
		endpoint:   types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		entityName: "IP Space",
	}

	// genericInitializerType[P any, C any] func(child *C) *P

	// initializer := genericIpSpaceInitializer(vcdClient)
	// g := new(GenericIpSpace)
	// // initializer := g.initialize(vcdClient)

	gggg := &GenericIpSpace{vcdClient: vcdClient}

	intermediate, err := genericInitializerCreateEntity[GenericIpSpace, types.IpSpace](&vcdClient.Client, ipSpaceConfig, c, gggg)

	return intermediate, err
	// return nil, nil
}

// GetIpSpaceById retrieves IP Space with a given ID
func (vcdClient *VCDClient) GenericGetIpSpaceById(id string) (*GenericIpSpace, error) {

	if id == "" { // TODO - `genericCrudConfig` or `genericGetSingleBareEntity` could do such validation?
		return nil, fmt.Errorf("empty NSX-T Segment Profile Template ID")
	}

	c := genericCrudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		endpointParams: []string{id},
		entityName:     "IP Space",
	}
	retrievedEntity, err := genericGetSingleBareEntity[types.IpSpace](&vcdClient.Client, c)
	if err != nil {
		return nil, err
	}

	wrappedEntry := genericWrappedResponse[GenericIpSpace, *types.IpSpace, *VCDClient, *Client](retrievedEntity, vcdClient, &vcdClient.Client)
	return wrappedEntry, nil
}

func (ipSpace *GenericIpSpace) Update(ipSpaceConfig *types.IpSpace) (*GenericIpSpace, error) {
	if ipSpaceConfig.ID == "" { // TODO - `genericCrudConfig` or `genericGetSingleBareEntity` could do such validation?
		return nil, fmt.Errorf("cannot update IP Space without ID")
	}

	c := genericCrudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		endpointParams: []string{ipSpaceConfig.ID},
		entityName:     "IP Space",
	}
	updatedEntity, err := genericUpdateBareEntity(&ipSpace.vcdClient.Client, ipSpaceConfig, c)
	if err != nil {
		return nil, err
	}

	wrappedEntry := genericWrappedResponse[GenericIpSpace, *types.IpSpace, *VCDClient, *Client](updatedEntity, ipSpace.vcdClient, &ipSpace.vcdClient.Client)
	return wrappedEntry, nil
}

func (ipSpace *GenericIpSpace) Delete() error {
	c := genericCrudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
		endpointParams: []string{ipSpace.IpSpace.ID},
		entityName:     "IP Space",
	}
	return deleteById(&ipSpace.vcdClient.Client, c)
}
