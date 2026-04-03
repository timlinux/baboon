// Package frontend provides the terminal user interface for the typing practice application.
package frontend

import (
	"math"
	"math/rand"
	"time"
)

// Celebration timing constants
const (
	CelebrationFireworksDuration = 8 * time.Second
	CelebrationMessageDuration   = 2 * time.Second
	CelebrationTickInterval      = 20 * time.Millisecond // 50 FPS
)

// CelebrationPhase represents the current phase of celebration
type CelebrationPhase int

const (
	PhaseFireworks CelebrationPhase = iota
	PhaseMessage
)

// BestType represents which types of personal bests were achieved
type BestType int

const (
	BestWPM BestType = 1 << iota
	BestAccuracy
	BestTime
)

// Particle represents a single particle in the celebration
type Particle struct {
	X, Y       float64 // Position
	VX, VY     float64 // Velocity
	Life       float64 // Remaining life (0-1)
	LifeDecay  float64 // How fast life decreases per second
	Color      int     // Terminal color code (256-color)
	Char       rune    // Character to display
	IsDebris   bool    // Whether this is text debris
	Brightness float64 // For fading effects (0-1)
}

// TextCell represents a character cell on screen that can be destroyed
type TextCell struct {
	Char      rune
	X, Y      int
	Destroyed bool
	Color     int // Original color
}

// ExplosionConfig controls firework explosion appearance
type ExplosionConfig struct {
	ParticleCount int
	Speed         float64
	Color         int
	Char          rune
}

// CelebrationState manages the celebration animation
type CelebrationState struct {
	Phase         CelebrationPhase
	StartTime     time.Time
	LastUpdate    time.Time
	Particles     []Particle
	TextGrid      [][]TextCell // Screen grid for collision detection
	BestTypes     BestType
	WPMValue      float64
	AccuracyValue float64
	TimeValue     float64 // Time in seconds
	Width, Height int

	// Scheduled explosion times
	explosionTimes []time.Duration
	nextExplosion  int

	// Physics constants
	gravity     float64
	groundLevel float64
}

// Celebratory colors for fireworks
var fireworkColors = []int{
	226, // Gold
	196, // Red
	208, // Orange
	201, // Magenta
	51,  // Cyan
	46,  // Green
	21,  // Blue
	129, // Purple
}

// NewCelebrationState creates a new celebration state
func NewCelebrationState(width, height int, bestTypes BestType, wpm, accuracy, timeSeconds float64) *CelebrationState {
	now := time.Now()
	c := &CelebrationState{
		Phase:         PhaseFireworks,
		StartTime:     now,
		LastUpdate:    now,
		Particles:     make([]Particle, 0, 500),
		BestTypes:     bestTypes,
		WPMValue:      wpm,
		AccuracyValue: accuracy,
		TimeValue:     timeSeconds,
		Width:         width,
		Height:        height,
		gravity:       30.0, // Gravity strength
		groundLevel:   float64(height) - 2,
	}

	// Schedule explosions throughout the fireworks phase
	// Create approximately 12 explosions over 8 seconds
	numExplosions := 12
	interval := CelebrationFireworksDuration / time.Duration(numExplosions+1)
	c.explosionTimes = make([]time.Duration, numExplosions)
	for i := range numExplosions {
		// Add some randomness to timing (±200ms)
		jitter := time.Duration(rand.Intn(400)-200) * time.Millisecond
		c.explosionTimes[i] = interval*time.Duration(i+1) + jitter
	}
	c.nextExplosion = 0

	// Trigger initial explosion
	c.TriggerExplosion(
		float64(width/2)+float64(rand.Intn(width/3)-width/6),
		float64(height/3)+float64(rand.Intn(height/4)),
	)

	return c
}

// Update advances the particle simulation
func (c *CelebrationState) Update(dt float64) {
	c.LastUpdate = time.Now()
	elapsed := time.Since(c.StartTime)

	// Check for phase transition
	if c.Phase == PhaseFireworks && elapsed >= CelebrationFireworksDuration {
		c.Phase = PhaseMessage
		c.Particles = nil // Clear particles for message phase
		return
	}

	// Check if we should trigger a scheduled explosion
	if c.Phase == PhaseFireworks && c.nextExplosion < len(c.explosionTimes) {
		if elapsed >= c.explosionTimes[c.nextExplosion] {
			// Trigger explosion at random position
			x := float64(c.Width/4) + float64(rand.Intn(c.Width/2))
			y := float64(c.Height/4) + float64(rand.Intn(c.Height/3))
			c.TriggerExplosion(x, y)
			c.nextExplosion++
		}
	}

	// Update particles
	newParticles := make([]Particle, 0, len(c.Particles))
	for i := range c.Particles {
		p := &c.Particles[i]

		// Apply gravity
		p.VY += c.gravity * dt

		// Update position
		p.X += p.VX * dt
		p.Y += p.VY * dt

		// Decay life
		p.Life -= p.LifeDecay * dt
		p.Brightness = p.Life

		// Bounce off ground
		if p.Y >= c.groundLevel {
			p.Y = c.groundLevel
			p.VY = -p.VY * 0.3 // Restitution
			p.VX *= 0.8        // Friction
		}

		// Bounce off walls
		if p.X < 0 {
			p.X = 0
			p.VX = -p.VX * 0.5
		} else if p.X >= float64(c.Width) {
			p.X = float64(c.Width) - 1
			p.VX = -p.VX * 0.5
		}

		// Check collision with text cells
		c.checkTextCollision(p)

		// Keep particle if still alive
		if p.Life > 0 {
			newParticles = append(newParticles, *p)
		}
	}
	c.Particles = newParticles
}

