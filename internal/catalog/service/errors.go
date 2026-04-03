package service

import "errors"

var (
	ErrCategoryCodeEmpty = errors.New("catalog category code is empty")
	ErrCategoryNameEmpty = errors.New("catalog category name is empty")
)
