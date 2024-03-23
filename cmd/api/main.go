package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	configs "github.com/allanCordeiro/pos-fc-clean-arch/config"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/event"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/event/handler"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/database"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/grpc/pb"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/grpc/service"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/web"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/web/webserver"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/usecases"
	"github.com/allanCordeiro/pos-fc-clean-arch/pkg/events"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepository, eventCreated)
	webserver.AddHandler("POST", "/order", webOrderHandler.Create)
	webserver.AddHandler("GET", "/order", webOrderHandler.List)
	log.Printf("starting webserver on port %s", configs.WebServerPort)
	go webserver.Start()

	createOrderUseCase := usecases.NewCreateOrderUseCase(orderRepository, eventCreated, eventDispatcher)
	listOrderUseCase := usecases.NewListOrderUseCase(orderRepository)

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	reflection.Register(grpcServer)
	pb.RegisterOrderServiceServer(grpcServer, orderService)

	log.Printf("starting grpc server on port %s", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	grpcServer.Serve(lis)

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
