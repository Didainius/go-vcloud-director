package govcd

import (
	"encoding/json"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_NsxtEdge(check *C) {
	// if vcd.config.VCD.EdgeGateway == "" {
	// 	check.Skip("Skipping test because no edge gateway given")
	// }
	edge, err := vcd.vdc.GetCloudAPIEdgeGatewayId("urn:vcloud:gateway:e934e76b-757c-4ea7-a89a-c160a6199689")
	check.Assert(err, IsNil)
	check.Assert(edge.EdgeGateway.Name, Equals, "nsxt-gw-dainius")
	// copyEdge := edge
	// err = edge.Refresh()
	// check.Assert(err, IsNil)
	// check.Assert(copyEdge.EdgeGateway.Name, Equals, edge.EdgeGateway.Name)
	// check.Assert(copyEdge.EdgeGateway.HREF, Equals, edge.EdgeGateway.HREF)
}

func (vcd *TestVCD) Test_NsxtEdgeCreate(check *C) {
	// if vcd.config.VCD.EdgeGateway == "" {
	// 	check.Skip("Skipping test because no edge gateway given")
	// }

	bd := new(types.CloudAPIEdgeGateway)

	data := []byte(`{
  "name": "nsxt-edge",
  "description": "Testing-nsxt-edge",
  "orgVdc": {
    "id": "urn:vcloud:vdc:170089ee-5d34-462c-826e-c9b88b25e893"
  },
  "edgeGatewayUplinks": [{
    "uplinkId": "urn:vcloud:network:ac9890b0-106d-469e-b1d5-0848dda198d6",
    "uplinkName": "nsxt-extnet-dainius",
    "subnets": {
      "values": [{
        "gateway": "10.150.191.253",
        "prefixLength": 19,
        "dnsSuffix": null,
        "dnsServer1": "",
        "dnsServer2": "",
        "ipRanges": {
          "values": []
        },
        "enabled": true,
        "totalIpCount": 2,
        "usedIpCount": 0
      }]
    },
    "dedicated": false
  }]
}
`)

	err := json.Unmarshal(data, bd)

	newEdge := &CloudAPIEdgeGateway{
		EdgeGateway: bd,
		client:      vcd.vdc.client,
	}

	// spew.Dump(newEdge.EdgeGateway)

	createdEdge, err := newEdge.Create(newEdge.EdgeGateway)

	check.Assert(err, IsNil)
	check.Assert(createdEdge.EdgeGateway.Name, Equals, newEdge.EdgeGateway.Name)

	// Check pagination stuff

	createdEdge.EdgeGateway.Name = "renamed-edge"
	updatedEdge, err := createdEdge.Update(createdEdge.EdgeGateway)
	check.Assert(err, IsNil)
	check.Assert(updatedEdge.EdgeGateway.Name, Equals, "renamed-edge")

	// FIQL filtering test
	queryParams := url.Values{}
	queryParams.Add("filter", "name==renamed-edge")
	//
	egws, err := vcd.vdc.GetCloudAPIEdgeGateways(queryParams)
	check.Assert(err, IsNil)
	check.Assert(len(egws) == 1, Equals, true)

	// check.Assert(edge.EdgeGateway.Name, Equals, "nsxt-gw-dainius")

	err = updatedEdge.Delete()
	check.Assert(err, IsNil)

	// copyEdge := edge
	// err = edge.Refresh()
	// check.Assert(err, IsNil)
	// check.Assert(copyEdge.EdgeGateway.Name, Equals, edge.EdgeGateway.Name)
	// check.Assert(copyEdge.EdgeGateway.HREF, Equals, edge.EdgeGateway.HREF)
}

func (vcd *TestVCD) Test_NsxtEdgeGetPages(check *C) {
	// if vcd.config.VCD.EdgeGateway == "" {
	// 	check.Skip("Skipping test because no edge gateway given")
	// }

	params := url.Values{}
	params.Add("pageSize", "1")

	_, err := vcd.vdc.GetCloudAPIEdgeGateways(params)
	// spew.Dump(edges)
	check.Assert(err, IsNil)
	// check.Assert(edges, )

	// spew.Dump(edges)
	// check.Assert(edge.EdgeGateway.Name, Equals, "nsxt-gw-dainius")
	// copyEdge := edge
	// err = edge.Refresh()
	// check.Assert(err, IsNil)
	// check.Assert(copyEdge.EdgeGateway.Name, Equals, edge.EdgeGateway.Name)
	// check.Assert(copyEdge.EdgeGateway.HREF, Equals, edge.EdgeGateway.HREF)
}

// func (vdc *Vdc) GetCloudAPIEdgeGateways(queryParameters url.Values) ([]*types.CloudAPIEdgeGateway, error) {
// 	urlString := vdc.client.VCDHREF.Scheme + "://" + vdc.client.VCDHREF.Host + "/cloudapi/1.0.0/edgeGateways"
// 	url, _ := url.ParseRequestURI(urlString)
//
// 	response := make([]*types.CloudAPIEdgeGateway, 1)
//
// 	// err := vdc.client.OpenApiGetAllItems(url, queryParameters, "error getting edge gateways %s", nil, &response)
// 	err := vdc.client.OpenApiGetAllItems(url, queryParameters, "error getting edge gateways %s", nil, &response)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return response, nil
// }
