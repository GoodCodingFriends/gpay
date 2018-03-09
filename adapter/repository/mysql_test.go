package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/entity/entitytest"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var db *sqlx.DB
var cfg *config.Config
var tables = []string{tableInvoices, tableTxs, tableUsers}

func init() {
	var err error
	cfg, err = config.Process()
	if err != nil {
		panic(err)
	}
	db, err = newMySQLDB(cfg)
	if err != nil {
		panic(err)
	}
}

func cleanup(t *testing.T) {
	_, err := db.DB.Exec("SET FOREIGN_KEY_CHECKS=0")
	require.NoError(t, err)
	for _, table := range tables {
		_, err := db.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
		require.NoError(t, err)
	}
	_, err = db.DB.Exec("SET FOREIGN_KEY_CHECKS=1")
	require.NoError(t, err)
}

func TestMySQLRepository_User(t *testing.T) {
	repo, err := NewMySQLRepository(cfg)
	require.NoError(t, err)

	t.Run("FindByID (not found)", func(t *testing.T) {
		defer cleanup(t)
		_, err := repo.User.FindByID(context.TODO(), entity.UserID("yuko"))
		require.Error(t, err)
	})

	t.Run("FindByID (found)", func(t *testing.T) {
		defer cleanup(t)
		u := entitytest.NewUser(t)
		err := repo.User.Store(context.TODO(), u)
		require.NoError(t, err)
		actual, err := repo.User.FindByID(context.TODO(), u.ID)
		require.NoError(t, err)
		require.Exactly(t, u, actual)
	})

	t.Run("FindAll (no entities)", func(t *testing.T) {
		defer cleanup(t)
		_, err := repo.User.FindAll(context.TODO())
		require.NoError(t, err)
	})

	t.Run("FindAll (found)", func(t *testing.T) {
		defer cleanup(t)
		err := repo.User.StoreAll(context.TODO(), []*entity.User{entitytest.NewUser(t), entitytest.NewUser(t)})
		require.NoError(t, err)
		actual, err := repo.User.FindAll(context.TODO())
		require.NoError(t, err)
		require.Len(t, actual, 2)
	})
}

func TestMySQLRepository_Invoice(t *testing.T) {
	repo, err := NewMySQLRepository(cfg)
	require.NoError(t, err)

	t.Run("FindByID (not found)", func(t *testing.T) {
		defer cleanup(t)
		_, err := repo.Invoice.FindByID(context.TODO(), entity.InvoiceID("yuko"))
		require.Error(t, err)
	})

	t.Run("FindByID / FindAll (found)", func(t *testing.T) {
		defer cleanup(t)
		i := entitytest.NewInvoice(t)
		err := repo.User.StoreAll(context.TODO(), []*entity.User{
			&entity.User{ID: i.FromID},
			&entity.User{ID: i.ToID},
		})
		require.NoError(t, err)

		err = repo.Invoice.Store(context.TODO(), i)
		require.NoError(t, err)
		actual, err := repo.Invoice.FindByID(context.TODO(), i.ID)
		require.NoError(t, err)
		require.Exactly(t, i, actual)

		actuals, err := repo.Invoice.FindAll(context.TODO())
		require.NoError(t, err)
		require.Len(t, actuals, 1)
	})

	t.Run("FindAll (no entities)", func(t *testing.T) {
		defer cleanup(t)
		_, err := repo.Invoice.FindAll(context.TODO())
		require.NoError(t, err)
	})
}

func TestMySQLRepository_Transaction(t *testing.T) {
	repo, err := NewMySQLRepository(cfg)
	require.NoError(t, err)

	t.Run("FindByID (not found)", func(t *testing.T) {
		defer cleanup(t)
		_, err := repo.Transaction.FindByID(context.TODO(), entity.TxID("yuko"))
		require.Error(t, err)
	})

	t.Run("FindByID / FindAll (found)", func(t *testing.T) {
		defer cleanup(t)
		tx := entitytest.NewTransaction(t)
		err := repo.User.StoreAll(context.TODO(), []*entity.User{
			&entity.User{ID: tx.From},
			&entity.User{ID: tx.To},
		})
		require.NoError(t, err)

		err = repo.Transaction.Store(context.TODO(), tx)
		require.NoError(t, err)
		actual, err := repo.Transaction.FindByID(context.TODO(), tx.ID)
		require.NoError(t, err)
		require.Exactly(t, tx, actual)
	})
}
