package lib

import "math/rand"

const APITokenLength = 32
const tokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewAPIToken() string {
	var apiToken string
		b := make([]byte, APITokenLength)
		for i := range b {
			b[i] = tokenChars[rand.Intn(len(tokenChars))]
		}
		apiToken = string(b)

	return apiToken
}


func ValidAPIToken(token string) bool {
	// TODO: do more checks
	return len(token) == APITokenLength
}