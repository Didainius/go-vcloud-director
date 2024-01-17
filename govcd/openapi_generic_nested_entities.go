package govcd

import "fmt"

// Generic type explanations
// Common generic parameter names seen in this code
// P - Parent. The parent "container" type that is not in the `types` package. E.g. 'IpSpace'
// C - Child. The child of parent - the data that is being marshalled unmarshaled residing in `types` package. E.g. `types.IpSpace`

type genericInitializerType2[P any, C any] interface {
	initialize(child *C) *P
}

func genericInitializerCreateEntity[P CustomConstructor[P, C], C any](client *Client, entityConfig *C, c genericCrudConfig, i genericInitializerType2[P, C]) (*P, error) {
	if entityConfig == nil {
		return nil, fmt.Errorf("entity config '%s' cannot be empty for create operation", c.entityName)
	}

	createdEntity, err := genericCreateBareEntity(client, entityConfig, c)
	if err != nil {
		return nil, err
	}

	return i.initialize(createdEntity), nil
}

func genericInitializerUpdateEntity[P CustomConstructor[P, C], C any](client *Client, entityConfig *C, c genericCrudConfig, i genericInitializerType2[P, C]) (*P, error) {
	if entityConfig == nil {
		return nil, fmt.Errorf("entity config '%s' cannot be empty for update operation", c.entityName)
	}

	createdEntity, err := genericUpdateBareEntity(client, entityConfig, c)
	if err != nil {
		return nil, err
	}

	return i.initialize(createdEntity), nil
}

func genericGetSingleEntity[P CustomConstructor[P, C], C any](client *Client, c genericCrudConfig, i genericInitializerType2[P, C]) (*P, error) {
	retrievedEntity, err := genericGetSingleBareEntity[C](client, c)
	if err != nil {
		return nil, err
	}

	return i.initialize(retrievedEntity), nil
}

func genericGetAllEntities[P CustomConstructor[P, C], C any](client *Client, c genericCrudConfig, i genericInitializerType2[P, C]) ([]*P, error) {
	retrievedEntities, err := genericGetAllBareFilteredEntities[C](client, c)
	if err != nil {
		return nil, err
	}

	/// TODO - double check if there are no issues to call initialize each time on the same entry
	wrappedResults := make([]*P, len(retrievedEntities))
	for index, singleChildEntity := range retrievedEntities {
		wrappedResults[index] = i.initialize(singleChildEntity)
	}

	return wrappedResults, nil
}
