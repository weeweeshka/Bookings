package models

type HotelRoom struct {
	Id       int  `json:"id"`
	HotelId  int  `json:"hotels_id"`
	Rooms    int  `json:"rooms"`
	Meals    bool `json:"meals"`
	Bar      bool `json:"bar"`
	Services bool `json:"services"`
	Busy     bool `json:"busy"`
}
