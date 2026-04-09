package service

import "errors"

var (
	ErrCategoryCodeEmpty = errors.New("catalog category code is empty")
	ErrCategoryNameEmpty = errors.New("catalog category name is empty")

	ErrCityCodeEmpty = errors.New("catalog city code is empty")
	ErrCityNameEmpty = errors.New("catalog city name is empty")

	ErrDistrictCityIDInvalid = errors.New("catalog district city id is invalid")
	ErrDistrictCodeEmpty     = errors.New("catalog district code is empty")
	ErrDistrictNameEmpty     = errors.New("catalog district name is empty")

	ErrProductCategoryIDInvalid = errors.New("catalog product category id is invalid")
	ErrProductCodeEmpty         = errors.New("catalog product code is empty")
	ErrProductNameEmpty         = errors.New("catalog product name is empty")
)
