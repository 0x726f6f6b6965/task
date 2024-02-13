package utils

import (
	"math/rand"
	"sync"
	"testing"
)

func TestNextMonotonic(t *testing.T) {
	gen, _ := NewGenerator(10)
	out := make([]string, 10000)

	for i := range out {
		seq, _ := gen.Next()
		out[i] = seq.String()
	}

	// ensure they are all distinct and increasing
	for i := range out[1:] {
		if out[i] >= out[i+1] {
			t.Fatal("bad entries:", out[i], out[i+1])
		}
	}
}

func TestMultiCall(t *testing.T) {
	gen, _ := NewGenerator(3)
	c := make(chan uint64)
	times := rand.Intn(100000) + 1000
	go func() {
		defer close(c)
		var wg sync.WaitGroup
		defer wg.Wait()
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < times; j++ {
					seq, _ := gen.Next()
					c <- seq.Uint64()
				}
			}()
		}
	}()
	show := map[uint64]bool{}
	for v := range c {
		if show[v] {
			t.Fatal("get repeat squence")
		}
		show[v] = true
	}
}

func BenchmarkCall(b *testing.B) {
	gen, _ := NewGenerator(7)
	c := make(chan uint64)
	go func() {
		for j := 0; j < b.N; j++ {
			seq, _ := gen.Next()
			c <- seq.Uint64()
		}
		close(c)
	}()
	show := map[uint64]bool{}
	for v := range c {
		if show[v] {
			b.Fatal("get repeat squence")
		}
		show[v] = true
	}
}
