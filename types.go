package main

import (
	"time"

	"github.com/pacerino/pr0gramm_music_backend/pr0gramm"
	"gorm.io/gorm"
)

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

type Items struct {
	gorm.Model
	ItemID   int      `json:"itemID" gorm:"not null;"`
	Title    string   `json:"title"`
	Album    string   `json:"album"`
	Artist   string   `json:"artist"`
	Url      string   `json:"url"`
	ACRID    string   `json:"acrID"`
	Metadata Metadata `gorm:"embedded"`
}

type Comments struct {
	CommentID int `json:"commentID" gorm:"primarykey;not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Up        int            `json:"up" gorm:"not null;"`
	Down      int            `json:"down" gorm:"not null;"`
	Content   string         `json:"content" gorm:"not null;"`
	Created   *time.Time     `json:"created" gorm:"not null;"`
	ItemID    int            `json:"itemid" gorm:"not null;"`
	Thumb     string         `json:"thumb" gorm:"not null;"`
}

type ItemResponse struct {
	Items
	Comments `json:"comment"`
}

type Metadata struct {
	DeezerURL     string `json:"deezerURL"`
	DeezerID      string `json:"deezerID" `
	SoundcloudURL string `json:"soundcloudURL"`
	SoundcloudID  string `json:"soundcloudID"`
	SpotifyURL    string `json:"spotifyURL"`
	SpotifyID     string `json:"spotifyID"`
	YoutubeURL    string `json:"youtubeURL" `
	YoutubeID     string `json:"youtubeID"`
	TidalURL      string `json:"tidalURL"`
	TidalID       string `json:"tidalID"`
	ApplemusicURL string `json:"applemusicURL"`
	ApplemusicID  string `json:"applemusicID"`
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

type Updater struct {
	Session *pr0gramm.Session
	After   pr0gramm.Timestamp
}
