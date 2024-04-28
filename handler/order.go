package handler

import (
	"crudservice/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type Orderpiece struct {
	ID    string `bson:"_id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Price int32  `bson:"price" json:"price"`
}

type CreateOrderRequest struct {
	Name  string `json:"name"`
	Price int32  `json:"price"`
}

func ReadAllHandler(c *gin.Context) {
	cursor, err := database.Orders.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not order found"})
	}
	var orders []Orderpiece
	if err = cursor.All(c, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no order found"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func ReadByID(c *gin.Context) {
	id := c.Param("id")
	_id := id

	cursor := database.Orders.FindOne(c, bson.M{"_id": _id})

	order := Orderpiece{}

	err := cursor.Decode(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not order found"})
	}

	c.JSON(http.StatusOK, order)
	return
}

func CreateOrder(c *gin.Context) {
	var body CreateOrderRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	res, err := database.Orders.InsertOne(c, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't order"})
		return
	}
	order := Orderpiece{
		ID:    res.InsertedID.(string),
		Name:  body.Name,
		Price: body.Price,
	}
	c.JSON(http.StatusCreated, order)
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	_id := id

	var body struct {
		Price int `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	_, err := database.Orders.UpdateOne(c, bson.M{"id": _id}, bson.M{"$set": bson.M{"price": body.Price}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "order added successfully"})

	return
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	_id := id
	res, _ := database.Orders.DeleteOne(c, bson.M{"_id": _id})
	if res.DeletedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order not found"})
	}
	c.JSON(http.StatusOK, gin.H{"success": "Successful delete"})
	return
}
