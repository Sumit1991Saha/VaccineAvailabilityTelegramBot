package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const maxLen = 4096

func main() {
	usingGetUpdates()
}

func readTokenFromFile() string  {
	byteSlice, err := ioutil.ReadFile("BotToken.txt")
	if err != nil {
		fmt.Println("Error :- ", err)
		os.Exit(1)
	}
	return string(byteSlice)
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

func GetDateInDDMMYYYYFormat(now time.Time) string {
	currentDate := now.String()
	dateInYYYYMMDDFormat := strings.Split(currentDate, " ")[0]
	date := strings.Split(dateInYYYYMMDDFormat, "-")
	dateInDDMMYYYYFormat := date[2] + "-" + date[1] + "-" + date[0]
	return dateInDDMMYYYYFormat
}

func ApiCall(url, method string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		_ = fmt.Errorf("Got error %s", err.Error())
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		_ = fmt.Errorf("Got error %s", err.Error())
		return "", err
	}

	data := ReadData(response)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(response.Body)
	return data, nil
}

func ReadData(response *http.Response) string {
	var jsonData Response
	err := json.NewDecoder(response.Body).Decode(&jsonData)
	if err != nil {
		return ""
	}
	return ParseJsonData(jsonData)
}

func ParseJsonData(jsonData Response) string {
	messageToBeReturned := ""
	for i := 0; i < len(jsonData.Centers); i++ {
		center := jsonData.Centers[i]
		sessions := center.Sessions
		for j := 0; j < len(sessions); j++ {
			session := sessions[j]
			noOfDosesAvailable := session.AvailableCapacity
			if noOfDosesAvailable > 0 {
				msg := "Vaccine available on :- " + session.Date + "\n" +
					"AvailableCapacity :- " + strconv.Itoa(noOfDosesAvailable) + "\n" +
					"At center :- " + center.Name + "\n" +
					"Address :- " + center.Address + "\n"
				if session.MinAgeLimit >= 45 {
					msg = "For 45+ years, " + "\n" + msg
				} else {
					msg = "For 18+ years, " + "\n" + msg
				}
				messageToBeReturned = messageToBeReturned + msg + "\n"
			}
		}
	}
	if messageToBeReturned == "" {
		messageToBeReturned = "No slots available, PLease try again tomorrow"
	}
	return messageToBeReturned
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

func SplitDataInChunks(dataToSend string) []string {
	splits := []string{}
	var l, r int
	for l, r = 0, maxLen; r < len(dataToSend); l, r = r, r+maxLen {
		for !utf8.RuneStart(dataToSend[r]) {
			r--
		}
		splits = append(splits, dataToSend[l:r])
	}
	splits = append(splits, dataToSend[l:])
	return splits
}

