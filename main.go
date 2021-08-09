package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/pacerino/pr0gramm_music_backend/pkg"

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

var db *gorm.DB

func initDB() {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE"))
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&Items{})
}

func main() {
	godotenv.Load()
	router := mux.NewRouter()
	// Read
	router.HandleFunc("/item/{Id}", getItem).Methods("GET")
	// Read-all
	router.HandleFunc("/items", getItems).Methods("GET")
	initDB()

	log.Println("Listen to :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
	defer db.Close()
}

func getItems(w http.ResponseWriter, r *http.Request) {
	var items []Items
	var pagination pkg.Pagination

	w.Header().Set("Content-Type", "application/json")

	if limit, _ := strconv.Atoi(r.FormValue("limit")); limit > 0 {
		pagination.Limit = limit
	} else {
		pagination.Limit = 10
	}

	if page, _ := strconv.Atoi(r.FormValue("page")); page > 0 {
		pagination.Page = page
	} else {
		pagination.Page = 1
	}

	sort := r.FormValue("sort")
	pagination.Sort = sort

	db.Find(&items)
	db.Scopes(paginate(items, &pagination, db)).Find(&items)
	pagination.Rows = items
	json.NewEncoder(w).Encode(pagination)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["Id"]

	var items Items
	db.First(&items, id)
	json.NewEncoder(w).Encode(items)
}

func paginate(value interface{}, pagination *pkg.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Round(float64(totalRows) / float64(pagination.Limit)))
	fmt.Println(pagination.Limit, totalRows)
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
