---
title: "Architecture"
description: "Understanding Baboon's technical design"
weight: 1
icon: "🏗️"
---

Baboon follows a clean backend/frontend separation with a well-defined REST API.

## High-Level Architecture

```
┌─────────────────┐     REST API      ┌─────────────────┐
│  Frontend (TUI) │ ◄───────────────► │     Backend     │
│   Bubble Tea    │                   │   Game Engine   │
└─────────────────┘                   └─────────────────┘
         ▲                                     │
         │                                     ▼
         │                            ┌─────────────────┐
┌─────────────────┐                   │   Statistics    │
│ Frontend (Web)  │                   │   ~/.config/    │
│     React       │ ◄─────────────────│     baboon/     │
└─────────────────┘                   └─────────────────┘
```

## Package Structure

```
baboon/
├── main.go           # Entry point
├── backend/
│   ├── api.go        # GameAPI interface & types
│   ├── engine.go     # Game logic implementation
│   └── server.go     # REST API server
├── frontend/
│   ├── model.go      # Bubble Tea model
│   ├── views.go      # Rendering
│   ├── styles.go     # Lipgloss styles
│   ├── animations.go # Spring animations
│   └── client.go     # REST client
├── web/              # React frontend
├── stats/            # Statistics & persistence
├── settings/         # User preferences
├── words/            # Word dictionary
└── font/             # Block letter font
```

## The GameAPI Interface

All frontend-backend communication goes through a single interface:

```go
type GameAPI interface {
    // Game Lifecycle
    StartRound()

    // Input Handling
    ProcessKeystroke(char string) KeystrokeResult
    ProcessKeystrokeWithTiming(char string, seekTimeMs int64) KeystrokeResult
    ProcessBackspace() bool
    ClearInput() bool
    ProcessSpace() SpaceResult
    ProcessSpaceWithTiming(seekTimeMs int64) SpaceResult
    SubmitTiming(startTime, endTime time.Time, durationMs int64)

    // State Queries
    GetGameState() GameState
    GetSessionStats() *stats.Stats
    GetHistoricalStats() *stats.HistoricalStats

    // Persistence
    SaveStats() error
}
```

## REST API Design

The backend exposes REST endpoints that both frontends consume:

### Session Management

- `POST /api/sessions` - Create session
- `DELETE /api/sessions/{id}` - Delete session
- `GET /api/sessions` - List sessions
- `GET /api/health` - Health check
- `GET /api/config` - Server configuration

### Game Operations

All scoped to `/api/sessions/{id}/`:

- `POST /round` - Start new round
- `POST /keystroke` - Process keystroke
- `POST /backspace` - Handle backspace
- `POST /clearinput` - Clear word
- `POST /space` - Handle space
- `POST /timing` - Submit timing
- `GET /state` - Get game state
- `GET /stats/session` - Session statistics
- `GET /stats/historical` - Historical statistics
- `POST /save` - Save statistics

## Frontend Timing

All timing-critical measurements are done on the frontend to avoid network latency:

1. Timer start/stop tracked locally
2. Seek times (between keystrokes) measured locally
3. Duration calculated from local timestamps
4. All timing data sent to backend at round end

## Statistics Persistence

Stats are saved to JSON files in `~/.config/baboon/`:

- `stats.json` - Historical statistics
- `settings.json` - User preferences

## Technologies Used

### Backend

- **Go 1.21+** - Core language
- **Standard library** - HTTP server, JSON

### Terminal Frontend

- **Bubble Tea** - TUI framework
- **Lipgloss** - Styling
- **Harmonica** - Spring animations

### Web Frontend

- **React 18** - UI framework
- **Chakra UI** - Component library
- **Framer Motion** - Animations
- **Vite** - Build tool
