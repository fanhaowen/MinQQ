package coreFunc

import (
	"crypto/md5"
	"fmt"
)
func MD5(str string) (md5str string) {
	md5str = str
	slat := "))*&^^$^%$aaa3@e"
	for i:=0; i<5;i++ {
		md5str += slat
		data := []byte(md5str)
		has := md5.Sum(data)
		md5str = fmt.Sprintf("%x", has)
	}
	return
}