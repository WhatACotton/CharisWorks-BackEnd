package funcs

import (
	"time"
	"unify/internal/database"
)

func StoreSession(sessionId string) {
	database.Storesession(sessionId)
}

func Getsession(sessionId string) (date time.Time) {
	return database.Getsession(sessionId)
}
