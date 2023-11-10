package main

import (
	"fmt"
	"log"
	"spider/Pornbest"
)

func main() {
	var Url string
	fmt.Println("请输入网址:")
	_, err := fmt.Scanln(&Url)
	if err != nil {
		log.Fatal(err)
		return
	}
	Pornbest.Start(Url)
}
