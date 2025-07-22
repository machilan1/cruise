package auctionhouse

import "time"

type AuctionHouse struct {
	ID        int
	Name      string
	Location  AuctionHouseLocation
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type AuctionHouseLocation struct {
	AddressDetail string
	CityID        int
	CityName      string
	TownID        int
	TownName      string
}

type NewAuctionHouse struct {
	Name          string
	AddressDetail string
	TownID        int
}

type UpdateAuctionHouse struct {
	Name          *string
	AddressDetail *string
	TownID        *int
}
