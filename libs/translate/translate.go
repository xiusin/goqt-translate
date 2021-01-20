package translate

import "math/rand"

type Intf interface {
	Translate(from, to, q string) string
}

func randSalt(length int, isNumber ...bool) string {
	var b = []byte("")
	salt := "abcdefghijklmnopqrstuvwxyz0123456789"
	if len(isNumber) > 0 && isNumber[0] {
		salt = "0123456789"
	}
	l := len(salt)
	for i := 0; i < length; i++ {
		b = append(b, salt[rand.Intn(l)])
	}
	return string(b)
}
