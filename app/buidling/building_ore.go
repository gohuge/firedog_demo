package buidling

import (
	"fmt"
	"github.com/gohuge/firedog/fd"
)

type BuildingOre struct {
	room *fd.Room
}

func (this *BuildingOre) New() {
	this.room = &fd.Room{}
	this.room.Name = "BuildingOre"
	this.room.OnInit = this.OnInit
	this.room.Handler = this.Handler
	this.room.Close = this.Close
}

func (this *BuildingOre) OnInit(r *fd.Room, args []interface{}) {
	fmt.Println(" Init ", r.Name, args)
}

func (this *BuildingOre) Handler(r *fd.Room, msg string, args []interface{}, ret chan []interface{}) {
	fmt.Println(" Handler ", r.Name, msg, args)

	if "ttt" == msg {
		ret <- nil
	}
}

func (this *BuildingOre) Close(r *fd.Room, reason string) {
	fmt.Println(" Close ", r.Name, reason)
}
