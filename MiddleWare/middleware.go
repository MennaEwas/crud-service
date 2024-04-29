package MiddleWare

import (
	"context"
	"crudservice/database"
	"crudservice/mapper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateOrderMiddleWare(ctx context.Context, body mapper.CreateOrderRequest) (*mongo.InsertOneResult, error) {
	res, err := database.Orders.InsertOne(ctx, body)

	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateOrderMiddleWare(c context.Context, id primitive.ObjectID, body mapper.CreateOrderRequest) (*mongo.UpdateResult, error) {

	res, err := database.Orders.UpdateOne(c, bson.M{"_id": id}, bson.M{"$set": bson.M{"price": body.Price}})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ReadOrderMiddleWare(c context.Context, id primitive.ObjectID) (mapper.Orderpiece, error) {
	cursor := database.Orders.FindOne(c, bson.M{"_id": id})
	var order mapper.Orderpiece
	err := cursor.Decode(&order)
	if err != nil {
		return order, err

	}
	return order, nil

}

func DeleteOrderMiddleWare(c context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error, int32) {
	res, _ := database.Orders.DeleteOne(c, bson.M{"_id": id})
	if res.DeletedCount == 0 {
		return nil, nil, 0
	}
	return res, nil, 1
}
