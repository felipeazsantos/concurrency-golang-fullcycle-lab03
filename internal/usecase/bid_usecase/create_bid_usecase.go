package bidusecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/logger"
	bidentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/bid_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
)

type BidInputDTO struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

type BidOutputDTO struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type BidUseCase struct {
	BidRepository       bidentity.BidRepositoryInterface
	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan bidentity.Bid
}

var bidBatch []bidentity.Bid

type BidUseCasesInterface interface {
	CreateBid(ctx context.Context, bidInput BidInputDTO) *internalerror.InternalError
	FindBidByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, *internalerror.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*BidOutputDTO, *internalerror.InternalError)
}

func NewBidUseCase(bidRepository bidentity.BidRepositoryInterface) BidUseCasesInterface {
	maxBatchSize := getMaxBatchSize()
	maxBatchSizeInterval := getMaxBatchSizeInterval()

	bidUseCase := &BidUseCase{
		BidRepository:       bidRepository,
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxBatchSizeInterval,
		timer:               time.NewTimer(maxBatchSizeInterval),
		bidChannel:          make(chan bidentity.Bid, maxBatchSize),
	}

	bidUseCase.triggerCreateBidRoutine(context.Background())

	return bidUseCase
}

func (bu *BidUseCase) CreateBid(ctx context.Context, bidInput BidInputDTO) *internalerror.InternalError {
	bidEntity, err := bidentity.CreateBid(bidInput.UserId, bidInput.AuctionId, bidInput.Amount)
	if err != nil {
		return err
	}

	bu.bidChannel <- *bidEntity

	return nil
}

func (bu *BidUseCase) triggerCreateBidRoutine(ctx context.Context) {
	go func() {
		defer close(bu.bidChannel)
		for {
			select {
			case bidEntity, ok := <-bu.bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
							logger.Error("error trying to process bid batch list", err)
						}
					}
					return
				}

				bidBatch = append(bidBatch, bidEntity)

				if len(bidBatch) > 0 && len(bidBatch) >= bu.maxBatchSize {
					if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}

					bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
				}
			case <-bu.timer.C:
				if len(bidBatch) > 0 {
					if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}

					bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
				}
			}
		}
	}()
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

func getMaxBatchSize() int {
	size, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}
	return size
}
