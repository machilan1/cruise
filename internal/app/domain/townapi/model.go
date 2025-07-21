package townapi

import (
	"strconv"

	"github.com/machilan1/cruise/internal/business/domain/town"
)

// AppTown represents an individual town.
type AppTown struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	City     AppCity `json:"city"`
	PostCode string  `json:"postCode"`
}

func ToAppTown(twn town.Town) AppTown {
	return AppTown{
		ID:       twn.ID,
		Name:     twn.Name,
		City:     toAppCity(twn.City),
		PostCode: strconv.Itoa(twn.PostCode),
	}
}

func ToAppTowns(twns []town.Town) []AppTown {
	items := make([]AppTown, len(twns))
	for i, twn := range twns {
		items[i] = ToAppTown(twn)
	}

	return items
}

type AppCity struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func toAppCity(city town.City) AppCity {
	return AppCity{
		ID:   city.ID,
		Name: city.Name,
	}
}

func toAppCities(cities []town.City) []AppCity {
	items := make([]AppCity, len(cities))
	for i, ct := range cities {
		items[i] = toAppCity(ct)
	}
	return items
}
