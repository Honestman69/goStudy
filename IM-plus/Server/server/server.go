package server

import (
	sqldb "IM-Server/sqlDB"
	"IM-Server/user"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	MessageChan chan string
	OnlineMap   map[string]*user.User
	MapLock     sync.RWMutex
}

func NewServer(Ip string, Port int) *Server {
	server := &Server{
		Ip:          Ip,
		Port:        Port,
		MessageChan: make(chan string),
		OnlineMap:   make(map[string]*user.User),
	}

	return server
}

// 私聊转发
func (server *Server) PrivateChat(msg string, user *user.User) {
	//对象名
	toName := strings.Split(msg, "|")[1]

	//发送信息
	toMsg := strings.Split(msg, "|")[2]
	for onlineName, toUser := range server.OnlineMap {
		if toName == onlineName {
			sendMsg := "[" + user.Name + "]" + toMsg
			toUser.Conn.Write([]byte(sendMsg))
			return
		}
	}
	sendMsg := "用户" + toName + "未上线..."
	user.Conn.Write([]byte(sendMsg))
}

// 在线用户查找
func (server *Server) OnlineFind(userName string) string {
	for onUserName := range server.OnlineMap {
		if onUserName == userName {
			return userName
		}
	}
	return ""
}

// 客户端改名处理
func (server *Server) UpdateName(sendMsg string, user *user.User) {
	newName := strings.Split(sendMsg, ":")[1]
	if server.OnlineFind(newName) != "" {
		_, err := user.Conn.Write([]byte("Fail:name"))
		if err != nil {
			fmt.Println("update name fail,err:", err)
			return
		}
	} else {
		dbMsg := sqldb.Update(newName, user)
		if dbMsg == "Success" {
			server.MapLock.Lock()
			delete(server.OnlineMap, user.Name)
			user.Name = newName
			server.OnlineMap[user.Name] = user
			server.MapLock.Unlock()
			_, err := user.Conn.Write([]byte("Success:name"))
			if err != nil {
				fmt.Println("conn write err:", err)
				return
			}
		} else if dbMsg == "have" {
			_, err := user.Conn.Write([]byte("Fail:name"))
			if err != nil {
				fmt.Println("update name fail,err:", err)
				return
			}
		} else {
			_, err := user.Conn.Write([]byte("服务器故障更新失败..."))
			if err != nil {
				fmt.Println("update name fail,err:", err)
				return
			}
		}

	}

}

// 客户端登录验证
func (server *Server) UserLogin(sendMsg string, user *user.User) {
	user.Name = strings.Split(sendMsg, ":")[1]
	user.Password = strings.Split(sendMsg, ":")[2]
	onlineUser := server.OnlineFind(user.Name)
	if onlineUser != "" {
		_, err := user.Conn.Write([]byte("Fail:login:have"))
		if err != nil {
			fmt.Println("user conn write fail,err:", err)
		}
		return
	}

	flag := sqldb.Select(user)
	if flag {
		_, err := user.Conn.Write([]byte("Success:login"))
		if err != nil {
			fmt.Println("user conn write fail,err:", err)
		}
		user.Password = ""
		server.OnlineUser(user)
		return
	}
	_, err := user.Conn.Write([]byte("Fail:login:err"))
	if err != nil {
		fmt.Println("user conn writefail,err:", err)
	}
}

// 客户端注册验证
func (server *Server) UserSign(sendMsg string, user *user.User) {
	user.Name = strings.Split(sendMsg, ":")[1]
	falg := sqldb.Select(user)
	if falg {
		_, err := user.Conn.Write([]byte("Fail:sign"))
		if err != nil {
			fmt.Println("user conn fail,err:", err)
		}
		return
	}
	user.Password = strings.Split(sendMsg, ":")[2]
	if true {
		err := sqldb.Insert(user)
		if err != nil {
			_, err := user.Conn.Write([]byte("用户创建失败..."))
			if err != nil {
				fmt.Println("user conn fail,err:", err)
				return
			}
			return
		}
	}

	_, err := user.Conn.Write([]byte("Success:sign"))
	if err != nil {
		fmt.Println("user conn fail,err:", err)
		return
	}
}

// 处理客户端消息
func (server *Server) DealMessage(sendMsg string, user *user.User) {
	if len(sendMsg) > 3 && sendMsg[:3] == "to|" { //私聊
		server.PrivateChat(sendMsg, user)
	} else if len(sendMsg) > 6 && sendMsg[:6] == "update" {
		server.UpdateName(sendMsg, user)
	} else if len(sendMsg) > 5 && sendMsg[:5] == "login" {
		server.UserLogin(sendMsg, user)
	} else if len(sendMsg) > 4 && sendMsg[:4] == "sign" {
		server.UserSign(sendMsg, user)
	} else if sendMsg == "Done" {
		server.OffUser(user)
	} else {
		msg := "[" + user.Name + "]" + sendMsg
		server.BroadcastMsg(msg)
	}

}

// 广播消息
func (server *Server) BroadcastMsg(sendMsg string) {
	for _, user := range server.OnlineMap {
		_, err := user.Conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("user write fail,err:", err)
			return
		}
	}

}

// 监听客户端消息
func (server *Server) ListenMessage(user *user.User) {
	for {
		var buf [4096]byte
		n, err := user.Conn.Read(buf[:])
		if err != nil {
			fmt.Println("conn fail,err:", err)
			return
		}
		sendMsg := string(buf[:n-2])
		server.DealMessage(sendMsg, user)
	}

}

// 处理客户端连接业务
func (server *Server) Handler(conn net.Conn) {
	user := user.NewUser(conn)
	server.ListenMessage(user)
}

// 用户上线功能
func (server *Server) OnlineUser(user *user.User) {
	server.MapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.MapLock.Unlock()
	sendMsg := "[" + user.Name + "]" + "已上线..."
	server.BroadcastMsg(sendMsg)
}

// 用户下线
func (server *Server) OffUser(user *user.User) {
	server.MapLock.Lock()
	delete(server.OnlineMap, user.Name)
	server.MapLock.Unlock()
	sendMsg := "[" + user.Name + "]" + "已下线..."
	server.BroadcastMsg(sendMsg)
}

// 服务器启动
func Start() {
	server := NewServer("127.0.0.1", 8888)
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("listen fail,err:", err)
		return
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen fail,err:", err)
			return
		}
		go server.Handler(conn)
	}

}
