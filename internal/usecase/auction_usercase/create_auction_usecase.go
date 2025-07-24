package auctionusercase

import (
	"context"
	"time"

	auctionentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/auction_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
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

type ProductCondition int64

type AuctionStatus int64

type AuctionUseCase struct {
	AuctionRepository auctionentity.AuctionRepositoryInterface
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

func (au *AuctionUseCase) FindAuctionById(ctx context.Context, id string) (*AuctionOutputDTO, *internalerror.InternalError) {
	auction, err := au.AuctionRepository.FindAuctionById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &AuctionOutputDTO{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   ProductCondition(auction.Condition),
		Status:      AuctionStatus(auction.Status),
		Timestamp:   auction.Timestamp,
	}, nil

}

func (au *AuctionUseCase) FindAuctions(ctx context.Context, status AuctionStatus, category string, productName string) ([]AuctionOutputDTO, *internalerror.InternalError) {
	auctions, err := au.AuctionRepository.FindAuctions(ctx, auctionentity.AuctionStatus(status), category, productName)
	if err != nil {
		return nil, err
	}

	var auctionsOutput []AuctionOutputDTO
	for _, auction := range auctions {
		auctionsOutput = append(auctionsOutput, AuctionOutputDTO{
			Id:          auction.Id,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Description: auction.Description,
			Condition:   ProductCondition(auction.Condition),
			Status:      AuctionStatus(auction.Status),
			Timestamp:   auction.Timestamp,
		})
	}

	return auctionsOutput, nil
}
