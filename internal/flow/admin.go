package flow

import "context"

// CategoryCreator defines the admin write use case required by flow.
//
// Flow depends on this narrow port to create catalog categories without
// depending on storage or transport-specific details.
type CategoryCreator interface {
	CreateCategory(ctx context.Context, params CreateCategoryParams) error
}

// CreateCategoryParams contains data required by flow admin category creation.
//
// This is a flow-local model passed through CategoryCreator.
type CreateCategoryParams struct {
	Code string
	Name string
}
