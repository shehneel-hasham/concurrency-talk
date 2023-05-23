package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	TFL                 = hitFakeEndpoint("TFL")
	StaticData          = hitFakeEndpoint("Static data")
	CitymapperAlgorithm = hitFakeEndpoint("Citymapper algorithm")
)

type Route func(query string) Result

type Result string

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()

	results := ComplexQuerySystem("London")

	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}

func ComplexQuerySystem(query string) (results []Result) {
	c := make(chan Result)
	go func() {
		c <- TFL(query)
	}()
	go func() {
		c <- StaticData(query)
	}()
	go func() {
		c <- CitymapperAlgorithm(query)
	}()

	for i := 0; i < 3; i++ {
		result := <-c
		results = append(results, result)
	}
	return
}

func hitFakeEndpoint(kind string) Route {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}
