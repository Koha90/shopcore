package telegram

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/manager"
)

func TestSpecAdminAccessResolver_DenyWhenListIsEmpty(t *testing.T) {
	t.Parallel()

	resolver := SpecAdminAccessResolver{}

	require.False(t, resolver.CanAdminTelegram(manager.BotSpec{}, 311485249))
}

func TestSpecAdminAccessResolver_AllowConfiguredUser(t *testing.T) {
	t.Parallel()

	resolver := SpecAdminAccessResolver{}
	spec := manager.BotSpec{
		ID:                   "shop-main",
		TelegramAdminUserIDs: []int64{311485249},
	}

	require.True(t, resolver.CanAdminTelegram(spec, 311485249))
	require.False(t, resolver.CanAdminTelegram(spec, 999999999))
}
