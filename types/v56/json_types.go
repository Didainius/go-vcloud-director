package types

//
// type ExternalNetworkV2 struct {
// 	Subnets struct {
// 		Values []struct {
// 			Gateway      string `json:"gateway"`
// 			PrefixLength int    `json:"prefixLength"`
// 			DNSSuffix    string `json:"dnsSuffix"`
// 			DNSServer1   string `json:"dnsServer1"`
// 			DNSServer2   string `json:"dnsServer2"`
// 			IPRanges     struct {
// 				Values []struct {
// 					StartAddress string `json:"startAddress"`
// 					EndAddress   string `json:"endAddress"`
// 				} `json:"values"`
// 			} `json:"ipRanges"`
// 			Enabled      bool `json:"enabled"`
// 			UsedIPCount  int  `json:"usedIpCount"`
// 			TotalIPCount int  `json:"totalIpCount"`
// 		} `json:"values"`
// 	} `json:"subnets"`
// 	Name            string      `json:"name"`
// 	Description     interface{} `json:"description"`
// 	NetworkBackings struct {
// 		Values []struct {
// 			BackingID       string `json:"backingId"`
// 			Name            string `json:"name"`
// 			BackingType     string `json:"backingType"`
// 			NetworkProvider struct {
// 				Name string `json:"name"`
// 				ID   string `json:"id"`
// 			} `json:"networkProvider"`
// 		} `json:"values"`
// 	} `json:"networkBackings"`
// }

type ExternalNetworkV2 struct {
	ID              string          `json:"id,omitempty"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Subnets         Subnets         `json:"subnets"`
	NetworkBackings NetworkBackings `json:"networkBackings"`
}
type IPRange2 struct {
	StartAddress string `json:"startAddress"`
	EndAddress   string `json:"endAddress"`
}
type IPRanges2 struct {
	Values []IPRange2 `json:"values"`
}

type Subnets struct {
	Values []Subnet `json:"values"`
}

type Subnet struct {
	Gateway      string    `json:"gateway"`
	PrefixLength int       `json:"prefixLength"`
	DNSSuffix    string    `json:"dnsSuffix"`
	DNSServer1   string    `json:"dnsServer1"`
	DNSServer2   string    `json:"dnsServer2"`
	IPRanges     IPRanges2 `json:"ipRanges"`
	Enabled      bool      `json:"enabled"`
	UsedIPCount  int       `json:"usedIpCount,omitempty"`
	TotalIPCount int       `json:"totalIpCount,omitempty"`
}

type NetworkBackings struct {
	Values []NetworkBacking `json:"values"`
}

type NetworkBacking struct {
	BackingID       string          `json:"backingId"`
	Name            string          `json:"name,omitempty"`
	BackingType     string          `json:"backingType"`
	NetworkProvider NetworkProvider `json:"networkProvider"`
}

type NetworkProvider struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id"`
}
