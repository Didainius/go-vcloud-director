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
