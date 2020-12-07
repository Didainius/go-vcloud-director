/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// NsxtLogicalSwitch is an NSX-T component which can be directly used to create NSX-T Imported Org Vdc networks.
//
// Note. With an NSX-T logical switch, you can create only an IPv4 isolated organization network. You cannot create a
// direct organization network based on an NSX-T logical switch.
type NsxtLogicalSwitch struct {
	NsxtLogicalSwitch *types.NsxtLogicalSwitch
	client            *Client
}

// GetNsxtLogicalSwitchByName retrieves a particular NSX-T Logical Switch by name available for that VDC
//
// Note. By default endpoint "/network/orgvdcnetworks/importableswitches" returns only unused NSX-T logical switches
// (the ones that are not already consumed in Org Vdc networks).
func (vdc *Vdc) GetNsxtLogicalSwitchByName(name string) (*NsxtLogicalSwitch, error) {
	if name == "" {
		return nil, fmt.Errorf("empty NSX-T Logical Switch name specified")
	}

	allNsxtImportableSwitches, err := vdc.GetAllNsxtLogicalSwitches()
	if err != nil {
		return nil, fmt.Errorf("error getting all Logical Switches for VDC '%s': %s", vdc.Vdc.Name, err)
	}

	var filteredNsxtImportableSwitches []*NsxtLogicalSwitch
	for _, nsxtImportableSwitch := range allNsxtImportableSwitches {
		if nsxtImportableSwitch.NsxtLogicalSwitch.Name == name {
			filteredNsxtImportableSwitches = append(filteredNsxtImportableSwitches, nsxtImportableSwitch)
		}
	}

	if len(filteredNsxtImportableSwitches) == 0 {
		// ErrorEntityNotFound is injected here for the ability to validate problem using ContainsNotFound()
		return nil, fmt.Errorf("%s: no NSX-T Logical Switch with name '%s' for Org Vdc with id '%s' found",
			ErrorEntityNotFound, name, vdc.Vdc.ID)
	}

	if len(filteredNsxtImportableSwitches) > 1 {
		return nil, fmt.Errorf("more than one (%d) NSX-T Logical Switch with name '%s' for Org Vdc with id '%s' found",
			len(filteredNsxtImportableSwitches), name, vdc.Vdc.ID)
	}

	return filteredNsxtImportableSwitches[0], nil
}

// GetAllNsxtLogicalSwitches retrieves all available logical switches which can be consumed for creating NSX-T
// "Imported" Org Vdc network
//
// Note. By default endpoint "/network/orgvdcnetworks/importableswitches" returns only unused NSX-T logical switches
// (the ones that are not already consumed in Org Vdc networks).
func (vdc *Vdc) GetAllNsxtLogicalSwitches() ([]*NsxtLogicalSwitch, error) {
	if vdc.Vdc.ID == "" {
		return nil, fmt.Errorf("VDC must have ID populated to retrieve NSX-T edge clusters")
	}

	apiEndpoint := vdc.client.VCDHREF
	endpoint := apiEndpoint.Scheme + "://" + apiEndpoint.Host + "/network/orgvdcnetworks/importableswitches"
	// error below is ignored because it is a static endpoint
	urlRef, _ := url.Parse(endpoint)

	// request requires Org Vdc ID to be specified as UUID, not as URN
	orgVdcId, err := getBareEntityUuid(vdc.Vdc.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get UUID from URN '%s': %s", vdc.Vdc.ID, err)
	}

	headAccept := http.Header{}
	headAccept.Set("Accept", types.JSONMime)
	request := vdc.client.newRequest(map[string]string{"orgVdc": orgVdcId}, nil, http.MethodGet, *urlRef, nil, vdc.client.APIVersion, headAccept)
	request.Header.Set("Accept", types.JSONMime)

	response, err := checkResp(vdc.client.Http.Do(request))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	nsxtLogicalSwitches := []*types.NsxtLogicalSwitch{}
	if err = decodeBody(types.BodyTypeJSON, response, &nsxtLogicalSwitches); err != nil {
		return nil, err
	}

	wrappedNsxtImportableSwitches := make([]*NsxtLogicalSwitch, len(nsxtLogicalSwitches))
	for sliceIndex := range nsxtLogicalSwitches {
		wrappedNsxtImportableSwitches[sliceIndex] = &NsxtLogicalSwitch{
			NsxtLogicalSwitch: nsxtLogicalSwitches[sliceIndex],
			client:            vdc.client,
		}
	}

	return wrappedNsxtImportableSwitches, nil
}
