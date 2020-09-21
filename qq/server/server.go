package main

import (
	"fmt"
	"golang.org/code.learn.fund/qq/coreFunc"
	"golang.org/code.learn.fund/qq/mysql"
	"net"
	"strings"
	"time"
)

var online []string

//ConnMap 存储在线用户列表
var ConnMap map[string]net.Conn

func init() {
	mysql.CreateTable(mysql.DB)
	mysql.CreateGroupTable(mysql.DB)
	mysql.CreateGroupUser(mysql.DB)
	mysql.CreateGroupApply(mysql.DB)
	mysql.CreateOffLineMsg(mysql.DB)
	mysql.CreateChatRecord(mysql.DB)
	online = make([]string, 20)
	ConnMap = make(map[string]net.Conn)
}

func main() {
	//运行结束时退出数据库，在mysql.go中已经开启
	defer mysql.DB.Close()
	//开启服务器
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("err:", err)
	}
	for {
		//检测到有用户连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("connect from:", conn.RemoteAddr().String())
		//开启线程处理该用户
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	// defer deleteUser(userID1)
	defer conn.Close()
	defer fmt.Println("用户退出了")
beginIndex:
	sysChoseMsg := `Msg From Server: 选择一项以继续:
	登录 : 1			注册 : 2
			:`
	_, _ = conn.Write([]byte(sysChoseMsg))
	var msgReader string
	for {
		coreFunc.Listen(conn, &msgReader)
		//开始登录
		if msgReader == "1" {
			for {
				loginOK, userID := coreFunc.Login(conn, ConnMap)
				if !loginOK {
					goto beginIndex
				}
				if loginOK {	//登录成功！
					for {
						time.Sleep(time.Second)
						ConnMap[userID] = conn		//该用户现在为已登录状态
						sysChoseMsg = fmt.Sprintf("Msg From Server: Welecome Here, %s\n", mysql.GetName(mysql.DB, userID, true))
						_, _ = conn.Write([]byte(sysChoseMsg))
						online = append(online, userID)
						coreFunc.ShowFriends(conn, userID, ConnMap)	//显示好友列表
						coreFunc.ShowWelcome(conn)					//欢迎界面
						for {
							var indexChose string
							coreFunc.Listen(conn, &indexChose)
							if indexChose == "ADD" || indexChose == "add" || indexChose == "添加" {
								coreFunc.AddFriend(conn, userID)
							} else if indexChose == "CHAT" || indexChose == "chat" || indexChose == "聊天" {
								coreFunc.Chat(conn, userID, ConnMap)
							} else if indexChose == "FLASH" || indexChose == "flash" || indexChose == "刷新" {
								coreFunc.ShowFriends(conn, userID, ConnMap)
							} else if indexChose == "create" || indexChose == "创建" {
								coreFunc.CreatGroup(conn, userID)
							} else if strings.HasPrefix(indexChose, "manage") {
								coreFunc.Manage(conn, indexChose, userID)
							} else if strings.HasPrefix(indexChose, "applygroup") {
								coreFunc.ApplyGroup(conn, indexChose, userID)
							} else if indexChose == "rev" {
								coreFunc.GetOffMSG(conn, mysql.DB, userID)
							} else if indexChose == "EXIT" || indexChose == "退出" || indexChose == "exit" {
								fmt.Println(userID,"退出连接")
								coreFunc.DeleteUser(userID, &ConnMap)
								return
							} else if strings.HasPrefix(indexChose, "sr"){
								coreFunc.GetRecord(conn, mysql.DB, userID, indexChose[3:])
							} else {
								sysChoseMsg = "你输入了无效码\n"
								_, _ = conn.Write([]byte(sysChoseMsg))
							}
							if !coreFunc.ErrFlag {
								sysChoseMsg = "请重新选择操作ADD,CHAT,EXIT,FLASH,APPLY,MANAGE\n"
								_, _ = conn.Write([]byte(sysChoseMsg))
							} else {
								coreFunc.DeleteUser(userID, &ConnMap)
								return
							}
						}
					}
				}
			}
		} else if msgReader == "2" {
			if !coreFunc.Register(conn) {
				sysChoseMsg = "注册失败了，可能是由于你输入了已存在的账号\n"
				_, _ = conn.Write([]byte(sysChoseMsg))
				goto beginIndex
			}
			goto beginIndex
		} else {
			sysChoseMsg = "你输入了无效码,请重新输入\n"
			_, _ = conn.Write([]byte(sysChoseMsg))
			goto beginIndex
		}
	}
}
