package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type flowDistrictListerStub struct {
	items []DistrictListItem
	err   error
}

func (s *flowDistrictListerStub) ListDistricts(ctx context.Context) ([]DistrictListItem, error) {
	if s == nil {
		return nil, nil
	}
	return s.items, s.err
}

func (s *flowDistrictListerStub) ListDistrictsByCity(ctx context.Context, cityID int) ([]DistrictListItem, error) {
	if s == nil {
		return nil, nil
	}
	return s.items, s.err
}

type flowDistrictVariantPriceUpdaterStub struct {
	called bool
	params UpdateDistrictVariantPriceParams
	err    error
}

func (s *flowDistrictVariantPriceUpdaterStub) UpdateDistrictVariantPrice(
	ctx context.Context,
	params UpdateDistrictVariantPriceParams,
) error {
	s.called = true
	s.params = params
	return s.err
}

func openAdminCatalog(t *testing.T, svc *Service, key SessionKey) {
	t.Helper()

	_, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminCatalogOpen,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
}

type flowDistrictPlacementReaderStub struct {
	categories        []CategoryListItem
	products          []ProductListItem
	variants          []DistrictPlacementVariantListItem
	availableVariants []VariantListItem

	categoriesErr        error
	productsErr          error
	variantsErr          error
	availableVariantsErr error
}

func (s *flowDistrictPlacementReaderStub) ListDistrictCategories(
	ctx context.Context,
	districtID int,
) ([]CategoryListItem, error) {
	return s.categories, s.categoriesErr
}

func (s *flowDistrictPlacementReaderStub) ListDistrictProducts(
	ctx context.Context,
	districtID, categoryID int,
) ([]ProductListItem, error) {
	return s.products, s.productsErr
}

func (s *flowDistrictPlacementReaderStub) ListDistrictVariants(
	ctx context.Context,
	districtID, productID int,
) ([]DistrictPlacementVariantListItem, error) {
	return s.variants, s.variantsErr
}

func (s *flowDistrictPlacementReaderStub) ListAvailableVariantsForDistrictProduct(
	ctx context.Context,
	districtID, productID int,
) ([]VariantListItem, error) {
	if s == nil {
		return nil, nil
	}

	return s.availableVariants, s.availableVariantsErr
}

func openAdminDistrictVariantPriceUpdate(t *testing.T, svc *Service, key SessionKey) {
	t.Helper()

	openAdminCatalog(t, svc, key)

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminDistrictVariantPriceUpdateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
}

func TestHandleAction_AdminDistrictVariantPriceUpdateStart_OpensDistrictSelect(t *testing.T) {
	t.Parallel()

	svc := NewServiceWithDeps(nil, nil, ServiceDeps{
		DistrictLister: &flowDistrictListerStub{
			items: []DistrictListItem{
				{ID: 7, Code: "center", Label: "Центр"},
			},
		},
	})
	key := testSessionKey("shop-admin-price-update-start")

	openAdminCatalog(t, svc, key)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionAdminDistrictVariantPriceUpdateStart,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Изменение цены варианта\n\nВыберите район:", vm.Text)
	require.NotNil(t, vm.Inline)
}

func TestHandleAction_AdminDistrictVariantPriceUpdate_SelectDistrict_OpensCategorySelect(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, ServiceDeps{
		DistrictLister: &flowDistrictListerStub{
			items: []DistrictListItem{
				{ID: 7, Code: "center", Label: "Центр"},
			},
		},
		DistrictPlacements: &flowDistrictPlacementReaderStub{
			categories: []CategoryListItem{
				{ID: 11, Code: "flowers", Label: "Цветы"},
			},
		},
	})
	key := testSessionKey("shop-admin-price-update-district")

	openAdminDistrictVariantPriceUpdate(t, svc, key)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictVariantSelectDistrictAction(7),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Изменение цены варианта\n\nРайон: Центр\n\nВыберите категорию:", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictVariantPriceUpdateCategorySelect, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
	require.Equal(t, "7", session.Pending.Value(PendingValueDistrictID))
	require.Equal(t, "Центр", session.Pending.Value(PendingValueDistrictName))
}

