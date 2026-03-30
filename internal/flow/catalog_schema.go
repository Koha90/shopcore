package flow

// CatalogLevel identifies one logical step inside catalog navigation.
//
// Levels describe business navigation order independently from transport.
// Example sequence:
//
//	city -> category -> district -> product -> variant
type CatalogLevel string

const (
	// LevelCity selects the city context for the catalog.
	LevelCity CatalogLevel = "city"

	// LevelCategory selects product category inside selected city.
	LevelCategory CatalogLevel = "category"

	// LevelDistrict selects district or area inside current context.
	LevelDistrict CatalogLevel = "district"

	// LevelDistrict selects a concrete product.
	LevelProduct CatalogLevel = "product"

	// LevelVariant selects a concrete product option, such as size or weight.
	LevelVariant CatalogLevel = "variant"
)

// CatalogSchema describes ordered catalog navigation levels.
//
// StartScenario defines how user enters catalog.
// CatalogSchema defines what user selects next inside catalog.
type CatalogSchema struct {
	Levels []CatalogLevel
}

// DemoCatalogSchema returns the default demo schema used by flow tests
// and in-memory catalog navigation.
func DemoCatalogSchema() CatalogSchema {
	return CatalogSchema{
		Levels: []CatalogLevel{
			LevelCity,
			LevelCategory,
			LevelDistrict,
			LevelProduct,
			LevelVariant,
		},
	}
}

// First returns the first level of the schema.
func (s CatalogSchema) First() (CatalogLevel, bool) {
	if len(s.Levels) == 0 {
		return "", false
	}
	return s.Levels[0], true
}

// Next returns the next level after the provided one.
func (s CatalogSchema) Next(level CatalogLevel) (CatalogLevel, bool) {
	for i, v := range s.Levels {
		if v != level {
			continue
		}
		if i+1 >= len(s.Levels) {
			return "", false
		}
		return s.Levels[i+1], true
	}
	return "", false
}
