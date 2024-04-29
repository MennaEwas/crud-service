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
	"log"
	"net"
	"os"
)

type myOrderServer struct {
	proto.UnimplementedOrderServer
}

// implement response failure
func protoOrderResponseFailure(err string, code proto.OrderResponse_FailureCode) *proto.OrderResponse {
	return &proto.OrderResponse{
		Result: &proto.OrderResponse_Failure_{
			Failure: &proto.OrderResponse_Failure{
				FailureMessage: err,
				FailureCode:    code,
			},
		},
	}
}

func mapErrorToFailureCode(err error) proto.OrderResponse_FailureCode {
	fmt.Println(err)
	return proto.OrderResponse_GENERAL_ERROR

}

func (s *myOrderServer) CreateOrder(ctx context.Context, server *proto.OrderCreateRequest) (*proto.OrderResponse, error) {
	//implement the Creation Process
	body := handler.CreateOrderRequest{
		Name:  server.Name,
		Price: server.Price,
	}

	_, err := database.Orders.InsertOne(ctx, body)
	if err != nil {
		failureCode := mapErrorToFailureCode(err)
		//return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		return protoOrderResponseFailure(err.Error(), failureCode), nil
	}
	//temp_id := res.InsertedID.(primitive.ObjectID)
	//order := handler.Orderpiece{
	//	ID:    temp_id.Hex(),
	//	Name:  server.Name,
	//	Price: server.Price,
	//}

	return &proto.OrderResponse{
		Result: &proto.OrderResponse_Success_{
			Success: &proto.OrderResponse_Success{},
		},
	}, nil
}

func (s *myOrderServer) UpdateOrder(ctx context.Context, server *proto.OrderUpdateRequest) (*proto.OrderResponse, error) {
	_id := server.GetId()
	newId, _ := primitive.ObjectIDFromHex(_id)
	var body struct {
		Price int `json:"price" binding:"required"`
	}
	_, err := database.Orders.UpdateOne(ctx, bson.M{"id": newId}, bson.M{"$set": bson.M{"price": body.Price}})
	if err != nil {
		failureCode := mapErrorToFailureCode(err)
		//return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		return protoOrderResponseFailure(err.Error(), failureCode), nil
	}
	return &proto.OrderResponse{
		Result: &proto.OrderResponse_Success_{
			Success: &proto.OrderResponse_Success{},
		},
	}, nil
}

func (s *myOrderServer) ReadOrder(ctx context.Context, server *proto.OrderReadRequest) (*proto.OrderResponse, error) {
	//read all orders implementation in mongo
	_id := server.GetId()

	newId, _ := primitive.ObjectIDFromHex(_id)
	cursor := database.Orders.FindOne(ctx, bson.M{"_id": newId})
	order := handler.Orderpiece{}

	err := cursor.Decode(&order)
	if err != nil {
		failureCode := mapErrorToFailureCode(err)
		//return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		return protoOrderResponseFailure(err.Error(), failureCode), nil
	}
	fmt.Println(order)
	return &proto.OrderResponse{
		Result: &proto.OrderResponse_Success_{
			Success: &proto.OrderResponse_Success{},
		},
	}, nil
}
func (s *myOrderServer) DeleteOrder(ctx context.Context, server *proto.OrderDeleteRequest) (*proto.OrderResponse, error) {
	_id := server.GetId()
	newId, _ := primitive.ObjectIDFromHex(_id)
	res, err := database.Orders.DeleteOne(ctx, bson.M{"_id": newId})
	if err != nil {
		failureCode := mapErrorToFailureCode(err)
		//return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		return protoOrderResponseFailure(err.Error(), failureCode), nil
	}

	//Ask Nasrallah
	if res.DeletedCount == 0 {
		// need to add notfound
		return protoOrderResponseFailure("IDNOTFOUND", 404), nil
	}

	return &proto.OrderResponse{
		Result: &proto.OrderResponse_Success_{
			Success: &proto.OrderResponse_Success{},
		},
	}, nil
}

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
	service := &myOrderServer{}
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
