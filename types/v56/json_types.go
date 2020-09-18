/*
 * Copyright 2019 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package types

type ExternalNetworkV2 struct {
	ID              string                    `json:"id,omitempty"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	Subnets         ExternalNetworkV2Subnets  `json:"subnets"`
	NetworkBackings ExternalNetworkV2Backings `json:"networkBackings"`
}
type ExternalNetworkV2IPRange struct {
	StartAddress string `json:"startAddress"`
	EndAddress   string `json:"endAddress"`
}
type ExternalNetworkV2IPRanges struct {
	Values []ExternalNetworkV2IPRange `json:"values"`
}

type ExternalNetworkV2Subnets struct {
	Values []ExternalNetworkV2Subnet `json:"values"`
}

type ExternalNetworkV2Subnet struct {
	Gateway      string                    `json:"gateway"`
	PrefixLength int                       `json:"prefixLength"`
	DNSSuffix    string                    `json:"dnsSuffix"`
	DNSServer1   string                    `json:"dnsServer1"`
	DNSServer2   string                    `json:"dnsServer2"`
	IPRanges     ExternalNetworkV2IPRanges `json:"ipRanges"`
	Enabled      bool                      `json:"enabled"`
	UsedIPCount  int                       `json:"usedIpCount,omitempty"`
	TotalIPCount int                       `json:"totalIpCount,omitempty"`
}

type ExternalNetworkV2Backings struct {
	Values []ExternalNetworkV2Backing `json:"values"`
}

type ExternalNetworkV2Backing struct {
	BackingID       string                  `json:"backingId"`
	Name            string                  `json:"name,omitempty"`
	BackingType     string                  `json:"backingType"`
	NetworkProvider NetworkProviderProvider `json:"networkProvider"`
}

type NetworkProviderProvider struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id"`
}
