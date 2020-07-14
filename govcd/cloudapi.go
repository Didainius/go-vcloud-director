package govcd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/peterhellberg/link"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"github.com/vmware/go-vcloud-director/v2/util"
)

// This file contains generalised low level methods to interact with VCD CloudAPI REST endpoints as documented in
// https://{VCD_HOST}/api-explorer/tenant/tenant-name and https://{VCD_HOST}/api-explorer/provider documentation. It has
// functions supporting below methods:
// GET /items (gets a slice of types like `[]*types.CloudAPIEdgeGateway` or even `[]*json.RawMessage` to process JSON as text.
// POST /items
// GET /items/urn
// PUT /items/urn
// DELETE /items/urn
//
// GET endpoints support FIQL for filtering in field `filter`. (FIQL IETF doc - https://tools.ietf.org/html/draft-nottingham-atompub-fiql-00)

// CloudAPI versioning.
// Versions in path (e.g. 1.0.0) should guarantee behavior while header versions shouldn't matter in long term

// CloudApiGetAllItems retrieves and accumulates all pages then parsing them to a single object. It works by at first
// crawling pages and accumulating all responses into []json.RawMessage (as strings). Because there is no intermediate
// unmarshalling to exact `outType` for every page it can actually unmarshal into response struct in one go. outType
// must be a slice of object (e.g. []*types.CloudAPIEdgeGateway) because this response contains list of structs
func (client *Client) CloudApiGetAllItems(urlRef *url.URL, queryParams url.Values, outType interface{}) error {
	util.Logger.Printf("[TRACE] Getting all items from endpoint %s for parsing into %s type\n",
		urlRef.String(), reflect.TypeOf(outType))

	// Perform API call to initial endpoint. The function call below is expected to follow pages using Link headers
	// "nextPage" until it crawls all results
	responses, err := client.cloudApiGetAllPages(nil, urlRef, queryParams, outType, nil)
	if err != nil {
		return fmt.Errorf("error getting all pages for endpoint %s: %s", urlRef.String(), err)
	}

	// Create a slice of raw JSON messages in text so that they can be unmarshalled to specified `outType` after multiple
	// calls are executed
	var rawJsonBodies []string
	for _, singleObject := range responses {
		jsonBody, err := singleObject.MarshalJSON()
		if err != nil {
			return fmt.Errorf("error marshalling single response object into raw JSON message: %s", err)
		}
		rawJsonBodies = append(rawJsonBodies, string(jsonBody))
	}

	// rawJsonBodies contains a slice of all response objects and they must be formatted as a JSON slice (wrapped
	// into `[]`, separated with semicolons) so that unmarshalling to specified `outType` works in one go
	allResponses := `[` + strings.Join(rawJsonBodies, ",") + `]`

	// Unmarshal all accumulated responses into `outType`
	if err = json.Unmarshal([]byte(allResponses), &outType); err != nil {
		return fmt.Errorf("error decoding values into type: %s", err)
	}

	return nil
}
func (client *Client) cloudApiPostItem(urlRef *url.URL, params url.Values, payload, outType interface{}) error {
	util.Logger.Printf("[TRACE] Posting %s item to endpoint %s with expected response of type %s",
		reflect.TypeOf(payload), urlRef.String(), reflect.TypeOf(outType))

	// Marshal payload if we have one
	var body *bytes.Buffer
	// Check if we have body
	if payload != nil {
		marshaledJson, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshalling JSON data %s", err)
		}
		// fmt.Println(string(marshaledJson))
		body = bytes.NewBufferString(string(marshaledJson))
	}

	// Exec request
	req := client.newCloudApiRequest(params, http.MethodPost, urlRef, body, "34.0")
	req.Header.Add("Content-Type", types.JSONMime)

	resp, err := client.Http.Do(req)
	if err != nil {
		return err
	}

	// resp is ignore below because it is the same as above
	_, err = checkRespWithErrType(resp, err, &types.CloudApiError{}, types.BodyTypeJSON)
	if err != nil {
		return fmt.Errorf("error in HTTP POST request: %s", err)
	}

	// log response explicitly because decodeBody() was not triggered here
	bd, _ := ioutil.ReadAll(resp.Body)
	util.ProcessResponseOutput(util.FuncNameCallStack(), resp, fmt.Sprintf("%s", bd))

	debugShowResponse(resp, bd)

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing response body: %s", err)
	}

	// Placeholder for newly created object ID
	var newObjectId string

	// if we have task - track it

	// CloudAPI may work synchronously or asynchronously. When working asynchronously - it will return HTTP 202 and
	// `Location` header will contain reference to task so that it can be tracked. Task tracking itself is in the legacy
	// XML API.
	if resp.StatusCode == http.StatusAccepted {
		taskUrl := resp.Header.Get("Location")
		task := NewTask(client)
		task.Task.HREF = taskUrl
		err = task.WaitTaskCompletion()
		if err != nil {
			return fmt.Errorf("error waiting completion of task (%s): %s", taskUrl, err)
		}
		// Task Owner ID is the ID of created object
		newObjectId = task.Task.Owner.ID

	}

	// Here we have to find the resource once more to return it populated
	newObjectUrl, _ := url.ParseRequestURI(urlRef.String() + "/" + newObjectId)
	err = client.cloudApiGetItem(newObjectUrl, nil, outType)
	if err != nil {
		return fmt.Errorf("error retrieving item after creation: %s", err)
	}

	// The request was successful
	return nil
}
func (client *Client) cloudApiGetItem(urlRef *url.URL, params url.Values, outType interface{}) error {
	util.Logger.Printf("[TRACE] Getting item from endpoint %s with expected response of type %s", urlRef.String(), reflect.TypeOf(outType))

	// Exec request
	req := client.newCloudApiRequest(params, http.MethodGet, urlRef, nil, "34.0")
	// acceptMime := types.JSONMime + ";version=" + "34.0"
	// req.Header.Add("Accept", acceptMime)
	req.Header.Add("Content-Type", types.JSONMime)

	resp, err := client.Http.Do(req)
	if err != nil {
		return err
	}

	// resp is ignore below because it is the same as above
	_, err = checkRespWithErrType(resp, err, &types.CloudApiError{}, types.BodyTypeJSON)
	if err != nil {
		return fmt.Errorf("error in HTTP GET request: %s", err)
	}

	if err = decodeBody(resp, outType, types.BodyTypeJSON); err != nil {
		return fmt.Errorf("error decoding JSON response: %s", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing response body: %s", err)
	}

	// The request was successful
	return nil
}

// cloudApiPutItem handles the PUT method for CloudAPI (PUT /items/urn)
//
func (client *Client) cloudApiPutItem(urlRef *url.URL, params url.Values, payload, outType interface{}) error {
	util.Logger.Printf("[TRACE] Putting item of type %s at endpoint %s with expected response of type %s",
		reflect.TypeOf(payload), urlRef.String(), reflect.TypeOf(outType))

	// Marshal payload if we have one
	var body *bytes.Buffer
	if payload != nil {
		marshaledJson, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshalling JSON data %s", err)
		}
		body = bytes.NewBuffer(marshaledJson)
	}

	// Exec request
	req := client.newCloudApiRequest(params, http.MethodPut, urlRef, body, "34.0")
	req.Header.Add("Content-Type", types.JSONMime)

	resp, err := client.Http.Do(req)
	if err != nil {
		return err
	}

	// resp is ignored below because it is the same as above
	_, err = checkRespWithErrType(resp, err, &types.CloudApiError{}, types.BodyTypeJSON)
	if err != nil {
		return fmt.Errorf("error in HTTP PUT request: %s", err)
	}

	// log response explicitly because decodeBody() was not triggered here
	bd, _ := ioutil.ReadAll(resp.Body)
	util.ProcessResponseOutput(util.FuncNameCallStack(), resp, fmt.Sprintf("%s", bd))
	debugShowResponse(resp, bd)

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing PUT response body: %s", err)
	}

	// if we have task - track it

	// CloudAPI may work synchronously or asynchronously. When working asynchronously - it will return HTTP 202 and
	// `Location` header will contain reference to task so that it can be tracked. Task tracking itself is in the legacy
	// XML API.
	if resp.StatusCode == http.StatusAccepted {
		taskUrl := resp.Header.Get("Location")
		task := NewTask(client)
		task.Task.HREF = taskUrl
		err = task.WaitTaskCompletion()
		if err != nil {
			return fmt.Errorf("error waiting completion of task (%s): %s", taskUrl, err)
		}
	}

	// Here we have to find the resource once more to return it populated
	newObjectUrl, _ := url.ParseRequestURI(urlRef.String())
	err = client.cloudApiGetItem(newObjectUrl, nil, outType)
	if err != nil {
		return fmt.Errorf("error retrieving item after PUT operation: %s", err)
	}

	// The request was successful
	return nil
}

