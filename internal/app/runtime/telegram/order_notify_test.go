package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

func TestBuildAdminOrderNotificationView(t *testing.T) {
	t.Parallel()

	vm := buildAdminOrderNotificationView(
		manager.BotSpec{
			ID:               "bot-1",
			Name:             "shop-main-bot",
			TelegramUsername: "Koha90_bot",
		},
		ordersvc.Order{
			ID:           42,
			BotID:        "bot-1",
			BotName:      "shop-main-bot",
			UserID:       123,
			ChatID:       456,
			UserName:     "Алексей",
			UserUsername: "koha90",
			CityName:     "Пермь",
			DistrictName: "Мотовилихинский",
			ProductName:  "«Мишки в лесу 🌳» Шишкин",
			VariantName:  "Мишки 🧸",
			PriceText:    "3000 ₽",
			Status:       ordersvc.OrderStatusNew,
		},
	)

	assert.Contains(t, vm.Text, "Новый заказ")
	assert.Contains(t, vm.Text, "Заказ: #42")
	assert.Contains(t, vm.Text, "Статус: new")
	assert.Contains(t, vm.Text, "Бот: @Koha90_bot")
	assert.Contains(t, vm.Text, "Пермь")
	assert.Contains(t, vm.Text, "Мотовилихинский")
	assert.Contains(t, vm.Text, "3000 ₽")
	assert.Contains(t, vm.Text, "123")
	assert.Contains(t, vm.Text, "456")

	if assert.NotNil(t, vm.Inline) {
		assert.Len(t, vm.Inline.Sections, 1)
		assert.Len(t, vm.Inline.Sections[0].Actions, 3)
		assert.Equal(t, "Взять в работу", vm.Inline.Sections[0].Actions[0].Label)
		assert.Equal(t, "Закрыть", vm.Inline.Sections[0].Actions[1].Label)
		assert.Equal(t, "Ответить клиенту", vm.Inline.Sections[0].Actions[2].Label)
		assert.Equal(
			t,
			flow.AdminCustomerReplyStartAction(456, 123),
			vm.Inline.Sections[0].Actions[2].ID,
		)
	}
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

func TestBuildAdminOrderNotificationView_InProgress(t *testing.T) {
	t.Parallel()

	vm := buildAdminOrderNotificationView(
		manager.BotSpec{
			Name:             "shop-main",
			TelegramUsername: "Koha90_bot",
		},
		ordersvc.Order{
			ID:           42,
			UserID:       123,
			ChatID:       456,
			UserName:     "Алексей",
			UserUsername: "koha90",
			CityName:     "Пермь",
			DistrictName: "Мотовилихинский",
			ProductName:  "Rose Box",
			VariantName:  "L / 25 шт",
			PriceText:    "5900 ₽",
			Status:       ordersvc.OrderStatusInProgress,
		},
	)

	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 1)
	require.Len(t, vm.Inline.Sections[0].Actions, 2)
	require.Equal(t, "Закрыть", vm.Inline.Sections[0].Actions[0].Label)
	require.Equal(t, "Ответить клиенту", vm.Inline.Sections[0].Actions[1].Label)
}

func TestBuildAdminOrderNotificationView_Closed(t *testing.T) {
	t.Parallel()

	vm := buildAdminOrderNotificationView(
		manager.BotSpec{Name: "shop-main"},
		ordersvc.Order{
			ID:           42,
			UserID:       123,
			ChatID:       456,
			UserName:     "Алексей",
			CityName:     "Пермь",
			DistrictName: "Мотовилихинский",
			ProductName:  "Rose Box",
			VariantName:  "L / 25 шт",
			PriceText:    "5900 ₽",
			Status:       ordersvc.OrderStatusClosed,
		},
	)

	require.NotNil(t, vm.Inline)
	require.Len(t, vm.Inline.Sections, 1)
	require.Len(t, vm.Inline.Sections[0].Actions, 1)
	require.Equal(t, "Ответить клиенту", vm.Inline.Sections[0].Actions[0].Label)
	require.Contains(t, vm.Text, "Статус: closed")
}
