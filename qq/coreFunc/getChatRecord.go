package coreFunc

import (
	"database/sql"
	"golang.org/code.learn.fund/qq/mysql"
	"net"
)

func GetRecord(conn net.Conn, DB *sql.DB,userID, friendID string) {
	msgS := mysql.GetRecord(DB, userID, friendID)
	msgWriter := ""
	for _, msg := range msgS {
		msgWriter += msg+"\n"
	}
	_, _ = conn.Write([]byte(msgWriter))
}
