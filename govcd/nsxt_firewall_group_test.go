// +build network nsxt functional openapi ALL

package govcd

import (
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

// Test_NsxtSecurityGroup tests out CRUD of NSX-T Security Group
//
// Note. Security Group is one type of Firewall Group
func (vcd *TestVCD) Test_NsxtSecurityGroup(check *C) {
	skipNoNsxtConfiguration(vcd, check)
	skipOpenApiEndpointTest(vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointFirewallGroups)

	org, err := vcd.client.GetOrgByName(vcd.config.VCD.Org)
	check.Assert(err, IsNil)

	nsxtVdc, err := org.GetVDCByName(vcd.config.VCD.Nsxt.Vdc, false)
	check.Assert(err, IsNil)

	edge, err := nsxtVdc.GetNsxtEdgeGatewayByName(vcd.config.VCD.Nsxt.EdgeGateway)
	check.Assert(err, IsNil)

	fwGroupDefinition := &types.NsxtFirewallGroup{
		Name:        check.TestName(),
		Description: check.TestName() + "-Description",
		Type:        types.FirewallGroupTypeSecurityGroup,
		OwnerRef:    &types.OpenApiReference{ID: edge.EdgeGateway.ID},
	}

	// Create firewall group and add to cleanup if it was created
	createdFwGroup, err := org.CreateNsxtFirewallGroup(fwGroupDefinition)
	check.Assert(err, IsNil)
	openApiEndpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointFirewallGroups + createdFwGroup.NsxtFirewallGroup.ID
	AddToCleanupListOpenApi(createdFwGroup.NsxtFirewallGroup.Name, check.TestName(), openApiEndpoint)

	check.Assert(createdFwGroup.NsxtFirewallGroup.ID, Not(Equals), "")
	check.Assert(createdFwGroup.NsxtFirewallGroup.EdgeGatewayRef.Name, Equals, vcd.config.VCD.Nsxt.EdgeGateway)

	// On creation one sets OwnerRef field, but in GET Edge Gateway is returned in EdgeGatewayRef
	// field
	check.Assert(createdFwGroup.NsxtFirewallGroup.EdgeGatewayRef.ID, Equals, fwGroupDefinition.OwnerRef.ID)
	check.Assert(createdFwGroup.NsxtFirewallGroup.Description, Equals, fwGroupDefinition.Description)
	check.Assert(createdFwGroup.NsxtFirewallGroup.Name, Equals, fwGroupDefinition.Name)
	check.Assert(createdFwGroup.NsxtFirewallGroup.Type, Equals, fwGroupDefinition.Type)

	// Update and compare
	createdFwGroup.NsxtFirewallGroup.Description = "updated-description"
	createdFwGroup.NsxtFirewallGroup.Name = check.TestName() + "-updated"

	updatedFwGroup, err := createdFwGroup.Update(createdFwGroup.NsxtFirewallGroup)
	check.Assert(err, IsNil)
	check.Assert(updatedFwGroup.NsxtFirewallGroup, DeepEquals, createdFwGroup.NsxtFirewallGroup)

	check.Assert(updatedFwGroup, DeepEquals, createdFwGroup)

	// Get all Firewall Groups and check if the created one is there
	allFwGroups, err := org.GetAllNsxtFirewallGroups(nil)
	check.Assert(err, IsNil)
	fwGroupFound := false
	for i := range allFwGroups {
		if allFwGroups[i].NsxtFirewallGroup.ID == updatedFwGroup.NsxtFirewallGroup.ID {
			fwGroupFound = true
			break
		}
	}
	check.Assert(fwGroupFound, Equals, true)

	// Get firewall group by name
	fwGroupByName, err := org.GetNsxtFirewallGroupByName(updatedFwGroup.NsxtFirewallGroup.Name, types.FirewallGroupTypeSecurityGroup)
	check.Assert(err, IsNil)

	fwGroupById, err := org.GetNsxtFirewallGroupById(updatedFwGroup.NsxtFirewallGroup.ID)
	check.Assert(err, IsNil)

	check.Assert(fwGroupById, DeepEquals, fwGroupByName)

	// Get Firewall Group using Edge Gateway
	edgeFirewallGroup, err := edge.GetNsxtFirewallGroupByName(updatedFwGroup.NsxtFirewallGroup.Name, types.FirewallGroupTypeSecurityGroup)
	check.Assert(err, IsNil)
	check.Assert(edgeFirewallGroup.NsxtFirewallGroup, DeepEquals, fwGroupByName.NsxtFirewallGroup)

	associatedVms, err := edgeFirewallGroup.GetAssociatedVms()
	// Try to list associated VMs and expect an empty list (because no Org VDC network is attached)
	check.Assert(err, IsNil)
	check.Assert(len(associatedVms), Equals, 0)

	// Remove
	err = createdFwGroup.Delete()
	check.Assert(err, IsNil)
}

// Test_NsxtSecurityGroupGetAssociatedVms tests if member routed Org VDC networks are added correctly to
// Security Groups and if associated VMs are correctly reported back
//
// Note. Security Group is one type of Firewall Group
func (vcd *TestVCD) Test_NsxtSecurityGroupGetAssociatedVms(check *C) {
	skipNoNsxtConfiguration(vcd, check)
	skipOpenApiEndpointTest(vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointFirewallGroups)

	org, err := vcd.client.GetOrgByName(vcd.config.VCD.Org)
	check.Assert(err, IsNil)

	nsxtVdc, err := org.GetVDCByName(vcd.config.VCD.Nsxt.Vdc, false)
	check.Assert(err, IsNil)

	edge, err := nsxtVdc.GetNsxtEdgeGatewayByName(vcd.config.VCD.Nsxt.EdgeGateway)
	check.Assert(err, IsNil)

	// Setup prerequisites - Routed Org VDC and add 2 VMs. With vApp and stadalone
	routedNet := createNsxtRoutedNetwork(check, vcd, nsxtVdc, edge.EdgeGateway.ID)
	openApiEndpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointOrgVdcNetworks + routedNet.OpenApiOrgVdcNetwork.ID
	AddToCleanupListOpenApi(routedNet.OpenApiOrgVdcNetwork.Name, check.TestName(), openApiEndpoint)

	// VMs are prependend to cleanup list to make sure they are removed before routed network
	standaloneVm := createStandaloneVm(check, vcd, nsxtVdc, routedNet)
	PrependToCleanupList(standaloneVm.VM.ID, "standaloneVm", "", check.TestName())

	vapp, vappVm := createVappVm(check, vcd, nsxtVdc, routedNet)
	PrependToCleanupList(vapp.VApp.Name, "vapp", vcd.nsxtVdc.Vdc.Name, check.TestName())

	fwGroupDefinition := &types.NsxtFirewallGroup{
		Name:        check.TestName(),
		Description: check.TestName() + "-Description",
		Type:        types.FirewallGroupTypeSecurityGroup,
		OwnerRef:    &types.OpenApiReference{ID: edge.EdgeGateway.ID},
		Members: []types.OpenApiReference{
			{ID: routedNet.OpenApiOrgVdcNetwork.ID},
		},
	}

	// Create firewall group and add to cleanup if it was created
	createdFwGroup, err := org.CreateNsxtFirewallGroup(fwGroupDefinition)
	check.Assert(err, IsNil)
	openApiEndpoint = types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointFirewallGroups + createdFwGroup.NsxtFirewallGroup.ID
	AddToCleanupListOpenApi(createdFwGroup.NsxtFirewallGroup.Name, check.TestName(), openApiEndpoint)

	// Expect to see VM created in associated VM query
	associatedVms, err := createdFwGroup.GetAssociatedVms()
	check.Assert(err, IsNil)

	check.Assert(len(associatedVms), Equals, 2)

	foundStandalone := false
	foundVappVm := false
	for i := range associatedVms {
		if associatedVms[i].Vmref.ID == standaloneVm.VM.ID {
			foundStandalone = true
		}

		if associatedVms[i].Vappref != nil && associatedVms[i].Vmref.ID == vappVm.VM.ID &&
			associatedVms[i].Vappref.ID == vapp.VApp.ID {
			foundVappVm = true
		}
	}

	check.Assert(foundStandalone, Equals, true)
	check.Assert(foundVappVm, Equals, true)

}

func createNsxtRoutedNetwork(check *C, vcd *TestVCD, vdc *Vdc, edgeGatewayId string) *OpenApiOrgVdcNetwork {
	orgVdcNetworkConfig := &types.OpenApiOrgVdcNetwork{
		Name:        check.TestName() + "routed-net",
		Description: check.TestName() + "-description",
		OrgVdc:      &types.OpenApiReference{ID: vcd.nsxtVdc.Vdc.ID},

		NetworkType: types.OrgVdcNetworkTypeRouted,

		// Connection is used for "routed" network
		Connection: &types.Connection{
			RouterRef: types.OpenApiReference{
				ID: edgeGatewayId,
			},
			ConnectionType: "INTERNAL",
		},
		Subnets: types.OrgVdcNetworkSubnets{
			Values: []types.OrgVdcNetworkSubnetValues{
				{
					Gateway:      "2.1.1.1",
					PrefixLength: 24,
					IPRanges: types.OrgVdcNetworkSubnetIPRanges{
						Values: []types.OrgVdcNetworkSubnetIPRangeValues{
							{
								StartAddress: "2.1.1.20",
								EndAddress:   "2.1.1.30",
							},
						}},
				},
			},
		},
	}

	orgVdcNet, err := vdc.CreateOpenApiOrgVdcNetwork(orgVdcNetworkConfig)
	check.Assert(err, IsNil)
	return orgVdcNet
}

func createStandaloneVm(check *C, vcd *TestVCD, vdc *Vdc, net *OpenApiOrgVdcNetwork) *VM {
	params := types.CreateVmParams{
		Name:    check.TestName(),
		PowerOn: false,
		// Description: description,
		CreateVm: &types.Vm{
			Name:                   check.TestName(),
			VirtualHardwareSection: nil,
			NetworkConnectionSection: &types.NetworkConnectionSection{
				Info:                          "Network Configuration for VM",
				PrimaryNetworkConnectionIndex: 0,
				NetworkConnection: []*types.NetworkConnection{
					&types.NetworkConnection{
						Network:                 net.OpenApiOrgVdcNetwork.Name,
						NeedsCustomization:      false,
						NetworkConnectionIndex:  0,
						IPAddress:               "any",
						IsConnected:             true,
						IPAddressAllocationMode: "DHCP",
						NetworkAdapterType:      "VMXNET3",
					},
				},
				Link: nil,
			},
			VmSpecSection: &types.VmSpecSection{
				Modified:          takeBoolPointer(true),
				Info:              "Virtual Machine specification",
				OsType:            "debian10Guest",
				NumCpus:           takeIntAddress(1),
				NumCoresPerSocket: takeIntAddress(1),
				CpuResourceMhz: &types.CpuResourceMhz{
					Configured: 0,
				},
				MemoryResourceMb: &types.MemoryResourceMb{
					Configured: 512,
				},
				DiskSection: &types.DiskSection{
					DiskSettings: []*types.DiskSettings{
						&types.DiskSettings{
							SizeMb:            1024,
							UnitNumber:        0,
							BusNumber:         0,
							AdapterType:       "5",
							ThinProvisioned:   takeBoolPointer(true),
							OverrideVmDefault: false,
						},
					},
				},

				HardwareVersion: &types.HardwareVersion{Value: "vmx-14"},
				VmToolsVersion:  "",
				VirtualCpuType:  "VM32",
			},
			GuestCustomizationSection: &types.GuestCustomizationSection{
				Info:         "Specifies Guest OS Customization Settings",
				ComputerName: "standalone1",
			},
		},
		Xmlns: types.XMLNamespaceVCloud,
	}

	vm, err := vdc.CreateStandaloneVm(&params)
	check.Assert(err, IsNil)
	check.Assert(vm, NotNil)

	return vm
}

func createVappVm(check *C, vcd *TestVCD, vdc *Vdc, net *OpenApiOrgVdcNetwork) (*VApp, *VM) {
	err := vdc.ComposeRawVApp(check.TestName())
	check.Assert(err, IsNil)

	vapp, err := vdc.GetVAppByName(check.TestName(), true)
	check.Assert(err, IsNil)

	desiredNetConfig := &types.NetworkConnectionSection{}
	desiredNetConfig.PrimaryNetworkConnectionIndex = 0
	desiredNetConfig.NetworkConnection = append(desiredNetConfig.NetworkConnection,
		&types.NetworkConnection{
			IsConnected:             true,
			IPAddressAllocationMode: types.IPAllocationModePool,
			Network:                 net.OpenApiOrgVdcNetwork.Name,
			NetworkConnectionIndex:  0,
		},
	)

	emptyVmDefinition := &types.RecomposeVAppParamsForEmptyVm{
		CreateItem: &types.CreateItem{
			Name:                      "Test_AddNewEmptyVMMultiNIC",
			NetworkConnectionSection:  desiredNetConfig,
			Description:               "created by Test_AddNewEmptyVMMultiNIC",
			GuestCustomizationSection: nil,
			VmSpecSection: &types.VmSpecSection{
				Modified:          takeBoolPointer(true),
				Info:              "Virtual Machine specification",
				OsType:            "debian10Guest",
				NumCpus:           takeIntAddress(2),
				NumCoresPerSocket: takeIntAddress(1),
				CpuResourceMhz:    &types.CpuResourceMhz{Configured: 1},
				MemoryResourceMb:  &types.MemoryResourceMb{Configured: 1024},
				DiskSection: &types.DiskSection{DiskSettings: []*types.DiskSettings{
					&types.DiskSettings{
						AdapterType:       "5",
						SizeMb:            int64(16384),
						BusNumber:         0,
						UnitNumber:        0,
						ThinProvisioned:   takeBoolPointer(true),
						OverrideVmDefault: true,
					},
				}},
				HardwareVersion:  &types.HardwareVersion{Value: "vmx-13"}, // need support older version vCD
				VmToolsVersion:   "",
				VirtualCpuType:   "VM32",
				TimeSyncWithHost: nil,
			},
		},
		AllEULAsAccepted: true,
	}

	createdVm, err := vapp.AddEmptyVm(emptyVmDefinition)
	check.Assert(err, IsNil)
	check.Assert(createdVm, NotNil)

	return vapp, createdVm
}
