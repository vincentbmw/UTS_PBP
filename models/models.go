package models

type NormalResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Rooms struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	IDGame int    `json:"id_game"`
}

type RoomsTwo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RoomResponse struct {
	Data []RoomsTwo `json:"rooms_data"`
}

type RoomsResponse struct {
	Status int          `json:"status"`
	Data   RoomResponse `json:"data"`
}

type Participant struct {
	ID        int    `json:"id"`
	IDAccount int    `json:"id_account"`
	Username  string `json:"username"`
}

type Room struct {
	ID           int           `json:"id"`
	RoomName     string        `json:"room_name"`
	Participants []Participant `json:"participants"`
}

type ResultOne struct {
	Status int       `json:"status"`
	Data   ResultTwo `json:"data"`
}

type ResultTwo struct {
	Data Room `json:"room"`
}

type InsertRoomResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type LeaveRoomResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}