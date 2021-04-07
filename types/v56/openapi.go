package types

import (
	"encoding/json"
	"fmt"
)

// OpenApiPages unwraps pagination for "Get All" endpoints in OpenAPI. Values kept in json.RawMessage helps to decouple
// marshalling paging related information from exact type related information. Paging can be handled dynamically this
// way while values can be marshaled into exact types.
type OpenApiPages struct {
	// ResultTotal reports total results available
	ResultTotal int `json:"resultTotal,omitempty"`
	// PageCount reports total result pages available
	PageCount int `json:"pageCount,omitempty"`
	// Page reports current page of result
	Page int `json:"page,omitempty"`
	// PageSize reports page size
	PageSize int `json:"pageSize,omitempty"`
	// Associations ...
	Associations interface{} `json:"associations,omitempty"`
	// Values holds types depending on the endpoint therefore `json.RawMessage` is used to dynamically unmarshal into
	// specific type as required
	Values json.RawMessage `json:"values,omitempty"`
}

// OpenApiError helps to marshal and provider meaningful `Error` for
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

// Role defines access roles in VCD
type Role struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BundleKey   string `json:"bundleKey"`
	ReadOnly    bool   `json:"readOnly"`
}

// NsxtTier0Router defines NSX-T Tier 0 router
type NsxtTier0Router struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
	DisplayName string `json:"displayName"`
}

// NsxtEdgeCluster is a struct to represent logical grouping of NSX-T Edge virtual machines.
type NsxtEdgeCluster struct {
	// ID contains edge cluster ID (UUID format)
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// NodeCount shows number of nodes in the edge cluster
	NodeCount int `json:"nodeCount"`
	// NodeType usually holds "EDGE_NODE"
	NodeType string `json:"nodeType"`
	// DeploymentType (e.g. "VIRTUAL_MACHINE")
	DeploymentType string `json:"deploymentType"`
}

// ExternalNetworkV2 defines a struct for OpenAPI endpoint which is capable of creating NSX-V or
// NSX-T external network based on provided NetworkBackings.
type ExternalNetworkV2 struct {
	// ID is unique for the network. This field is read-only.
	ID string `json:"id,omitempty"`
	// Name of the network.
	Name string `json:"name"`
	// Description of the network
	Description string `json:"description"`
	// Subnets define one or more subnets and IP allocation pools in edge gateway
	Subnets ExternalNetworkV2Subnets `json:"subnets"`
	// NetworkBackings for this external network. Describes if this external network is backed by
	// port groups, vCenter standard switch or an NSX-T Tier-0 router.
	NetworkBackings ExternalNetworkV2Backings `json:"networkBackings"`
}

// ExternalNetworkV2IPRange defines allocated IP pools for a subnet in external network
type ExternalNetworkV2IPRange struct {
	// StartAddress holds starting IP address in the range
	StartAddress string `json:"startAddress"`
	// EndAddress holds ending IP address in the range
	EndAddress string `json:"endAddress"`
}

// ExternalNetworkV2IPRanges contains slice of ExternalNetworkV2IPRange
type ExternalNetworkV2IPRanges struct {
	Values []ExternalNetworkV2IPRange `json:"values"`
}

// ExternalNetworkV2Subnets contains slice of ExternalNetworkV2Subnet
type ExternalNetworkV2Subnets struct {
	Values []ExternalNetworkV2Subnet `json:"values"`
}

// ExternalNetworkV2Subnet defines one subnet for external network with assigned static IP ranges
type ExternalNetworkV2Subnet struct {
	// Gateway for the subnet
	Gateway string `json:"gateway"`
	// PrefixLength holds prefix length of the subnet
	PrefixLength int `json:"prefixLength"`
	// DNSSuffix is the DNS suffix that VMs attached to this network will use (NSX-V only)
	DNSSuffix string `json:"dnsSuffix"`
	// DNSServer1 - first DNS server that VMs attached to this network will use (NSX-V only)
	DNSServer1 string `json:"dnsServer1"`
	// DNSServer2 - second DNS server that VMs attached to this network will use (NSX-V only)
	DNSServer2 string `json:"dnsServer2"`
	// Enabled indicates whether the external network subnet is currently enabled
	Enabled bool `json:"enabled"`
	// UsedIPCount shows number of IP addresses defined by the static IP ranges
	UsedIPCount int `json:"usedIpCount,omitempty"`
	// TotalIPCount shows number of IP address used from the static IP ranges
	TotalIPCount int `json:"totalIpCount,omitempty"`
	// IPRanges define allocated static IP pools allocated from a defined subnet
	IPRanges ExternalNetworkV2IPRanges `json:"ipRanges"`
}

type ExternalNetworkV2Backings struct {
	Values []ExternalNetworkV2Backing `json:"values"`
}

// ExternalNetworkV2Backing defines which networking subsystem is used for external network (NSX-T or NSX-V)
type ExternalNetworkV2Backing struct {
	// BackingID must contain either Tier-0 router ID for NSX-T or PortGroup ID for NSX-V
	BackingID string `json:"backingId"`
	Name      string `json:"name,omitempty"`
	// BackingType can be either ExternalNetworkBackingTypeNsxtTier0Router in case of NSX-T or one
	// of ExternalNetworkBackingTypeNetwork or ExternalNetworkBackingDvPortgroup in case of NSX-V
	BackingType string `json:"backingType"`
	// NetworkProvider defines backing network manager
	NetworkProvider NetworkProvider `json:"networkProvider"`
}

// NetworkProvider can be NSX-T manager or vCenter. ID is sufficient for creation purpose.
type NetworkProvider struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id"`
}

