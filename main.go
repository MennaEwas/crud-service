package main

import (
	"context"
	"crudservice/database"
	"crudservice/handler"
	"crudservice/proto"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

type myOrderServer struct {
	proto.UnimplementedOrderServer
}

func (s *myOrderServer) CreateOrder(ctx context.Context, server *proto.OrderRequest) (*proto.OrderResponse, error) {
	//implement the Creation Process
	body := handler.CreateOrderRequest{
		Name:  server.Name,
		Price: server.Price,
	}

	res, err := database.Orders.InsertOne(ctx, body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	fmt.Println(res)
	temp_id := res.InsertedID.(primitive.ObjectID)

	order := handler.Orderpiece{
		ID:    temp_id.Hex(),
		Name:  server.Name,
		Price: server.Price,
	}

	return &proto.OrderResponse{
		ServerMessage: fmt.Sprintf("Creation is Done for %v", order),
	}, nil
}

func (s *myOrderServer) UpdateOrder(ctx context.Context, server *proto.OrderRequest) (*proto.OrderResponse, error) {
	_id := server.GetId()
	new_id, _ := primitive.ObjectIDFromHex(_id)
	var body struct {
		Price int `json:"price" binding:"required"`
	}
	_, err := database.Orders.UpdateOne(ctx, bson.M{"id": new_id}, bson.M{"$set": bson.M{"price": body.Price}})
	if err != nil {

		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	return &proto.OrderResponse{
		ServerMessage: "Updating is Done",
	}, nil
}

func (s *myOrderServer) ReadOrder(ctx context.Context, server *proto.OrderRequest) (*proto.OrderResponse, error) {
	//read all orders implementation in mongo
	_id := server.GetId()

	new_id, _ := primitive.ObjectIDFromHex(_id)
	cursor := database.Orders.FindOne(ctx, bson.M{"_id": new_id})
	order := handler.Orderpiece{}

	err := cursor.Decode(&order)
	if err != nil {

		return nil, status.Errorf(codes.DataLoss, fmt.Sprintf("OrderNotFound: %v", err))
	}
	fmt.Println(order)
	return &proto.OrderResponse{
		ServerMessage: fmt.Sprintf("Reading is Done %v", order),
	}, nil
}
func (s *myOrderServer) DeleteOrder(ctx context.Context, server *proto.OrderRequest) (*proto.OrderResponse, error) {
	_id := server.GetId()
	new_id, _ := primitive.ObjectIDFromHex(_id)
	res, err := database.Orders.DeleteOne(ctx, bson.M{"_id": new_id})
	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.DataLoss, fmt.Sprintf("OrderNotFound: %v", err))
	}
	return &proto.OrderResponse{
		ServerMessage: "Deletion is Done",
	}, nil
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
	service := &myOrderServer{}
	//takes register and server
	proto.RegisterOrderServer(ServerRegister, service)

	//launch grpc server
	go func() {
		err = ServerRegister.Serve(lis)
		if err != nil {
			log.Fatalf("A7eeh %s", err)
		}
	}()

	//starting http server
	r := gin.Default()
	// create new functions to list all requests
	r.GET("/order", handler.ReadAllHandler)
	r.POST("/order", handler.CreateOrder)
	r.GET("/order/:id", handler.ReadByID)
	r.PATCH("/order/:id", handler.UpdateOrder)
	r.DELETE("/order/:id", handler.DeleteOrder)

	err = r.Run(":8085")
	if err != nil {
		fmt.Printf("Mooseeba %v", err)
	} else {
	}
	fmt.Printf("Server succesfully started on port :%v", lis.Addr())

}
