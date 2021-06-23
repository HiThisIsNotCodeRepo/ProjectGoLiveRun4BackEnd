package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strings"
)

func main() {
	accessKey := uuid.NewV4().String()
	fmt.Println(len(accessKey))
	fmt.Println(accessKey)
	strArr := strings.Split(accessKey, "-")
	lastStrByteArr := []byte(strArr[len(strArr)-1])
	var last4ByteArr []byte
	for i := len(lastStrByteArr) - 4; i < len(lastStrByteArr); i++ {
		last4ByteArr = append(last4ByteArr, lastStrByteArr[i])
	}
	fmt.Println(string(last4ByteArr))
}
