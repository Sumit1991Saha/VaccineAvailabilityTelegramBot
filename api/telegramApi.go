package api

import "VaccineAvailabilityTelegramBot/utils"

func GetUrlToSetWebhook() string {
	botToken := utils.ReadTokenFromFile()
	return "https://api.telegram.org/bot" + botToken + "/setWebhook"
}

func GetUrlToDeleteWebhook() string {
	botToken := utils.ReadTokenFromFile()
	return "https://api.telegram.org/bot" + botToken + "/deleteWebhook"
}
