package govcd

import "fmt"

// Generic type explanations
// Common generic parameter names seen in this code
// P - Parent. The parent "container" type that is not in the `types` package. E.g. 'IpSpace'
// C - Child. The child of parent - the data that is being marshalled unmarshaled residing in `types` package. E.g. `types.IpSpace`

// type genericInitializerType2[P any, C any] interface {
// 	initialize(child *C) *P
// }

type GenericParentConstructor[P any, C any] interface {
	wrap(child *C) *P
}

func genericInitializerCreateEntity[P GenericParentConstructor[P, C], C any](client *Client, parentEntity P, c genericCrudConfig, entityConfig *C) (*P, error) {
	if entityConfig == nil {
		return nil, fmt.Errorf("entity config '%s' cannot be empty for create operation", c.entityName)
	}

	createdBareEntity, err := genericCreateBareEntity(client, c, entityConfig)
	if err != nil {
		return nil, err
	}

	return parentEntity.wrap(createdBareEntity), nil
}

func genericInitializerUpdateEntity[P GenericParentConstructor[P, C], C any](client *Client, parentEntity P, c genericCrudConfig, entityConfig *C) (*P, error) {
	if entityConfig == nil {
		return nil, fmt.Errorf("entity config '%s' cannot be empty for update operation", c.entityName)
	}

	updatedBareEntity, err := genericUpdateBareEntity(client, c, entityConfig)
	if err != nil {
		return nil, err
	}

	return parentEntity.wrap(updatedBareEntity), nil
}

func genericGetSingleEntity[P GenericParentConstructor[P, C], C any](client *Client, parentEntity P, c genericCrudConfig) (*P, error) {
	retrievedBareEntity, err := genericGetSingleBareEntity[C](client, c)
	if err != nil {
		return nil, err
	}

	return parentEntity.wrap(retrievedBareEntity), nil
}

func genericGetAllEntities[P GenericParentConstructor[P, C], C any](client *Client, parentEntity P, c genericCrudConfig) ([]*P, error) {
	retrievedAllBareEntities, err := genericGetAllBareFilteredEntities[C](client, c)
	if err != nil {
		return nil, err
	}

	/// TODO - double check if there are no issues to call initialize each time on the same entry
	wrappedResults := make([]*P, len(retrievedAllBareEntities))
	for index, singleChildEntity := range retrievedAllBareEntities {
		wrappedResults[index] = parentEntity.wrap(singleChildEntity)
	}

	return wrappedResults, nil
}
