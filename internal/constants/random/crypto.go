package random

import (
	"crypto/rand"

	"github.com/jeffrpowell/listaway/internal/constants"
)

// Returns random [n]byte, limited to the provided set of byte values provided
// If an error is returned from the random read, this function returns it.
func Bytes(n int, charset []byte) (data []byte, err error) {
	if n < 1 {
		n = constants.DefaultN
	}

	data = make([]byte, n)

	if _, err = rand.Read(data); err != nil {
		return nil, err
	}

	t := len(charset)

	if t > 0 {
		for i := 0; i < n; i++ {
			data[i] = charset[data[i]%byte(t)]
		}
	}

	return data, nil
}

// Returns random string of length n, limited to the provided set of character values provided
// If an error is returned from the underlying random byte read, this function returns it.
func String(n int, characters string) (data string, err error) {
	var d []byte

	if d, err = Bytes(n, []byte(characters)); err != nil {
		return "", err
	}

	return string(d), nil
}
