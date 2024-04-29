package handler

import (
	"crudservice/proto"
	"fmt"
)

// implement response failure
func ProtoOrderResponseFailure(err string, code proto.OrderResponse_FailureCode) *proto.OrderResponse {
	return &proto.OrderResponse{
		Result: &proto.OrderResponse_Failure_{
			Failure: &proto.OrderResponse_Failure{
				FailureMessage: err,
				FailureCode:    code,
			},
		},
	}
}

func ProtoOrderDeleteResponseFailure(err string, code proto.OrderDeleteResponse_FailureCode) *proto.OrderDeleteResponse {
	return &proto.OrderDeleteResponse{
		Result: &proto.OrderDeleteResponse_Failure_{
			Failure: &proto.OrderDeleteResponse_Failure{
				FailureMessage: err,
				FailureCode:    code,
			},
		},
	}
}

func MapErrorToFailureCode(err error) proto.OrderResponse_FailureCode {
	fmt.Println(err)
	return proto.OrderResponse_GENERAL_ERROR

}
