package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/pacerino/pr0gramm_music_backend/pkg"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/now"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

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
	// Read-all
	router.HandleFunc("/stats", getStats).Methods("GET")
	// Crawl Links by ID
	router.HandleFunc("/crawl/{Id}", crawlLinks).Methods("GET")
	// Crawl Links by ID
	router.HandleFunc("/info/{Id}", getLinks).Methods("GET")
	initDB()

	log.Println("Listen to :8080")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
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

func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bearer := "Bearer " + os.Getenv("ACR_JWT")
	client := &http.Client{}
	var dateRange = DateRange{}
	if start, err := time.Parse("2006-01-02", r.FormValue("start")); err != nil {
		dateRange.Start = now.BeginningOfMonth().Format("2006-01-02")
	} else {
		dateRange.Start = start.Format("2006-01-02")
	}

	if end, err := time.Parse("2006-01-02", r.FormValue("end")); err != nil {
		dateRange.End = now.EndOfMonth().Format("2006-01-02")
	} else {
		dateRange.End = end.Format("2006-01-02")
	}
	fmt.Println(dateRange.Start, dateRange.End)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://eu-api-v2.acrcloud.com/api/base-projects/27121/day-stat?start=%s&end=%s", dateRange.Start, dateRange.End), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Android 4.4; Tablet; rv:41.0) Gecko/41.0 Firefox/41.0")
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	io.Copy(w, resp.Body)
}

func getLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var metadata Metadata
	params := mux.Vars(r)
	id := params["Id"]
	db.Where("acrID = ?", id).First(&metadata)
	json.NewEncoder(w).Encode(metadata)
}

func crawlLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["Id"]
	var Links CrawlLinks
	var response ApiResponse
	var ahaResponse AHAAPIResponse

	if len(id) > 0 {
		reqLink := fmt.Sprintf(`https://aha-music.com/%s`, id)

		page := rod.New().ControlURL(os.Getenv("CHROME_WS")).MustConnect().MustPage(reqLink)
		section := page.MustWaitLoad().MustElements("a.resource-external-link")
		spotifyRegex, err := regexp.Compile(`(?mi)^(https:\/\/open.spotify.com\/track\/)(.*)$`)
		if err != nil {
			log.Fatal(err)
		}

		deezerRegex, err := regexp.Compile(`(?mi)^https?:\/\/(?:www\.)?deezer\.com\/(track|album|playlist)\/(\d+)$`)
		if err != nil {
			log.Fatal(err)
		}

		for _, s := range section {
			link := s.MustProperty("href").String()
			if spotifyLink := spotifyRegex.FindString(link); len(spotifyLink) > 0 {
				Links.Spotify = spotifyLink
			}

			if deezerLink := deezerRegex.FindString(link); len(deezerLink) > 0 {
				Links.Deezer = deezerLink
			}

		}
		if len(Links.Spotify) > 0 {
			ahaResponse = callAHAAPI(Links.Spotify)
			response.Success = true
		} else if len(Links.Deezer) > 0 {
			ahaResponse = callAHAAPI(Links.Deezer)
			response.Success = true
		} else {
			response.Success = false
			response.Message = "Could not find Deezer or Spotify link"
		}
	} else {
		response.Success = false
		response.Message = "Missing ID"
	}
	if response.Success {
		response.Data = ahaResponse
	}
	json.NewEncoder(w).Encode(response)
}

func callAHAAPI(sourceLink string) AHAAPIResponse {
	requestString := fmt.Sprintf(`https://metadata.aha-music.com/v1-alpha.1/links?url=%s`, sourceLink)
	page := rod.New().ControlURL(os.Getenv("CHROME_WS")).MustConnect().MustPage(requestString)
	preTag, err := page.MustElement("pre").Text()
	if err != nil {
		log.Fatal(err)
	}
	response := AHAAPIResponse{}
	jsonErr := json.Unmarshal([]byte(preTag), &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return response
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
