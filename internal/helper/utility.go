package helper

import (
	"crypto/rand"
	"strconv"
)

func RandomNumbers(length int) (int, error) {
	// generate random numbers of length
	const charset = "0123456789"

	buffer := make([]byte, length)

	_, err := rand.Read(buffer)

	if err != nil {
		return 0, err
	}

	numLen := len(charset)

	for i := 0; i < length; i++ {
		buffer[i] = charset[int(buffer[i])%numLen]

	}

	return strconv.Atoi(string(buffer))
}
