package helper

import (
	"math/rand"
	"time"
)

const (
	charset  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLength = 8
)

func generateRandomID() string {
	rand.Seed(time.Now().UnixNano())
	id := make([]byte, idLength)
	for i := range id {
		id[i] = charset[rand.Intn(len(charset))]
	}

	return string(id)
}

func GenerateUserID() string {
	return "USR-" + generateRandomID()
}

func GenerateVenueID() string {
	return "VNE-" + generateRandomID()
}

func GenerateReviewID() string {
	return "RVW-" + generateRandomID()
}

func GenerateImageID() string {
	return "IMG-" + generateRandomID()
}
