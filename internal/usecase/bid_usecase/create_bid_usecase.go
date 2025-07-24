package bidusecase

import (
	"context"
	"time"

	bidentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/bid_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
)

type BidInputDTO struct {
}

type BidOutputDTO struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type BidUseCase struct {
	BidRepository bidentity.BidRepositoryInterface
}

func (bu *BidUseCase) CreateBid(ctx context.Context, bidEntities []BidInputDTO) *internalerror.InternalError {

}
