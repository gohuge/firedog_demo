// author：gouhongjie gohuge@qq.com
package main

import (
	_ "./app"
	_ "./app/pb"
	"github.com/gohuge/firedog/boot"
)

func init() {

}

func main() {

	go boot.Start("./app/#/")

	boot.Console(0)
}
