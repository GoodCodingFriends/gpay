package usecase

import (
	"context"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository"
)

type ListTransactionsParam struct {
	User *entity.User
}

// ListTransactions lists all transactions related to User
// if User is nil, list all transactions
func ListTransactions(repo *repository.Repository, p *ListTransactionsParam) ([]*entity.Transaction, error) {
	if p.User != nil {
		return repo.Transaction.FindAllByUserID(context.Background(), p.User.ID)
	}
	return repo.Transaction.FindAll(context.Background())
}
