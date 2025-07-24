package auctionusecase

import (
	"context"
	"time"

	auctionentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/auction_entity"
	bidentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/bid_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
	bidusecase "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/usecase/bid_usecase"
)

type AuctionInputDTO struct {
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
}

type AuctionOutputDTO struct {
	Id          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	Timestamp   time.Time        `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO         `json:"auction"`
	Bid     *bidusecase.BidOutputDTO `json:"bid,omitempty"`
}

type ProductCondition int64

type AuctionStatus int64

type AuctionUseCase struct {
	AuctionRepository auctionentity.AuctionRepositoryInterface
	BidRepository     bidentity.BidRepositoryInterface
}

func (au *AuctionUseCase) CreateAuction(ctx context.Context, auctionEntity AuctionInputDTO) *internalerror.InternalError {
	auction, err := auctionentity.CreateAuction(
		auctionEntity.ProductName,
		auctionEntity.Category,
		auctionEntity.Description,
		auctionentity.ProductCondition(auctionEntity.Condition),
	)

	if err != nil {
		return err
	}

	return au.AuctionRepository.CreateAuction(ctx, auction)
}
