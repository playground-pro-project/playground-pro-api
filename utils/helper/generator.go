package helper

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
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

func GenerateReservationID() string {
	return uuid.New().String()
}

func GenerateOTP(length int) string {
	rand.Seed(time.Now().UnixNano())

	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		otp[i] = byte(rand.Intn(10)) + '0'
	}

	return string(otp)
}

func GenerateIdentifier() string {
	return uuid.New().String()
}
