package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestHandleAction_Back_ClearsPendingInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{}
	svc := NewServiceWithDeps(store, nil, ServiceDeps{Categories: creator})
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

func TestHandleAction_ProductScreen_ShowsVariantButtonsWithPrice(t *testing.T) {
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

	var vm ViewModel

	for _, action := range []ActionID{
		catalogSelectAction(LevelCity, "moscow"),
		catalogSelectAction(LevelCategory, "flowers"),
		catalogSelectAction(LevelDistrict, "center"),
		catalogSelectAction(LevelProduct, "rose-box"),
	} {
		vm, err = svc.HandleAction(context.Background(), ActionRequest{
			BotID:         "shop-inline",
			BotName:       "Inline Shop",
			StartScenario: string(StartScenarioInlineCatalog),
			ActionID:      action,
			SessionKey:    key,
		})
		require.NoError(t, err)
	}

	require.Equal(t, "Rose Box\n\nКомпозиция из роз для центрального района.\n\nВыберите вариант:", vm.Text)
	require.NotNil(t, vm.Inline)
	require.True(t, hasInlineActionLabel(vm, "L / 25 шт - 5900 ₽"))
}

func TestHandleAction_NonVariantButtons_KeepPlainLabels(t *testing.T) {
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

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      catalogSelectAction(LevelCity, "moscow"),
		SessionKey:    key,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      catalogSelectAction(LevelCategory, "flowers"),
		SessionKey:    key,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      catalogSelectAction(LevelDistrict, "center"),
		SessionKey:    key,
	})
	require.NoError(t, err)

	require.Equal(t, "Центр\n\nВыберите товар:", vm.Text)
	require.NotNil(t, vm.Inline)
	require.True(t, hasInlineActionLabel(vm, "Rose Box"))
	require.False(t, hasInlineActionLabel(vm, "Rose Box • 5900 ₽"))
}

func hasInlineActionLabel(vm ViewModel, want string) bool {
	if vm.Inline == nil {
		return false
	}

	for _, section := range vm.Inline.Sections {
		for _, action := range section.Actions {
			if action.Label == want {
				return true
			}
		}
	}

	return false
}

func TestSyncSessionAccess_UpdatesAdminFlagWhenGranted(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewService(store)

	key := testSessionKey("shop-inline")
	session := Session{
		Current:  ScreenRootExtended,
		CanAdmin: false,
	}

	got := svc.syncSessionAccess(key, session, true, string(StartScenarioInlineCatalog))
	require.True(t, got.CanAdmin)
}

func TestSyncSessionAccess_DropsAdminStateWhenRevoked(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewService(store)

	key := testSessionKey("shop-inline")
	session := Session{
		Current:  ScreenAdminCatalog,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputCategoryName,
		},
	}

	got := svc.syncSessionAccess(key, session, false, string(StartScenarioInlineCatalog))

	require.False(t, got.CanAdmin)
	require.Equal(t, ScreenRootExtended, got.Current)
	require.Nil(t, got.History)
	require.Equal(t, PendingInputNone, got.Pending.Kind)
}
