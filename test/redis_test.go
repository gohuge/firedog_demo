package main

import (
	"../src/firedog/db"
	"errors"
	"fmt"
	"testing"
	"time"
)

func Division(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为0")
	}

	return a / b, nil
}

func Test_Division_1(t *testing.T) {
	if i, e := Division(6, 2); i != 3 || e != nil { //try a unit test on function
		t.Error("除法函数测试没通过") // 如果不是如预期的那么就报错
	} else {
		t.Log("第一个测试通过了") //记录一些你期望记录的信息
	}
}

func Test_Division_2(t *testing.T) {
	t.Error("就是不通过")
}

func Test_ont(t *testing.T) {
	fmt.Println("test 1")
	t.Log("a")
}
func Test_tow(t *testing.T) {
	fmt.Println("test 2")

	t.Log("b")
}

func Test_redis(t *testing.T) {
	redisCfg := db.RedisCfg{
		MaxIdle:     8,
		MaxActive:   64,
		IdleTimeout: 300,
		Address:     "x",
		Password:    "",
		DbNum:       1,
	}
	fmt.Println("Test Start ")
	bm := db.Redis.Init(redisCfg)
	timeoutDuration := 200 * time.Second
	//
	//c, err := redis.Dial("tcp", redisCfg.RedisServer)
	//t.Log("connect ",c)
	//t.Log("err ",err)
	//var err error
	//r1,e1:=c.Do("set", "ghj",3231313)
	err := bm.Put("cctv", 1, timeoutDuration)
	t.Log("put ", err, bm.Get("cctv"))
	//r2,e2:=c.Do("get","ghj")
	//t.Log("ret ",r1,e1,r2,e2)
	//if err = bm.Put("ghj", 1, timeoutDuration); err != nil {
	//	t.Error("set Error", err)
	//}
	//if !bm.IsExist("ghj") {
	//	t.Error("check err")
	//}
	//time.Sleep(3 * time.Second)
	//if bm.IsExist("ghj") {
	//	t.Error("check err")
	//}
	//if err = bm.Put("ghj", 1, timeoutDuration); err != nil {
	//	t.Error("set Error", err)
	//}
	//if v, _ := redis.Int(bm.Get("ghj"), err); v != 1 {
	//	t.Error("get err")
	//}
	//if err = bm.Incr("ghj"); err != nil {
	//	t.Error("Incr Error", err)
	//}
	//if v, _ := redis.Int(bm.Get("ghj"), err); v != 2 {
	//	t.Error("get err")
	//}
	//if err = bm.Decr("ghj"); err != nil {
	//	t.Error("Decr Error", err)
	//}
	//if v, _ := redis.Int(bm.Get("ghj"), err); v != 1 {
	//	t.Error("get err")
	//}
	//bm.Delete("ghj")
	//if bm.IsExist("ghj") {
	//	t.Error("delete err")
	//}
	////test string
	//if err = bm.Put("ghj", "author", timeoutDuration); err != nil {
	//	t.Error("set Error", err)
	//}
	//if !bm.IsExist("ghj") {
	//	t.Error("check err")
	//}
	//if v, _ := redis.String(bm.Get("ghj"), err); v != "author" {
	//	t.Error("get err")
	//}
	////test GetMulti
	//if err = bm.Put("ghj", "author1", timeoutDuration); err != nil {
	//	t.Error("set Error", err)
	//}
	//if !bm.IsExist("ghj") {
	//	t.Error("check err")
	//}
	//vv := bm.GetMulti([]string{"ghj", "ghj1"})
	//if len(vv) != 2 {
	//	t.Error("GetMulti ERROR")
	//}
	//if v, _ := redis.String(vv[0], nil); v != "author" {
	//	t.Error("GetMulti ERROR")
	//}
	//if v, _ := redis.String(vv[1], nil); v != "author1" {
	//	t.Error("GetMulti ERROR")
	//}
	//fmt.Println(redis.String(bm.Hget("1", "a"), nil))
	//fmt.Println(bm.Hset("1", "b", "123123"))
	//fmt.Println(bm.Hset("1", "c", 999999))
	//fmt.Println(bm.Hset("1", "d", 123))
	//for _, v := range bm.HgetMulti("1", []interface{}{"a", "b", "c"}) {
	//	fmt.Println(redis.String(v, nil))
	//}
	// test clear all
	//	if err = bm.ClearAll(); err != nil {
	//		t.Error("clear all err")
	//	}
	//fmt.Println("开启一个事务操作")
	//bm.Put("test", 203, db.INFINITE)
	//bm.Do("MULTI")
	//v, _ := redis.Int(bm.Get("test"), err)
	//bm.Put("test", v+20, db.INFINITE)
	//bm.Do("EXEC")
	//fmt.Println("开启一个乐观锁事务操作")
	//bm.Do("WATCH", "test")
	//v, _ = redis.Int(bm.Get("test"), err)
	//bm.Do("MULTI")
	//bm.Put("test", v-30, db.INFINITE)
	//bm.Do("EXEC")
	//v, _ = redis.Int(bm.Get("test"), err)
	//fmt.Println("使用Transaction来操作")
	//bm.Hset("hehe", "age", 123)
	//fmt.Println("---------hehe----------", bm.Get("enen1"))
	//bm.Put("enen1", 123, db.INFINITE)
	//fmt.Println("---------enen1----------", bm.Get("enen1"))
	////使用封装的事务
	//bm.Transaction(func() (int, error) {
	//	// check
	//	a, _ := redis.Uint64(bm.Get("enen1"), nil)
	//	fmt.Println(a, a > 123)
	//	if a > 123 {
	//		return 1, fmt.Errorf("eeeee")
	//	}
	//	return 0, nil
	//}, func() {
	//	fmt.Println("2222222222222222", bm.Get("enen1"))
	//	bm.Send("PUT", "enen1", 1234)
	//	fmt.Println("2222222222222222", bm.Get("enen1"))
	//}, "enen1")

}
