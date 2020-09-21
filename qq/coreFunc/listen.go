package coreFunc

import (
	"net"
)
var ErrFlag bool
func Listen(conn net.Conn, msgReader *string) bool {
	// buf是缓冲区
	buf := make([]byte, 1024)
	//将信息读取，存放在buf里面，reader是信息的长度
	reader, err := conn.Read(buf)
	if err != nil {
		ErrFlag = true
	}
	//将内容取出放入*msgReader
	*msgReader = string(buf[:reader])
	
	//检查该信息是否为结束信号，特征为.end 或 Quit!
	if *msgReader == "EXIT" || *msgReader == "quit"{
		return false
	}
	return true
}