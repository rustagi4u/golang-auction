package repository

import (
	"fmt"
	"model"
	database "shared"
)

// "database/sql"
// create a bid
func CreateBid(bid model.Bid) model.Bid {
	DB, err := database.NewOpen()

	var lastInsertID uint64
	sqlStatement := `INSERT INTO bid (bid_price, offer_id, created_by) VALUES ($1, $2, $3) 
RETURNING bid_id`

	err = DB.QueryRow(sqlStatement, bid.Bid_Price, bid.Offer_Id, bid.Created_By).Scan(&lastInsertID)
	if err != nil {
		panic(err)
	}

	// DB.Close()

	return bid
}

func UpdateBid(bidId int) (bool, string) {
	DB, err := database.NewOpen()
	var bid model.Bid
	// var offer model.Offer
	row := DB.QueryRow("SELECT offer_id FROM bid where bid_id=$1", bidId)
	err = row.Scan(&bid.Offer_Id)

	if err != nil {
		value := err.Error()
		if value == "sql: no rows in result set" {
			return false, "Bid id doesn't exist"
		}

	}

	stmt, err := DB.Prepare("update offer set sold = True where id = $1")
	if err != nil {
		return false, "Connection Error"
	}

	res, err := stmt.Exec(bid.Offer_Id)
	if err != nil {
		return false, "Connection Error"
	}

	affect, err := res.RowsAffected()
	fmt.Println(affect)
	if err != nil {
		return false, "Connection Error"
	}
	stmt1, err := DB.Prepare("update bid set bid_accepted = True where bid_id = $1")
	if err != nil {
		return false, "Connection Error"
	}

	res1, err := stmt1.Exec(bidId)
	if err != nil {
		return false, "Connection Error"
	}

	affect1, err := res1.RowsAffected()
	fmt.Println(affect1)
	if err != nil {
		return false, "Connection Error"
	}

	return true, "Bid Updated Successfully"
}
