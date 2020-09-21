package coreFunc

import (
	"fmt"
	"golang.org/code.learn.fund/qq/mysql"
	"net"
)

func Login(conn net.Conn, ConnMap map[string]net.Conn) (bool, string)  {
	sysChoseMsg := "Msg From Server:请输入你的 ID:"
	msgReader := ""
	_, _ = conn.Write([]byte(sysChoseMsg))
	Listen(conn, &msgReader)
	userID := msgReader
	v := ConnMap[userID]
	if v != nil {
		msyWriter := "该用户已在线"
		_, _ = conn.Write([]byte(msyWriter))
		return false,""
	}
	//从数据库获得密码进行匹配
	sysChoseMsg = "Msg From Server:请输入你的密码:"
	_, _ = conn.Write([]byte(sysChoseMsg))
	Listen(conn, &msgReader)
	
	password := msgReader
	password = MD5(password)	//密码加密
	pwFromDB := mysql.GetPW(mysql.DB, userID)
	
	if pwFromDB == password {
		fmt.Println(userID, "Login Succeed!")
		return true, userID
	} else if pwFromDB != password {
		sysChoseMsg := "Msg From Server: 账号密码不匹配或错误 请重新输入\n"
		_, _ = conn.Write([]byte(sysChoseMsg))
	}
	return false, ""
}