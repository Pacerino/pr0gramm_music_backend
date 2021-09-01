package pr0gramm

import "time"

type Id uint64

type LoginResponse struct {
	Success    bool   `json:"success"`
	Identifier string `json:"identifier"`

	Ban *BanInfo `json:"ban,omitempty"`
}

type BanInfo struct {
	Banned  bool      `json:"banned"`
	Reason  string    `json:"reason"`
	EndTime Timestamp `json:"till"`
}

type Response struct {
	Timestamp    Timestamp     `json:"ts"`
	ResponseTime time.Duration `json:"rt"`
	QueryCount   uint          `json:"qt"`
}

type CommentResponse struct {
	Comments []Comment `json:"comments"`
	HasOlder bool      `json:"hasOlder"`
	HasNewer bool      `json:"hasNewer"`
	Response
}

type Comment struct {
	Id        Id        `json:"id"`
	Created   Timestamp `json:"created"`
	Up        int       `json:"up"`
	Down      int       `json:"down"`
	Content   string    `json:"content"`
	Thumbnail string    `json:"thumb"`
	ItemId    Id        `json:"itemId"`
}
