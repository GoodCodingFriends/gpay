package usecase

import (
	"context"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository"
	multierror "github.com/hashicorp/go-multierror"
)

// errUserNotFound has UserID which identifies the missing user
type errUserNotFound struct {
	err error
	id  entity.UserID
}

func (e errUserNotFound) Error() string {
	return e.err.Error()
}

func FindBothUsers(repo *repository.Repository, fromID, toID entity.UserID) (*entity.User, *entity.User, error) {
	var result error
	from, err := repo.User.FindByID(context.Background(), fromID)
	if err != nil {
		result = multierror.Append(result, errUserNotFound{err: err, id: fromID})
	}
	to, err := repo.User.FindByID(context.Background(), toID)
	if err != nil {
		result = multierror.Append(result, errUserNotFound{err: err, id: toID})
	}
	return from, to, result
}

// for testing
var findBothUsers func(*repository.Repository, entity.UserID, entity.UserID) (*entity.User, *entity.User, error) = FindBothUsers

func FindBothUsersWithUserCreation(cfg *config.Config, repo *repository.Repository, fromID, toID entity.UserID) (*entity.User, *entity.User, error) {
	from, to, baseErr := findBothUsers(repo, fromID, toID)
	if baseErr == nil {
		return from, to, nil
	}

	merr, ok := baseErr.(*multierror.Error)
	if !ok {
		return nil, nil, baseErr
	}

	users := make([]*entity.User, 0, len(merr.Errors))
	for _, err := range merr.Errors {
		uerr, ok := err.(errUserNotFound)
		if !ok {
			return nil, nil, baseErr
		}

		// TODO: first, lastname, display name
		u := entity.NewUser(cfg, uerr.id, "", "", "")
		users = append(users, u)

		if fromID == uerr.id {
			from = u
		} else if toID == uerr.id {
			to = u
		}
	}

	if err := repo.User.StoreAll(context.Background(), users); err != nil {
		return nil, nil, err
	}

	return from, to, nil
}
