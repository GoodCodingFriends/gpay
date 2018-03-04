package usecase

import (
	"context"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository"
)

type PayParam struct {
	From, To *entity.User
	Amount   entity.Amount
	Message  string
}

func Pay(repo repository.Repository, p *PayParam) (*entity.Transaction, error) {
	if result, err := p.From.Pay(p.To, p.Amount, p.Message); err != nil {
		return nil, err
	}

	tx := repo.BeginTx()

	var err error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = tx.User.StoreAll(context.Background(), []*entity.User{from, to})
	if err != nil {
		return nil, err
	}

	err = tx.Transaction.Store(context.Background(), result)
	if err != nil {
		return nil, err
	}

	return result, tx.Commit()
}
