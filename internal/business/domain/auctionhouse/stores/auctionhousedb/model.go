package auctionhousedb

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/auctionhouse"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb/dbjson"
)

type dbAuctionHouse struct {
	ID            int                                       `db:"auction_house_id"`
	Name          string                                    `db:"auction_house_name"`
	Location      dbjson.JSONColumn[dbAuctionHouseLocation] `db:"location"`
	TownID        int                                       `db:"town_id"`
	AddressDetail string                                    `db:"address_detail"`
	CreatedAt     time.Time                                 `db:"created_at"`
	UpdatedAt     time.Time                                 `db:"updated_at"`
	DeletedAt     *time.Time                                `db:"deleted_at"`
}

func toCoreAuctionHouse(dah dbAuctionHouse) (auctionhouse.AuctionHouse, error) {
	return auctionhouse.AuctionHouse{
		ID:   dah.ID,
		Name: dah.Name,
		Location: auctionhouse.AuctionHouseLocation{
			AddressDetail: dah.Location.Get().AddressDetail,
			CityID:        dah.Location.Get().CityID,
			CityName:      dah.Location.Get().CityName,
			TownID:        dah.Location.Get().TownID,
			TownName:      dah.Location.Get().TownName,
		},
		CreatedAt: dah.CreatedAt,
		UpdatedAt: dah.UpdatedAt,
		DeletedAt: dah.DeletedAt,
	}, nil
}

func toCoreAuctionHouses(dahs []dbAuctionHouse) ([]auctionhouse.AuctionHouse, error) {
	ahs := make([]auctionhouse.AuctionHouse, len(dahs))
	for i, v := range dahs {
		ah, err := toCoreAuctionHouse(v)
		if err != nil {
			return nil, err
		}
		ahs[i] = ah
	}
	return ahs, nil
}

func toDBAuctionHouse(ah auctionhouse.AuctionHouse) dbAuctionHouse {
	return dbAuctionHouse{
		ID:            ah.ID,
		Name:          ah.Name,
		TownID:        ah.Location.TownID,
		AddressDetail: ah.Location.AddressDetail,
		CreatedAt:     ah.CreatedAt,
		UpdatedAt:     ah.UpdatedAt,
		DeletedAt:     ah.DeletedAt,
	}
}

type dbAuctionHouseLocation struct {
	AddressDetail string `json:"address_detail"`
	CityID        int    `json:"city_id"`
	CityName      string `json:"city_name"`
	TownID        int    `json:"town_id"`
	TownName      string `json:"town_name"`
}
