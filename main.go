package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"os"
)

type moodOperatorStruct struct {
	Username string `json:"username"`
	Mood string `json:"mood"`
	Operation string `json:"operation"`
}

var (
	clients = make(map[*websocket.Conn]bool)
	broadcast = make(chan *moodOperatorStruct)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	db = Database()
	port = flag.String("port", "8844", "port to listen on")

)

type App struct {
	Port string
	Logf func(string, ...interface{})
}

//go:embed index.html
var index string

func main() {
	flag.Parse()
	app := App{
		Port: *port,
	}
	port := os.Getenv("PORT")
	if port != "" {
		app.Port = port
	}

	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/mood", addMoodHandler).Methods("POST")
	router.HandleFunc("/ws", wsHandler)
	router.HandleFunc("/all", allHandler).Methods("GET")
	router.HandleFunc("/delete", deleteMoodHandler).Methods("POST")
	go echo()

	log.Printf("Now open http://localhost:%s", app.Port)
	addr := net.JoinHostPort("", app.Port)
	log.Fatal(http.ListenAndServe(addr, router))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, index)
}

func writer(mood *moodOperatorStruct) {
	broadcast <- mood
}

func addMoodHandler(w http.ResponseWriter, r *http.Request) {
	var mood UserMoodStruct
	if err := json.NewDecoder(r.Body).Decode(&mood); err != nil {
		log.Printf("ERROR: %s", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	log.Println(mood)
	Save(mood, db)
	defer r.Body.Close()
	go writer(&moodOperatorStruct{Username: mood.Username, Mood: mood.Mood, Operation: "Save"})
}

func deleteMoodHandler(w http.ResponseWriter, r * http.Request) {
	var mood UserMoodStruct
	if err := json.NewDecoder(r.Body).Decode(&mood); err != nil {
		log.Printf("ERROR: %s", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	Delete(mood.Username, db)
	defer r.Body.Close()
	go writer(&moodOperatorStruct{Username: mood.Username, Mood: mood.Mood, Operation: "Delete"})
}


func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client
	clients[ws] = true
}

func allHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	userMoods := GetAll(db)
	json.NewEncoder(w).Encode(userMoods)
}

func echo() {
	for {
		val := <-broadcast
		mood, _ := json.Marshal(val)
		// send to every client that is currently connected
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, mood)
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}