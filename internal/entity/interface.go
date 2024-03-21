package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	ListAll() ([]Order, error)
}
