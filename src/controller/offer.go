package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"model"
	"net/http"
	"repository"
	utils "shared"
	"strconv"
)

//Create Offer...
func CreateOffer(w http.ResponseWriter, req *http.Request) {
	usr := req.Context().Value("user")
	fmt.Println(usr)
	body, _ := ioutil.ReadAll(req.Body)
	requestBody := []byte(body)
	var offer model.Offer
	err := json.Unmarshal(requestBody, &offer)
	if err != nil {
		msg := fmt.Sprintf("Invalid request payload. Error: %s", err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, msg)
		return
	}
	if response, ok := ValidateOffer(offer); !ok {
		utils.Respond(w, utils.Message(false, response["message"].(string)))
		return
	}
	offer.Created_By = usr.(uint64)
	u := repository.CreateOffer(offer)
	go Writer(&offer)
	utils.RespondWithJSON(w, http.StatusCreated, u)
}

func GetOffer(w http.ResponseWriter, req *http.Request) {
	page, err := strconv.Atoi(req.FormValue("page"))
	if err != nil {
		page = 0
	}
	size, err := strconv.Atoi(req.FormValue("size"))
	if err != nil {
		size = 10
	}
	sortKey := req.FormValue("sort")
	fmt.Println("Query Parameters", page, size, sortKey)
	data, err := repository.GetOffer(page, size, sortKey)
	utils.RespondWithJSON(w, http.StatusOK, data)

}

func SoldOffers(w http.ResponseWriter, req *http.Request) {
	var sold = true
	data, err := repository.GetSoldOffers(sold)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, data)
}

//Validate incoming user details...
func ValidateOffer(offer model.Offer) (map[string]interface{}, bool) {

	if offer.Bid_Price <= 0 {
		return utils.Message(false, "Bid Price is required"), false
	}

	if offer.Title == "" {
		return utils.Message(false, "Title is required"), false
	}
	return utils.Message(false, "Requirement passed"), true
}
