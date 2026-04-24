package flow

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// TestServiceHandleAction_OrderStartFromVariantLeaf verifies that order flow
// starts only after full catalog drill-down to selected variant leaf.
func TestServiceHandleAction_OrderStartFromVariantLeaf(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewMemoryStore()
	svc := NewService(store)

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 100,
		UserID: 200,
	}

	_, err := svc.Start(ctx, StartRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
	})
	if err != nil {
		t.Fatalf("start flow: %v", err)
	}

	selectDemoVariantLeaf(t, ctx, svc, key)

	vm, err := svc.HandleAction(ctx, ActionRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		ActionID:      ActionOrderStart,
	})
	if err != nil {
		t.Fatalf("start order: %v", err)
	}

	session, ok := store.Get(key)
	if !ok {
		t.Fatal("session not found after order start")
	}
	if session.Current != ScreenOrderConfirm {
		t.Fatalf("current screen = %q, want %q", session.Current, ScreenOrderConfirm)
	}

	assertTextContains(t, vm.Text, "Москва")
	assertTextContains(t, vm.Text, "Центр")
	assertTextContains(t, vm.Text, "Rose Box")
	assertTextContains(t, vm.Text, "S / 9 шт")
	assertTextContains(t, vm.Text, "2500 ₽")

	if !viewHasAction(vm, ActionOrderConfirm) {
		t.Fatal("confirm action not found in order confirm screen")
	}
	if !viewHasAction(vm, ActionBack) {
		t.Fatal("back action not found in order confirm screen")
	}
}

// TestServiceHandleAction_OrderStartOutsideLeafReturnsError verifies that
// order flow cannot start before reaching variant leaf.
func TestServiceHandleAction_OrderStartOutsideLeafReturnsError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewMemoryStore()
	svc := NewService(store)

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 101,
		UserID: 201,
	}

	_, err := svc.Start(ctx, StartRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
	})
	if err != nil {
		t.Fatalf("start flow: %v", err)
	}

	_, err = svc.HandleAction(ctx, ActionRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		ActionID:      catalogSelectAction(LevelCity, "moscow"),
	})
	if err != nil {
		t.Fatalf("select city: %v", err)
	}

	_, err = svc.HandleAction(ctx, ActionRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		ActionID:      ActionOrderStart,
	})
	if !errors.Is(err, ErrUnknownAction) {
		t.Fatalf("error = %v, want %v", err, ErrUnknownAction)
	}
}

// TestServiceHandleAction_OrderConfirm verifies terminal transition from
// confirmation screen to done screen.
func TestServiceHandleAction_OrderConfirm(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewMemoryStore()
	svc := NewService(store)

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 102,
		UserID: 202,
	}

	_, err := svc.Start(ctx, StartRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
	})
	if err != nil {
		t.Fatalf("start flow: %v", err)
	}

	selectDemoVariantLeaf(t, ctx, svc, key)

	_, err = svc.HandleAction(ctx, ActionRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		ActionID:      ActionOrderStart,
	})
	if err != nil {
		t.Fatalf("start order: %v", err)
	}

	vm, err := svc.HandleAction(ctx, ActionRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		ActionID:      ActionOrderConfirm,
	})
	if err != nil {
		t.Fatalf("confirm order: %v", err)
	}

	session, ok := store.Get(key)
	if !ok {
		t.Fatal("session not found after order confirm")
	}
	if session.Current != ScreenOrderDone {
		t.Fatalf("current screen = %q, want %q", session.Current, ScreenOrderDone)
	}

	assertTextContains(t, vm.Text, "Заявка принята")

	if !viewHasAction(vm, ActionCatalogStart) {
		t.Fatal("root action not found in order done screen")
	}
}

