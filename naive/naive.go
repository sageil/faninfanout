package naive

import (
	"log"
	"math/rand"
	"time"
)

// generic func to generate random int and return a readonly channel
func generate[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

func take[T any, K any](done <-chan K, in <-chan T, n int) <-chan T {
	takeStream := make(chan T)
	go func() {
		defer close(takeStream)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-in:
			}
		}
	}()
	return takeStream
}

func primeIntFinder(done <-chan int, randomStream <-chan int) <-chan int {
	isPrime := func(rndInt int) bool {
		for i := rndInt - 1; i > 1; i-- {
			if rndInt%i == 0 {
				return false
			}
		}
		return true
	}

	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case randomInt := <-randomStream:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()
	return primes
}

func Run_Naive() {
	log.SetPrefix("Naive: ")
	start := time.Now()
	done := make(chan int)
	defer close(done)

	randomintFetcher := func() int {
		return rand.Intn(int(1e9))
	}
	randStream := generate(done, randomintFetcher)
	primeStream := primeIntFinder(done, randStream)

	// only take 10 out of the infinitly generating stream
	for i := range take(done, primeStream, 10) {
		log.Print(i)
	}
	log.Println("Finished in ", time.Since(start))
}
