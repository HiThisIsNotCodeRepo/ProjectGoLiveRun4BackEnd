package main

import (
	"fmt"
	"time"
)

func main() {
	if time.Now().Hour() < 12 {
		newTime := time.Now().Add(time.Hour * (12 - time.Duration(time.Now().Hour())))
		fmt.Println(newTime)
	} else {
		newTime := time.Now().Add(-time.Hour * (time.Duration(time.Now().Hour()) - 12))
		fmt.Println(newTime)
	}
}
