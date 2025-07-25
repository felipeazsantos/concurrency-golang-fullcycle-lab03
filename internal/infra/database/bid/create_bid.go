package bid

import (
	"context"
	"sync"

	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/logger"
	auctionentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/auction_entity"
	bidentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/bid_entity"
	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/infra/database/auction"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidEntityMongo struct {
	Id        string  `bson:"_id"`
	UserId    string  `bson:"user_id"`
	AuctionId string  `bson:"auction_id"`
	Amount    float64 `bson:"amount"`
	Timestamp int64   `bson:"timestamp"`
}

type BidRepository struct {
	Collection        *mongo.Collection
	AuctionRepository *auction.AuctionRepository
}

func NewBidRepository(database *mongo.Database, auctionRepository *auction.AuctionRepository) *BidRepository {
	return &BidRepository{
		Collection:        database.Collection("bids"),
		AuctionRepository: auctionRepository,
	}
}

func (br *BidRepository) CreateBid(
	ctx context.Context,
	bidEntities []bidentity.Bid,
) *internalerror.InternalError {
	var wg sync.WaitGroup

	for _, bid := range bidEntities {
		wg.Add(1)

		go func(bidValue bidentity.Bid) {
			defer wg.Done()

			auctionEntity, err := br.AuctionRepository.FindAuctionById(ctx, bidValue.AuctionId)
			if err != nil {
				logger.Error("Error trying to find auction by id", err)
				return
			}

			if auctionEntity.Status != auctionentity.Active {
				return
			}

			bidEntityMongo := &BidEntityMongo{
				Id:        bidValue.Id,
				UserId:    bidValue.UserId,
				AuctionId: bidValue.AuctionId,
				Amount:    bidValue.Amount,
				Timestamp: bidValue.Timestamp.Unix(),
			}

			if _, err := br.Collection.InsertOne(ctx, bidEntityMongo); err != nil {
				logger.Error("Error trying to insert bid", err)
				return
			}
		}(bid)
	}

	wg.Wait()
	return nil
}
