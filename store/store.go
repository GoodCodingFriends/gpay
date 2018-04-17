package store

import "github.com/GoodCodingFriends/gpay/entity"

type Store struct {
	Eupho EuphoStore
}

func (s *Store) Close() error {
	return s.Eupho.Close()
}

type EuphoStore interface {
	Get() (*entity.EuphoImage, error)
	Close() error
}
