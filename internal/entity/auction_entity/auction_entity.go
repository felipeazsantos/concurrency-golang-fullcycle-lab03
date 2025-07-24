package auctionentity

import (
	"context"
	"time"

	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
	"github.com/google/uuid"
)

type Auction struct {
	Id          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	Timestamp   time.Time
}

type ProductCondition int
type AuctionStatus int

const (
	Active AuctionStatus = iota
	Completed
)

const (
	New ProductCondition = iota
	Used
	Refurbished
)

type AuctionRepositoryInterface interface {
	CreateAuction(ctx context.Context, auctionEntity *Auction) *internalerror.InternalError
	FindAuctionById(ctx context.Context, id string) (*Auction, *internalerror.InternalError)
	FindAuctions(ctx context.Context, status AuctionStatus, category string, productName string) ([]Auction, *internalerror.InternalError)
}

func CreateAuction(productName, category, description string, condition ProductCondition) (*Auction, *internalerror.InternalError) {
	auction := &Auction{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}

	if err := auction.Validate(); err != nil {
		return nil, err
	}

	return auction, nil
}

func (a *Auction) Validate() *internalerror.InternalError {
	if len(a.ProductName) <= 1 ||
		len(a.Category) <= 2 ||
		len(a.Description) <= 10 &&
			(a.Condition != New && a.Condition != Refurbished && a.Condition != Used) {
		return internalerror.NewBadRequestError("invalid auction object")
	}

	return nil
}
