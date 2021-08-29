package main

type AHAAPIResponse struct {
	EntityUniqueID  string `json:"entityUniqueId"`
	UserCountry     string `json:"userCountry"`
	PageURL         string `json:"pageUrl"`
	LinksByPlatform struct {
		Deezer struct {
			Country        string `json:"country"`
			URL            string `json:"url"`
			EntityUniqueID string `json:"entityUniqueId"`
		} `json:"deezer"`
		AppleMusic struct {
			Country             string `json:"country"`
			URL                 string `json:"url"`
			NativeAppURIMobile  string `json:"nativeAppUriMobile"`
			NativeAppURIDesktop string `json:"nativeAppUriDesktop"`
			EntityUniqueID      string `json:"entityUniqueId"`
		} `json:"appleMusic"`
		Itunes struct {
			Country             string `json:"country"`
			URL                 string `json:"url"`
			NativeAppURIMobile  string `json:"nativeAppUriMobile"`
			NativeAppURIDesktop string `json:"nativeAppUriDesktop"`
			EntityUniqueID      string `json:"entityUniqueId"`
		} `json:"itunes"`
		Soundcloud struct {
			Country        string `json:"country"`
			URL            string `json:"url"`
			EntityUniqueID string `json:"entityUniqueId"`
		} `json:"soundcloud"`
		Spotify struct {
			Country             string `json:"country"`
			URL                 string `json:"url"`
			NativeAppURIDesktop string `json:"nativeAppUriDesktop"`
			EntityUniqueID      string `json:"entityUniqueId"`
		} `json:"spotify"`
		Tidal struct {
			Country        string `json:"country"`
			URL            string `json:"url"`
			EntityUniqueID string `json:"entityUniqueId"`
		} `json:"tidal"`
		Yandex struct {
			Country        string `json:"country"`
			URL            string `json:"url"`
			EntityUniqueID string `json:"entityUniqueId"`
		} `json:"yandex"`
		Youtube struct {
			Country        string `json:"country"`
			URL            string `json:"url"`
			EntityUniqueID string `json:"entityUniqueId"`
		} `json:"youtube"`
		YoutubeMusic struct {
			Country        string `json:"country"`
			URL            string `json:"url"`
			EntityUniqueID string `json:"entityUniqueId"`
		} `json:"youtubeMusic"`
	} `json:"linksByPlatform"`
}

type Tabler interface {
	TableName() string
}

func (Items) TableName() string {
	return "Items"
}
func (Metadata) TableName() string {
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

type Metadata struct {
	DeezerURL     string `json:"deezerURL" gorm:"column:deezerUrl"`
	DeezerID      string `json:"deezerID" gorm:"column:deezerID"`
	SoundcloudURL string `json:"soundcloudURL" gorm:"column:soundcloudURL"`
	SoundcloudID  string `json:"soundcloudID" gorm:"column:soundcloudID"`
	SpotifyURL    string `json:"spotifyURL" gorm:"column:spotifyURL"`
	SpotifyID     string `json:"spotifyID" gorm:"column:spotifyID"`
	YoutubeURL    string `json:"youtubeURL" gorm:"column:youtubeURL"`
	YoutubeID     string `json:"youtubeID" gorm:"column:youtubeID"`
	TidalURL      string `json:"tidalURL" gorm:"column:tidalURL"`
	TidalID       string `json:"tidalID" gorm:"column:tidalID"`
	ApplemusicURL string `json:"applemusicURL" gorm:"column:applemusicURL"`
	ApplemusicID  string `json:"applemusicID" gorm:"column:applemusicID"`
	ACRID         string `json:"acrID" gorm:"column:acrID"`
}

type DateRange struct {
	Start string
	End   string
}

type ApiResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    AHAAPIResponse `json:"data"`
}

type CrawlLinks struct {
	Spotify string `json:"spotify"`
	Deezer  string `json:"deezer"`
}
