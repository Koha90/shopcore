package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
		CustomerMessageNotification{
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
	assert.Nil(t, vm.Inline)
}

func TestBuildAdminCustomerMessageNotificationView_WithoutOptionalUserLabels(t *testing.T) {
	t.Parallel()

	vm := buildAdminCustomerMessageNotificationView(
		manager.BotSpec{Name: "shop-main-bot"},
		CustomerMessageNotification{
			UserID: 123,
			ChatID: 456,
			Text:   "Нужна помощь",
		},
	)

	assert.Contains(t, vm.Text, "Бот: shop-main-bot")
	assert.NotContains(t, vm.Text, "Пользователь:")
	assert.NotContains(t, vm.Text, "Логин:")
	assert.Contains(t, vm.Text, "Нужна помощь")
}
