package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math"

	"../src/firedog/util"
)

func main() {

	sendBytes := []byte("abcde")
	fmt.Println("sendBytes : ", sendBytes)

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

	fmt.Println("byte hello", result, byte(uint16(412)>>8), byte(uint16(412)&0xFF))

	n := sendCrc
	// 大端数据测试
	bytes := make([]byte, 4)
	bytes[0] = byte(n >> 24 & 0xff)
	bytes[1] = byte(n >> 16 & 0xff)
	bytes[2] = byte(n >> 8 & 0xff)
	bytes[3] = byte(n & 0xff)

	// binary.Read(b_buf, binary.BigEndian, &x)

	x := binary.BigEndian.Uint32(bytes)
	fmt.Println("byte : ", bytes, n, x, math.Float64bits(333.3))
	// b_buf  : =  bytes .NewBuffer(b)

	// 小端数据测试
	bytes[0] = byte(n & 0xff)
	bytes[1] = byte(n >> 8 & 0xff)
	bytes[2] = byte(n >> 16 & 0xff)
	bytes[3] = byte(n >> 24 & 0xff)

	x = binary.LittleEndian.Uint32(bytes)
	fmt.Println("small  : ", bytes, n, x)

	// 16位
	nn := 823
	bytes = make([]byte, 4)
	bytes[2] = byte(nn >> 8 & 0xff)
	bytes[3] = byte(nn & 0xff)
	// bytes[2] = byte(nn & 0xff)
	//   bytes[3] = 0xFF

	x = binary.BigEndian.Uint32(bytes)
	fmt.Println("16 test  : ", bytes, nn, x)
	// binary.LittleEndian.PutUint64(bytes, bits)
	// fmt.Println("byte : ",bytes,n,x,math.Int32f)

	// 复制测试
	port := 65533
	msgLen := 2312312
	bytes = make([]byte, 5)
	bytes[0] = byte(port >> 8 & 0xff)
	bytes[1] = byte(port & 0xff)

	bytes[2] = byte(msgLen >> 16 & 0xff)
	bytes[3] = byte(msgLen >> 8 & 0xff)
	bytes[4] = byte(msgLen & 0xff)

	var tmp = make([]byte, 4)

	copy(tmp[2:4], bytes[:2])

	n1 := binary.BigEndian.Uint32(tmp)

	copy(tmp[1:4], bytes[2:5])
	n2 := binary.BigEndian.Uint32(tmp)

	fmt.Println("port msglen  : ", n1, n2, bytes[:2], bytes[2:5])

	b := util.NewBuffer(nil)
	b.WriteByte(12)
	b.WriteString("gouhongjie")

	p := b.Pos()
	v, _ := b.ReadByte()
	s, _ := b.ReadString()

	fmt.Println(" ... ", b.Length(), p, b.Pos(), v, s)

	fmt.Println("bytes:", []byte("Jst a test!\r\n"), []byte("abc"), []byte("汉字"))

	bytes = make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = byte(i)
	}
	pos := 4

	data := make([]byte, len(bytes)-pos)
	copy(data, bytes[pos:])
	fmt.Println("bytes : ", bytes, " data : ", data)
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}
