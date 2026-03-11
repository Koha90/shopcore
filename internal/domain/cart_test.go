package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCart(t *testing.T) {
	_, err := NewCart(0)
	require.ErrorIs(t, err, ErrInvalidUserID)

	c, err := NewCart(1)
	require.NoError(t, err)
	require.Equal(t, CartStatusActive, c.Status())
}

func TestCart_AddItem(t *testing.T) {
	c, _ := NewCart(1)

	err := c.AddItem(1, 2, 100)
	require.NoError(t, err)
	require.Len(t, c.Items(), 1)

	err = c.AddItem(1, 1, 100)
	require.ErrorIs(t, err, ErrItemAlreadyExists)
}

func TestCart_RemoveItem(t *testing.T) {
	c, _ := NewCart(1)
	_ = c.AddItem(1, 2, 100)

	err := c.RemoveItem(1)
	require.NoError(t, err)
	require.Len(t, c.Items(), 0)

	err = c.RemoveItem(1)
	require.ErrorIs(t, err, ErrItemNotFound)
}

func TestCart_ChangeQuantity(t *testing.T) {
	c, _ := NewCart(1)
	_ = c.AddItem(1, 2, 100)

	err := c.ChangeQuantity(1, 5)
	require.NoError(t, err)
	require.Equal(t, 5, c.Items()[0].quantity)

	err = c.ChangeQuantity(1, 0)
	require.ErrorIs(t, err, ErrInvalidItemQuality)
}

func TestCart_Total(t *testing.T) {
	c, _ := NewCart(1)
	_ = c.AddItem(1, 2, 100)
	_ = c.AddItem(2, 1, 50)

	require.Equal(t, int64(250), c.Total())
}

func TestCart_Checkout(t *testing.T) {
	c, _ := NewCart(1)

	err := c.Checkout()
	require.ErrorIs(t, err, ErrCartEmpty)

	_ = c.AddItem(1, 1, 100)

	require.NoError(t, c.Checkout())
	require.Equal(t, CartStatusCheckedOut, c.Status())

	err = c.AddItem(2, 1, 100)
	require.ErrorIs(t, err, ErrCartNotActive)
}
