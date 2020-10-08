package types

const (
	gatewayTypeNsxT = "NSXT_BACKED"
	gatewayTypeNsxV = "NSXV_BACKED"
)

type NsxtEdgeGateway struct {
	Status      string `json:"status,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OrgVdc      struct {
		ID string `json:"id"`
	} `json:"orgVdc"`
	EdgeGatewayUplinks []struct {
		UplinkID   string `json:"uplinkId"`
		UplinkName string `json:"uplinkName"`
		Subnets    struct {
			Values []struct {
				Gateway      string      `json:"gateway"`
				PrefixLength int         `json:"prefixLength"`
				DNSSuffix    interface{} `json:"dnsSuffix"`
				DNSServer1   string      `json:"dnsServer1"`
				DNSServer2   string      `json:"dnsServer2"`
				IPRanges     struct {
					Values []struct {
						StartAddress string `json:"startAddress"`
						EndAddress   string `json:"endAddress"`
					} `json:"values"`
				} `json:"ipRanges"`
				Enabled      bool   `json:"enabled"`
				TotalIPCount int    `json:"totalIpCount"`
				UsedIPCount  int    `json:"usedIpCount"`
				PrimaryIp    string `json:"primaryIp,omitempty"`
			} `json:"values"`
		} `json:"subnets"`
		Dedicated bool `json:"dedicated"`
	} `json:"edgeGatewayUplinks"`
}

type CloudAPIEdgeGateway2 struct {
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
