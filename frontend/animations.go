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

// CarouselAnimator handles smooth spring-based animations for word transitions
type CarouselAnimator struct {
	// Springs for each element
	prevSpring    harmonica.Spring
	currentSpring harmonica.Spring
	nextSpring    harmonica.Spring

	// Positions (0.0 to 1.0 representing animation progress)
	PrevPos    float64
	CurrentPos float64
	NextPos    float64

	// Velocities for spring physics
	prevVel    float64
	currentVel float64
	nextVel    float64

	// Animation state
	IsAnimating bool
	frame       int
}

// NewCarouselAnimator creates a new animator for typing screen carousel
func NewCarouselAnimator() *CarouselAnimator {
	return &CarouselAnimator{
		// Snappy spring for smooth but quick animations
		prevSpring:    harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.6),
		currentSpring: harmonica.NewSpring(harmonica.FPS(60), 7.0, 0.5),
		nextSpring:    harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.6),
		// Start with everything in place
		PrevPos:     1.0,
		CurrentPos:  1.0,
		NextPos:     1.0,
		IsAnimating: false,
	}
}

// TriggerTransition starts the carousel animation when moving to next word
func (c *CarouselAnimator) TriggerTransition() {
	c.IsAnimating = true
	c.frame = 0

	// Previous word: start visible, will scroll up and fade out
	c.PrevPos = 0.0
	c.prevVel = 0.0

	// Current word: start below, will slide up into view
	c.CurrentPos = 0.0
	c.currentVel = 0.0

	// Next word: start hidden below, will fade in
	c.NextPos = 0.0
	c.nextVel = 0.0
}

// Update advances all springs by one frame
func (c *CarouselAnimator) Update() {
	if !c.IsAnimating {
		return
	}

	c.frame++

	// Update all springs toward target position of 1.0
	c.PrevPos, c.prevVel = c.prevSpring.Update(c.PrevPos, c.prevVel, 1.0)
	c.CurrentPos, c.currentVel = c.currentSpring.Update(c.CurrentPos, c.currentVel, 1.0)

	// Next word starts slightly delayed for stagger effect
	if c.frame > 2 {
		c.NextPos, c.nextVel = c.nextSpring.Update(c.NextPos, c.nextVel, 1.0)
	}

	// Check if animation is complete
	if c.PrevPos > 0.98 && c.CurrentPos > 0.98 && c.NextPos > 0.98 {
		if abs(c.prevVel) < 0.01 && abs(c.currentVel) < 0.01 && abs(c.nextVel) < 0.01 {
			c.IsAnimating = false
			c.PrevPos = 1.0
			c.CurrentPos = 1.0
			c.NextPos = 1.0
		}
	}
}

// GetPrevOffset returns the vertical offset for the previous word (scrolls up)
func (c *CarouselAnimator) GetPrevOffset() int {
	// Starts at bottom of its area (offset 2), scrolls up to final position (offset 0)
	maxOffset := 2
	return int(float64(maxOffset) * (1.0 - c.PrevPos))
}

// GetPrevOpacity returns opacity for previous word (0.0 to 1.0)
func (c *CarouselAnimator) GetPrevOpacity() float64 {
	// Fades from 0 to target opacity as it scrolls up
	return c.PrevPos * 0.5 // Max 50% opacity for dimmed previous word
}

// GetCurrentOffset returns the vertical offset for current word (slides up from below)
func (c *CarouselAnimator) GetCurrentOffset() int {
	// Starts below (offset 3), slides up to center (offset 0)
	maxOffset := 3
	return int(float64(maxOffset) * (1.0 - c.CurrentPos))
}

// GetCurrentScale returns a scale factor for the current word (grows as it enters)
func (c *CarouselAnimator) GetCurrentScale() float64 {
	// Starts at 0.7, grows to 1.0
	return 0.7 + (c.CurrentPos * 0.3)
}

// GetNextOffset returns the vertical offset for next word (fades in below)
func (c *CarouselAnimator) GetNextOffset() int {
	// Starts further below (offset 2), moves up slightly (offset 0)
	maxOffset := 2
	return int(float64(maxOffset) * (1.0 - c.NextPos))
}

// GetNextOpacity returns opacity for next word (0.0 to 1.0)
func (c *CarouselAnimator) GetNextOpacity() float64 {
	// Fades in to target opacity
	return c.NextPos * 0.6 // Max 60% opacity for next word preview
}

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
		a.positions[i] = 0 // Start at 0 (invisible)
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
