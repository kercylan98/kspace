package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go A(ctx)

	select {
	case <-time.After(1 * time.Second):
		fmt.Println(ctx.Deadline())
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}

func A(ctx context.Context) {
	fmt.Println("A:RUN")
	select {
	case <-time.After(1 * time.Second):
		break
	case <-ctx.Done():
		return
	}
	fmt.Println("A:END")
	go B(ctx)
}

func B(ctx context.Context) {
	fmt.Println("B:RUN")
	select {
	case <-time.After(1 * time.Second):
		break
	case <-ctx.Done():
		return
	}
	fmt.Println("B:END")
	go C(ctx)
}

func C(ctx context.Context) {
	fmt.Println("C:RUN")
	select {
	case <-time.After(1 * time.Second):
		break
	case <-ctx.Done():
		return
	}
	fmt.Println("C:END")
}
