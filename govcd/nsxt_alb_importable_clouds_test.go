// +build alb functional ALL

package govcd

import (
	"fmt"

	. "gopkg.in/check.v1"
)

// Tests VDC storage profile update
func (vcd *TestVCD) Test_GetAllAlbImportableClouds(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	controllers, err := vcd.client.GetAllAlbControllers(nil)
	check.Assert(err, IsNil)
	check.Assert(len(controllers), Equals, 1)

	clientImportableClouds, err := vcd.client.GetAllAlbImportableCloud(controllers[0].NsxtAlbController.ID, nil)
	check.Assert(err, IsNil)
	check.Assert(len(clientImportableClouds), Equals, 1)

	controllerImportableClouds, err := controllers[0].GetAllAlbImportableCloud(nil)
	check.Assert(err, IsNil)
	check.Assert(len(controllerImportableClouds), Equals, 1)
}
