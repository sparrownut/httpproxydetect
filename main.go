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
var threadsMax = 500

func main() {
	port := "8080"
	protocol := "http"
	app := &cli.App{
		Name:      "protocaldetect",
		Usage:     "judg protocol\n protocol:\nhttp\nssh\nmysql\nshiro\nyonyou", // 这里写协议
		UsageText: "lazy to write...",
		Version:   "0.5.9",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Destination: &port, Value: "8080", Usage: "port", Required: true},
			&cli.StringFlag{Name: "protocol", Aliases: []string{"P"}, Destination: &protocol, Value: "ssh", Usage: "protocol", Required: true},
			&cli.IntFlag{Name: "threads", Aliases: []string{"T"}, Destination: &threadsMax, Value: 500, Usage: "DBG MOD", Required: false},
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
waitToRun:
	if DBG {
		println(fmt.Sprintf("当前进程%v个", threads))
	}
	if threads <= threadsMax {
		go func() {
			_ = dofunc(port, protocol, file, host)
		}()
	} else {
		goto waitToRun
	}
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
		if dial == nil {
			return nil
		}
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
		if dial == nil {
			return nil
		}
		_ = dial.SetReadDeadline(time.Now().Add(timeout))
		buf := [64]byte{}
		n, _ := dial.Read(buf[:])
		if strings.Contains(string(buf[:n]), "mysql") {
			fmt.Printf(host + "\n")
			_, _ = file.WriteString("[MYSQL]" + host + ":" + port + "\n")

		}
		_ = dial.Close()
	} else if protocol == "shiro" {
		if DBG {
			println("shiro MOD")
		}
		client := goreq.NewClient()
		var req *goreq.Request
		if port == "80" { //判断加密
			req = goreq.Get(fmt.Sprintf("http://%v", host)).SetClient(client)
		} else if port == "443" {
			req = goreq.Get(fmt.Sprintf("https://%v", host)).SetClient(client)
		} else {
			req = goreq.Get(fmt.Sprintf("http://%v:%v", host, port)).SetClient(client)
		}
		ret := goreq.Do(req)
		if DBG {
			println(ret.Text)
		}
		if strings.Contains(ret.Header.Get("Set-Cookie"), "rememberMe") {
			fmt.Printf("%v\n", host)
		}

	} else if protocol == "yonyou" {
		if DBG {
			println("yonyou MOD")
		}
		client := goreq.NewClient()
		var req *goreq.Request

		host = strings.ReplaceAll(host, "https://", "") //过滤http前缀
		host = strings.ReplaceAll(host, "http://", "")

		suffix := "/servlet/~ic/bsh.servlet.BshServlet"
		if strings.Contains(host, ":") { // 如果输入有端口
			req = goreq.Get(fmt.Sprintf("http://%v%v", host, suffix)).SetClient(client)
		} else {
			req = goreq.Get(fmt.Sprintf("http://%v:%v%v", host, port, suffix)).SetClient(client)
		}
		ret := goreq.Do(req)
		if DBG {
			println(ret.Text)
		}
		if strings.Contains(strings.ToUpper(strings.ReplaceAll(ret.Text, " ", "")), "BEANSHELLTESTSERVLET") {
			fmt.Printf("%v\n", host)
		}

	} else {
		println("无此协议")
	}

	return nil
}
