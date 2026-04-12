package flow

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type cityListerStub struct {
	items []CityListItem
	err   error
}

func (s *cityListerStub) ListCities(ctx context.Context) ([]CityListItem, error) {
	return s.items, s.err
}

type districtCreatorStub struct {
	called bool
	params CreateDistrictParams
	err    error
}

func (s *districtCreatorStub) CreateDistrict(ctx context.Context, params CreateDistrictParams) error {
	s.called = true
	s.params = params
	return s.err
}

func openAdminDistrictCreate(t *testing.T, svc *Service, key SessionKey) {
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
		ActionID:      ActionAdminDistrictCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
}

func TestHandleAction_AdminDistrictCreateStart_ShowsCitySelect(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	lister := &cityListerStub{
		items: []CityListItem{
			{ID: 1, Code: "moskva", Label: "Москва"},
			{ID: 2, Code: "spb", Label: "Санкт-Петербург"},
		},
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, lister, nil, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

	openAdminDistrictCreate(t, svc, key)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictCitySelect, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)

	vm := svc.buildAdminDistrictCitySelectScreen()
	require.Equal(t, "Новый район\n\nВыберите город:", vm.Text)
	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 1)
	require.Len(t, vm.Inline.Sections[0].Actions, 3)
	require.Equal(t, adminDistrictSelectCityAction(1), vm.Inline.Sections[0].Actions[0].ID)
	require.Equal(t, "Москва", vm.Inline.Sections[0].Actions[0].Label)
}

func TestHandleAction_AdminDistrictSelectCity_StartsNameInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	lister := &cityListerStub{
		items: []CityListItem{
			{ID: 7, Code: "moskva", Label: "Москва"},
		},
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, lister, nil, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

	openAdminDistrictCreate(t, svc, key)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictSelectCityAction(7),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый район\n\nГород: Москва\n\nВведите название района сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictCreate, session.Current)
	require.Equal(t, PendingInputDistrictName, session.Pending.Kind)
	require.Equal(t, "7", session.Pending.Value(PendingValueCityID))
	require.Equal(t, "Москва", session.Pending.Value(PendingValueCityName))
}

