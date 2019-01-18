package main

import (
	"fmt"
	"unsafe"
)

type Home struct {
	address string
}

type Object struct {
	data int64
	name string
	Home *Home
}

type Post struct {
	port uint16
	from string
	info *interface{}
}

type data struct {
	addr uintptr
	len  int
	cap  int
}

func (this *Post) Set(port uint16, from string, info interface{}) {
	this.port = port
	this.from = from
	this.info = &info
}

func (this *Post) ToBytes() []byte {
	Len := unsafe.Sizeof(*this)
	sd := &data{
		addr: uintptr(unsafe.Pointer(this)),
		cap:  int(Len),
		len:  int(Len),
	}
	return *(*[]byte)(unsafe.Pointer(sd))
}

func (this *Post) Parse(data []byte) *Post {
	var p *Post = *(**Post)(unsafe.Pointer(&data))
	return p
}

func main() {

	//var obj = &Object{100,"gouhongjie",&Home{"chengdu"}}
	post := &Post{}
	post.Set(111, "from xxx", nil)
	data := post.ToBytes()
	fmt.Println("[1]byte is : ", data)

	fmt.Println("ptestStruct.data is : ", post.Parse(data).info)
}
