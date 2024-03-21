package usecases

import (
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/entity"
	"github.com/allanCordeiro/pos-fc-clean-arch/pkg/events"
	"github.com/google/uuid"
)

type OrderInput struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutput struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(orderRepository entity.OrderRepositoryInterface,
	event events.EventInterface, eventDispatcher events.EventDispatcherInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
		OrderCreated:    event,
		EventDispatcher: eventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInput) (OrderOutput, error) {
	order, err := entity.NewOrder(uuid.NewString(), input.Price, input.Tax)
	if err != nil {
		return OrderOutput{}, err
	}
	err = order.CalculateFinalPrice()
	if err != nil {
		return OrderOutput{}, err
	}

	err = c.OrderRepository.Save(order)
	if err != nil {
		return OrderOutput{}, err
	}
	output := OrderOutput{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}
	c.OrderCreated.SetPayload(output)
	c.EventDispatcher.Dispatch(c.OrderCreated)

	return output, nil
}
