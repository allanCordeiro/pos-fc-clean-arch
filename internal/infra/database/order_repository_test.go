package database

import (
	"database/sql"
	"testing"

	"github.com/allanCordeiro/pos-fc-clean-arch/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	_, err = db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.NoError(err)
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	suite.Db.Close()
}

func (suite *OrderRepositoryTestSuite) TestCreateOrder() {
	expectedFinalPrice := 12.0
	expectedOrder, _ := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(expectedOrder.CalculateFinalPrice())

	repo := NewOrderRepository(suite.Db)
	err := repo.Save(expectedOrder)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("SELECT id, price, tax, final_price FROM orders WHERE id = ?", expectedOrder.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)
	suite.NoError(err)
	suite.Equal(expectedOrder.ID, orderResult.ID)
	suite.Equal(expectedOrder.Price, orderResult.Price)
	suite.Equal(expectedOrder.Tax, orderResult.Tax)
	suite.Equal(expectedOrder.FinalPrice, expectedFinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestListOrders() {
	expectedLenght := 2
	order1, _ := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(order1.CalculateFinalPrice())
	order2, _ := entity.NewOrder("456", 10.0, 2.0)
	suite.NoError(order2.CalculateFinalPrice())

	repo := NewOrderRepository(suite.Db)
	suite.NoError(repo.Save(order1))
	suite.NoError(repo.Save(order2))
	orders, err := repo.ListAll()

	suite.NoError(err)
	suite.Equal(expectedLenght, len(orders))
	suite.Equal(order1.ID, orders[0].ID)
	suite.Equal(order2.ID, orders[1].ID)
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}
