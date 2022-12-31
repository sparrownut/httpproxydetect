package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/zhshch2002/goreq"
	"net"
	"os"
	"strings"
	"time"
)

var DBG bool
var threads = 0

func main() {
	port := "8080"
	protocol := "http"
	app := &cli.App{
		Name:      "protocaldetect",
		Usage:     "judg protocol\n protocol:\nhttp\nssh\nmysql", // 这里写协议
		UsageText: "lazy to write...",
		Version:   "0.4.8",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Destination: &port, Value: "8080", Usage: "port", Required: true},
			&cli.StringFlag{Name: "protocol", Aliases: []string{"P"}, Destination: &protocol, Value: "ssh", Usage: "protocol", Required: true},
			&cli.BoolFlag{Name: "DBG", Aliases: []string{"D"}, Destination: &DBG, Value: false, Usage: "DBG MOD", Required: false},
		},
		HideHelpCommand: true,
		Action: func(c *cli.Context) error {
			err := do(port, protocol)
			if err != nil {

			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		//panic(err)
	}

	//fmt.Printf(os.Args[1])

}
func do(port string, protocol string) error {
	file, _ := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	// 文件初始化
start: // 在这里循环
	host := ""
	_, _ = fmt.Scanln(&host)

	go func() {

		_ = dofunc(port, protocol, file, host)
	}()
	goto start
}
func dofunc(port string, protocol string, file *os.File, host string) error {
	threads++ // 线程+1
	defer func() {
		threads--
		if r := recover(); r != nil {
			if DBG {
				fmt.Println("recover value is", r)
				fmt.Printf("ERROR INFO host:%v", host)
			}
		}
	}() //清理线程计数 处理异常
	if DBG {
		println(fmt.Sprintf("当前进程%v个", threads))
	}
	timeout := 2 * time.Second

	if protocol == "http" {

		proxyStr := "http://" + host + ":" + port
		client := goreq.NewClient()
		req := goreq.Get("http://icanhazip.com/").SetClient(client).SetProxy(proxyStr).SetTimeout(timeout)
		//fmt.Printf(proxyStr)
		if req.Err == nil {
			ret := goreq.Do(req)
			if strings.Contains(ret.Text, host) {
				fmt.Printf(proxyStr + "\n") // 输出
				_, _ = file.WriteString(proxyStr + "\n")

			}
		}

	} else if protocol == "ssh" {
		dial, connecterr := net.Dial("tcp", host+":"+port)
		_ = dial.SetReadDeadline(time.Now().Add(timeout))
		if connecterr != nil {
		}
		buf := [64]byte{}
		n, _ := dial.Read(buf[:])

		if strings.Contains(string(buf[:n]), "SSH") {
			fmt.Printf(host + "\n")
			_, _ = file.WriteString("[SSH]" + host + ":" + port + "\n")

		}
		_ = dial.Close()
	} else if protocol == "mysql" {
		dial, _ := net.Dial("tcp", host+":"+port)
		_ = dial.SetReadDeadline(time.Now().Add(timeout))
		buf := [64]byte{}
		n, _ := dial.Read(buf[:])
		if strings.Contains(string(buf[:n]), "mysql") {
			fmt.Printf(host + "\n")
			_, _ = file.WriteString("[MYSQL]" + host + ":" + port + "\n")

		}
		_ = dial.Close()
	} else {
		println("无此协议")
	}

	return nil
}
