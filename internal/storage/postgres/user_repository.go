package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"botmanager/internal/domain"
	"botmanager/internal/service"
)

var _ service.UserRepository = (*UserRepository)(nil)

// UserRepository stores users in PostgreSQL.
type UserRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewUserRepository creates a new PostgreSQL user repository.
//
// logger may be nil. In that case slog.Default() is used.
func NewUserRepository(db *sql.DB, logger *slog.Logger) *UserRepository {
	if db == nil {
		panic("postgres: db is nil")
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// Save persists user state.
//
// New users are inserted, existing users are updated.
func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	if user.ID() == 0 {
		return r.insert(ctx, user)
	}
	return r.update(ctx, user)
}

// ByID returns user by its identifier.
func (r *UserRepository) ByID(ctx context.Context, id int) (*domain.User, error) {
	const q = `
		SELECT id, tg_id, tg_name, email, password_hash, role, balance,
					 is_enabled, admin_access_expires_at, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var (
		userID       int
		tgID         sql.NullInt64
		tgName       sql.NullString
		email        sql.NullString
		passwordHash sql.NullString
		role         string
		balance      int64
		isEnabled    bool
		adminExpires sql.NullTime
		creaedAt     time.Time
		updatedAt    time.Time
	)

	row := r.queryRow(ctx, q, id)
	if err := row.Scan(
		&userID,
		&tgID,
		&tgName,
		&email,
		&passwordHash,
		&role,
		&balance,
		&isEnabled,
		&adminExpires,
		&creaedAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	var tgIDPtr *int64
	if tgID.Valid {
		v := tgID.Int64
		tgIDPtr = &v
	}

	params := domain.NewUserParams{
		TgID:         tgIDPtr,
		TgName:       tgName.String,
		Email:        email.String,
		PasswordHash: passwordHash.String,
		Role:         domain.Role(role),
	}

	user, err := domain.NewUser(params)
	if err != nil {
		return nil, err
	}

	user.SetID(userID)

	if !isEnabled {
		user.Disable()
	}
	if balance > 0 {
		if err := user.AddBalance(balance); err != nil {
			return nil, err
		}
	}
	if adminExpires.Valid {
		user.GrantAdminAccess(adminExpires.Time)
	}
	_ = creaedAt
	_ = updatedAt

	return user, nil
}

func (r *UserRepository) insert(ctx context.Context, user *domain.User) error {
	const q = `
		INSERT INTO users (
			tg_id, tg_name, email, password_hash, role, balance,
			is_enabled, admin_access_expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id int
	if err := r.queryRow(
		ctx,
		q,
		nullableTelegramID(user),
		nullableTelegramName(user),
		user.Email(),
		user.PasswordHash(),
		user.Role(),
		user.Balance(),
		user.IsEnabled(),
		user.AdminAccessExpiresAt(),
	).Scan(&id); err != nil {
		return err
	}

	user.SetID(id)
	return nil
}

func (r *UserRepository) update(ctx context.Context, user *domain.User) error {
	const q = `
		UPDATE users
		SET tg_id = $1,
		    tg_name = $2,
		    email = $3,
		    password_hash = $4,
		    role = $5,
		    balance = $6,
		    is_enabled = $7,
		    admin_access_expires_at = $8,
		    updated_at = NOW()
		WHERE id = $9
	`
	_, err := r.exec(
		ctx,
		q,
		nullableTelegramID(user),
		nullableTelegramName(user),
		user.Email(),
		user.PasswordHash(),
		user.Role(),
		user.Balance(),
		user.IsEnabled(),
		user.AdminAccessExpiresAt(),
		user.ID(),
	)
	return err
}

func (r *UserRepository) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.ExecContext(ctx, query, args...)
	}
	return r.db.ExecContext(ctx, query, args...)
}

func (r *UserRepository) query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.QueryContext(ctx, query, args...)
	}
	return r.db.QueryContext(ctx, query, args...)
}

func (r *UserRepository) queryRow(ctx context.Context, query string, args ...any) *sql.Row {
	if tx, ok := txFromContext(ctx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return r.db.QueryRowContext(ctx, query, args...)
}

func nullableTelegramID(user *domain.User) any {
	if id, ok := user.TelegramID(); ok {
		return id
	}
	return nil
}

func nullableTelegramName(user *domain.User) any {
	if name, ok := user.TelegramName(); ok {
		return name
	}
	return nil
}
