package api

import (
	"VaccineAvailabilityTelegramBot/utils"
	"strconv"
	"time"
)

func GetUrlForAllStateData() string {
	return "https://cdn-api.co-vin.in/api/v2/admin/location/states"
}

func GetUrlByStateIdForDistrictData(stateId int) string {
	return "https://cdn-api.co-vin.in/api/v2/admin/location/districts/" + strconv.Itoa(stateId)
}

func GetUrlByDistrictId(districtId int) string {
	dateInDDMMYYYYFormat := utils.GetDateInDDMMYYYYFormat(time.Now())
	return "https://cdn-api.co-vin.in/api/v2/appointment/sessions/calendarByDistrict?" +
		"district_id=" + strconv.Itoa(districtId) +
		"&date=" + dateInDDMMYYYYFormat
}

func GetUrlByPincode(pincode int) string {
	dateInDDMMYYYYFormat := utils.GetDateInDDMMYYYYFormat(time.Now())
	return "https://cdn-api.co-vin.in/api/v2/appointment/sessions/calendarByPin?" +
		"pincode=" + strconv.Itoa(pincode) +
		"&date=" + dateInDDMMYYYYFormat
}
