package controllers

import (
	m "UTS/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func sendMessage(w http.ResponseWriter, code int, message string){
	var response m.NormalResponse
	response.Status = code
	response.Message = message
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAllRooms(w http.ResponseWriter, r *http.Request){
	db := connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, room_name FROM Rooms")
	if err != nil {
		log.Println(err)
		sendMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	var room m.RoomsTwo
	var rooms []m.RoomsTwo
	for rows.Next() {
		if err := rows.Scan(&room.ID, &room.Name); err != nil {
			log.Println(err)
			sendMessage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		} else {
			rooms = append(rooms, room)
		}
	}
	w.Header().Set("content-Type", "application/json")

	var response m.RoomsResponse
	response.Status = 200
	response.Data.Data = rooms
	json.NewEncoder(w).Encode(response)
}

func GetDetailRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	// id_room, _ := strconv.Atoi(r.Form.Get("id_room"))
	temp := r.URL.Query()["id_room"]
	// Tidak dapat menggunakan Form.Get / FormValue dikarenakan selalu mendapatkan nilai nil dari inputan
	fmt.Println("Room ID: ", temp[0])

	rows, err := db.Query(`SELECT Rooms.id, room_name, Participants.id AS participant_id, id_account, username
		FROM Rooms
		JOIN Participants ON Rooms.id = Participants.id_room
		JOIN Accounts ON Participants.id_account = Accounts.id
		WHERE Rooms.id = ?`, temp[0])

	if err != nil {
		log.Println("Error querying database:", err)
		sendMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var response m.ResultOne
	var room m.Room
	
	for rows.Next(){
		var participant m.Participant
		if err := rows.Scan(
			&room.ID,
			&room.RoomName,
			&participant.ID,
			&participant.IDAccount,
			&participant.Username,
		); err != nil {
			log.Println(err)
			sendMessage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		room.Participants = append(room.Participants, participant)
	}

	w.Header().Set("Content-Type", "application/json")
	response.Status = 200
	response.Data.Data = room
	json.NewEncoder(w).Encode(response)
}

func InsertRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var err error
	
	// id_room, _ := strconv.Atoi(r.Form.Get("id_room"))
	// id_account, _ := strconv.Atoi(r.Form.Get("id_account"))
	// Tidak dapat menggunakan Form.Get / FormValue dikarenakan selalu mendapatkan nilai nil dari inputan
	
	id_room := r.URL.Query()["id_room"]
	id_account := r.URL.Query()["id_account"]

	fmt.Println("Room ID: ", id_room[0])
	fmt.Println("Account ID: ", id_account[0])

	var maxPlayer int
    err = db.QueryRow("SELECT max_player FROM Games WHERE id = (SELECT id_game FROM Rooms WHERE id = ?)", id_room[0]).Scan(&maxPlayer)
    if err != nil {
        log.Println("Error querying max player:", err)
        sendMessage(w, http.StatusInternalServerError, "Internal Server Error")
        return
    }

    var currentParticipants int
    err = db.QueryRow("SELECT COUNT(*) FROM Participants WHERE id_room = ?", id_room[0]).Scan(&currentParticipants)
    if err != nil {
        log.Println("Error querying current participants:", err)
        sendMessage(w, http.StatusInternalServerError, "Internal Server Error")
        return
    }

    if currentParticipants >= maxPlayer {
        sendMessage(w, http.StatusForbidden, "Room is full. Cannot join.")
        return
    }

    _, err = db.Exec("INSERT INTO Participants (id_room, id_account) VALUES (?, ?)", id_room[0], id_account[0])
    if err != nil {
        log.Println("Error inserting participant:", err)
        sendMessage(w, http.StatusInternalServerError, "Internal Server Error")
        return
    }

	var response m.InsertRoomResponse
	response.Status = 200
	response.Message = "Sukses Membuat Room"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LeaveRoom(w http.ResponseWriter, r *http.Request) {
    db := connect()
    defer db.Close()
	var err error

    // id_room, _ := strconv.Atoi(r.Form.Get("id_room"))
	// id_account, _ := strconv.Atoi(r.Form.Get("id_account"))
	// Tidak dapat menggunakan Form.Get / FormValue dikarenakan selalu mendapatkan nilai nil dari inputan
	
	id_room := r.URL.Query()["id_room"]
	id_account := r.URL.Query()["id_account"]

	fmt.Println("Room ID: ", id_room[0])
	fmt.Println("Account ID: ", id_account[0])

    _, err = db.Exec("DELETE FROM Participants WHERE id_room = ? AND id_account = ?", id_room[0], id_account[0])
    if err != nil {
        log.Println("Error leaving room:", err)
        sendMessage(w, http.StatusInternalServerError, "Internal Server Error")
        return
    }

    var response m.LeaveRoomResponse
    response.Status = 200
    response.Message = "Keluar room sukses"
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}