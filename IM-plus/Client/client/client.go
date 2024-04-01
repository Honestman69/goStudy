package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var flag bool

type Client struct {
	Name     string
	Password string

	ServerIP   string
	ServerPort int

	conn net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIp,
		ServerPort: serverPort,
	}
	return client
}

// 服务器回传消息处理
func (client *Client) DealBackMessage(conn net.Conn) {
	for {
		var buf [4096]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("conn read fail,err:", err)
			return
		}
		msg := string(buf[:n])
		if len(msg) > 4 && msg[:4] == "Fail" {
			sign := strings.Split(msg, ":")[1]
			if sign == "name" {
				fmt.Println("用户名被占用...")
			} else if sign == "login" {
				sign = strings.Split(msg, ":")[2]
				if sign == "have" {
					fmt.Println("该账户已在别的终端登录...")
				} else {
					fmt.Println("用户名或密码错误,请重新登录...")
				}

			} else if sign == "sign" {
				fmt.Println("用户名已被注册，请重新注册...")
			}
			flag = false
		} else if len(msg) > 7 && msg[:7] == "Success" {
			sign := strings.Split(msg, ":")[1]
			if sign == "name" {
				fmt.Println("用户名修改成功...")
			} else if sign == "login" {
				fmt.Println("登录成功...")
			} else if sign == "sign" {
				fmt.Println("注册成功...")
			}
			flag = true
		} else {
			fmt.Println(msg)
		}
	}
}

// 改名
func (client *Client) UpdateName() {
	var sendMsg string
	fmt.Printf("请输入新的用户名：")
	fmt.Scanln(&sendMsg)
	if sendMsg == "exit" {
		return
	}
	_, err := client.conn.Write([]byte("update:" + sendMsg + "\r\n"))
	if err != nil {
		fmt.Println("用户名修改失败，conn write err:", err)
	}
}

// 群聊功能
func (client *Client) GroupChat() {
	fmt.Println("欢迎进入聊天室...")
	for {
		reader := bufio.NewReader(os.Stdin)
		sendMsg, _ := reader.ReadString('\n')
		if sendMsg == "exit\r\n" {
			return
		} else {
			client.SendMsg(sendMsg)
		}

	}
}

// 私聊功能
func (client *Client) PrivateChat() {
	var toName string
	for {
		if toName == "" {
			fmt.Println("请输入私聊对象用户名：")
			fmt.Scanln(&toName)
		}
		if toName == "exit" {
			return
		}

		fmt.Printf("请输入聊天内容：")
		reader := bufio.NewReader(os.Stdin)
		toMsg, _ := reader.ReadString('\n')
		if toMsg == "\r\n" {
			fmt.Println("发送内容不能为空...")
			reader := bufio.NewReader(os.Stdin)
			toMsg, _ = reader.ReadString('\n')
		}
		if len(toMsg) > 4 && toMsg[:4] == "exit" {
			return
		}
		sendMsg := "to|" + toName + "|" + toMsg
		client.SendMsg(sendMsg)
		time.Sleep(1 * time.Second)
	}
}

// 客户端消息发送
func (client *Client) SendMsg(sendMsg string) {
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("信息发送失败，conn write err:", err)
	}
}

// 登录
func (client *Client) Login() {
	flag = false
	for !flag {
		var name, password string
		fmt.Printf("请输入用户名：")
		fmt.Scanln(&name)
		if name == "exit" {
			return
		}
		fmt.Printf("请输入密码：")
		fmt.Scanln(&password)
		if password == "exit" {
			return
		}
		sendMsg := "login:" + name + ":" + password + "\r\n"
		client.SendMsg(sendMsg)
		time.Sleep(1 * time.Second)
	}

}

// 注册
func (client *Client) Sign() {
	flag = false
	for !flag {
		var name, password string
		fmt.Printf("请输入用户名：")
		fmt.Scanln(&name)
		if name == "exit" {
			return
		}
		fmt.Printf("请输入密码：")
		fmt.Scanln(&password)
		if password == "exit" {
			return
		}
		sendMsg := "sign:" + name + ":" + password + "\r\n"
		client.SendMsg(sendMsg)
		time.Sleep(1 * time.Second)
	}

}

// 登出功能
func (client *Client) Done() {
	sendMsg := "Done" + "\r\n"
	client.SendMsg(sendMsg)
}

// 登录注册菜单
func (client *Client) StartMenu() {
	fmt.Println("   1.注册")
	fmt.Println("   2.登录")
	fmt.Println("   3.退出")
}

// 功能菜单
func (client *Client) Menu() {
	fmt.Println(">>>>菜单<<<<")
	fmt.Println("   1.群聊")
	fmt.Println("   2.私聊")
	fmt.Println("   3.改名")
	fmt.Println("   4.退出")
	fmt.Println(">>>>菜单<<<<")
	fmt.Println("")
}

// 选项表
func (client *Client) Run() {
	var option int
	//功能选项
	for {
		client.Menu()
		fmt.Scanln(&option)
		if option < 1 || option > 4 {
			fmt.Println("请输入正确的数字...")
			fmt.Scanln(&option)
		}
		switch option {
		case 1:
			client.GroupChat()
		case 2:
			client.PrivateChat()
		case 3:
			client.UpdateName()
		case 4:
			client.Done()
			return
		}
	}
}

func (client *Client) StartRun() {
	//登录注册选项
	var option int
	for {
		client.StartMenu()
		fmt.Scanln(&option)
		if option < 1 || option > 3 {
			fmt.Println("请输入正确的数字...")
			fmt.Scanln(&option)
		}
		switch option {
		case 1:
			client.Sign()
		case 2:
			client.Login()
			client.Run()
		case 3:
			return
		}
	}
}

// 客户端开启
func Start() {
	client := NewClient("127.0.0.1", 8888)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIP, client.ServerPort))
	client.conn = conn
	if err != nil {
		fmt.Println("dial fail,err:", err)
		return
	}

	go client.DealBackMessage(conn)

	client.StartRun()
}
