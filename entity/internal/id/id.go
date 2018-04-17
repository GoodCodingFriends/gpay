package id

import uuid "github.com/satori/go.uuid"

func New() string {
	return uuid.Must(uuid.NewV4()).String()
}
