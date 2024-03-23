package service

import (
	"context"

	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/grpc/pb"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/usecases"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecases.CreateOrderUseCase
	ListOrderUseCase   usecases.ListOrderUseCase
}

func NewOrderService(createOrderUseCase usecases.CreateOrderUseCase,
	listOrderUseCase usecases.ListOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecases.OrderInput{
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrder(ctx context.Context, in *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	ouput, err := s.ListOrderUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var orderList []*pb.CreateOrderResponse

	for _, order := range ouput.Orders {
		orderList = append(orderList, &pb.CreateOrderResponse{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
	}
	return &pb.ListOrderResponse{OrderList: orderList}, nil
}
