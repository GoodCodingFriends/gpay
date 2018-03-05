package usecase

import (
	"testing"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository/repositorytest"
	"github.com/stretchr/testify/require"
)

func TestClaim(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := repositorytest.NewInMemory()

		from, to := createUser(t, repo), createUser(t, repo)
		amount := entity.Amount(300)
		p := &ClaimParam{
			FromID:  from.ID,
			ToID:    to.ID,
			Amount:  amount,
			Message: "msg",
		}

		_, err := Claim(repo, p)
		require.NoError(t, err)
	})
}

func TestAcceptInvoice(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		repo := repositorytest.NewInMemory()
		repo.Transaction = &txRepositoryWithError{}

		from, to := createUser(t, repo), createUser(t, repo)
		param := &AcceptInvoiceParam{
			InvoiceID: entity.InvoiceID("id"),
			FromID:    from.ID,
			ToID:      to.ID,
		}

		_, err := AcceptInvoice(repo, param)
		require.Error(t, err)

		from, to = getUser(t, repo, from.ID), getUser(t, repo, to.ID)

		require.Equal(t, entity.Amount(0), entity.GetBalanceAmount(from))
		require.Equal(t, entity.Amount(0), entity.GetBalanceAmount(to))
	})
}
