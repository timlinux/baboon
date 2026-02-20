package frontend

import (
	"strings"
	"time"

	"github.com/charmbracelet/harmonica"
)

// Animation configuration constants
const (
	NumAnimatedRows   = 25 // Number of rows to animate on results screen
	AnimationInterval = 50 * time.Millisecond
	StaggerDelay      = 3 // Frames between each row starting
)

// Animator handles spring-based animations for the results screen
type Animator struct {
	springs    []harmonica.Spring
	positions  []float64
	velocities []float64
	frame      int
}

// NewAnimator creates a new animator for results screen animations
func NewAnimator() *Animator {
	a := &Animator{
		springs:    make([]harmonica.Spring, NumAnimatedRows),
		positions:  make([]float64, NumAnimatedRows),
		velocities: make([]float64, NumAnimatedRows),
		frame:      0,
	}

	// Create springs with a bouncy feel
	for i := range a.springs {
		a.springs[i] = harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.5)
		a.positions[i] = 0  // Start at 0 (invisible)
		a.velocities[i] = 0
	}

	return a
}

// Update advances all active springs by one frame
func (a *Animator) Update() {
	a.frame++

	for i := range a.springs {
		// Stagger the start: each row starts staggerDelay frames after the previous
		startFrame := i * StaggerDelay
		if a.frame >= startFrame {
			// Target is 1.0 (fully visible)
			a.positions[i], a.velocities[i] = a.springs[i].Update(
				a.positions[i], a.velocities[i], 1.0,
			)
		}
	}
}

// IsComplete returns true if all animations have finished
func (a *Animator) IsComplete() bool {
	for i := range a.positions {
		// Check if position is close to target (1.0) and velocity is near zero
		if a.positions[i] < 0.99 || abs(a.velocities[i]) > 0.01 {
			return false
		}
	}
	return true
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// ApplyAnimation applies slide-in animation to a row based on position (0-1)
func (a *Animator) ApplyAnimation(row string, animIdx int) string {
	if animIdx >= len(a.positions) {
		return row
	}

	pos := a.positions[animIdx]
	if pos >= 0.99 {
		return row // Fully visible
	}
	if pos <= 0.01 {
		return "" // Hidden
	}

	// Slide in from right: offset decreases as position increases
	maxOffset := 40
	offset := int(float64(maxOffset) * (1.0 - pos))
	if offset > 0 {
		return strings.Repeat(" ", offset) + row
	}
	return row
}

// GetInterval returns the animation tick interval
func GetAnimationInterval() time.Duration {
	return AnimationInterval
}
