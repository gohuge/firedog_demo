package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

const (
	//------------------测试数据配置-------------------//
	//在线玩家数
	defaultACU = 2000
	//一个玩家的数据总量，单位：KB，该值影响执行redis命令时value的大小
	defaultOneUserDataTotal = 20
	//一个玩家的数据条目数量，该值影响执行redis命令时key的数量
	defaultOneUserDataCount = 100
	//与redis服务器连接的连接数，由于一个协程处理一个连接，因此该值也是创建的协程数
	//由于协程数影响到数据并发处理，因此
	//需满足 ACU * OneUserDataCount / RedisConnectCount 为正整数
	defaultRedisConnectCount = 10
	//测试数据存储路径，若文件不存在，则随机生成一份数据并写入到该文件中
	TestDataFileName = "./test-data.txt"

	//---------------redis服务器连接设置---------------//
	//连接方式
	defaultRedisNetwork = "tcp"
	//地址
	defaultRedisAddress = "xxxx:6379"
	//验证密码
	defaultRedisPassword = "test"
	//连接超时时间
	defaultRedisConnectTimeout = time.Second * 5
	//读取超时时间
	defaultRedisReadTimeout = time.Second * 5
	//写入超时时间
	defaultRedisWriteTimeout = time.Second * 5

	//------------------程序启动参数-------------------//
	//默认执行mget指令
	defaultTestCmd = "mset"
	//干扰协程的key数量，0表示不进行干扰
	defaultNoiseKeyCount = 0
	//是否总是随机生成数据，默认为false
	defaultIsGen = false
	//是否打印详细信息，默认为false
	defaultIsPrint = false
)

var (
	ACU               int
	OneUserDataTotal  int
	OneUserDataCount  int
	RedisConnectCount int

	TestCmd       string
	NoiseKeyCount int
	IsGen         bool
	IsPrint       bool

	//启动参数，使用方法：在程序启动时加上 -a xxx -ot xxx -oc xxx -cc xxx -t xxx -n xxx -g=xxx -p=xxx
	inputACU               = flag.Int("a", defaultACU, "ACU")
	inputOneUserDataTotal  = flag.Int("ot", defaultOneUserDataTotal, "OneUserDataTotal")
	inputOneUserDataCount  = flag.Int("oc", defaultOneUserDataCount, "OneUserDataCount")
	inputRedisConnectCount = flag.Int("cc", defaultRedisConnectCount, "RedisConnectCount")
	inputTestCmd           = flag.String("t", defaultTestCmd, "TestCmd")
	inputNoiseKeyCount     = flag.Int("n", defaultNoiseKeyCount, "NoiseKeyCount")
	inputIsGen             = flag.Bool("g", defaultIsGen, "IsGen")
	inputIsPrint           = flag.Bool("p", defaultIsPrint, "IsPrint")
)

func mainq() {
	//初始化命令参数
	InitArg()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//运行干扰协程，使redis服务器处于忙碌状态
		WorkForever(NoiseKeyCount)
	}()
	wg.Wait()
}

func main() {
	//初始化命令参数
	InitArg()
	//当前cpu数
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Println("current cpu number:", numCPU)

	var wg sync.WaitGroup
	var MaxTime, MinTime, AverageTime time.Duration
	conns := make([]*redis.Conn, RedisConnectCount)
	goroutinesRunTime := make([]time.Duration, RedisConnectCount)

	//创建连接
	for i := 0; i < RedisConnectCount; i++ {
		conns[i] = Connect()
	}
	//延迟关闭连接
	defer func() {
		for i := 0; i < RedisConnectCount; i++ {
			(*conns[i]).Close()
		}
	}()
	//生成测试数据
	testData := genUserData()
	fmt.Printf("generate testData size(byte): %d, len: %d\n", sumDataSize(testData), len(testData))
	//定义协程执行函数
	workerFunc := func(index int, workerType string) {
		defer wg.Done()
		startTime := time.Now()

		switch workerType {
		case "mset", "MSET":
			WorkMSET(index, &testData, conns[index])
		case "mget", "MGET":
			WorkMGET(index, conns[index])
		default:
			fmt.Println("worker type err. type:", workerType)
		}

		goroutinesRunTime[index] = time.Now().Sub(startTime)

		if IsPrint {
			fmt.Println(workerType, "worker id:", index, "finished. time:", goroutinesRunTime[index])
		}
	}

	//开始执行协程函数并等待执行完毕
	wg.Add(RedisConnectCount)
	for i := 0; i < RedisConnectCount; i++ {
		go workerFunc(i, TestCmd)
	}
	wg.Wait()
	//计算最小时间、最大时间、平均时间
	MinTime, MaxTime, AverageTime = goroutinesRunTime[0], time.Duration(0), time.Duration(0)
	for i := 0; i < RedisConnectCount; i++ {
		if goroutinesRunTime[i] > MaxTime {
			MaxTime = goroutinesRunTime[i]
		}
		if goroutinesRunTime[i] < MinTime {
			MinTime = goroutinesRunTime[i]
		}
		AverageTime += goroutinesRunTime[i]
	}
	fmt.Println("all work finished. worker number:", RedisConnectCount, "TestCmd:", TestCmd, "MaxTime:", MaxTime, "MinTime:", MinTime, "AverageTime:", AverageTime/time.Duration(RedisConnectCount))
}

