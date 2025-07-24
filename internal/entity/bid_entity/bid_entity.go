package bidentity

import (
	"context"
	"time"

	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
)

type Bid struct {
	Id        string
	UserId    string
	AuctionId string
	Amount    float64
	Timestamp time.Time
}

type BidRepositoryInterface interface {
	CreateBid(ctx context.Context, bidEntities []Bid) *internalerror.InternalError
	FindBidByAuctionId(ctx context.Context, auctionId string) ([]Bid, *internalerror.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*Bid, *internalerror.InternalError)
}
