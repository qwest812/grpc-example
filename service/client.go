package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc/external/orders"
	"io"
	"time"
)

const RetryTime = time.Second * 5

func main() {
	conn, err := grpc.Dial("127.0.0.1:5443", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Error conn Err: %v", err)
	}
	client := orders.NewOrdersServerClient(conn)
	ctx := context.Background()
	ConnectGRPC(ctx, client)

}
func ConnectGRPC(ctx context.Context, client orders.OrdersServerClient) {
	req := &orders.Empty{}
	stream, err := client.OrdersStream(ctx, req)

	if err != nil {
		fmt.Printf("New order stream Err: %v, Retry in %v \n", err, RetryTime)
		time.Sleep(RetryTime)
		ConnectGRPC(ctx, client)
	}

	for {
		resp, err := stream.Recv()
		if errors.As(err, &io.EOF) {
			fmt.Printf("New Order stream recv EDF %v \n", err)
			time.Sleep(RetryTime)
			ConnectGRPC(ctx, client)
		}
		if err != nil {
			fmt.Printf("New Order stream ERR %v \n", err)
		}
		if resp != nil {
			fmt.Printf("New Orde:  %v \n", resp)
		}
	}
}
