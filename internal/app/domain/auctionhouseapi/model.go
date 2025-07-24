package auctionhouseapi

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/auctionhouse"
)

type AppAuctionHouse struct {
	ID        int                     `json:"id"`
	Name      string                  `json:"name"`
	Location  AppAuctionHouseLocation `json:"location"`
	CreatedAt time.Time               `json:"createdAt"`
	UpdatedAt time.Time               `json:"updatedAt"`
	DeletedAt *time.Time              `json:"deletedAt"`
}

func toAppAuctionHouse(ah auctionhouse.AuctionHouse) AppAuctionHouse {
	return AppAuctionHouse{
		ID:   ah.ID,
		Name: ah.Name,
		Location: AppAuctionHouseLocation{
			AddressDetail: ah.Location.AddressDetail,
			CityID:        ah.Location.CityID,
			CityName:      ah.Location.CityName,
			TownID:        ah.Location.TownID,
			TownName:      ah.Location.TownName,
		},
		CreatedAt: ah.CreatedAt,
		UpdatedAt: ah.CreatedAt,
		DeletedAt: ah.DeletedAt,
	}
}

func toAppAuctionHouses(ahs []auctionhouse.AuctionHouse) []AppAuctionHouse {
	aahs := make([]AppAuctionHouse, len(ahs))
	for i, v := range ahs {
		aahs[i] = toAppAuctionHouse(v)
	}
	return aahs
}

type AppAuctionHouseLocation struct {
	AddressDetail string `json:"addressDetail"`
	CityID        int    `json:"cityId"`
	CityName      string `json:"cityName"`
	TownID        int    `json:"townId"`
	TownName      string `json:"townName"`
}

type AppNewAuctionHouse struct {
	Name          string `json:"name" validate:"required"`
	AddressDetail string `json:"addressDetail" validate:"required"`
	TownID        int    `json:"townId" validate:"required"`
}

func toCoreNewAuctionHouse(anah AppNewAuctionHouse) (auctionhouse.NewAuctionHouse, error) {
	return auctionhouse.NewAuctionHouse{
		Name:          anah.Name,
		AddressDetail: anah.AddressDetail,
		TownID:        anah.TownID,
	}, nil
}

type AppUpdateAuctionHouse struct {
	Name          *string `json:"name"`
	AddressDetail *string `json:"addressDetail"`
	TownID        *int    `json:"townId"`
}

func toCoreUpdateAuctionHouse(auah AppUpdateAuctionHouse) (auctionhouse.UpdateAuctionHouse, error) {
	return auctionhouse.UpdateAuctionHouse{
		Name:          auah.Name,
		AddressDetail: auah.AddressDetail,
		TownID:        auah.TownID,
	}, nil
}
