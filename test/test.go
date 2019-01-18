package main

import (
	"fmt"
)

func f() {
	defer func() {
		fmt.Println("b")
		if err := recover(); err != nil {
			fmt.Println("ddd ", err)
		}
		fmt.Println("d")
	}()
	fmt.Println("a")
	panic("a bug occur")
	fmt.Println("c")
}

func main() {
	f()
	fmt.Println("x")
}
