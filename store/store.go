package store

import "github.com/GoodCodingFriends/gpay/entity"

type Store struct {
	Eupho EuphoStore
}

type EuphoStore interface {
	Get() (*entity.EuphoImage, error)
}
