// +build alb functional ALL

package govcd

import (
	"fmt"

	. "gopkg.in/check.v1"
)

// Tests VDC storage profile update
func (vcd *TestVCD) Test_GetAllAlbControllers(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}
	controllers, err := vcd.client.GetAllAlbControllers(nil)
	check.Assert(err, IsNil)
	check.Assert(len(controllers), Equals, 1)
}

// Test_CreateAlbController
func (vcd *TestVCD) Test_CreateAlbController(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	skipNoNsxtAlbConfiguration(vcd, check)

	// Check if ALB controller is already added to the configuration
	existingController, err := vcd.client.GetAlbControllerByUrl(vcd.config.VCD.Nsxt.NsxtAlbControllerUrl)
	if !ContainsNotFound(err) {
		check.Assert(err, IsNil)
		check.Assert(existingController.NsxtAlbController.Name, Equals, "aviController1")
	}

	if ContainsNotFound(err) {

	}

	controllers, err := vcd.client.GetAllAlbControllers(nil)
	check.Assert(err, IsNil)
	check.Assert(len(controllers), Equals, 1)

	err = controllers[0].Delete()
	check.Assert(err, IsNil)

	// Create it again
	controllers[0].NsxtAlbController.Username = vcd.config.VCD.Nsxt.NsxtAlbControllerUser
	controllers[0].NsxtAlbController.Password = vcd.config.VCD.Nsxt.NsxtAlbControllerPassword
	newController, err := vcd.client.CreateNsxtAlbController(controllers[0].NsxtAlbController)
	check.Assert(err, IsNil)
	check.Assert(newController.NsxtAlbController.ID, Not(Equals), "")

	// GetByName
	controllerByName, err := vcd.client.GetAlbControllerByName(newController.NsxtAlbController.Name)
	check.Assert(err, IsNil)
	check.Assert(controllerByName.NsxtAlbController.ID, Equals, newController.NsxtAlbController.ID)

}

func createAlbControllerIfNotExists(vcd *TestVCD, check *C) {

}
