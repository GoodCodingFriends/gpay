package entity

import uuid "github.com/satori/go.uuid"

func newID() string {
	return uuid.NewV4().String()
}
