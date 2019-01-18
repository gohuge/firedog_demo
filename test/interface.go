package main

import (
	"errors"
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

type nodeseq struct {
	port uint16
	from string
	info *interface{}
}

type seqdata struct {
	addr uintptr
	len  int
	cap  int
}

func main() {

	var obj = &Object{100, "gouhongjie", &Home{"chengdu"}}
	seq := NodeSeq(111, "from xxx", obj)
	//var msg = &Message{65534,"from me",nil}
	//msg.SetInfo(obj)
	//
	//fmt.Println("type  ",*obj,unsafe.Sizeof(*msg), unsafe.Sizeof(Message{}))
	//
	////WriteObject(&obj)
	//Len := unsafe.Sizeof(*msg)
	//testBytes := &SliceMock{
	//	addr: uintptr(unsafe.Pointer(msg)),
	//	cap:  int(Len),
	//	len:  int(Len),
	//}
	//data := *(*[]byte)(unsafe.Pointer(testBytes))
	//data2 := []byte{100,0,0,0,0,0,0,0,184,252,75,0,0,0,0,0,10,0,0,0,0,0,0,0,192,225,4,0,192,0,0,0}
	data, _ := SeqData(seq)
	fmt.Println("[1]byte is : ", data)
	//fmt.Println("[2]byte is : ", data2)
	var p *nodeseq = *(**nodeseq)(unsafe.Pointer(&data))

	fmt.Println("ptestStruct.data is : ", p.port, p.from, *p.info)
}

func NodeSeq(port uint16, from string, info interface{}) *nodeseq {
	return &nodeseq{port, from, &info}
}

func SeqData(seq *nodeseq) (data []byte, err error) {
	if seq == nil {
		err = errors.New("bad nodeseq")
	}
	Len := unsafe.Sizeof(*seq)
	sd := &seqdata{
		addr: uintptr(unsafe.Pointer(seq)),
		cap:  int(Len),
		len:  int(Len),
	}
	data = *(*[]byte)(unsafe.Pointer(sd))
	return
}
