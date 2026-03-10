package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"botmanager/internal/domain"
)

func TestUserService_Create(t *testing.T) {
	repo := &stubUserRepository{}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewUserService(repo, tx, bus, nil)

	tgID := int64(123456789)

	user, err := svc.Create(context.Background(), domain.NewUserParams{
		TgID:   &tgID,
		TgName: "koha",
	})

	require.NoError(t, err)
	require.NotNil(t, user)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, repo.saveCalls)
	require.NotNil(t, repo.savedUser)
	require.Same(t, user, repo.savedUser)
}

func TestUserService_Create_DomainError(t *testing.T) {
	repo := &stubUserRepository{}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewUserService(repo, tx, bus, nil)

	user, err := svc.Create(context.Background(), domain.NewUserParams{})

	require.Nil(t, user)
	require.ErrorIs(t, err, domain.ErrInvalidCredentials)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, repo.saveCalls)
	require.Nil(t, repo.savedUser)
}

func TestUserService_Create_SaveError(t *testing.T) {
	repo := &stubUserRepository{
		saveErr: errors.New("save failed"),
	}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewUserService(repo, tx, bus, nil)

	tgID := int64(123456789)

	user, err := svc.Create(context.Background(), domain.NewUserParams{
		TgID:   &tgID,
		TgName: "koha",
	})

	require.Nil(t, user)
	require.EqualError(t, err, "save failed")

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, repo.saveCalls)
	require.NotNil(t, repo.savedUser)
}

func TestUserService_Create_TransactionError(t *testing.T) {
	repo := &stubUserRepository{}
	tx := &stubTxManager{
		err: errors.New("transaction failed"),
	}
	bus := &stubEventBus{}

	svc := NewUserService(repo, tx, bus, nil)

	tgID := int64(123456789)

	user, err := svc.Create(context.Background(), domain.NewUserParams{
		TgID:   &tgID,
		TgName: "koha",
	})

	require.Nil(t, user)
	require.EqualError(t, err, "transaction failed")

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, repo.saveCalls)
	require.Nil(t, repo.savedUser)
}
