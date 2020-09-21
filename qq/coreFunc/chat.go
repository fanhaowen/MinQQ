package coreFunc

import (
	"fmt"
	"golang.org/code.learn.fund/qq/mysql"
	"net"
	"strings"
	"time"
)

func Chat(conn net.Conn, userID string, ConnMap map[string]net.Conn){
	var msgWriter string
	sysChoseMsg := "Msg From Server: 输入好友的 ID:"
	_, _ = conn.Write([]byte(sysChoseMsg))
	var friendID string
	Listen(conn, &friendID)
	if strings.HasPrefix(friendID, "g") {	//选择了群组
		if mysql.JudgeGroupExit(mysql.DB, friendID[1:]) {
			sysChoseMsg = "您已进入该群聊\nSend:"
			_, _ = conn.Write([]byte(sysChoseMsg))
			for {
				isQuit := Listen(conn, &msgWriter)
				if isQuit  {
					return
				}
				if len(msgWriter) > 0 {
					sysChoseMsg = ("\n" + mysql.GetName(mysql.DB, userID, true)) + time.Now().Format("2006/1/2 15:04:05") + "From Group" + mysql.GetName(mysql.DB, friendID, false) + " Send:\n" + msgWriter + "\n"
					for vID := range ConnMap {
						if userID != vID {
							_, _ = ConnMap[vID].Write([]byte(sysChoseMsg))	//给除了自己的人发
						}
					}
				}
			}
		}
	} else {
		//判断该ID是否为好友
		friendExist := false
		friendsList := mysql.GetFriends(mysql.DB, userID)
		for _, fid := range friendsList {
			if friendID == fid {
				friendExist = true
			}
		}
		if friendExist {
			_, ok := ConnMap[friendID]
			if ok {
				for {
					_, _ = conn.Write([]byte("Send:"))
					isOK := Listen(conn, &msgWriter)
					if !isOK {
						return
					}
					if len(msgWriter) > 0 {
						sysChoseMsg = ("\n" + mysql.GetName(mysql.DB, userID, true)) + time.Now().Format("2006/1/2 15:04:05") + " Send:\n" + msgWriter + "\n"
						_, _ = ConnMap[friendID].Write([]byte(sysChoseMsg))
						mysql.InsertChatRecord(mysql.DB, userID, friendID, msgWriter)
						
					}
				}
			} else {
				sysChoseMsg = "Msg From Server: 对方不在线，可以继续发送离线消息\n"
				_, _ = conn.Write([]byte(sysChoseMsg))
				for {
					_, _ = conn.Write([]byte("Send:"))
					isOK := Listen(conn, &msgWriter)
					if !isOK {
						fmt.Println("isQuit:",isOK,"\n msg:",msgWriter)
						return
					}
					if len(msgWriter) > 0 {
						mysql.InsertChatRecord(mysql.DB, userID, friendID, msgWriter)
						mysql.InsertOffLineMsg(mysql.DB, userID, friendID, msgWriter)
					}
				}
			}
			
		} else {
			sysChoseMsg = "Msg From Server: 对方不是您的好友 请在您的好友列表中选择或添加对方为好友\n"
			_, _ = conn.Write([]byte(sysChoseMsg))
		}
	}
}
