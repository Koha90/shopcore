package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatalogImageInputTargetProduct(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithCatalogProvider(store, NewStaticCatalogProvider(DemoCatalog()))

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 10,
		UserID: 20,
	}

	store.Put(key, Session{
		Current: ScreenAdminProductImageInput,
		Pending: PendingInput{
			Kind: PendingInputProductImageURL,
			Payload: PendingInputPayload{
				PendingValueProductID: "123",
				PendingValueCode:      "rose-box",
			},
		},
		CanAdmin: true,
	})

	target, ok := svc.CatalogImageInputTarget(key)

	require.True(t, ok)
	assert.Equal(t, CatalogImageTargetProduct, target.Kind)
	assert.Equal(t, 123, target.EntityID)
	assert.Equal(t, "rose-box", target.EntityCode)
}

func TestCatalogImageInputTargetVariant(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithCatalogProvider(store, NewStaticCatalogProvider(DemoCatalog()))

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 10,
		UserID: 20,
	}

	store.Put(key, Session{
		Current: ScreenAdminVariantImageInput,
		Pending: PendingInput{
			Kind: PendingInputVariantImageURL,
			Payload: PendingInputPayload{
				PendingValueVariantID: "456",
				PendingValueCode:      "large-red",
			},
		},
		CanAdmin: true,
	})

	target, ok := svc.CatalogImageInputTarget(key)

	require.True(t, ok)
	assert.Equal(t, CatalogImageTargetVariant, target.Kind)
	assert.Equal(t, 456, target.EntityID)
	assert.Equal(t, "large-red", target.EntityCode)
}

func TestCatalogImageInputTargetNoPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithCatalogProvider(store, NewStaticCatalogProvider(DemoCatalog()))

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 10,
		UserID: 20,
	}

	store.Put(key, Session{
		Current:  ScreenAdminCatalog,
		Pending:  PendingInput{},
		CanAdmin: true,
	})

	_, ok := svc.CatalogImageInputTarget(key)

	assert.False(t, ok)
}
