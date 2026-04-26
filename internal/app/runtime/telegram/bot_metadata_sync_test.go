package telegram

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/manager"
)

type botMetadataUpdaterStub struct {
	called   bool
	id       string
	botID    int64
	username string
	name     string
	err      error
}

func (s *botMetadataUpdaterStub) UpdateTelegramBotMetadata(
	ctx context.Context,
	id string,
	telegramBotID int64,
	telegramUsername string,
	telegramBotName string,
) error {
	s.called = true
	s.id = id
	s.botID = telegramBotID
	s.username = telegramUsername
	s.name = telegramBotName
	return s.err
}

func TestRunnerApplyBotMetadata(t *testing.T) {
	t.Parallel()

	updater := &botMetadataUpdaterStub{}
	runner := NewRunnerWithDeps(Config{}, nil, nil, nil, updater, nil)

	err := runner.applyBotMetadata(
		context.Background(),
		manager.BotSpec{ID: "shop-main"},
		BotMetadata{
			ID:       777000123,
			Username: "shop_main_bot",
			Name:     "Shop Main Bot",
		},
	)
	require.NoError(t, err)

	require.True(t, updater.called)
	require.Equal(t, "shop-main", updater.id)
	require.Equal(t, int64(777000123), updater.botID)
	require.Equal(t, "shop_main_bot", updater.username)
	require.Equal(t, "Shop Main Bot", updater.name)
}
