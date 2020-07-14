/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

var (
	username        string
	password        string
	org             string
	apiEndpoint     string
	customAdfsRptId string
)

func init() {
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")
	flag.StringVar(&org, "org", "System", "Org name. Default is 'System'")
	flag.StringVar(&apiEndpoint, "endpoint", "", "API endpoint (e.g. 'https://hostname/api')")
}

// Usage:
// # go build -o cloudapi
// # ./cloudapi --username my_user --password my_secret_password --org my-org --endpoint https://192.168.1.160/api
func main() {
	flag.Parse()

	if username == "" || password == "" || org == "" || apiEndpoint == "" {
		fmt.Printf("At least 'username', 'password', 'org' and 'endpoint' must be specified\n")
		os.Exit(1)
	}

	vcdURL, err := url.Parse(apiEndpoint)
	if err != nil {
		fmt.Printf("Error parsing supplied endpoint %s: %s", apiEndpoint, err)
		os.Exit(2)
	}

	vcdCli := govcd.NewVCDClient(*vcdURL, true)
	err = vcdCli.Authenticate(username, password, org)
	if err != nil {

		fmt.Println(err)
		os.Exit(3)
	}

}

// cloudAPIGetRawAuditTrail is an example function how to use low level function to interact with CloudAPI in VCD
func cloudAPIGetRawAuditTrail(vcdClient *govcd.VCDClient) {
	urlRef, err := BuildCloudAPIEndpoint(vcdClient, "1.0.0/auditTrail")

	response := []json.RawMessage{}
	err = vcdClient.CloudApiGetAllItems(urlRef, nil, &response)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Got %d results", len(response))

	for index, value := range response {
		jsonBody, err := value.MarshalJSON()
		if err != nil {

		}
		fmt.Printf("%d, :%s", index, jsonBody)
	}
}

// cloudAPIGetStructAuditTrail is an example function how to use low level function to interact with CloudAPI in VCD and
// marshal responses into defined struct with tags.
func cloudAPIGetStructAuditTrail(vcdClient *govcd.VCDClient) {
	urlRef, err := BuildCloudAPIEndpoint(vcdClient, "1.0.0/auditTrail")

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

	response := []*AudiTrail{}

	err = vcdClient.CloudApiGetAllItems(urlRef, nil, &response)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Got %d results", len(response))

	for _, value := range response {
		fmt.Printf("%s - %s, -%s", value.Timestamp, value.User.Name, value.EventType)
	}
}
