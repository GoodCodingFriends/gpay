package adapter

import "github.com/GoodCodingFriends/gpay/entity"

type SlackUserIDGenerator struct {
	UserID entity.UserID
}

func (g *SlackUserIDGenerator) Generate() entity.UserID {
	return g.UserID
}
