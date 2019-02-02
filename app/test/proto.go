package main;
 
import (
    "github.com/golang/protobuf/proto"
    "../src/pb"
    "reflect"
    "fmt"
)
 
func buftest() {
    p1 := &pb.Task{}
    p1.Tid = 100
    p1.Number = "3000"
    //编码数据
    data, _ := proto.Marshal(p1)
    //把数据写入文件
    fmt.Println(" write : ",data,proto.MessageType("pb.task"))

    p2 := &pb.Task{}
    proto.Unmarshal(data, p2);

    mt := proto.MessageType("pb.task")
    m := reflect.New(mt.Elem())
    // p3.Tid = 300
    // p3 := new(proto.Message)
    // p3 := proto.DecodeMessage(data)
    // fmt.
     fmt.Println(" interface : ",m.Interface().(proto.Message))
    p3 :=  m.Interface().(proto.Message)
    proto.Unmarshal(data,p3);


    fmt.Println(" read : ",p3,m.Interface())

}
 
func main() {
    buftest();
}