package cowin

import (
	"VaccineAvailabilityTelegramBot/utils"
	"fmt"
	"strconv"
	"time"
)

func FetchDataByPinCode(pincode int) string {
	dateInDDMMYYYYFormat := utils.GetDateInDDMMYYYYFormat(time.Now())
	url := "https://cdn-api.co-vin.in/api/v2/appointment/sessions/calendarByPin?" +
		"pincode=" + strconv.Itoa(pincode) +
		"&date=" + dateInDDMMYYYYFormat
	return FetchData(url)
}

func FetchDataByDistrictId(districtId int) string {
	dateInDDMMYYYYFormat := utils.GetDateInDDMMYYYYFormat(time.Now())
	//vadodara_corporation_dist_id := 777
	url := "https://cdn-api.co-vin.in/api/v2/appointment/sessions/calendarByDistrict?" +
		"district_id=" + strconv.Itoa(districtId) +
		"&date=" + dateInDDMMYYYYFormat
	return FetchData(url)
}

func FetchData(url string) string {
	response, err := utils.ApiCall(url, "GET")
	if err != nil || response == nil {
		fmt.Println("Error (FetchData) :- ", err)
		return ""
	}
	vacancies := utils.ReadDataForCenters(response)
	return vacancies
	//call("http://pokeapi.co/api/v2/pokedex/kanto/", "GET")
}