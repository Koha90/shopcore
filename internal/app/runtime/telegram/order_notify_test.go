package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

// TestBuildAdminOrderNotificationView verifies admin notification content.
func TestBuildAdminOrderNotificationView(t *testing.T) {
	t.Parallel()

	vm := buildAdminOrderNotificationView(
		manager.BotSpec{
			ID:               "bot-1",
			Name:             "shop-main-bot",
			TelegramUsername: "Koha90_bot",
		},
		OrderNotificationMeta{
			BotUsername: "shop_main_bot",
			UserID:      123,
			ChatID:      456,
			UserName:    "Алексей",
			UserLogin:   "koha90",
		},
		flow.OrderContext{
			CityName:      "Пермь",
			DistrictName:  "Мотовилихинский",
			ProductLabel:  "«Мишки в лесу 🌳» Шишкин",
			VariantLabel:  "Мишки 🧸",
			BasePriceText: "3000 ₽",
		},
	)

	assert.Contains(t, vm.Text, "Новый заказ")
	assert.Contains(t, vm.Text, "Бот: @Koha90_bot")
	assert.NotContains(t, vm.Text, "Бот: shop-main-bot")
	assert.Contains(t, vm.Text, "Пермь")
	assert.Contains(t, vm.Text, "Мотовилихинский")
	assert.Contains(t, vm.Text, "3000 ₽")
	assert.Contains(t, vm.Text, "123")
	assert.Contains(t, vm.Text, "456")
}

func TestFormatBotLabel_FallbackToTelegramBotName(t *testing.T) {
	t.Parallel()

	got := formatBotLabel(manager.BotSpec{
		Name:            "shop-main-bot",
		TelegramBotName: "Koha90 Bot",
	})

	assert.Equal(t, "Koha90 Bot", got)
}

func TestFormatBotLabel_FallbackToSpecName(t *testing.T) {
	t.Parallel()

	got := formatBotLabel(manager.BotSpec{
		Name: "shop-main-bot",
	})

	assert.Equal(t, "shop-main-bot", got)
}

func TestFormatBotLabel_PrefersTelegramUsername(t *testing.T) {
	t.Parallel()

	got := formatBotLabel(manager.BotSpec{
		Name:             "shop-main-bot",
		TelegramBotName:  "Koha90 Bot",
		TelegramUsername: "Koha90_bot",
	})

	assert.Equal(t, "@Koha90_bot", got)
}
