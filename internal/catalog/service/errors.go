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

	ErrProductIDInvalid         = errors.New("catalog product id is invalid")
	ErrProductCategoryIDInvalid = errors.New("catalog product category id is invalid")
	ErrProductCodeEmpty         = errors.New("catalog product code is empty")
	ErrProductNameEmpty         = errors.New("catalog product name is empty")

	ErrVariantIDInvalid        = errors.New("catalog variant id is invalid")
	ErrVariantProductIDInvalid = errors.New("catalog variant product id is invalid")
	ErrVariantCodeEmpty        = errors.New("catalog variant code is empty")
	ErrVariantNameEmpty        = errors.New("catalog variant name is empty")

	ErrDistrictVariantDistrictIDInvalid = errors.New("catalog district variant district id is invalid")
	ErrDistrictVariantVariantIDInvalid  = errors.New("catalog district variant variant id is invalid")
	ErrDistrictVariantPriceInvalid      = errors.New("catalog district variant price is invalid")
	ErrDistrictVariantAlreadyExists     = errors.New("catalog district variant already exists")

	ErrImageURLInvalid = errors.New("catalog image url is invalid")
)
