package coreFunc

import (
	"database/sql"
	"golang.org/code.learn.fund/qq/mysql"
	"net"
)

func GetOffMSG(conn net.Conn, DB *sql.DB, userID string){
	msgS := mysql.GetUserMsgOffLine(DB, userID)
	msgWriter := ""
	for _, msg := range msgS {
		msgWriter += msg+"\n"
	}
	_, _ = conn.Write([]byte(msgWriter))
}
