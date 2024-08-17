package random

import (
	"crypto/rand"
	"io"
)

// Random is the production random.Provider which uses crypto/rand.
type Random struct{}

// Read implements the io.Reader interface.
func (r *Random) Read(p []byte) (n int, err error) {
	return io.ReadFull(rand.Reader, p)
}

// BytesCustomErr returns random data as bytes with n length and can contain only byte values from the provided
// values. If n is less than 1 then DefaultN is used instead. If an error is returned from the random read this function
// returns it.
func (r *Random) BytesCustomErr(n int, charset []byte) (data []byte, err error) {
	if n < 1 {
		n = DefaultN
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

// StringCustomErr is an overload of BytesCustomWithErr which takes a characters string and returns a string.
func (r *Random) StringCustomErr(n int, characters string) (data string, err error) {
	var d []byte

	if d, err = r.BytesCustomErr(n, []byte(characters)); err != nil {
		return "", err
	}

	return string(d), nil
}
