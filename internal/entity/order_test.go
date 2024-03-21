package entity

import "testing"

func TestCreateOrders(t *testing.T) {
	scenarios := []struct {
		Name               string
		ID                 string
		Price              float64
		Tax                float64
		ExpectedFinalPrice float64
		ExpectedErr        error
	}{
		{
			Name:        "Given an empty ID when create a new order then should receive an error",
			ID:          "",
			Price:       5.55,
			Tax:         5.55,
			ExpectedErr: ErrInvalidID,
		},
		{
			Name:        "Given an empty Price when create a new order then should receive an error",
			ID:          "123",
			ExpectedErr: ErrInvalidPrice,
		},
		{
			Name:        "Given an empty Tax when create a new order then should receive an error",
			ID:          "123",
			Price:       5.55,
			ExpectedErr: ErrInvalidTax,
		},
		{
			Name:               "Given valid params when create a new order then should create order",
			ID:                 "123",
			Price:              5.55,
			Tax:                5.55,
			ExpectedFinalPrice: 11.1,
			ExpectedErr:        nil,
		},
	}

	for _, test := range scenarios {
		t.Run(test.Name, func(t *testing.T) {
			order := Order{
				ID:    test.ID,
				Price: test.Price,
				Tax:   test.Tax,
			}
			err := order.IsValid()
			if err != test.ExpectedErr {
				t.Errorf("expected %s but found %s", test.ExpectedErr, err)
			}
			if err == nil {
				err = order.CalculateFinalPrice()
				if err != nil {
					t.Errorf("expected no error but found %s", err)
				}
				if order.FinalPrice != test.ExpectedFinalPrice {
					t.Errorf("final price error. expected %f but found %f", test.ExpectedFinalPrice, order.FinalPrice)
				}
			}
		})
	}

}

func TestCreateNewOrderFunc(t *testing.T) {
	ID := "123"
	Price := 5.55
	Tax := 5.55
	ExpectedFinalPrice := 11.1

	t.Run("Given valid params when i call NewOrder function then should create order", func(t *testing.T) {
		order, err := NewOrder(ID, Price, Tax)
		if err != nil {
			t.Errorf("expected no error but found %s", err)
		}

		err = order.CalculateFinalPrice()
		if err != nil {
			t.Errorf("expected no error but found %s", err)
		}

		if order.FinalPrice != ExpectedFinalPrice {
			t.Errorf("expected final price error: expected %f but found %f", ExpectedFinalPrice, order.FinalPrice)
		}
		if order.ID != ID {
			t.Errorf("expected id error: expected %s but found %s", ID, order.ID)
		}
		if order.Price != Price {
			t.Errorf("expected price error: expected %f but found %f", Price, order.Price)
		}
		if order.Tax != Tax {
			t.Errorf("expected final price error: expected %f but found %f", Tax, order.Tax)
		}
	})
}
