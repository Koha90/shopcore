package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminCustomerReplyStartAction(t *testing.T) {
	t.Parallel()

	actionID := AdminCustomerReplyStartAction(456, 123)

	chatID, userID, ok := parseAdminCustomerReplyStartAction(actionID)

	require.True(t, ok)
	assert.Equal(t, int64(456), chatID)
	assert.Equal(t, int64(123), userID)
}

func TestHandleActionAdminCustomerReplyStartSetsPendingInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewMemoryStore()
	svc := NewServiceWithCatalogProvider(store, NewStaticCatalogProvider(DemoCatalog()))

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 1000,
		UserID: 2000,
	}

	_, err := svc.Start(ctx, StartRequest{
		BotID:         "bot-1",
		BotName:       "shop-main",
		StartScenario: string(StartScenarioInlineCatalog),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	vm, err := svc.HandleAction(ctx, ActionRequest{
		BotID:         "bot-1",
		BotName:       "shop-main",
		StartScenario: string(StartScenarioInlineCatalog),
		ActionID:      AdminCustomerReplyStartAction(456, 123),
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	assert.Contains(t, vm.Text, "Ответ пользователю")
	assert.Contains(t, vm.Text, "Введите ответ")

	session, ok := store.Get(key)
	require.True(t, ok)

	assert.Equal(t, ScreenAdminCustomerReply, session.Current)
	assert.Equal(t, PendingInputAdminCustomerReply, session.Pending.Kind)
	assert.Equal(t, "456", session.Pending.Value(PendingValueCustomerChatID))
	assert.Equal(t, "123", session.Pending.Value(PendingValueCustomerUserID))
}

func TestHandleTextAdminCustomerReplyReturnsSendTextEffect(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewMemoryStore()
	svc := NewServiceWithCatalogProvider(store, NewStaticCatalogProvider(DemoCatalog()))

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 1000,
		UserID: 2000,
	}

	store.Put(key, Session{
		Current: ScreenAdminCustomerReply,
		Pending: PendingInput{
			Kind: PendingInputAdminCustomerReply,
			Payload: PendingInputPayload{
				PendingValueCustomerChatID: "456",
				PendingValueCustomerUserID: "123",
			},
		},
		CanAdmin: true,
	})

	vm, err := svc.HandleText(ctx, TextRequest{
		BotID:         "bot-1",
		BotName:       "shop-main",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "Здравствуйте, доставка сегодня есть.",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	assert.Contains(t, vm.Text, "Ответ отправлен")
	require.Len(t, vm.Effects, 1)

	effect := vm.Effects[0]
	assert.Equal(t, EffectSendText, effect.Kind)
	assert.Equal(t, int64(456), effect.Target.ChatID)
	assert.Equal(t, int64(123), effect.Target.UserID)
	assert.Equal(t, "Здравствуйте, доставка сегодня есть.", effect.Text)

	session, ok := store.Get(key)
	require.True(t, ok)

	assert.Equal(t, ScreenAdminCustomerReplyDone, session.Current)
	assert.False(t, session.Pending.Active())
}

func TestHandleTextAdminCustomerReplyRejectsEmptyText(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewMemoryStore()
	svc := NewServiceWithCatalogProvider(store, NewStaticCatalogProvider(DemoCatalog()))

	key := SessionKey{
		BotID:  "bot-1",
		ChatID: 1000,
		UserID: 2000,
	}

	store.Put(key, Session{
		Current: ScreenAdminCustomerReply,
		Pending: PendingInput{
			Kind: PendingInputAdminCustomerReply,
			Payload: PendingInputPayload{
				PendingValueCustomerChatID: "456",
				PendingValueCustomerUserID: "123",
			},
		},
		CanAdmin: true,
	})

	vm, err := svc.HandleText(ctx, TextRequest{
		BotID:         "bot-1",
		BotName:       "shop-main",
		StartScenario: string(StartScenarioInlineCatalog),
		Text:          "   ",
		SessionKey:    key,
		CanAdmin:      true,
	})
	require.NoError(t, err)

	assert.Contains(t, vm.Text, "Ответ не может быть пустым")

	session, ok := store.Get(key)
	require.True(t, ok)

	assert.Equal(t, PendingInputAdminCustomerReply, session.Pending.Kind)
}
