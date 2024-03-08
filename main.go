package main

import (
	"UTS/controllers"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/rooms", controllers.GetAllRooms).Methods("GET")
	router.HandleFunc("/detailroom", controllers.GetDetailRooms).Methods("GET")
	router.HandleFunc("/insertroom", controllers.InsertRooms).Methods("POST")
	router.HandleFunc("/leaveroom", controllers.LeaveRoom).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
