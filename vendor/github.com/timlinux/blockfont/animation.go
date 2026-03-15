// SPDX-FileCopyrightText: 2026 Tim Sutton / Kartoza
// SPDX-License-Identifier: MIT

package blockfont

import (
	"time"

	"github.com/charmbracelet/harmonica"
)

// AnimationInterval is the default tick interval for animations
const AnimationInterval = 50 * time.Millisecond

// TransitionType represents different animation transition types
type TransitionType int

const (
	// TransitionSlideUp slides content up
	TransitionSlideUp TransitionType = iota
	// TransitionSlideDown slides content down
	TransitionSlideDown
	// TransitionFadeIn fades content in
	TransitionFadeIn
	// TransitionFadeOut fades content out
	TransitionFadeOut
	// TransitionScale scales content
	TransitionScale
)

// Animator handles spring-based animations for smooth transitions
type Animator struct {
	// Springs for animation physics
	positionSpring harmonica.Spring
	opacitySpring  harmonica.Spring
	scaleSpring    harmonica.Spring

	// Current values (0.0 to 1.0 representing animation progress)
	Position float64
	Opacity  float64
	Scale    float64

	// Velocities for spring physics
	positionVel float64
	opacityVel  float64
	scaleVel    float64

	// Target values
	targetPosition float64
	targetOpacity  float64
	targetScale    float64

	// Animation state
	IsAnimating bool
	frame       int
}

// NewAnimator creates a new animator with default spring settings
func NewAnimator() *Animator {
	return &Animator{
		// Snappy springs for smooth but quick animations
		positionSpring: harmonica.NewSpring(harmonica.FPS(60), 7.0, 0.5),
		opacitySpring:  harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.6),
		scaleSpring:    harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.6),
		// Start at rest
		Position:       1.0,
		Opacity:        1.0,
		Scale:          1.0,
		targetPosition: 1.0,
		targetOpacity:  1.0,
		targetScale:    1.0,
		IsAnimating:    false,
	}
}

// TriggerTransition starts an animation of the specified type
func (a *Animator) TriggerTransition(t TransitionType) {
	a.IsAnimating = true
	a.frame = 0

	switch t {
	case TransitionSlideUp:
		a.Position = 0.0
		a.positionVel = 0.0
		a.targetPosition = 1.0
	case TransitionSlideDown:
		a.Position = 1.0
		a.positionVel = 0.0
		a.targetPosition = 0.0
	case TransitionFadeIn:
		a.Opacity = 0.0
		a.opacityVel = 0.0
		a.targetOpacity = 1.0
	case TransitionFadeOut:
		a.Opacity = 1.0
		a.opacityVel = 0.0
		a.targetOpacity = 0.0
	case TransitionScale:
		a.Scale = 0.7
		a.scaleVel = 0.0
		a.targetScale = 1.0
	}
}

// Update advances all springs by one frame.
// Returns true if animation is still in progress.
func (a *Animator) Update() bool {
	if !a.IsAnimating {
		return false
	}

	a.frame++

	// Update all springs toward their targets
	a.Position, a.positionVel = a.positionSpring.Update(a.Position, a.positionVel, a.targetPosition)
	a.Opacity, a.opacityVel = a.opacitySpring.Update(a.Opacity, a.opacityVel, a.targetOpacity)
	a.Scale, a.scaleVel = a.scaleSpring.Update(a.Scale, a.scaleVel, a.targetScale)

	// Check if animation is complete
	if a.isNearTarget(a.Position, a.targetPosition) &&
		a.isNearTarget(a.Opacity, a.targetOpacity) &&
		a.isNearTarget(a.Scale, a.targetScale) {
		if abs(a.positionVel) < 0.01 && abs(a.opacityVel) < 0.01 && abs(a.scaleVel) < 0.01 {
			a.IsAnimating = false
			a.Position = a.targetPosition
			a.Opacity = a.targetOpacity
			a.Scale = a.targetScale
		}
	}

	return a.IsAnimating
}

// isNearTarget checks if a value is close to its target
func (a *Animator) isNearTarget(current, target float64) bool {
	return abs(current-target) < 0.02
}

// GetOffset returns the vertical offset based on position (for slide animations)
func (a *Animator) GetOffset(maxOffset int) int {
	return int(float64(maxOffset) * (1.0 - a.Position))
}

// GetOpacityLevel returns opacity as a value from 0.0 to maxOpacity
func (a *Animator) GetOpacityLevel(maxOpacity float64) float64 {
	return a.Opacity * maxOpacity
}

