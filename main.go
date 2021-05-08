package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	//usingGetUpdates()
	usingWebhook()
}

func usingWebhook() {
	http.ListenAndServe(":3000", http.HandlerFunc(Handler))
}

// Handler This handler is called everytime telegram sends us a webhook event
func Handler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	body := &WebhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	if val, err := strconv.Atoi(body.Message.Text); err == nil {
		dataToSend := ""
		/*if len(update.Message.Text) != 6 {
			dataToSend = "PLease enter a valid pincode"
		} else {
			dataToSend = FetchDataByPinCode(val)
		}*/
		dataToSend = FetchDataByDistrictId(val)
		if dataToSend == "" {
			fmt.Println("Error empty string, so not sending data")
		} else {
			fmt.Println(dataToSend)
			SendTelegramMessageUsingWebhook(body.Message.Chat.ID, dataToSend)
		}
	}

	// log a confirmation message if the message is sent successfully
	fmt.Println("reply sent")
}

func usingGetUpdates() {
	botToken := readTokenFromFile()
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("Error (usingGetUpdates) :- ", err)
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
			dataToSend := ""
			/*if len(update.Message.Text) != 6 {
				dataToSend = "PLease enter a valid pincode"
			} else {
				dataToSend = FetchDataByPinCode(val)
			}*/
			dataToSend = FetchDataByDistrictId(val)
			if dataToSend == "" {
				fmt.Println("Error empty string, so not sending data")
			} else {
				SendTelegramMessageUsingBotApi(bot, update.Message.Chat.ID, dataToSend)
			}
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
		fmt.Println("Error (FetchData) :- ", err)
		return ""
	}
	return vacancies
	//call("http://pokeapi.co/api/v2/pokedex/kanto/", "GET")
}

func SendTelegramMessageUsingBotApi(bot *tgbotapi.BotAPI, chatId int64, dataToSend string) {
	data := SplitDataInChunks(dataToSend)
	for i := 0; i < len(data); i++ {
		msg := tgbotapi.NewMessage(chatId, data[i])
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error while sending message :- ", err)
		}
	}
}

func SendTelegramMessageUsingWebhook(chatID int64, dataToSend string) {
	data := SplitDataInChunks(dataToSend)
	for i := 0; i < len(data); i++ {
		// Create the request body struct
		reqBody := &SendMessageReqBody{
			ChatID: chatID,
			Text:   data[i],
		}
		// Create the JSON body from the struct
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Println("Error Marshal:- ", err)
		} else {
			botToken := readTokenFromFile()
			url := "https://api.telegram.org/bot" + botToken + "/sendMessage"

			// Send a post request with your token
			res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
			if err != nil {
				fmt.Println("Error Post:- ", err)
			}
			if res.StatusCode != http.StatusOK {
				fmt.Println("Error non 200 status code:- ", err)
			}
		}
	}
}