// TestServiceHandleAction_BackFromOrderConfirmReturnsToVariantLeaf verifies
// history-based back behavior for order flow.
func TestServiceHandleAction_BackFromOrderConfirmReturnsToVariantLeaf(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewMemoryStore()
	svc := NewService(store)

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 103,
		UserID: 203,
	}

	_, err := svc.Start(ctx, StartRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
	})
	if err != nil {
		t.Fatalf("start flow: %v", err)
	}

	selectDemoVariantLeaf(t, ctx, svc, key)

	_, err = svc.HandleAction(ctx, ActionRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		ActionID:      ActionOrderStart,
	})
	if err != nil {
		t.Fatalf("start order: %v", err)
	}

	_, err = svc.HandleAction(ctx, ActionRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		ActionID:      ActionBack,
	})
	if err != nil {
		t.Fatalf("back from order confirm: %v", err)
	}

	session, ok := store.Get(key)
	if !ok {
		t.Fatal("session not found after back")
	}

	want := catalogScreen(CatalogPath{
		{Level: LevelCity, ID: "moscow"},
		{Level: LevelCategory, ID: "flowers"},
		{Level: LevelDistrict, ID: "center"},
		{Level: LevelProduct, ID: "rose-box"},
		{Level: LevelVariant, ID: "small"},
	})

	if session.Current != want {
		t.Fatalf("current screen = %q, want %q", session.Current, want)
	}
}

// selectDemoVariantLeaf walks demo catalog until selected variant leaf.
//
// This helper keeps tests short and makes flow path explicit.
func selectDemoVariantLeaf(
	t *testing.T,
	ctx context.Context,
	svc *Service,
	key SessionKey,
) {
	t.Helper()

	actions := []ActionID{
		catalogSelectAction(LevelCity, "moscow"),
		catalogSelectAction(LevelCategory, "flowers"),
		catalogSelectAction(LevelDistrict, "center"),
		catalogSelectAction(LevelProduct, "rose-box"),
		catalogSelectAction(LevelVariant, "small"),
	}

	for _, actionID := range actions {
		_, err := svc.HandleAction(ctx, ActionRequest{
			SessionKey:    key,
			StartScenario: "inline_catalog",
			ActionID:      actionID,
		})
		if err != nil {
			t.Fatalf("apply action %q: %v", actionID, err)
		}
	}
}

// viewHasAction reports whether inline keyboard contains provided action.
func viewHasAction(vm ViewModel, want ActionID) bool {
	if vm.Inline == nil {
		return false
	}

	for _, section := range vm.Inline.Sections {
		for _, action := range section.Actions {
			if action.ID == want {
				return true
			}
		}
	}

	return false
}

// assertTextContains verifies that response text includes expected fragment.
func assertTextContains(t *testing.T, text, want string) {
	t.Helper()

	if !strings.Contains(text, want) {
		t.Fatalf("text %q does not contain %q", text, want)
	}
}

// TestServiceHandleAction_OrderDoneMainMenuKeepsExtendedRootForInlineScenario verifies
// that done-screen return action uses scenario-aware catalog root instead of hardcoded compact root.
func TestServiceHandleAction_OrderDoneMainMenuKeepsExtendedRootForInlineScenario(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewMemoryStore()
	svc := NewService(store)

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 104,
		UserID: 204,
	}

	_, err := svc.Start(ctx, StartRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		CanAdmin:      true,
	})
	if err != nil {
		t.Fatalf("start flow: %v", err)
	}

	selectDemoVariantLeaf(t, ctx, svc, key)

	_, err = svc.HandleAction(ctx, ActionRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		CanAdmin:      true,
		ActionID:      ActionOrderStart,
	})
	if err != nil {
		t.Fatalf("start order: %v", err)
	}

	_, err = svc.HandleAction(ctx, ActionRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		CanAdmin:      true,
		ActionID:      ActionOrderConfirm,
	})
	if err != nil {
		t.Fatalf("confirm order: %v", err)
	}

	vm, err := svc.HandleAction(ctx, ActionRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
		CanAdmin:      true,
		ActionID:      ActionCatalogStart,
	})
	if err != nil {
		t.Fatalf("return to main menu: %v", err)
	}

	session, ok := store.Get(key)
	if !ok {
		t.Fatal("session not found after return to main menu")
	}

	if session.Current != ScreenRootExtended {
		t.Fatalf("current screen = %q, want %q", session.Current, ScreenRootExtended)
	}

	if !viewHasAction(vm, ActionAdminOpen) {
		t.Fatal("admin action not found on extended root after returning from order done screen")
	}
}
