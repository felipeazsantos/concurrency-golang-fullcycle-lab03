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
	ProductName string           `json:"product_name" biding:"required,min=1"`
	Category    string           `json:"category" biding:"required,min=2"`
	Description string           `json:"description"  biding:"required,min=10"`
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

type AuctionUseCaseInterface interface {
	CreateAuction(ctx context.Context, auctionInputDTO AuctionInputDTO) *internalerror.InternalError
	FindAuctionById(ctx context.Context, id string) (*AuctionOutputDTO, *internalerror.InternalError)
	FindAuctions(ctx context.Context, status AuctionStatus, category string, productName string) ([]AuctionOutputDTO, *internalerror.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*WinningInfoOutputDTO, *internalerror.InternalError)
}

type AuctionUseCase struct {
	AuctionRepository auctionentity.AuctionRepositoryInterface
	BidRepository     bidentity.BidRepositoryInterface
}

func NewAuctionUseCase(
	auctionRepository auctionentity.AuctionRepositoryInterface,
	bidRepository bidentity.BidRepositoryInterface,
) AuctionUseCaseInterface {
	return &AuctionUseCase{
		AuctionRepository: auctionRepository,
		BidRepository:     bidRepository,
	}
}

func (au *AuctionUseCase) CreateAuction(ctx context.Context, auctionInputDTO AuctionInputDTO) *internalerror.InternalError {
	auction, err := auctionentity.CreateAuction(
		auctionInputDTO.ProductName,
		auctionInputDTO.Category,
		auctionInputDTO.Description,
		auctionentity.ProductCondition(auctionInputDTO.Condition),
	)

	if err != nil {
		return err
	}

	return au.AuctionRepository.CreateAuction(ctx, auction)
}
