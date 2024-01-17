package govcd

import "fmt"

// Generic type explanations
// Common generic parameter names seen in this code
// O - Outer type that is in the `govcd` package. (e.g. 'IpSpace')
// I - Inner type the type that is being marshalled/unmarshalled (usually in `types` package. E.g. `types.IpSpace`)

type genericOuterConstructor[O any, I any] interface {
	wrap(child *I) *O
}

func genericCreateEntity[O genericOuterConstructor[O, I], I any](client *Client, outerEntity O, c crudConfig, innerConfig *I) (*O, error) {
	if innerConfig == nil {
		return nil, fmt.Errorf("entity config '%s' cannot be empty for create operation", c.entityName)
	}

	createdInnerEntity, err := genericCreateInnerEntity(client, c, innerConfig)
	if err != nil {
		return nil, err
	}

	return outerEntity.wrap(createdInnerEntity), nil
}

func genericUpdateEntity[O genericOuterConstructor[O, I], I any](client *Client, outerEntity O, c crudConfig, innerConfig *I) (*O, error) {
	if innerConfig == nil {
		return nil, fmt.Errorf("entity config '%s' cannot be empty for update operation", c.entityName)
	}

	updatedInnerEntity, err := genericUpdateInnerEntity(client, c, innerConfig)
	if err != nil {
		return nil, err
	}

	return outerEntity.wrap(updatedInnerEntity), nil
}

func genericGetSingleEntity[O genericOuterConstructor[O, I], I any](client *Client, outerEntity O, c crudConfig) (*O, error) {
	retrievedInnerEntity, err := genericGetInnerEntity[I](client, c)
	if err != nil {
		return nil, err
	}

	return outerEntity.wrap(retrievedInnerEntity), nil
}

func genericGetAllEntities[O genericOuterConstructor[O, I], I any](client *Client, outerEntity O, c crudConfig) ([]*O, error) {
	retrievedAllInnerEntities, err := genericGetAllInnerEntities[I](client, c)
	if err != nil {
		return nil, err
	}

	wrappedOuterEntities := make([]*O, len(retrievedAllInnerEntities))
	for index, singleInnerEntity := range retrievedAllInnerEntities {
		// outerEntity.wrap is a value receiver, therefore it creates a shallow copy for each call
		wrappedOuterEntities[index] = outerEntity.wrap(singleInnerEntity)
	}

	return wrappedOuterEntities, nil
}
