package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type categoryWriterStub struct {
	params CreateCategoryParams
	called bool
	err    error
}

func (s *categoryWriterStub) CreateCategory(ctx context.Context, params CreateCategoryParams) error {
	s.called = true
	s.params = params
	return s.err
}

func TestCreateCategory_Valid(t *testing.T) {
	writer := &categoryWriterStub{}
	svc := New(writer, nil, nil, nil)

	err := svc.CreateCategory(context.Background(), CreateCategoryParams{
		Code:        " Flowers ",
		Name:        " Цветы ",
		NameLatin:   " Flowers ",
		Description: " Букеты и композиции ",
		SortOrder:   10,
		// IsActive:    true,
	})
	require.NoError(t, err)

	require.True(t, writer.called)
	require.Equal(t, "flowers", writer.params.Code)
	require.Equal(t, "Цветы", writer.params.Name)
	require.Equal(t, "Flowers", writer.params.NameLatin)
	require.Equal(t, "Букеты и композиции", writer.params.Description)
	require.Equal(t, 10, writer.params.SortOrder)
	// require.True(t, writer.params.IsActive)
}

func TestCreateCategory_CodeEmpty(t *testing.T) {
	writer := &categoryWriterStub{}
	svc := New(writer, nil, nil, nil)

	err := svc.CreateCategory(context.Background(), CreateCategoryParams{
		Code: "",
		Name: "Цветы",
	})
	require.ErrorIs(t, err, ErrCategoryCodeEmpty)
	require.False(t, writer.called)
}

func TestCreateCategory_NameEmpty(t *testing.T) {
	writer := &categoryWriterStub{}
	svc := New(writer, nil, nil, nil)

	err := svc.CreateCategory(context.Background(), CreateCategoryParams{
		Code: "flowers",
		Name: "   ",
	})
	require.ErrorIs(t, err, ErrCategoryNameEmpty)
	require.False(t, writer.called)
}

func TestCreateCategory_NilWriter(t *testing.T) {
	svc := New(nil, nil, nil, nil)

	err := svc.CreateCategory(context.Background(), CreateCategoryParams{
		Code: "flowers",
		Name: "Цветы",
	})
	require.EqualError(t, err, "category writer is nil")
}
