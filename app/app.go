package app

import (
	"../firedog/fd"
	"../firedog/log"
	"../firedog/util"
	_ "./user"
)

type Module struct {
	Name string
}

var UserDispatcher *fd.EventDispatcher

func init() {
	util.Module("app", new(Module))
}

func (this *Module) Start() {
	// fmt.Println("start")
}

func (this *Module) Print(a string, b int) {
	log.Info("App Print ", a, b)
}

func (this *Module) IsEnable() bool {
	return true
}
