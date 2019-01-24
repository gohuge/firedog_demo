package main

import (
	"fmt"
	"os"
	"path"
	"time"
)

type key struct {
	offset int
	size   int
	vsn    int
}

type storage struct {
	keys   map[string]key
	values map[string]interface{}
}

func (this *storage) Load(pathname string) {
	file, err := os.OpenFile(path.Join(pathname, "table.dat"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	//fin, err := os.Open(path)
	//defer fin.Close()
	if err != nil {
		fmt.Println(pathname, err)
		return
	}
	timeUnixNano := time.Now().UnixNano()
	buf := make([]byte, 1024*1024*6)
	count := 0
	for {
		n, _ := file.Read(buf)
		if 0 == n {
			break
		}
		count = count + n
		//os.Stdout.Write(buf[:n])
	}
	fmt.Println("size  ... ", count, (time.Now().UnixNano()-timeUnixNano)/1000000)
}

func main() {
	s := &storage{}
	s.Load("./")
}
