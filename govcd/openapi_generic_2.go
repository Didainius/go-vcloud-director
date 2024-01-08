package govcd

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
func genericCreateEntity[P AnyConstructor[P, X, Y, Z], X, Y, Z any](client *Client, entityConfig *X, c genericCrudConfig, gVcdClient Y, gClient Z) (*P, error) {
	// c := genericCrudConfig{
	// 	endpoint:   types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSpaces,
	// 	entityName: "IP Space",
	// }

	createdEntity, err := genericCreateBareEntity(client, entityConfig, c)
	if err != nil {
		return nil, err
	}

	//P AnyConstructor[P, X, Y, Z], X, Y, Z any
	// Response wrapper
	wrappedEntry := genericWrappedResponse[P, X, Y, Z](*createdEntity, gVcdClient, gClient)

	return wrappedEntry, nil

	// Dublis 2
	// ppp := new(P)

	// // vcdClient := VCDClient{}

	// y := new(Y)
	// z := new(Z)

	// return P.New(*ppp, *createdEntity, *y, *z), nil
	// return wrappedEntry, nil
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
