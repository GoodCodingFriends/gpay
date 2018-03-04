package entitytest

import "github.com/GoodCodingFriends/gpay/entity"

func NewUser() *entity.User {
	return &entity.User{
		ID:          UserID(newID()),
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
