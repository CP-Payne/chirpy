package database

import (
	"errors"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string
}

func (db *DB) CreateUser(email string, passwordHash string) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	newId := 0
	for key := range dbData.Users {
		if key > newId {
			newId = key
		}
	}

	newId++

	newUser := User{
		ID:       newId,
		Email:    email,
		Password: passwordHash,
	}
	dbData.Users[newId] = newUser
	err = db.writeDB(dbData)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (db *DB) ValidateUserExist(email string) (bool, error) {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	exist := false
	for _, val := range dbData.Users {
		if val.Email == email {
			exist = true
			break
		}
	}

	return exist, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}
	for _, user := range dbData.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("user does not exist")

}
