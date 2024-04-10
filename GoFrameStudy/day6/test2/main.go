package main

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	TraceIdName = "trace-id"
)

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			ctx := context.WithValue(r.Context(), TraceIdName, "HBm876TFCde435Tgf")
			r.SetCtx(ctx)
			r.Middleware.Next()
		})
		group.ALL("/", func(r *ghttp.Request) {
			r.Response.Write(r.Context().Value(TraceIdName))
			// 也可以使用
			// r.Response.Write(r.GetCtxVar(TraceIdName))
		})
	})
	s.SetPort(8888)
	s.Run()
}
