package database

import (
	"errors"
	"fmt"
	"sort"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

var ErrNotFound = errors.New("resource not found")

func (db *DB) GetChirps() ([]Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if dbData.Chirps == nil {
		return []Chirp{}, nil
	}

	// Making the length zero will initialize the slice with its null default values
	// Appending then to the slice will add the elements and keep the previously initialized default values
	// Therefore, make the slice's length 0
	chirps := make([]Chirp, 0, len(dbData.Chirps))

	for _, val := range dbData.Chirps {
		chirps = append(chirps, val)
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})
	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return Chirp{}, err
	}
	chirp, ok := dbData.Chirps[id]
	if !ok {
		return Chirp{}, nil
	}
	return chirp, nil
}
func (db *DB) DeleteChirp(id int) error {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, ok := dbData.Chirps[id]
	if !ok {
		return ErrNotFound
	}
	delete(dbData.Chirps, id)
	err = db.writeDB(dbData)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return Chirp{}, err
	}

	newId := 0
	for key := range dbData.Chirps {
		if key > newId {
			newId = key
		}
	}

	newId++

	newChirp := Chirp{
		ID:       newId,
		Body:     body,
		AuthorID: authorID,
	}

	// if storedChirps.Chirps == nil {
	// 	storedChirps.Chirps = make(map[int]Chirp)
	// }

	dbData.Chirps[newId] = newChirp
	err = db.writeDB(dbData)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}
