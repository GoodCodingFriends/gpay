package usecase

import (
	"testing"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/entity/entitytest"
	"github.com/stretchr/testify/require"
)

func TestPay(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &inMemoryRepository{
			user: &userInMemoryRepository,
		}

		from, to := entitytest.NewUserWithBalance(1000), entitytest.NewUser()
		param := &PayParam{
			From:    from,
			To:      to,
			Amount:  entity.Amount(300),
			Message: "msg",
		}

		_, err := Pay(repo, param)
		require.NoError(t, err)
	})
}