//初始化命令参数
func InitArg() {
	flag.Parse()

	ACU = *inputACU
	OneUserDataTotal = *inputOneUserDataTotal * 1000
	OneUserDataCount = *inputOneUserDataCount
	RedisConnectCount = *inputRedisConnectCount
	TestCmd = *inputTestCmd
	NoiseKeyCount = *inputNoiseKeyCount
	IsGen = *inputIsGen
	IsPrint = *inputIsPrint
}

//创建一个连接redis服务器的连接
func Connect() (conn *redis.Conn) {
	c, err := redis.Dial(defaultRedisNetwork, defaultRedisAddress,
		redis.DialPassword(defaultRedisPassword),
		redis.DialConnectTimeout(defaultRedisConnectTimeout),
		redis.DialReadTimeout(defaultRedisReadTimeout),
		redis.DialWriteTimeout(defaultRedisWriteTimeout))
	if err != nil {
		fmt.Println("Connect to redis error:", err)
		return &c
	}
	return &c
}

//MSET
func WorkMSET(index int, testData *[][]byte, conn *redis.Conn) {
	iter := ACU * OneUserDataCount / RedisConnectCount
	offset := index * iter
	buf := make([]interface{}, iter*2)
	for i := offset; i < offset+iter; i++ {
		buf[(i-offset)*2] = i
		buf[(i-offset)*2+1] = (*testData)[i]
	}
	if IsPrint {
		fmt.Println(index, "WorkMSET cmd args len:", len(buf))
	}

	_, err := (*conn).Do("MSET", buf...)
	if err != nil {
		fmt.Println(index, "WorkSET MSET to redis error:", err)
	}
}

//MGET
func WorkMGET(index int, conn *redis.Conn) {
	iter := ACU * OneUserDataCount / RedisConnectCount
	offset := index * iter
	buf := make([]interface{}, iter)
	for i := offset; i < offset+iter; i++ {
		buf[i-offset] = i
	}
	if IsPrint {
		fmt.Println(index, "WorkMGET cmd args len:", len(buf))
	}

	_, err := (*conn).Do("MGET", buf...)
	if err != nil {
		fmt.Println(index, "WorkGET MGET to redis error:", err)
	}
}

func WorkForever(keyCount int) {
	if keyCount < 1 {
		return
	}

	var keys []interface{}
	conn := Connect()
	defer (*conn).Close()
	for i := 0; i < keyCount; i++ {
		keys = append(keys, randByteSlice(10), randByteSlice(100))
	}

	fmt.Println("WorkForever start...")

	for i := 0; ; i++ {
		_, err := (*conn).Do("MSET", keys...)
		if err != nil {
			fmt.Println("WorkForever MSET to redis error:", err)
			break
		}
		if IsPrint {
			fmt.Println(i, "turn WorkForever finished.")
		}
	}
}

func PipeTest(c *redis.Conn, testData [][]byte) {
	var startTime time.Time
	var valueBuf [][]byte
	var ret [][]byte

	//pipe Testing
	startTime = time.Now()
	PipeSET(c, testData)
	fmt.Println("PipeSET Time:", time.Now().Sub(startTime))

	startTime = time.Now()
	ret = PipeGET(c, testData)
	fmt.Println("PipeGET Time:", time.Now().Sub(startTime))
	fmt.Println("PipeGET data length:", sumDataSize(ret))
	for _, v := range testData {
		valueBuf = append(valueBuf, v)
	}
	fmt.Println("PipeGET data check:", checkData(&valueBuf, &ret))

	startTime = time.Now()
	PipeDEL(c, testData)
	fmt.Println("PipeDEL Time:", time.Now().Sub(startTime))

	fmt.Println("Testing all over.")
}

