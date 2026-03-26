package flow

import (
	"context"
	"testing"
)

func TestNormalizeStartScenario_Default(t *testing.T) {
	t.Parallel()

	got := NormalizeStartScenario("")
	if got != StartScenarioReplyWelcome {
		t.Fatalf("expected default scenario %q, got %q", StartScenarioReplyWelcome, got)
	}
}

func TestStart_ReplyWelcome(t *testing.T) {
	t.Parallel()

	svc := NewService()

	vm, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-main",
		BotName:       "Main Shop",
		StartScenario: string(StartScenarioReplyWelcome),
	})
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}

	if vm.Reply == nil {
		t.Fatal("expected reply keyboard")
	}
	if vm.Inline != nil {
		t.Fatal("did not expected inline keyboard")
	}
	if len(vm.Reply.Rows) != 2 {
		t.Fatalf("expected 2 reply rows, got %d", len(vm.Reply.Rows))
	}
}

func TestStart_InlineCatalog(t *testing.T) {
	t.Parallel()

	svc := NewService()

	vm, err := svc.Start(context.Background(), StartRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
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
	if len(vm.Inline.Sections) != 2 {
		t.Fatalf("expected 2 inline Sections, got %d", len(vm.Inline.Sections))
	}
	if vm.Inline.Sections[0].Columns != 2 {
		t.Fatalf("expected first section columns=2, got %d", vm.Inline.Sections[0].Columns)
	}
	if vm.Inline.Sections[1].Columns != 1 {
		t.Fatalf("expected first section columns=2, got %d", vm.Inline.Sections[1].Columns)
	}
}

func TestHandleAction_BackReturnsStartScenario(t *testing.T) {
	t.Parallel()

	svc := NewService()

	vm, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-inline",
		BotName:       "Inline Shop",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      ActionBack,
	})
	if err != nil {
		t.Fatalf("HandleAction returned error: %v", err)
	}

	if vm.Inline == nil {
		t.Fatalf("expected catalog start text, go %q", vm.Text)
	}

	if vm.Text != "Каталог" {
		t.Fatalf("expected catalog start text, got %q", vm.Text)
	}
}

func TestResolveReplyAction(t *testing.T) {
	t.Parallel()

	svc := NewService()

	tests := []struct {
		name string
		text string
		want ActionID
		ok   bool
	}{
		{name: "catalog", text: "Каталог", want: ActionCatalogStart, ok: true},
		{name: "cabinet", text: "Мой кабинет", want: ActionCabinetOpen, ok: true},
		{name: "support", text: "Поддержка", want: ActionSupportOpen, ok: true},
		{name: "reviews", text: "Отзывы", want: ActionReviewsOpen, ok: true},
		{name: "unknown", text: "Что-то ещё", want: "", ok: false},

		{name: "catalog", text: "♻️ Каталог", want: ActionCatalogStart, ok: true},
		{name: "cabinet", text: "⚙️ Мой кабинет", want: ActionCabinetOpen, ok: true},
		{name: "support", text: "🤷‍♂️ Поддержка", want: ActionSupportOpen, ok: true},
		{name: "reviews", text: "📨 Отзывы", want: ActionReviewsOpen, ok: true},
	}

	for _, tt := range tests {
		// tt := tt
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

	svc := NewService()

	_, err := svc.HandleAction(context.Background(), ActionRequest{
		BotID:         "shop-main",
		BotName:       "Main Shop",
		StartScenario: string(StartScenarioReplyWelcome),
		ActionID:      ActionID("unknown:action"),
	})
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
	if err != ErrUnknownAction {
		t.Fatalf("expected ErrUnknownAction, got %v", err)
	}
}
