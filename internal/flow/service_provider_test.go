package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type failingCatalogProvider struct {
	err error
}

func (p failingCatalogProvider) Catalog(ctx context.Context) (Catalog, error) {
	return Catalog{}, p.err
}

func TestStart_ReturnsCatalogProviderError(t *testing.T) {
	svc := NewServiceWithCatalogProvider(nil, failingCatalogProvider{err: ErrUnknownAction})

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    testSessionKey("shop-inline"),
	})

	require.ErrorIs(t, err, ErrUnknownAction)
}

func TestHandleAction_ReturnsCatalogProviderError(t *testing.T) {
	svc := NewServiceWithCatalogProvider(nil, failingCatalogProvider{err: ErrUnknownAction})

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionCatalogStart,
		SessionKey:    testSessionKey("shop-inline"),
	})

	require.ErrorIs(t, err, ErrUnknownAction)
}
