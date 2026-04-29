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

// CategoryListItem contains one category option for admin selection flows.
type CategoryListItem struct {
	ID    int
	Code  string
	Label string
}

// CategoryLister defines the admin read use case required by flow
// to select an existing category before nested catalog creation steps.
type CategoryLister interface {
	ListCategories(ctx context.Context) ([]CategoryListItem, error)
}

// CityCreator defines the admin write use case required by flow.
type CityCreator interface {
	CreateCity(ctx context.Context, params CreateCityParams) error
}

// CreateCityParams contains data required by flow admin city creation.
type CreateCityParams struct {
	Code string
	Name string
}

// CityListItem contains one city option for admin selection flows.
type CityListItem struct {
	ID    int
	Code  string
	Label string
}

// CityLister defines the admin read use case required by flow
// to select an existing city before nested catalog creation steps.
type CityLister interface {
	ListCities(ctx context.Context) ([]CityListItem, error)
}

// DistrictCreator defines the admin write use case required by flow.
type DistrictCreator interface {
	CreateDistrict(ctx context.Context, params CreateDistrictParams) error
}

// CreateDistrictParams contains data required by flow admin district creation.
type CreateDistrictParams struct {
	CityID int
	Code   string
	Name   string
}

// DistrictListItem contains one district option for admon selection flows.
type DistrictListItem struct {
	ID     int
	Code   string
	Label  string
	CityID int
}

// DistrictLister defines the admin read use case required by flow
// to select an existing district before nested catalog creation steps.
type DistrictLister interface {
	ListDistricts(ctx context.Context) ([]DistrictListItem, error)
	ListDistrictsByCity(ctx context.Context, cityID int) ([]DistrictListItem, error)
}

// ProductCreator defines the admin write use case required by flow.
type ProductCreator interface {
	CreateProduct(ctx context.Context, params CreateProductParams) error
}

// CreateProductParams contains data required by flow admin product creation.
type CreateProductParams struct {
	CategoryID int
	Code       string
	Name       string
}

// ProductListItem contains one product option for admin selection flows.
type ProductListItem struct {
	ID    int
	Code  string
	Label string
}

// ProductLister defines the admin read use case required by flow
// to select an existing product before nested catalog creation steps.
type ProductLister interface {
	ListProducts(ctx context.Context) ([]ProductListItem, error)
	ListProductsByCategory(ctx context.Context, categoryID int) ([]ProductListItem, error)
}

// VariantCreator defines the admin write use case required by flow.
type VariantCreator interface {
	CreateVariant(ctx context.Context, params CreateVariantParams) error
}

// CreateVariantParams contains data required by flow admin variant creation.
type CreateVariantParams struct {
	ProductID int
	Code      string
	Name      string
}

// VariantListItem contains one variant option for admin selection flows.
type VariantListItem struct {
	ID           int
	Code         string
	Label        string
	ProductLabel string
}

// DistrictPlacementVariantListItem contains one placed variant option
// for district placement edit flows.
type DistrictPlacementVariantListItem struct {
	ID        int
	Code      string
	Label     string
	Price     int
	PriceText string
}

// VariantLister defines the admin read use case required by flow
// to select an existing variant before nested catalog creation steps.
type VariantLister interface {
	ListVariants(ctx context.Context) ([]VariantListItem, error)
	ListVariantsByProduct(ctx context.Context, productID int) ([]VariantListItem, error)
}

// DistrictVariantCreator defines the admin write use case required by flow.
type DistrictVariantCreator interface {
	CreateDistrictVariant(ctx context.Context, params CreateDistrictVariantParams) error
}

// CreateDistrictVariantParams contains data required by flow admin placement creation.
type CreateDistrictVariantParams struct {
	DistrictID int
	VariantID  int
	Price      int
}

// DistrictVariantPriceUpdater updates price for an existing district variant placement
type DistrictVariantPriceUpdater interface {
	UpdateDistrictVariantPrice(ctx context.Context, params UpdateDistrictVariantPriceParams) error
}

// UpdateDistrictVariantPriceParams contains input for district variant price update.
type UpdateDistrictVariantPriceParams struct {
	DistrictID int
	VariantID  int
	Price      int
}

// DistrictPlacementReader defines filtered admin read use cases
// for district placement edit flows.
type DistrictPlacementReader interface {
	ListDistrictCategories(ctx context.Context, districtID int) ([]CategoryListItem, error)
	ListDistrictProducts(ctx context.Context, districtID, categoryID int) ([]ProductListItem, error)
	ListDistrictVariants(ctx context.Context, districtID, productID int) ([]DistrictPlacementVariantListItem, error)

	ListAvailableVariantsForDistrictProduct(ctx context.Context, districtID, productID int) ([]VariantListItem, error)
}

// ProductImageUpdater defines the admin write use case required by flow.
type ProductImageUpdater interface {
	UpdateProductImage(ctx context.Context, params UpdateProductImageParams) error
}

// UpdateProductImageParams contains input for product image update.
type UpdateProductImageParams struct {
	ProductID int
	ImageURL  string
}

// VariantImageUpdater defines the admin write use case required by flow.
type VariantImageUpdater interface {
	UpdateVariantImage(ctx context.Context, params UpdateVariantImageParams) error
}

// UpdateVariantImageParams contains input for variant image update.
type UpdateVariantImageParams struct {
	VariantID int
	ImageURL  string
}
