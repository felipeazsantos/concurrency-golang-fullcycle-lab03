package bidusecase

import (
	"context"

	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
)

func (bu *BidUseCase) FindBidByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, *internalerror.InternalError) {
	bidList, err := bu.BidRepository.FindBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	var bidOutputList []BidOutputDTO
	for _, bid := range bidList {
		bidOutputList = append(bidOutputList, BidOutputDTO{
			Id:        bid.Id,
			UserId:    bid.UserId,
			AuctionId: bid.AuctionId,
			Amount:    bid.Amount,
			Timestamp: bid.Timestamp,
		})
	}

	return bidOutputList, nil
}

func (bu *BidUseCase) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*BidOutputDTO, *internalerror.InternalError) {
	bid, err := bu.BidRepository.FindWinningBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	return &BidOutputDTO{
		Id:        bid.Id,
		UserId:    bid.UserId,
		AuctionId: bid.AuctionId,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp,
	}, nil
}
