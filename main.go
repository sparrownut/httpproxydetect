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
			file, fileerrerr := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if fileerrerr != nil {
				return fileerrerr
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {

				}
			}(file)
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

					_, writeerr := file.WriteString(proxyStr + "\n")
					if writeerr != nil {
						return writeerr
					}

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
