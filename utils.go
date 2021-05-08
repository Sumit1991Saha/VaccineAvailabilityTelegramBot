package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const maxLen = 4096

func readTokenFromFile() string  {
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

func ApiCall(url, method string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error (ApiCall - NewRequest) :- ", err)
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error (ApiCall - Do) :- ", err)
		return "", err
	}

	data := ""
	if response.StatusCode == 200 {
		data = ReadDataForCenters(response)
	} else {
		fmt.Println("Error calling api :- " , response.Status)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(response.Body)
	return data, nil
}

func ReadDataForCenters(response *http.Response) string {
	var jsonData ResponseCenters

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

	return ParseJsonData(jsonData)
}

func ParseJsonData(jsonData ResponseCenters) string {
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
