package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

func TestBuildAdminCustomerMessageNotificationView(t *testing.T) {
	t.Parallel()

	vm := buildAdminCustomerMessageNotificationView(
		manager.BotSpec{
			ID:               "bot-1",
			Name:             "shop-main-bot",
			TelegramUsername: "Koha90_bot",
		},
		customerMessageNotification{
			UserID:       123,
			ChatID:       456,
			UserName:     "Алексей",
			UserUsername: "koha90",
			Text:         "Здравствуйте, есть доставка сегодня?",
		},
	)

	assert.Contains(t, vm.Text, "Новое сообщение")
	assert.Contains(t, vm.Text, "Бот: @Koha90_bot")
	assert.Contains(t, vm.Text, "Пользователь: Алексей")
	assert.Contains(t, vm.Text, "Логин: @koha90")
	assert.Contains(t, vm.Text, "User ID: 123")
	assert.Contains(t, vm.Text, "Chat ID: 456")
	assert.Contains(t, vm.Text, "Здравствуйте, есть доставка сегодня?")
	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 1)
	require.Len(t, vm.Inline.Sections[0].Actions, 1)

	action := vm.Inline.Sections[0].Actions[0]
	assert.Equal(t, "Ответить", action.Label)
	assert.Equal(t, flow.AdminCustomerReplyStartAction(456, 123), action.ID)
}

func TestBuildAdminCustomerMessageNotificationView_WithoutOptionalUserLabels(t *testing.T) {
	t.Parallel()

	vm := buildAdminCustomerMessageNotificationView(
		manager.BotSpec{Name: "shop-main-bot"},
		customerMessageNotification{
			UserID: 123,
			ChatID: 456,
			Text:   "Нужна помощь",
		},
	)

	assert.Contains(t, vm.Text, "Бот: shop-main-bot")
	assert.NotContains(t, vm.Text, "Пользователь:")
	assert.NotContains(t, vm.Text, "Логин:")
	assert.Contains(t, vm.Text, "Нужна помощь")
	require.NotNil(t, vm.Inline)
}
