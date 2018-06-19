package store

import (
	"github.com/GoodCodingFriends/gpay/entity"
	multierror "github.com/hashicorp/go-multierror"
)

type Store struct {
	Eupho   EuphoStore
	Hamachi HamachiStore
}

func (s *Store) Close() error {
	var result error
	err := s.Eupho.Close()
	if err != nil {
		result = multierror.Append(result, err)
	}
	err = s.Hamachi.Close()
	if err != nil {
		result = multierror.Append(result, err)
	}
	return result
}

type EuphoStore interface {
	Get() (*entity.EuphoImage, error)
	Close() error
}

type HamachiStore interface {
	Get() (*entity.HamachiImage, error)
	Close() error
}
