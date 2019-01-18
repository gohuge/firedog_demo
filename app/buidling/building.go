package buidling

import (
	"fmt"
	"github.com/gohuge/firedog/fd"
)

var room *fd.Room

func init() {
	room = &fd.Room{}
	room.Name = "building"
	room.OnInit = OnInit
	room.Handler = Handler
	room.Close = Close
}

func OnInit(r *fd.Room, args []interface{}) {
	fmt.Println(" Init ", r.Name, args)
}

func Handler(r *fd.Room, msg string, args []interface{}, ret chan []interface{}) {
	fmt.Println(" Handler ", r.Name, msg, args)

	if "ttt" == msg {
		ret <- nil
	}
}

func Close(r *fd.Room, reason string) {
	fmt.Println(" Close ", r.Name, reason)
}

//func main() {
//	fd.Family.Start(room)
//
//	fd.Family.Cast("activity", "ghj", nil)
//
//	ret, err := fd.Family.Call("activity", "ttt", nil, 1000)
//
//	fmt.Println("ret ", fd.Family.IsAlive("activity"), fd.Family.Count(), ret, err)
//}
