package main

import (
	"bitcoin/config"
	"fmt"
)

func main() {
	fmt.Println(config.Config.ApiKey)
	fmt.Println(config.Config.ApiSercet)
}