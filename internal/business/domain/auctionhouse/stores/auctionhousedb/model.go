package auctionhousedb

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/auctionhouse"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb/dbjson"
)

type dbAuctionHouse struct {
	ID        int                                       `db:"auction_house_id"`
	Name      string                                    `db:"auction_house_name"`
	Location  dbjson.JSONColumn[dbAuctionHouseLocation] `db:"location"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func toCoreAuctionHouse(dah dbAuctionHouse) (auctionhouse.AuctionHouse, error) {

	return auctionhouse.AuctionHouse{
		ID:   dah.ID,
		Name: dah.Name,
	}, nil
}

type dbAuctionHouseLocation struct {
	Address  string `json:"address"`
	CityID   int    `json:"city_id"`
	CityName string `json:"city_name"`
	TownID   int    `json:"town_id"`
	TownName string `json:"town_name"`
}
