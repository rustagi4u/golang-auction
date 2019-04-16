package repository

import (
	// "database/sql"
	"fmt"
	"time"

	// "time"
	// "github.com/lib/pq"
	"model"
	database "shared"
)

// Post inserts all the specific record into the database
func CreateOffer(offer model.Offer) model.Offer {
	DB, err := database.NewOpen()
	var lastInsertID uint64
	sqlStatement := `INSERT INTO offer (bid_price, go_live, lifetime, photo_url, title, created_by, sold) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`
	err = DB.QueryRow(sqlStatement, offer.Bid_Price, offer.Go_Live, offer.Lifetime, offer.Photo_Url, offer.Title, offer.Created_By, offer.Sold).Scan(&lastInsertID)
	if err != nil {
		panic(err)
	}

	// DB.Close()
	return offer
}

func GetOfferPrice(offer_id uint64) (model.Offer, error) {
	DB, err := database.NewOpen()
	offer := model.Offer{}
	row := DB.QueryRow("SELECT bid_price FROM offer where id=$1 limit $2", offer_id, 1)
	err = row.Scan(&offer.Bid_Price)
	if err != nil {
		return offer, err
	}
	// DB.Close()
	return offer, nil
}

func UpdateOffer(bid model.Bid) bool {
	DB, err := database.NewOpen()

	stmt, err := DB.Prepare("update offer set bid_price = $1 where id = $2")
	if err != nil {
		fmt.Println(err)
		return false
	}

	res, err := stmt.Exec(bid.Bid_Price, bid.Offer_Id)
	if err != nil {
		fmt.Println(err)
		return false
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(affect, "rows changed")

	// DB.Close()

	return true
}

// Get offer using pagination...
func GetOffer(page int, size int, sortKey string) ([]model.Offer, error) {
	DB, err := database.NewOpen()

	if size == 0 {
		size = 10
	}
	var offsetVal int
	if page > 0 {
		offsetVal = (page * size)
	} else {
		offsetVal = 0
	}

	if sortKey == "" {
		sortKey = "golive"
	}
	// size = 1
	// offsetVal = 1
	rows, err := DB.Query("SELECT id, bid_price, go_live, lifetime, photo_url, title, created_by, sold FROM offer order by $1 limit $2 offset $3", sortKey, size, offsetVal)
	if err != nil {
		return nil, err
	}
	// defer rows.Close()
	offers := make([]model.Offer, 0)
	for rows.Next() {
		var id uint64
		var bidprice float64
		var golive time.Time
		var lifetime int64
		var photourl string
		var title string
		var createby uint64
		var sold bool
		offer := model.Offer{}
		err := rows.Scan(&id, &bidprice, &golive, &lifetime, &photourl, &title, &createby, &sold)
		if err != nil {
			return nil, err
		}
		offer.Id = id
		offer.Bid_Price = bidprice
		offer.Go_Live = golive
		offer.Lifetime = lifetime
		offer.Photo_Url = photourl
		offer.Title = title
		offer.Created_By = createby
		offer.Sold = sold

		offers = append(offers, offer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// DB.Close()
	return offers, nil

}

// sold offers list...
func GetSoldOffers(sold bool) ([]model.Offer, error) {
	DB, err := database.NewOpen()
	rows, err := DB.Query("SELECT id, bid_price, go_live, lifetime, photo_url, title, created_by, sold FROM offer where sold = $1", sold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	offers := make([]model.Offer, 0)
	for rows.Next() {
		var id uint64
		var bidprice float64
		var golive time.Time
		var lifetime int64
		var photourl string
		var title string
		var createby uint64
		var sold bool
		offer := model.Offer{}
		err := rows.Scan(&id, &bidprice, &golive, &lifetime, &photourl, &title, &createby, &sold)
		if err != nil {
			return nil, err
		}
		offer.Id = id
		offer.Bid_Price = bidprice
		offer.Go_Live = golive
		offer.Lifetime = lifetime
		offer.Photo_Url = photourl
		offer.Title = title
		offer.Created_By = createby
		offer.Sold = sold

		offers = append(offers, offer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// DB.Close()
	return offers, nil

}
