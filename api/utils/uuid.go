package utils

import (
	"log"

	"github.com/satori/go.uuid"
)

func NewUUID() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Printf("sth. went wrong: %s", err)
		return "", err
	}
	return uuid.String(), nil
}
