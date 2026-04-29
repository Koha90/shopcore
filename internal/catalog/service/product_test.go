package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type productWriterStub struct {
	params CreateProductParams
	called bool
	err    error
}

func (s *productWriterStub) CreateProduct(ctx context.Context, params CreateProductParams) error {
	s.called = true
	s.params = params
	return s.err
}

func TestCreateProduct_Valid(t *testing.T) {
	writer := &productWriterStub{}
	svc := New(nil, nil, nil, writer, nil, nil, nil, nil, nil)

	err := svc.CreateProduct(context.Background(), CreateProductParams{
		CategoryID:  7,
		Code:        " Rose-Box ",
		Name:        " Розы в коробке ",
		NameLatin:   " Rose Box ",
		Description: " Композиция из роз ",
		SortOrder:   10,
	})
	require.NoError(t, err)

	require.True(t, writer.called)
	require.Equal(t, 7, writer.params.CategoryID)
	require.Equal(t, "rose-box", writer.params.Code)
	require.Equal(t, "Розы в коробке", writer.params.Name)
	require.Equal(t, "Rose Box", writer.params.NameLatin)
	require.Equal(t, "Композиция из роз", writer.params.Description)
	require.Equal(t, 10, writer.params.SortOrder)
}

func TestCreateProduct_CategoryIDInvalid(t *testing.T) {
	writer := &productWriterStub{}
	svc := New(nil, nil, nil, writer, nil, nil, nil, nil, nil)

	err := svc.CreateProduct(context.Background(), CreateProductParams{
		CategoryID: 0,
		Code:       "rose-box",
		Name:       "Розы в коробке",
	})
	require.ErrorIs(t, err, ErrProductCategoryIDInvalid)
	require.False(t, writer.called)
}

func TestCreateProduct_CodeEmpty(t *testing.T) {
	writer := &productWriterStub{}
	svc := New(nil, nil, nil, writer, nil, nil, nil, nil, nil)

	err := svc.CreateProduct(context.Background(), CreateProductParams{
		CategoryID: 1,
		Code:       "",
		Name:       "Розы в коробке",
	})
	require.ErrorIs(t, err, ErrProductCodeEmpty)
	require.False(t, writer.called)
}

func TestCreateProduct_NameEmpty(t *testing.T) {
	writer := &productWriterStub{}
	svc := New(nil, nil, nil, writer, nil, nil, nil, nil, nil)

	err := svc.CreateProduct(context.Background(), CreateProductParams{
		CategoryID: 1,
		Code:       "rose-box",
		Name:       "   ",
	})
	require.ErrorIs(t, err, ErrProductNameEmpty)
	require.False(t, writer.called)
}

func TestCreateProduct_NilWriter(t *testing.T) {
	svc := New(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	err := svc.CreateProduct(context.Background(), CreateProductParams{
		CategoryID: 1,
		Code:       "rose-box",
		Name:       "Розы в коробке",
	})
	require.EqualError(t, err, "product writer is nil")
}
