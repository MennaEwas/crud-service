package main

import (
	"crudservice/database"
	"crudservice/handler"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func Inithttp() {

	r := gin.Default()
	// create new functions to list all requests
	r.POST("/order", handler.CreateOrder)
	r.GET("/order/:id", handler.ReadByID)
	r.PUT("/order/:id", handler.UpdateOrder)
	r.DELETE("/order/:id", handler.DeleteOrder)
	err := r.Run(":8085")
	if err != nil {
		fmt.Printf("Mooseeba %v", err)
		return
	} else {
	}
	fmt.Println("Server succesfully started on port 8085")
}

func main() {
	err2 := database.Init(os.Getenv("MONGO_URI"), "orderdb")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("Connected to mongo")

	defer func() {
		// this is related to database and should be included with the grpc
		err := database.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	Inithttp()
}
