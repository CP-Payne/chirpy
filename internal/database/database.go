package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
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

func (db *DB) GetChirps() ([]Chirp, error) {
	storedChirps, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if storedChirps.Chirps == nil {
		return []Chirp{}, nil
	}

	// Making the length zero will initialize the slice with its null default values
	// Appending then to the slice will add the elements and keep the previously initialized default values
	// Therefore, make the slice's length 0
	chirps := make([]Chirp, 0, len(storedChirps.Chirps))

	for _, val := range storedChirps.Chirps {
		chirps = append(chirps, val)
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})
	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	storedChirps, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return Chirp{}, err
	}
	chirp, ok := storedChirps.Chirps[id]
	if !ok {
		return Chirp{}, nil
	}
	return chirp, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	storedChirps, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return Chirp{}, err
	}

	newId := 0
	for key := range storedChirps.Chirps {
		if key > newId {
			newId = key
		}
	}

	newId++

	newChirp := Chirp{
		ID:   newId,
		Body: body,
	}

	// if storedChirps.Chirps == nil {
	// 	storedChirps.Chirps = make(map[int]Chirp)
	// }

	storedChirps.Chirps[newId] = newChirp
	err = db.writeDB(storedChirps)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
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
		return DBStructure{Chirps: make(map[int]Chirp)}, nil
	}

	chirps := DBStructure{}
	err = json.Unmarshal(data, &chirps)
	if err != nil {
		// Maybe handle the error within the calling function instead of this function
		fmt.Printf("Error unmarshalling JSON: %s", err)
		return DBStructure{}, err
	}
	return chirps, nil
}
