package repository

import (
	"awesomeProject5/OMS/database"
	"awesomeProject5/OMS/orders/requests"
	"awesomeProject5/OMS/orders/responses"
	"context"
	"time"
)

type OrderService interface {
	//method create order in that interface
	CreateOrder(ctx context.Context, request *requests.CreateOrderSvcRequest) (*responses.CreateOrderSvcResponse, error)
}

func CreateOrder(ctx context.Context, request *requests.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Database("OMS_service").Collection("orders")

	_, err := collection.InsertOne(ctx, request)
	return err
}
