package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/onrik/logrus/sentry"
	log "github.com/sirupsen/logrus"

	"github.com/go-rod/rod"
	"github.com/mileusna/crontab"
	"github.com/pacerino/pr0gramm_music_backend/pkg"
	"github.com/pacerino/pr0gramm_music_backend/pr0gramm"

	"github.com/gorilla/mux"
	"github.com/jinzhu/now"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func initDB() {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE"))
	db, err = gorm.Open(mysql.Open(dataSourceName))

	if err != nil {
		log.WithError(err).Fatal("Could not open connection to DB")
	}
	if err := db.AutoMigrate(&Comments{}, &Items{}); err != nil {
		log.WithError(err).Fatal("Could not Migrate Models!")
	}
}

func main() {
	godotenv.Load()
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetReportCaller(true)

	sentryHook, err := sentry.NewHook(sentry.Options{
		Dsn: os.Getenv("SENTRY_DSN"),
	}, log.PanicLevel, log.FatalLevel, log.ErrorLevel)
	if err != nil {
		log.Error(err)
		return
	}
	defer sentryHook.Flush()

	log.AddHook(sentryHook)

	router := mux.NewRouter()
	// Read
	router.HandleFunc("/item/{Id}", getItem).Methods("GET")
	// Read-all
	router.HandleFunc("/items", getItems).Methods("GET")
	// Return Stats
	router.HandleFunc("/stats", getStats).Methods("GET")
	// Crawl Links by ACR ID
	router.HandleFunc("/crawl/{Id}", crawlLinks).Methods("GET")
	// Get Metadata by ACR ID
	router.HandleFunc("/info/{Id}", getLinks).Methods("GET")
	initDB()

	ctab := crontab.New()
	err = ctab.AddJob(os.Getenv("CRONJOB"), goThroughBotComments)
	if err != nil {
		log.WithError(err).Fatal("Could not add cronjob!")
	}

	log.Println("Listen to :8080")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func getItems(w http.ResponseWriter, r *http.Request) {
	var itemResponse []ItemResponse
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
	db.Table("items").Select("items.*, comments.*").Scopes(paginate(items, &pagination, db)).Joins("LEFT JOIN comments ON items.item_id = comments.item_id").Where("items.title OR items.album OR items.artist IS NOT NULL").Find(&itemResponse)
	pagination.Rows = itemResponse
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
	req, err := http.NewRequest("GET", fmt.Sprintf("https://eu-api-v2.acrcloud.com/api/base-projects/27121/day-stat?start=%s&end=%s", dateRange.Start, dateRange.End), nil)
	if err != nil {
		log.WithError(err).Fatal("Could not create HTTP Request for ACR Stats!")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Android 4.4; Tablet; rv:41.0) Gecko/41.0 Firefox/41.0")
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Could not execute HTTP Request for ACR Stats!")
	}
	defer r.Body.Close()
	io.Copy(w, resp.Body)
}

func getLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var items ItemResponse
	params := mux.Vars(r)
	id := params["Id"]
	db.Table("items").Select("items.*, comments.*").Where("items.acr_id = ?", id).Joins("LEFT JOIN comments ON items.item_id = comments.item_id").First(&items)
	json.NewEncoder(w).Encode(items)
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
			log.WithError(err).Fatal("Could not compile Regex for Spotify!")
		}

		deezerRegex, err := regexp.Compile(`(?mi)^https?:\/\/(?:www\.)?deezer\.com\/(track|album|playlist)\/(\d+)$`)
		if err != nil {
			log.WithError(err).Fatal("Could not compile Regex for Deezer!")
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
		log.WithError(err).Fatal("Could not retrieve Links from aha-music!")
	}
	response := AHAAPIResponse{}
	err = json.Unmarshal([]byte(preTag), &response)
	if err != nil {
		log.WithError(err).Fatal("Could not unmarshal JSON from aha-music!")
	}
	return response
}

func paginate(value interface{}, pagination *pkg.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Round(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func goThroughBotComments() {
	session := pr0gramm.NewSession(http.Client{Timeout: 10 * time.Second})
	if resp, err := session.Login(os.Getenv("PR0_USER"), os.Getenv("PR0_PASSWORD")); err != nil {
		log.WithError(err).Fatal("Could not login into pr0gramm!")
		return
	} else {
		if !resp.Success {
			log.WithError(err).Fatal("Could not login into pr0gramm!")
			return
		}
	}
	up := Updater{Session: session, After: pr0gramm.Timestamp{time.Unix(1623837600, 0)}}
	go up.Update()
}

func (u *Updater) Update() {
	for {
		data, err := u.Session.GetUserComments("Sauce", 15, int(u.After.Unix()))
		if err != nil {
			log.WithError(err).Fatal("Could not get user comments!")
			break
		}
		for _, c := range data.Comments {
			if !strings.Contains(c.Content, "Es wurden") {
				continue
			}
			comm := Comments{
				CommentID: int(c.Id),
				Up:        c.Up,
				Down:      c.Down,
				Content:   c.Content,
				Created:   &c.Created.Time,
				ItemID:    int(c.ItemId),
				Thumb:     c.Thumbnail,
			}
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&comm)
			log.Printf("Inserted Comment with ID: %d", c.Id)
		}
		if !data.HasNewer {
			break
		}
		u.After = data.Comments[len(data.Comments)-1].Created
	}
}
