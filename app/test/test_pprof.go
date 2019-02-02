package main

import (
    "flag"
    "log"
    "net/http"
    _ "net/http/pprof"
    "sync"
    "time"
)
//go run test_pprof.go
// overview:http://localhost:8080/debug/pprof/
// 查看栈信息
//go tool pprof http://localhost:8080/debug/pprof/heap
// 查看CPU
//go tool pprof http://localhost:8080/debug/pprof/profile
func Counter(wg *sync.WaitGroup) {
    time.Sleep(time.Second)

    var counter int
    for i := 0; i < 1000000; i++ {
        time.Sleep(time.Millisecond * 200)
        counter++
    }
    wg.Done()
}

func main() {
    flag.Parse()

    //远程获取pprof数据
    go func() {
        log.Println(http.ListenAndServe("localhost:8080", nil))
    }()

    var wg sync.WaitGroup
    wg.Add(10)
    for i := 0; i < 10; i++ {
        go Counter(&wg)
    }
    wg.Wait()

    // sleep 10mins, 在程序退出之前可以查看性能参数.
    time.Sleep(60 * time.Second)
}
