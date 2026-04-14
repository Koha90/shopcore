package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleAction_CatalogStart_ReplyWelcomeOpensCompactRoot(t *testing.T) {
	svc := NewService(nil)
	key := SessionKey{BotID: "bot-1", ChatID: 1, UserID: 1}

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "bot-1",
		StartScenario: string(StartScenarioReplyWelcome),
		SessionKey:    key,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "bot-1",
		StartScenario: string(StartScenarioReplyWelcome),
		ActionID:      ActionCatalogStart,
		SessionKey:    key,
	})
	require.NoError(t, err)
	require.Equal(t, "Каталог\n\nВыберите раздел:", vm.Text)
	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 1)
	require.Equal(t, 1, vm.Inline.Sections[0].Columns)
}

func TestHandleAction_CatalogStart_InlineCatalogOpensExtendedRoot(t *testing.T) {
	svc := NewService(nil)
	key := SessionKey{BotID: "bot-2", ChatID: 2, UserID: 2}

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "bot-2",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    key,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "bot-2",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionCatalogStart,
		SessionKey:    key,
	})
	require.NoError(t, err)
	require.Equal(t, "Каталог\n\nВыберите раздел:", vm.Text)
	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 2)
	require.Equal(t, 2, vm.Inline.Sections[0].Columns)
	require.Equal(t, 1, vm.Inline.Sections[1].Columns)
}

func TestHandleAction_CatalogStart_UsesProvidedCatalogRootInCompactMode(t *testing.T) {
	custom := Catalog{
		Schema: CatalogSchema{
			Levels: []CatalogLevel{LevelCity},
		},
		Roots: []CatalogNode{
			{
				Level: LevelCity,
				ID:    "custom-city",
				Label: "Кастомный город",
			},
		},
	}

	svc := NewServiceWithCatalogProvider(
		nil,
		NewStaticCatalogProvider(custom),
	)

	key := testSessionKey("shop-reply")

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-reply",
		BotName:       "Reply Shop",
		StartScenario: string(StartScenarioReplyWelcome),
		SessionKey:    key,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-reply",
		BotName:       "Reply Shop",
		StartScenario: string(StartScenarioReplyWelcome),
		ActionID:      ActionCatalogStart,
		SessionKey:    key,
	})
	require.NoError(t, err)

	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 1)
	require.Len(t, vm.Inline.Sections[0].Actions, 1)
	require.Equal(t, "Кастомный город", vm.Inline.Sections[0].Actions[0].Label)
	require.Equal(t, ActionID("catalog:select:city:custom-city"), vm.Inline.Sections[0].Actions[0].ID)
}

func TestHandleAction_CatalogSelectionChain(t *testing.T) {
	svc := NewService(nil)
	key := testSessionKey("shop-inline")

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    key,
	})
	require.NoError(t, err)

	steps := []struct {
		name   string
		action ActionID
		want   string
	}{
		{
			name:   "city",
			action: catalogSelectAction(LevelCity, "moscow"),
			want:   "Москва\n\nВыберите категорию:",
		},
		{
			name:   "category",
			action: catalogSelectAction(LevelCategory, "flowers"),
			want:   "Цветы\n\nВыберите район:",
		},
		{
			name:   "district",
			action: catalogSelectAction(LevelDistrict, "center"),
			want:   "Центр\n\nВыберите товар:",
		},
		{
			name:   "product",
			action: catalogSelectAction(LevelProduct, "rose-box"),
			want:   "Rose Box\n\nКомпозиция из роз для центрального района.\n\nВыберите вариант:",
		},
		{
			name:   "variant",
			action: catalogSelectAction(LevelVariant, "large"),
			want:   "L / 25 шт\n\n5900 ₽\n\nБольшая упаковка.",
		},
	}

	for _, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			vm, err := svc.HandleAction(context.Background(), ActionRequest{
				BotID:         "shop-inline",
				BotName:       "Inline Shop",
				StartScenario: string(StartScenarioInlineCatalog),
				ActionID:      step.action,
				SessionKey:    key,
			})
			require.NoError(t, err)
			require.Equal(t, step.want, vm.Text)
			require.NotNil(t, vm.Inline)
		})
	}
}

func TestHandleAction_CatalogSelection_WrongLevelOrder(t *testing.T) {
	svc := NewService(nil)
	key := testSessionKey("shop-inline")

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    key,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      catalogSelectAction(LevelDistrict, "center"),
		SessionKey:    key,
	})
	require.ErrorIs(t, err, ErrUnknownAction)
}

func TestHandleAction_CatalogSelection_UnknownNode(t *testing.T) {
	svc := NewService(nil)
	key := testSessionKey("shop-inline")

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    key,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      catalogSelectAction(LevelCity, "unknown"),
		SessionKey:    key,
	})
	require.ErrorIs(t, err, ErrUnknownAction)
}

func TestHandleAction_CatalogBackChain(t *testing.T) {
	svc := NewService(nil)
	key := testSessionKey("shop-inline")

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    key,
	})
	require.NoError(t, err)

	for _, action := range []ActionID{
		catalogSelectAction(LevelCity, "moscow"),
		catalogSelectAction(LevelCategory, "flowers"),
		catalogSelectAction(LevelDistrict, "center"),
		catalogSelectAction(LevelProduct, "rose-box"),
		catalogSelectAction(LevelVariant, "large"),
	} {
		_, err := svc.HandleAction(context.Background(), ActionRequest{
			BotID:         "shop-inline",
			BotName:       "Inline Shop",
			StartScenario: string(StartScenarioInlineCatalog),
			ActionID:      action,
			SessionKey:    key,
		})
		require.NoError(t, err)
	}

	backChecks := []string{
		"Rose Box\n\nКомпозиция из роз для центрального района.\n\nВыберите вариант:",
		"Центр\n\nВыберите товар:",
		"Цветы\n\nВыберите район:",
		"Москва\n\nВыберите категорию:",
		"Каталог\n\nВыберите раздел:",
	}

	for _, want := range backChecks {
		vm, err := svc.HandleAction(context.Background(), ActionRequest{
			BotID:         "shop-inline",
			BotName:       "Inline Shop",
			StartScenario: string(StartScenarioInlineCatalog),
			ActionID:      ActionBack,
			SessionKey:    key,
		})
		require.NoError(t, err)
		require.Equal(t, want, vm.Text)
	}
}

func TestHandleAction_BackFromRoot_StaysOnRoot(t *testing.T) {
	svc := NewService(nil)
	key := testSessionKey("shop-inline")

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    key,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionBack,
		SessionKey:    key,
	})
	require.NoError(t, err)
	require.Equal(t, "Каталог\n\nВыберите раздел:", vm.Text)
}
