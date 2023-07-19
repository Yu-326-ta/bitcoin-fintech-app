package main

import (
	"bitcoin/app/controllers"
	"bitcoin/config"
	"bitcoin/utils"
	"log"
)

func main() {
	utils.LoggingSetting(config.Config.Logfile)
	controllers.StreamIngestionData()
	log.Println(controllers.StartWebServer())
}

