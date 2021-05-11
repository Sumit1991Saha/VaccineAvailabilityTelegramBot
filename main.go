package main

import (
	"VaccineAvailabilityTelegramBot/telegram"
)

func main() {
	//telegram.StartServiceUsingGetUpdates()
	telegram.StartServiceUsingWebhook()

	//cowin.FetchStateDataAlongWithDistrictData()
}

