package main

import (
	"github.com/zhshch2002/goreq"
	"os"
	"strings"
	"time"
)

func main() {
	//fmt.Printf(os.Args[1])
	proxyStr := "http://" + os.Args[1] + "8080"
	client := goreq.NewClient()
	req := goreq.Get("https://www.baidu.com").SetClient(client).SetProxy(proxyStr).SetTimeout(5 * time.Second)
	if req.Err != nil {
		return
	}
	ret := goreq.Do(req)
	if ret.Err != nil {
		return
	}
	//fmt.Printf(ret.Text)
	if strings.Contains(ret.Text, "www.baidu.com") {
		println(proxyStr)
	}
}
