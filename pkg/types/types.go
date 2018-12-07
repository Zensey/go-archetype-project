package types

import "time"

type Review struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	ProductID string `json:"productid,omitempty"`
	Review    string `json:"review,omitempty"`
}

type ReviewResponse struct {
	Success  bool  `json:"success"`
	ReviewID int64 `json:"reviewid,omitempty"`
}

type MsgReview struct {
	ReviewID int64  `json:"reviewid"`
	Review   string `json:"review,omitempty"`
}

type RecReview struct {
	Name      string `db:"reviewername"`
	Email     string `db:"emailaddress"`
	ProductID string `db:"productid"`
	Status    bool   `db:"approved"`
}

type RowReview struct {
	Name      string    `db:"reviewername"`
	Email     string    `db:"emailaddress"`
	ProductID string    `db:"productid"`
	Status    *bool     `db:"approved"`
	Created   time.Time `db:"reviewdate"`
}

type RowCount struct {
	Count  *int  `db:"count"`
	Status *bool `db:"approved"`
}

type ReportTable struct {
	Title string
	Items []RowReview
	Count []RowCount
}
