package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	utils "shared"
	"strings"

	// "strconv"
	// "github.com/gorilla/mux"
	"model"
	"repository"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// // Get all Clients
// func GetClients(w http.ResponseWriter, req *http.Request) {
// 	clients, err := repository.GetAllClient()
// 	if err != nil {
// 		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	utils.RespondWithJSON(w, http.StatusOK, clients)
// }

// // "Create Client"
// func CreateClient(w http.ResponseWriter, req *http.Request) {
// 	body, _ := ioutil.ReadAll(req.Body)
// 	requestBody := []byte(body)
// 	var client model.Client
// 	err := json.Unmarshal(requestBody, &client)
// 	if err != nil {
// 		msg := fmt.Sprintf("Invalid request payload. Error: %s", err.Error())
// 		utils.RespondWithError(w, http.StatusBadRequest, msg)
// 		return
// 		// fmt.Println("error:", err)
// 	}

// 	u := repository.CreateClient(client)
// 	utils.RespondWithJSON(w, http.StatusCreated, u)

// }

//Validate incoming user details...
func Validate(account model.Account) (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return utils.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return utils.Message(false, "Password Length Should be greated than 6"), false
	}

	//check for errors and duplicate emails
	email := string(account.Email)
	temp, err := repository.GetUserDetail(email)
	if err != nil {
		value := err.Error()
		if value != "sql: no rows in result set" {
			return utils.Message(false, "Connection error. Please retry"), false
		}
	}
	if temp.Email != "" {
		return utils.Message(false, "Email address already in use by another user."), false
	}

	return utils.Message(false, "Requirement passed"), true
}

// Create a new user...
var CreateAccount = func(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	requestBody := []byte(body)
	var account model.Account
	err := json.Unmarshal(requestBody, &account)

	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	if repsonse, ok := Validate(account); !ok {
		utils.Respond(w, utils.Message(false, repsonse["message"].(string)))
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account = repository.CreateAccount(account)
	fmt.Println(account)
	if account.A_ID <= 0 {
		fmt.Println(account.A_ID)
	}
	//Create new JWT token for the newly registered account
	tk := model.Token{UserId: account.A_ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password

	response := utils.Message(true, "Account has been created")
	response["account"] = account
	utils.Respond(w, response)

}

// Authenticate a user...
var Authenticate = func(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	requestBody := []byte(body)
	var account model.Account
	err := json.Unmarshal(requestBody, &account) //decode the request body into struct and failed if any error occur
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	resp := GetUserInfo(account.Email, account.Password)
	utils.Respond(w, resp)
}

func GetUserInfo(email, password string) map[string]interface{} {

	var account model.Account
	account, err := repository.GetUser(email)
	fmt.Println(account)
	if err != nil {
		value := err.Error()
		if value == "sql: no rows in result set" {
			return utils.Message(false, "Email address not found")
		}
		return utils.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	fmt.Println(err)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return utils.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := model.Token{UserId: account.A_ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	resp := utils.Message(true, "Logged In")
	resp["account"] = account
	return resp
}
