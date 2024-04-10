package main

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/os/gtime"
)

func SessionFile() {
	s := g.Server()
	s.SetSessionMaxAge(time.Minute)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/set", func(r *ghttp.Request) {
			r.Session.Set("time", gtime.Timestamp())
			r.Response.Write("ok")
		})
		group.ALL("/get", func(r *ghttp.Request) {
			r.Response.Write(r.Session.Data())
		})
		group.ALL("/delect", func(r *ghttp.Request) {
			_ = r.Session.RemoveAll()
			r.Response.Write("ok")
		})
	})
	s.SetPort(8888)
	s.Run()
}

func SessionMemory() {
	s := g.Server()
	s.SetSessionMaxAge(time.Minute)
	s.SetSessionStorage(gsession.NewStorageMemory())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/set", func(r *ghttp.Request) {
			r.Session.MustSet("time", gtime.Timestamp())
			r.Response.Write("ok")
		})
		group.ALL("/get", func(r *ghttp.Request) {
			r.Response.Write(r.Session.Data())
		})
		group.ALL("/delect", func(r *ghttp.Request) {
			_ = r.Session.RemoveAll()
			r.Response.Write("ok")
		})
	})
	s.SetPort(8888)
	s.Run()

}

func SessionRedisKV() {
	s := g.Server()
	s.SetSessionMaxAge(time.Minute)
	s.SetSessionStorage(gsession.NewStorageRedis(g.Redis()))
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/set", func(r *ghttp.Request) {
			r.Session.Set("time", gtime.Timestamp())
			r.Response.Write("ok")
		})
		group.ALL("/get", func(r *ghttp.Request) {
			r.Response.Write(r.Session.Data())
		})
		group.ALL("/delect", func(r *ghttp.Request) {
			_ = r.Session.RemoveAll()
			r.Response.Write("ok")
		})
	})
	s.SetPort(8888)
	s.Run()
}

func main() {
	//SessionMemory()
	SessionRedisKV()
}
