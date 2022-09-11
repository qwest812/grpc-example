package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc/external/orders"
	"io"
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
	//order := orders.OrderStatus{
	//	Id: "123",
	//}
	//if err := stream.Send(&order); err != nil {
	//	fmt.Printf("Error send Order, Err: %v", err)
	//}

	done := make(chan struct{})

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- struct{}{}
				return
			}
			if err != nil {
				fmt.Printf("Error Order status stream , Err: %v", err)
				continue
			}
			if resp != nil {
				fmt.Printf(resp.String())
			}
			status := &orders.OrderStatus{
				Id: "123",
			}
			err = stream.Send(status)
			if err != nil {
				fmt.Printf("Error Order status stream , Err: %v", err)
			}
		}
	}()
	//we are waiting EOF
	<-done
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
