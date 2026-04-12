package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type districtVariantPriceUpdaterStub struct {
	params UpdateDistrictVariantPriceParams
	called bool
	err    error
}

func (s *districtVariantPriceUpdaterStub) UpdateDistrictVariantPrice(
	ctx context.Context,
	params UpdateDistrictVariantPriceParams,
) error {
	s.called = true
	s.params = params
	return s.err
}

func TestUpdateDistrictVariantPrice_Valid(t *testing.T) {
	updater := &districtVariantPriceUpdaterStub{}
	svc := New(nil, nil, nil, nil, nil, nil, updater)

	err := svc.UpdateDistrictVariantPrice(context.Background(), UpdateDistrictVariantPriceParams{
		DistrictID: 7,
		VariantID:  9,
		Price:      6100,
	})
	require.NoError(t, err)

	require.True(t, updater.called)
	require.Equal(t, 7, updater.params.DistrictID)
	require.Equal(t, 9, updater.params.VariantID)
	require.Equal(t, 6100, updater.params.Price)
}

func TestUpdateDistrictVariantPrice_DistrictIDInvalid(t *testing.T) {
	updater := &districtVariantPriceUpdaterStub{}
	svc := New(nil, nil, nil, nil, nil, nil, updater)

	err := svc.UpdateDistrictVariantPrice(context.Background(), UpdateDistrictVariantPriceParams{
		DistrictID: 0,
		VariantID:  9,
		Price:      6100,
	})
	require.ErrorIs(t, err, ErrDistrictVariantDistrictIDInvalid)
	require.False(t, updater.called)
}

func TestUpdateDistrictVariantPrice_VariantIDInvalid(t *testing.T) {
	updater := &districtVariantPriceUpdaterStub{}
	svc := New(nil, nil, nil, nil, nil, nil, updater)

	err := svc.UpdateDistrictVariantPrice(context.Background(), UpdateDistrictVariantPriceParams{
		DistrictID: 7,
		VariantID:  0,
		Price:      6100,
	})
	require.ErrorIs(t, err, ErrDistrictVariantVariantIDInvalid)
	require.False(t, updater.called)
}

func TestUpdateDistrictVariantPrice_PriceInvalid(t *testing.T) {
	updater := &districtVariantPriceUpdaterStub{}
	svc := New(nil, nil, nil, nil, nil, nil, updater)

	err := svc.UpdateDistrictVariantPrice(context.Background(), UpdateDistrictVariantPriceParams{
		DistrictID: 7,
		VariantID:  9,
		Price:      0,
	})
	require.ErrorIs(t, err, ErrDistrictVariantPriceInvalid)
	require.False(t, updater.called)
}

func TestUpdateDistrictVariantPrice_NilUpdater(t *testing.T) {
	svc := New(nil, nil, nil, nil, nil, nil, nil)

	err := svc.UpdateDistrictVariantPrice(context.Background(), UpdateDistrictVariantPriceParams{
		DistrictID: 7,
		VariantID:  9,
		Price:      6100,
	})
	require.EqualError(t, err, "catalog district variant price updater is nil")
}
