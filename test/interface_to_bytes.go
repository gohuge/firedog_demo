package main

import(
	"fmt"
	//"github.com/gohuge/firedog/util"
	"github.com/gohuge/firedog/fd"
)


func main(){
	msg := fd.NewMsg(255)
	msg.Info["name"]="gouhongjie"
	msg.Info["age"]= 28
	buf := msg.ToBuffer()

	fmt.Println("msg ",buf.Length())

	// ********************* Unmarshal *********************
	msg2,err := fd.ReadMsg(buf)

	fmt.Println("msg2 typ ",msg2.Typ,err)
	fmt.Println("msg2 info ",msg2.Info)

}