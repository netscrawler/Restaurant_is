package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	orderv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type testEnv struct {
	client orderv1.OrderServiceClient
	conn   *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func setupTestEnv(t *testing.T) *testEnv {
	// Setup gRPC client
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	client := orderv1.NewOrderServiceClient(conn)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)

	return &testEnv{
		client: client,
		conn:   conn,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (env *testEnv) cleanup() {
	env.cancel()
	env.conn.Close()
}

func TestCreateOrder(t *testing.T) {
	env := setupTestEnv(t)
	defer env.cleanup()

	// Create test data
	userID := uuid.New()
	dishID := uuid.New()
	orderType := orderv1.OrderType_ORDER_TYPE_DELIVERY
	deliveryAddress := []byte("Test Address 123")

	// Create order request
	req := &orderv1.CreateOrderRequest{
		UserId: &orderv1.UUID{
			Value: userID.String(),
		},
		OrderType:       orderType,
		DeliveryAddress: deliveryAddress,
		Items: []*orderv1.OrderItemCreation{
			{
				DishId: &orderv1.UUID{
					Value: dishID.String(),
				},
				Quantity: 2,
			},
		},
	}

	// Call CreateOrder
	resp, err := env.client.CreateOrder(env.ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.NotNil(t, resp.Id)
	assert.Equal(t, "created", resp.Status)
	assert.Positive(t, resp.TotalAmount)

	// Verify order ID is valid UUID
	orderID, err := uuid.Parse(resp.Id.Value)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, orderID)
}

func TestCreateOrderWithInvalidDish(t *testing.T) {
	env := setupTestEnv(t)
	defer env.cleanup()

	// Create test data with non-existent dish
	userID := uuid.New()
	nonExistentDishID := uuid.New()
	orderType := orderv1.OrderType_ORDER_TYPE_DELIVERY
	deliveryAddress := []byte("Test Address 123")

	// Create order request
	req := &orderv1.CreateOrderRequest{
		UserId: &orderv1.UUID{
			Value: userID.String(),
		},
		OrderType:       orderType,
		DeliveryAddress: deliveryAddress,
		Items: []*orderv1.OrderItemCreation{
			{
				DishId: &orderv1.UUID{
					Value: nonExistentDishID.String(),
				},
				Quantity: 1,
			},
		},
	}

	// Call CreateOrder
	resp, err := env.client.CreateOrder(env.ctx, req)
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestCreateOrderWithMultipleItems(t *testing.T) {
	env := setupTestEnv(t)
	defer env.cleanup()

	// Create test data
	userID := uuid.New()
	dishID1 := uuid.New()
	dishID2 := uuid.New()
	orderType := orderv1.OrderType_ORDER_TYPE_DELIVERY
	deliveryAddress := []byte("Test Address 123")

	// Create order request with multiple items
	req := &orderv1.CreateOrderRequest{
		UserId: &orderv1.UUID{
			Value: userID.String(),
		},
		OrderType:       orderType,
		DeliveryAddress: deliveryAddress,
		Items: []*orderv1.OrderItemCreation{
			{
				DishId: &orderv1.UUID{
					Value: dishID1.String(),
				},
				Quantity: 2,
			},
			{
				DishId: &orderv1.UUID{
					Value: dishID2.String(),
				},
				Quantity: 3,
			},
		},
	}

	// Call CreateOrder
	resp, err := env.client.CreateOrder(env.ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.NotNil(t, resp.Id)
	assert.Equal(t, "created", resp.Status)
	assert.Positive(t, resp.TotalAmount)

	// Verify order ID is valid UUID
	orderID, err := uuid.Parse(resp.Id.Value)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, orderID)
}
