package main

import (
	"bitcoin/app/controllers"
	"bitcoin/app/models"
	"bitcoin/config"
	"bitcoin/utils"
	"fmt"
	"log"
	"time"
)

func main() {
	df, _ := models.GetAllCandle(config.Config.ProductCode, time.Minute, 365)
	fmt.Printf("%+v¥n", df.OptimizeParams())
	utils.LoggingSetting(config.Config.Logfile)
	controllers.StreamIngestionData()
	log.Println(controllers.StartWebServer())
}

