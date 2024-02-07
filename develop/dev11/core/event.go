package core

import "time"

type Event struct {
	ID     string    `json:"-"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
	UserID string    `json:"user_id"`
}

type EventResp struct {
	ID   string    `json:"-"`
	Text string    `json:"text"`
	Date time.Time `json:"date"`
}
