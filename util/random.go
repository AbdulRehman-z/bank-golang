package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GenerateRandomInt generates a random integer between min and max (inclusive).
func GenerateRandomInt(min, max int32) int32 {
	return min + random.Int31n(max-min+1)
}

// GenerateRandomString generates a random string of given length using the alphabet.
func GenerateRandomString(length int) string {
	var result strings.Builder
	for i := 0; i < length; i++ {
		randomChar := alphabet[random.Intn(len(alphabet))]
		result.WriteByte(randomChar)
	}
	return result.String()
}

// GenerateRandomOwnerName generates a random owner name consisting of 6 characters.
func GenerateRandomOwnerName() string {
	return GenerateRandomString(6)
}

// GenerateRandomMoney generates a random amount of money between 0 and 1000.
func GenerateRandomMoney() int32 {
	return GenerateRandomInt(0, 1000)
}

// GenerateRandomCurrencyCode generates a random currency code (USD, EUR, CAD).
func GenerateRandomCurrencyCode() string {
	currencyCodes := []string{"USD", "EUR", "CAD"}
	randomIndex := random.Intn(len(currencyCodes))
	return currencyCodes[randomIndex]
}

// GenerateRandomEmail generates a random email address with the format "randomstring@email.com".
func GenerateRandomEmail() string {
	return fmt.Sprintf("%s@email.com", GenerateRandomString(6))
}
