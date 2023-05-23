package main

import (
	"fmt"
	"math/rand"
	"sync"
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
	var wg sync.WaitGroup

	c := make(chan Result, 3)

	wg.Add(3)

	go func() {
		defer wg.Done()
		c <- TFL(query)
	}()
	go func() {
		defer wg.Done()
		c <- StaticData(query)
	}()
	go func() {
		defer wg.Done()
		c <- CitymapperAlgorithm(query)
	}()

	wg.Wait()
	close(c)

	for result := range c {
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