func TestHandleAction_AdminDistrictVariantPriceUpdate_SelectCategory_OpensProductSelect(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, ServiceDeps{
		DistrictLister: &flowDistrictListerStub{
			items: []DistrictListItem{
				{ID: 7, Code: "center", Label: "Центр"},
			},
		},
		DistrictPlacements: &flowDistrictPlacementReaderStub{
			categories: []CategoryListItem{
				{ID: 11, Code: "flowers", Label: "Цветы"},
			},
			products: []ProductListItem{
				{ID: 21, Code: "rose-box", Label: "Rose Box"},
			},
		},
	})
	key := testSessionKey("shop-admin-price-update-category")

	openAdminDistrictVariantPriceUpdate(t, svc, key)

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictVariantSelectDistrictAction(7),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictVariantPriceUpdateSelectCategoryAction(11),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Изменение цены варианта\n\nРайон: Центр\n\nКатегория: Цветы\n\nВыберите товар:", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictVariantPriceUpdateProductSelect, session.Current)
	require.Equal(t, "11", session.Pending.Value(PendingValueCategoryID))
	require.Equal(t, "Цветы", session.Pending.Value(PendingValueCategoryName))
}

func TestHandleAction_AdminDistrictVariantPriceUpdate_SelectProduct_OpensVariantSelect(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, ServiceDeps{
		DistrictLister: &flowDistrictListerStub{
			items: []DistrictListItem{
				{ID: 7, Code: "center", Label: "Центр"},
			},
		},
		DistrictPlacements: &flowDistrictPlacementReaderStub{
			categories: []CategoryListItem{
				{ID: 11, Code: "flowers", Label: "Цветы"},
			},
			products: []ProductListItem{
				{ID: 21, Code: "rose-box", Label: "Rose Box"},
			},
			variants: []DistrictPlacementVariantListItem{
				{ID: 9, Code: "large", Label: "L / 25 шт", Price: 5900, PriceText: "5900 ₽"},
			},
		},
	})
	key := testSessionKey("shop-admin-price-update-product")

	openAdminDistrictVariantPriceUpdate(t, svc, key)

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictVariantSelectDistrictAction(7),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	_, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictVariantPriceUpdateSelectCategoryAction(11),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictVariantPriceUpdateSelectProductAction(21),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Изменение цены варианта\n\nРайон: Центр\n\nТовар: Rose Box\n\nВыберите вариант:", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictVariantPriceUpdateVariantSelect, session.Current)
	require.Equal(t, "21", session.Pending.Value(PendingValueProductID))
	require.Equal(t, "Rose Box", session.Pending.Value(PendingValueProductName))
	require.NotNil(t, vm.Inline)
	require.True(t, hasInlineActionLabel(vm, "Rose Box - L / 25 шт - 5900 ₽"))
}

func TestHandleAction_AdminDistrictVariantPriceUpdate_SelectProduct_OpensVariantSelect_UsesPriceFallback(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, ServiceDeps{
		DistrictLister: &flowDistrictListerStub{
			items: []DistrictListItem{
				{ID: 7, Code: "center", Label: "Центр"},
			},
		},
		DistrictPlacements: &flowDistrictPlacementReaderStub{
			categories: []CategoryListItem{
				{ID: 11, Code: "flowers", Label: "Цветы"},
			},
			products: []ProductListItem{
				{ID: 21, Code: "rose-box", Label: "Rose Box"},
			},
			variants: []DistrictPlacementVariantListItem{
				{ID: 9, Code: "large", Label: "L / 25 шт", Price: 5900},
			},
		},
	})
	key := testSessionKey("shop-admin-price-update-product-fallback")

	openAdminDistrictVariantPriceUpdate(t, svc, key)

	for _, action := range []ActionID{
		adminDistrictVariantSelectDistrictAction(7),
		adminDistrictVariantPriceUpdateSelectCategoryAction(11),
	} {
		_, err := svc.HandleAction(context.Background(), ActionRequest{
			BotID:         "shop-admin",
			BotName:       "Admin Shop",
			StartScenario: string(StartScenarioInlineCatalog),
			ActionID:      action,
			SessionKey:    key,
			CanAdmin:      true,
		})
		require.NoError(t, err)
	}

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictVariantPriceUpdateSelectProductAction(21),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.True(t, hasInlineActionLabel(vm, "Rose Box - L / 25 шт - 5900 ₽"))
}

