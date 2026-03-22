package bootstrap

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"botmanager/internal/botconfig"
	"botmanager/internal/manager"
)

// stubRuntimeSource implements RuntimeSource for bootstrap tests.
type stubRuntimeSource struct {
	bots []botconfig.RuntimeBot
	err  error
}

// ListEnabledRuntimeBots returns configured runtime bots or preset error.
func (s *stubRuntimeSource) ListEnabledRuntimeBots(ctx context.Context) ([]botconfig.RuntimeBot, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.bots, nil
}

// stubRuntimeManager implements RuntimeManager for bootstrap tests.
type stubRuntimeManager struct {
	registerErrs map[string]error
	startErrs    map[string]error
	renameErrs   map[string]error

	registeredSpecs []manager.BotSpec
	startedIDs      []string
	renamed         map[string]string
}

// Register records bot spec and returns configured error for bot ID if any.
func (m *stubRuntimeManager) Register(spec manager.BotSpec) error {
	m.registeredSpecs = append(m.registeredSpecs, spec)

	if err, ok := m.registerErrs[spec.ID]; ok {
		return err
	}
	return nil
}

// Start records bot start attempt and returns configured error for bot ID if any.
func (m *stubRuntimeManager) Start(ctx context.Context, id string) error {
	m.startedIDs = append(m.startedIDs, id)

	if err, ok := m.startErrs[id]; ok {
		return err
	}
	return nil
}

// Rename records runtime rename attempt and returns configured error for bot ID if any.
func (m *stubRuntimeManager) Rename(id, name string) error {
	if m.renamed == nil {
		m.renamed = make(map[string]string)
	}
	m.renamed[id] = name

	if err, ok := m.renameErrs[id]; ok {
		return err
	}
	return nil
}

// TestNewStarter_PanicsOnNilDependencies verifies constructor guards.
func TestNewStarter_PanicsOnNilDependencies(t *testing.T) {
	source := &stubRuntimeSource{}
	mgr := &stubRuntimeManager{}

	require.Panics(t, func() {
		NewStarter(nil, mgr)
	})

	require.Panics(t, func() {
		NewStarter(source, nil)
	})
}

// TestStarter_StartEnabled verifies that starter loads enabled bots,
// registers them in manager, and starts their runtimes.
func TestStarter_StartEnabled(t *testing.T) {
	source := &stubRuntimeSource{
		bots: []botconfig.RuntimeBot{
			{
				ID:         "shop-main",
				Name:       "Shop Main",
				Token:      "token-main",
				DatabaseID: "main-db",
				IsEnabled:  true,
			},
			{
				ID:         "slow-bot",
				Name:       "Slow Bot",
				Token:      "token-slow",
				DatabaseID: "analytics-db",
				IsEnabled:  true,
			},
		},
	}
	mgr := &stubRuntimeManager{}
	starter := NewStarter(source, mgr)

	results, err := starter.StartEnabled(context.Background())
	require.NoError(t, err)

	require.Len(t, results, 2)

	require.Equal(t, "shop-main", results[0].ID)
	require.True(t, results[0].Registered)
	require.True(t, results[0].Started)
	require.NoError(t, results[0].Err)

	require.Equal(t, "slow-bot", results[1].ID)
	require.True(t, results[1].Registered)
	require.True(t, results[1].Started)
	require.NoError(t, results[1].Err)

	require.Len(t, mgr.registeredSpecs, 2)
	require.Equal(t, "shop-main", mgr.registeredSpecs[0].ID)
	require.Equal(t, "Shop Main", mgr.registeredSpecs[0].Name)
	require.Equal(t, "token-main", mgr.registeredSpecs[0].Token)

	require.Equal(t, []string{"shop-main", "slow-bot"}, mgr.startedIDs)
}

// TestStarter_StartEnabled_PropagatesSourceError verifies that
// bootstrap fails immediately when runtime source cannot list bots.
func TestStarter_StartEnabled_PropagatesSourceError(t *testing.T) {
	source := &stubRuntimeSource{
		err: errors.New("list bots failed"),
	}
	mgr := &stubRuntimeManager{}
	starter := NewStarter(source, mgr)

	results, err := starter.StartEnabled(context.Background())
	require.Error(t, err)
	require.ErrorContains(t, err, "list bots failed")
	require.Nil(t, results)
}

