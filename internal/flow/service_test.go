package flow

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type categoryCreatorStub struct {
	called bool
	params CreateCategoryParams
	err    error
}

func (s *categoryCreatorStub) CreateCategory(ctx context.Context, params CreateCategoryParams) error {
	s.called = true
	s.params = params
	return s.err
}

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

func TestStart_InlineCatalog_UsesProvidedCatalogRoot(t *testing.T) {
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

	vm, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    testSessionKey("shop-inline"),
		CanAdmin:      true,
	})
	require.NoError(t, err)

	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 2)

	require.Len(t, vm.Inline.Sections[0].Actions, 1)
	require.Equal(t, 2, vm.Inline.Sections[0].Columns)
	require.Equal(t, "Кастомный город", vm.Inline.Sections[0].Actions[0].Label)
	require.Equal(t, ActionID("catalog:select:city:custom-city"), vm.Inline.Sections[0].Actions[0].ID)

	require.Equal(t, 1, vm.Inline.Sections[1].Columns)
	require.Len(t, vm.Inline.Sections[1].Actions, 4)
	require.Equal(t, ActionBalanceOpen, vm.Inline.Sections[1].Actions[0].ID)
	require.Equal(t, ActionBotsMine, vm.Inline.Sections[1].Actions[1].ID)
	require.Equal(t, ActionOrderLast, vm.Inline.Sections[1].Actions[2].ID)
	require.Equal(t, ActionAdminOpen, vm.Inline.Sections[1].Actions[3].ID)
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

func TestHandleText_AdminCategoryCreate_Success(t *testing.T) {
	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCategoryCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          " Цветы ",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новая категория\n\nКатегория создана.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCreateDone, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)

	backVM, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionBack,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Админка · Каталог\n\nВыберите действие:", backVM.Text)
}

func TestHandleText_WithoutPending_RendersCurrentScreen(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)
	key := testSessionKey("shop-inline")

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    key,
	})
	require.NoError(t, err)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "какой-то текст",
		SessionKey:    key,
	})
	require.NoError(t, err)
	require.Equal(t, "Каталог\n\nВыберите раздел:", vm.Text)
	require.NotNil(t, vm.Inline)
	require.Nil(t, vm.Reply)
}

func TestHandleText_AdminCategoryCreate_EmptyTextKeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCategoryCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "   ",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новая категория\n\nНазвание категории не может быть пустым.\n\nВведите название категории сообщением.", vm.Text)
	require.NotNil(t, vm.Inline)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCreate, session.Current)
	require.Equal(t, PendingInputCategoryName, session.Pending.Kind)
}

func TestHandleAction_AdminCategoryCreateStart_PendingIsClearedByRegularAction(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCategoryCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, PendingInputCategoryName, session.Pending.Kind)
	require.Equal(t, ScreenAdminCategoryCreate, session.Current)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Админка\n\nВыберите раздел:", vm.Text)

	session, ok = store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminRoot, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}

func TestHandleAction_Back_ClearsPendingInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCategoryCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionBack,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Админка · Каталог\n\nВыберите действие:", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCatalog, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}

func TestHandleText_AdminCategoryCreate_CallsCategoryCreator(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCategoryCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          " Цветы ",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новая категория\n\nКатегория создана.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, "tsvety", creator.params.Code)
	require.Equal(t, "Цветы", creator.params.Name)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCreateDone, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}

func TestHandleText_AdminCategoryCreate_AutoCodeError_OpensManualCodeInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{
		err: errors.New("create category failed"),
	}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCategoryCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Цветы",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новая категория\n\nНе удалось создать категорию с автоматическим code.\n\nВведите code категории сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCode, session.Current)
	require.Equal(t, PendingInputCategoryCode, session.Pending.Kind)
	require.Equal(t, "Цветы", session.Pending.Value(PendingValueName))
	require.Equal(t, "tsvety", session.Pending.Value(PendingValueCode))
}

func TestHandleText_AdminCategoryCreate_NilCategoryCreator(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, nil, nil)
	key := testSessionKey("shop-admin")

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCategoryCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Цветы",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.EqualError(t, err, "flow category creator is nil")
}

func TestHandleAction_AdminCategoryCreateStart_InitializesPendingPayload(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCategoryCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCreate, session.Current)
	require.Equal(t, PendingInputCategoryName, session.Pending.Kind)
	require.Nil(t, session.Pending.Payload)
	require.Equal(t, "", session.Pending.Value(PendingValueName))
}

func TestHandleText_AdminCategoryCreate_StoresNameInPendingPayload(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCategoryCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          " Цветы ",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	require.True(t, creator.called)
	require.Equal(t, "tsvety", creator.params.Code)
	require.Equal(t, "Цветы", creator.params.Name)
}

func TestPendingInput_SetValueAndValue(t *testing.T) {
	t.Parallel()

	var pending PendingInput

	require.False(t, pending.Active())
	require.Equal(t, "", pending.Value(PendingValueName))

	pending.Kind = PendingInputCategoryName
	pending.SetValue(PendingValueName, "Цветы")

	require.True(t, pending.Active())
	require.Equal(t, "Цветы", pending.Value(PendingValueName))
}

