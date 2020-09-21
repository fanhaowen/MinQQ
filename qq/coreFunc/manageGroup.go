package coreFunc

import (
	"golang.org/code.learn.fund/qq/mysql"
	"net"
	"strings"
)

func Manage(conn net.Conn, cmd, userID string){
	groupManageID := cmd[7:]
	if mysql.JudgeGroupExit(mysql.DB, groupManageID) {
		if mysql.JudgeUserAdmin(mysql.DB, userID, groupManageID) {
			sysChoseMsg := `grouplist:查询群员名单
			delete userID:删除该用户
			apply groupID 查询申请表单
			accept userID 同意user进群`
			_, _ = conn.Write([]byte(sysChoseMsg))
			var msgWriter string
			Listen(conn, &msgWriter)
			if msgWriter == "grouplist" {
				members := mysql.GetAllMember(mysql.DB, groupManageID)
				sysChoseMsg = "Members:"
				for _, x := range members {
					if x != "" {
						sysChoseMsg += x + " "
					}
				}
				_, _ = conn.Write([]byte(sysChoseMsg + "\n"))
			}
			if strings.HasPrefix(msgWriter, "delete") {
				groupDeleteID := msgWriter[7:]
				mysql.DeleteUser(mysql.DB, groupDeleteID)
				sysChoseMsg = "删除成功\n"
				_, _ = conn.Write([]byte(sysChoseMsg))
			}
			if msgWriter == "apply" {
				applyList := mysql.GetApplyUserS(mysql.DB, userID)
				for _, v := range applyList {
					_, _ = conn.Write([]byte(v))
				}
			}
			if strings.HasPrefix(msgWriter, "accept") {
				
			}
		} else {
			sysChoseMsg := "你不是群主，无权操作"
			_, _ = conn.Write([]byte(sysChoseMsg))
		}
		
	}
}