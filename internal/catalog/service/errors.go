package service

import "errors"

var (
	ErrCategoryCodeEmpty = errors.New("catalog category code is empty")
	ErrCategoryNameEmpty = errors.New("catalog category name is empty")
)

var (
	ErrCityCodeEmpty = errors.New("catalog city code is empty")
	ErrCityNameEmpty = errors.New("catalog city name is empty")
)
