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

var ErrAlreadyExists = errors.New("already exists")

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

func (db *DB) UpdateUser(id int, email, hashedPassword string) error {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return err
	}

	userEmailExist, err := db.ValidateUserExist(email)
	if err != nil {
		fmt.Println(err)
		return err
	}

	user, ok := dbData.Users[id]
	if !ok {
		return errors.New("could not find user")
	}

	// Check if the email does not exist or that it is not equal to its current email
	// If the provided email and the current email are the same, the user is probably just changing the password
	if userEmailExist && user.Email != email {
		return ErrAlreadyExists
	}
	user.Email = email
	user.Password = hashedPassword

	dbData.Users[id] = user

	err = db.writeDB(dbData)
	if err != nil {
		fmt.Printf("Failed to write to database: %v", err)
		return err

	}
	return nil

}
