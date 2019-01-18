package user

import (
	"fmt"
	"github.com/gohuge/firedog/fd"
)

func init() {
	fd.SetHandler("UserPort", UserPorthandler)
}

func UserPorthandler(session *fd.Session, msg ...interface{}) {

	fmt.Println("logic ", session, msg)

}
