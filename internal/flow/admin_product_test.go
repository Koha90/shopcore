package flow

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type categoryListerStub struct {
	items []CategoryListItem
	err   error
}

func (s *categoryListerStub) ListCategories(ctx context.Context) ([]CategoryListItem, error) {
	return s.items, s.err
}

type productCreatorStub struct {
	called bool
	params CreateProductParams
	err    error
}

func (s *productCreatorStub) CreateProduct(ctx context.Context, params CreateProductParams) error {
	s.called = true
	s.params = params
	return s.err
}

func openAdminProductCreate(t *testing.T, svc *Service, key SessionKey) {
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
		ActionID:      ActionAdminProductCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
}

func TestHandleAction_AdminProductCreateStart_ShowsCategorySelect(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	lister := &categoryListerStub{
		items: []CategoryListItem{
			{ID: 1, Code: "flowers", Label: "Цветы"},
			{ID: 2, Code: "gifts", Label: "Подарки"},
		},
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, lister, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

	openAdminProductCreate(t, svc, key)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminProductCategorySelect, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)

	vm := svc.buildAdminProductCategorySelectScreen()
	require.Equal(t, "Новый товар\n\nВыберите категорию:", vm.Text)
	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 1)
	require.Len(t, vm.Inline.Sections[0].Actions, 3)
	require.Equal(t, adminProductSelectCategoryAction(1), vm.Inline.Sections[0].Actions[0].ID)
	require.Equal(t, "Цветы", vm.Inline.Sections[0].Actions[0].Label)
}

func TestHandleAction_AdminProductSelectCategory_StartsNameInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	lister := &categoryListerStub{
		items: []CategoryListItem{
			{ID: 7, Code: "flowers", Label: "Цветы"},
		},
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, lister, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

	openAdminProductCreate(t, svc, key)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminProductSelectCategoryAction(7),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый товар\n\nКатегория: Цветы\n\nВведите название товара сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminProductCreate, session.Current)
	require.Equal(t, PendingInputProductName, session.Pending.Kind)
	require.Equal(t, "7", session.Pending.Value(PendingValueCategoryID))
	require.Equal(t, "Цветы", session.Pending.Value(PendingValueCategoryName))
}

