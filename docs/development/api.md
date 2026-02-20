# API Reference

Complete reference for Baboon's REST API.

## Overview

- **Base URL**: `http://127.0.0.1:8787` (configurable via `-port`)
- **Format**: JSON
- **Authentication**: None (local use)

## Session Management

### Create Session

Creates a new game session.

```http
POST /api/sessions
Content-Type: application/json

{
  "punctuation_mode": false
}
```

**Response** (201 Created):

```json
{
  "session_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

### Delete Session

Removes a session and cleans up resources.

```http
DELETE /api/sessions/{session_id}
```

**Response** (204 No Content)

### List Sessions

Lists all active sessions.

```http
GET /api/sessions
```

**Response**:

```json
{
  "sessions": [
    {
      "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
      "created_at": "2024-01-15T10:30:00Z",
      "last_used": "2024-01-15T10:35:00Z"
    }
  ]
}
```

### Health Check

Check server status.

```http
GET /api/health
```

**Response**:

```json
{
  "status": "healthy",
  "active_sessions": 3
}
```

## Game Operations

All game operations are scoped to a session: `/api/sessions/{session_id}/...`

### Start Round

Starts a new 30-word round.

```http
POST /api/sessions/{session_id}/round
```

**Response**:

```json
{
  "words": ["hello", "world", "typing", ...],
  "total_words": 30,
  "total_characters": 150
}
```

### Process Keystroke

Submits a typed character with timing data.

```http
POST /api/sessions/{session_id}/keystroke
Content-Type: application/json

{
  "char": "a",
  "seek_time_ms": 150
}
```

**Response**:

```json
{
  "is_correct": true,
  "timer_started": true,
  "char_index": 0
}
```

| Field | Type | Description |
|-------|------|-------------|
| `is_correct` | boolean | Whether the character matches |
| `timer_started` | boolean | Whether this keystroke started the timer |
| `char_index` | int | Position in current word |

### Process Backspace

Removes the last typed character.

```http
POST /api/sessions/{session_id}/backspace
```

**Response**:

```json
{
  "success": true,
  "new_length": 3
}
```

### Process Space

Attempts to advance to the next word.

```http
POST /api/sessions/{session_id}/space
Content-Type: application/json

{
  "seek_time_ms": 200
}
```

**Response**:

```json
{
  "advanced": true,
  "round_complete": false,
  "treated_as_error": false
}
```

| Field | Type | Description |
|-------|------|-------------|
| `advanced` | boolean | Whether advanced to next word |
| `round_complete` | boolean | Whether this was the last word |
| `treated_as_error` | boolean | Whether space was counted as an error |

### Submit Timing

Submits final timing data when round completes.

```http
POST /api/sessions/{session_id}/timing
Content-Type: application/json

