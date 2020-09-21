package coreFunc

import (
	"golang.org/code.learn.fund/qq/mysql"
	"net"
)

func CreatGroup(conn net.Conn, userID string){
	sysChoseMsg := "请输入你要创建的群号:"
	_, _ = conn.Write([]byte(sysChoseMsg))
	var groupIDCreate string
	Listen(conn, &groupIDCreate)
	if mysql.JudgeGroupExit(mysql.DB, groupIDCreate) {
		sysChoseMsg = "群号已存在啦"
		_, _ = conn.Write([]byte(sysChoseMsg))
	} else {
		sysChoseMsg = "请输入群的名字:"
		_, _ = conn.Write([]byte(sysChoseMsg))
		var groupCreatName string
		Listen(conn, &groupCreatName)
		mysql.CreateGroup(mysql.DB, userID, groupIDCreate, groupCreatName)
		
		sysChoseMsg = "创建群完成\n"
		_, _ = conn.Write([]byte(sysChoseMsg))
		
	}
	
}
