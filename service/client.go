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
	go OrderStreamGRPC(ctx, client)
	go OrderStatusStreamGRPC(ctx, client)

}
func OrderStreamGRPC(ctx context.Context, client orders.OrdersServerClient) {
	req := &orders.Empty{}
	stream, err := client.OrdersStream(ctx, req)

	if err != nil {
		fmt.Printf("New order stream Err: %v, Retry in %v \n", err, RetryTime)
		time.Sleep(RetryTime)
		OrderStreamGRPC(ctx, client)
	}

	for {
		resp, err := stream.Recv()
		if errors.As(err, &io.EOF) {
			fmt.Printf("New Order stream recv EDF %v \n", err)
			time.Sleep(RetryTime)
			OrderStreamGRPC(ctx, client)
		}
		if err != nil {
			fmt.Printf("New Order stream ERR %v \n", err)
		}
		if resp != nil {
			fmt.Printf("New Orde:  %v \n", resp)
		}
	}
}

func OrderStatusStreamGRPC(ctx context.Context, client orders.OrdersServerClient) {
	stream, err := client.OrderStatusStream(ctx)

	if err != nil {
		fmt.Printf("New order status stream Err: %v, Retry in %v \n", err, RetryTime)
		time.Sleep(RetryTime)
		OrderStatusStreamGRPC(ctx, client)
	}

	go func() {

		err := stream.Send(&orders.OrderRequest{
			Id: "Ok",
		})
		if err != nil {
			fmt.Printf("Send order Request:  %v \n", err)
		}
	}()

	for {
		resp, err := stream.Recv()
		if errors.As(err, &io.EOF) {
			fmt.Printf("New Order status stream recv EDF %v \n", err)
			time.Sleep(RetryTime)
			OrderStatusStreamGRPC(ctx, client)

			if err != nil {
				fmt.Printf("New Order stream ERR %v \n", err)
			}
			if resp != nil {
				fmt.Printf("New Order Stream:  %v \n", resp)
			}
		}
	}
}
