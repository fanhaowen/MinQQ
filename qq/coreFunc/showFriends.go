package coreFunc

import (
	"golang.org/code.learn.fund/qq/mysql"
	"net"
)

func ShowFriends(conn net.Conn,userID string, ConnMap map[string]net.Conn){
	sysChoseMsg := ""
	friendsList := mysql.GetFriends(mysql.DB, userID)
	sysChoseMsg = `
		Your friends list:	`
	for _, value := range friendsList {
		if value != "" {
			sysChoseMsg += value
			isOnline := ConnMap[value]
			if isOnline != nil {
				sysChoseMsg += "(在线)"
			}
			sysChoseMsg += " "
		}
	}
	sysChoseMsg += "\n"
	groupList, _ := mysql.GetUserGroups(mysql.DB, userID)
	sysChoseMsg += `
		Your Group list:	`
	for _, value := range groupList {
		sysChoseMsg += value
		isAdmin := mysql.JudgeUserAdmin(mysql.DB, userID, value)
		if isAdmin {
			sysChoseMsg += "(群主)"
		}
		sysChoseMsg += " "
	}
	sysChoseMsg += "\n"
	_, _ = conn.Write([]byte(sysChoseMsg))
}
func ShowWelcome(conn net.Conn){
	sysChoseMsg :=   `
		*	与好友进行对话 (input CHAT'聊天')	          *
		*	添加好友 (input ADD'添加')		           *
		*	退出 (input EXIT or .end'退出')		           *
		*	刷新好友列表 (input FLASH'刷新')	           *
		*	创建群聊 (input create or "创建")	           *
		*	接收离线信息(input rev)					*
		*	申请加入群聊 (input "apply groupID/userID")*
		*	处理群聊相关内容(群主专属 input"manage groupID")   *
		*	查询聊天记录(只支持用户之间的记录 sr userID)*
			:`
	_, _ = conn.Write([]byte(sysChoseMsg))
}