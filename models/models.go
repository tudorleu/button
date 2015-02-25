package models

import (
	"errors"
	"fmt"
	"gopkg.in/gorp.v1"
)

// Models corresponding to the db tables for users and transfers.
// NOTE(tudor): At the moment these also get serialized directly to the API response,
// since the API objects correspond directly to the db schema.

type User struct {
	UserId    int    `db:"id" json:"userId"`
	Email     string `db:"email" json:"email"`
	FirstName string `db:"firstname" json:"firstName"`
	LastName  string `db:"lastname" json:"lastName"`
	Points    int    `db:"points" json:"points"`
}

type Transfer struct {
	TransferId int `db:"id" json:"id"`
	UserId     int `db:"userid" json:"userId"`
	Amount     int `db:"amount" json:"amount"`
}

func NewUser(email, firstName, lastName string) (User, error) {
	user := User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Points:    0,
	}
	err := DbMap.Insert(&user)
	return user, err
}

func GetUser(userId int) (User, error) {
	result, err := DbMap.Get(User{}, userId)
	if err != nil {
		return User{}, err
	}
	if result == nil {
		return User{}, fmt.Errorf("No user with id %d", userId)
	}
	user := result.(*User)
	return *user, nil
}

func NewTransfer(userId, amount int) (Transfer, error) {
	_, err := GetUser(userId)
	if err != nil {
		return Transfer{}, fmt.Errorf("Couldn't retrieve user %d. Error: %s",
			userId, err)
	}
	tx, err := DbMap.Begin()
	if err != nil {
		return Transfer{}, err
	}

	transfer, err := addTransferAndUpdateBalance(tx, userId, amount)
	if err != nil {
		tx.Rollback()
		return Transfer{}, err
	}
	err = tx.Commit()
	if err != nil {
		return Transfer{}, err
	}
	return transfer, nil
}

func addTransferAndUpdateBalance(
	tx *gorp.Transaction, userId, amount int) (Transfer, error) {
	result, err := tx.Exec(
		"UPDATE users SET points = points + $1 WHERE id = $2 AND points + $3 >= 0",
		amount, userId, amount)
	if err != nil {
		return Transfer{}, err
	}
	if rows, err := result.RowsAffected(); rows == 0 || err != nil {
		// NOTE(tudor): We assume the userId passed into this function is valid, so then
		// the query would've failed to update a row if the points balance of the user
		// was insufficient to perform the transfer.
		return Transfer{}, errors.New("Insufficient points")
	}

	transfer := Transfer{Amount: amount, UserId: userId}
	err = tx.Insert(&transfer)
	return transfer, err
}

func GetTransfers(userId int) ([]Transfer, error) {
	var transfers []Transfer
	_, err := DbMap.Select(&transfers, "SELECT * FROM transfers ORDER BY id")
	return transfers, err
}
