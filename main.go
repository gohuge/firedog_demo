// author：gouhongjie gohuge@qq.com
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

import (
	_ "./app"

	"github.com/gohuge/firedog/src/firedog/boot"
	//"./firedog/boot"
	_ "./pb"
)

func init() {
	fmt.Println("【APP START !!!】")
	fmt.Println("=================================================================")
	for _, arg := range os.Args[1:] {
		fmt.Println("[boot] args = ", arg)
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[boot] current cpu dir :", dir)
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Println("[boot] current cpu number :", numCPU)
	fmt.Println("=================================================================")
}

func main() {

	go boot.Start("./conf")
	boot.Console(0)
}
