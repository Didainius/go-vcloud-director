package govcd

import (
	"encoding/json"
	"fmt"
	"net/url"

	. "gopkg.in/check.v1"
)

// Test_CloudAPIAudiTrail uses low level Get function to test out that pagination really works. It is an example how to
// fetch response from multiple pages in RAW json messages without having defined a clear struct.
func (vcd *TestVCD) Test_CloudAPIAudiTrail(check *C) {
	urlString := vcd.vdc.client.VCDHREF.Scheme + "://" + vcd.vdc.client.VCDHREF.Host + "/cloudapi/1.0.0/auditTrail"
	url, _ := url.ParseRequestURI(urlString)

	// Prep response struct
	response := make([]*json.RawMessage, 1)
	err := vcd.vdc.client.cloudApiGetAllItems(url, nil, "error getting audittrail %s", nil, &response)

	check.Assert(err, IsNil)

	for _, v := range response {
		s, _ := v.MarshalJSON()
		fmt.Println(string(s))

	}

	fmt.Println(len(response))
}

// func (vdc *Vdc) GetCloudAPIEdgeGateways(queryParameters url.Values) ([]*types.CloudAPIEdgeGateway, error) {
// 	urlString := vdc.client.VCDHREF.Scheme + "://" + vdc.client.VCDHREF.Host + "/cloudapi/1.0.0/edgeGateways"
// 	url, _ := url.ParseRequestURI(urlString)
//
// 	response := make([]*types.CloudAPIEdgeGateway, 1)
//
// 	// err := vdc.client.cloudApiGetAllItems(url, queryParameters, "error getting edge gateways %s", nil, &response)
// 	err := vdc.client.cloudApiGetAllItems(url, queryParameters, "error getting edge gateways %s", nil, &response)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return response, nil
// }
