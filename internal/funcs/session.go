package funcs

import (
	"unify/internal/database"
)

func StoreSession(sessionId string) {
	database.Storesession(sessionId)
}

func Getsession(sessionId string) (date []uint8) {
	return database.Getsession(sessionId)
}