{
  "start_time_unix_ms": 1705312200000,
  "end_time_unix_ms": 1705312257000,
  "duration_ms": 57000
}
```

**Response**:

```json
{
  "success": true
}
```

### Get Game State

Retrieves current game state.

```http
GET /api/sessions/{session_id}/state
```

**Response**:

```json
{
  "words": ["hello", "world", "typing", ...],
  "current_word_idx": 5,
  "current_input": "typ",
  "timer_started": true,
  "punctuation_mode": false,
  "word_number": 6,
  "total_words": 30,
  "live_wpm": 52.3,
  "current_word": "typing",
  "previous_word": "world",
  "next_word": "practice",
  "next_words": ["practice", "test", "words"]
}
```

### Get Session Statistics

Retrieves statistics for the current session.

```http
GET /api/sessions/{session_id}/stats/session
```

**Response**:

```json
{
  "wpm": 52.3,
  "accuracy": 95.5,
  "time_seconds": 57.2,
  "correct_chars": 143,
  "total_chars": 150,
  "is_new_best_wpm": false,
  "is_new_best_accuracy": false,
  "is_new_best_time": true
}
```

### Get Historical Statistics

Retrieves cumulative statistics across all sessions.

```http
GET /api/sessions/{session_id}/stats/historical
```

**Response**:

```json
{
  "best_wpm": 65.5,
  "best_accuracy": 98.2,
  "best_time": 45.3,
  "average_wpm": 52.1,
  "average_accuracy": 94.8,
  "average_time": 58.7,
  "total_sessions": 15,
  "letter_accuracy": {
    "a": { "presented": 150, "correct": 145, "accuracy": 96.7 },
    "b": { "presented": 45, "correct": 42, "accuracy": 93.3 }
  },
  "letter_seek_time": {
    "a": { "total_time_ms": 22500, "count": 145, "average_ms": 155.2 },
    "b": { "total_time_ms": 9450, "count": 42, "average_ms": 225.0 }
  },
  "finger_stats": {
    "0": { "presented": 200, "correct": 195, "accuracy": 97.5 }
  },
  "row_stats": {
    "0": { "presented": 400, "correct": 380, "accuracy": 95.0 }
  },
  "hand_balance": { "left": 47.2, "right": 52.8 },
  "alternation_rate": 68.5,
  "sfb_stats": { "count": 150, "average_ms": 245.0 },
  "rhythm_stddev": 85.3,
  "error_patterns": [
    { "expected": "e", "typed": "r", "count": 5 },
    { "expected": "a", "typed": "s", "count": 3 }
  ]
}
```

### Save Statistics

Persists statistics to disk.

```http
POST /api/sessions/{session_id}/save
```

**Response**:

```json
{
  "success": true
}
```

## Error Responses

All errors follow this format:

```json
{
  "error": "Session not found",
  "code": "SESSION_NOT_FOUND"
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `SESSION_NOT_FOUND` | 404 | Session ID doesn't exist |
| `ROUND_NOT_STARTED` | 400 | No active round |
| `INVALID_REQUEST` | 400 | Malformed request body |
| `INTERNAL_ERROR` | 500 | Server error |

## Data Types

### KeystrokeResult

```typescript
interface KeystrokeResult {
  is_correct: boolean;
  timer_started: boolean;
  char_index: number;
}
```

### SpaceResult

```typescript
interface SpaceResult {
  advanced: boolean;
  round_complete: boolean;
  treated_as_error: boolean;
}
```

### GameState

```typescript
interface GameState {
  words: string[];
  current_word_idx: number;
  current_input: string;
  timer_started: boolean;
  punctuation_mode: boolean;
  word_number: number;
  total_words: number;
  live_wpm: number;
  current_word: string;
  previous_word: string;
  next_word: string;
  next_words: string[];
}
```

### LetterStats

```typescript
interface LetterStats {
  presented: number;
  correct: number;
  accuracy?: number;  // Calculated: correct/presented * 100
}
```

### SeekTimeStats

```typescript
interface SeekTimeStats {
  total_time_ms: number;
  count: number;
  average_ms?: number;  // Calculated: total_time_ms/count
}
```

## Usage Examples

### JavaScript

```javascript
// Create session
const response = await fetch('/api/sessions', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ punctuation_mode: false })
});
const { session_id } = await response.json();

// Start round
await fetch(`/api/sessions/${session_id}/round`, { method: 'POST' });

// Process keystroke
const result = await fetch(`/api/sessions/${session_id}/keystroke`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ char: 'a', seek_time_ms: 150 })
}).then(r => r.json());
```

### cURL

```bash
# Create session
SESSION=$(curl -s -X POST http://localhost:8787/api/sessions \
  -H "Content-Type: application/json" \
  -d '{"punctuation_mode":false}' | jq -r '.session_id')

# Start round
curl -X POST "http://localhost:8787/api/sessions/$SESSION/round"

# Get state
curl "http://localhost:8787/api/sessions/$SESSION/state"

# Process keystroke
curl -X POST "http://localhost:8787/api/sessions/$SESSION/keystroke" \
  -H "Content-Type: application/json" \
  -d '{"char":"h","seek_time_ms":0}'
```

### Go

```go
import (
    "bytes"
    "encoding/json"
    "net/http"
)

// Create session
reqBody, _ := json.Marshal(map[string]bool{"punctuation_mode": false})
resp, _ := http.Post("http://localhost:8787/api/sessions",
    "application/json", bytes.NewBuffer(reqBody))

var result map[string]string
json.NewDecoder(resp.Body).Decode(&result)
sessionID := result["session_id"]
```

## Rate Limits

No rate limits for local use. For networked deployments, consider adding rate limiting.

## Versioning

The API is currently v1. Breaking changes will increment the version.

## Next Steps

- [Architecture](architecture.md) - System design
- [Building](building.md) - Build instructions
- [Contributing](contributing.md) - How to contribute
