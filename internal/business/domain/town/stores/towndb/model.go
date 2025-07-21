package towndb

import (
	"fmt"

	"github.com/machilan1/cruise/internal/business/domain/town"
)

// dbTown represents an individual town.
type dbTown struct {
	ID       int    `db:"town_id"`
	Name     string `db:"town_name"`
	CityID   int    `db:"city_id"`
	CityName string `db:"city_name"`
	PostCode int    `db:"post_code"`
}

func toCoreTown(dbTwn dbTown) (town.Town, error) {
	ct := town.City{
		ID:   dbTwn.CityID,
		Name: dbTwn.CityName,
	}
	twn := town.Town{
		ID:       dbTwn.ID,
		Name:     dbTwn.Name,
		City:     ct,
		PostCode: dbTwn.PostCode,
	}

	return twn, nil
}

func toCoreTowns(dbTowns []dbTown) ([]town.Town, error) {
	twns := make([]town.Town, len(dbTowns))
	for i, dbTwn := range dbTowns {
		var err error
		twns[i], err = toCoreTown(dbTwn)
		if err != nil {
			return nil, fmt.Errorf("parse type: %w", err)
		}
	}

	return twns, nil
}

type dbCity struct {
	ID   int    `db:"city_id"`
	Name string `db:"city_name"`
}

func toCoreCity(dbCity dbCity) town.City {
	return town.City{
		ID:   dbCity.ID,
		Name: dbCity.Name,
	}
}

func toCoreCities(dbCities []dbCity) []town.City {
	cities := make([]town.City, len(dbCities))
	for i, dbCt := range dbCities {
		cities[i] = toCoreCity(dbCt)
	}

	return cities
}