func PipeSET(c *redis.Conn, testData [][]byte) {
	for i := range testData {
		(*c).Send("SET", string(i), testData[i])
	}
	(*c).Flush()
	for range testData {
		_, err := (*c).Receive()
		if err != nil {
			fmt.Println("PipeSET to redis error:", err)
		}
	}
}

func PipeGET(c *redis.Conn, testData [][]byte) [][]byte {
	var buf [][]byte
	for i := range testData {
		(*c).Send("GET", string(i))
	}
	(*c).Flush()
	for range testData {
		ret, err := redis.Bytes((*c).Receive())
		if err != nil {
			fmt.Println("PipeGET to redis error:", err)
		}
		buf = append(buf, ret)
	}
	return buf
}

func PipeDEL(c *redis.Conn, testData [][]byte) {
	for i := range testData {
		(*c).Send("DEL", string(i))
	}
	(*c).Flush()
	for range testData {
		_, err := (*c).Receive()
		if err != nil {
			fmt.Println("PipeDEL to redis error:", err)
		}
	}
}

//生成一份测试数据，若已有测试数据文件则读取文件，否则生成一份数据并存入文件
func genUserData() [][]byte {
	var buf [][]byte
	var f *os.File
	var err error

	if IsGen {
		fmt.Printf("data gen...\n")
		return newUserData()
	}

	defer f.Close()
	if checkFileExist(TestDataFileName) {
		f, err = os.Open(TestDataFileName)
		if err != nil {
			fmt.Printf("open %s err: %s\n", TestDataFileName, err)
			return buf
		}

		s := bufio.NewScanner(f)
		for s.Scan() {
			line := s.Bytes()
			buf = append(buf, line)
		}
		//r := bufio.NewReader(f)
		//for {
		//  line, err := r.ReadSlice(byte('\n'))
		//  if err != nil {
		//    fmt.Printf("read %s err: %s\n", TestDataFileName, err)
		//    break
		//  }
		//  buf = append(buf, line)
		//}
		if sumDataSize(buf) == ACU*OneUserDataTotal && len(buf) == ACU*OneUserDataCount {
			return buf
		}

		f.Close()
		err := os.Remove(TestDataFileName)
		if err != nil {
			fmt.Printf("remove %s err: %s\n", TestDataFileName, err)
		}
		fmt.Printf("remove %s ok.\n", TestDataFileName)
	}

	f, err = os.Create(TestDataFileName)
	buf = newUserData()
	w := bufio.NewWriter(f)
	for _, v := range buf {
		_, err := w.Write(v)
		if err != nil {
			fmt.Printf("write %s err: %s\n", TestDataFileName, err)
			break
		}
		err = w.WriteByte(byte('\n'))
		if err != nil {
			fmt.Printf("write byte %s err: %s\n", TestDataFileName, err)
			break
		}
	}
	w.Flush()

	return buf
}

//生成ACU份数据
func newUserData() [][]byte {
	var buf [][]byte
	for i := 0; i < ACU; i++ {
		tmp := newOneUserData()
		buf = append(buf, tmp...)
	}
	return buf
}

//生成一份数据，数据总大小固定，内容随机
func newOneUserData() [][]byte {
	buf := make([][]byte, OneUserDataCount)
	dataSize := randFixedSum(OneUserDataCount, OneUserDataTotal)
	for i := range buf {
		buf[i] = randByteSlice(dataSize[i])
	}
	return buf
}

//生成数量和总和固定的随机数组
func randFixedSum(count, sum int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]int, count)
	for i := 0; i < count-1; i++ {
		buf[i] = r.Intn(sum)
	}
	buf[count-1] = sum
	sort.Ints(buf)
	for i := count - 1; i > 0; i-- {
		buf[i] -= buf[i-1]
	}
	return buf
}

//生成给定大小的随机byte数组
func randByteSlice(size int) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = byte(r.Intn(127-33) + 33)
	}
	return buf
}

//计算数据总大小
func sumDataSize(testData [][]byte) int {
	sum := 0
	for i := range testData {
		sum += len(testData[i])
	}
	return sum
}

//检查两个测试数据是否相同
func checkData(a, b *[][]byte) bool {
	for i := 0; i < len(*a); i++ {
		if !bytes.Equal((*a)[i], (*b)[i]) {
			return false
		}
	}
	return true
}

//检查文件是否存在
func checkFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
