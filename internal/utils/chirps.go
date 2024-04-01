package utils

import (
	"sort"

	"github.com/CP-Payne/chirpy/internal/database"
)

func SortChirpSlice(chirps []database.Chirp, order string) []database.Chirp {

	sort.Slice(chirps, func(i, j int) bool {
		if order == "desc" {
			return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
		}
		// If "desc" is not passed then it will default to asc
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})
	return chirps

}
