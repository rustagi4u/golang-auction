package model

// "User model"
type Client struct {
	CID          uint64
	ClientName   string `json:"username" bson:"username"`
	Password     string `json:”password” bson:”password"`
	EmailAddress string `json:”emailaddress” bson:”emailaddress”`
}

// CREATE TABLE Client (
// 	CID BIGSERIAL PRIMARY KEY NOT NULL,
// 	ClientName TEXT NOT NULL,
// 	Password TEXT NOT NULL,
// 	email_address VARCHAR(45) NOT NULL
// );

//https://github.com/AmundsenJunior/rest-go-mux-pq

// "Account model"
type Account struct {
	A_ID     uint64
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

// CREATE TABLE account (
// 	A_ID BIGSERIAL PRIMARY KEY NOT NULL,
// Email VARCHAR(45) NOT NULL,
// Password TEXT NOT NULL,
// Token VARCHAR(45)
// );