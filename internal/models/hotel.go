package models

type Hotel struct {
	Id        int    `json:"id"`
	Country   string `json:"country"`
	City      string `json:"city"`
	HotelName string `json:"hotel_name"`
	Stars     int    `json:"stars"`
}
