package common

import (
	"math/rand"
	"sync/atomic"
	"time"
)

// Generator is a snowflake generator.
// For compatibility with other snowflake implementations, both worker and PID are set,
// but they are randomized for every generator.
type IDGenerator struct {
	inc         *uint64
	worker, pid uint64
}

var defaultGenerator = NewIDGenerator(0, 0)

// NewIDGenerator creates a new ID generator with the given worker and pid.
// If worker or pid is empty, it will be set to a random number.
func NewIDGenerator(worker, pid uint64) *IDGenerator {
	if worker == 0 {
		worker = rand.Uint64()
	}
	if pid == 0 {
		pid = rand.Uint64()
	}

	g := &IDGenerator{
		inc:    new(uint64),
		worker: worker % 32,
		pid:    pid % 32,
	}

	return g
}

// GenerateID generates a new snowflake with the default generator.
// If you need to customize the worker and PID, manually call (*Generator).Generate.
func GenerateID() Snowflake {
	return defaultGenerator.Generate()
}

// GenerateID generates a new snowflake with the given time with the default generator.
// If you need to customize the worker and PID, manually call (*Generator).GenerateWithTime.
func GenerateIDWithTime(t time.Time) Snowflake {
	return defaultGenerator.GenerateWithTime(t)
}

// Generate generates a snowflake with the current time.
func (g *IDGenerator) Generate() Snowflake {
	return g.GenerateWithTime(time.Now())
}

// GenerateWithTime generates a snowflake with the given time.
// To generate a snowflake for comparison, use the top-level New function instead.
func (g *IDGenerator) GenerateWithTime(t time.Time) Snowflake {
	increment := atomic.AddUint64(g.inc, 1)
	ts := uint64(t.UnixMilli() - Epoch)

	worker := g.worker << 17
	pid := g.pid << 12

	return Snowflake(ts<<22 | worker | pid | (increment % 4096))
}
