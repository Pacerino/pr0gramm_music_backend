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
