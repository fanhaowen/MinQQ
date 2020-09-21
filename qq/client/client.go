package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	run()
	for {
		fmt.Println("You could try again,if you input AGAIN")
		var msgWriter string
		Scanf(&msgWriter)
		if msgWriter == "AGAIN" || msgWriter == "again" {
			run()
		} else {
			fmt.Println("You exit this program because you didn't input AGAIN")
			break
		}
	}
}

func run() {
	conn, err := net.Dial("tcp", "localhost:8080")
	// fmt.Printf("%T", conn)
	if err != nil {
		fmt.Println("err:", err)
	}
	errFlag := false
	defer conn.Close()
	buf := make([]byte, 1024)
	go readFromServer(conn, buf, &errFlag)
	for {
		var msgWriter string
		if errFlag {
			break
		}
		Scanf(&msgWriter)
		// fmt.Println(msg)
		if msgWriter == "EXIT" {
			conn.Write([]byte(msgWriter))
			time.Sleep(time.Second)
			conn.Close()
			return
		}
		conn.Write([]byte(msgWriter))
	}
}

//Scanf 代替之前的fmt.Scanf 可以输入空格
func Scanf(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*a = string(data)
}

func readFromServer(conn net.Conn, buf []byte, errFlag *bool) {

	reader, err := conn.Read(buf)
	if err != nil {
		fmt.Println("CONNECT FAILED!")
		*errFlag = true
		fmt.Println("An Error Occured.	*ENTER to CONTINUE")
		return
	}
	msgReader := string(buf[:reader])
	if msgReader != "" {
		fmt.Print(msgReader)
	}
	time.Sleep(time.Second)

	if !*errFlag {
		readFromServer(conn, buf, errFlag)
	}
}