// cloudApiDeleteItem
func (client *Client) cloudApiDeleteItem(urlRef *url.URL, params url.Values) error {
	util.Logger.Printf("[TRACE] Deleting item at endpoint %s", urlRef.String())

	// Exec request
	req := client.newCloudApiRequest(params, http.MethodDelete, urlRef, nil, "34.0")
	req.Header.Add("Content-Type", types.JSONMime)

	resp, err := client.Http.Do(req)
	if err != nil {
		return err
	}

	// resp is ignored below because it is the same as above
	_, err = checkRespWithErrType(resp, err, &types.CloudApiError{}, types.BodyTypeJSON)
	if err != nil {
		return fmt.Errorf("error in HTTP DELETE request: %s", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing response body: %s", err)
	}

	// TODO add trace logging to see which HTTP status code was returned

	// CloudAPI may work synchronously or asynchronously. When working asynchronously - it will return HTTP 202 and
	// `Location` header will contain reference to task so that it can be tracked
	if resp.StatusCode == http.StatusAccepted {
		taskUrl := resp.Header.Get("Location")
		task := NewTask(client)
		task.Task.HREF = taskUrl
		err = task.WaitTaskCompletion()
		if err != nil {
			return fmt.Errorf("error waiting completion of task (%s): %s", taskUrl, err)
		}
	}

	// The request was successful
	return nil
}

// cloudApiGetAllPages helps to accumulate responses from multiple pages for GET query. It works by at first crawling
// pages and accumulating all responses into []json.RawMessage (as strings). Because there are no intermediate
// unmarshalling to exact `outType` for every page it can actually unmarshal into direct type passed.
// outType must be a slice of object (e.g. []*types.CloudAPIEdgeGateway) because accumulated responses are wrapped into
// JSON slice
func (client *Client) cloudApiGetAllPages(pageSize *int, urlRef *url.URL, queryParams url.Values, outType interface{}, responses []*json.RawMessage) ([]*json.RawMessage, error) {
	if responses == nil {
		responses = []*json.RawMessage{}
	}

	// Reuse existing queryParams struct to fill in pages or create a new one if nil was passed
	queryParameters := url.Values{}
	if queryParams != nil {
		queryParameters = queryParams
	}

	// if page != nil {
	// 	queryParameters.Set("page", strconv.Itoa(*page))
	// }

	if pageSize != nil {
		queryParameters.Set("pageSize", strconv.Itoa(*pageSize))
	}

	// Execute request
	req := client.newCloudApiRequest(queryParams, http.MethodGet, urlRef, nil, "34.0")
	req.Header.Add("Content-Type", types.JSONMime)

	resp, err := client.Http.Do(req)
	if err != nil {
		return nil, err
	}

	// resp is ignored below because it is the same as above
	_, err = checkRespWithErrType(resp, err, &types.CloudApiError{}, types.BodyTypeJSON)
	if err != nil {
		return nil, fmt.Errorf("error in HTTP GET request: %s", err)
	}

	// Pages will unwrap pagination and keep a slice of raw json message to marshal to specific types
	pages := &types.CloudApiPages{}

	if err = decodeBody(resp, pages, types.BodyTypeJSON); err != nil {
		return nil, fmt.Errorf("error decoding JSON page response: %s", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing response body: %s", err)
	}

	// Accumulate all responses in a single query
	// After pages are unwrapped one can marshal response into specified type
	// singleQueryResponses := &json.RawMessage{}
	var singleQueryResponses []*json.RawMessage
	if err = json.Unmarshal(pages.Values, &singleQueryResponses); err != nil {
		return nil, fmt.Errorf("error decoding values into accumulation type: %s", err)
	}
	responses = append(responses, singleQueryResponses...)

	// If there is a "nextPage" link in headers - follow it:
	links := link.ParseHeader(resp.Header)
	nextPage, _ := links["nextPage"]

	for k, v := range links {
		if strings.Contains(k, "nextPage") {
			nextPage = v
		}
	}

	// This must be tuned with proper checks

	// nextPage exists, follow it recursively and continue accumulating responses into single variable
	if nextPage != nil && nextPage.URI != "" {
		urlRef, _ = url.Parse(nextPage.String())

		responses, err = client.cloudApiGetAllPages(nil, urlRef, url.Values{}, outType, responses)
		if err != nil {
			return nil, fmt.Errorf("got error on page %d: %s", pages.Page, err)
		}
	}

	return responses, nil
}

// newCloudApiRequest is a low level function used in upstream CloudAPI functions which handles logging and
// authentication for each API request
func (client *Client) newCloudApiRequest(params url.Values, method string, reqUrl *url.URL, body io.Reader, apiVersion string) *http.Request {

	// Add the params to our URL
	reqUrl.RawQuery += params.Encode()

	// If the body contains data - try to read all contents for logging and re-create another
	// io.Reader with all contents to use it down the line
	var readBody []byte
	if body != nil {
		readBody, _ = ioutil.ReadAll(body)
		body = bytes.NewReader(readBody)
	}

	// Build the request, no point in checking for errors here as we're just
	// passing a string version of an url.URL struct and http.NewRequest returns
	// error only if can't process an url.ParseRequestURI().
	req, _ := http.NewRequest(method, reqUrl.String(), body)

	if client.VCDAuthHeader != "" && client.VCDToken != "" {
		// Add the authorization header
		req.Header.Add(client.VCDAuthHeader, client.VCDToken)
	}
	if client.VCDAuthHeader != "" && client.VCDToken != "" {
		// Add the Accept header for VCD
		acceptMime := types.JSONMime + ";version=" + apiVersion
		req.Header.Add("Accept", acceptMime)
	}

	// Avoids passing data if the logging of requests is disabled
	if util.LogHttpRequest {
		payload := ""
		if req.ContentLength > 0 {
			payload = string(readBody)
		}
		util.ProcessRequestOutput(util.FuncNameCallStack(), method, reqUrl.String(), payload, req)
		debugShowRequest(req, payload)
	}

	return req

}

func (client *Client) BuildCloudAPIEndpoint(endpoint string) (*url.URL, error) {
	endpointString := client.VCDHREF.Scheme + "://" + client.VCDHREF.Host + "/cloudapi/" + endpoint
	urlRef, err := url.ParseRequestURI(endpointString)
	if err != nil {
		return nil, fmt.Errorf("error formatting CloudAPI: %s", err)
	}
	return urlRef, nil
}
