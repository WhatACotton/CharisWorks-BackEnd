package funcs

import (
	"unify/internal/database"
)

func StoreSession(sessionId string) {
	database.Storesession(sessionId)
}

func Getsession(sessionId string) (date string) {
	return database.Getsession(sessionId)
}
