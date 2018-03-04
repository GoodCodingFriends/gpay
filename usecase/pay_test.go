package usecase

import (
	"testing"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/entity/entitytest"
	"github.com/GoodCodingFriends/gpay/repository/repositorytest"
	"github.com/stretchr/testify/require"
)

func TestPay(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := repositorytest.NewInMemory()

		from, to := entitytest.NewUser(t), entitytest.NewUser(t)
		amount := entity.Amount(300)
		param := &PayParam{
			From:    from,
			To:      to,
			Amount:  amount,
			Message: "msg",
		}

		tx, err := Pay(repo, param)
		require.NoError(t, err)

		assertTx(t, tx, from, to, entity.TxTypePay, amount, "msg")
	})
}

func assertTx(t *testing.T, tx *entity.Transaction, from, to *entity.User, txType entity.TxType, amount entity.Amount, message string) {
	t.Helper()
	require.Equal(t, txType, tx.Type)
	require.Equal(t, from.ID, tx.From)
	require.Equal(t, to.ID, tx.To)
	require.Equal(t, amount, tx.Amount)
	require.Equal(t, message, tx.Message)
}