// VdcComputePolicy is represented as VM sizing policy in UI
type VdcComputePolicy struct {
	ID                         string   `json:"id,omitempty"`
	Description                string   `json:"description,omitempty"`
	Name                       string   `json:"name"`
	CPUSpeed                   *int     `json:"cpuSpeed,omitempty"`
	Memory                     *int     `json:"memory,omitempty"`
	CPUCount                   *int     `json:"cpuCount,omitempty"`
	CoresPerSocket             *int     `json:"coresPerSocket,omitempty"`
	MemoryReservationGuarantee *float64 `json:"memoryReservationGuarantee,omitempty"`
	CPUReservationGuarantee    *float64 `json:"cpuReservationGuarantee,omitempty"`
	CPULimit                   *int     `json:"cpuLimit,omitempty"`
	MemoryLimit                *int     `json:"memoryLimit,omitempty"`
	CPUShares                  *int     `json:"cpuShares,omitempty"`
	MemoryShares               *int     `json:"memoryShares,omitempty"`
	ExtraConfigs               *struct {
		AdditionalProp1 string `json:"additionalProp1,omitempty"`
		AdditionalProp2 string `json:"additionalProp2,omitempty"`
		AdditionalProp3 string `json:"additionalProp3,omitempty"`
	} `json:"extraConfigs,omitempty"`
	PvdcComputePolicyRef *struct {
		Name string `json:"name,omitempty"`
		ID   string `json:"id,omitempty"`
	} `json:"pvdcComputePolicyRef,omitempty"`
	PvdcComputePolicy *struct {
		Name string `json:"name,omitempty"`
		ID   string `json:"id,omitempty"`
	} `json:"pvdcComputePolicy,omitempty"`
	CompatibleVdcTypes []string `json:"compatibleVdcTypes,omitempty"`
	IsSizingOnly       bool     `json:"isSizingOnly,omitempty"`
	PvdcID             string   `json:"pvdcId,omitempty"`
	NamedVMGroups      [][]struct {
		Name string `json:"name,omitempty"`
		ID   string `json:"id,omitempty"`
	} `json:"namedVmGroups,omitempty"`
	LogicalVMGroupReferences []struct {
		Name string `json:"name,omitempty"`
		ID   string `json:"id,omitempty"`
	} `json:"logicalVmGroupReferences,omitempty"`
	IsAutoGenerated bool `json:"isAutoGenerated,omitempty"`
}

// OpenApiReference is a generic reference type commonly used throughout OpenAPI endpoints
type OpenApiReference struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type OpenApiReferences []OpenApiReference

// VdcCapability can be used to determine VDC capabilities, including such:
// * Is it backed by NSX-T or NSX-V pVdc
// * Does it support BGP routing
type VdcCapability struct {
	// Name of capability
	Name string `json:"name"`
	// Description of capability
	Description string `json:"description"`
	// Value can be any value. Sometimes it is a JSON bool (true, false), sometimes it is a JSON array (["custom", "default"])
	// and sometimes just a string ("NSX_V"). It is up for the consumer to handle values as per the Type field.
	Value interface{} `json:"value"`
	// Type of field (e.g. "Boolean", "String", "List")
	Type string `json:"type"`
	// Category of capability (e.g. "Security", "EdgeGateway", "OrgVdcNetwork")
	Category string `json:"category"`
}

// NsxtFirewallGroup allows to set either SECURITY_GROUP or IP_SET which is defined by Type field.
type NsxtFirewallGroup struct {
	// ID contains Firewall Group ID (URN format)
	// e.g. urn:vcloud:firewallGroup:d7f4e0b4-b83f-4a07-9f22-d242c9c0987a
	ID string `json:"id"`
	// Name of firewall group
	Name        string `json:"name"`
	Description string `json:"description"`
	// IP Addresses included in the group. This is only applicable for IP_SET Firewall Groups. This
	// can support IPv4 and IPv6 addresses in single, range, and CIDR formats.
	// E.g [
	//     "12.12.12.1",
	//     "10.10.10.0/24",
	//     "11.11.11.1-11.11.11.2"
	// ],
	IpAddresses []string `json:"ipAddresses,omitempty"`

	// Members define list of Org VDC networks belonging to this security group
	Members []OpenApiReference `json:"members,omitempty"`

	// OwnerRef replaces EdgeGatewayRef and can accept both - NSX-T Edge Gateway or a VDC group ID
	// Sample VDC Group URN - urn:vcloud:vdcGroup:89a53000-ef41-474d-80dc-82431ff8a020
	// Sample Edge Gateway URN - urn:vcloud:gateway:71df3e4b-6da9-404d-8e44-0865751c1c38
	//
	// Note. Using API V34.0 Firewall Groups can be created for VDC groups, but on a GET operation
	// there will be no VDC group ID returned.
	OwnerRef *OpenApiReference `json:"ownerRef,omitempty"`

	// EdgeGatewayRef is a deprecated field (use OwnerRef) for setting value, but during read only
	// it returns value
	EdgeGatewayRef *OpenApiReference `json:"edgeGatewayRef,omitempty"`

	// Type is either SECURITY_GROUP or IP_SET
	Type string `json:"type"`
}

// type NsxtFirewallGroup struct {
// 	// ID contains edge cluster ID (UUID format)
// 	ID          string `json:"id"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	// NodeCount shows number of nodes in the edge cluster
// 	NodeCount int `json:"nodeCount"`
// 	// NodeType usually holds "EDGE_NODE"
// 	NodeType string `json:"nodeType"`
// 	// DeploymentType (e.g. "VIRTUAL_MACHINE")
// 	DeploymentType string `json:"deploymentType"`
// }
