package handler

import (
	"time"

	"github.com/google/uuid"
)

func GetDate() string {
	const template = "2006-01-02 15:04:05"
	t := time.Now()
	t.Format(template)
	return t.String()
}
func GettransactionId() string {
	uuidObj, _ := uuid.NewUUID()
	return uuidObj.String()
}
