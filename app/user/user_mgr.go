
package user

import (
	"../../firedog/fd"
	"../../firedog/util"
	"fmt"
	"time"
)

var Dispatcher *fd.EventDispatcher

type UserMgr struct {
}

func init() {
	Dispatcher = new(fd.EventDispatcher)
	util.Module("user_mgr", new(UserMgr))
}

func (this *UserMgr) Test() {

	Dispatcher.AddEventListener("ghj", handler)
	fmt.Println("..test ...")
}

func handler(event fd.Event) {
	fmt.Println("..evet ...", event)
}

func (this *UserMgr) Notify() {
	Dispatcher.Notify("ghj", nil)
}

func (this *UserMgr) TestTimer() {
	fd.Timer.SetInterval(1*time.Second, 0, timer, nil)
	fmt.Println("Test Timer")
}

func timer(args []interface{}) {
	fmt.Println("timer args ", args)
}