func TestHandleAction_AdminDistrictVariantPriceUpdate_SelectVariant_OpensPriceInput(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	svc := NewServiceWithDeps(store, nil, ServiceDeps{
		DistrictLister: &flowDistrictListerStub{
			items: []DistrictListItem{
				{ID: 7, Code: "center", Label: "Центр"},
			},
		},
		DistrictPlacements: &flowDistrictPlacementReaderStub{
			categories: []CategoryListItem{
				{ID: 11, Code: "flowers", Label: "Цветы"},
			},
			products: []ProductListItem{
				{ID: 21, Code: "rose-box", Label: "Rose Box"},
			},
			variants: []DistrictPlacementVariantListItem{
				{ID: 9, Code: "large", Label: "L / 25 шт", Price: 5900, PriceText: "5900 ₽"},
			},
		},
	})
	key := testSessionKey("shop-admin-price-update-variant")

	openAdminDistrictVariantPriceUpdate(t, svc, key)

	for _, action := range []ActionID{
		adminDistrictVariantSelectDistrictAction(7),
		adminDistrictVariantPriceUpdateSelectCategoryAction(11),
		adminDistrictVariantPriceUpdateSelectProductAction(21),
	} {
		_, err := svc.HandleAction(context.Background(), ActionRequest{
			BotID:         "shop-admin",
			BotName:       "Admin Shop",
			StartScenario: string(StartScenarioInlineCatalog),
			ActionID:      action,
			SessionKey:    key,
			CanAdmin:      true,
		})
		require.NoError(t, err)
	}

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      adminDistrictVariantSelectVariantAction(9),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Изменение цены варианта\n\nРайон: Центр\n\nВариант: Rose Box - L / 25 шт\n\nТекущая цена: 5900 ₽\n\nВведите новую цену сообщением.", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictVariantPriceUpdatePrice, session.Current)
	require.Equal(t, PendingInputDistrictVariantPriceUpdate, session.Pending.Kind)
}

func TestHandleText_AdminDistrictVariantPriceUpdate_InvalidPrice(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	updater := &flowDistrictVariantPriceUpdaterStub{}
	svc := NewServiceWithDeps(store, nil, ServiceDeps{
		DistrictLister: &flowDistrictListerStub{
			items: []DistrictListItem{
				{ID: 7, Code: "center", Label: "Центр"},
			},
		},
		DistrictPlacements: &flowDistrictPlacementReaderStub{
			categories: []CategoryListItem{
				{ID: 11, Code: "flowers", Label: "Цветы"},
			},
			products: []ProductListItem{
				{ID: 21, Code: "rose-box", Label: "Rose Box"},
			},
			variants: []DistrictPlacementVariantListItem{
				{ID: 9, Code: "large", Label: "L / 25 шт", Price: 5900, PriceText: "5900 ₽"},
			},
		},
		DistrictVariantPrices: updater,
	})
	key := testSessionKey("shop-admin-price-update-invalid")

	openAdminDistrictVariantPriceUpdate(t, svc, key)

	for _, action := range []ActionID{
		adminDistrictVariantSelectDistrictAction(7),
		adminDistrictVariantPriceUpdateSelectCategoryAction(11),
		adminDistrictVariantPriceUpdateSelectProductAction(21),
		adminDistrictVariantSelectVariantAction(9),
	} {
		_, err := svc.HandleAction(context.Background(), ActionRequest{
			BotID:         "shop-admin",
			BotName:       "Admin Shop",
			StartScenario: string(StartScenarioInlineCatalog),
			ActionID:      action,
			SessionKey:    key,
			CanAdmin:      true,
		})
		require.NoError(t, err)
	}

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "abc",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(
		t,
		"Изменение цены варианта\n\nРайон: Центр\n\nВариант: Rose Box - L / 25 шт\n\nТекущая цена: 5900 ₽\n\nЦена должна быть положительным числом.\n\nВведите новую цену сообщением.",
		vm.Text,
	)

	require.False(t, updater.called)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictVariantPriceUpdatePrice, session.Current)
	require.Equal(t, PendingInputDistrictVariantPriceUpdate, session.Pending.Kind)
}

