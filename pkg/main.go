package main

import (
	"fmt"
	"time"
)

func main() {

	loc, _ := time.LoadLocation("Europe/Moscow")

	fmt.Println(time.Now().In(loc).Round(time.Minute))

}
