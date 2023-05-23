package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

var (
	TFL                 = hitFakeEndpoint("TFL")
	StaticData          = hitFakeEndpoint("Static data")
	CitymapperAlgorithm = hitFakeEndpoint("Citymapper algorithm")
)

type Route func(ctx context.Context, query string) Result

type Result string

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
	defer cancel()

	results := ComplexQuerySystem(ctx, "London")

	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}

func FasterSearching(ctx context.Context, query string, replicates ...Route) Result {
	c := make(chan Result)

	searches := func(ctx context.Context, i int) { c <- replicates[i](ctx, query) }

	for i := range replicates {
		go searches(ctx, i)
	}
	return <-c
}

func ComplexQuerySystem(ctx context.Context, query string) (results []Result) {
	c := make(chan Result)
	go func() {
		c <- FasterSearching(ctx, query, TFL, TFL)
	}()
	go func() {
		c <- FasterSearching(ctx, query, StaticData, StaticData)
	}()
	go func() {
		c <- FasterSearching(ctx, query, CitymapperAlgorithm, CitymapperAlgorithm)
	}()

	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-ctx.Done():
			fmt.Println("timed out")
			return
		}
	}
	return
}

func hitFakeEndpoint(kind string) Route {
	return func(ctx context.Context, query string) Result {
		select {
		case <-ctx.Done():
			return Result(fmt.Sprintf("Request canceled for %s\n", kind))
		default:
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			return Result(fmt.Sprintf("%s result for %q\n", kind, query))
		}
	}
}
