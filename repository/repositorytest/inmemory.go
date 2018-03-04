package repositorytest

import (
	"context"
	"sync"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository"
)

type inMemoryUserRepository struct {
	m sync.Map
	repository.UserRepository
}

func (r *inMemoryUserRepository) StoreAll(ctx context.Context, users []*entity.User) error {
	for _, user := range users {
		r.m.Store(user.ID, user)
	}
	return nil
}

type inMemoryTxRepository struct {
	m sync.Map
	repository.TxRepository
}

func (r *inMemoryTxRepository) Store(ctx context.Context, tx *entity.Transaction) error {
	r.m.Store(tx.ID, tx)
	return nil
}

type inMemoryTxBeginner struct {
}

func (b *inMemoryTxBeginner) BeginTx(ctx context.Context) (repository.TxCommitter, error) {
	return &inMemoryTxCommitter{}, nil
}

type inMemoryTxCommitter struct {
}

func (c *inMemoryTxCommitter) Commit() error {
	return nil
}

func (c *inMemoryTxCommitter) Rollback() error {
	return nil
}

func NewInMemory() *repository.Repository {
	return repository.New(
		&inMemoryTxBeginner{},
		&inMemoryUserRepository{},
		nil,
		&inMemoryTxRepository{},
	)
}
