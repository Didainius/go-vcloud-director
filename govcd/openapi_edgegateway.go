package govcd

import (
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"net/url"
	"time"
)

type NsxtEdgeGateway struct {
	EdgeGateway *types.CloudAPIEdgeGateway
	client      *Client
}

func (vdc *Vdc) GetOpenApiEdgeGatewayById(id string) (*NsxtEdgeGateway, error) {
	endpoint := "1.0.0/edgeGateways/"
	minimumApiVersion, err := vdc.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}


	time.Sleep(10 * time.Second)

	if id == "" {
		return nil, fmt.Errorf("empty edge gateway id")
	}

	urlRef, err := vdc.client.OpenApiBuildEndpoint(endpoint, id)
	if err != nil {
		return nil, err
	}

	egw := &NsxtEdgeGateway{
		EdgeGateway:   &types.CloudAPIEdgeGateway{},
		client: vdc.client,
	}

	err = vdc.client.OpenApiGetItem(minimumApiVersion, urlRef, nil, egw.EdgeGateway)
	if err != nil {
		return nil, err
	}

	return egw, nil
}

func (vdc *Vdc) GetAllOpenApiEdgeGateways(queryParameters url.Values) ([]*types.CloudAPIEdgeGateway, error) {
	endpoint := "1.0.0/edgeGateways/"
	minimumApiVersion, err := vdc.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := vdc.client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	responses := []*types.CloudAPIEdgeGateway{{}}

	err = vdc.client.OpenApiGetAllItems(minimumApiVersion, urlRef, queryParameters, &responses)
	if err != nil {
		return nil, err
	}

	return responses, nil
}


func (egw *NsxtEdgeGateway) Create(e *types.CloudAPIEdgeGateway) (*NsxtEdgeGateway, error) {
	endpoint := "1.0.0/edgeGateways"
	minimumApiVersion, err := egw.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := egw.client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	returnEgw := &NsxtEdgeGateway{
		EdgeGateway:   &types.CloudAPIEdgeGateway{},
		client: egw.client,
	}

	err = egw.client.OpenApiPostItem(minimumApiVersion, urlRef, nil, e, returnEgw.EdgeGateway)
	if err != nil {
		return nil, fmt.Errorf("error creating edge gateway: %s", err)
	}

	return returnEgw, nil
}

func (egw *NsxtEdgeGateway) Update(e *types.CloudAPIEdgeGateway) (*NsxtEdgeGateway, error) {
	endpoint := "1.0.0/edgeGateways/"
	minimumApiVersion, err := egw.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	if e.ID == "" {
		return nil, fmt.Errorf("cannot update edge gateway without id")
	}

	urlRef, err := egw.client.OpenApiBuildEndpoint(endpoint, e.ID)
	if err != nil {
		return nil, err
	}

	returnEgw := &NsxtEdgeGateway{
		EdgeGateway:   &types.CloudAPIEdgeGateway{},
		client: egw.client,
	}

	err = egw.client.OpenApiPutItem(minimumApiVersion, urlRef, nil, e, returnEgw.EdgeGateway)
	if err != nil {
		return nil, fmt.Errorf("error updating edge gateway: %s", err)
	}

	return returnEgw, nil
}

func (egw *NsxtEdgeGateway) Delete() error {
	endpoint := "1.0.0/edgeGateways/"
	minimumApiVersion, err := egw.client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return err
	}

	if egw.EdgeGateway.ID == "" {
		return fmt.Errorf("cannot delete Edge Gateway without id")
	}

	urlRef, err := egw.client.OpenApiBuildEndpoint(endpoint, egw.EdgeGateway.ID)
	if err != nil {
		return err
	}

	err = egw.client.OpenApiDeleteItem(minimumApiVersion, urlRef, nil)

	if err != nil {
		return fmt.Errorf("error deleting role: %s", err)
	}

	return nil
}
