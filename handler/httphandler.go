package handler

import (
	"crudservice/MiddleWare"
	"crudservice/mapper"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Response struct {
	Status int               `bson:"status" json:"status"`
	Order  mapper.Orderpiece `bson:"order" json:"order"`
}

func CreateOrder(c *gin.Context) {
	var body mapper.CreateOrderRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "can't order"})
		return
	}

	res, _ := MiddleWare.CreateOrderMiddleWare(c, body)

	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "can't order"})
		return
	}

	tempId := res.InsertedID.(primitive.ObjectID)
	order := mapper.Orderpiece{
		ID:    tempId.Hex(),
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

	var body mapper.CreateOrderRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	// Update query corrected to search by "_id" instead of "id"
	_, err1 := MiddleWare.UpdateOrderMiddleWare(c, _id, body)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't update"})
		return
	}

	order := mapper.Orderpiece{
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

func ReadByID(c *gin.Context) {
	id := c.Param("id")
	_id, _ := primitive.ObjectIDFromHex(id)

	order, err := MiddleWare.ReadOrderMiddleWare(c, _id)

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

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	_id, _ := primitive.ObjectIDFromHex(id)

	_, err, count := MiddleWare.DeleteOrderMiddleWare(c, _id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order not found"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": "Deletion Done Successfully"})
	return
}
