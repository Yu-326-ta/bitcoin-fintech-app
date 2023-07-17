package main

import (
	"bitcoin/app/controllers"
	"bitcoin/config"
	"bitcoin/utils"
)

func main() {
	utils.LoggingSetting(config.Config.Logfile)
	controllers.StreamIngestionData()
	controllers.StartWebServer()
}
