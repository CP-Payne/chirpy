package database

import (
	"errors"
	"fmt"
	"time"
)

type Status struct {
	CreatedAt time.Time
	RevokedAt time.Time
	Revoked   bool
}

var ErrAlreadyRevoked = errors.New("token already revoked")
var ErrDoesNotExist = errors.New("refresh token not found")

func (db *DB) AddTokenToDB(refreshToken string) error {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return err
	}

	tokenInfo := Status{
		CreatedAt: time.Now(),
		RevokedAt: time.Now(),
		Revoked:   false,
	}

	dbData.RefreshTokens[refreshToken] = tokenInfo

	err = db.writeDB(dbData)
	if err != nil {
		fmt.Printf("Failed to write to database: %v", err)
		return err

	}
	return nil

}

func (db *DB) RevokeRefreshToken(refreshToken string) error {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	val, exist := dbData.RefreshTokens[refreshToken]
	if !exist {
		return ErrDoesNotExist
	}

	val.Revoked = true
	val.RevokedAt = time.Now()

	dbData.RefreshTokens[refreshToken] = val
	err = db.writeDB(dbData)
	if err != nil {
		fmt.Printf("Failed to write to database: %v", err)
		return err

	}
	return nil

}

func (db *DB) IsTokenRevoked(refreshToken string) (bool, error) {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return true, err
	}
	if val, exist := dbData.RefreshTokens[refreshToken]; exist {
		if val.Revoked {
			return true, nil
		}
		return false, nil
	}

	return false, ErrDoesNotExist
}
