package main

import (
	"bitcoin/bitflyer"
	"bitcoin/config"
	"bitcoin/utils"
	"fmt"
	"time"
)

func main() {
	utils.LoggingSetting(config.Config.Logfile)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSercet)
	
	tickerChannel := make(chan bitflyer.Ticker)
	go apiClient.GetRealTimeTicker(config.Config.ProductCode, tickerChannel)
	for ticker := range tickerChannel {
		fmt.Println(ticker)
		fmt.Println(ticker.GetMidPrice())
		fmt.Println(ticker.DateTime())
		fmt.Println(ticker.TruncateDateTime(time.Second))
		fmt.Println(ticker.TruncateDateTime(time.Minute))
		fmt.Println(ticker.TruncateDateTime(time.Hour))
	}
}
