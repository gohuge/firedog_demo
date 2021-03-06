package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"math/rand"
	"net"
	"os"
	"time"
	// "sync"
	"../pb"
	"encoding/xml"
	"github.com/gohuge/firedog/util"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
)

//数据包类型
const (
	HEART_BEAT_PACKET = 0x00
	REPORT_PACKET     = 0x01
)

type NetClient struct {
	XMLName xml.Name `xml:"client"`
	count   int      `xml:"count,attr"` //连接数
	//next int		 `xml:"next,attr"`//下一个链接的创建毫秒
	address string   `xml:"address,attr"` //地址 127.0.0.1:3344
	sets    []MsgSet `xml:"set"`          // 消息设置
}

type MsgSet struct {
	XMLName xml.Name `xml:"msg"`
	typ     int      `xml:"type,attr"` //类型，循环间隔还是定时间隔，1，2
	time    int      `xml:"time,attr"` //时间，秒
	info    []byte   `xml:"info,attr"` //数据
	file    string   `xml:"file,attr"` //文件
}

//数据包
type Packet struct {
	PacketType    byte
	PacketContent []byte
}

//心跳包
type HeartPacket struct {
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

//数据包
type ReportPacket struct {
	Content   string `json:"content"`
	Rand      int    `json:"rand"`
	Timestamp int64  `json:"timestamp"`
}

//客户端对象
type TcpClient struct {
	connection *net.TCPConn
	hawkServer *net.TCPAddr
}

func main() {
	file, err := os.Open("./net.xml")
	if err != nil {
		fmt.Printf("error1: %v", err)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error2: %v", err)
		return
	}
	c := NetClient{}
	err = xml.Unmarshal(data, &c)
	if err != nil {
		fmt.Printf("error3: %v", err)
	}
	print("xml:", c)
	ticker := time.NewTicker(1 * time.Second)
	var Count = 0
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
			NewConnect(c)
			Count++
			if Count >= c.count {
				ticker.Stop()
				print("Net Creat Over Count = ", Count)
			}
		}
	}()
	// 发送消息
	ticker2 := time.NewTicker(1 * time.Second)
	go func() {
		for t := range ticker2.C {

		}
	}()
	<-stopChan
}

func NewConnect(v NetClient) {
	//拿到服务器地址信息
	hawkServer, err := net.ResolveTCPAddr("tcp", v.address)
	if err != nil {
		fmt.Printf("hawk server [%s] resolve error: [%s]", v.address, err.Error())
		os.Exit(1)
	}
	//连接服务器
	connection, err := net.DialTCP("tcp", nil, hawkServer)
	if err != nil {
		fmt.Printf("connect to hawk server error: [%s]", err.Error())
		os.Exit(1)
	}
	client := &TcpClient{
		connection: connection,
		hawkServer: hawkServer,
	}
	//启动接收
	go client.receivePackets()

	//发送心跳的goroutine
	// go func() {
	//     heartBeatTick := time.Tick(2 * time.Second)
	//     for{
	//         select {
	//         case <-heartBeatTick:
	//             client.sendHeartPacket()
	//         case <-client.stopChan:
	//             return
	//         }
	//     }
	// }()

	//测试用的，开300个goroutine每秒发送一个包
	// for i:=0;i<300;i++ {
	//     go func() {
	//         sendTimer := time.After(1 * time.Second)
	//         for{
	//             select {
	//             case <-sendTimer:
	//                 client.sendReportPacket()
	//                 sendTimer = time.After(1 * time.Second)
	//             case <-client.stopChan:
	//                 return
	//             }
	//         }
	//     }()
	// }
	for i := 0; i < 1; i++ {
		go func() {
			sendTimer := time.After(10 * time.Second)
			for {
				select {
				case <-sendTimer:
					client.sendProto()
					sendTimer = time.After(10 * time.Second)
				}
			}
		}()
	}
}

// 接收数据包
func (client *TcpClient) receivePackets() {
	reader := bufio.NewReader(client.connection)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(msg)
	}
}

func (client *TcpClient) sendReportPacket() {
	reportPacket := ReportPacket{
		Content:   getRandString(),
		Timestamp: time.Now().Unix(),
		Rand:      rand.Int(),
	}
	packetBytes, err := json.Marshal(reportPacket)
	if err != nil {
		fmt.Println(err.Error())
	}
	packet := Packet{
		PacketType:    REPORT_PACKET,
		PacketContent: packetBytes,
	}
	sendBytes, err := json.Marshal(packet)
	if err != nil {
		fmt.Println(err.Error())
	}
	client.connection.Write(EnPackSendData(sendBytes))
	fmt.Println("Send metric data success!", sendBytes)
}

func (client *TcpClient) sendProto() {
	p1 := &pb.Task{}
	p1.Tid = 100
	p1.Number = "3000"

	//编码数据
	data, err := proto.Marshal(p1)
	if err != nil {
		fmt.Println(err.Error())
	}

	port := 1234

	buf := util.NewBuffer(nil)
	buf.WriteU16(uint16(port))
	buf.WriteBytes(data)

	msgLen := buf.Length()
	bytes := make([]byte, msgLen+2)
	bytes[0] = byte(msgLen >> 8 & 0xff)
	bytes[1] = byte(msgLen & 0xff)

	copy(bytes[2:], buf.Data())
	//
	////发送
	client.connection.Write(bytes)
	fmt.Println("Send x metric data success!", buf.Length(), msgLen, len(bytes), bytes)
	fmt.Println("Send x metric data success!", msgLen, data)
}

//使用的协议与服务器端保持一致
func EnPackSendData(sendBytes []byte) []byte {
	packetLength := len(sendBytes) + 8
	result := make([]byte, packetLength)
	result[0] = 0xFF
	result[1] = 0xFF
	result[2] = byte(uint16(len(sendBytes)) >> 8)
	result[3] = byte(uint16(len(sendBytes)) & 0xFF)
	copy(result[4:], sendBytes)
	sendCrc := crc32.ChecksumIEEE(sendBytes)
	result[packetLength-4] = byte(sendCrc >> 24)
	result[packetLength-3] = byte(sendCrc >> 16 & 0xFF)
	result[packetLength-2] = 0xFF
	result[packetLength-1] = 0xFE
	fmt.Println(result)
	return result
}

//发送心跳包，与发送数据包一样
func (client *TcpClient) sendHeartPacket() {
	heartPacket := HeartPacket{
		Version:   "1.0",
		Timestamp: time.Now().Unix(),
	}
	packetBytes, err := json.Marshal(heartPacket)
	if err != nil {
		fmt.Println(err.Error())
	}
	packet := Packet{
		PacketType:    HEART_BEAT_PACKET,
		PacketContent: packetBytes,
	}
	sendBytes, err := json.Marshal(packet)
	if err != nil {
		fmt.Println(err.Error())
	}
	client.connection.Write(EnPackSendData(sendBytes))
	fmt.Println("Send heartbeat data success!")
}

//拿一串随机字符
func getRandString() string {
	length := rand.Intn(10)
	strBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		strBytes[i] = byte(rand.Intn(26) + 97)
	}
	return string(strBytes)
}
