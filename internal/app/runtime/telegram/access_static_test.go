package telegram

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStaticAdminAccessResolver_DenyWhenAllowlistIsNil(t *testing.T) {
	t.Parallel()

	resolver := StaticAdminAccessResolver{}

	require.False(t, resolver.CanAdminTelegram("shop-main", 311485249))
}

func TestStaticAdminAccessResolver_DenyWhenBotIsNotConfigured(t *testing.T) {
	t.Parallel()

	resolver := StaticAdminAccessResolver{
		Allow: map[string]map[int64]struct{}{
			"shop-main": {
				311485249: {},
			},
		},
	}

	require.False(t, resolver.CanAdminTelegram("shop-other", 311485249))
}

func TestStaticAdminAccessResolver_AllowOnlyConfiguredUserForBot(t *testing.T) {
	t.Parallel()

	resolver := StaticAdminAccessResolver{
		Allow: map[string]map[int64]struct{}{
			"shop-main": {
				311485249: {},
			},
		},
	}

	require.True(t, resolver.CanAdminTelegram("shop-main", 311485249))
	require.False(t, resolver.CanAdminTelegram("shop-main", 999999999))
}
