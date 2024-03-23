package graph

import "github.com/allanCordeiro/pos-fc-clean-arch/internal/usecases"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateOrderUseCase usecases.CreateOrderUseCase
	ListOrderUseCase   usecases.ListOrderUseCase
}
