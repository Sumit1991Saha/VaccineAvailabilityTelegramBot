package utils

import (
	"VaccineAvailabilityTelegramBot/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

const maxLen = 4096

func ReadTokenFromFile() string  {
	byteSlice, err := ioutil.ReadFile("BotToken.txt")
	if err != nil {
		fmt.Println("Error :- ", err)
		os.Exit(1)
	}
	return string(byteSlice)
}

func GetDateInDDMMYYYYFormat(now time.Time) string {
	currentDate := now.String()
	dateInYYYYMMDDFormat := strings.Split(currentDate, " ")[0]
	date := strings.Split(dateInYYYYMMDDFormat, "-")
	dateInDDMMYYYYFormat := date[2] + "-" + date[1] + "-" + date[0]
	return dateInDDMMYYYYFormat
}

func ApiCall(url, method string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error (ApiCall - NewRequest) :- ", err)
		return nil, err
	}
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error (ApiCall - Do) :- ", err)
		return nil, err
	}

	return response, err
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

func GetPublicUrlNgrok(response *http.Response) string {
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error calling api :- " , response.Status)
		return ""
	}

	var jsonData models.ResponseTunnels

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error ReadData (ReadAll) :- ", err)
		return ""
	}
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		fmt.Println("Error ReadData (Unmarshal) :- ", err)
		return ""
	}

	return jsonData.Tunnel[0].PublicUrl
}
