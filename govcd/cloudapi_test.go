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

	responses := []json.RawMessage{{}}
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

	respp := []*AudiTrail{{}}

	// FIQL filtering test
	queryParams := url.Values{}

	// Find audi trail logs for the last 12 hours
	filterTime := time.Now().Add(-24 * 30 * time.Hour).Format(types.FiqlQueryTimestampFormat)
	queryParams.Add("filter", "timestamp=gt="+filterTime+";"+"user.id!=urn:vcloud:user:2e15afcf-ae56-4287-b7c9-4139d4592d0c")
	// queryParams.Add("filter", "timestamp=gt="+filterTime)
	// queryParams.Add("filter", "user.name==administrator")
	// queryParams.Add("filter", "description=administrator")
	// queryParams.Add("filter", "timestamp=gt="+filterTime+";"+"title==foo*;(updated=lt=-P1D,title==*bar)")
	// queryParams.Add("filter", "title==foo*;(updated=lt=-P1D,title==*bar)")

	err = vcd.vdc.client.CloudApiGetAllItems(urll, queryParams, &respp)

	spew.Dump(respp)
	check.Assert(err, IsNil)

	for _, v := range respp {
		fmt.Println(v.Timestamp, " ", v.User.Name, "-", v.AdditionalProperties.UserRoles, "", " ", v.EventType)
	}

	fmt.Println(len(respp))
}

// Test_CloudAPIInlineStructCRUDRoles test aims to test out low level CloudAPI functions to check if all of them work as
// expected. It uses a very simple "Roles" endpoint which does not have bigger prerequisites and therefore is not
// dependent one more deployment specific features.
// Actions of the test:
// 1. Get all available roles using "Get all endpoint"
// 2.
func (vcd *TestVCD) Test_CloudAPIInlineStructCRUDRoles(check *C) {
	// Step 1 - Get all roles
	urlRef, err := vcd.client.Client.BuildCloudAPIEndpoint("1.0.0/roles")
	check.Assert(err, IsNil)

	type Roles struct {
		ID          string `json:"id,omitempty"`
		Name        string `json:"name"`
		Description string `json:"description"`
		BundleKey   string `json:"bundleKey"`
		ReadOnly    bool   `json:"readOnly"`
	}

	allExistingRoles := []*Roles{{}}
	err = vcd.vdc.client.CloudApiGetAllItems(urlRef, nil, &allExistingRoles)

	// Step 2 - Get all roles using query filters
	for _, oneRole := range allExistingRoles {
		// Step 2.1 - retrieve specific role by using FIQL filter
		urlRef2, err := vcd.client.Client.BuildCloudAPIEndpoint("1.0.0/roles")
		check.Assert(err, IsNil)

		queryParams := url.Values{}
		queryParams.Add("filter", "id=="+oneRole.ID)

		expectOneRoleResultById := []*Roles{{}}

		// Use the same urlRef inject FIQL filter to test filtering by ID
		err = vcd.vdc.client.CloudApiGetAllItems(urlRef2, queryParams, &expectOneRoleResultById)
		check.Assert(err, IsNil)
		check.Assert(len(expectOneRoleResultById) == 1, Equals, true)

		// Step 2.2 - retrieve specific role by using endpoint
		singleRef, err := vcd.client.Client.BuildCloudAPIEndpoint("1.0.0/roles/" + oneRole.ID)
		oneRole := &Roles{}
		err = vcd.vdc.client.CloudApiGetItem(singleRef, nil, oneRole)
		check.Assert(err, IsNil)
		check.Assert(oneRole, NotNil)

		// Step 2.3 - compare struct retrieve by using filter and the one retrieve by exact endpoint ID
		check.Assert(oneRole, DeepEquals, expectOneRoleResultById[0])

	}

	// Step 3 - Create a new role

	createUrl, err := vcd.client.Client.BuildCloudAPIEndpoint("1.0.0/roles")
	check.Assert(err, IsNil)

	newRole := &Roles{
		Name:        check.TestName(),
		Description: "Role created by test",
		BundleKey:   "com.vmware.vcloud.undefined.key",
		ReadOnly:    false,
	}
	newRoleResponse := &Roles{}
	err = vcd.client.Client.CloudApiPostItem(createUrl, nil, newRole, newRoleResponse)
	check.Assert(err, IsNil)

	// Ensure supplied and created structs differ only by ID
	newRole.ID = newRoleResponse.ID
	check.Assert(newRoleResponse, DeepEquals, newRole)

	// Delete

	deleteUrlRef, err := vcd.client.Client.BuildCloudAPIEndpoint("1.0.0/roles/" + newRoleResponse.ID)
	check.Assert(err, IsNil)

	err = vcd.client.Client.CloudApiDeleteItem(deleteUrlRef, nil)
	check.Assert(err, IsNil)

	// Read is tricky - it throws an error ACCESS_TO_RESOURCE_IS_FORBIDDEN when the resource with ID does not
	// exist therefore one cannot know what kind of error occurred.
	lostRole := &Roles{}
	err = vcd.client.Client.CloudApiGetItem(deleteUrlRef, nil, lostRole)
	check.Assert(ContainsNotFound(err), Equals, true)

	//
	// // FIQL filtering test
	// queryParams := url.Values{}
	//
	// // Find existing Org Vdc network defined in config
	// queryParams.Add("filter", "name=="+vcd.config.VCD.Network.Net1)
	//
	//
	//
	// spew.Dump(respp)
	// check.Assert(err, IsNil)
	// //
	// // for _, v := range respp {
	// // 	fmt.Println(v.Timestamp, " ", v.User.Name, "-", v.AdditionalProperties.UserRoles, "", " ", v.EventType)
	// // }
	//
	// fmt.Println(len(respp))
}
