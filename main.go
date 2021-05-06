package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"time"
)

func main() {
	usingGetUpdates()
}

func usingGetUpdates() {
	botToken := readTokenFromFile()
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if val, err := strconv.Atoi(update.Message.Text); err == nil {
			dataToSend:= ""
			/*if len(update.Message.Text) < 6 {
				dataToSend = "PLease enter a valid pincode"
			} else {
				dataToSend = fetchDataByPinCode(val)
			}*/
			dataToSend = FetchDataByDistrictId(val)
			SendTelegramMessage(bot, update.Message.Chat.ID, dataToSend)
		}
	}
}

func FetchDataByPinCode(pincode int) string {
	dateInDDMMYYYYFormat := GetDateInDDMMYYYYFormat(time.Now())
	url := "https://cdn-api.co-vin.in/api/v2/appointment/sessions/calendarByPin?" +
		"pincode=" + strconv.Itoa(pincode) +
		"&date=" + dateInDDMMYYYYFormat
	return FetchData(url)
}

func FetchDataByDistrictId(districtId int) string {
	dateInDDMMYYYYFormat := GetDateInDDMMYYYYFormat(time.Now())
	//vadodara_corporation_dist_id := 777
	url := "https://cdn-api.co-vin.in/api/v2/appointment/sessions/calendarByDistrict?" +
		"district_id=" + strconv.Itoa(districtId) +
		"&date=" + dateInDDMMYYYYFormat
	return FetchData(url)
}

func FetchData(url string) string {
	vacancies, err := ApiCall(url, "GET")
	if err != nil {
		fmt.Println("Error :- ", err)
		return ""
	}
	return vacancies
	//call("http://pokeapi.co/api/v2/pokedex/kanto/", "GET")
}

func SendTelegramMessage(bot *tgbotapi.BotAPI, chatId int64, dataToSend string) {
	data := SplitDataInChunks(dataToSend)
	for i := 0; i < len(data); i++ {
		msg := tgbotapi.NewMessage(chatId, data[i])
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error while sending message :- ", err)
		}
	}
}



