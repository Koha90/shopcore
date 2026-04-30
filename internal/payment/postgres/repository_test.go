package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListBotPaymentMethodsRejectsNilRepository(t *testing.T) {
	t.Parallel()

	var repo *Repository

	_, err := repo.ListBotPaymentMethods(context.Background(), "shop-main")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "repository is nil")
}

func TestListBotPaymentMethodsRejectsNilPool(t *testing.T) {
	t.Parallel()

	repo := NewRepository(nil)

	_, err := repo.ListBotPaymentMethods(context.Background(), "shop-main")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "pool is nil")
}
