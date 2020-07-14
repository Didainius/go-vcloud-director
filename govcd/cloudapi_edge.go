package govcd

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

type CloudAPIEdgeGateway struct {
	EdgeGateway *types.CloudAPIEdgeGateway
	client      *Client
}

func (vdc *Vdc) GetCloudAPIEdgeGatewayId(id string) (*CloudAPIEdgeGateway, error) {
	if id == "" {
		return nil, fmt.Errorf("empty edge gateway id")
	}

	edge := &CloudAPIEdgeGateway{
		EdgeGateway: new(types.CloudAPIEdgeGateway),
		client:      vdc.client,
	}

	urlString := vdc.client.VCDHREF.Scheme + "://" + vdc.client.VCDHREF.Host + "/cloudapi/1.0.0/edgeGateways/" + id
	url, _ := url.ParseRequestURI(urlString)
	// acceptMime := types.JSONMime + ";version=" + vdc.client.APIVersion
	// acceptMime := types.JSONMime + ";version=34.0"

	err := vdc.client.cloudApiGetItem(url, nil, edge.EdgeGateway)

	if err != nil {
		return nil, err
	}
	return edge, nil
}

func (vdc *Vdc) GetCloudAPIEdgeGateways(queryParameters url.Values) ([]*types.CloudAPIEdgeGateway, error) {
	urlString := vdc.client.VCDHREF.Scheme + "://" + vdc.client.VCDHREF.Host + "/cloudapi/1.0.0/edgeGateways"
	url, _ := url.ParseRequestURI(urlString)

	response := make([]*types.CloudAPIEdgeGateway, 1)

	// err := vdc.client.CloudApiGetAllItems(url, queryParameters, "error getting edge gateways %s", nil, &response)
	err := vdc.client.CloudApiGetAllItems(url, queryParameters, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (egw *CloudAPIEdgeGateway) Create(e *types.CloudAPIEdgeGateway) (*CloudAPIEdgeGateway, error) {
	urlString := egw.client.VCDHREF.Scheme + "://" + egw.client.VCDHREF.Host + "/cloudapi/1.0.0/edgeGateways"
	url, _ := url.ParseRequestURI(urlString)

	returnEdge := &CloudAPIEdgeGateway{
		EdgeGateway: &types.CloudAPIEdgeGateway{},
		client:      egw.client,
	}

	err := egw.client.cloudApiPostItem(url, nil, e, returnEdge.EdgeGateway)

	if err != nil {
		return nil, fmt.Errorf("error creating edge gateway: %s", err)
	}

	return returnEdge, nil
}

func (egw *CloudAPIEdgeGateway) Update(e *types.CloudAPIEdgeGateway) (*CloudAPIEdgeGateway, error) {
	urlString := egw.client.VCDHREF.Scheme + "://" + egw.client.VCDHREF.Host + "/cloudapi/1.0.0/edgeGateways/" + egw.EdgeGateway.ID
	url, _ := url.ParseRequestURI(urlString)

	returnEdge := &CloudAPIEdgeGateway{
		EdgeGateway: new(types.CloudAPIEdgeGateway),
		client:      egw.client,
	}

	err := egw.client.cloudApiPutItem(url, nil, e, returnEdge.EdgeGateway)

	if err != nil {
		return nil, fmt.Errorf("error creating updating gateway: %s", err)
	}

	return returnEdge, nil
}

func (egw *CloudAPIEdgeGateway) Delete() error {
	urlString := egw.client.VCDHREF.Scheme + "://" + egw.client.VCDHREF.Host + "/cloudapi/1.0.0/edgeGateways/" + egw.EdgeGateway.ID
	url, _ := url.ParseRequestURI(urlString)

	err := egw.client.cloudApiDeleteItem(url, nil)

	if err != nil {
		return fmt.Errorf("error deleting edge gateway: %s", err)
	}

	return nil
}
