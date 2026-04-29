package postgres

import "errors"

var (
	ErrDistrictVariantNotFound = errors.New("catalog district variant not found")
	ErrProductNotFound         = errors.New("catalog product not found")
	ErrVariantNotFound         = errors.New("catalog variant not found")
)
