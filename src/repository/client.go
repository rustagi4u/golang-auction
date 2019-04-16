package repository

import (
	// "database/sql"

	// "time"
	// "github.com/lib/pq"
	"model"
	database "shared"
)

// GetAll gets all the records from the database...
func GetUserDetail(email string) (model.Account, error) {
	DB, err := database.NewOpen()
	account := model.Account{}
	row := DB.QueryRow("SELECT email FROM account where email=$1", email)
	err = row.Scan(&account.Email)
	if err != nil {
		return account, err
	}
	// DB.Close()

	return account, nil
}

// Create Account for a new user
func CreateAccount(account model.Account) model.Account {
	DB, err := database.NewOpen()

	sqlStatement := `INSERT INTO account (Email, Password) VALUES ($1, $2) returning A_ID`

	err = DB.QueryRow(sqlStatement, account.Email, account.Password).Scan(&account.A_ID)
	if err != nil {
		panic(err)
	}
	// DB.Close()
	return account
}

// GetAll gets all the records from the database...
func GetUser(email string) (model.Account, error) {
	DB, err := database.NewOpen()
	account := model.Account{}
	row := DB.QueryRow("SELECT a_id, email, password FROM account where email=$1", email)
	err = row.Scan(&account.A_ID, &account.Email, &account.Password)
	if err != nil {
		return account, err
	}
	// DB.Close()
	return account, nil
}
