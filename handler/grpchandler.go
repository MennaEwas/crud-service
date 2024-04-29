package handler

import (
	"context"
	"crudservice/MiddleWare"
	"crudservice/mapper"
	"crudservice/proto"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MyOrderServer struct {
	proto.UnimplementedOrderServer
}

func (s *MyOrderServer) CreateOrder(ctx context.Context, server *proto.OrderCreateRequest) (*proto.OrderResponse, error) {
	//implement the Creation Process
	body := mapper.CreateOrderRequest{
		Name:  server.Name,
		Price: server.Price,
	}
	res, err := MiddleWare.CreateOrderMiddleWare(ctx, body)

	if err != nil {
		failureCode := MapErrorToFailureCode(err)
		//return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		return ProtoOrderResponseFailure(err.Error(), failureCode), nil
	}
	tempId := res.InsertedID.(primitive.ObjectID)

	return &proto.OrderResponse{
		Result: &proto.OrderResponse_Success_{
			Success: &proto.OrderResponse_Success{
				Id:    tempId.Hex(),
				Name:  body.Name,
				Price: body.Price},
		},
	}, nil
}

func (s *MyOrderServer) UpdateOrder(ctx context.Context, server *proto.OrderUpdateRequest) (*proto.OrderResponse, error) {
	_id := server.GetId()
	newId, _ := primitive.ObjectIDFromHex(_id)

	body := mapper.CreateOrderRequest{
		Name:  server.GetName(),
		Price: server.GetPrice(),
	}

	_, err := MiddleWare.UpdateOrderMiddleWare(ctx, newId, body)
	if err != nil {
		failureCode := MapErrorToFailureCode(err)
		//return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		return ProtoOrderResponseFailure(err.Error(), failureCode), nil
	}
	return &proto.OrderResponse{
		Result: &proto.OrderResponse_Success_{
			Success: &proto.OrderResponse_Success{
				Id:    _id,
				Name:  server.GetName(),
				Price: server.GetPrice(),
			},
		},
	}, nil
}

func (s *MyOrderServer) ReadOrder(ctx context.Context, server *proto.OrderReadRequest) (*proto.OrderResponse, error) {
	//read all orders implementation in mongo
	_id := server.GetId()

	newId, _ := primitive.ObjectIDFromHex(_id)
	order, err := MiddleWare.ReadOrderMiddleWare(ctx, newId)

	if err != nil {
		failureCode := MapErrorToFailureCode(err)
		//return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		return ProtoOrderResponseFailure(err.Error(), failureCode), nil
	}
	fmt.Println(order)
	return &proto.OrderResponse{
		Result: &proto.OrderResponse_Success_{
			Success: &proto.OrderResponse_Success{
				Id:    _id,
				Name:  order.Name,
				Price: order.Price,
			},
		},
	}, nil
}

func (s *MyOrderServer) DeleteOrder(ctx context.Context, server *proto.OrderDeleteRequest) (*proto.OrderDeleteResponse, error) {
	_id := server.GetId()
	newId, _ := primitive.ObjectIDFromHex(_id)
	_, err, count := MiddleWare.DeleteOrderMiddleWare(ctx, newId)
	if err != nil {
		failureCode := MapErrorToFailureCode(err)
		//return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		return ProtoOrderDeleteResponseFailure(err.Error(), proto.OrderDeleteResponse_FailureCode(failureCode)), nil
	}

	//Ask Nasrallah
	if count == 0 {
		// need to add notfound
		return ProtoOrderDeleteResponseFailure("IDNOTFOUND", 404), nil
	}

	return &proto.OrderDeleteResponse{
		Result: &proto.OrderDeleteResponse_Success_{
			Success: &proto.OrderDeleteResponse_Success{},
		},
	}, nil
}
