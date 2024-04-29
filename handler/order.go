package handler

import (
	"crudservice/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Orderpiece struct {
	ID    string `bson:"_id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Price int32  `bson:"price" json:"price"`
}

type Response struct {
	Status int        `bson:"status" json:"status"`
	Order  Orderpiece `bson:"order" json:"order"`
}

type CreateOrderRequest struct {
	Name  string `json:"name"`
	Price int32  `json:"price"`
}

func ReadByID(c *gin.Context) {
	id := c.Param("id")
	_id, _ := primitive.ObjectIDFromHex(id)

	cursor := database.Orders.FindOne(c, bson.M{"_id": _id})

	var order Orderpiece

	err := cursor.Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no order found"})
		return
	}
	fmt.Println(order)

	res := Response{
		Status: http.StatusOK,
		Order:  order,
	}
	c.JSON(http.StatusOK, res)
	return
}

func CreateOrder(c *gin.Context) {
	var body CreateOrderRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid"})
		return
	}
	res, err := database.Orders.InsertOne(c, body)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "can't order"})
		return
	}

	temp_id := res.InsertedID.(primitive.ObjectID)
	order := Orderpiece{
		ID:    temp_id.Hex(),
		Name:  body.Name,
		Price: body.Price,
	}
	response := Response{
		Status: http.StatusCreated,
		Order:  order,
	}

	c.JSON(http.StatusCreated, response)
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	_id, _ := primitive.ObjectIDFromHex(id)

	var body struct {
		Price int32  `json:"price" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	// Update query corrected to search by "_id" instead of "id"
	_, err := database.Orders.UpdateOne(c, bson.M{"_id": _id}, bson.M{"$set": bson.M{"price": body.Price}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't update"})
		return
	}

	order := Orderpiece{
		ID:    _id.Hex(),
		Name:  body.Name, // Update to use Name from the request body
		Price: body.Price,
	}
	response := Response{
		Status: http.StatusOK,
		Order:  order,
	}
	c.JSON(http.StatusOK, response)
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	_id, _ := primitive.ObjectIDFromHex(id)

	res, _ := database.Orders.DeleteOne(c, bson.M{"_id": _id})
	if res.DeletedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": "Deletion Done Successfully"})
	return
}
