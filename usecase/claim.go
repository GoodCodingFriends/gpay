package usecase

import (
	"context"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository"
)

type ClaimParam struct {
	FromID, ToID entity.UserID
	Amount       entity.Amount
	Message      string
}

func Claim(repo *repository.Repository, p *ClaimParam) (*entity.Invoice, error) {
	from, to, err := findBothUsers(repo, p.FromID, p.ToID)
	if err != nil {
		return nil, err
	}

	invoice, err := from.Claim(to, p.Amount, p.Message)
	if err != nil {
		return nil, err
	}

	err = repo.Invoice.Store(context.TODO(), invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

type AcceptInvoiceParam struct {
	InvoiceID    entity.InvoiceID
	FromID, ToID entity.UserID
}

func AcceptInvoice(repo *repository.Repository, p *AcceptInvoiceParam) (*entity.Transaction, error) {
	invoice, err := repo.Invoice.FindByID(context.Background(), p.InvoiceID)
	if err != nil {
		return nil, err
	}

	from, to, err := findBothUsers(repo, p.FromID, p.ToID)
	if err != nil {
		return nil, err
	}

	return to.AcceptInvoice(invoice, from)
}

func findBothUsers(repo *repository.Repository, fromID, toID entity.UserID) (*entity.User, *entity.User, error) {
	from, err := repo.User.FindByID(context.Background(), fromID)
	if err != nil {
		return nil, nil, err
	}
	to, err := repo.User.FindByID(context.Background(), toID)
	if err != nil {
		return nil, nil, err
	}
	return from, to, nil
}
