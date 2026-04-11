package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleAction_AdminRoot_BackReturnsToCatalog(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-root")

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Админка\n\nВыберите раздел:", vm.Text)

	vm, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionBack,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Каталог\n\nВыберите раздел:", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenRootExtended, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}
