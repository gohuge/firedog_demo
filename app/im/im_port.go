package im

import (
	"fmt"
	"github.com/gohuge/firedog/fd"
)

func init() {
	fd.SetHandler("im_register", RegisterHandler)
	fd.SetHandler("im_login", LoginHandler)
	fd.SetHandler("im_exit", ExitHandler)
	fd.SetHandler("im_message", MessageHandler)
	fd.SetHandler("im_relation", RelationHandler)
	fd.SetHandler("im_room", RoomHandler)
}

func RegisterHandler(session *fd.Session, msg ...interface{}) {

	fmt.Println("logic ", session, msg)

}

func LoginHandler(session *fd.Session, msg ...interface{}) {

	fmt.Println("logic ", session, msg)

}

func ExitHandler(session *fd.Session, msg ...interface{}) {

	fmt.Println("logic ", session, msg)

}

func MessageHandler(session *fd.Session, msg ...interface{}) {

	fmt.Println("logic ", session, msg)

}

func RelationHandler(session *fd.Session, msg ...interface{}) {

	fmt.Println("logic ", session, msg)

}

func RoomHandler(session *fd.Session, msg ...interface{}) {

	fmt.Println("logic ", session, msg)

}