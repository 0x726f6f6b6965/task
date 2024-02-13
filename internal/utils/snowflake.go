package utils

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

const (
	maxEpoch    uint64 = 1<<41 - 1
	maxNode     uint64 = 1<<8 - 1
	maxSequence uint64 = 1<<14 - 1
	shiftEpoch  uint8  = 22
	shiftNode   uint8  = 14
	// 2022-01-01 00:00:00
	baseEpoch uint64 = 1640966400000
)

var (
	// onceInitGenerator guarantee initialize generator only once
	onceInitGenerator sync.Once
	// rootGenerator - the root generator
	rootGenerator *generator
)

type Generator interface {
	// Next - get an unused sequence
	Next() (Sequence, error)
}

type generator struct {
	// nodeID is the node ID that the Snowflake generator will use for the next 8 bits
	nodeID uint64
	// sequence is the last 14 bits, usually an incremented number but can be anything. If set to 0, it will be random.
	sequence chan uint64
}

type Sequence interface {
	Uint64() uint64
	String() string
	Float64() float64
}

type sequence struct {
	num *big.Int
}

func NewGenerator(node uint64) (Generator, error) {
	if node > maxNode {
		return nil, fmt.Errorf("invalid node id; must be 0 ≤ id < %d", node)
	}
	// singleton
	onceInitGenerator.Do(func() {
		rootGenerator = &generator{
			nodeID:   node,
			sequence: make(chan uint64),
		}
		go func() {
			var (
				seq chan uint64
			)

			for {
				var reset <-chan time.Time
				if seq == nil {
					reset = time.After(time.Millisecond)
				}
				select {
				case <-reset:
					seq = make(chan uint64, 1)
					seq <- 0
				case current := <-seq:
					seq = nil
					for i := current; current <= maxSequence; i++ {
						rootGenerator.sequence <- i
					}
				}
			}
		}()
	})

	return rootGenerator, nil
}

func (g *generator) Next() (Sequence, error) {
	current := uint64(time.Now().UnixMilli())
	if (current - baseEpoch) > maxEpoch {
		return nil, fmt.Errorf("timestamp overflow")
	}

	seq := <-g.sequence

	nodeId := g.nodeID << shiftNode
	result := (current-baseEpoch)<<shiftEpoch + nodeId + seq
	num := big.NewInt(0)
	return &sequence{num: num.SetUint64(result)}, nil
}

func (s *sequence) String() string {
	return s.num.String()
}

func (s *sequence) Uint64() uint64 {
	return s.num.Uint64()
}

func (s *sequence) Float64() float64 {
	result, _ := s.num.Float64()
	return result
}