// GetScaleFactor returns scale factor (e.g., 0.7 to 1.0)
func (a *Animator) GetScaleFactor(minScale float64) float64 {
	return minScale + (a.Scale * (1.0 - minScale))
}

// Reset stops any animation and resets to default state
func (a *Animator) Reset() {
	a.IsAnimating = false
	a.Position = 1.0
	a.Opacity = 1.0
	a.Scale = 1.0
	a.positionVel = 0.0
	a.opacityVel = 0.0
	a.scaleVel = 0.0
	a.targetPosition = 1.0
	a.targetOpacity = 1.0
	a.targetScale = 1.0
	a.frame = 0
}

// GetAnimationInterval returns the recommended tick interval
func GetAnimationInterval() time.Duration {
	return AnimationInterval
}

// WordCarouselAnimator handles the three-word carousel animation
// Used for speed reading applications
type WordCarouselAnimator struct {
	// Individual springs for each word
	prevSpring    harmonica.Spring
	currentSpring harmonica.Spring
	nextSpring    harmonica.Spring

	// Positions (0.0 to 1.0)
	PrevPos    float64
	CurrentPos float64
	NextPos    float64

	// Velocities
	prevVel    float64
	currentVel float64
	nextVel    float64

	// State
	IsAnimating bool
	frame       int
}

// NewWordCarouselAnimator creates a new carousel animator
func NewWordCarouselAnimator() *WordCarouselAnimator {
	return &WordCarouselAnimator{
		prevSpring:    harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.6),
		currentSpring: harmonica.NewSpring(harmonica.FPS(60), 7.0, 0.5),
		nextSpring:    harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.6),
		PrevPos:       1.0,
		CurrentPos:    1.0,
		NextPos:       1.0,
		IsAnimating:   false,
	}
}

// TriggerTransition starts the word carousel animation
func (w *WordCarouselAnimator) TriggerTransition() {
	w.IsAnimating = true
	w.frame = 0

	w.PrevPos = 0.0
	w.prevVel = 0.0
	w.CurrentPos = 0.0
	w.currentVel = 0.0
	w.NextPos = 0.0
	w.nextVel = 0.0
}

// Update advances all springs by one frame
func (w *WordCarouselAnimator) Update() bool {
	if !w.IsAnimating {
		return false
	}

	w.frame++

	w.PrevPos, w.prevVel = w.prevSpring.Update(w.PrevPos, w.prevVel, 1.0)
	w.CurrentPos, w.currentVel = w.currentSpring.Update(w.CurrentPos, w.currentVel, 1.0)

	// Next word starts slightly delayed for stagger effect
	if w.frame > 2 {
		w.NextPos, w.nextVel = w.nextSpring.Update(w.NextPos, w.nextVel, 1.0)
	}

	// Check completion
	if w.PrevPos > 0.98 && w.CurrentPos > 0.98 && w.NextPos > 0.98 {
		if abs(w.prevVel) < 0.01 && abs(w.currentVel) < 0.01 && abs(w.nextVel) < 0.01 {
			w.IsAnimating = false
			w.PrevPos = 1.0
			w.CurrentPos = 1.0
			w.NextPos = 1.0
		}
	}

	return w.IsAnimating
}

// GetPrevOffset returns vertical offset for previous word
func (w *WordCarouselAnimator) GetPrevOffset() int {
	maxOffset := 2
	return int(float64(maxOffset) * (1.0 - w.PrevPos))
}

// GetPrevOpacity returns opacity for previous word (max 50%)
func (w *WordCarouselAnimator) GetPrevOpacity() float64 {
	return w.PrevPos * 0.5
}

// GetCurrentOffset returns vertical offset for current word
func (w *WordCarouselAnimator) GetCurrentOffset() int {
	maxOffset := 3
	return int(float64(maxOffset) * (1.0 - w.CurrentPos))
}

// GetCurrentScale returns scale factor for current word
func (w *WordCarouselAnimator) GetCurrentScale() float64 {
	return 0.7 + (w.CurrentPos * 0.3)
}

// GetNextOffset returns vertical offset for next word
func (w *WordCarouselAnimator) GetNextOffset() int {
	maxOffset := 2
	return int(float64(maxOffset) * (1.0 - w.NextPos))
}

// GetNextOpacity returns opacity for next word (max 60%)
func (w *WordCarouselAnimator) GetNextOpacity() float64 {
	return w.NextPos * 0.6
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
