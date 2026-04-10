package flow

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type cityCreatorStub struct {
	called bool
	params CreateCityParams
	err    error
}

func (s *cityCreatorStub) CreateCity(ctx context.Context, params CreateCityParams) error {
	s.called = true
	s.params = params
	return s.err
}

func openAdminCityCreate(t *testing.T, svc *Service, key SessionKey) {
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
		ActionID:      ActionAdminCityCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
}

func TestHandleAction_AdminCityCreateStart_InitializesPendingPayload(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &cityCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, creator, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	openAdminCityCreate(t, svc, key)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCityCreate, session.Current)
	require.Equal(t, PendingInputCityName, session.Pending.Kind)
	require.Nil(t, session.Pending.Payload)
	require.Equal(t, "", session.Pending.Value(PendingValueName))
}

func TestHandleText_AdminCityCreate_EmptyTextKeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &cityCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, creator, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	openAdminCityCreate(t, svc, key)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "   ",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый город\n\nНазвание города не может быть пустым.\n\nВведите название города сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCityCreate, session.Current)
	require.Equal(t, PendingInputCityName, session.Pending.Kind)
}

func TestHandleText_AdminCityCreate_AutoCodeSuccess(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &cityCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, creator, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	openAdminCityCreate(t, svc, key)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Москва",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый город\n\nГород создан.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, "moskva", creator.params.Code)
	require.Equal(t, "Москва", creator.params.Name)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCityCreateDone, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}

func TestHandleText_AdminCityCreate_AutoCodeFailure_OpensManualCodeInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &cityCreatorStub{
		err: errors.New("duplicate city code"),
	}
	svc := NewServiceWithDeps(store, nil, nil, creator, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	openAdminCityCreate(t, svc, key)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Тестовый город",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый город\n\nНе удалось создать город с автоматическим code.\n\nВведите code города сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCityCode, session.Current)
	require.Equal(t, PendingInputCityCode, session.Pending.Kind)
	require.Equal(t, "Тестовый город", session.Pending.Value(PendingValueName))
	require.Equal(t, "testovyy-gorod", session.Pending.Value(PendingValueCode))
}

func TestHandleText_AdminCityCreate_EmptySuggestedCode_OpensManualCodeInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &cityCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, creator, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	openAdminCityCreate(t, svc, key)

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "!!!",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый город\n\nНе удалось автоматически подобрать code.\n\nВведите code города сообщением.", vm.Text)

	require.False(t, creator.called)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCityCode, session.Current)
	require.Equal(t, PendingInputCityCode, session.Pending.Kind)
	require.Equal(t, "!!!", session.Pending.Value(PendingValueName))
	require.Equal(t, "", session.Pending.Value(PendingValueCode))
}

func TestHandleText_AdminCityCreate_NilCityCreator(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	openAdminCityCreate(t, svc, key)

	_, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Москва",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.EqualError(t, err, "flow city creator is nil")
}

func TestHandleText_AdminCityCode_EmptyTextKeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &cityCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, creator, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	store.Put(key, Session{
		Current:  ScreenAdminCityCode,
		History:  []ScreenID{ScreenAdminCatalog},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputCityCode,
			Payload: PendingInputPayload{
				PendingValueName: "Москва",
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
	require.Equal(t, "Новый город\n\nCode города не может быть пустым.\n\nВведите code города сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCityCode, session.Current)
	require.Equal(t, PendingInputCityCode, session.Pending.Kind)
	require.Equal(t, "Москва", session.Pending.Value(PendingValueName))
}

func TestHandleText_AdminCityCode_Success_UsesManualCode(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &cityCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, creator, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	store.Put(key, Session{
		Current:  ScreenAdminCityCode,
		History:  []ScreenID{ScreenAdminCatalog},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputCityCode,
			Payload: PendingInputPayload{
				PendingValueName: "Москва",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "moscow-manual",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый город\n\nГород создан.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, "moscow-manual", creator.params.Code)
	require.Equal(t, "Москва", creator.params.Name)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCityCreateDone, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}

func TestHandleText_AdminCityCode_CreateError_KeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &cityCreatorStub{
		err: errors.New("duplicate city code"),
	}
	svc := NewServiceWithDeps(store, nil, nil, creator, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	store.Put(key, Session{
		Current:  ScreenAdminCityCode,
		History:  []ScreenID{ScreenAdminCatalog},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputCityCode,
			Payload: PendingInputPayload{
				PendingValueName: "Москва",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "moscow-manual",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый город\n\nНе удалось создать город. Попробуйте другой code.\n\nВведите code города сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminCityCode, session.Current)
	require.Equal(t, PendingInputCityCode, session.Pending.Kind)
	require.Equal(t, "Москва", session.Pending.Value(PendingValueName))
	require.Equal(t, "moscow-manual", session.Pending.Value(PendingValueCode))
}

func TestHandleAction_AdminCityCreateStart_PendingIsClearedByRegularAction(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &cityCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, creator, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	openAdminCityCreate(t, svc, key)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, PendingInputCityName, session.Pending.Kind)
	require.Equal(t, ScreenAdminCityCreate, session.Current)

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

func TestHandleAction_AdminCityBack_ClearsPendingInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &cityCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, creator, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-city")

	openAdminCityCreate(t, svc, key)

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

func TestHandleAction_AdminCityCreateStart_NonAdmin_ReturnsUnknownAction(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCityCreateStart,
		SessionKey:    testSessionKey("shop-inline-city"),
		CanAdmin:      false,
	})
	require.ErrorIs(t, err, ErrUnknownAction)
}
