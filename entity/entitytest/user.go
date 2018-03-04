package entitytest

import (
	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/entity/internal/id"
)

func NewUser() *entity.User {
	return &entity.User{
		ID:          entity.UserID(id.NewID()),
		FirstName:   "kumiko",
		LastName:    "omae",
		DisplayName: "omae-chan",
		balance:     &entity.Balance{},
	}
}

func NewUserWithBalance(amount int64) *entity.User {
	u := NewUser()
	u.(*entity.Balance).Amount = Amount(amount)
	return u
}
