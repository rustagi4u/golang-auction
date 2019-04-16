package controller

import (
	"encoding/json"
	"io/ioutil"
	"model"
	"net/http"
	"repository"
	utils "shared"
	"strconv"
)

// Authenticate a user...
var CreateBid = func(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	requestBody := []byte(body)
	var bid model.Bid
	err := json.Unmarshal(requestBody, &bid) //decode the request body into struct and failed if any error occur
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}
	offerDetails, err := repository.GetOfferPrice(bid.Offer_Id)
	if err != nil {
		value := err.Error()
		if value == "sql: no rows in result set" {
			utils.Respond(w, utils.Message(false, "Offer:Invalid"))
			return
		}
	}
	if bid.Bid_Price < offerDetails.Bid_Price {
		utils.Respond(w, utils.Message(false, "Bid Price is lesser than offer price"))
		return
	}
	usr := req.Context().Value("user")
	bid.Created_By = usr.(uint64)
	bid_data := repository.CreateBid(bid)
	u := repository.UpdateOffer(bid)
	if u != true {
		utils.Respond(w, utils.Message(false, "Offer: Not Updated"))
		return
	}
	resp := utils.Message(true, "Bid Created")
	resp["bid"] = bid_data
	utils.Respond(w, resp)
}

// Accept a bid successfully...
func BidAccepted(w http.ResponseWriter, req *http.Request) {
	bidId, err := strconv.Atoi(req.FormValue("bidId"))
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}
	response, data := repository.UpdateBid(bidId)
	if response != true {
		utils.Respond(w, utils.Message(false, data))
		return
	}
	resp := utils.Message(true, data)
	utils.Respond(w, resp)
}
