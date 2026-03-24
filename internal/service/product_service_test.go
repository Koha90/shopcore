package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/domain"
)

func TestProductService_Create(t *testing.T) {
	repo := &stubProductRepository{}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewProductService(repo, tx, bus, nil)

	product, err := svc.Create(
		context.Background(),
		"Amnesia",
		1,
		"good stuff",
		"/tmp/img.png",
	)
	require.NoError(t, err)
	require.NotNil(t, product)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, repo.saveCalls)
	require.NotNil(t, repo.savedProduct)
	require.Same(t, product, repo.savedProduct)
}

func TestProductService_Create_DomainError(t *testing.T) {
	repo := &stubProductRepository{}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewProductService(repo, tx, bus, nil)

	product, err := svc.Create(
		context.Background(),
		"",
		1,
		"good stuff",
		"/tmp/img.png",
	)
	require.Nil(t, product)
	require.Error(t, err)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, repo.saveCalls)
	require.Nil(t, repo.savedProduct)
}

func TestProductService_Create_SaveError(t *testing.T) {
	repo := &stubProductRepository{
		saveErr: errors.New("save failed"),
	}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewProductService(repo, tx, bus, nil)

	product, err := svc.Create(
		context.Background(),
		"Amnesia",
		1,
		"good stuff",
		"/tmp/img.png",
	)

	require.Nil(t, product)
	require.EqualError(t, err, "save failed")

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, repo.saveCalls)
	require.NotNil(t, repo.savedProduct)
}

func TestProductService_Create_TransactionError(t *testing.T) {
	repo := &stubProductRepository{}
	tx := &stubTxManager{
		err: errors.New("transaction failed"),
	}
	bus := &stubEventBus{}

	svc := NewProductService(repo, tx, bus, nil)

	product, err := svc.Create(
		context.Background(),
		"Amnesia",
		1,
		"good stuff",
		"/tmp/img.png",
	)

	require.Nil(t, product)
	require.EqualError(t, err, "transaction failed")

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, repo.saveCalls)
	require.Nil(t, repo.savedProduct)
}

func TestProductService_AddVariant(t *testing.T) {
	product, err := domain.NewProduct(
		"Amnesia",
		1,
		"good stuff",
		"/tmp/img.png",
	)
	require.NoError(t, err)

	repo := &stubProductRepository{
		product: product,
	}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewProductService(repo, tx, bus, nil)

	err = svc.AddVariant(context.Background(), product.ID(), "1g", 10, 1500)
	require.NoError(t, err)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, repo.saveCalls)
	require.NotNil(t, repo.savedProduct)

	variants := repo.savedProduct.Variants()
	require.Len(t, variants, 1)
	require.Equal(t, "1g", variants[0].PackSize())
	require.Equal(t, 10, variants[0].DistrictID())
	require.EqualValues(t, 1500, variants[0].Price())
}

func TestProductService_AddVariant_LoadError(t *testing.T) {
	repo := &stubProductRepository{
		byIDErr: errors.New("load failed"),
	}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewProductService(repo, tx, bus, nil)

	err := svc.AddVariant(context.Background(), 100, "1g", 10, 1500)

	require.EqualError(t, err, "load failed")
	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, repo.saveCalls)
}

func TestProductService_AddVariant_DomainError(t *testing.T) {
	product, err := domain.NewProduct(
		"Amnesia",
		1,
		"good stuff",
		"/tmp/img.png",
	)
	require.NoError(t, err)

	repo := &stubProductRepository{
		product: product,
	}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewProductService(repo, tx, bus, nil)

	err = svc.AddVariant(context.Background(), product.ID(), "", 10, 1500)

	require.Error(t, err)
	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, repo.saveCalls)
}

func TestProductService_AddVariant_SaveError(t *testing.T) {
	product, err := domain.NewProduct(
		"Amnesia",
		1,
		"good stuff",
		"/tmp/img.png",
	)
	require.NoError(t, err)

	repo := &stubProductRepository{
		product: product,
		saveErr: errors.New("save failed"),
	}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewProductService(repo, tx, bus, nil)

	err = svc.AddVariant(context.Background(), product.ID(), "1g", 10, 1500)

	require.EqualError(t, err, "save failed")
	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, repo.saveCalls)
	require.NotNil(t, repo.savedProduct)
}

func TestNewProductService_PanicsOnNilRepo(t *testing.T) {
	require.Panics(t, func() {
		NewProductService(nil, &stubTxManager{}, &stubEventBus{}, nil)
	})
}

func TestNewProductService_PanicsOnNilTx(t *testing.T) {
	require.Panics(t, func() {
		NewProductService(&stubProductRepository{}, nil, &stubEventBus{}, nil)
	})
}

func TestNewProductService_PanicsOnNilBus(t *testing.T) {
	require.Panics(t, func() {
		NewProductService(&stubProductRepository{}, &stubTxManager{}, nil, nil)
	})
}

func TestNewProductService_UsesDefaultLogger(t *testing.T) {
	svc := NewProductService(&stubProductRepository{}, &stubTxManager{}, &stubEventBus{}, nil)
	require.NotNil(t, svc)
	require.NotNil(t, svc.logger)
}
