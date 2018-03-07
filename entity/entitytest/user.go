package entitytest

import (
	"testing"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/stretchr/testify/require"
)

func NewUser(t *testing.T) *entity.User {
	cfg, err := config.Process()
	require.NoError(t, err)
	return entity.NewUser(cfg, "kumiko", "omae", "omae-chan")
}
