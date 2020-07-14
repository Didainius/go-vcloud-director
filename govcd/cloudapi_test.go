package govcd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/davecgh/go-spew/spew"

	. "gopkg.in/check.v1"
)

// Test_CloudAPIAudiTrail uses low level GET function to test out that pagination really works. It is an example how to
// fetch response from multiple pages in RAW json messages without having defined a clear struct.
func (vcd *TestVCD) Test_CloudAPIRawJsonAudiTrail(check *C) {
	urlRef, err := vcd.client.Client.BuildCloudAPIEndpoint("1.0.0/auditTrail")
	check.Assert(err, IsNil)

	responses := []json.RawMessage{}
	err = vcd.vdc.client.CloudApiGetAllItems(urlRef, nil, &responses)
	check.Assert(err, IsNil)

	check.Assert(len(responses) > 1, Equals, true)
}

func (vcd *TestVCD) Test_CloudAPIInlineStructAudiTrail(check *C) {
	urll, err := vcd.client.Client.BuildCloudAPIEndpoint("1.0.0/auditTrail")
	check.Assert(err, IsNil)

	// Inline type
	type AudiTrail struct {
		EventID      string `json:"eventId"`
		Description  string `json:"description"`
		OperatingOrg struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"operatingOrg"`
		User struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"user"`
		EventEntity struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"eventEntity"`
		TaskID               interface{} `json:"taskId"`
		TaskCellID           string      `json:"taskCellId"`
		CellID               string      `json:"cellId"`
		EventType            string      `json:"eventType"`
		ServiceNamespace     string      `json:"serviceNamespace"`
		EventStatus          string      `json:"eventStatus"`
		Timestamp            string      `json:"timestamp"`
		External             bool        `json:"external"`
		AdditionalProperties struct {
			UserRoles                         string `json:"user.roles"`
			UserSessionID                     string `json:"user.session.id"`
			CurrentContextUserProxyAddress    string `json:"currentContext.user.proxyAddress"`
			CurrentContextUserClientIPAddress string `json:"currentContext.user.clientIpAddress"`
		} `json:"additionalProperties"`
	}

	respp := []*AudiTrail{}

	// FIQL filtering test
	queryParams := url.Values{}

	// Find audi trail logs for the last 12 hours
	filterTime := time.Now().Add(-4 * time.Hour).Format(types.FiqlQueryTimestampFormat)
	queryParams.Add("filter", "timestamp=gt="+filterTime+";"+"user.name==administrator")
	// queryParams.Add("filter", "user.name==administrator")
	// queryParams.Add("filter", "description=administrator")
	// queryParams.Add("filter", "timestamp=gt="+filterTime+";"+"title==foo*;(updated=lt=-P1D,title==*bar)")
	// queryParams.Add("filter", "title==foo*;(updated=lt=-P1D,title==*bar)")

	err = vcd.vdc.client.CloudApiGetAllItems(urll, queryParams, &respp)

	spew.Dump(respp)
	check.Assert(err, IsNil)

	for _, v := range respp {
		fmt.Println(v.Timestamp, " ", v.EventType)
	}

	fmt.Println(len(respp))
}

// func (vdc *Vdc) GetCloudAPIEdgeGateways(queryParameters url.Values) ([]*types.CloudAPIEdgeGateway, error) {
// 	urlString := vdc.client.VCDHREF.Scheme + "://" + vdc.client.VCDHREF.Host + "/cloudapi/1.0.0/edgeGateways"
// 	url, _ := url.ParseRequestURI(urlString)
//
// 	response := make([]*types.CloudAPIEdgeGateway, 1)
//
// 	// err := vdc.client.CloudApiGetAllItems(url, queryParameters, "error getting edge gateways %s", nil, &response)
// 	err := vdc.client.CloudApiGetAllItems(url, queryParameters, "error getting edge gateways %s", nil, &response)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return response, nil
// }
