package models

type Visitor struct {
	Id        int    `json:"visitor_id"`
	HotelId   int    `json:"hotel_id"`
	HotelRoom int    `json:"hotel_room_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}