func TestHandleText_AdminDistrictCreate_EmptyTextKeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &districtCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

	store.Put(key, Session{
		Current:  ScreenAdminDistrictCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputDistrictName,
			Payload: PendingInputPayload{
				PendingValueCityID:   "7",
				PendingValueCityName: "Москва",
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
	require.Equal(t, "Новый район\n\nГород: Москва\n\nНазвание района не может быть пустым.\n\nВведите название района сообщением.", vm.Text)
	require.False(t, creator.called)
}

func TestHandleText_AdminDistrictCreate_AutoCodeSuccess(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &districtCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

	store.Put(key, Session{
		Current:  ScreenAdminDistrictCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputDistrictName,
			Payload: PendingInputPayload{
				PendingValueCityID:   "7",
				PendingValueCityName: "Москва",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Центр",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый район\n\nРайон создан.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, 7, creator.params.CityID)
	require.Equal(t, "tsentr", creator.params.Code)
	require.Equal(t, "Центр", creator.params.Name)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictCreateDone, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}

func TestHandleText_AdminDistrictCreate_AutoCodeFailure_OpensManualCodeInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &districtCreatorStub{
		err: errors.New("duplicate district code"),
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

	store.Put(key, Session{
		Current:  ScreenAdminDistrictCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputDistrictName,
			Payload: PendingInputPayload{
				PendingValueCityID:   "7",
				PendingValueCityName: "Москва",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Центр",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый район\n\nГород: Москва\n\nНе удалось создать район с автоматическим code.\n\nВведите code района сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictCode, session.Current)
	require.Equal(t, PendingInputDistrictCode, session.Pending.Kind)
	require.Equal(t, "7", session.Pending.Value(PendingValueCityID))
	require.Equal(t, "Москва", session.Pending.Value(PendingValueCityName))
	require.Equal(t, "Центр", session.Pending.Value(PendingValueName))
	require.Equal(t, "tsentr", session.Pending.Value(PendingValueCode))
}

func TestHandleText_AdminDistrictCreate_EmptySuggestedCode_OpensManualCodeInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &districtCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

	store.Put(key, Session{
		Current:  ScreenAdminDistrictCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputDistrictName,
			Payload: PendingInputPayload{
				PendingValueCityID:   "7",
				PendingValueCityName: "Москва",
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
	require.Equal(t, "Новый район\n\nГород: Москва\n\nНе удалось автоматически подобрать code.\n\nВведите code района сообщением.", vm.Text)

	require.False(t, creator.called)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictCode, session.Current)
	require.Equal(t, PendingInputDistrictCode, session.Pending.Kind)
	require.Equal(t, "!!!", session.Pending.Value(PendingValueName))
	require.Equal(t, "", session.Pending.Value(PendingValueCode))
}

func TestHandleText_AdminDistrictCreate_NilDistrictCreator(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

	store.Put(key, Session{
		Current:  ScreenAdminDistrictCreate,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputDistrictName,
			Payload: PendingInputPayload{
				PendingValueCityID:   "7",
				PendingValueCityName: "Москва",
			},
		},
	})

	_, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Центр",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.EqualError(t, err, "flow district creator is nil")
}

func TestHandleText_AdminDistrictCode_EmptyTextKeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &districtCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

	store.Put(key, Session{
		Current:  ScreenAdminDistrictCode,
		History:  []ScreenID{ScreenAdminDistrictCitySelect},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputDistrictCode,
			Payload: PendingInputPayload{
				PendingValueCityID:   "7",
				PendingValueCityName: "Москва",
				PendingValueName:     "Центр",
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
	require.Equal(t, "Новый район\n\nГород: Москва\n\nCode района не может быть пустым.\n\nВведите code района сообщением.", vm.Text)
}

func TestHandleText_AdminDistrictCode_Success_UsesManualCode(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &districtCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

	store.Put(key, Session{
		Current:  ScreenAdminDistrictCode,
		History:  []ScreenID{ScreenAdminDistrictCitySelect},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputDistrictCode,
			Payload: PendingInputPayload{
				PendingValueCityID:   "7",
				PendingValueCityName: "Москва",
				PendingValueName:     "Центр",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "center-manual",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый район\n\nРайон создан.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, 7, creator.params.CityID)
	require.Equal(t, "center-manual", creator.params.Code)
	require.Equal(t, "Центр", creator.params.Name)
}

func TestHandleText_AdminDistrictCode_CreateError_KeepsPending(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &districtCreatorStub{
		err: errors.New("duplicate district code"),
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, creator, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

	store.Put(key, Session{
		Current:  ScreenAdminDistrictCode,
		History:  []ScreenID{ScreenAdminDistrictCitySelect},
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputDistrictCode,
			Payload: PendingInputPayload{
				PendingValueCityID:   "7",
				PendingValueCityName: "Москва",
				PendingValueName:     "Центр",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "center-manual",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Новый район\n\nГород: Москва\n\nНе удалось создать район. Попробуйте другой code.\n\nВведите code района сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictCode, session.Current)
	require.Equal(t, PendingInputDistrictCode, session.Pending.Kind)
	require.Equal(t, "center-manual", session.Pending.Value(PendingValueCode))
}

func TestHandleAction_AdminDistrictCreateStart_NilCityLister(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	key := testSessionKey("shop-admin-district")

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
		ActionID:      ActionAdminDistrictCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.EqualError(t, err, "flow city lister is nil")
}

func TestHandleAction_AdminDistrictCreateStart_NonAdmin_ReturnsUnknownAction(t *testing.T) {
	t.Parallel()

	svc := NewService(nil)

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminDistrictCreateStart,
		SessionKey:    testSessionKey("shop-inline-district"),
		CanAdmin:      false,
	})
	require.ErrorIs(t, err, ErrUnknownAction)
}