func TestHandleText_AdminProductCreate_EmptyTextKeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &productCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

	store.Put(key, Session{
		Current:  ScreenAdminProductCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputProductName,
			Payload: PendingInputPayload{
				PendingValueCategoryID:   "7",
				PendingValueCategoryName: "Цветы",
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
	require.Equal(t, "Новый товар\n\nКатегория: Цветы\n\nНазвание товара не может быть пустым.\n\nВведите название товара сообщением.", vm.Text)
	require.False(t, creator.called)
}

func TestHandleText_AdminProductCreate_AutoCodeSuccess(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &productCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

	store.Put(key, Session{
		Current:  ScreenAdminProductCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputProductName,
			Payload: PendingInputPayload{
				PendingValueCategoryID:   "7",
				PendingValueCategoryName: "Цветы",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Розы в коробке",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый товар\n\nТовар создан.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, 7, creator.params.CategoryID)
	require.Equal(t, "rozy-v-korobke", creator.params.Code)
	require.Equal(t, "Розы в коробке", creator.params.Name)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminProductCreateDone, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}

func TestHandleText_AdminProductCreate_AutoCodeFailure_OpensManualCodeInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &productCreatorStub{
		err: errors.New("duplicate product code"),
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

	store.Put(key, Session{
		Current:  ScreenAdminProductCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputProductName,
			Payload: PendingInputPayload{
				PendingValueCategoryID:   "7",
				PendingValueCategoryName: "Цветы",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Розы в коробке",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый товар\n\nКатегория: Цветы\n\nНе удалось создать товар с автоматическим code.\n\nВведите code товара сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminProductCode, session.Current)
	require.Equal(t, PendingInputProductCode, session.Pending.Kind)
	require.Equal(t, "7", session.Pending.Value(PendingValueCategoryID))
	require.Equal(t, "Цветы", session.Pending.Value(PendingValueCategoryName))
	require.Equal(t, "Розы в коробке", session.Pending.Value(PendingValueName))
	require.Equal(t, "rozy-v-korobke", session.Pending.Value(PendingValueCode))
}

func TestHandleText_AdminProductCreate_EmptySuggestedCode_OpensManualCodeInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &productCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

	store.Put(key, Session{
		Current:  ScreenAdminProductCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputProductName,
			Payload: PendingInputPayload{
				PendingValueCategoryID:   "7",
				PendingValueCategoryName: "Цветы",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "!!!",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый товар\n\nКатегория: Цветы\n\nНе удалось автоматически подобрать code.\n\nВведите code товара сообщением.", vm.Text)

	require.False(t, creator.called)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminProductCode, session.Current)
	require.Equal(t, PendingInputProductCode, session.Pending.Kind)
	require.Equal(t, "!!!", session.Pending.Value(PendingValueName))
	require.Equal(t, "", session.Pending.Value(PendingValueCode))
}

func TestHandleText_AdminProductCreate_NilProductCreator(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

	store.Put(key, Session{
		Current:  ScreenAdminProductCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputProductName,
			Payload: PendingInputPayload{
				PendingValueCategoryID:   "7",
				PendingValueCategoryName: "Цветы",
			},
		},
	})

	_, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Розы в коробке",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.EqualError(t, err, "flow product creator is nil")
}

func TestHandleText_AdminProductCode_EmptyTextKeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &productCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

	store.Put(key, Session{
		Current:  ScreenAdminProductCode,
		History:  []ScreenID{ScreenAdminProductCategorySelect},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputProductCode,
			Payload: PendingInputPayload{
				PendingValueCategoryID:   "7",
				PendingValueCategoryName: "Цветы",
				PendingValueName:         "Розы в коробке",
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
	require.Equal(t, "Новый товар\n\nКатегория: Цветы\n\nCode товара не может быть пустым.\n\nВведите code товара сообщением.", vm.Text)
}

func TestHandleText_AdminProductCode_Success_UsesManualCode(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &productCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

	store.Put(key, Session{
		Current:  ScreenAdminProductCode,
		History:  []ScreenID{ScreenAdminProductCategorySelect},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputProductCode,
			Payload: PendingInputPayload{
				PendingValueCategoryID:   "7",
				PendingValueCategoryName: "Цветы",
				PendingValueName:         "Розы в коробке",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "rose-box",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый товар\n\nТовар создан.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, 7, creator.params.CategoryID)
	require.Equal(t, "rose-box", creator.params.Code)
	require.Equal(t, "Розы в коробке", creator.params.Name)
}

func TestHandleText_AdminProductCode_CreateError_KeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &productCreatorStub{
		err: errors.New("duplicate product code"),
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

	store.Put(key, Session{
		Current:  ScreenAdminProductCode,
		History:  []ScreenID{ScreenAdminProductCategorySelect},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputProductCode,
			Payload: PendingInputPayload{
				PendingValueCategoryID:   "7",
				PendingValueCategoryName: "Цветы",
				PendingValueName:         "Розы в коробке",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "rose-box",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый товар\n\nКатегория: Цветы\n\nНе удалось создать товар. Попробуйте другой code.\n\nВведите code товара сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminProductCode, session.Current)
	require.Equal(t, PendingInputProductCode, session.Pending.Kind)
	require.Equal(t, "rose-box", session.Pending.Value(PendingValueCode))
}

func TestHandleAction_AdminProductCreateStart_NilCategoryLister(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-product")

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
		ActionID:      ActionAdminProductCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.EqualError(t, err, "flow category lister is nil")
}

func TestHandleAction_AdminProductCreateStart_NonAdmin_ReturnsUnknownAction(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminProductCreateStart,
		SessionKey:    testSessionKey("shop-inline-product"),
		CanAdmin:      false,
	})
	require.ErrorIs(t, err, ErrUnknownAction)
}
