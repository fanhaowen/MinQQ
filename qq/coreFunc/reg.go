package coreFunc

import (
	"fmt"
	"golang.org/code.learn.fund/qq/mysql"
	"net"
	"strconv"
)
func Register(conn net.Conn) bool {
	sysChoseMsg := ""
	sysChoseMsg = "Msg From Server: 输入你的名字:"
	_, _ = conn.Write([]byte(sysChoseMsg))
	var regName string
	Listen(conn, &regName)
	sysChoseMsg = "Msg From Server: 输入ID账号:"
	_, _ = conn.Write([]byte(sysChoseMsg))
	var regID string
	Listen(conn, &regID)
	ok := mysql.JudgeUserExit(mysql.DB, regID)
	if ok {
		return false
	}
	sysChoseMsg = "Msg From Server: 你的城市:"
	_, _ = conn.Write([]byte(sysChoseMsg))
	var regLoc string
	Listen(conn, &regLoc)
	sysChoseMsg = "sex:"
	_, _ = conn.Write([]byte(sysChoseMsg))
	var regSex string
	Listen(conn, &regSex)
	sysChoseMsg = "age:"
	_, _ = conn.Write([]byte(sysChoseMsg))
	var regAge string
	Listen(conn, &regAge)
	age, _ := strconv.Atoi(regAge)
	sysChoseMsg = "Msg From Server: 你的密码:"
	_, _ = conn.Write([]byte(sysChoseMsg))
	var regPW string
	Listen(conn, &regPW)
	regPW = MD5(regPW)
	register(regName, age, regLoc, regSex, regID, regPW)
	sysChoseMsg = "Msg From Server: 注册好了\n"
	_, _ = conn.Write([]byte(sysChoseMsg))
	return true
}
func register(name string, age int, location string, sex string, userID string, pw string) {
	//调用mysql.go里面写好的函数
	registerSQL := mysql.InsertData(mysql.DB)
	registerSQL(name, age, location, sex, userID, pw)
	fmt.Println("Register finished!")
	//注册完成后为用户添加一张表 来存储好友信息
	mysql.CreateFriendTable(mysql.DB, userID)
	fmt.Println("Create Friend Table Finished!")
}