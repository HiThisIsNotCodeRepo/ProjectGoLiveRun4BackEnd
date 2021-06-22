package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
)

func main() {
	accessKey := uuid.NewV4().String()
	fmt.Println(len(accessKey))
}
