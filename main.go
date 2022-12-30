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

func main() {
	port := "8080"
	protocol := "http"
	app := &cli.App{
		Name:      "protocaldetect",
		Usage:     "judg protocol\n protocol:\nhttp\nssh\nmysql", // 这里写协议
		UsageText: "lazy to write...",
		Version:   "0.4.4",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Destination: &port, Value: "8080", Usage: "port", Required: true},
			&cli.StringFlag{Name: "protocol", Aliases: []string{"P"}, Destination: &protocol, Value: "8080", Usage: "protocol", Required: true},
		},
		HideHelpCommand: true,
		Action: func(c *cli.Context) error {
			err := do(port, protocol)
			if err != nil {
				return err
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

	//fmt.Printf(os.Args[1])

}
func do(port string, protocol string) error {
	file, fileerrerr := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if fileerrerr != nil {
		return fileerrerr
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	// 文件初始化
start: // 在这里循环
	host := ""
	_, err := fmt.Scanln(&host)
	if err != nil {
		return err
	}
	go func() {
		err := dofunc(port, protocol, file, host)
		if err != nil {

		}
	}()
	goto start
}
func dofunc(port string, protocol string, file *os.File, host string) error {
	if protocol == "http" {

		proxyStr := "http://" + host + ":" + port
		client := goreq.NewClient()
		req := goreq.Get("http://icanhazip.com/").SetClient(client).SetProxy(proxyStr).SetTimeout(5 * time.Second)
		//fmt.Printf(proxyStr)
		if req.Err == nil {
			ret := goreq.Do(req)
			if strings.Contains(ret.Text, host) {
				println(proxyStr) // 输出
				_, writeerr := file.WriteString(proxyStr + "\n")
				if writeerr != nil {
					return writeerr
				}

			}
		}

		//return nil
	} else if protocol == "ssh" {
		dial, err := net.Dial("tcp", host+":"+port)
		if err != nil {
			return err
		}
		_, err = dial.Write([]byte("")) // 发送空消息
		if err != nil {
			return err
		}
		buf := [512]byte{}
		n, err := dial.Read(buf[:])
		//println(string(buf[:n]))
		if strings.Contains(string(buf[:n]), "SSH") {
			println(host)
			_, writeerr := file.WriteString("[SSH]" + host + ":" + port + "\n")
			if writeerr != nil {
				return writeerr
			}
		}
	} else if protocol == "mysql" {
		dial, err := net.Dial("tcp", host+":"+port)
		if err != nil {
			return err
		}
		_, err = dial.Write([]byte("")) // 发送空消息
		if err != nil {
			return err
		}
		buf := [512]byte{}
		n, err := dial.Read(buf[:])
		//println(string(buf[:n]))
		if strings.Contains(string(buf[:n]), "mysql") {
			println(host)
			_, writeerr := file.WriteString("[mysql]" + host + ":" + port + "\n")
			if writeerr != nil {
				return writeerr
			}
		}
	} else {
		fmt.Printf("无此协议")
		return nil
	}
	return nil
}
