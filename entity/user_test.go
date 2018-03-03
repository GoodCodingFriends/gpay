package entity

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func newUser() *User {
	return &User{
		ID:          UserID(uuid.NewV4().String()),
		FirstName:   "kumiko",
		LastName:    "omae",
		DisplayName: "omae-chan",
		balance:     &balance{},
	}
}

func TestUser_Pay_errors(t *testing.T) {
	t.Run("same user", func(t *testing.T) {
		u := newUser()
		_, err := u.Pay(u, 100, "")
		require.Equal(t, ErrSameUser, err)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		from := newUser()
		to := newUser()
		_, err := from.Pay(to, 100, "")
		require.Equal(t, ErrInsufficientBalance, err)
	})
}

func TestUser_Pay_success(t *testing.T) {
	from, to := newUser(), newUser()
	from.balance.amount = 1000
	tx, err := from.Pay(to, 300, "msg")
	require.NoError(t, err)
	require.Equal(t, TxTypePay, tx.Type)
	require.Equal(t, from.ID, tx.From)
	require.Equal(t, to.ID, tx.To)
	require.Equal(t, Amount(300), tx.Amount)
	require.Equal(t, "msg", tx.Message)
}
