package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrOrderNotFound is returned when the specified order cannot be found
	ErrOrderNotFound = errors.New("order not found")

	// ErrOrderExists is returned when a new order is being inserted but an order
	// with the same ID already exists
	ErrOrderExists = errors.New("order already exists")
)

////////////////////////////////////////////////////////////////////////////////

// GetOrder should return the order with the given ID. If that ID isn't found then
// the special ErrOrderNotFound error should be returned.
func (i *Instance) GetOrder(ctx context.Context, id string) (Order, error) {
	var order Order
	filter := bson.M{"id": id}
	err := i.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Order{}, ErrOrderNotFound
		}
		return Order{}, err
	}

	return order, nil
}

////////////////////////////////////////////////////////////////////////////////

// GetOrders should return all orders with the given status. If status is the
// special -1 value then it should return all orders regardless of their status.
func (i *Instance) GetOrders(ctx context.Context, status OrderStatus) ([]Order, error) {
	// TODO: get orders from DB based based on the status sent, unless status is -1
	var orders []Order
	var filter bson.M
	if status == -1 {
		filter = bson.M{}
	} else {
		filter = bson.M{"status": status}
	}
	cursor, err := i.collection.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return orders, ErrOrderNotFound
		}
		return orders, err
	}

	// Iterate over the cursor and decode each document into an Order
	for cursor.Next(ctx) {
		var order Order
		if err := cursor.Decode(&order); err != nil {
			return nil, fmt.Errorf("error decoding document: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

////////////////////////////////////////////////////////////////////////////////

// SetOrderStatus should update the order with the given ID and set the status
// field. If that ID isn't found then the special ErrOrderNotFound error should
// be returned.
func (i *Instance) SetOrderStatus(ctx context.Context, id string, status OrderStatus) error {
	// TODO: update the order's status field to status for the id
	filter := bson.D{{Key: "id", Value: id}}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: status}}}}
	result, err := i.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrOrderNotFound
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

// InsertOrder should fill in the order's ID with a unique identifier if it's not
// already set and then insert it into the database. It should return the order's
// ID. If the order already exists then ErrOrderExists should be returned.
func (i *Instance) InsertOrder(ctx context.Context, order Order) (string, error) {
	// TODO: if the order's ID field is empty, generate a random ID, then insert
	if order.ID == "" {
		id := uuid.New()
		order.ID = id.String()
	} else {
		// Check if document with the same ID already exists
		filter := bson.D{{Key: "id", Value: order.ID}}
		//filter := bson.M{"id": order.ID}
		err := i.collection.FindOne(ctx, filter).Decode(&Order{})
		if err == nil {
			return "", ErrOrderExists
		}
	}

	// If no document with the same ID exists, insert the new document
	_, err := i.collection.InsertOne(ctx, order)
	if err != nil {
		return "", fmt.Errorf("error inserting document: %e", err)
	}

	//originally added this type assertion and conversion but now I think this is overkill
	//insertedID, ok := result.InsertedID.(primitive.ObjectID)
	//if !ok {
	//	return "", errors.New("was not able to type assert the returned ID from mongo")
	//}

	return order.ID, nil
}
