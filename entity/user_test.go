package entity

import (
	"testing"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity/internal/id"
	"github.com/stretchr/testify/require"
)

func newUser() *User {
	return &User{
		ID:          UserID(id.New()),
		FirstName:   "kumiko",
		LastName:    "omae",
		DisplayName: "omae-chan",
		balance: balance{
			conf: &config.Config{
				Entity: &config.Entity{
					BalanceLowerLimit: -5000,
				},
			},
		},
	}
}

func TestUser_Pay_errors(t *testing.T) {
	t.Run("same user", func(t *testing.T) {
		u := newUser()
		_, err := u.Pay(u, 100, "")
		require.Equal(t, ErrSameUser, err)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		cases := []struct {
			amount int64
			isErr  bool
		}{
			{-5000, true},
			{-4901, true},
			{-4900, false},
		}
		for _, c := range cases {
			from := newUser()
			from.balance.amount = Amount(c.amount)
			to := newUser()
			_, err := from.Pay(to, 100, "")
			if c.isErr {
				require.Equal(t, ErrInsufficientBalance, err)
			} else {
				require.NoError(t, err)
			}
		}
	})

	t.Run("amount is 0", func(t *testing.T) {
		from := newUser()
		to := newUser()
		_, err := from.Pay(to, 0, "")
		require.Equal(t, ErrZeroAmount, err)
	})
}

func TestUser_Pay_success(t *testing.T) {
	from, to := newUser(), newUser()
	from.balance.amount = 1000

	tx, err := from.Pay(to, 300, "msg")
	require.NoError(t, err)

	assertTx(t, tx, from, to, TxTypePay, Amount(300), "msg")
}

func TestUser_Claim_errors(t *testing.T) {
	t.Run("same user", func(t *testing.T) {
		u := newUser()
		_, err := u.Claim(u, 100, "")
		require.Equal(t, ErrSameUser, err)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		cases := []struct {
			amount int64
			isErr  bool
		}{
			{-5000, true},
			{-4901, true},
			{-4900, false},
		}
		for _, c := range cases {
			from := newUser()
			to := newUser()
			to.balance.amount = Amount(c.amount)
			_, err := from.Claim(to, 100, "")
			if c.isErr {
				require.Equal(t, ErrInsufficientBalance, err)
			} else {
				require.NoError(t, err)
			}
		}
	})

	t.Run("amount is 0", func(t *testing.T) {
		from := newUser()
		to := newUser()
		_, err := from.Claim(to, 0, "")
		require.Equal(t, ErrZeroAmount, err)
	})
}

func TestUser_Claim_success(t *testing.T) {
	from, to := newUser(), newUser()
	to.balance.amount = 1000

	invoice, err := from.Claim(to, 300, "msg")
	require.NoError(t, err)

	require.Equal(t, from.ID, invoice.FromID)
	require.Equal(t, to.ID, invoice.ToID)
	require.Equal(t, Amount(300), invoice.Amount)
	require.Equal(t, "msg", invoice.Message)
}

func TestUser_AcceptInvoice_errors(t *testing.T) {
	t.Run("wrong destination", func(t *testing.T) {
		invoice := &Invoice{
			ToID: "dummy",
		}
		u := newUser()
		_, err := u.AcceptInvoice(invoice, nil)
		require.Equal(t, ErrWrongDestination, err)
	})
}

func TestUser_AcceptInvoice_success(t *testing.T) {
	from, to := newUser(), newUser()
	to.balance.amount = 1000

	invoice := &Invoice{
		ID:      InvoiceID(id.New()),
		FromID:  from.ID,
		ToID:    to.ID,
		Amount:  Amount(300),
		Message: "msg",
	}

	tx, err := to.AcceptInvoice(invoice, from)
	require.NoError(t, err)

	assertTx(t, tx, from, to, TxTypeClaim, Amount(300), "msg")
}

func assertTx(t *testing.T, tx *Transaction, from, to *User, txType TxType, amount Amount, message string) {
	t.Helper()
	require.Equal(t, txType, tx.Type)
	require.Equal(t, from.ID, tx.From)
	require.Equal(t, to.ID, tx.To)
	require.Equal(t, amount, tx.Amount)
	require.Equal(t, message, tx.Message)
}
