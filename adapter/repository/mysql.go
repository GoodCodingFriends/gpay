package repository

import (
	"context"
	"fmt"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity"
	repo "github.com/GoodCodingFriends/gpay/repository"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	tableUsers    = "users"
	tableInvoices = "invoices"
	tableTxs      = "transactions"
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
	db      *sqlx.DB
	user    *mySQLUserRepository
	invoice *mySQLInvoiceRepository
	tx      *mySQLTxRepository
}

func (b *mySQLTxBeginner) BeginTx(ctx context.Context) (repo.TxCommitter, context.Context, error) {
	tx, err := b.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return &mySQLTxCommitter{
		dbtx:    tx,
		user:    b.user,
		invoice: b.invoice,
		tx:      b.tx,
	}, ctx, nil
}

func (b *mySQLTxBeginner) Close() error {
	return b.db.Close()
}

type mySQLTxCommitter struct {
	dbtx *sqlx.Tx

	user    *mySQLUserRepository
	invoice *mySQLInvoiceRepository
	tx      *mySQLTxRepository
}

func (c *mySQLTxCommitter) Commit() error {
	return c.dbtx.Commit()
}

func (c *mySQLTxCommitter) Rollback() error {
	return c.dbtx.Rollback()
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
	q := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, tableUsers)
	err := r.client.Get(&user, q, string(id))
	if err != nil {
		return nil, err
	}
	return toUserEntity(r.cfg, &user), nil
}

func (r *mySQLUserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	var dbUsers []*user
	q := fmt.Sprintf(`SELECT * FROM %s`, tableUsers)
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
	q := fmt.Sprintf(`INSERT INTO %s(
		id, first_name, last_name, display_name, amount)
		VALUES(?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		id = ?, first_name = ?, last_name = ?, display_name = ?, amount = ?`, tableUsers)
	_, err := r.client.Exec(
		q,
		string(user.ID), user.FirstName, user.LastName, user.DisplayName, int64(user.BalanceAmount()),
		string(user.ID), user.FirstName, user.LastName, user.DisplayName, int64(user.BalanceAmount()),
	)
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
	q := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, tableInvoices)
	err := r.client.Get(&invoice, q, string(id))
	if err != nil {
		return nil, err
	}
	return toInvoiceEntity(&invoice), nil
}

func (r *mySQLInvoiceRepository) FindAll(ctx context.Context) ([]*entity.Invoice, error) {
	var dbInvoices []*invoice
	q := fmt.Sprintf(`SELECT * FROM %s`, tableInvoices)
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
	q := fmt.Sprintf(`INSERT INTO %s(
		id, status, from_id, to_id, amount, message)
		VALUES(?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		id = ?, status = ?, from_id = ?, to_id = ?, amount = ?, message = ?`, tableInvoices)
	_, err := r.client.Exec(
		q,
		string(invoice.ID), int(invoice.Status), string(invoice.FromID), string(invoice.ToID), int64(invoice.Amount), invoice.Message,
		string(invoice.ID), int(invoice.Status), string(invoice.FromID), string(invoice.ToID), int64(invoice.Amount), invoice.Message,
	)
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
	q := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, tableTxs)
	err := r.client.Get(&tx, q, string(id))
	if err != nil {
		return nil, err
	}
	return toTransactionEntity(&tx), nil
}

func (r *mySQLTxRepository) FindAll(ctx context.Context) ([]*entity.Transaction, error) {
	var dbTransactions []*transaction
	q := fmt.Sprintf(`SELECT * FROM %s`, tableTxs)
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
	q := fmt.Sprintf(`INSERT INTO %s(
		id, transaction_type, from_id, to_id, amount, message)
		VALUES(?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		id = ?, transaction_type = ?, from_id = ?, to_id = ?, amount = ?, message = ?`, tableTxs)
	_, err := r.client.Exec(
		q,
		string(tx.ID), int(tx.Type), string(tx.From), string(tx.To), int64(tx.Amount), tx.Message,
		string(tx.ID), int(tx.Type), string(tx.From), string(tx.To), int64(tx.Amount), tx.Message,
	)
	return err
}

func newMySQLDB(cfg *config.Config) (*sqlx.DB, error) {
	dsn := &mysql.Config{
		User:   cfg.Repository.MySQL.UserName,
		Passwd: cfg.Repository.MySQL.Password,
		Net:    cfg.Repository.MySQL.Net,
		Addr:   cfg.Repository.MySQL.Address,
		DBName: cfg.Repository.MySQL.DBName,
	}
	return sqlOpen("mysql", dsn.FormatDSN())
}

func NewMySQLRepository(cfg *config.Config) (*repo.Repository, error) {
	db, err := newMySQLDB(cfg)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	user := &mySQLUserRepository{cfg, db}
	invoice := &mySQLInvoiceRepository{cfg, db}
	tx := &mySQLTxRepository{cfg, db}
	return repo.New(
		&mySQLTxBeginner{db, user, invoice, tx},
		user,
		invoice,
		tx,
	), nil
}
