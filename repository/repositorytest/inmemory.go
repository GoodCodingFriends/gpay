package repositorytest

import (
	"context"
	"sync"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository"
)

var (
	txFlagKey = "in_transaction"
)

type inMemoryUserRepository struct {
	m           sync.Map
	uncommitted sync.Map
	repository.UserRepository
}

func (r *inMemoryUserRepository) FindByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	iv, ok := r.m.Load(id)
	if !ok {
		return nil, repository.ErrNotFound
	}
	v := iv.(entity.User)
	return &v, nil
}

func (r *inMemoryUserRepository) Store(ctx context.Context, user *entity.User) error {
	store(ctx, &r.m, &r.uncommitted, user.ID, *user)
	return nil
}

func (r *inMemoryUserRepository) StoreAll(ctx context.Context, users []*entity.User) error {
	for _, user := range users {
		if err := r.Store(ctx, user); err != nil {
			return err
		}
	}
	return nil
}

type inMemoryInvoiceRepository struct {
	m           sync.Map
	uncommitted sync.Map
	repository.InvoiceRepository
}

func (r *inMemoryInvoiceRepository) FindByID(ctx context.Context, id entity.InvoiceID) (*entity.Invoice, error) {
	iv, ok := r.m.Load(id)
	if !ok {
		return nil, repository.ErrNotFound
	}
	v := iv.(entity.Invoice)
	return &v, nil
}

func (r *inMemoryInvoiceRepository) Store(ctx context.Context, invoice *entity.Invoice) error {
	store(ctx, &r.m, &r.uncommitted, invoice.ID, *invoice)
	return nil
}

type inMemoryTxRepository struct {
	m           sync.Map
	uncommitted sync.Map
	repository.TxRepository
}

func (r *inMemoryTxRepository) Store(ctx context.Context, tx *entity.Transaction) error {
	store(ctx, &r.m, &r.uncommitted, tx.ID, *tx)
	return nil
}

type inMemoryTxBeginner struct {
	user *inMemoryUserRepository
	tx   *inMemoryTxRepository
}

func (b *inMemoryTxBeginner) BeginTx(ctx context.Context) (repository.TxCommitter, context.Context, error) {
	return &inMemoryTxCommitter{b.user, b.tx}, context.WithValue(ctx, txFlagKey, true), nil
}

type inMemoryTxCommitter struct {
	user *inMemoryUserRepository
	tx   *inMemoryTxRepository
}

func (c *inMemoryTxCommitter) Commit() error {
	c.user.uncommitted.Range(func(k, v interface{}) bool {
		c.user.m.Store(k, v)
		return true
	})
	c.tx.uncommitted.Range(func(k, v interface{}) bool {
		c.tx.m.Store(k, v)
		return true
	})
	clearUncommitted(c)
	return nil
}

// Rollback clears uncomitted changes
func (c *inMemoryTxCommitter) Rollback() error {
	clearUncommitted(c)
	return nil
}

func store(ctx context.Context, m *sync.Map, um *sync.Map, k, v interface{}) {
	f := ctx.Value(txFlagKey)
	// in a transaction
	if f != nil && f.(bool) == true {
		um.Store(k, v)
	} else {
		m.Store(k, v)
	}
}

func clearUncommitted(c *inMemoryTxCommitter) {
	c.user.uncommitted = sync.Map{}
	c.tx.uncommitted = sync.Map{}
}

func NewInMemory() *repository.Repository {
	user := &inMemoryUserRepository{}
	invoice := &inMemoryInvoiceRepository{}
	tx := &inMemoryTxRepository{}
	return repository.New(
		&inMemoryTxBeginner{user, tx},
		user,
		invoice,
		tx,
	)
}
