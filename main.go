package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

type Tabler interface {
	TableName() string
}

func (Items) TableName() string {
	return "Items"
}

type Items struct {
	// gorm.Model
	Id     int    `json:"id" gorm:"primary_key;autoIncrement;column:id"`
	ItemID int    `json:"itemID" gorm:"not null;column:itemID"`
	Title  string `json:"title" gorm:"column:title"`
	Album  string `json:"album" gorm:"column:album"`
	Artist string `json:"artist" gorm:"column:artist"`
	Url    string `json:"url" gorm:"column:url"`
	NoData int    `json:"noData" gorm:"column:noData"`
}

type Votes struct {
	gorm.Model
	SauceID uint   `json:"sauceID" gorm:"not null"`
	VoteID  string `json:"voteID" gorm:"not null;unique"`
}

type GormErr struct {
	Number  int    `json:"Number"`
	Message string `json:"Message"`
}

var db *gorm.DB

func initDB() {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE"))
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&Items{}, &Votes{})
}

func main() {
	godotenv.Load()
	router := mux.NewRouter()
	// Read
	router.HandleFunc("/items/{Id}", getItem).Methods("GET")
	// Read-all
	router.HandleFunc("/items", getItems).Methods("GET")
	// Put-Vote
	router.HandleFunc("/vote", putVote).Methods("POST")
	initDB()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var items []Items
	db.Find(&items)
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["Id"]

	var items Items
	db.First(&items, id)
	json.NewEncoder(w).Encode(items)
}

func putVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var vote Votes
	json.NewDecoder(r.Body).Decode(&vote)
	if err := db.Create(&Votes{SauceID: vote.SauceID, VoteID: fmt.Sprintf("%s_%d", vote.VoteID, vote.SauceID)}); err != nil {
		// Source: https://github.com/go-gorm/gorm/issues/4037#issuecomment-881834378
		byteErr, _ := json.Marshal(err.Error)
		var newError GormErr
		json.Unmarshal((byteErr), &newError)
		if newError.Number == 1062 {
			json.NewEncoder(w).Encode(ErrAlreadyVoted)
		} else {
			log.Fatalln(newError)
		}
	} else {
		json.NewEncoder(w).Encode(vote)
	}
}