func TestHandleAction_AdminDistrictVariantPriceUpdate_BackFromDone_ReturnsAdminCatalog(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	updater := &flowDistrictVariantPriceUpdaterStub{}
	svc := NewServiceWithDeps(store, nil, ServiceDeps{
		DistrictLister: &flowDistrictListerStub{
			items: []DistrictListItem{
				{ID: 7, Code: "center", Label: "Центр"},
			},
		},
		DistrictPlacements: &flowDistrictPlacementReaderStub{
			categories: []CategoryListItem{
				{ID: 11, Code: "flowers", Label: "Цветы"},
			},
			products: []ProductListItem{
				{ID: 21, Code: "rose-box", Label: "Rose Box"},
			},
			variants: []DistrictPlacementVariantListItem{
				{ID: 9, Code: "large", Label: "L / 25 шт", Price: 5900, PriceText: "5900 ₽"},
			},
		},
		DistrictVariantPrices: updater,
	})
	key := testSessionKey("shop-admin-price-update-back-done")

	openAdminDistrictVariantPriceUpdate(t, svc, key)

	for _, action := range []ActionID{
		adminDistrictVariantSelectDistrictAction(7),
		adminDistrictVariantPriceUpdateSelectCategoryAction(11),
		adminDistrictVariantPriceUpdateSelectProductAction(21),
		adminDistrictVariantSelectVariantAction(9),
	} {
		_, err := svc.HandleAction(context.Background(), ActionRequest{
			BotID:         "shop-admin",
			BotName:       "Admin Shop",
			StartScenario: string(StartScenarioInlineCatalog),
			ActionID:      action,
			SessionKey:    key,
			CanAdmin:      true,
		})
		require.NoError(t, err)
	}

	_, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "6100",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
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

func TestHandleAction_AdminDistrictVariantPriceUpdate_BackChainAfterDone_ReturnsAdminRoot(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	updater := &flowDistrictVariantPriceUpdaterStub{}
	svc := NewServiceWithDeps(store, nil, ServiceDeps{
		DistrictLister: &flowDistrictListerStub{
			items: []DistrictListItem{
				{ID: 7, Code: "center", Label: "Центр"},
			},
		},
		DistrictPlacements: &flowDistrictPlacementReaderStub{
			categories: []CategoryListItem{
				{ID: 11, Code: "flowers", Label: "Цветы"},
			},
			products: []ProductListItem{
				{ID: 21, Code: "rose-box", Label: "Rose Box"},
			},
			variants: []DistrictPlacementVariantListItem{
				{ID: 9, Code: "large", Label: "L / 25 шт", Price: 5900, PriceText: "5900 ₽"},
			},
		},
		DistrictVariantPrices: updater,
	})
	key := testSessionKey("shop-admin-price-update-back-root")

	openAdminDistrictVariantPriceUpdate(t, svc, key)

	for _, action := range []ActionID{
		adminDistrictVariantSelectDistrictAction(7),
		adminDistrictVariantPriceUpdateSelectCategoryAction(11),
		adminDistrictVariantPriceUpdateSelectProductAction(21),
		adminDistrictVariantSelectVariantAction(9),
	} {
		_, err := svc.HandleAction(context.Background(), ActionRequest{
			BotID:         "shop-admin",
			BotName:       "Admin Shop",
			StartScenario: string(StartScenarioInlineCatalog),
			ActionID:      action,
			SessionKey:    key,
			CanAdmin:      true,
		})
		require.NoError(t, err)
	}

	_, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "6100",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionBack,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Админка · Каталог\n\nВыберите действие:", vm.Text)

	vm, err = svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionBack,
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Админка\n\nВыберите раздел:", vm.Text)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminRoot, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}

func TestHandleText_AdminDistrictVariantPriceUpdate_Success(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	updater := &flowDistrictVariantPriceUpdaterStub{}
	svc := NewServiceWithDeps(store, nil, ServiceDeps{
		DistrictLister: &flowDistrictListerStub{
			items: []DistrictListItem{
				{ID: 7, Code: "center", Label: "Центр"},
			},
		},
		DistrictPlacements: &flowDistrictPlacementReaderStub{
			categories: []CategoryListItem{
				{ID: 11, Code: "flowers", Label: "Цветы"},
			},
			products: []ProductListItem{
				{ID: 21, Code: "rose-box", Label: "Rose Box"},
			},
			variants: []DistrictPlacementVariantListItem{
				{ID: 9, Code: "large", Label: "L / 25 шт", Price: 5900, PriceText: "5900 ₽"},
			},
		},
		DistrictVariantPrices: updater,
	})
	key := testSessionKey("shop-admin-price-update-success")

	openAdminDistrictVariantPriceUpdate(t, svc, key)

	for _, action := range []ActionID{
		adminDistrictVariantSelectDistrictAction(7),
		adminDistrictVariantPriceUpdateSelectCategoryAction(11),
		adminDistrictVariantPriceUpdateSelectProductAction(21),
		adminDistrictVariantSelectVariantAction(9),
	} {
		_, err := svc.HandleAction(context.Background(), ActionRequest{
			BotID:         "shop-admin",
			BotName:       "Admin Shop",
			StartScenario: string(StartScenarioInlineCatalog),
			ActionID:      action,
			SessionKey:    key,
			CanAdmin:      true,
		})
		require.NoError(t, err)
	}

	vm, err := svc.HandleText(context.Background(), TextRequest{
		BotID:         "shop-admin",
		BotName:       "Admin Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "6100",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)
	require.Equal(t, "Изменение цены варианта\n\nЦена варианта обновлена.", vm.Text)

	require.True(t, updater.called)
	require.Equal(t, 7, updater.params.DistrictID)
	require.Equal(t, 9, updater.params.VariantID)
	require.Equal(t, 6100, updater.params.Price)

	session, ok := store.Get(key)
	require.True(t, ok)
	require.Equal(t, ScreenAdminDistrictVariantPriceUpdateDone, session.Current)
	require.Equal(t, PendingInputNone, session.Pending.Kind)
}
