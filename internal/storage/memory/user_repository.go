package memory

import (
	"context"
	"sync"

	"github.com/koha90/shopcore/internal/domain"
)

// UserRepository stores users in process memory.
//
// It is intended for local development, tests and simple runtime scenarios.
// Reposutory assigns incremental IDs to new users on first save.
type UserRepository struct {
	mu     *sync.Mutex
	users  map[int]*domain.User
	nextID int
}

// NewUserRepository creates a new in-memory user repository.
//
// mu must point to shared storage mutes used by all repositories
// participating in the same logical transaction.
func NewUserRepository(mu *sync.Mutex) *UserRepository {
	return &UserRepository{
		mu:     mu,
		users:  make(map[int]*domain.User),
		nextID: 1,
	}
}

// Save stores user in memory.
//
// If user does not yet have an ID, repository assigns a new one.
func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.ID() == 0 {
		user.SetID(r.nextID)
		r.nextID++
	}

	r.users[user.ID()] = user
	return nil
}

// ByID returns user by its identifier.
func (r *UserRepository) ByID(ctx context.Context, id int) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}
