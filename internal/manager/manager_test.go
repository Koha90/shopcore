package manager

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type stubRunner struct {
	runErr     error
	calls      int
	block      bool
	done       chan struct{}
	readyDelay time.Duration
	skipReady  bool
}

func (r *stubRunner) Run(ctx context.Context, spec BotSpec, ready func()) error {
	r.calls++

	if !r.skipReady {
		if r.readyDelay > 0 {
			select {
			case <-time.After(r.readyDelay):
			case <-ctx.Done():
				if r.done != nil {
					close(r.done)
				}
				return r.runErr
			}
		}
		ready()
	}

	if r.block {
		<-ctx.Done()
		if r.done != nil {
			close(r.done)
		}
		return r.runErr
	}

	return r.runErr
}

func TestNew_PanicOnNilRunner(t *testing.T) {
	require.Panics(t, func() { New(nil) })
}

func TestRegister(t *testing.T) {
	m := New(&stubRunner{})

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	status, err := m.Status("bot-1")
	require.NoError(t, err)
	require.Equal(t, StatusStopped, status)
}

func TestRegister_DuplicateID(t *testing.T) {
	m := New(&stubRunner{})

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	err = m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "duplicate",
		Token: "another-token",
	})
	require.ErrorIs(t, err, ErrDuplicateBotID)
}

func TestStart(t *testing.T) {
	runner := &stubRunner{}
	m := New(runner)

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	err = m.Start(context.Background(), "bot-1")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		status, err := m.Status("bot-1")
		return err == nil && status == StatusStopped
	}, time.Second, 10*time.Millisecond)

	require.Equal(t, 1, runner.calls)
}

func TestStart_UnknownBot(t *testing.T) {
	m := New(&stubRunner{})

	err := m.Start(context.Background(), "missing")
	require.ErrorIs(t, err, ErrBotNotFound)
}

func TestStart_AlreadyRunning(t *testing.T) {
	runner := &stubRunner{
		block: true,
	}
	m := New(runner)

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	err = m.Start(context.Background(), "bot-1")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		status, err := m.Status("bot-1")
		return err == nil && (status == StatusStarting || status == StatusRunning)
	}, time.Second, 10*time.Millisecond)

	err = m.Start(context.Background(), "bot-1")
	require.ErrorIs(t, err, ErrBotAlreadyRunning)

	err = m.Stop("bot-1")
	require.NoError(t, err)
}

func TestStart_RunErrorSetsFailedStatus(t *testing.T) {
	runner := &stubRunner{
		runErr:    errors.New("boom"),
		skipReady: true,
	}
	m := New(runner)

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	err = m.Start(context.Background(), "bot-1")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		info, err := m.Info("bot-1")
		return err == nil && info.Status == StatusFailed && info.LastError == "boom"
	}, time.Second, 10*time.Millisecond)
}

func TestStop(t *testing.T) {
	runner := &stubRunner{
		block: true,
	}
	m := New(runner)

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	err = m.Start(context.Background(), "bot-1")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		status, err := m.Status("bot-1")
		return err == nil && (status == StatusStarting || status == StatusRunning)
	}, time.Second, 10*time.Millisecond)

	err = m.Stop("bot-1")
	require.NoError(t, err)

	status, err := m.Status("bot-1")
	require.NoError(t, err)
	require.Equal(t, StatusStopped, status)
}

func TestStop_UnknownBot(t *testing.T) {
	m := New(&stubRunner{})

	err := m.Stop("missing")
	require.ErrorIs(t, err, ErrBotNotFound)
}

func TestStop_NotRunning(t *testing.T) {
	m := New(&stubRunner{})

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	err = m.Stop("bot-1")
	require.ErrorIs(t, err, ErrBotNotRunning)
}

func TestRestart(t *testing.T) {
	runner := &stubRunner{
		block: true,
	}
	m := New(runner)

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	err = m.Start(context.Background(), "bot-1")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		status, err := m.Status("bot-1")
		return err == nil && (status == StatusStarting || status == StatusRunning)
	}, time.Second, 10*time.Millisecond)

	err = m.Restart(context.Background(), "bot-1")
	require.NoError(t, err)

	err = m.Stop("bot-1")
	require.NoError(t, err)

	require.GreaterOrEqual(t, runner.calls, 2)
}

func TestStatus_UnknownBot(t *testing.T) {
	m := New(&stubRunner{})

	_, err := m.Status("missing")
	require.ErrorIs(t, err, ErrBotNotFound)
}

func TestInfo(t *testing.T) {
	m := New(&stubRunner{})

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	info, err := m.Info("bot-1")
	require.NoError(t, err)
	require.Equal(t, "bot-1", info.ID)
	require.Equal(t, "main", info.Name)
	require.Equal(t, StatusStopped, info.Status)
	require.Empty(t, info.LastError)
}

func TestList(t *testing.T) {
	m := New(&stubRunner{})

	err := m.Register(BotSpec{
		ID:    "bot-2",
		Name:  "second",
		Token: "token-2",
	})
	require.NoError(t, err)

	err = m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "first",
		Token: "token-1",
	})
	require.NoError(t, err)

	list := m.List()
	require.Len(t, list, 2)
	require.Equal(t, "bot-1", list[0].ID)
	require.Equal(t, "bot-2", list[1].ID)
}

func TestRename(t *testing.T) {
	m := New(&stubRunner{})

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	err = m.Rename("bot-1", "renamed")
	require.NoError(t, err)

	info, err := m.Info("bot-1")
	require.NoError(t, err)
	require.Equal(t, "renamed", info.Name)
	require.Equal(t, StatusStopped, info.Status)
	require.Empty(t, info.LastError)
}

func TestRename_UnknownBot(t *testing.T) {
	m := New(&stubRunner{})

	err := m.Rename("missing", "renamed")
	require.ErrorIs(t, err, ErrBotNotFound)
}

func TestRename_ListReflectsUpdatedName(t *testing.T) {
	m := New(&stubRunner{})

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	err = m.Rename("bot-1", "renamed")
	require.NoError(t, err)

	list := m.List()
	require.Len(t, list, 1)
	require.Equal(t, "renamed", list[0].Name)
}

func TestStart_RemainsStartingUntilReady(t *testing.T) {
	runner := &stubRunner{
		block:      true,
		readyDelay: 200 * time.Millisecond,
	}
	m := New(runner)

	err := m.Register(BotSpec{
		ID:    "bot-1",
		Name:  "main",
		Token: "token",
	})
	require.NoError(t, err)

	err = m.Start(context.Background(), "bot-1")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		status, err := m.Status("bot-1")
		return err == nil && status == StatusStarting
	}, 100*time.Millisecond, 10*time.Millisecond)

	require.Eventually(t, func() bool {
		status, err := m.Status("bot-1")
		return err == nil && status == StatusRunning
	}, time.Second, 10*time.Millisecond)

	err = m.Stop("bot-1")
	require.NoError(t, err)
}
