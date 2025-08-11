package clock

import "time"

// GameClock provides a fixed-step accumulator for core logic and a monotonically
// increasing elapsedMs for animation and time-based effects.
type GameClock struct {
	fixedStepMs   int64
	accumulatorMs int64
	elapsedMs     int64
	lastTime      time.Time
}

// NewGameClock returns a clock configured for a fixed step (e.g., 16ms ~ 60 updates/sec).
func NewGameClock(fixedStepMs int64) *GameClock {
	return &GameClock{
		fixedStepMs: fixedStepMs,
		lastTime:    time.Now(),
	}
}

// Advance accumulates real frame time. It returns how many fixed steps to process.
func (c *GameClock) Advance() (steps int) {
	now := time.Now()
	dt := now.Sub(c.lastTime)
	c.lastTime = now

	dtMs := dt.Milliseconds()
	if dtMs < 0 {
		dtMs = 0 // guard against clock skew
	}

	c.accumulatorMs += dtMs
	c.elapsedMs += dtMs

	for c.accumulatorMs >= c.fixedStepMs {
		c.accumulatorMs -= c.fixedStepMs
		steps++
	}
	return steps
}

// ElapsedMs returns total elapsed milliseconds since creation/reset.
func (c *GameClock) ElapsedMs() int64 {
	return c.elapsedMs
}

// Reset clears the accumulator and elapsed time.
func (c *GameClock) Reset() {
	c.accumulatorMs = 0
	c.elapsedMs = 0
	c.lastTime = time.Now()
}
