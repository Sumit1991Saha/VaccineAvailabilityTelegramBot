package main

import (
	"VaccineAvailabilityTelegramBot/cowin"
)

func main() {
	cowin.FetchStateDataAlongWithDistrictData()
	//telegram.StartServiceUsingGetUpdates()
	//telegram.StartServiceUsingWebhook()
}

