package main

import (
	"bitcoin/app/controllers"
	"bitcoin/app/models"
	"bitcoin/config"
	"bitcoin/utils"
	"fmt"
)

func main() {
	utils.LoggingSetting(config.Config.Logfile)
	fmt.Println(models.DbConnection)
	controllers.StreamIngestionData()

}
