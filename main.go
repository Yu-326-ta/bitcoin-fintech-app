package main

import (
	"bitcoin/bitflyer"
	"bitcoin/config"
	"bitcoin/utils"
	"fmt"
)

func main() {
	utils.LoggingSetting(config.Config.Logfile)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSercet)
	
	order := &bitflyer.Order{
		ProductCode:     config.Config.ProductCode,
		ChildOrderType:  "MARKET",
		Side:            "BUY",
		Size:            0.0001,
		MinuteToExpires: 1,
		TimeInForce:     "GTC",
	}
	res, _ := apiClient.SendOrder(order)
	fmt.Println(res)

}