func openAdminCategoryCreate(t *testing.T, svc *Service, key SessionKey) {
	t.Helper()

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCategoryCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
}

func TestHandleText_AdminCategoryCreate_AutoCodeSuccess(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	openAdminCategoryCreate(t, svc, key)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Цветы",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новая категория\n\nКатегория создана.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, "tsvety", creator.params.Code)
	require.Equal(t, "Цветы", creator.params.Name)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCreateDone, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}

func TestHandleText_AdminCategoryCreate_AutoCodeFailure_OpensManualCodeInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{
		err: errors.New("duplicate category code"),
	}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	openAdminCategoryCreate(t, svc, key)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Тестовая категория",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новая категория\n\nНе удалось создать категорию с автоматическим code.\n\nВведите code категории сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCode, session.Current)
	require.Equal(t, PendingInputCategoryCode, session.Pending.Kind)
	require.Equal(t, "Тестовая категория", session.Pending.Value(PendingValueName))
	require.Equal(t, "testovaya-kategoriya", session.Pending.Value(PendingValueCode))
}

func TestHandleText_AdminCategoryCreate_EmptySuggestedCode_OpensManualCodeInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	openAdminCategoryCreate(t, svc, key)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "!!!",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новая категория\n\nНе удалось автоматически подобрать code.\n\nВведите code категории сообщением.", vm.Text)

	require.False(t, creator.called)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCode, session.Current)
	require.Equal(t, PendingInputCategoryCode, session.Pending.Kind)
	require.Equal(t, "!!!", session.Pending.Value(PendingValueName))
	require.Equal(t, "", session.Pending.Value(PendingValueCode))
}

func TestHandleText_AdminCategoryCode_EmptyTextKeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	store.Put(key, Session{
		Current:  ScreenAdminCategoryCode,
		History:  []ScreenID{ScreenAdminCatalog},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputCategoryCode,
			Payload: PendingInputPayload{
				PendingValueName: "Цветы",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "   ",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новая категория\n\nCode категории не может быть пустым.\n\nВведите code категории сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCode, session.Current)
	require.Equal(t, PendingInputCategoryCode, session.Pending.Kind)
	require.Equal(t, "Цветы", session.Pending.Value(PendingValueName))
}

func TestHandleText_AdminCategoryCode_Success_UsesManualCode(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	store.Put(key, Session{
		Current:  ScreenAdminCategoryCode,
		History:  []ScreenID{ScreenAdminCatalog},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputCategoryCode,
			Payload: PendingInputPayload{
				PendingValueName: "Цветы",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "flowers-manual",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новая категория\n\nКатегория создана.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, "flowers-manual", creator.params.Code)
	require.Equal(t, "Цветы", creator.params.Name)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCreateDone, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}

func TestHandleText_AdminCategoryCode_CreateError_KeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{
		err: errors.New("duplicate category code"),
	}
	svc := NewServiceWithDeps(store, nil, creator, nil)
	key := testSessionKey("shop-admin")

	store.Put(key, Session{
		Current:  ScreenAdminCategoryCode,
		History:  []ScreenID{ScreenAdminCatalog},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputCategoryCode,
			Payload: PendingInputPayload{
				PendingValueName: "Цветы",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "flowers-manual",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новая категория\n\nНе удалось создать категорию. Попробуйте другой code.\n\nВведите code категории сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCategoryCode, session.Current)
	require.Equal(t, PendingInputCategoryCode, session.Pending.Kind)
	require.Equal(t, "Цветы", session.Pending.Value(PendingValueName))
	require.Equal(t, "flowers-manual", session.Pending.Value(PendingValueCode))
}

func TestStart_InlineCatalog_NonAdmin_HidesAdminButton(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	vm, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    testSessionKey("shop-inline"),
		CanAdmin:      false,
	})
	require.NoError(t, err)

	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 2)
	require.Len(t, vm.Inline.Sections[1].Actions, 3)
	require.Equal(t, ActionBalanceOpen, vm.Inline.Sections[1].Actions[0].ID)
	require.Equal(t, ActionBotsMine, vm.Inline.Sections[1].Actions[1].ID)
	require.Equal(t, ActionOrderLast, vm.Inline.Sections[1].Actions[2].ID)
}

func TestStart_InlineCatalog_Admin_ShowsAdminButton(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	vm, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    testSessionKey("shop-inline-admin"),
		CanAdmin:      true,
	})
	require.NoError(t, err)

	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 2)
	require.Len(t, vm.Inline.Sections[1].Actions, 4)
	require.Equal(t, ActionAdminOpen, vm.Inline.Sections[1].Actions[3].ID)
}

func TestHandleAction_AdminOpen_NonAdmin_ReturnsUnknownAction(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    testSessionKey("shop-inline"),
		CanAdmin:      false,
	})
	require.ErrorIs(t, err, ErrUnknownAction)
}
