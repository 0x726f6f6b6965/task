package utils

import (
	"fmt"
	"math/big"
	"sync/atomic"
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

type Generator interface {
	// Next - get an unused sequence
	Next() (Sequence, error)
}

type generator struct {
	// epoch is the first 41 bits
	epoch uint64
	// nodeID is the node ID that the Snowflake generator will use for the next 8 bits
	nodeID uint64
	// sequence is the last 14 bits, usually an incremented number but can be anything. If set to 0, it will be random.
	sequence uint64
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

	return &generator{
		nodeID: node,
	}, nil
}

func (g *generator) Next() (Sequence, error) {
	current := uint64(time.Now().UnixMilli())
	if current > maxEpoch {
		return nil, fmt.Errorf("timestamp overflow")
	}

	prev := atomic.LoadUint64(&g.epoch)
	var seq uint64 = 0
	if current > prev {
		atomic.StoreUint64(&g.sequence, 0)
		atomic.StoreUint64(&g.epoch, current)
	} else {
		seq = atomic.AddUint64(&g.sequence, 1)
	}

	if seq > maxSequence {
		return nil, fmt.Errorf("sequence overflow")
	}
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
