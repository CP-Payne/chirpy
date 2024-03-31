package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp     `json:"chirps"`
	Users         map[int]User      `json:"users"`
	RefreshTokens map[string]Status `json:"tokens"`
}

func NewDB(path string) (*DB, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, err // Unable to create file
		}
		file.Close()
	} else if err != nil {
		return nil, err // Other error while accessing file
	}
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	return db, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	jsonData, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	err = os.WriteFile(db.path, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	if len(data) == 0 {
		return DBStructure{Chirps: make(map[int]Chirp), Users: make(map[int]User), RefreshTokens: make(map[string]Status)}, nil
	}
	dbData := DBStructure{}
	err = json.Unmarshal(data, &dbData)
	if err != nil {
		// Maybe handle the error within the calling function instead of this function
		fmt.Printf("Error unmarshalling JSON: %s", err)
		return DBStructure{}, err
	}

	if dbData.Chirps == nil {
		dbData.Chirps = make(map[int]Chirp)
	}
	if dbData.Users == nil {
		dbData.Users = make(map[int]User)
	}
	if dbData.RefreshTokens == nil {
		dbData.RefreshTokens = make(map[string]Status)
	}
	return dbData, nil
}
