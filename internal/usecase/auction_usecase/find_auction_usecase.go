package auctionusecase

import (
	"context"

	auctionentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/auction_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
	bidusecase "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/usecase/bid_usecase"
)

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

func (au *AuctionUseCase) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*WinningInfoOutputDTO, *internalerror.InternalError) {
	auction, err := au.AuctionRepository.FindAuctionById(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	auctionOutput := AuctionOutputDTO{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   ProductCondition(auction.Condition),
		Status:      AuctionStatus(auction.Status),
		Timestamp:   auction.Timestamp,
	}

	bidWinning, err := au.BidRepository.FindWinningBidByAuctionId(ctx, auction.Id)
	if err != nil {
		return &WinningInfoOutputDTO{
			Auction: auctionOutput,
			Bid:     nil,
		}, err
	}

	bidOutput := &bidusecase.BidOutputDTO{
		Id:        bidWinning.Id,
		UserId:    bidWinning.UserId,
		AuctionId: bidWinning.AuctionId,
		Amount:    bidWinning.Amount,
		Timestamp: bidWinning.Timestamp,
	}

	return &WinningInfoOutputDTO{
		Auction: auctionOutput,
		Bid:     bidOutput,
	}, nil
}
