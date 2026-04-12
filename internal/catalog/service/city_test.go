package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type cityWriterStub struct {
	params CreateCityParams
	called bool
	err    error
}

func (s *cityWriterStub) CreateCity(ctx context.Context, params CreateCityParams) error {
	s.called = true
	s.params = params
	return s.err
}

func TestCreateCity_Valid(t *testing.T) {
	writer := &cityWriterStub{}
	svc := New(nil, writer, nil, nil, nil, nil, nil)

	err := svc.CreateCity(context.Background(), CreateCityParams{
		Code:      " Moscow ",
		Name:      " Москва ",
		NameLatin: " Moscow ",
		SortOrder: 10,
		// IsActive:    true,
	})
	require.NoError(t, err)

	require.True(t, writer.called)
	require.Equal(t, "moscow", writer.params.Code)
	require.Equal(t, "Москва", writer.params.Name)
	require.Equal(t, "Moscow", writer.params.NameLatin)
	require.Equal(t, 10, writer.params.SortOrder)
	// require.True(t, writer.params.IsActive)
}

func TestCreateCity_CodeEmpty(t *testing.T) {
	writer := &cityWriterStub{}
	svc := New(nil, writer, nil, nil, nil, nil, nil)

	err := svc.CreateCity(context.Background(), CreateCityParams{
		Code: "",
		Name: "Москва",
	})
	require.ErrorIs(t, err, ErrCityCodeEmpty)
	require.False(t, writer.called)
}

func TestCreateCity_NameEmpty(t *testing.T) {
	writer := &cityWriterStub{}
	svc := New(nil, writer, nil, nil, nil, nil, nil)

	err := svc.CreateCity(context.Background(), CreateCityParams{
		Code: "moskva",
		Name: "   ",
	})
	require.ErrorIs(t, err, ErrCityNameEmpty)
	require.False(t, writer.called)
}

func TestCreateCity_NilWriter(t *testing.T) {
	svc := New(nil, nil, nil, nil, nil, nil, nil)

	err := svc.CreateCity(context.Background(), CreateCityParams{
		Code: "moskva",
		Name: "Москва",
	})
	require.EqualError(t, err, "city writer is nil")
}
