package entitytest

import (
	"testing"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/entity/internal/id"
	"github.com/stretchr/testify/require"
)

type generator struct{}

func (g *generator) Generate() entity.UserID {
	return entity.UserID(id.New())
}

func NewUser(t *testing.T) *entity.User {
	cfg, err := config.Process()
	require.NoError(t, err)
	return entity.NewUser(cfg, entity.UserID(id.New()), "kumiko", "omae", "omae-chan")
}
