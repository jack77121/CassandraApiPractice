package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jack77121/CassandraApiPractice/Users"

	"github.com/gorilla/mux"
	"github.com/jack77121/CassandraApiPractice/Cassandra"
)

type heartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(heartbeatResponse{Status: "OK", Code: 200})
}

func main() {
	CassandraSession := Cassandra.Session
	defer CassandraSession.Close()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", heartbeat)
	router.HandleFunc("/users/new", Users.Post)
	router.HandleFunc("/users", Users.Get)
	router.HandleFunc("/users/{user_uuid}", Users.GetOne)
	log.Fatal(http.ListenAndServe(":8080", router))
}
