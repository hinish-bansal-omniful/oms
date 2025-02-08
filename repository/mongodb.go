package repository

import (
	"context"
	"oms-service/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOrders(ctx context.Context, orders []*domain.Order, db *mongo.Client) error {
	collection := db.Database("oms-service-db").Collection("orders")

	documents := make([]interface{}, len(orders))
	for i, order := range orders {
		documents[i] = order
	}

	_, err := collection.InsertMany(ctx, documents)
	return err
}
