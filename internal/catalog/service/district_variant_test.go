package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type districtVariantWriterStub struct {
	params CreateDistrictVariantParams
	called bool
	err    error
}

func (s *districtVariantWriterStub) CreateDistrictVariant(ctx context.Context, params CreateDistrictVariantParams) error {
	s.called = true
	s.params = params
	return s.err
}

func TestCreateDistrictVariant_Valid(t *testing.T) {
	writer := &districtVariantWriterStub{}
	svc := New(nil, nil, nil, nil, nil, writer, nil, nil, nil)

	err := svc.CreateDistrictVariant(context.Background(), CreateDistrictVariantParams{
		DistrictID: 7,
		VariantID:  9,
		Price:      5900,
	})
	require.NoError(t, err)

	require.True(t, writer.called)
	require.Equal(t, 7, writer.params.DistrictID)
	require.Equal(t, 9, writer.params.VariantID)
	require.Equal(t, 5900, writer.params.Price)
}

func TestCreateDistrictVariant_DistrictIDInvalid(t *testing.T) {
	writer := &districtVariantWriterStub{}
	svc := New(nil, nil, nil, nil, nil, writer, nil, nil, nil)

	err := svc.CreateDistrictVariant(context.Background(), CreateDistrictVariantParams{
		DistrictID: 0,
		VariantID:  9,
		Price:      5900,
	})
	require.ErrorIs(t, err, ErrDistrictVariantDistrictIDInvalid)
	require.False(t, writer.called)
}

func TestCreateDistrictVariant_VariantIDInvalid(t *testing.T) {
	writer := &districtVariantWriterStub{}
	svc := New(nil, nil, nil, nil, nil, writer, nil, nil, nil)

	err := svc.CreateDistrictVariant(context.Background(), CreateDistrictVariantParams{
		DistrictID: 7,
		VariantID:  0,
		Price:      5900,
	})
	require.ErrorIs(t, err, ErrDistrictVariantVariantIDInvalid)
	require.False(t, writer.called)
}

func TestCreateDistrictVariant_PriceInvalid(t *testing.T) {
	writer := &districtVariantWriterStub{}
	svc := New(nil, nil, nil, nil, nil, writer, nil, nil, nil)

	err := svc.CreateDistrictVariant(context.Background(), CreateDistrictVariantParams{
		DistrictID: 7,
		VariantID:  9,
		Price:      0,
	})
	require.ErrorIs(t, err, ErrDistrictVariantPriceInvalid)
	require.False(t, writer.called)
}

func TestCreateDistrictVariant_NilWriter(t *testing.T) {
	svc := New(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	err := svc.CreateDistrictVariant(context.Background(), CreateDistrictVariantParams{
		DistrictID: 7,
		VariantID:  9,
		Price:      5900,
	})
	require.EqualError(t, err, "district variant writer is nil")
}
