package testkit

import (
	"context"
	"sync"

	"botmanager/internal/domain"
)

// ---- Tx ----

// TxMock is a tiny TxManager implementation for tests.
// It records whether transaction wrapper was used.
type TxMock struct {
	Called bool
	Err    error
}

func (t *TxMock) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	t.Called = true
	if t.Err != nil {
		return t.Err
	}

	return fn(ctx)
}

// ---- EventBus ----

// EventBusSpy records publish events.
type EventBusSpy struct {
	mu sync.Mutex

	Called    bool
	Published []domain.Event
	Err       error
}

func (b *EventBusSpy) Publish(ctx context.Context, events ...domain.Event) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Called = true
	if b.Err != nil {
		return b.Err
	}

	b.Published = append(b.Published, events...)
	return nil
}

// Subscribe If yor EventBus interface also has Subscribe, keep it as noop in tests.
func (b *EventBusSpy) Subscribe(
	enventName string,
	handler func(context.Context, domain.Event) error,
) {
	// noop
}

func (b *EventBusSpy) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Called = false
	b.Published = nil
	b.Err = nil
}

// ---- UserRepository ----

// UserRepoMock is a function-field mock for UserRepository.
// Set ByIDFn/SaveFn per-test for strict behavior.
type UserRepoMock struct {
	SaveFn func(ctx context.Context, u *domain.User) error
	ByIDFn func(ctx context.Context, id int) (*domain.User, error)

	SaveCalls int
	ByIDCalls int
}

func (m *UserRepoMock) Save(ctx context.Context, u *domain.User) error {
	m.SaveCalls++
	if m.SaveFn == nil {
		panic("UserRepoMock.SaveFn is nil (unexpected call)")
	}
	return m.SaveFn(ctx, u)
}

func (m *UserRepoMock) ByID(ctx context.Context, id int) (*domain.User, error) {
	m.ByIDCalls++
	if m.ByIDFn == nil {
		panic("UserRepoMock.ByIDFn is nil (unexpected call)")
	}
	return m.ByIDFn(ctx, id)
}

// ---- OrderRepository ----

type OrderRepoMock struct {
	SaveFn func(ctx context.Context, o *domain.Order) error
	ByIDFn func(ctx context.Context, id int) (*domain.Order, error)

	SaveCalls int
	ByIDCalls int
}

func (m *OrderRepoMock) Save(ctx context.Context, o *domain.Order) error {
	m.SaveCalls++
	if m.SaveFn == nil {
		panic("OrderRepoMock.SaveFn is nil (unexpected call)")
	}
	return m.SaveFn(ctx, o)
}

func (m *OrderRepoMock) ByID(ctx context.Context, id int) (*domain.Order, error) {
	m.ByIDCalls++
	if m.ByIDFn == nil {
		panic("OrderRepoMock.ByIDFn is nil (unexpected call)")
	}
	return m.ByIDFn(ctx, id)
}

// ---- ProductRepository (read side) ----
// Adjust signature to match your actual service interfaces.

type ProductRepoMock struct {
	ByIDFn func(ctx context.Context, id int) (*domain.Product, error)

	ByIDCalls int
}

func (m *ProductRepoMock) ByID(ctx context.Context, id int) (*domain.Product, error) {
	m.ByIDCalls++
	if m.ByIDFn == nil {
		panic("ProductRepoMock.ByID is nil (unexpected call)")
	}
	return m.ByIDFn(ctx, id)
}
