package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
}

// generates random integer btw min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generates random string
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "NGN"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomUUID() uuid.UUID {
	uuid := uuid.New()

	return uuid
}

func RandomUUIDR() (id *uuid.UUID, err error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &uuid, nil
}

func Test() uuid.UUID {
	id, _ := RandomUUIDR()

	return *id
}

func RandomRole() string {
	roles := []string{"seller", "admin", "buyer"}
	n := len(roles)
	return roles[rand.Intn(n)]
}
