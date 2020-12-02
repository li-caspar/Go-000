package main

import (
	"app/service"
	"fmt"
	"log"
	"os"
)

func main() {
	user, err := service.GetUser(1)
	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
	fmt.Printf("db select user name:%s", user.Name)
}