// checkTextCollision checks if a particle collides with text and destroys it
func (c *CelebrationState) checkTextCollision(p *Particle) {
	if c.TextGrid == nil {
		return
	}

	cellX := int(p.X)
	cellY := int(p.Y)

	if cellY < 0 || cellY >= len(c.TextGrid) {
		return
	}
	row := c.TextGrid[cellY]
	if cellX < 0 || cellX >= len(row) {
		return
	}

	cell := &row[cellX]
	if !cell.Destroyed && cell.Char != ' ' && cell.Char != 0 {
		// Destroy the cell
		cell.Destroyed = true

		// Reduce particle life
		p.Life -= 0.3

		// Spawn debris particles using the destroyed character
		c.spawnDebris(float64(cellX), float64(cellY), cell.Char, cell.Color)
	}
}

// spawnDebris creates debris particles from destroyed text
func (c *CelebrationState) spawnDebris(x, y float64, char rune, color int) {
	// Spawn 3-5 debris particles
	count := 3 + rand.Intn(3)
	for i := 0; i < count; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := 5.0 + rand.Float64()*10.0
		c.Particles = append(c.Particles, Particle{
			X:          x,
			Y:          y,
			VX:         math.Cos(angle) * speed,
			VY:         math.Sin(angle)*speed - 5, // Bias upward
			Life:       0.5 + rand.Float64()*0.5,
			LifeDecay:  1.5,
			Color:      color,
			Char:       char,
			IsDebris:   true,
			Brightness: 1.0,
		})
	}
}

// TriggerExplosion spawns a firework explosion at the given position
func (c *CelebrationState) TriggerExplosion(x, y float64) {
	// Choose random explosion color
	color := fireworkColors[rand.Intn(len(fireworkColors))]

	// Create radial explosion
	particleCount := 30 + rand.Intn(20)
	baseSpeed := 15.0 + rand.Float64()*10.0

	for i := range particleCount {
		angle := float64(i) * (2 * math.Pi / float64(particleCount))
		// Add some randomness to angle and speed
		angle += (rand.Float64() - 0.5) * 0.3
		speed := baseSpeed * (0.7 + rand.Float64()*0.6)

		// Choose particle character
		chars := []rune{'*', '.', '+', '·', '•', '✦', '✧'}
		char := chars[rand.Intn(len(chars))]

		c.Particles = append(c.Particles, Particle{
			X:          x,
			Y:          y,
			VX:         math.Cos(angle) * speed,
			VY:         math.Sin(angle)*speed - 8, // Initial upward boost
			Life:       0.8 + rand.Float64()*0.5,
			LifeDecay:  0.6 + rand.Float64()*0.4,
			Color:      color,
			Char:       char,
			IsDebris:   false,
			Brightness: 1.0,
		})
	}

	// Add some trail/sparkle particles
	sparkleCount := 15 + rand.Intn(10)
	for i := range sparkleCount {
		_ = i // Silence unused variable warning
		angle := rand.Float64() * 2 * math.Pi
		speed := baseSpeed * 0.3 * (0.5 + rand.Float64())

		c.Particles = append(c.Particles, Particle{
			X:          x + (rand.Float64()-0.5)*2,
			Y:          y + (rand.Float64()-0.5)*2,
			VX:         math.Cos(angle) * speed,
			VY:         math.Sin(angle)*speed - 3,
			Life:       0.3 + rand.Float64()*0.3,
			LifeDecay:  1.5,
			Color:      226, // Gold sparkles
			Char:       '·',
			IsDebris:   false,
			Brightness: 1.0,
		})
	}
}

// InitTextGrid initializes the text grid from screen content for collision detection
func (c *CelebrationState) InitTextGrid(lines []string) {
	c.TextGrid = make([][]TextCell, c.Height)
	for y := range c.Height {
		c.TextGrid[y] = make([]TextCell, c.Width)
		var line string
		if y < len(lines) {
			line = lines[y]
		}
		runes := []rune(line)
		for x := range c.Width {
			var char rune
			if x < len(runes) {
				char = runes[x]
			} else {
				char = ' '
			}
			c.TextGrid[y][x] = TextCell{
				Char:      char,
				X:         x,
				Y:         y,
				Destroyed: false,
				Color:     252, // Default white
			}
		}
	}
}

// IsComplete returns true when the celebration has finished
func (c *CelebrationState) IsComplete() bool {
	elapsed := time.Since(c.StartTime)
	totalDuration := CelebrationFireworksDuration + CelebrationMessageDuration
	return elapsed >= totalDuration
}
