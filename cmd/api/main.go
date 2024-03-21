package main

import (
	"database/sql"
	"fmt"
	"log"

	configs "github.com/allanCordeiro/pos-fc-clean-arch/config"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/event"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/event/handler"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/database"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/web"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/web/webserver"
	"github.com/allanCordeiro/pos-fc-clean-arch/pkg/events"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			configs.DBUser,
			configs.DBPassword,
			configs.DBHost,
			configs.DBPort,
			configs.DBName,
		))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	orderRepository := database.NewOrderRepository(db)
	rabbitMqChannel := getRabbitMQChannel()
	eventDispatcher := events.NewEventDispatcher()
	eventCreated := event.NewOrderCreated()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMqChannel,
	})

	//createOrderUseCase := usecases.NewCreateOrderUseCase(orderRepository, eventCreated, eventDispatcher)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepository, eventCreated)
	webserver.AddHandler("POST", "/order", webOrderHandler.Create)
	log.Printf("starting webserver on port %s", configs.WebServerPort)
	webserver.Start()
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") //KLUDGE::colocar isso no config
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}