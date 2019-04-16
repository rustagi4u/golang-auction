package model

// Offer model
type Bid struct {
	Bid_Id       uint64
	Bid_Price    float64 `json:"bid_price" bson:"bid_price"`
	Offer_Id     uint64  `json:"offer_id" bson:"offer_id"`
	Created_By   uint64  `json:"created_by" bson:"created_by"`
	Bid_Accepted bool    `json:"bid_accepted" bson:"bid_accepted"`
}

// CREATE TABLE bid (
// 	Bid_Id BIGSERIAL PRIMARY KEY NOT NULL,
// 	Bid_Price float,
// 	Offer_Id int,
// 	Created_By int,
// 	Bid_Accepted boolean Default false
// );
