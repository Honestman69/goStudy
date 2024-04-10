package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	Main = &gcmd.Command{
		Name:        "main",
		Brief:       "start server",
		Description: "服务启动函数",
	}
	Http = &gcmd.Command{
		Name:        "http",
		Brief:       "start http server",
		Description: "启动http服务",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			fmt.Println("grcp server start")
			return
		},
	}
	Grpc = &gcmd.Command{
		Name:        "grpc",
		Brief:       "start grpc server",
		Description: "启动grpc服务",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			fmt.Println("grcp server start")
			return
		},
	}
)

func main() {
	err := Main.AddCommand(Http, Grpc)
	if err != nil {
		panic(err)
	}
	Main.Run(gctx.New())
}
