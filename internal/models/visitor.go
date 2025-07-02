package models

type Visitor struct {
	HotelRoom int    `json:"hotel_room_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}
