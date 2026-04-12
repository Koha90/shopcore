package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type districtListStub struct {
	items []DistrictListItem
	err   error
}

func (s *districtListStub) ListDistricts(ctx context.Context) ([]DistrictListItem, error) {
	return s.items, s.err
}

type variantListStub struct {
	items []VariantListItem
	err   error
}

func (s *variantListStub) ListVariants(ctx context.Context) ([]VariantListItem, error) {
	return s.items, s.err
}

type districtVariantCreatorStub struct {
	called bool
	params CreateDistrictVariantParams
	err    error
}

func (s *districtVariantCreatorStub) CreateDistrictVariant(ctx context.Context, params CreateDistrictVariantParams) error {
	s.called = true
	s.params = params
	return s.err
}

func openAdminDistrictVariantCreate(t *testing.T, svc *Service, key SessionKey) {
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
		ActionID:      ActionAdminDistrictVariantCreateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
}

func TestHandleAction_AdminDistrictVariantCreateStart_ShowsDistrictSelect(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	districts := &districtListStub{
		items: []DistrictListItem{
			{ID: 1, Code: "center", Label: "Центр"},
			{ID: 2, Code: "south", Label: "Юг"},
		},
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, nil, districts, nil, nil)
	key := testSessionKey("shop-admin-district-variant")

	openAdminDistrictVariantCreate(t, svc, key)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictVariantDistrictSelect, session.Current)

	vm := svc.buildAdminDistrictVariantDistrictSelectScreen()
	require.Equal(t, "Размещение варианта\n\nВыберите район:", vm.Text)
	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections[0].Actions, 3)
}

func TestHandleAction_AdminDistrictVariantSelectDistrict_ShowsVariantSelect(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	districts := &districtListStub{
		items: []DistrictListItem{
			{ID: 7, Code: "center", Label: "Центр"},
		},
	}
	variants := &variantListStub{
		items: []VariantListItem{
			{ID: 9, Code: "large", Label: "L / 25 шт"},
		},
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, nil, districts, variants, nil)
	key := testSessionKey("shop-admin-district-variant")

	openAdminDistrictVariantCreate(t, svc, key)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictVariantSelectDistrictAction(7),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Размещение варианта\n\nРайон: Центр\n\nВыберите вариант:", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictVariantVariantSelect, session.Current)
	require.Equal(t, "7", session.Pending.Value(PendingValueDistrictID))
	require.Equal(t, "Центр", session.Pending.Value(PendingValueDistrictName))
}

func TestHandleAction_AdminDistrictVariantSelectVariant_StartsPriceInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	variants := &variantListStub{
		items: []VariantListItem{
			{ID: 9, Code: "large", Label: "L / 25 шт"},
		},
	}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, variants, nil)
	key := testSessionKey("shop-admin-district-variant")

	store.Put(key, Session{
		Current:  ScreenAdminDistrictVariantVariantSelect,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputNone,
			Payload: PendingInputPayload{
				PendingValueDistrictID:   "7",
				PendingValueDistrictName: "Центр",
			},
		},
	})

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictVariantSelectVariantAction(9),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Размещение варианта\n\nРайон: Центр\n\nВариант: L / 25 шт\n\nВведите цену сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictVariantPrice, session.Current)
	require.Equal(t, PendingInputDistrictVariantPrice, session.Pending.Kind)
	require.Equal(t, "9", session.Pending.Value(PendingValueVariantID))
}

func TestHandleText_AdminDistrictVariantPrice_Success(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	creator := &districtVariantCreatorStub{}
	svc := NewServiceWithDeps(store, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, creator)
	key := testSessionKey("shop-admin-district-variant")

	store.Put(key, Session{
		Current:  ScreenAdminDistrictVariantPrice,
		CanAdmin: true,
		Pending: PendingInput{
			Kind: PendingInputDistrictVariantPrice,
			Payload: PendingInputPayload{
				PendingValueDistrictID:   "7",
				PendingValueDistrictName: "Центр",
				PendingValueVariantID:    "9",
				PendingValueVariantName:  "L / 25 шт",
			},
		},
	})

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "5900",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Размещение варианта\n\nВариант размещён в районе.", vm.Text)

	require.True(t, creator.called)
	require.Equal(t, 7, creator.params.DistrictID)
	require.Equal(t, 9, creator.params.VariantID)
	require.Equal(t, 5900, creator.params.Price)
}
