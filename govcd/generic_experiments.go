/*
 * Copyright 2021 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

//	func NewNsxtAlbController(albCtrl *types.NsxtAlbController, vcdClient *VCDClient) *NsxtAlbController {
//		return &NsxtAlbController{
//			NsxtAlbController: albCtrl,
//			vcdClient:         vcdClient,
//		}
//	}
type Querier[T any] struct {
	client *Client
}

func NewQuerier[T any](c *Client) *Querier[T] {
	return &Querier[T]{
		client: c,
	}
}

func newContainer[T, Y, Z any]() *Z {
	return new(Z)
}

// T - main type (container)
// Y - inside type
func genericConstructor[T, Y any](typesField Y, vcdClient *VCDClient, client *Client) T {
	container := new(T)

	err := SetField(container, "types.NsxtAlbController", &typesField)
	fmt.Printf("err: %s", err)

	if vcdClient != nil {
		err := SetField(container, "VCDClient", &typesField)
		fmt.Printf("err: %s", err)
	}

	return *container
}

func SetField(item interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}
	// It's possible we can cache this, which is why precompute all these ahead of time.
	findJsonName := func(t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("vcd"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", fmt.Errorf("tag %s provided does not define a vcd tag", fieldName)
	}
	fieldNames := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		jname, _ := findJsonName(tag)
		fieldNames[jname] = i
	}

	fieldNum, ok := fieldNames[fieldName]
	if !ok {
		return fmt.Errorf("field %s does not exist within the provided item", fieldName)
	}
	fieldVal := v.Field(fieldNum)
	if fieldVal.CanSet() { // Exported fields
		fieldVal.Set(reflect.ValueOf(value))
	} else { // attempt setting unexported fields using unsafe

	}
	return nil
}

// func testInPackage() {

// 	vcdClient := &VCDClient{
// 		Client: Client{
// 			APIVersion: "2000",
// 		},
// 	}

// 	internalField := types.NsxtAlbController{ID: "testing-id"}

// 	asd := genericConstructor[NsxtAlbController, types.NsxtAlbController](internalField, vcdClient, nil)
// 	spew.Dump(asd)
// }

// Experiment 2

type NsxtAlbControllerExp2 struct {
	NsxtAlbController *types.NsxtAlbController `vcd:"types.NsxtAlbController" json:"nsxt_alb_controller,omitempty"`
	vcdClient         *VCDClient               `vcd:"VCDClient"`
}

func (t NsxtAlbControllerExp2) New(ct *types.NsxtAlbController, vcdClient *VCDClient, client *Client) *NsxtAlbControllerExp2 {
	t.vcdClient = vcdClient
	t.NsxtAlbController = ct

	return &t
}

// func [T any] New(tt *T) {

// }

// Experiment 3

// https://stackoverflow.com/questions/71801376/generics-constraint-to-a-type-with-a-function-returning-itself

type Builder[F, E any] interface {
	SetFoo(F) E
}

func demo[E Builder[F, E], F any](bldr E, foo F) E {
	return bldr.SetFoo(foo)
}

type AnyConstructor[T any, X any, Y any, Z any] interface {
	New(X, Y, Z) *T
}

func genericNew[T, X, Y, Z any](constr AnyConstructor[T, X, Y, Z], child X, vvv Y, ccc Z) *T {
	ret := new(T)

	ret = constr.New(child, vvv, ccc)

	return ret

}

func genericNew22[T AnyConstructor[T, X, Y, Z], X, Y, Z any](child X, vvv Y, ccc Z) *T {
	aaa := new(T)
	return T.New(*aaa, child, vvv, ccc)
}

// func genConstructorExp2[T any, Y any](a Constructor[T, Y]) *T {
// 	ret := new(T)

// 	return ret
// }

type NsxtAlbControllerExp3 struct {
	NsxtAlbController *types.NsxtAlbController `vcd:"types.NsxtAlbController" json:"nsxt_alb_controller,omitempty"`
	vcdClient         *VCDClient               `vcd:"VCDClient"`
}

func (t NsxtAlbControllerExp3) New(ct *types.NsxtAlbController) *NsxtAlbControllerExp3 {
	// t.vcdClient = vcdClient
	t.NsxtAlbController = ct

	return &t
}

type tinyConstructor[P any, C any] interface {
	New(*C) *P
}

func wrappedResponse[P tinyConstructor[P, C], C any](child *C) *P {
	ppp := new(P)
	return P.New(*ppp, child)
}

func wrappedResponses[P tinyConstructor[P, C], C any](children []*C) []*P {
	ppp := new(P)

	res := make([]*P, len(children))

	for index, value := range children {
		res[index] = P.New(*ppp, value)
	}

	return res
}

// wrappedResponses := make([]*NsxtAlbController, len(typeResponses))
// for sliceIndex := range typeResponses {
// 	wrappedResponses[sliceIndex] = &NsxtAlbController{
// 		NsxtAlbController: typeResponses[sliceIndex],
// 		vcdClient:         vcdClient,
// 	}
// // }
