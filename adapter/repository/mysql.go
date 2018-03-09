package repository

import (
	"context"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity"
	repo "github.com/GoodCodingFriends/gpay/repository"
	"github.com/jmoiron/sqlx"
)

var sqlOpen func(driveName, dataSourceName string) (*sqlx.DB, error) = sqlx.Open

type client interface {
	sqlx.Execer
	sqlx.Queryer

	// sqlx
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

type mySQLTxBeginner struct {
	db   *sqlx.DB
	user *mySQLUserRepository
}

func (b *mySQLTxBeginner) BeginTx(ctx context.Context) (repo.TxCommitter, context.Context, error) {
	tx, err := b.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return &mySQLTxCommitter{
		tx:   tx,
		user: b.user,
	}, ctx, nil
}

type mySQLTxCommitter struct {
	tx   *sqlx.Tx
	user *mySQLUserRepository
}

func (c *mySQLTxCommitter) Commit() error {
	return c.tx.Commit()
}

func (c *mySQLTxCommitter) Rollback() error {
	return c.tx.Rollback()
}

type user struct {
	ID          string
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	DisplayName string `db:"display_name"`
	Amount      int64
}

type mySQLUserRepository struct {
	cfg    *config.Config
	client client
}

func toUserEntity(cfg *config.Config, user *user) *entity.User {
	return entity.NewUser(
		cfg,
		entity.UserID(user.ID),
		user.FirstName,
		user.LastName,
		user.DisplayName,
		user.Amount,
	)
}

func (r *mySQLUserRepository) FindByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	var user user
	q := `SELECT * FROM users WHERE id = ?`
	err := r.client.Get(&user, q, id)
	if err != nil {
		return nil, err
	}
	return toUserEntity(r.cfg, &user), nil
}

func (r *mySQLUserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	var dbUsers []*user
	q := `SELECT * FROM users`
	err := r.client.Select(&dbUsers, q)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, 0, len(dbUsers))
	for _, user := range dbUsers {
		users = append(
			users,
			toUserEntity(r.cfg, user),
		)
	}
	return users, nil
}

func (r *mySQLUserRepository) Store(ctx context.Context, user *entity.User) error {
	q := `INSERT INTO users(
		id, first_name, last_name, display_name, amount)
		VALUES(?, ?, ?, ?, ?)`
	_, err := r.client.Exec(q, user.ID, user.FirstName, user.LastName, user.DisplayName, user.BalanceAmount())
	return err
}

func (r *mySQLUserRepository) StoreAll(ctx context.Context, users []*entity.User) error {
	for _, user := range users {
		if err := r.Store(ctx, user); err != nil {
			return err
		}
	}
	return nil
}

type invoice struct {
	ID      string
	Status  int8
	FromID  string `db:"from_id"`
	ToID    string `db:"to_id"`
	Amount  int64
	Message string
}

func toInvoiceEntity(invoice *invoice) *entity.Invoice {
	return &entity.Invoice{
		ID:      entity.InvoiceID(invoice.ID),
		FromID:  entity.UserID(invoice.FromID),
		ToID:    entity.UserID(invoice.ToID),
		Amount:  entity.Amount(invoice.Amount),
		Message: invoice.Message,
	}
}

type mySQLInvoiceRepository struct {
	cfg    *config.Config
	client client
}

func (r *mySQLInvoiceRepository) FindByID(ctx context.Context, id entity.InvoiceID) (*entity.Invoice, error) {
	var invoice invoice
	q := `SELECT * FROM invoices WHERE id = ?`
	err := r.client.Get(&invoice, q, id)
	if err != nil {
		return nil, err
	}
	return toInvoiceEntity(&invoice), nil
}

func (r *mySQLInvoiceRepository) FindAll(ctx context.Context) ([]*entity.Invoice, error) {
	var dbInvoices []*invoice
	q := `SELECT * FROM invoices`
	err := r.client.Select(&dbInvoices, q)
	if err != nil {
		return nil, err
	}

	invoices := make([]*entity.Invoice, 0, len(dbInvoices))
	for _, invoice := range dbInvoices {
		invoices = append(
			invoices,
			toInvoiceEntity(invoice),
		)
	}
	return invoices, nil
}

func (r *mySQLInvoiceRepository) Store(ctx context.Context, invoice *entity.Invoice) error {
	q := `INSERT INTO invoices(
		id, status, from_id, to_id, amount, message)
		VALUES(?, ?, ?, ?, ?, ?)`
	_, err := r.client.Exec(q, invoice.ID, invoice.Status, invoice.FromID, invoice.ToID, invoice.Amount, invoice.Message)
	return err
}

func (r *mySQLInvoiceRepository) StoreAll(ctx context.Context, invoices []*entity.Invoice) error {
	for _, invoice := range invoices {
		if err := r.Store(ctx, invoice); err != nil {
			return err
		}
	}
	return nil
}

type mySQLTxRepository struct {
	cfg    *config.Config
	client client
}

type transaction struct {
	ID      string
	Type    int8   `db:"transaction_type"`
	FromID  string `db:"from_id"`
	ToID    string `db:"to_id"`
	Amount  int64
	Message string
}

func toTransactionEntity(tx *transaction) *entity.Transaction {
	return &entity.Transaction{
		ID:      entity.TxID(tx.ID),
		Type:    entity.TxType(tx.Type),
		From:    entity.UserID(tx.FromID),
		To:      entity.UserID(tx.ToID),
		Amount:  entity.Amount(tx.Amount),
		Message: tx.Message,
	}
}

func (r *mySQLTxRepository) FindByID(ctx context.Context, id entity.TxID) (*entity.Transaction, error) {
	var tx transaction
	q := `SELECT * FROM transactions WHERE id = ?`
	err := r.client.Get(&tx, q, id)
	if err != nil {
		return nil, err
	}
	return toTransactionEntity(&tx), nil
}

func (r *mySQLTxRepository) FindAll(ctx context.Context) ([]*entity.Transaction, error) {
	var dbTransactions []*transaction
	q := `SELECT * FROM transactions`
	err := r.client.Select(&dbTransactions, q)
	if err != nil {
		return nil, err
	}

	transactions := make([]*entity.Transaction, 0, len(dbTransactions))
	for _, tx := range dbTransactions {
		transactions = append(
			transactions,
			toTransactionEntity(tx),
		)
	}
	return transactions, nil
}

func (r *mySQLTxRepository) Store(ctx context.Context, tx *entity.Transaction) error {
	q := `INSERT INTO transactions(
		id, transaction_type, from_id, to_id, amount, message)
		VALUES(?, ?, ?, ?, ?, ?)`
	_, err := r.client.Exec(q, tx.ID, tx.Type, tx.From, tx.To, tx.Amount, tx.Message)
	return err
}

func NewMySQLRepository(cfg *config.Config) (*repo.Repository, error) {
	db, err := sqlOpen("", "")
	if err != nil {
		return nil, err
	}
	user := &mySQLUserRepository{cfg, db}
	invoice := &mySQLInvoiceRepository{}
	tx := &mySQLTxRepository{}
	return repo.New(
		&mySQLTxBeginner{db, user},
		user,
		invoice,
		tx,
	), nil
}
