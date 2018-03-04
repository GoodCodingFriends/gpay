package entitytest

import (
	"testing"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/stretchr/testify/require"
)

func NewUser(t *testing.T) *entity.User {
	conf, err := config.Process()
	require.NoError(t, err)
	return entity.NewUser(conf.Entity, "kumiko", "omae", "omae-chan")
}
