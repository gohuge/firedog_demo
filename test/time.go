// +build OMIT

package main

import (
	"fmt"
	"time"
)

func main() {
	stop := time.After(3 * time.Second)
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for {
		select {
		case a:=<-tick.C:
			fmt.Println(time.Now(),"|||",a)
		case <-stop:
			return
		}
	}
}
