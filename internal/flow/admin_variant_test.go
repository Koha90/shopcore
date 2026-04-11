package flow

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type variantProductListerStub struct {
	items []ProductListItem
	err   error
}

func (s *variantProductListerStub) ListProducts(ctx context.Context) ([]ProductListItem, error) {
	return s.items, s.err
}

type variantCreatorStub struct {
	called bool
	params CreateVariantParams
	err    error
}

func (s *variantCreatorStub) CreateVariant(ctx context.Context, params CreateVariantParams) error {
	s.called = true
	s.params = params
	return s.err
}

func openAdminVariantCreate(t *testing.T, svc *Service, key SessionKey) {
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
		ActionID:      ActionAdminVariantCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
}

func TestHandleAction_AdminVariantCreateStart_ShowsProductSelect(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	lister := &variantProductListerStub{
		items: []ProductListItem{
			{ID: 1, Code: "rose-box", Label: "Rose Box"},
			{ID: 2, Code: "gift-box", Label: "Gift Box"},
		},
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, lister, nil)
	key := testSessionKey("shop-admin-variant")

	openAdminVariantCreate(t, svc, key)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminVariantProductSelect, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)

	vm := svc.buildAdminVariantProductSelectScreen()
	require.Equal(t, "Новый вариант\n\nВыберите товар:", vm.Text)
	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 1)
	require.Len(t, vm.Inline.Sections[0].Actions, 3)
	require.Equal(t, adminVariantSelectProductAction(1), vm.Inline.Sections[0].Actions[0].ID)
	require.Equal(t, "Rose Box", vm.Inline.Sections[0].Actions[0].Label)
}

func TestHandleAction_AdminVariantSelectProduct_StartsNameInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	lister := &variantProductListerStub{
		items: []ProductListItem{
			{ID: 7, Code: "rose-box", Label: "Rose Box"},
		},
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, lister, nil)
	key := testSessionKey("shop-admin-variant")

	openAdminVariantCreate(t, svc, key)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminVariantSelectProductAction(7),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый вариант\n\nТовар: Rose Box\n\nВведите название варианта сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminVariantCreate, session.Current)
	require.Equal(t, PendingInputVariantName, session.Pending.Kind)
	require.Equal(t, "7", session.Pending.Value(PendingValueProductID))
	require.Equal(t, "Rose Box", session.Pending.Value(PendingValueProductName))
}

func TestHandleText_AdminVariantCreate_AutoCodeSuccess(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &variantCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, creator)
	key := testSessionKey("shop-admin-variant")

	store.Put(key, Session{
		Current:  ScreenAdminVariantCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputVariantName,
			Payload: PendingInputPayload{
				PendingValueProductID:   "7",
				PendingValueProductName: "Rose Box",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "L / 25 шт",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый вариант\n\nВариант создан.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, 7, creator.params.ProductID)
	require.Equal(t, "l-25-sht", creator.params.Code)
	require.Equal(t, "L / 25 шт", creator.params.Name)
}

func TestHandleText_AdminVariantCode_CreateError_KeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &variantCreatorStub{
		err: errors.New("duplicate variant code"),
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, creator)
	key := testSessionKey("shop-admin-variant")

	store.Put(key, Session{
		Current:  ScreenAdminVariantCode,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputVariantCode,
			Payload: PendingInputPayload{
				PendingValueProductID:   "7",
				PendingValueProductName: "Rose Box",
				PendingValueName:        "L / 25 шт",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "l-25",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый вариант\n\nТовар: Rose Box\n\nНе удалось создать вариант. Попробуйте другой code.\n\nВведите code варианта сообщением.", vm.Text)
}
