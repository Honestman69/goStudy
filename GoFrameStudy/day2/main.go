package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	// test1()
	test2()
}

func test1() {
	s := g.Server()
	s.BindHandler("/:name", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.BindHandler("/:name/update", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.BindHandler("/:name/:action", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.BindHandler("/:name/*any", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.BindHandler("/user/list/{field}.html", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.SetPort(8888)
	s.Run()
}

func test2() {
	s := g.Server()
	// 一个简单的分页路由示例
	s.BindHandler("/user/list/{page}.html", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("page"))
	})
	// {xxx} 规则与 :xxx 规则混合使用
	s.BindHandler("/{object}/:attr/{act}.php", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("object"))
		r.Response.Writeln(r.Get("attr"))
		r.Response.Writeln(r.Get("act"))
	})
	// 多种模糊匹配规则混合使用
	s.BindHandler("/{class}-{course}/:name/*act", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("class"))
		r.Response.Writeln(r.Get("course"))
		r.Response.Writeln(r.Get("name"))
		r.Response.Writeln(r.Get("act"))
	})
	s.SetPort(8199)
	s.Run()

}
