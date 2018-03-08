package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/entity/entitytest"
	"github.com/GoodCodingFriends/gpay/repository"
	"github.com/GoodCodingFriends/gpay/repository/repositorytest"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

type userRepositoryWithError struct {
	repository.UserRepository
	err error
}

func (r userRepositoryWithError) FindByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	return nil, r.err
}

func TestFindBothUsers(t *testing.T) {
	t.Run("not found users", func(t *testing.T) {
		repo := repositorytest.NewInMemory()
		repo.User = &userRepositoryWithError{err: errors.New("an error")}

		from, to := entity.UserID("hazuki"), entity.UserID("sapphire")
		_, _, err := FindBothUsers(repo, from, to)
		require.Error(t, err)
		merr, ok := err.(*multierror.Error)
		require.True(t, ok)
		require.Len(t, merr.Errors, 2)

		ids := []entity.UserID{from, to}
		for i, err := range merr.Errors {
			err, ok := err.(errUserNotFound)
			require.True(t, ok)
			require.Equal(t, err.id, ids[i])
		}
	})

	t.Run("found", func(t *testing.T) {
		repo := repositorytest.NewInMemory()

		from, to := entitytest.NewUser(t), entitytest.NewUser(t)
		repo.User.StoreAll(context.TODO(), []*entity.User{from, to})

		from2, to2, err := FindBothUsers(repo, from.ID, to.ID)
		require.NoError(t, err)
		require.Exactly(t, from, from2)
		require.Exactly(t, to, to2)
	})
}

func TestFindBothUsersWithUserCreation(t *testing.T) {
	cfg, err := config.Process()
	require.NoError(t, err)

	setFindBothUsers := func(err error) {
		// replace global findBothUsers function
		findBothUsers = func(_ *repository.Repository, _, _ entity.UserID) (*entity.User, *entity.User, error) {
			return nil, nil, err
		}
	}

	fromID, toID := entity.UserID("kumiko"), entity.UserID("reina")

	t.Run("baseErr is nil", func(t *testing.T) {
		repo := repositorytest.NewInMemory()
		setFindBothUsers(nil)

		from, to, err := FindBothUsersWithUserCreation(cfg, repo, fromID, toID)
		require.Nil(t, from)
		require.Nil(t, to)
		require.Nil(t, err)
	})

	t.Run("baseErr is not multierror", func(t *testing.T) {
		repo := repositorytest.NewInMemory()
		baseErr := errors.New("an error")
		setFindBothUsers(baseErr)

		_, _, err := FindBothUsersWithUserCreation(cfg, repo, fromID, toID)
		require.Equal(t, baseErr, err)
	})

	t.Run("baseErr is multierror, but not errUserNotFound", func(t *testing.T) {
		var baseErr error
		baseErr = multierror.Append(baseErr, errors.New("an error"))
		repo := repositorytest.NewInMemory()
		setFindBothUsers(baseErr)

		_, _, err := FindBothUsersWithUserCreation(cfg, repo, fromID, toID)
		require.Equal(t, baseErr, err)
	})

	t.Run("baseErr is multierror, and errUserNotFound", func(t *testing.T) {
		var baseErr error
		uerr := errUserNotFound{id: entity.UserID("kumiko")}
		baseErr = multierror.Append(baseErr, uerr)
		repo := repositorytest.NewInMemory()
		setFindBothUsers(baseErr)

		_, _, err := FindBothUsersWithUserCreation(cfg, repo, fromID, toID)
		require.NoError(t, err)
		users, err := repo.User.FindAll(context.TODO())
		require.NoError(t, err)
		require.Len(t, users, 1)
	})
}
