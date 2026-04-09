package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type variantWriterStub struct {
	params CreateVariantParams
	called bool
	err    error
}

func (s *variantWriterStub) CreateVariant(ctx context.Context, params CreateVariantParams) error {
	s.called = true
	s.params = params
	return s.err
}

func TestCreateVariant_Valid(t *testing.T) {
	writer := &variantWriterStub{}
	svc := New(nil, nil, nil, nil, writer)

	err := svc.CreateVariant(context.Background(), CreateVariantParams{
		ProductID:   7,
		Code:        " L-25 ",
		Name:        " L / 25 шт ",
		NameLatin:   " L / 25 pcs ",
		Description: " Большая упаковка ",
		SortOrder:   10,
	})
	require.NoError(t, err)

	require.True(t, writer.called)
	require.Equal(t, 7, writer.params.ProductID)
	require.Equal(t, "l-25", writer.params.Code)
	require.Equal(t, "L / 25 шт", writer.params.Name)
	require.Equal(t, "L / 25 pcs", writer.params.NameLatin)
	require.Equal(t, "Большая упаковка", writer.params.Description)
	require.Equal(t, 10, writer.params.SortOrder)
}

func TestCreateVariant_ProductIDInvalid(t *testing.T) {
	writer := &variantWriterStub{}
	svc := New(nil, nil, nil, nil, writer)

	err := svc.CreateVariant(context.Background(), CreateVariantParams{
		ProductID: 0,
		Code:      "l-25",
		Name:      "L / 25 шт",
	})
	require.ErrorIs(t, err, ErrVariantProductIDInvalid)
	require.False(t, writer.called)
}

func TestCreateVariant_CodeEmpty(t *testing.T) {
	writer := &variantWriterStub{}
	svc := New(nil, nil, nil, nil, writer)

	err := svc.CreateVariant(context.Background(), CreateVariantParams{
		ProductID: 1,
		Code:      "",
		Name:      "L / 25 шт",
	})
	require.ErrorIs(t, err, ErrVariantCodeEmpty)
	require.False(t, writer.called)
}

func TestCreateVariant_NameEmpty(t *testing.T) {
	writer := &variantWriterStub{}
	svc := New(nil, nil, nil, nil, writer)

	err := svc.CreateVariant(context.Background(), CreateVariantParams{
		ProductID: 1,
		Code:      "l-25",
		Name:      "   ",
	})
	require.ErrorIs(t, err, ErrVariantNameEmpty)
	require.False(t, writer.called)
}

func TestCreateVariant_NilWriter(t *testing.T) {
	svc := New(nil, nil, nil, nil, nil)

	err := svc.CreateVariant(context.Background(), CreateVariantParams{
		ProductID: 1,
		Code:      "l-25",
		Name:      "L / 25 шт",
	})
	require.EqualError(t, err, "variant writer is nil")
}
