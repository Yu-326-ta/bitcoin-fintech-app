package main

import (
	"bitcoin/config"
	"bitcoin/utils"
	"fmt"
	"log"
)

func main() {
	utils.LoggingSetting(config.Config.Logfile)
	log.Println("test")
	fmt.Println(config.Config.ApiKey)
	fmt.Println(config.Config.ApiSercet)
}
