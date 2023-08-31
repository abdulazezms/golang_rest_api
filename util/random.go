package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(int(RandomInt(2, 10))))
}

func RandomFullName() string {
	return fmt.Sprintf("%s %s", RandomOwner(), RandomOwner())
}

// RandomString generates a random string of length n consisting of the alphabet characters in lowercase
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		randomChar := alphabet[rand.Intn(k)]
		sb.WriteByte(randomChar)
	}
	return sb.String()
}

// RandomOwner generates a random owner
func RandomOwner() string {
	return RandomString(6)
}

// RandomAmount generates a random amount of money
func RandomAmount() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency
func RandomCurrency() string {
	currencies := []string{SAR, USD, EUR}
	return currencies[rand.Intn(len(currencies))]
}
