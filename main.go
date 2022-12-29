package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/zhshch2002/goreq"
	"os"
	"strings"
	"time"
)

func main() {
	port := "8080"
	app := &cli.App{
		Name:      "proxydetect",
		Usage:     "judg proxy can use",
		UsageText: "lazy to write...",
		Version:   "0.4.4",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Destination: &port, Value: "8080", Usage: "port"},
		},
		HideHelpCommand: true,
		Action: func(c *cli.Context) error {
		start:
			host := ""
			_, err := fmt.Scanln(&host)
			if err != nil {
				return err
			}
			proxyStr := "http://" + host + ":" + port
			client := goreq.NewClient()
			req := goreq.Get("http://icanhazip.com/").SetClient(client).SetProxy(proxyStr).SetTimeout(5 * time.Second)
			//fmt.Printf(proxyStr)
			if req.Err == nil {
				ret := goreq.Do(req)
				if strings.Contains(ret.Text, host) {
					println(proxyStr)
				}
			}

			goto start
			//return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

	//fmt.Printf(os.Args[1])

}
