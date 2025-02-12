package main

//gorrila mux
import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Game struct {
	ID        string     `json:"id"`
	GID       string     `json:"gid"`
	Title     string     `json:"title"`
	Publisher *Publisher `json:"publisher"`
}

type Publisher struct {
	Name string `json:"name"`
}

var games []Game

func getGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func deleteGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range games {
		if item.ID == params["id"] {
			//append everything before this index and after this index where the ID is found
			games = append(games[:index], games[index+1:]...)
			break
		}
	}
	//return rest of the games
	json.NewEncoder(w).Encode(games)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range games {
		if item.ID == params["id"] {
			//if ID equals simply return that item
			json.NewEncoder(w).Encode(item)
			break
		}
	}
}

func createGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var game Game
	_ = json.NewDecoder(r.Body).Decode(&game)
	game.ID = strconv.Itoa(rand.Intn(1000000000000000000))
	games = append(games, game)
	json.NewEncoder(w).Encode(game)
}

func updateGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//loop over the games and delete the game with the ID sent
	for index, item := range games {
		if item.ID == params["id"] {
			games = append(games[:index], games[index+1:]...)
			//add a new games with the requested updates in the body of the UPDATE
			var game Game
			_ = json.NewDecoder(r.Body).Decode(&game)
			game.ID = params["id"]
			games = append(games, game)
			json.NewEncoder(w).Encode(game)
			break
		}
	}
}

func main() {
	r := mux.NewRouter()

	games = append(games, Game{ID: "1", GID: "48", Title: "Sekiro", Publisher: &Publisher{Name: "FromSoft"}})
	games = append(games, Game{ID: "2", GID: "14", Title: "Terraria", Publisher: &Publisher{Name: "ReLogic"}})
	r.HandleFunc("/games", getGames).Methods("GET")
	r.HandleFunc("/games/{id}", getGame).Methods("GET")
	r.HandleFunc("/games", createGame).Methods("POST")
	r.HandleFunc("/games/{id}", updateGame).Methods("PUT")
	r.HandleFunc("/games/{id}", deleteGame).Methods("DELETE")

	fmt.Println("Server starting at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
