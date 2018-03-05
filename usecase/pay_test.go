package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/entity/entitytest"
	"github.com/GoodCodingFriends/gpay/repository"
	"github.com/GoodCodingFriends/gpay/repository/repositorytest"
	"github.com/stretchr/testify/require"
)

type txRepositoryWithError struct {
	repository.TxRepository
}

func (r *txRepositoryWithError) Store(_ context.Context, _ *entity.Transaction) error {
	return errors.New("an error")
}

func TestPay(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := repositorytest.NewInMemory()

		from, to := createUser(t, repo), createUser(t, repo)
		amount := entity.Amount(300)
		param := &PayParam{
			FromID:  from.ID,
			ToID:    to.ID,
			Amount:  amount,
			Message: "msg",
		}

		tx, err := Pay(repo, param)
		require.NoError(t, err)

		from, to = getUser(t, repo, from.ID), getUser(t, repo, to.ID)

		require.Equal(t, entity.Amount(-300), entity.GetBalanceAmount(from))
		require.Equal(t, entity.Amount(300), entity.GetBalanceAmount(to))

		assertTx(t, tx, from, to, entity.TxTypePay, amount, "msg")
	})

	t.Run("rollback", func(t *testing.T) {
		repo := repositorytest.NewInMemory()
		repo.Transaction = &txRepositoryWithError{}

		from, to := createUser(t, repo), createUser(t, repo)
		amount := entity.Amount(300)
		param := &PayParam{
			FromID:  from.ID,
			ToID:    to.ID,
			Amount:  amount,
			Message: "msg",
		}

		_, err := Pay(repo, param)
		require.Error(t, err)

		from, to = getUser(t, repo, from.ID), getUser(t, repo, to.ID)

		require.Equal(t, entity.Amount(0), entity.GetBalanceAmount(from))
		require.Equal(t, entity.Amount(0), entity.GetBalanceAmount(to))
	})
}

func createUser(t *testing.T, repo *repository.Repository) *entity.User {
	u := entitytest.NewUser(t)
	repo.User.Store(context.TODO(), u)
	return u
}

func getUser(t *testing.T, repo *repository.Repository, id entity.UserID) *entity.User {
	u, err := repo.User.FindByID(context.TODO(), id)
	require.NoError(t, err)
	return u
}

func assertTx(t *testing.T, tx *entity.Transaction, from, to *entity.User, txType entity.TxType, amount entity.Amount, message string) {
	t.Helper()
	require.Equal(t, txType, tx.Type)
	require.Equal(t, from.ID, tx.From)
	require.Equal(t, to.ID, tx.To)
	require.Equal(t, amount, tx.Amount)
	require.Equal(t, message, tx.Message)
}
