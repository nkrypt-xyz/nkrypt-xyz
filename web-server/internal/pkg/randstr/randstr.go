package randstr

import (
	"crypto/rand"
	"math/big"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateID returns a securely generated random string of the given length
// using the allowed charset (0-9a-zA-Z). This is used for entity IDs.
func GenerateID(length int) (string, error) {
	return generateFromCharset(length, charset)
}

// GenerateAPIKey returns a securely generated random API key string of the
// given length using the same charset as IDs.
func GenerateAPIKey(length int) (string, error) {
	return generateFromCharset(length, charset)
}

func generateFromCharset(length int, chars string) (string, error) {
	result := make([]byte, length)
	max := big.NewInt(int64(len(chars)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = chars[n.Int64()]
	}

	return string(result), nil
}

