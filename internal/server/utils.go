package server

import (
	"math/rand"
	"time"
)

type mainTemplate struct {
	IDInstanceName      string
	APITokenName        string
	PhoneNumberName     string
	MessageBodyName     string
	PhoneNumberFileName string
	FileUrlName         string
}

func generateMessageID() int64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Int63n(1000000)
}
