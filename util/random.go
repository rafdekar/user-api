package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const alphabetAndNumbers = "1234567890abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates random integer between min and max
func RandomInt(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomWord generates random word of length n
func RandomWord(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		letter := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(letter)
	}

	return sb.String()
}

func RandomWordWithNumbers(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		letter := alphabetAndNumbers[rand.Intn(len(alphabetAndNumbers))]
		sb.WriteByte(letter)
	}

	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s", RandomWord(10), RandomProvider())
}

func RandomProvider() string {
	domains := []string{
		"gmail.com",
		"google.com",
		"microsoft.com",
		"youtube.com",
		"facebook.com",
	}
	return domains[rand.Intn(len(domains))]
}

func RandomCountry() string {
	countries := []string{
		"UK",
		"DE",
		"NL",
		"PL",
		"ET",
		"LT",
		"LV",
		"EE",
	}
	return countries[rand.Intn(len(countries))]
}
