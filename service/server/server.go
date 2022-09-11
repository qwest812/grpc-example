package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc/external/orders"
	"net"
	"time"
)

func main() {

	lis, err := net.Listen("tcp", "127.0.0.1:5443")
	if err != nil {
		fmt.Printf("%v", err)
	}

	grpcServer := grpc.NewServer()
	orders.RegisterOrdersServerServer(grpcServer, NewServer())

	for {
		err = grpcServer.Serve(lis)
		if err != nil {
			fmt.Printf("Error GRPC Serve: %v", err)
		}
		time.Sleep(time.Second * 5)
	}

}

func NewServer() *server {
	return &server{}
}

type server struct {
	orders.UnimplementedOrdersServerServer
}

func (s server) OrderStatusStream(stream orders.OrdersServer_OrderStatusStreamServer) error {
	fmt.Println("OrderStatusStream established")

	return nil
}
func (s server) OrdersStream(req *orders.Empty, stream orders.OrdersServer_OrdersStreamServer) error {
	fmt.Println("OrdersStream established")
	order := orders.NewOrder{
		Id:   "123",
		Type: "order",
	}
	if err := stream.Send(&order); err != nil {
		fmt.Printf("Error send Order, Err: %v", err)
	}

	return nil
}
