package govcd

import (
	"github.com/vmware/go-vcloud-director/v3/types/v56"
)

const labelTmOrgNetworkingSettings = "Org Networking Settings"

// TmOrgNetworkingSettings retrieves and updates networking-specific settings for the given
// organization
type TmOrgNetworkingSettings struct {
	TmOrgNetworkingSettings *types.TmOrgNetworkingSettings
	vcdClient               *VCDClient
	OrgId                   string
}

// wrap is a hidden helper that facilitates the usage of a generic CRUD function
//
//lint:ignore U1000 this method is used in generic functions, but annoys staticcheck
func (g TmOrgNetworkingSettings) wrap(inner *types.TmOrgNetworkingSettings) *TmOrgNetworkingSettings {
	g.TmOrgNetworkingSettings = inner
	return &g
}

// GetTmOrgNetworkingSettingsById retrieves Org Networking Settings by ID
func (vcdClient *VCDClient) GetTmOrgNetworkingSettingsByOrgId(orgId string) (*TmOrgNetworkingSettings, error) {
	c := crudConfig{
		entityLabel:    labelTmOrgNetworkingSettings,
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointTmOrgNetworkingSettings,
		endpointParams: []string{orgId},
		requiresTm:     true,
	}

	outerType := TmOrgNetworkingSettings{vcdClient: vcdClient, OrgId: orgId}
	return getOuterEntity(&vcdClient.Client, outerType, c)
}

// Update Org Networking Settings with a given config
func (o *TmOrgNetworkingSettings) Update(config *types.TmOrgNetworkingSettings) (*TmOrgNetworkingSettings, error) {
	c := crudConfig{
		entityLabel:    labelTmOrgNetworkingSettings,
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointTmOrgNetworkingSettings,
		endpointParams: []string{o.OrgId},
		requiresTm:     true,
	}
	outerType := TmOrgNetworkingSettings{vcdClient: o.vcdClient, OrgId: o.OrgId}
	return updateOuterEntity(&o.vcdClient.Client, outerType, c, config)
}

// Delete Org Networking Settings
func (o *TmOrgNetworkingSettings) Delete() error {
	_, err := o.Update(&types.TmOrgNetworkingSettings{})
	return err
}
