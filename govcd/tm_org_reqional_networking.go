package govcd

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v3/types/v56"
)

const labelTmRegionalNetworkingSettings = "Tm Regional Networking Settings"

type TmRegionalNetworkingSettings struct {
	TmRegionalNetworkingSettings *types.TmRegionalNetworkingSettings
	vcdClient                    *VCDClient
}

// wrap is a hidden helper that facilitates the usage of a generic CRUD function
//
//lint:ignore U1000 this method is used in generic functions, but annoys staticcheck
func (g TmRegionalNetworkingSettings) wrap(inner *types.TmRegionalNetworkingSettings) *TmRegionalNetworkingSettings {
	g.TmRegionalNetworkingSettings = inner
	return &g
}

// CreateTmRegionalNetworkingSettings creates a new Tm Regional Networking Settings with a given configuration
func (vcdClient *VCDClient) CreateTmRegionalNetworkingSettings(config *types.TmRegionalNetworkingSettings) (*TmRegionalNetworkingSettings, error) {
	c := crudConfig{
		entityLabel: labelTmRegionalNetworkingSettings,
		endpoint:    types.OpenApiPathVcf + types.OpenApiEndpointTmRegionalNetworkSetting,
		requiresTm:  true,
	}
	outerType := TmRegionalNetworkingSettings{vcdClient: vcdClient}
	return createOuterEntity(&vcdClient.Client, outerType, c, config)
}

// GetAllTmRegionalNetworkingSettings retrieves all Tm Regional Networking Settings with an optional filter
func (vcdClient *VCDClient) GetAllTmRegionalNetworkingSettings(queryParameters url.Values) ([]*TmRegionalNetworkingSettings, error) {
	c := crudConfig{
		entityLabel:     labelTmRegionalNetworkingSettings,
		endpoint:        types.OpenApiPathVcf + types.OpenApiEndpointTmRegionalNetworkSetting,
		queryParameters: queryParameters,
		requiresTm:      true,
	}

	outerType := TmRegionalNetworkingSettings{vcdClient: vcdClient}
	return getAllOuterEntities(&vcdClient.Client, outerType, c)
}

// GetTmRegionalNetworkingSettingsByName retrieves Tm Regional Networking Settings by Name
func (vcdClient *VCDClient) GetTmRegionalNetworkingSettingsByName(name string) (*TmRegionalNetworkingSettings, error) {
	if name == "" {
		return nil, fmt.Errorf("%s lookup requires name", labelTmRegionalNetworkingSettings)
	}

	queryParams := url.Values{}
	queryParams.Add("filter", "name=="+name)

	filteredEntities, err := vcdClient.GetAllTmRegionalNetworkingSettings(queryParams)
	if err != nil {
		return nil, err
	}

	singleEntity, err := oneOrError("name", name, filteredEntities)
	if err != nil {
		return nil, err
	}

	return vcdClient.GetTmRegionalNetworkingSettingsById(singleEntity.TmRegionalNetworkingSettings.ID)
}

// GetTmRegionalNetworkingSettingsById retrieves Tm Regional Networking Settings by ID
func (vcdClient *VCDClient) GetTmRegionalNetworkingSettingsById(id string) (*TmRegionalNetworkingSettings, error) {
	c := crudConfig{
		entityLabel:    labelTmRegionalNetworkingSettings,
		endpoint:       types.OpenApiPathVcf + types.OpenApiEndpointTmRegionalNetworkSetting,
		endpointParams: []string{id},
		requiresTm:     true,
	}

	outerType := TmRegionalNetworkingSettings{vcdClient: vcdClient}
	return getOuterEntity(&vcdClient.Client, outerType, c)
}

// GetTmRegionalNetworkingSettingsByNameAndOrgId retrieves Tm Regional Networking Settings by Name and Org ID
func (vcdClient *VCDClient) GetTmRegionalNetworkingSettingsByNameAndOrgId(name, orgId string) (*TmRegionalNetworkingSettings, error) {
	if name == "" || orgId == "" {
		return nil, fmt.Errorf("%s lookup requires name and Org ID", labelTmRegionalNetworkingSettings)
	}

	queryParams := url.Values{}
	queryParams.Add("filter", "name=="+name)
	queryParams = queryParameterFilterAnd("orgRef.id=="+orgId, queryParams)

	filteredEntities, err := vcdClient.GetAllTmRegionalNetworkingSettings(queryParams)
	if err != nil {
		return nil, err
	}

	singleEntity, err := oneOrError("name", name, filteredEntities)
	if err != nil {
		return nil, err
	}

	return vcdClient.GetTmRegionalNetworkingSettingsById(singleEntity.TmRegionalNetworkingSettings.ID)
}

// Update Tm Regional Networking Settings with a given config
func (o *TmRegionalNetworkingSettings) Update(TmRegionalNetworkingSettingsConfig *types.TmRegionalNetworkingSettings) (*TmRegionalNetworkingSettings, error) {
	c := crudConfig{
		entityLabel:    labelTmRegionalNetworkingSettings,
		endpoint:       types.OpenApiPathVcf + types.OpenApiEndpointTmRegionalNetworkSetting,
		endpointParams: []string{o.TmRegionalNetworkingSettings.ID},
		requiresTm:     true,
	}
	outerType := TmRegionalNetworkingSettings{vcdClient: o.vcdClient}
	return updateOuterEntity(&o.vcdClient.Client, outerType, c, TmRegionalNetworkingSettingsConfig)
}

// Delete Tm Regional Networking Settings
func (o *TmRegionalNetworkingSettings) Delete() error {
	c := crudConfig{
		entityLabel:    labelTmRegionalNetworkingSettings,
		endpoint:       types.OpenApiPathVcf + types.OpenApiEndpointTmRegionalNetworkSetting,
		endpointParams: []string{o.TmRegionalNetworkingSettings.ID},
		requiresTm:     true,
	}
	return deleteEntityById(&o.vcdClient.Client, c)
}
