package auctionapi

import "github.com/machilan1/cruise/internal/business/domain/auction"

var orderByFields = map[string]string{
	"createdAt": auction.OrderByCreatedAt,
	"updatedAt": auction.OrderByUpdatedAt,
}
