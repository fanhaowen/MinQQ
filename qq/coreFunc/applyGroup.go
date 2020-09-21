package coreFunc

import (
	"golang.org/code.learn.fund/qq/mysql"
	"net"
)

func ApplyGroup(conn net.Conn, cmd string, userID string){
	msg := ""
	JoinID := cmd[6:]
	if JoinID[0] == 'g' {
		if mysql.JudgeGroupExit(mysql.DB, JoinID[1:]) {
			//mysql.InsertMember(mysql.DB, userID, JoinID)
			mysql.InsertApplyGroup(mysql.DB,userID, JoinID[1:])
			msg = "申请成功\n"
			_, _ = conn.Write([]byte(msg))
		} else {
			msg = "该群不存在"
			_, _ = conn.Write([]byte(msg))
		}
	} else {
		userFriends := mysql.GetFriends(mysql.DB, JoinID)
		var friendExist bool
		for _, value := range userFriends {
			if value == JoinID {
				friendExist = true
			}
		}
		if mysql.JudgeUserExit(mysql.DB, JoinID) {
			if !friendExist {
				mysql.AddFriend(mysql.DB, userID, JoinID)
				msg = "Msg From Server: 添加成功\n"
				_, _ = conn.Write([]byte(msg))
			} else {
				msg = "Msg From Server: ta已经是你的好友了 请勿重复添加\n"
				_, _ = conn.Write([]byte(msg))
			}
		} else {
			msg = "Msg From Server: 查无此人\n"
			_, _ = conn.Write([]byte(msg))
		}
	}

}
