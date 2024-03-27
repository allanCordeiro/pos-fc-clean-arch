package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	configs "github.com/allanCordeiro/pos-fc-clean-arch/config"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/entity"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/event"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/event/handler"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/database"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/graph"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/grpc/pb"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/grpc/service"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/web"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/infra/web/webserver"
	"github.com/allanCordeiro/pos-fc-clean-arch/internal/usecases"
	"github.com/allanCordeiro/pos-fc-clean-arch/pkg/events"
	_ "github.com/lib/pq"
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
		fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable",
			configs.DBDriver,
			configs.DBUser,
			configs.DBPassword,
			configs.DBHost,
			configs.DBName,
		))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	orderRepository := database.NewOrderRepository(db)

	rabbitMqChannel := getRabbitMQChannel(configs.RabbitMqHost,
		configs.RabbitMqPort,
		configs.RabbitMqUser,
		configs.RabbitMqPassword)
	eventDispatcher := events.NewEventDispatcher()
	eventCreated := event.NewOrderCreated()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMqChannel,
	})
	createOrderUseCase := usecases.NewCreateOrderUseCase(orderRepository, eventCreated, eventDispatcher)
	listOrderUseCase := usecases.NewListOrderUseCase(orderRepository)

	go startWebServer(configs.WebServerPort, eventDispatcher, orderRepository, eventCreated)
	go startGrpcServer(configs.GRPCServerPort, *createOrderUseCase, *listOrderUseCase)

	startGraphQlServer(configs.GraphQLServerPort, *createOrderUseCase, *listOrderUseCase)
}

func startWebServer(port string, ed events.EventDispatcherInterface,
	repository entity.OrderRepositoryInterface,
	event events.EventInterface) {

	webserver := webserver.NewWebServer(port)
	webOrderHandler := web.NewWebOrderHandler(ed, repository, event)
	webserver.AddHandler("POST", "/order", webOrderHandler.Create)
	webserver.AddHandler("GET", "/order", webOrderHandler.List)
	log.Printf("starting webserver on port %s", port)
	webserver.Start()
}

func startGrpcServer(port string,
	orderUseCase usecases.CreateOrderUseCase,
	listOrderUseCase usecases.ListOrderUseCase) {

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(orderUseCase, listOrderUseCase)
	reflection.Register(grpcServer)
	pb.RegisterOrderServiceServer(grpcServer, orderService)

	log.Printf("starting grpc server on port %s", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	grpcServer.Serve(lis)
}

func startGraphQlServer(port string,
	orderUseCase usecases.CreateOrderUseCase,
	listOrderUseCase usecases.ListOrderUseCase) {
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: orderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getRabbitMQChannel(host string, port string, user string, password string) *amqp.Channel {
	//usually rabbitMQ take some time to cold start. Not a sexy way to do this, but it works
	time.Sleep(5 * time.Second)
	rabbitMQUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	conn, err := amqp.Dial(rabbitMQUrl)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