// TestStarter_StartEnabled_DuplicateRegistrationRenamesAndStarts verifies that
// duplicate registration is treated as non-fatal and runtime name is synchronized.
func TestStarter_StartEnabled_DuplicateRegistrationRenamesAndStarts(t *testing.T) {
	source := &stubRuntimeSource{
		bots: []botconfig.RuntimeBot{
			{
				ID:         "shop-main",
				Name:       "Shop Main Renamed",
				Token:      "token-main",
				DatabaseID: "main-db",
				IsEnabled:  true,
			},
		},
	}
	mgr := &stubRuntimeManager{
		registerErrs: map[string]error{
			"shop-main": manager.ErrDuplicateBotID,
		},
	}
	starter := NewStarter(source, mgr)

	results, err := starter.StartEnabled(context.Background())
	require.NoError(t, err)
	require.Len(t, results, 1)

	require.Equal(t, "shop-main", results[0].ID)
	require.False(t, results[0].Registered)
	require.True(t, results[0].Started)
	require.NoError(t, results[0].Err)

	require.Equal(t, "Shop Main Renamed", mgr.renamed["shop-main"])
	require.Equal(t, []string{"shop-main"}, mgr.startedIDs)
}

// TestStarter_StartEnabled_AlreadyRunningIsNonFatal verifies that
// already running runtime is treated as successful bootstrap outcome.
func TestStarter_StartEnabled_AlreadyRunningIsNonFatal(t *testing.T) {
	source := &stubRuntimeSource{
		bots: []botconfig.RuntimeBot{
			{
				ID:         "shop-main",
				Name:       "Shop Main",
				Token:      "token-main",
				DatabaseID: "main-db",
				IsEnabled:  true,
			},
		},
	}
	mgr := &stubRuntimeManager{
		startErrs: map[string]error{
			"shop-main": manager.ErrBotAlreadyRunning,
		},
	}
	starter := NewStarter(source, mgr)

	results, err := starter.StartEnabled(context.Background())
	require.NoError(t, err)
	require.Len(t, results, 1)

	require.Equal(t, "shop-main", results[0].ID)
	require.True(t, results[0].Registered)
	require.True(t, results[0].Started)
	require.NoError(t, results[0].Err)
}

// TestStarter_StartEnabled_CollectsPerBotErrors verifies that bootstrap
// continues after individual bot failures and returns aggregated error.
func TestStarter_StartEnabled_CollectsPerBotErrors(t *testing.T) {
	source := &stubRuntimeSource{
		bots: []botconfig.RuntimeBot{
			{
				ID:         "broken-register",
				Name:       "Broken Register",
				Token:      "token-1",
				DatabaseID: "db-1",
				IsEnabled:  true,
			},
			{
				ID:         "broken-start",
				Name:       "Broken Start",
				Token:      "token-2",
				DatabaseID: "db-2",
				IsEnabled:  true,
			},
			{
				ID:         "healthy-bot",
				Name:       "Healthy Bot",
				Token:      "token-3",
				DatabaseID: "db-3",
				IsEnabled:  true,
			},
		},
	}
	mgr := &stubRuntimeManager{
		registerErrs: map[string]error{
			"broken-register": errors.New("register failed"),
		},
		startErrs: map[string]error{
			"broken-start": errors.New("start failed"),
		},
	}
	starter := NewStarter(source, mgr)

	results, err := starter.StartEnabled(context.Background())
	require.Error(t, err)
	require.Len(t, results, 3)

	require.Equal(t, "broken-register", results[0].ID)
	require.False(t, results[0].Registered)
	require.False(t, results[0].Started)
	require.ErrorContains(t, results[0].Err, "register failed")

	require.Equal(t, "broken-start", results[1].ID)
	require.True(t, results[1].Registered)
	require.False(t, results[1].Started)
	require.ErrorContains(t, results[1].Err, "start failed")

	require.Equal(t, "healthy-bot", results[2].ID)
	require.True(t, results[2].Registered)
	require.True(t, results[2].Started)
	require.NoError(t, results[2].Err)

	require.ErrorContains(t, err, "register failed")
	require.ErrorContains(t, err, "start failed")
}

// TestStarter_StartEnabled_DuplicateRenameFailureStopsThatBot verifies that
// rename failure after duplicate registration is treated as fatal for that bot.
func TestStarter_StartEnabled_DuplicateRenameFailureStopsThatBot(t *testing.T) {
	source := &stubRuntimeSource{
		bots: []botconfig.RuntimeBot{
			{
				ID:         "shop-main",
				Name:       "Shop Main",
				Token:      "token-main",
				DatabaseID: "main-db",
				IsEnabled:  true,
			},
		},
	}
	mgr := &stubRuntimeManager{
		registerErrs: map[string]error{
			"shop-main": manager.ErrDuplicateBotID,
		},
		renameErrs: map[string]error{
			"shop-main": errors.New("rename failed"),
		},
	}
	starter := NewStarter(source, mgr)

	results, err := starter.StartEnabled(context.Background())
	require.Error(t, err)
	require.Len(t, results, 1)

	require.Equal(t, "shop-main", results[0].ID)
	require.False(t, results[0].Registered)
	require.False(t, results[0].Started)
	require.ErrorContains(t, results[0].Err, "rename failed")

	require.Empty(t, mgr.startedIDs)
}
