// +build alb functional ALL

package govcd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/vmware/go-vcloud-director/v2/types/v56"

	. "gopkg.in/check.v1"
)

// Test_NsxtAlbController tests out NSX-T ALB Controller capabilities.
// Having multiple ALB Controllers for testing is not easy, and it may be that the controller is
// already attached to VCD with child resources. This test deletes and recreates ALB controller by default,
// but it is possible to skip this step by setting GOVCD_TEST_LEAVE_ALB env variable
func (vcd *TestVCD) Test_NsxtAlbController(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}
	skipNoNsxtAlbConfiguration(vcd, check)

	existingController, err := vcd.client.GetAlbControllerByUrl(vcd.config.VCD.Nsxt.NsxtAlbControllerUrl)
	if !ContainsNotFound(err) { // Not empty error which is not ErrorEntityNotFound is not expected
		check.Assert(err, IsNil)
	}

	// If the ALB Controller is not found - it is expected to be created (and cleaned up afterwards to preserve state)
	if ContainsNotFound(err) {
		newControllerDef := &types.NsxtAlbController{
			Name:        "aviController1",
			Url:         vcd.config.VCD.Nsxt.NsxtAlbControllerUrl,
			Username:    vcd.config.VCD.Nsxt.NsxtAlbControllerUser,
			Password:    vcd.config.VCD.Nsxt.NsxtAlbControllerPassword,
			LicenseType: "ENTERPRISE",
		}

		newController, err := vcd.client.CreateNsxtAlbController(newControllerDef)
		check.Assert(err, IsNil)
		check.Assert(newController.NsxtAlbController.ID, Not(Equals), "")

		// The ALB Controller was not present during this test run therefore it is going to be cleaned up in the end of tests
		openApiEndpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbController + newController.NsxtAlbController.ID
		AddToCleanupListOpenApi(newController.NsxtAlbController.Name, check.TestName(), openApiEndpoint)
	}

	if os.Getenv("GOVCD_TEST_LEAVE_ALB") != "" && existingController != nil {
		fmt.Printf("# NSX-T ALB Controller '%s' is already configured. Leaving it.\n", existingController.NsxtAlbController.Name)
		fmt.Println("# Unset env variable 'GOVCD_TEST_LEAVE_ALB' to force recreation of it")
	}

	// If ALB Controller was found, but we explicitly want to recreate it
	if os.Getenv("GOVCD_TEST_LEAVE_ALB") == "" && existingController != nil {
		fmt.Printf("# env variable 'GOVCD_TEST_LEAVE_ALB' is not set. Recreating NSX-T ALB Controller '%s'\n",
			existingController.NsxtAlbController.Name)
		err = existingController.Delete()
		check.Assert(err, IsNil)

		newControllerDef := &types.NsxtAlbController{
			Name:        "aviController1",
			Url:         vcd.config.VCD.Nsxt.NsxtAlbControllerUrl,
			Username:    vcd.config.VCD.Nsxt.NsxtAlbControllerUser,
			Password:    vcd.config.VCD.Nsxt.NsxtAlbControllerPassword,
			LicenseType: "ENTERPRISE",
		}

		newController, err := vcd.client.CreateNsxtAlbController(newControllerDef)
		check.Assert(err, IsNil)
		check.Assert(newController.NsxtAlbController.ID, Not(Equals), "")
	}

	// Get by Url
	controllerByUrl, err := vcd.client.GetAlbControllerByUrl(vcd.config.VCD.Nsxt.NsxtAlbControllerUrl)
	check.Assert(err, IsNil)
	// Get by Name
	controllerByName, err := vcd.client.GetAlbControllerByName(controllerByUrl.NsxtAlbController.Name)
	check.Assert(err, IsNil)
	check.Assert(controllerByName.NsxtAlbController.ID, Equals, controllerByUrl.NsxtAlbController.ID)

	// Get all Controllers and expect to find at least the known one
	allControllers, err := vcd.client.GetAllAlbControllers(nil)
	check.Assert(err, IsNil)
	check.Assert(len(allControllers) > 0, Equals, true)
	var foundController bool
	for controllerIndex := range allControllers {
		if allControllers[controllerIndex].NsxtAlbController.ID == controllerByUrl.NsxtAlbController.ID {
			foundController = true
		}
	}
	check.Assert(foundController, Equals, true)

	// Check filtering for GetAllAlbControllers works
	filter := url.Values{}
	filter.Add("filter", "name=="+controllerByUrl.NsxtAlbController.Name)
	filteredControllers, err := vcd.client.GetAllAlbControllers(nil)
	check.Assert(err, IsNil)
	check.Assert(len(filteredControllers), Equals, 1)
	check.Assert(filteredControllers[0].NsxtAlbController.ID, Equals, controllerByUrl.NsxtAlbController.ID)

	// Test update of ALB controller
	updateControllerDef := &types.NsxtAlbController{
		ID:          controllerByUrl.NsxtAlbController.ID,
		Name:        controllerByUrl.NsxtAlbController.Name + "-update",
		Description: "Description set",
		Url:         vcd.config.VCD.Nsxt.NsxtAlbControllerUrl,
		Username:    vcd.config.VCD.Nsxt.NsxtAlbControllerUser,
		Password:    vcd.config.VCD.Nsxt.NsxtAlbControllerPassword,
		LicenseType: "BASIC",
	}
	updatedController, err := controllerByUrl.Update(updateControllerDef)
	check.Assert(err, IsNil)
	check.Assert(updatedController.NsxtAlbController.Name, Equals, updateControllerDef.Name)
	check.Assert(updatedController.NsxtAlbController.Description, Equals, updateControllerDef.Description)
	check.Assert(updatedController.NsxtAlbController.Url, Equals, updateControllerDef.Url)
	check.Assert(updatedController.NsxtAlbController.Username, Equals, updateControllerDef.Username)
	check.Assert(updatedController.NsxtAlbController.LicenseType, Equals, updateControllerDef.LicenseType)

	// Revert settings to original ones
	_, err = controllerByUrl.Update(controllerByUrl.NsxtAlbController)
	check.Assert(err, IsNil)

}
