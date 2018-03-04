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

func Pay(repo *repository.Repository, p *PayParam) (*entity.Transaction, error) {
	result, err := p.From.Pay(p.To, p.Amount, p.Message)
	if err != nil {
		return nil, err
	}

	tx, err := repo.BeginTx(context.Background())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.User.StoreAll(context.Background(), []*entity.User{p.From, p.To})
	if err != nil {
		return nil, err
	}

	err = tx.Transaction.Store(context.Background(), result)
	if err != nil {
		return nil, err
	}

	return result, tx.Commit()
}
