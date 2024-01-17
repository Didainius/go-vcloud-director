package govcd

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

// layer 2 abstractions

// create
/*
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
return wrappedEntry, nil */

// P - parent
// C - child
// [P AnyConstructor[P, X, Y, Z], X, Y, Z any](child X, vvv Y, ccc Z) *P {
func genericCreateEntity[P AnyParentConstructorClientVcdClient[P, X, Y, Z], X, Y, Z any](client *Client, entityConfig *X, c genericCrudConfig, gVcdClient Y, gClient Z) (*P, error) {
	createdEntity, err := genericCreateBareEntity(client, entityConfig, c)
	if err != nil {
		return nil, err
	}

	fmt.Println("created entity")
	spew.Dump(createdEntity)

	// Works, but is more confusing
	wrappedEntry := genericWrappedResponse[P, X, Y, Z](*createdEntity, gVcdClient, gClient)

	// construcFunc := func() GenericIpSpace {
	// 	return P.New()
	// }

	// wrappedEntry := genericWrapper[GenericIpSpace](construcFunc)

	return wrappedEntry, nil
}

func genericCreateEntityVcdClient[P ParentConstructorVcdClient[P, C, *VCDClient], C any](vcdClient *VCDClient, entityConfig *C, c genericCrudConfig) (*P, error) {
	createdEntity, err := genericCreateBareEntity(&vcdClient.Client, entityConfig, c)
	if err != nil {
		return nil, err
	}

	fmt.Println("created entity")
	spew.Dump(createdEntity)

	// Works, but is more confusing
	wrappedEntry := genericWrappedResponseVcdClient[P, C](*createdEntity, vcdClient)

	// construcFunc := func() GenericIpSpace {
	// 	return P.New()
	// }

	// wrappedEntry := genericWrapper[GenericIpSpace](construcFunc)

	return wrappedEntry, nil
}

func genericWrappedResponseVcdClient[P ParentConstructorVcdClient[P, C, *VCDClient], C any](child C, vcdClient *VCDClient) *P {
	ppp := new(P)
	return P.New2(*ppp, child)
}

// func genericWrappedResponse[P AnyConstructor[P, X, Y, Z], X, Y, Z any](child X, vvv Y, ccc Z) *P {
// 	ppp := new(P)
// 	return P.New(*ppp, child, vvv, ccc)
// }

func genericWrapper[T any](wrapFunc func() T) *T {

	return nil
}

// func genericCreateBareEntity[T any](client *Client, entityConfig *T, c genericCrudConfig) (*T, error) {
// 	apiVersion, err := client.getOpenApiHighestElevatedVersion(c.endpoint)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting API version for creating entity '%s': %s", c.entityName, err)
// 	}

// 	exactEndpoint, err := urlFromEndpoint(c.endpoint, c.endpointParams)
// 	if err != nil {
// 		return nil, fmt.Errorf("error building endpoint '%s' with given params '%s' for entity '%s': %s", c.endpoint, strings.Join(c.endpointParams, ","), c.entityName, err)
// 	}

// 	urlRef, err := client.OpenApiBuildEndpoint(exactEndpoint)
// 	if err != nil {
// 		return nil, fmt.Errorf("error building API endpoint for entity '%s' creation: %s", c.entityName, err)
// 	}

// 	createdEntityConfig := new(T)
// 	err = client.OpenApiPostItem(apiVersion, urlRef, c.queryParameters, entityConfig, createdEntityConfig, c.additionalHeader)
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating entity of type '%s': %s", c.entityName, err)
// 	}

// 	return createdEntityConfig, nil
// }
