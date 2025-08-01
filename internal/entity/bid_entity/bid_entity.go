package bidentity

import (
	"context"
	"time"

	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
	"github.com/google/uuid"
)

type Bid struct {
	Id        string
	UserId    string
	AuctionId string
	Amount    float64
	Timestamp time.Time
}

func CreateBid(userId, auctionId string, amount float64) (*Bid, *internalerror.InternalError) {
	bid := &Bid{
		Id:        uuid.New().String(),
		UserId:    userId,
		AuctionId: auctionId,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := bid.Validate(); err != nil {
		return nil, err
	}

	return bid, nil
}

func (b *Bid) Validate() *internalerror.InternalError {
	if err := uuid.Validate(b.UserId); err != nil {
		return internalerror.NewBadRequestError("UserId is not a valid uuid")
	}

	if err := uuid.Validate(b.AuctionId); err != nil {
		return internalerror.NewBadRequestError("AuctionId is not a valid uuid")
	}

	if b.Amount <= 0 {
		return internalerror.NewBadRequestError("Amount is not a valid value")
	}

	return nil
}

type BidRepositoryInterface interface {
	CreateBid(ctx context.Context, bidEntities []Bid) *internalerror.InternalError
	FindBidByAuctionId(ctx context.Context, auctionId string) ([]Bid, *internalerror.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*Bid, *internalerror.InternalError)
}
