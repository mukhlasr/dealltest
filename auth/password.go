package auth

import (
	"crypto/rand"
	"fmt"
	"hash"
	"io"
	"math/big"
	"strings"
)

// GeneratePassword will generate password with len(result) = length
// It will use the randReader as the source of the randomness.
// The password characters will be between 0x30 - 0x7A in the ascii table
func GeneratePassword(length int, randReader io.Reader) (string, error) {
	var res []byte
	// get ascii randomly from 0x30 '0' - 0x7A 'z' in ascii table
	// by getting random number from the random generator at certain range
	// 0x7A - 0x30 = 0x4a = dec(74)
	for i := 0; i < length; i++ {
		n, err := rand.Int(randReader, big.NewInt(75)) // [0-75)
		if err != nil {
			return "", err
		}
		res = append(res, byte(n.Int64())+'0') // shift the num, 0 = ascii '0', 1 = ascii '1', etc...
	}
	return string(res), nil
}

func PasswordMatched(password, hashString string, hash hash.Hash) bool {
	return GetPasswordHash(password, hash) == hashString
}

func GetPasswordHash(password string, hash hash.Hash) string {
	hash.Reset()

	_, _ = io.Copy(hash, strings.NewReader(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
