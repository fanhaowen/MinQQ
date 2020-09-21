package coreFunc

import (
	"fmt"
	"net"
)

func DeleteUser(userID string, ConnMap *map[string]net.Conn) {
	fmt.Println("准备清除用户在线信息:", userID)
	delete(*ConnMap, userID)
	fmt.Println("清除完毕")
}