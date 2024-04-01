package user

import (
	"net"
)

type User struct {
	Name     string
	Password string
	Addr     string

	MsgChan chan string
	Conn    net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,

		MsgChan: make(chan string),
		Conn:    conn,
	}
	return user
}
