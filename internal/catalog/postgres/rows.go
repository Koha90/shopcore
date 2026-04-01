package postgres

type cityRow struct {
	ID        int
	Code      string
	Name      string
	NameLatin string
	SortOrder int
}

type categoryRow struct {
	ID          int
	Code        string
	Name        string
	NameLatin   string
	Description string
	SortOrder   int
}

type cityCategoryRow struct {
	CityID     int
	CategoryID int
	SortOrder  int
}

type districtRow struct {
	ID        int
	CityID    int
	Code      string
	Name      string
	NameLatin string
	SortOrder int
}

type productRow struct {
	ID          int
	CategoryID  int
	DistrictID  int
	Code        string
	Name        string
	NameLatin   string
	Description string
	SortOrder   int
}

type variantRow struct {
	ID           int
	ProductID    int
	Code         string
	Name         string
	NameLatin    string
	Description  string
	PriceMinor   int64
	CurrencyCode string
	SortOrder    int
}
