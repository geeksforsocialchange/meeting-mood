package main

import (
	"embed"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"text/template"
)

type moodOperatorStruct struct {
	Username  string `json:"username"`
	Mood      string `json:"mood"`
	Operation string `json:"operation"`
	Room      string `json:"string"`
}

var (
	clients   = make(map[Client]bool)
	broadcast = make(chan *moodOperatorStruct)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	db   = Database()
	port = flag.String("port", "8844", "port to listen on")
)

type Client struct {
	Room   string
	Client *websocket.Conn
}

type App struct {
	Port string
	Logf func(string, ...interface{})
}

//go:embed index.html
var index string

//go:embed room.html
var room string

//go:embed assets/*
var assets embed.FS

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
	router.StrictSlash(true)
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/{room:[0-9]+}", roomHandler).Methods("GET")
	router.HandleFunc("/{room:[0-9]+}/mood", addMoodHandler).Methods("POST")
	router.HandleFunc("/{room:[0-9]+}/ws", wsHandler)
	router.HandleFunc("/{room:[0-9]+}/all", allHandler).Methods("GET")
	router.HandleFunc("/{room:[0-9]+}/delete", deleteMoodHandler).Methods("POST")

	router.PathPrefix("/assets/").Handler(assetsHandler())
	go echo()

	log.Printf("Now open http://localhost:%s", app.Port)
	addr := net.JoinHostPort("", app.Port)
	log.Fatal(http.ListenAndServe(addr, router))
}

func assetsHandler() http.Handler {
	fsys := fs.FS(assets)
	contentStatic, _ := fs.Sub(fsys, "assets")
	return http.StripPrefix("/assets/", http.FileServer(http.FS(contentStatic)))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("template").Delims("[[", "]]").Parse(index))
	footer, _ := ioutil.ReadFile("footer.html")
	tmpl.Execute(w, string(footer))
}

func roomHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("template").Delims("[[", "]]").Parse(room))
	footer, _ := ioutil.ReadFile("footer.html")
	tmpl.Execute(w, string(footer))
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
	Save(mood, db)
	defer r.Body.Close()
	go writer(&moodOperatorStruct{Username: mood.Username, Mood: mood.Mood, Operation: "Save", Room: mood.Room})
}

func deleteMoodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var mood UserMoodStruct
	if err := json.NewDecoder(r.Body).Decode(&mood); err != nil {
		log.Printf("ERROR: %s", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	Delete(mood.RoomUser, db)
	defer r.Body.Close()
	go writer(&moodOperatorStruct{Username: mood.Username, Mood: mood.Mood, Operation: "Delete", Room: vars["room"]})
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client
	clients[Client{Client: ws, Room: vars["room"]}] = true
}

func allHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	userMoods := GetAll(vars["room"], db)
	json.NewEncoder(w).Encode(userMoods)
}

func echo() {
	for {
		val := <-broadcast
		mood, _ := json.Marshal(val)
		// @TODO there's probably a better way of only sending to clients in the current room
		for client := range clients {
			if client.Room == val.Room {
				err := client.Client.WriteMessage(websocket.TextMessage, mood)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Client.Close()
					delete(clients, client)
				}
			}

		}
	}
}
