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
			ID:   "bot-1",
			Name: "shop-main-bot",
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
	assert.Contains(t, vm.Text, "shop-main-bot")
	assert.Contains(t, vm.Text, "Пермь")
	assert.Contains(t, vm.Text, "Мотовилихинский")
	assert.Contains(t, vm.Text, "3000 ₽")
	assert.Contains(t, vm.Text, "123")
	assert.Contains(t, vm.Text, "456")
}
