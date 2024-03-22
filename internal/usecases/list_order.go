package usecases

import (
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/entity"
)

type ListOrdersOutput struct {
	Orders []OrderOutput `json:"orders"`
	Count  int64         `json:"count"`
}

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(orderRepository entity.OrderRepositoryInterface) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (l *ListOrderUseCase) Execute() (ListOrdersOutput, error) {

	output, err := l.OrderRepository.ListAll()
	if err != nil {
		return ListOrdersOutput{}, err
	}

	var listOrder []OrderOutput
	for _, order := range output {
		listOrder = append(listOrder, OrderOutput{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}
	count := len(listOrder)

	return ListOrdersOutput{
		Orders: listOrder,
		Count:  int64(count),
	}, nil
}
