package main

import (
	"context"
	"net/http"

	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

type HelloReq struct {
	g.Meta `path:"/hello" method:"get"`
	Name   string `v:"required" dc:"Your name"`
}
type HelloRes struct {
	Reply string `dc:"Reply content"`
}

type Hello struct{}

func (Hello) Say(ctx context.Context, req *HelloReq) (res *HelloRes, err error) {
	g.Log().Debugf(ctx, `receive say: %+v`, req)
	res = &HelloRes{
		Reply: fmt.Sprintf(`Hi %s`, req.Name),
	}
	return
}

func UserGet(r *ghttp.Request) {
	r.Response.Write("get")
}

func UserDelete(r *ghttp.Request) {
	r.Response.Write("delete")
}

func MiddlewareAuth(r *ghttp.Request) {
	token := r.Get("token")
	if gconv.String(token) == "123456" {
		r.Middleware.Next()
	} else {
		r.Response.WriteStatus(http.StatusForbidden)
	}
}

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func MiddlewareLog(r *ghttp.Request) {
	r.Middleware.Next()
	fmt.Println(r.Response.Status, r.URL.Path)
}

func main() {
	// LevelEnroll()
	// MapEnroll()
	SwaggerTest()
}

func LevelEnroll() {
	s := g.Server()
	s.Use(MiddlewareLog)
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareAuth, MiddlewareCORS)
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
		group.Group("/order", func(group *ghttp.RouterGroup) {
			group.GET("/list", func(r *ghttp.Request) {
				r.Response.Write("list")
			})
			group.PUT("/update", func(r *ghttp.Request) {
				r.Response.Write("update")
			})
		})
		group.Group("/user", func(group *ghttp.RouterGroup) {
			group.GET("/info", func(r *ghttp.Request) {
				r.Response.Write("info")
			})
			group.POST("/edit", func(r *ghttp.Request) {
				r.Response.Write("edit")
			})
			group.DELETE("/drop", func(r *ghttp.Request) {
				r.Response.Write("drop")
			})
		})
		group.Group("/hook", func(group *ghttp.RouterGroup) {
			group.Hook("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
				r.Response.Write("hook any")
			})
			group.Hook("/:name", ghttp.HookBeforeServe, func(r *ghttp.Request) {
				r.Response.Write("hook name")
			})
		})
	})
	s.Run()
}

func MapEnroll() {
	s := g.Server()
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Map(g.Map{
			"Get:/user":    UserGet,
			"Delete:/user": UserDelete,
		})
	})
	s.Run()
}

func SwaggerTest() {
	s := g.Server()
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			new(Hello),
		)
	})
	s.Run()
}
