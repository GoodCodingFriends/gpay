package entitytest

import (
	"github.com/GoodCodingFriends/gpay/entity"
)

func NewUser() *entity.User {
	return entity.NewUser("kumiko", "omae", "omae-chan")
}
