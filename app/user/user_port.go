package user

import (
	"../../firedog/fd"
	"fmt"
)

func init() {
	fd.SetHandler("UserPort", UserPorthandler)
}

func UserPorthandler(session *fd.Session, msg ...interface{}) {

	fmt.Println("logic ", session, msg)

}
