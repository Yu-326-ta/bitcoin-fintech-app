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
	fmt.Println(apiClient.GetBalance())
}
