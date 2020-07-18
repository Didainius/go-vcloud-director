package types

import (
	"encoding/json"
	"fmt"
)

const (
	gatewayTypeNsxT = "NSXT_BACKED"
	gatewayTypeNsxV = "NSXV_BACKED"
)

// CloudApiPageValues is a slice of json.RawMessage. json.RawMessage itself allows to partially marshal responses and is
// used to decouple API paging handling from particular returned types.
type CloudApiPageValues []json.RawMessage

// CloudApiPages unwraps pagination for "Get All" endpoints in CloudAPI. It uses a type "CloudApiPageValues" for values
// which are kept int []json.RawMessage. json.RawMessage helps to decouple marshalling paging related information from
// exact type related information. Paging can be handled dynamically this way while values can be marshaled into exact
// types.
type CloudApiPages struct {
	// ResultTotal reports total results available
	ResultTotal int `json:"resultTotal,omitempty"`
	// PageCount reports total result pages available
	PageCount int `json:"pageCount,omitempty"`
	// Page reports current page of result
	Page int `json:"page,omitempty"`
	// PageSize reports pagesize
	PageSize int `json:"pageSize,omitempty"`
	// Associations ...
	Associations interface{} `json:"associations,omitempty"`
	// Values holds types depending on the endpoint therefore `json.RawMessage` is used to dynamically unmarshal into
	// specific type as required
	Values json.RawMessage `json:"values,omitempty"`
}

////
// AccumulatePageResponses helps to accumulate Raw JSON objects during a pagination query
type AccumulatePageResponses []json.RawMessage

type CloudAPIEdgeGateway struct {
	Status                    string               `json:"status,omitempty"`
	ID                        string               `json:"id,omitempty"`
	Name                      string               `json:"name,omitempty"`
	Description               string               `json:"description,omitempty"`
	EdgeGatewayUplinks        []EdgeGatewayUplinks `json:"edgeGatewayUplinks,omitempty"`
	DistributedRoutingEnabled bool                 `json:"distributedRoutingEnabled,omitempty"`
	OrgVdcNetworkCount        int                  `json:"orgVdcNetworkCount,omitempty"`
	// GatewayBacking            GatewayBacking       `json:"gatewayBacking,omitempty"`
	OrgVdc                   OrgVdc            `json:"orgVdc,omitempty"`
	OrgRef                   OrgRef            `json:"orgRef,omitempty"`
	ServiceNetworkDefinition string            `json:"serviceNetworkDefinition,omitempty"`
	EdgeClusterConfig        EdgeClusterConfig `json:"edgeClusterConfig,omitempty"`
}
type IPRanges2 struct {
	Values []interface{} `json:"values,omitempty"`
}
type Values struct {
	Gateway              string      `json:"gateway,omitempty"`
	PrefixLength         int         `json:"prefixLength,omitempty"`
	DNSSuffix            interface{} `json:"dnsSuffix,omitempty"`
	DNSServer1           string      `json:"dnsServer1,omitempty"`
	DNSServer2           string      `json:"dnsServer2,omitempty"`
	IPRanges             IPRanges2   `json:"ipRanges,omitempty"`
	Enabled              bool        `json:"enabled,omitempty"`
	TotalIPCount         int         `json:"totalIpCount,omitempty"`
	UsedIPCount          interface{} `json:"usedIpCount,omitempty"`
	PrimaryIP            string      `json:"primaryIp,omitempty"`
	AutoAllocateIPRanges bool        `json:"autoAllocateIpRanges,omitempty"`
}
type Subnets struct {
	Values []Values `json:"values,omitempty"`
}
type EdgeGatewayUplinks struct {
	UplinkID                 string      `json:"uplinkId,omitempty"`
	UplinkName               string      `json:"uplinkName,omitempty"`
	Subnets                  Subnets     `json:"subnets,omitempty"`
	Connected                bool        `json:"connected,omitempty"`
	QuickAddAllocatedIPCount interface{} `json:"quickAddAllocatedIpCount,omitempty"`
	Dedicated                bool        `json:"dedicated,omitempty"`
}
type NetworkProvider struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}
type GatewayBacking struct {
	BackingID       string          `json:"backingId,omitempty"`
	GatewayType     string          `json:"gatewayType,omitempty"`
	NetworkProvider NetworkProvider `json:"networkProvider,omitempty"`
}
type OrgVdc struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}
type OrgRef struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}
type EdgeClusterRef struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}
type PrimaryEdgeCluster struct {
	EdgeClusterRef EdgeClusterRef `json:"edgeClusterRef,omitempty"`
	BackingID      string         `json:"backingId,omitempty"`
}
type EdgeClusterConfig struct {
	PrimaryEdgeCluster   PrimaryEdgeCluster `json:"primaryEdgeCluster,omitempty"`
	SecondaryEdgeCluster interface{}        `json:"secondaryEdgeCluster,omitempty"`
}

// OpenApiError helpes to marshal and provider meaningful `Error` for
type OpenApiError struct {
	MinorErrorCode string `json:"minorErrorCode"`
	Message        string `json:"message"`
	StackTrace     string `json:"stackTrace"`
}

// Error method implements Go's default `error` interface for CloudAPI errors formats them for human readable output.
func (openApiError OpenApiError) Error() string {
	return fmt.Sprintf("%s - %s", openApiError.MinorErrorCode, openApiError.Message)
}

// ErrorWithStack is the same as `Error()`, but also includes stack trace returned by API which is usually lengthy.
func (openApiError OpenApiError) ErrorWithStack() string {
	return fmt.Sprintf("%s - %s. Stack: %s", openApiError.MinorErrorCode, openApiError.Message,
		openApiError.StackTrace)
}
