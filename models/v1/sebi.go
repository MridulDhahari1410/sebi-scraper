package modelsv1

import "time"

type Report struct {
	Date       string
	Title      string
	Content    []byte
	Department string
}

type GetPublicReportsResponse struct {
	ID         int64     `json:"id"`
	Date       time.Time `json:"date"`
	Title      string    `json:"title"`
	Department string    `json:"department"`
	Status     string    `json:"status"`
}

type RecordsSize struct {
	Size int `json:"size"`
}
