/*
 * Copyright 2021 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"github.com/vmware/go-vcloud-director/v2/util"
)

type NsxtIpSecVpn struct {
	NsxtIpSecVpn *types.NsxtIpSecVpn
	client       *Client
	// edgeGatewayId is stored here so that pointer receiver functions can embed edge gateway ID into path
	edgeGatewayId string
}

// GetAllIpSecVpns returns all IP Sec VPN configurations
func (egw *NsxtEdgeGateway) GetAllIpSecVpns(queryParameters url.Values) ([]*NsxtIpSecVpn, error) {
	client := egw.client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSecVpn
	apiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(fmt.Sprintf(endpoint, egw.EdgeGateway.ID))
	if err != nil {
		return nil, err
	}

	typeResponses := []*types.NsxtIpSecVpn{{}}
	err = client.OpenApiGetAllItems(apiVersion, urlRef, queryParameters, &typeResponses)
	if err != nil {
		return nil, err
	}

	// Wrap all typeResponses into NsxtIpSecVpn types with client
	wrappedResponses := make([]*NsxtIpSecVpn, len(typeResponses))
	for sliceIndex := range typeResponses {
		wrappedResponses[sliceIndex] = &NsxtIpSecVpn{
			NsxtIpSecVpn:  typeResponses[sliceIndex],
			client:        client,
			edgeGatewayId: egw.EdgeGateway.ID,
		}
	}

	return wrappedResponses, nil
}

func (egw *NsxtEdgeGateway) GetIpSecVpnById(id string) (*NsxtIpSecVpn, error) {
	if id == "" {
		return nil, fmt.Errorf("canot find NSX-T IP Sec VPN configuration without ID")
	}

	client := egw.client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSecVpn
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(fmt.Sprintf(endpoint, egw.EdgeGateway.ID), id)
	if err != nil {
		return nil, err
	}

	returnObject := &NsxtIpSecVpn{
		NsxtIpSecVpn:  &types.NsxtIpSecVpn{},
		client:        client,
		edgeGatewayId: egw.EdgeGateway.ID,
	}

	err = client.OpenApiGetItem(minimumApiVersion, urlRef, nil, returnObject.NsxtIpSecVpn)
	if err != nil {
		return nil, err
	}

	return returnObject, nil
}

func (egw *NsxtEdgeGateway) GetIpSecVpnByName(name string) (*NsxtIpSecVpn, error) {
	if name == "" {
		return nil, fmt.Errorf("canot find NSX-T IP Sec VPN configuration without Name")
	}

	allVpns, err := egw.GetAllIpSecVpns(nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving all NSX-T IP Sec VPN configurations: %s", err)
	}

	var allResults []*NsxtIpSecVpn

	for _, vpnConfig := range allVpns {
		if vpnConfig.NsxtIpSecVpn.Name == name {
			allResults = append(allResults, vpnConfig)
		}
	}

	if len(allResults) > 1 {
		return nil, fmt.Errorf("error - found %d NSX-T IP Sec VPN configuratios with Name '%s'. Expected 1", len(allResults), name)
	}

	if len(allResults) == 0 {
		return nil, ErrorEntityNotFound
	}

	// Retrieving the object by ID, because only it includes Pre-shared Key
	return egw.GetIpSecVpnById(allResults[0].NsxtIpSecVpn.ID)
}

func (egw *NsxtEdgeGateway) CreateIpSecVpn(ipSecVpnConfig *types.NsxtIpSecVpn) (*NsxtIpSecVpn, error) {
	client := egw.client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSecVpn
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(fmt.Sprintf(endpoint, egw.EdgeGateway.ID))
	if err != nil {
		return nil, err
	}

	task, err := client.OpenApiPostItemAsync(minimumApiVersion, urlRef, nil, ipSecVpnConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating NSX-T IP Sec VPN configuration: %s", err)
	}

	err = task.WaitTaskCompletion()
	if err != nil {
		return nil, fmt.Errorf("task failed while creating NSX-T IP Sec VPN configuration: %s", err)
	}

	// filtering even by Name is not supported
	allVpns, err := egw.GetAllIpSecVpns(nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving all NSX-T IP Sec VPN configuration after creation: %s", err)
	}

	for index, singleConfig := range allVpns {
		if singleConfig.IsEqualTo(ipSecVpnConfig) {
			// retrieve exact value by ID, because only this endpoint includes private key
			ipSecVpn, err := egw.GetIpSecVpnById(allVpns[index].NsxtIpSecVpn.ID)
			if err != nil {
				return nil, fmt.Errorf("error retrieving NSX-T IP Sec VPN configuration: %s", err)
			}

			return ipSecVpn, nil
		}
	}

	return nil, fmt.Errorf("error finding NSX-T IP Sec VPN configuration after creation: %s", ErrorEntityNotFound)
}

func (ipSecVpn *NsxtIpSecVpn) Update(ipSecVpnConfig *types.NsxtIpSecVpn) (*NsxtIpSecVpn, error) {
	client := ipSecVpn.client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSecVpn
	apiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if ipSecVpn.NsxtIpSecVpn.ID == "" {
		return nil, fmt.Errorf("cannot update NSX-T IP Sec VPN configuration without ID")
	}

	urlRef, err := client.OpenApiBuildEndpoint(fmt.Sprintf(endpoint, ipSecVpn.edgeGatewayId), ipSecVpn.NsxtIpSecVpn.ID)
	if err != nil {
		return nil, err
	}

	returnObject := &NsxtIpSecVpn{
		NsxtIpSecVpn:  &types.NsxtIpSecVpn{},
		client:        client,
		edgeGatewayId: ipSecVpn.edgeGatewayId,
	}

	err = client.OpenApiPutItem(apiVersion, urlRef, nil, ipSecVpnConfig, returnObject.NsxtIpSecVpn)
	if err != nil {
		return nil, fmt.Errorf("error updating NSX-T IP Sec VPN configuration: %s", err)
	}

	return returnObject, nil
}

// Delete allows users to delete NSX-T Application Port Profile
func (ipSecVpn *NsxtIpSecVpn) Delete() error {
	client := ipSecVpn.client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointIpSecVpn
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return err
	}

	if ipSecVpn.NsxtIpSecVpn.ID == "" {
		return fmt.Errorf("cannot delete NSX-T IP Sec VPN configuration without ID")
	}

	urlRef, err := ipSecVpn.client.OpenApiBuildEndpoint(fmt.Sprintf(endpoint, ipSecVpn.edgeGatewayId), ipSecVpn.NsxtIpSecVpn.ID)
	if err != nil {
		return err
	}

	err = ipSecVpn.client.OpenApiDeleteItem(minimumApiVersion, urlRef, nil)

	if err != nil {
		return fmt.Errorf("error deleting NSX-T IP Sec VPN configuration: %s", err)
	}

	return nil
}

// IsEqualTo helps to find NSX-T IP Sec Configuration
// Combination of LocalAddress and RemoteAddress has to be Unique. This is a list of fields compared:
// * Name
// * Description
// * Enabled
// * LocalEndpoint.LocalAddress
// * RemoteEndpoint.RemoteAddress
func (ipSecVpn *NsxtIpSecVpn) IsEqualTo(vpnConfig *types.NsxtIpSecVpn) bool {
	return ipSetVpnRulesEqual(ipSecVpn.NsxtIpSecVpn, vpnConfig)
}

// ipSetVpnRulesEqual performs comparison of two rules to ease lookup. This is a list of fields compared:
//// * Name
//// * Description
//// * Enabled
//// * LocalEndpoint.LocalAddress
//// * RemoteEndpoint.RemoteAddress
func ipSetVpnRulesEqual(first, second *types.NsxtIpSecVpn) bool {
	util.Logger.Println("comparing NSX-T IP Sev VPN configuration:")
	util.Logger.Printf("%+v\n", first)
	util.Logger.Println("against:")
	util.Logger.Printf("%+v\n", second)

	// These fields should be enough to cover uniqueness
	if first.Name == second.Name &&
		first.Description == second.Description &&
		first.Enabled == second.Enabled &&
		first.LocalEndpoint.LocalAddress == second.LocalEndpoint.LocalAddress &&
		first.RemoteEndpoint.RemoteAddress == second.RemoteEndpoint.RemoteAddress {
		return true
	}

	return false
}
