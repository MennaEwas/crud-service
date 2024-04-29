package main

import (
	"crudservice/database"
	"crudservice/handler"
	"crudservice/proto"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
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

	//open connection with mongo
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

	//starting grpc
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Connect Create Listner %s", err)

	}
	//creating new server
	ServerRegister := grpc.NewServer()
	service := &handler.MyOrderServer{}
	//takes register and server
	proto.RegisterOrderServer(ServerRegister, service)
	go func() {
		err = ServerRegister.Serve(lis)
		if err != nil {
			log.Fatalf("A7eeh %s", err)
		}
	}()
	//starting http server
	Inithttp()

}
