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

func TestHandleText_AdminCategoryCreate_Success(t *testing.T) {
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

func TestHandleText_AdminCategoryCreate_EmptyTextKeepsPending(t *testing.T) {
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

func TestHandleText_AdminCategoryCreate_CallsCategoryCreator(t *testing.T) {
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
	svc := NewServiceWithDeps(store, nil, ServiceDeps{})
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

func TestHandleText_AdminCategoryCreate_AutoCodeFailure_OpensManualCodeInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &categoryCreatorStub{
		err: errors.New("duplicate category code"),
	}
	svc := NewServiceWithDeps(store, nil, ServiceDeps{Categories: creator})
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
	svc := NewServiceWithDeps(store, nil, ServiceDeps{Categories: creator})
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
	svc := NewServiceWithDeps(store, nil, ServiceDeps{Categories: creator})
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
	svc := NewServiceWithDeps(store, nil, ServiceDeps{Categories: creator})
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
	svc := NewServiceWithDeps(store, nil, ServiceDeps{Categories: creator})
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
