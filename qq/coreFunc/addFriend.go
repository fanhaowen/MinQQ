package coreFunc

import (
	"golang.org/code.learn.fund/qq/mysql"
	"net"
)

func AddFriend(conn net.Conn, userID string){
	sysChoseMsg := "Msg From Server:输入对方的 ID:"
	_, _ = conn.Write([]byte(sysChoseMsg))
	var addID string
	Listen(conn, &addID)
	userFriends := mysql.GetFriends(mysql.DB, addID)
	var friendExist bool
	for _, value := range userFriends {
		if value == addID {
			friendExist = true
		}
	}
	if mysql.JudgeUserExit(mysql.DB, addID) {
		if !friendExist {
			mysql.AddFriend(mysql.DB, userID, addID)
			sysChoseMsg = "Msg From Server: 添加成功\n"
			_, _ = conn.Write([]byte(sysChoseMsg))
		} else {
			sysChoseMsg = "Msg From Server: ta已经是你的好友了 请勿重复添加\n"
			_, _ = conn.Write([]byte(sysChoseMsg))
		}
	} else {
		sysChoseMsg = "Msg From Server: 查无此人\n"
		_, _ = conn.Write([]byte(sysChoseMsg))
	}
}