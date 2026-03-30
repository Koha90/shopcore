package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func testSessionKey(botID string) SessionKey {
	return SessionKey{
		BotID:  botID,
		ChatID: 1,
		UserID: 1,
	}
}

func TestNormalizeStartScenario_Default(t *testing.T) {
	t.Parallel()

	got := NormalizeStartScenario("")
	if got != StartScenarioReplyWelcome {
		t.Fatalf("expected default scenario %q, got %q", StartScenarioReplyWelcome, got)
	}
}

func TestStart_ReplyWelcome(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	vm, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-reply",
		BotName:       "Reply Shop",
		StartScenario: string(StartScenarioReplyWelcome),
		SessionKey:    testSessionKey("shop-reply"),
	})
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}

	if vm.Reply == nil {
		t.Fatal("expected reply keyboard")
	}
	if vm.Inline != nil {
		t.Fatal("did not expect inline keyboard")
	}
	if vm.Text != "Добро пожаловать 👋\nВыберите раздел:" {
		t.Fatalf("expected welcome text, got %q", vm.Text)
	}
	if len(vm.Reply.Rows) != 2 {
		t.Fatalf("expected 2 reply rows, got %d", len(vm.Reply.Rows))
	}
}

func TestStart_InlineCatalog(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	vm, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    testSessionKey("shop-inline"),
	})
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}

	if vm.Inline == nil {
		t.Fatal("expected inline keyboard")
	}
	if vm.Reply != nil {
		t.Fatal("did not expect reply keyboard")
	}
	if vm.Text != "Каталог\n\nВыберите раздел:" {
		t.Fatalf("expected catalog text, got %q", vm.Text)
	}
	if len(vm.Inline.Sections) != 2 {
		t.Fatalf("expected 2 inline sections, got %d", len(vm.Inline.Sections))
	}
	if vm.Inline.Sections[0].Columns != 2 {
		t.Fatalf("expected first section columns=2, got %d", vm.Inline.Sections[0].Columns)
	}
	if vm.Inline.Sections[1].Columns != 1 {
		t.Fatalf("expected second section columns=1, got %d", vm.Inline.Sections[1].Columns)
	}
}

func TestHandleAction_RootExtendedReturnsExtendedCatalog(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionRootExtended,
		SessionKey:    testSessionKey("shop-inline"),
	})
	if err != nil {
		t.Fatalf("HandleAction returned error: %v", err)
	}

	if vm.Inline == nil {
		t.Fatal("expected inline keyboard")
	}
	if vm.Text != "Каталог\n\nВыберите раздел:" {
		t.Fatalf("expected catalog text, got %q", vm.Text)
	}
	if len(vm.Inline.Sections) != 2 {
		t.Fatalf("expected 2 inline sections, got %d", len(vm.Inline.Sections))
	}
	if vm.Inline.Sections[0].Columns != 2 {
		t.Fatalf("expected first section columns=2, got %d", vm.Inline.Sections[0].Columns)
	}
	if vm.Inline.Sections[1].Columns != 1 {
		t.Fatalf("expected second section columns=1, got %d", vm.Inline.Sections[1].Columns)
	}
}

func TestResolveReplyAction(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	tests := []struct {
		name string
		text string
		want ActionID
		ok   bool
	}{
		{name: "catalog plain", text: "Каталог", want: ActionCatalogStart, ok: true},
		{name: "catalog emoji", text: "♻️ Каталог", want: ActionCatalogStart, ok: true},
		{name: "cabinet plain", text: "Мой кабинет", want: ActionCabinetOpen, ok: true},
		{name: "cabinet emoji", text: "⚙️ Мой кабинет", want: ActionCabinetOpen, ok: true},
		{name: "support plain", text: "Поддержка", want: ActionSupportOpen, ok: true},
		{name: "support emoji", text: "🤷‍♂️ Поддержка", want: ActionSupportOpen, ok: true},
		{name: "reviews plain", text: "Отзывы", want: ActionReviewsOpen, ok: true},
		{name: "reviews emoji", text: "📨 Отзывы", want: ActionReviewsOpen, ok: true},
		{name: "unknown", text: "Что-то ещё", want: "", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := svc.ResolveReplyAction(tt.text)
			if ok != tt.ok {
				t.Fatalf("expected ok=%v, got %v", tt.ok, ok)
			}
			if got != tt.want {
				t.Fatalf("expected action %q, got %q", tt.want, got)
			}
		})
	}
}

func TestHandleAction_Unknown(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-main",
		BotName:       "Main Shop",
		StartScenario: string(StartScenarioReplyWelcome),
		ActionID:      ActionID("unknown:action"),
		SessionKey:    testSessionKey("shop-reply"),
	})
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
	if err != ErrUnknownAction {
		t.Fatalf("expected ErrUnknownAction, got %v", err)
	}
}

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
