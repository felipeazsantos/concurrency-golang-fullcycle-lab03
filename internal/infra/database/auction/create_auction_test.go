package auction

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	auctionentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/auction_entity"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockCollection struct {
	mu           sync.Mutex
	insertCalled bool
	updateCalled bool
	insertError  bool
	updateError  bool
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	m.mu.Lock()
	m.insertCalled = true
	m.mu.Unlock()

	if m.insertError {
		return nil, assert.AnError
	}
	return &mongo.InsertOneResult{}, nil
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	m.mu.Lock()
	m.updateCalled = true
	m.mu.Unlock()

	if m.updateError {
		return nil, assert.AnError
	}
	return &mongo.UpdateResult{}, nil
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return nil
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return nil, nil
}

func TestCreateAuction_GoroutineClosesAuction(t *testing.T) {
	os.Setenv("AUCTION_INTERVAL", "100ms")
	defer os.Unsetenv("AUCTION_INTERVAL")

	mockCollection := &MockCollection{}

	auctionRepo := NewAuctionRepositoryWithCollection(mockCollection)

	auctionEntity := &auctionentity.Auction{
		Id:          "test-auction-id",
		ProductName: "Test Product",
		Category:    "Electronics",
		Description: "Test Description",
		Condition:   auctionentity.New,
		Status:      auctionentity.Active,
		Timestamp:   time.Now(),
	}

	err := auctionRepo.CreateAuction(context.Background(), auctionEntity)
	assert.Nil(t, err)

	mockCollection.mu.Lock()
	insertCalled := mockCollection.insertCalled
	mockCollection.mu.Unlock()

	assert.True(t, insertCalled, "Expected InsertOne to be called")

	time.Sleep(150 * time.Millisecond)

	mockCollection.mu.Lock()
	updateCalled := mockCollection.updateCalled
	mockCollection.mu.Unlock()

	assert.True(t, updateCalled, "Expected UpdateOne to be called by goroutine")
}

func TestCreateAuction_GoroutineHandlesError(t *testing.T) {
	os.Setenv("AUCTION_INTERVAL", "50ms")
	defer os.Unsetenv("AUCTION_INTERVAL")

	mockCollection := &MockCollection{
		updateError: true,
	}

	auctionRepo := NewAuctionRepositoryWithCollection(mockCollection)

	auctionEntity := &auctionentity.Auction{
		Id:          "test-auction-error",
		ProductName: "Test Product Error",
		Category:    "Books",
		Description: "Test Description Error",
		Condition:   auctionentity.Used,
		Status:      auctionentity.Active,
		Timestamp:   time.Now(),
	}

	err := auctionRepo.CreateAuction(context.Background(), auctionEntity)
	assert.Nil(t, err)

	time.Sleep(100 * time.Millisecond)

	mockCollection.mu.Lock()
	updateCalled := mockCollection.updateCalled
	mockCollection.mu.Unlock()

	assert.True(t, updateCalled, "Expected UpdateOne to be called by goroutine even if it fails")
}
