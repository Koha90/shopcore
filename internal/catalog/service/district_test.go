package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type districtWriterStub struct {
	params CreateDistrictParams
	called bool
	err    error
}

func (s *districtWriterStub) CreateDistrict(ctx context.Context, params CreateDistrictParams) error {
	s.called = true
	s.params = params
	return s.err
}

func TestCreateDistrict_Valid(t *testing.T) {
	writer := &districtWriterStub{}
	svc := New(nil, nil, writer, nil, nil, nil, nil, nil, nil)

	err := svc.CreateDistrict(context.Background(), CreateDistrictParams{
		CityID:    7,
		Code:      " Center ",
		Name:      " Центр ",
		NameLatin: " Center ",
		SortOrder: 10,
	})
	require.NoError(t, err)

	require.True(t, writer.called)
	require.Equal(t, 7, writer.params.CityID)
	require.Equal(t, "center", writer.params.Code)
	require.Equal(t, "Центр", writer.params.Name)
	require.Equal(t, "Center", writer.params.NameLatin)
	require.Equal(t, 10, writer.params.SortOrder)
}

func TestCreateDistrict_CityIDInvalid(t *testing.T) {
	writer := &districtWriterStub{}
	svc := New(nil, nil, writer, nil, nil, nil, nil, nil, nil)

	err := svc.CreateDistrict(context.Background(), CreateDistrictParams{
		CityID: 0,
		Code:   "center",
		Name:   "Центр",
	})
	require.ErrorIs(t, err, ErrDistrictCityIDInvalid)
	require.False(t, writer.called)
}

func TestCreateDistrict_CodeEmpty(t *testing.T) {
	writer := &districtWriterStub{}
	svc := New(nil, nil, writer, nil, nil, nil, nil, nil, nil)

	err := svc.CreateDistrict(context.Background(), CreateDistrictParams{
		CityID: 1,
		Code:   "",
		Name:   "Центр",
	})
	require.ErrorIs(t, err, ErrDistrictCodeEmpty)
	require.False(t, writer.called)
}

func TestCreateDistrict_NameEmpty(t *testing.T) {
	writer := &districtWriterStub{}
	svc := New(nil, nil, writer, nil, nil, nil, nil, nil, nil)

	err := svc.CreateDistrict(context.Background(), CreateDistrictParams{
		CityID: 1,
		Code:   "center",
		Name:   "   ",
	})
	require.ErrorIs(t, err, ErrDistrictNameEmpty)
	require.False(t, writer.called)
}

func TestCreateDistrict_NilWriter(t *testing.T) {
	svc := New(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	err := svc.CreateDistrict(context.Background(), CreateDistrictParams{
		CityID: 1,
		Code:   "center",
		Name:   "Центр",
	})
	require.EqualError(t, err, "district writer is nil")
}
