package flow

import "context"

// CategoryCreator creates catatolog categories for admin flow.
type CategoryCreator interface {
	CreateCategory(ctx context.Context, params CreateCategoryParams) error
}

// CreateCategoryParams contains data required by flow admin category creation.
type CreateCategoryParams struct {
	Code string
	Name string
}
