package main

type ErrReturn struct {
	ErrCode int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrAlreadyVoted = &ErrReturn{22, "Already Voted!"}
)
