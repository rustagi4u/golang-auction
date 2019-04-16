package model

import "time"

// Offer model
type Offer struct {
	Id         uint64    `json:"id" bson:"_id,omitempty"`
	Bid_Price  float64   `json:"bid_price"`
	Go_Live    time.Time `json:"go_live"`
	Lifetime   int64     `json:"lifetime"`
	Photo_Url  string    `json:"photo_url"`
	Title      string    `json:"title"`
	Created_By uint64    `json:"created_by"`
	Sold       bool      `json:"sold"`
}

// CREATE TABLE Offer (
// 	id BIGSERIAL PRIMARY KEY NOT NULL,
// 	bid_price FLOAT NOT NULL Default 0.00,
// 	go_live TIMESTAMP,
// 	lifetime Integer,
// 	photo_rul VARCHAR(200),
// 	title VARCHAR(50) NOT NULL,
// 	created_by int,
// 	sold BOOLEAN Default false
// );
