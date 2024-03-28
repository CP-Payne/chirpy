package database

import "fmt"

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
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
		ID:    newId,
		Email: email,
	}
	dbData.Users[newId] = newUser
	err = db.writeDB(dbData)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}
