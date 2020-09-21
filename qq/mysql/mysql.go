package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//User struct
type User struct {
	Name   string
	UserID int
	PW     string
}
//DB database
var DB *sql.DB
var dbErr error

func init() {
	DB, dbErr = sql.Open("mysql",
		"root:123456@tcp(127.0.0.1:3306)/test")
	if dbErr != nil {
		fmt.Println(dbErr)
	}
}
//InsertData  to insert data
func InsertData(DB *sql.DB) func(string, int, string, string, string, string) {
	return func(name string, age int, location string, sex string, userID string, pw string) {
		queryString := fmt.Sprintf("INSERT INTO go_users (NAME, AGE, LOCATION, SEX, USERID, PW) VALUES (\"%s\", %d, \"%s\", \"%s\", %s, \"%s\")",
			name, age, location, sex, userID, pw)
		rows, err := DB.Query(queryString)
		if err != nil {
			fmt.Println("insert data error:", err)
			return
		}
		defer rows.Close()
	}
}
//AddFriend can add other user as friend
func AddFriend(DB *sql.DB, userID string, friendID string) {
	userTable := fmt.Sprintf("friends_%s", userID)
	queryString := fmt.Sprintf("INSERT INTO %s (friendID) VALUES (%s);", userTable, friendID)
	rows, err := DB.Query(queryString)

	if err != nil {
		fmt.Println("insert data error:", err)
	}
	fmt.Println("添加好友成功")
	//同时好友的好友列表也添加
	friendTable := fmt.Sprintf("friends_%s", friendID)
	queryString = fmt.Sprintf("INSERT INTO %s (friendID) VALUES (%s);", friendTable, userID)
	rows, _ = DB.Query(queryString)
	defer rows.Close()

	if err != nil {
		fmt.Println("insert data error:", err)
	}
	fmt.Println("对方添加好友成功")
}
//InsertMember 加入成员到群聊
func InsertMember(DB *sql.DB, memberID string, groupID string)  {
	defer func() {
		errR := recover()
		if errR != nil {
			fmt.Println("加入成员到群聊表中时发生错误")
		}
	}()
	queryString := fmt.Sprintf("insert into go_group_user (group_id, userID) values (%s,%s);", groupID, memberID)
	rows, err := DB.Query(queryString)

	
	if err == nil {
		fmt.Println("成功创建群聊-用户表")
	} else {
		fmt.Println("创建群聊-用户表失败")
		return
	}
	defer rows.Close()
}
//InsertOffLineMsg 插入一条离线信息
func InsertOffLineMsg(DB *sql.DB, userID, friendID, msg string){
	queryString := fmt.Sprintf("insert into go_offline_msg (srcID, dstID, msg, sendTime) VALUES (%s, %s, '%s', now());",userID,friendID,msg)
	fmt.Println(queryString)
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Println("插入离线信息时出错：",err)
		return
	}
	rows.Close()
}
//InsertChatRecord 插入一条消息到聊天记录
func InsertChatRecord(DB *sql.DB, userID, friendID, msg string){
	queryString := fmt.Sprintf("insert into go_chat_record (srcID, dstID, msg, sendTime) VALUES (%s, %s, '%s', now());",userID,friendID,msg)
	fmt.Println(queryString)
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Println("插入聊天记录信息时出错：",err)
		return
	}
	rows.Close()
}
//插入一条申请到群聊申请表
func InsertApplyGroup(DB *sql.DB, userID, groupID string){
	adminName := GetAdmin(DB, groupID)
	queryString := fmt.Sprintf("insert into go_group_apply () VALUES (%s %s %s)", groupID, adminName, userID)
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Println("插入群聊申请表时出错：",err)
		return
	}
	rows.Close()
	return
}
//获取群组的admin
func GetAdmin(DB *sql.DB, groupID string) (res string) {
	queryString := fmt.Sprintf("select admin_id from go_groups where group_id=%s", groupID)
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Println("获取群聊admin时发生错误")
		return ""
	}
	defer rows.Close()
	rows.Next()
	_ = rows.Scan(&res)
	return
}
//获取该用户收到的离线信息
func GetUserMsgOffLine(DB *sql.DB, userID string) []string {
	var msgS []string
	queryString := fmt.Sprintf("select srcID, msg, sendTime from go_offline_msg where dstID=%s ORDER BY sendTime;",userID)
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Println("获取离线消息时发生了错误:",err)
		return msgS
	}
	var msg string
	var time string
	var srcID string
	for rows.Next(){
		err = rows.Scan(&srcID, &msg, &time)
		if err != nil {
			fmt.Printf("扫描离线信息时错误：%v\n", err)
			return msgS
		}
		name := GetName(DB, srcID,true)
		msgS = append(msgS, name+" "+time+" Send:\n"+msg+"\n")
	}
	
	return msgS
}
//返回该用户与目标用的聊天记录
func GetRecord(DB *sql.DB, userID, friendID string) []string {
	var msgS []string
	queryString := fmt.Sprintf("select srcID, msg, sendTime from go_chat_record where (srcID=%s and dstID=%s) OR (srcID=%s and dstID=%s)ORDER BY sendTime;", userID, friendID, friendID, userID)
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Println("获取离线消息时发生了错误:",err)
		return msgS
	}
	defer rows.Close()
	var msg string
	var time string
	var srcID string
	for rows.Next(){
		err = rows.Scan(&srcID, &msg, &time)
		if err != nil {
			fmt.Printf("扫描离线信息时错误：%v\n", err)
			return msgS
		}
		name := GetName(DB, srcID,true)
		msgS = append(msgS, name+" "+time+" Send:\n"+msg+"\n")
	}
	return msgS
}
//GetPW 返回用户的密码MD5
func GetPW(DB *sql.DB, userID string) string {
	user := new(User)
	queryString := fmt.Sprintf("select pw from go_users where userID=%s;", userID)
	rows, _ := DB.Query(queryString)
	rows.Next()
	if err := rows.Scan(&user.PW); err != nil {
		fmt.Printf("scan failed, err:%v", err)
		return "failed scan"
	}
	defer rows.Close()
	return user.PW
}
//GetFriends 返回用户的所有好友
func GetFriends(DB *sql.DB, userID string) []string {
	friends := make([]string, 10, 20)
	queryString := fmt.Sprintf("select * from friends_%s;", userID)
	rows, er := DB.Query(queryString)
	if er != nil {
		fmt.Println("Error occurred in mysql:", er)
		return friends
	}
	defer rows.Close()
	var err error
	var friend string
	
	for rows.Next() {
		err = rows.Scan(&friend)
		if err != nil {
			fmt.Printf("scan failed, err:%v", err)
			return friends
		}
		friends = append(friends, friend)
	}
	return friends
}
//GetName return user's name
func GetName(DB *sql.DB, ID string, isUser bool) string {
	user := new(User)
	if isUser {
		queryString := fmt.Sprintf("select name from go_users where userID=%s;", ID)
		rows, _ := DB.Query(queryString)
		defer rows.Close()
		rows.Next()
		if err := rows.Scan(&user.Name); err != nil {
			fmt.Printf("scan failed, err:%v", err)
			return "failed scan"
		}
	} else {
		queryString := fmt.Sprintf("select group_name from go_groups where group_id=%s;", ID)
		rows, _ := DB.Query(queryString)
		defer rows.Close()
		rows.Next()
		if err := rows.Scan(&user.Name); err != nil {
			fmt.Printf("scan failed, err:%v", err)
			return "failed scan"
		}
	}
	return user.Name
}
//GetUserGroups 返回用户的群聊列表
func GetUserGroups(DB *sql.DB, userID string) (groups []string, res string) {
	defer func() {
		errR := recover()
		if errR != nil {
			res = fmt.Sprint("未成功获取群聊列表")
			fmt.Println("发生了什么？", errR)
		}
	}()
	queryString := fmt.Sprintf("select group_id from go_group_user where userID=%s;", userID)
	var gr string
	rows, err := DB.Query(queryString)
	if err != nil {
		return
	}
	defer rows.Close()
	for  rows.Next() {
		err = rows.Scan(&gr)
		if err != nil {
			res = fmt.Sprintf("scan failed, err:%v", err)
			return
		}
		groups = append(groups, gr)
	}
	return
}
//GetAllMember 返回该组的所有成员
func GetAllMember(DB *sql.DB, groupID string) []string {
	userIDs := make([]string, 10, 20)
	queryString := fmt.Sprintf("select userID from go_group_user where group_id=%s;", groupID)
	rows, err := DB.Query(queryString)
	if err != nil {
		return make([]string,1)
	}
	var userID string
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return userIDs
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs
}
//GetApplyUserS 获取群聊申请表单
func GetApplyUserS(DB *sql.DB, userID string) []string {
	queryString := fmt.Sprintf("select userID, group_id from go_group_apply where group_admin='%s'", userID)
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Println("获取群聊申请表时发送错误:",err)
		return make([]string,0)
	}
	defer rows.Close()
	groupID, applyID := "", ""
	var res []string
	for rows.Next() {
		err = rows.Scan(&applyID, &groupID)
		if err != nil {
			fmt.Println("获取时发生错误:",err)
			return make([]string,0)
		}
		res = append(res, applyID+"申请加入"+groupID+"\n")
	}
	return res
}
//JudgeUserAdmin 判断用户是否为管理员
func JudgeUserAdmin(DB *sql.DB, userID string, groupID string) (res bool) {
	defer func() {
		errR := recover()
		if errR != nil {
			fmt.Println(errR)
		}
	}()
	queryString := fmt.Sprintf("select adminId from go_groups where group_id=%s;", groupID)
	rows, _ := DB.Query(queryString)
	defer rows.Close()
	rows.Next()
	var adminId string
	_ = rows.Scan(&adminId)
	if userID == adminId {
		res = true
	} else {
		res = false
	}
	return
}
//JudgeGroupExit 判断该群组是否存在
func JudgeGroupExit(DB *sql.DB, groupID string) (res bool) {
	defer func() {
		errR := recover()
		if errR != nil {
			fmt.Println(errR)
		}
	}()
	queryString := fmt.Sprintf("select group_id from go_groups where group_id=%s;", groupID)
	rows, _ := DB.Query(queryString)
	defer rows.Close()
	rows.Next()
	var gr string
	_ = rows.Scan(&gr)
	if groupID == gr {
		res = true
	} else {
		res = false
	}
	return
}
//JudgeUserExit 判断userID 是否存在雨数据库中，返回bool型
func JudgeUserExit(DB *sql.DB, userID string) (res bool) {
	queryString := fmt.Sprintf("select userID from go_users where userID=%s;", userID)
	rows, _ := DB.Query(queryString)
	rows.Next()
	var getString string
	_ = rows.Scan(&getString)
	if len(getString) > 0 {
		res = true
	} else {
		res = false
	}
	return
}
//DeleteUser 删除
func DeleteUser(DB *sql.DB, userID string) {
	queryString := fmt.Sprintf("delete from go_group_user where userID=%s;", userID)
	rows, _ := DB.Query(queryString)
	defer rows.Close()
	return
}
//CreateGroupUser 创建群聊-用户 对应表
func CreateGroupUser(DB *sql.DB) {
	rows, err := DB.Query(
		`
		CREATE TABLE go_group_user (
			group_id VARCHAR(80) NOT NULL,
			userID VARCHAR(80) NOT NULL
			);
		`)
	if err == nil {
		fmt.Println("成功创建群聊-用户表")
	} else {
		fmt.Println("群聊-用户表已存在")
		return
	}
	defer rows.Close()
}
//CreateGroup 创建一个群聊
func CreateGroup(DB *sql.DB, userID string, groupID string, groupName string) {
	defer func() {
		errR := recover()
		if errR != nil {
			fmt.Println("创建群聊时发生错误，err:", errR)
			return
		}
	}()
	queryString := fmt.Sprintf("insert into go_groups (group_id, group_name, adminId) values (%s,\"%s\",%s);", groupID, groupName, userID)
	rows, err := DB.Query(queryString)
	
	
	if err == nil {
		fmt.Println("成功创建该群聊")
		InsertMember(DB, userID, groupID)
	} else {
		fmt.Println("创建群聊失败")
		return
	}
	defer rows.Close()
}
//CreateOffLineMsg 创建离线消息缓存表
func CreateOffLineMsg(DB *sql.DB){
	rows, err := DB.Query(
		`
		CREATE TABLE go_offline_msg (
			srcID VARCHAR(80) NOT NULL,
			dstID VARCHAR(80) NOT NULL,
			msg VARCHAR(200) NOT NULL,
			sendTime DATETIME NOT NULL
			);
		`)
	if err == nil {
		fmt.Println("成功创建离线消息表")
	} else {
		fmt.Println("离线消息表已存在")
		return
	}
	defer rows.Close()
}
//CreateGroupTable 创建群聊表
func CreateGroupTable(DB *sql.DB) {
	var rows, err = DB.Query(
		`
		CREATE TABLE go_groups (
			group_id VARCHAR(80) NOT NULL,
			group_name VARCHAR(45) NOT NULL,
			adminId VARCHAR(80) NOT NULL,
			PRIMARY KEY (group_id)
			);
		`)
	
	if err == nil {
		fmt.Println("成功创建群聊表")
	} else {
		fmt.Println("群聊表已存在")
		return
	}
	defer rows.Close()
	
}
//CreateGroupApply 该表保存的是加入群聊的申请
func CreateGroupApply(DB *sql.DB){
	rows, err := DB.Query(
		`
		CREATE TABLE go_group_apply (
			group_id VARCHAR(80) NOT NULL,
			group_admin VARCHAR(80) NOT NULL,
			userID VARCHAR(80) NOT NULL
			);
		`)
	if err == nil {
		fmt.Println("成功创建群聊申请表")
	} else {
		fmt.Println("群聊申请表已存在")
		return
	}
	defer rows.Close()
}
//CreateTable 创建一个用户表
func CreateTable(DB *sql.DB) {
	defer func() {
		errR := recover()
		if errR != nil {
			fmt.Println("已创建")
		}
	}()
	rows, _ := DB.Query(
		`create table go_users(
		name varchar(40) not null,
		age int not null,
		location varchar(40),
		sex varchar(20),
		userID varchar(80) not null,
		pw varchar(80) not null,
		primary key (userID));`)
	defer rows.Close()
}
//CreateFriendTable create a table of friends after register
func CreateFriendTable(DB *sql.DB, userID string) {
	queryString := fmt.Sprintf("create table friends_%s (friendID int);", userID)
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Println("创建好友列表失败")
		return
	}
	defer rows.Close()
}
//CreateChatRecord 聊天记录表
func CreateChatRecord(DB *sql.DB){
	rows, err := DB.Query(
		`
		CREATE TABLE go_chat_record (
			srcID VARCHAR(80) NOT NULL,
			dstID VARCHAR(80) NOT NULL,
			msg VARCHAR(200) NOT NULL,
			sendTime DATETIME NOT NULL
			);
		`)
	if err == nil {
		fmt.Println("成功创建消息记录表")
	} else {
		fmt.Println("消息记录表已存在")
		return
	}
	defer rows.Close()
}