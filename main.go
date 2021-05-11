package main

import (
	"VaccineAvailabilityTelegramBot/cowin"
	"VaccineAvailabilityTelegramBot/telegram"
)

func main() {
	cowin.FetchStateDataAlongWithDistrictData()
	//telegram.StartServiceUsingGetUpdates()
	telegram.StartServiceUsingWebhook()
}

