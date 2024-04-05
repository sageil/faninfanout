package efficient

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
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

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func Run_Efficient() {
	log.SetPrefix("Efficient: ")
	start := time.Now()

	done := make(chan int)
	defer close(done)

	randomintFetcher := func() int {
		return rand.Intn(int(1e9))
	}
	randStream := generate(done, randomintFetcher)
	cpuCount := MaxParallelism()
	primeFinderChans := make([]<-chan int, cpuCount)
	for i := 0; i < cpuCount; i++ {
		primeFinderChans[i] = primeIntFinder(done, randStream)
	}
	fannedInStream := fanIn(done, primeFinderChans...)
	for i := range take(done, fannedInStream, 10) {
		log.Println(i)
	}

	log.Println("Effient fanIn completed in", time.Since(start))
}

// fan in
func fanIn[T any](done <-chan int, chans ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)
	output := func(c <-chan T) {
		defer wg.Done()
		for v := range c {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}
	wg.Add(len(chans))
	for _, c := range chans {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
