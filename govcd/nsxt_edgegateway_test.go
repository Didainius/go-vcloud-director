// +build network functional openapi ALL

package govcd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

/*
func (vcd *TestVCD) Test_NsxtEdge(check *C) {
	// if vcd.config.VCD.EdgeGateway == "" {
	// 	check.Skip("Skipping test because no edge gateway given")
	// }
	edge, err := vcd.vdc.GetNsxtEdgeGatewayById("urn:vcloud:gateway:e934e76b-757c-4ea7-a89a-c160a6199689")
	check.Assert(err, IsNil)
	check.Assert(edge.EdgeGateway.Name, Equals, "nsxt-gw-dainius")
	// copyEdge := edge
	// err = edge.Refresh()
	// check.Assert(err, IsNil)
	// check.Assert(copyEdge.EdgeGateway.Name, Equals, edge.EdgeGateway.Name)
	// check.Assert(copyEdge.EdgeGateway.HREF, Equals, edge.EdgeGateway.HREF)
}
*/
func (vcd *TestVCD) Test_NsxtEdgeCreate(check *C) {
	// if vcd.config.VCD.EdgeGateway == "" {
	// 	check.Skip("Skipping test because no edge gateway given")
	// }

	// Lookup data. Get Org Vdc

	nsxtVdcName := "vdc-dainius-nsxt"

	adminOrg, err := vcd.client.GetAdminOrgByName(vcd.config.VCD.Org)
	check.Assert(err, IsNil)

	nsxtVdc, err := adminOrg.GetVDCByName(nsxtVdcName, false)
	check.Assert(err, IsNil)

	bd := new(types.NsxtEdgeGateway)

	data := []byte(`{
  "name": "nsx-t-edge",
  "description": "nsx-t-edge-description",
  "orgVdc": {
    "id": "` + nsxtVdc.Vdc.ID + `"
  },
  "edgeGatewayUplinks": [
    {
      "uplinkId": "urn:vcloud:network:3f1f6081-c151-412d-8933-6f25bd896032",
      "subnets": {
        "values": [
          {
            "gateway": "1.1.1.1",
            "prefixLength": 24,
            "dnsSuffix": null,
            "dnsServer1": "",
            "dnsServer2": "",
            "ipRanges": {
              "values": [
                {
                  "startAddress": "1.1.1.10",
                  "endAddress": "1.1.1.15"
                }
              ]
            },
            "enabled": true
          }
        ]
      },
      "dedicated": false
    }
  ]
}
`)

	fmt.Println(string(data))

	err = json.Unmarshal(data, bd)
	check.Assert(err, IsNil)

	// newEdge := &NsxtEdgeGateway{
	// 	EdgeGateway: bd,
	// 	client:      vcd.vdc.client,
	// }

	// spew.Dump(newEdge.EdgeGateway)
	// GetAdminOrgByName()

	createdEdge, err := adminOrg.CreateNsxtEdgeGateway(bd)

	check.Assert(err, IsNil)
	check.Assert(createdEdge.EdgeGateway.Name, Equals, bd.Name)

	// Check pagination stuff

	createdEdge.EdgeGateway.Name = "renamed-edge"
	updatedEdge, err := createdEdge.Update(createdEdge.EdgeGateway)
	check.Assert(err, IsNil)
	check.Assert(updatedEdge.EdgeGateway.Name, Equals, "renamed-edge")

	// FIQL filtering test
	queryParams := url.Values{}
	queryParams.Add("filter", "name==renamed-edge")
	//
	egws, err := adminOrg.GetAllNsxtEdgeGateways(queryParams)
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

/*
func (vcd *TestVCD) Test_NsxtEdgeGetPages(check *C) {
	// if vcd.config.VCD.EdgeGateway == "" {
	// 	check.Skip("Skipping test because no edge gateway given")
	// }

	params := url.Values{}
	params.Add("pageSize", "1")

	_, err := vcd.vdc.GetAllNsxtEdgeGateways(params)
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

// func (vdc *Vdc) GetAllNsxtEdgeGateways(queryParameters url.Values) ([]*types.NsxtEdgeGateway, error) {
// 	urlString := vdc.client.VCDHREF.Scheme + "://" + vdc.client.VCDHREF.Host + "/cloudapi/1.0.0/edgeGateways"
// 	url, _ := url.ParseRequestURI(urlString)
//
// 	response := make([]*types.NsxtEdgeGateway, 1)
//
// 	// err := vdc.client.OpenApiGetAllItems(url, queryParameters, "error getting edge gateways %s", nil, &response)
// 	err := vdc.client.OpenApiGetAllItems(url, queryParameters, "error getting edge gateways %s", nil, &response)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return response, nil
// }
*/
