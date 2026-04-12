# Baboon - Typing Practice Application

## Overview

Baboon is a cross-platform typing practice application built with Go. It helps users improve their typing speed and accuracy by presenting common English words in large block letter format. The application features two frontends:

1. **Terminal UI (TUI)**: Built with Bubble Tea and Lipgloss, displaying words using Unicode block characters (█) that change colour in real-time as the user types. The TUI runs as a **single embedded binary** with no client-server architecture required.

2. **Web UI**: Built with React and Chakra UI, featuring physics-based animations with Framer Motion, large chunky block letters, and a beautiful dark theme. The web frontend communicates with the Go backend via REST API.

The TUI uses an embedded game engine directly, while the web frontend uses the REST API for multi-client support.

## User Stories

### US-001: Practice Typing Common Words
**As a** user wanting to improve my typing skills
**I want to** practice typing common English words
**So that** I can increase my typing speed and accuracy

### US-002: Visual Feedback During Typing
**As a** user practicing typing
**I want to** see immediate visual feedback on my keystrokes
**So that** I can identify and correct errors quickly

### US-003: Track My Progress
**As a** regular user
**I want to** see my typing statistics after each round
**So that** I can measure my improvement over time

### US-004: Compare to Personal Best
**As a** competitive user
**I want to** compare my current performance to my historical best and average
**So that** I can see how I'm improving

### US-005: Monitor Speed in Real-Time
**As a** user during a typing session
**I want to** see my current WPM as I type
**So that** I can adjust my pace accordingly

### US-006: Access via Web Browser
**As a** user who prefers a web-based interface
**I want to** practice typing through a beautiful web application
**So that** I can use any device with a modern browser

### US-007: Cross-Device Stats Synchronisation
**As a** registered user
**I want to** have my typing statistics synchronised across devices
**So that** I can practice on multiple devices and see my overall progress

### US-008: Sign In with SSO
**As a** user who doesn't want to create yet another account
**I want to** sign in with my existing Google, GitHub, Apple, or Microsoft account
**So that** I can quickly start using authenticated features without a new password

### US-009: Compete on Monthly Leaderboard
**As a** competitive typist
**I want to** see my ranking on a monthly leaderboard
**So that** I can compete with other users and track my standing

### US-010: Enter Arcade Name for Leaderboard
**As a** user who achieves a top 10 score
**I want to** enter a classic arcade-style name (up to 10 characters)
**So that** my achievement is displayed prominently on the leaderboard

### US-011: Share My Leaderboard Achievement
**As a** proud user who made the leaderboard
**I want to** share a badge image showing my score and rank
**So that** I can share my achievement on social media

## Functional Requirements

### FR-001: Word Display
- The application SHALL display words from a dictionary of common English words (British English spelling)
- Each word SHALL be rendered in large block characters using Unicode block elements (█)
- Words SHALL be displayed centered horizontally and vertically on the terminal screen
- The word display SHALL show progress indicator: "Word X/30"
- Letters SHALL change colour in-place as the user types (no separate input display line)
- All words SHALL be displayed in UPPERCASE in the block font for maximum readability
- User input is case-insensitive; lowercase input matches uppercase display
- Words SHALL be displayed in a carousel layout:
  - The previous word SHALL be displayed ABOVE the current word in dimmed text
  - The next 3 upcoming words SHALL be displayed BELOW the current word in dimmed text
  - When advancing to the next word, smooth carousel animation SHALL scroll words upward
- Console: Previous word uses greyscale colour (240), next words use decreasing greyscale (from 245), with decorative markers on first upcoming word
- Web: Previous/next words displayed at 50% scale with blur effects

### FR-002: Block Letter Font
- Each letter SHALL be 6 lines tall
- Letters SHALL be constructed using Unicode block elements for smooth edges:
  - █ (full block) for solid letter bodies
  - ◢ ◣ ◤ ◥ (filled triangles) for smooth rounded corners at curved edges
- The font SHALL support letters A-Z (uppercase), numbers 0-9, and punctuation: , . ; : ! ? - ' "
- Unknown characters SHALL render as spaces
- Letters SHALL have 1 character spacing between them

### FR-003: Typing Input and Colour Feedback
- The application SHALL accept keyboard input character by character
- Each character typed SHALL immediately change the corresponding letter's colour:
  - **Green (colour 10)**: Character matches the expected character at that position
  - **Red (colour 9)**: Character does not match the expected character
  - **Gray (colour 8)**: Characters not yet typed
- The backspace key SHALL remove the last typed character (undoing its colour)
- Ctrl+W (or Ctrl+Backspace) SHALL clear all typed characters for the current word
- The space key SHALL only advance to the next word when ALL letters have been typed
- If space is pressed before the word is complete, it SHALL be treated as an incorrect character (red)
- Extra characters beyond word length SHALL count as incorrect (red)
- The last word (word 30) SHALL auto-complete when the final character is typed correctly
  - No space press is required for the final word
  - The round immediately transitions to results upon correct final character
  - The word counter SHALL remain at "30/30" (never exceed total words)

### FR-004: Round Structure
- Each round SHALL consist of exactly 30 words totalling exactly 150 characters
- Words SHALL be randomly selected using stratified selection to meet both constraints
- Word selection algorithm SHALL:
  - Calculate ideal word length based on remaining characters and words
  - Allow variance of ±2 characters from ideal to maintain variety
  - Ensure feasibility by checking remaining capacity
  - Retry up to 100 times if constraints cannot be met
- After completing all 30 words, the application SHALL display the results screen
- The user SHALL be able to start a new round by pressing Enter on results screen

### FR-015: Adaptive Word Selection
- Word selection SHALL be weighted based on two factors:
  1. **Frequency balancing**: Favour words with underrepresented letters
  2. **Accuracy practice**: Favour words with letters the user frequently mistypes
- Each word SHALL be scored using a combined algorithm:
  - Frequency score = 1 - (letter_presented / max_letter_presented)
  - Accuracy score = 1 - (letter_correct / letter_presented)
  - Letter score = (frequency_score + accuracy_score) / 2
  - Word score = average letter score across all letters in the word
- Words with higher scores SHALL have higher selection probability
- This adaptive selection helps users practice their weakest letters
- Frequency balancing aims to achieve spread within 10% from highest to lowest

### FR-005: Timer Behaviour
- The timer SHALL NOT start when the application launches
- The timer SHALL start when the user types the first CORRECT character of the first word
- If the first character typed is incorrect, the timer SHALL NOT start
- The timer SHALL stop when the 30th word is completed (space pressed)

### FR-006: Live WPM Bar (During Typing)
- The application SHALL display a gradient WPM bar at the bottom of the screen during typing
- The bar SHALL be 50 characters wide
- The bar SHALL update every 100ms to show current typing speed
- The bar SHALL use a gradient colour scheme from red (slow) through yellow to green (fast):
  - Colours: 196, 202, 208, 214, 220, 226, 190, 154, 118, 82, 46, 47
- The bar SHALL scale from 0 to 120 WPM maximum
- The bar SHALL display numeric WPM value with colour coding:
  - Red (colour 196): Below 40 WPM
  - Yellow (colour 226): 40-60 WPM
  - Green (colour 46): Above 60 WPM
- The bar SHALL show scale markers: "0", "60", "120"
- Empty portion of bar SHALL use character ░ in colour 236

### FR-007: Results Screen Layout
- The results screen SHALL display "Round Complete!" title in cyan (colour 14), bold
- Statistics SHALL be displayed in a grid layout with three columns:
  - Label column: 18 characters wide, right-aligned, gray (colour 7)
  - Value column: 8 characters wide, right-aligned, white (colour 15)
  - Bar column: 30 characters wide gradient bar + 2 character star column

### FR-008: Results Statistics Display
- WPM section:
  - "WPM this run:" with current session WPM and bar
  - "WPM best:" with historical best WPM and bar
  - "WPM average:" with calculated average WPM and bar
- Time section (blank line before):
  - "Time this run:" with session time in seconds (e.g., "147.2s") and bar
  - "Time best:" with historical best (fastest) time and bar
  - "Time average:" with calculated average time and bar
- Accuracy section (blank line before):
  - "Accuracy this run:" with percentage (e.g., "95.5%") and bar
  - "Accuracy best:" with historical best accuracy and bar
  - "Accuracy average:" with calculated average accuracy and bar
- Sessions section (blank line before):
  - "Total sessions:" label in cyan (colour 6) with count
- Legend (blank line before):
  - "* = New personal best!" in yellow (colour 226), bold

### FR-009: Results Bar Rendering
- WPM bars: Scale 0-120, higher is better (more fill = better)
- Time bars: Scale 0-180 seconds, INVERTED (lower time = more fill = better)
- Accuracy bars: Scale 0-100%, higher is better (more fill = better)
- New personal best SHALL show " *" after the bar
- Non-best bars SHALL show "  " (two spaces) to maintain alignment
- All bars SHALL use same gradient colours as live WPM bar

### FR-010: Historical Statistics Persistence
- Historical data SHALL be stored in `~/.config/baboon/stats.json`
- The config directory SHALL be created automatically if it doesn't exist
- The application SHALL track:
  - `best_wpm`: Highest WPM achieved (float64)
  - `best_accuracy`: Highest accuracy percentage achieved (float64)
  - `best_time`: Fastest (lowest) completion time in seconds (float64)
  - `total_wpm`: Sum of all session WPMs for averaging (float64)
  - `total_accuracy`: Sum of all session accuracies for averaging (float64)
  - `total_time`: Sum of all session times for averaging (float64)
  - `total_sessions`: Count of completed sessions (int)
  - `last_session_date`: Timestamp of last session (RFC3339)
  - `letter_accuracy`: Per-letter accuracy tracking (map of letter to stats)
  - `letter_seek_time`: Per-letter seek time tracking (map of letter to timing stats)
  - `bigram_seek_time`: Per-bigram (letter pair) seek time tracking (map of bigram to timing stats)

### FR-013: Per-Letter Accuracy Tracking
- When a round starts, all letters in all 30 words SHALL be recorded as "presented"
- When a user types a correct letter, that letter SHALL be recorded as "correct"
- Letter statistics SHALL be tracked per individual letter (a-z)
- For each letter, the application SHALL track:
  - `presented`: Number of times this letter was presented to the user
  - `correct`: Number of times the user typed this letter correctly
- Letter accuracy data SHALL persist across sessions (cumulative)
- Letter accuracy SHALL be calculated as: (correct / presented) × 100

### FR-016: Per-Letter Seek Time Tracking
- The application SHALL track the time between keystrokes (seek time)
- Seek time SHALL only be recorded for CORRECT keystrokes
- Seek time SHALL be recorded against the EXPECTED letter (not the typed character)
- The FIRST letter of each word SHALL be excluded from seek time tracking (includes word-reading time)
- Seek times > 5000ms SHALL be filtered out (assumed user pauses)
- For each letter, the application SHALL track:
  - `total_time_ms`: Total seek time in milliseconds
  - `count`: Number of measurements
- Average seek time = total_time_ms / count
- Seek time data SHALL persist across sessions (cumulative)

### FR-017: Bigram (Letter Pair) Seek Time Tracking
- The application SHALL track seek time for letter pairs (bigrams)
- A bigram is formed from the previous correctly typed letter + current correctly typed letter
- Bigrams SHALL only be recorded for consecutive correct keystrokes
- Bigrams SHALL reset at word boundaries (first letter of new word has no preceding letter)
- For each bigram (e.g., "th", "he", "in"), the application SHALL track:
  - `total_time_ms`: Total seek time in milliseconds
  - `count`: Number of measurements
- Bigram data SHALL persist across sessions (cumulative)
- Common slow bigrams indicate letter combinations the user struggles with

### FR-014: Letter Statistics Display
- The results screen SHALL display a letter statistics matrix:
  1. **Header row**: 26 uppercase letters (A-Z) as column labels, white bold text
  2. **Accuracy row**: Filled circles (●) coloured by typing accuracy
  3. **Frequency row**: Filled circles (●) coloured by presentation count
  4. **Seek time row**: Filled circles (●) coloured by average typing speed
- Each circle SHALL be coloured using **relative scaling** (see BR-006)
- Letters with no data SHALL be displayed in gray (colour 240)
- Letters in header row are spaced to align with circles below
- Seek time is measured as milliseconds between keystrokes
- Seek times > 5 seconds are filtered out (assumed pauses)

### FR-027: Relative Colour Scaling for Statistics
- Finger, row, and letter statistics SHALL use **relative colour scaling**
- Relative scaling ensures meaningful visual differentiation when values are clustered
- The colour gradient spreads across the actual data range (min to max)
- Statistics with best values get green, worst values get red
- This prevents all items appearing identical when they have similar absolute values
- Example: If finger accuracies range from 92% to 97%, the 92% finger gets red and 97% gets green
- Absolute thresholds (FR-009) still apply to session-level WPM, Time, Accuracy bars

**Gradient colours (index 0-11 in GradientColours array):**
  - 196: Red (worst in range)
  - 202, 208, 214, 220, 226: Red → Yellow gradient
  - 190, 154, 118, 82: Yellow → Green gradient
  - 46, 47: Bright green (best in range)

### FR-018: Results Screen Animation
- Results screen elements SHALL animate in sequentially using spring physics
- The harmonica library SHALL be used for smooth spring-based animations
- Each stat row SHALL slide in from the right with staggered timing
- Animation interval SHALL be 50ms per frame
- Stagger delay SHALL be 3 frames between each row starting
- Spring parameters: 60 FPS, frequency 6.0, damping 0.5
- Total of 25 animated rows (core stats + typing theory stats + letter matrix)

### FR-020: Finger-Specific Statistics
- The application SHALL track per-finger typing accuracy and speed
- Standard touch typing finger assignments SHALL be used:
  - Left pinky (0): q, a, z
  - Left ring (1): w, s, x
  - Left middle (2): e, d, c
  - Left index (3): r, f, v, t, g, b
  - Right index (6): y, h, n, u, j, m
  - Right middle (7): i, k
  - Right ring (8): o, l
  - Right pinky (9): p
- For each finger, the application SHALL track:
  - `presented`: Times a key for this finger was needed
  - `correct`: Times the correct key was pressed
  - `total_time_ms`: Total seek time for keys typed with this finger
  - `count`: Number of timed keypresses
- Results screen SHALL display finger accuracy row with colour-coded indicators
- Finger labels: LP, LR, LM, LI (left hand), RI, RM, RR, RP (right hand)

### FR-021: Keyboard Row Statistics
- The application SHALL track per-row typing accuracy and speed
- Row assignments:
  - Top row (0): q, w, e, r, t, y, u, i, o, p
  - Home row (1): a, s, d, f, g, h, j, k, l
  - Bottom row (2): z, x, c, v, b, n, m
- For each row, the application SHALL track:
  - `presented`: Times a key on this row was needed
  - `correct`: Times the correct key was pressed
  - `total_time_ms`: Total seek time for keys on this row
  - `count`: Number of timed keypresses
- Results screen SHALL display row accuracy with labels: Top, Home, Bot

### FR-022: Hand Balance and Alternation Tracking
- The application SHALL track hand usage balance (left vs right)
- Hand assignments: Left (q-t, a-g, z-b), Right (y-p, h-l, n-m)
- The application SHALL track hand alternations vs same-hand runs:
  - `hand_alternations`: Count of transitions between hands
  - `same_hand_runs`: Count of consecutive same-hand keypresses
- Alternation rate = hand_alternations / (hand_alternations + same_hand_runs) × 100
- Higher alternation rate indicates better typing flow
- Results screen SHALL display hand balance (L:X% R:Y%) and alternation rate

### FR-023: Same-Finger Bigram (SFB) Tracking
- The application SHALL detect and track same-finger bigrams
- An SFB occurs when consecutive letters use the same finger
- For each SFB occurrence, the application SHALL track:
  - Count of SFBs encountered
  - Total seek time for SFBs
  - Average seek time = total_time / count
- SFBs are inherently slower than alternating-finger bigrams
- Results screen SHALL display SFB count and average time per session

### FR-024: Rhythm Consistency (Variance) Tracking
- The application SHALL track typing rhythm consistency
- Rhythm is measured as the standard deviation of seek times
- For rhythm calculation, the application SHALL track:
  - All seek times during the session
  - Sum of seek times
  - Sum of squared seek times (for variance calculation)
- Variance = (sum_of_squares / count) - (mean²)
- Standard deviation = √variance
- Lower standard deviation indicates more consistent rhythm
- Results screen SHALL display session StdDev and historical average

### FR-025: Error Substitution Pattern Tracking
- The application SHALL track which letters are commonly confused
- When a letter is mistyped, the application SHALL record:
  - The expected letter
  - The typed letter
  - Increment the count for this (expected → typed) pair
- Error substitution data SHALL persist across sessions
- Results screen SHALL display top 5 most common error patterns
- Format: "a→s(12)" means 'a' was typed as 's' 12 times

### FR-019: Punctuation Mode
- The application SHALL support a `-p` command line flag for punctuation mode
- When enabled, words SHALL be separated by random punctuation followed by space
- Supported punctuation characters: , . ; : ! ?
- Punctuation SHALL be appended to each word except the last word in the round
- The user SHALL type the punctuation character before pressing space to advance
- Letter accuracy tracking SHALL only count letters (a-z), not punctuation
- Letter seek time tracking SHALL only measure letters (a-z), not punctuation
- Punctuation mode persists for subsequent rounds until the application exits

### FR-011: Statistics Validation
- On load, the application SHALL validate historical statistics for corruption
- If totals are 0 but bests exist, data SHALL be reset using best values as estimates
- If average WPM is less than half of best WPM, data SHALL be reset
- Reset formula: total = best × total_sessions

### FR-012: Navigation
- ESC or Ctrl+C SHALL exit the application at any time
- Ctrl+W SHALL clear all typed characters for the current word (undo all input)
- The advance key (configurable: Space, Enter, or Either) SHALL advance to the next word during typing (when input length > 0)
- ENTER SHALL start a new round when viewing results screen
- TAB SHALL restart the current round at any time (typing or results screen)
- Ctrl+O SHALL open the options screen before the timer starts (from typing screen) or at any time (from results screen)
- The application SHALL use alternate screen buffer (fullscreen mode)

### FR-031: Accessibility (Web UI)
- All interactive elements SHALL have visible focus indicators
- All buttons SHALL have ARIA labels describing their function and keyboard shortcuts
- The typing area SHALL include ARIA live regions to announce progress to screen readers
- Letter status (correct/incorrect) SHALL include visual indicators beyond color:
  - Correct letters SHALL display a checkmark (✓) indicator
  - Incorrect letters SHALL display an X (×) indicator
  - Current letter SHALL display an arrow (▸) indicator
- Progress and word counter SHALL have descriptive ARIA labels
- Screen reader announcements SHALL include current word and progress
- Focus styles SHALL use brand colors with 2px outline offset

### FR-028: Web Frontend Local Storage
- The web frontend SHALL store user historical statistics in browser local storage
- The storage key SHALL be `baboon_historical_stats` for statistics
- The storage key SHALL be `baboon_settings` for user preferences (e.g., punctuation mode)
- The welcome screen SHALL display the user's best WPM, best accuracy, and total sessions if available
- Settings SHALL persist across browser sessions
- Local storage data SHALL be updated after each completed round

### FR-032: User Authentication System
- The application SHALL support optional user authentication via OAuth2/OIDC
- Authentication SHALL be enabled when a database DSN is configured (`BABOON_DATABASE_DSN`)
- Supported OAuth providers: Google, GitHub, Apple, Microsoft
- Only configured providers (with client ID and secret) SHALL be shown in the login UI
- If no database or OAuth providers are configured, authentication features SHALL be hidden
- Anonymous users SHALL continue using localStorage for stats (no degradation of existing functionality)

### FR-033: OAuth Login Flow
- Each OAuth provider SHALL have a dedicated login endpoint: `GET /api/auth/{provider}/login`
- The login endpoint SHALL redirect the user to the OAuth provider's authorization page
- The OAuth callback endpoint `GET /api/auth/{provider}/callback` SHALL:
  - Validate the OAuth state parameter to prevent CSRF
  - Exchange the authorization code for access token
  - Retrieve user info from the provider
  - Create or update the user record in the database
  - Generate JWT access token and refresh token
  - Set HTTP-only cookies with the tokens
  - Redirect to the frontend with success or error status

### FR-034: JWT Token Management
- Access tokens SHALL be JWTs signed with HS256 using a configurable secret
- Access tokens SHALL expire after 15 minutes by default
- Refresh tokens SHALL be stored in the database with hashed values
- Refresh tokens SHALL expire after 7 days by default
- The refresh endpoint `POST /api/auth/refresh` SHALL:
  - Validate the current refresh token
  - Revoke the old refresh token
  - Generate new access and refresh tokens
- Access token cookies SHALL be accessible by JavaScript for API calls
- Refresh token cookies SHALL be HTTP-only for security

### FR-035: User Stats Synchronisation
- Authenticated users SHALL have their stats stored in the database
- The sync endpoint `POST /api/user/stats/sync` SHALL:
  - Accept local stats from the client
  - Merge with existing server stats using the merge algorithm
  - Save the merged stats to the database
  - Return the merged stats
- Merge algorithm:
  - Cumulative values (total sessions, total time): ADD both together
  - Best values (best WPM, best accuracy, best time): Take MAX (or MIN for time)
  - Map values (letter stats, finger stats): Merge additively
- When a user logs in with existing local stats, the frontend SHALL offer to sync

### FR-036: User Data Privacy
- Users SHALL be able to export all their data via `GET /api/user/export`
- Export SHALL include user profile and all statistics
- Users SHALL be able to delete all their stats via `DELETE /api/user/stats`
- Users SHALL be able to delete their account via `DELETE /api/auth/account`
- Account deletion SHALL cascade delete all user data (stats, tokens)
- No keystroke logging or raw typing data SHALL be stored - only aggregate statistics

### FR-029: Google AdSense Integration
- The web server mode SHALL accept an `-adsense` flag with a Google AdSense publisher ID
- When AdSense is enabled, the `/api/config` endpoint SHALL return the publisher ID
- The web frontend SHALL display an ad component beneath the typing game when a publisher ID is configured
- The ad component SHALL load the Google AdSense script dynamically
- The ad SHALL be styled to match the application's dark theme
- Ads SHALL NOT interfere with the typing experience or keyboard input

### FR-030: Personal Best Celebration
- When the user achieves a new personal best WPM or accuracy, the application SHALL display a celebration animation
- The celebration SHALL consist of two phases:
  1. **Fireworks Phase** (8 seconds):
     - Particle-based fireworks explosions across the screen
     - Approximately 12 explosions scheduled throughout the phase
     - Particles use celebratory colours: gold, red, orange, magenta, cyan, green, blue, purple
     - Each explosion spawns 30-50 radial particles plus sparkle particles
     - Particles are affected by gravity, bounce off screen boundaries
     - Text destruction effect: particles colliding with text cells destroy them, spawning debris
  2. **Message Phase** (2 seconds):
     - Large block font message displays "PERSONAL BEST" centered on screen
     - The achieved WPM value is displayed below in green
     - Decorative sparkles and "Press any key to continue" hint
- Any keypress during celebration SHALL skip directly to the results screen
- Ctrl+C during celebration SHALL exit the application
- Celebration uses 50 FPS (20ms tick interval) for smooth animation
- After celebration completes, the normal results screen is displayed

### FR-026: Options Screen
- The application SHALL provide an options screen accessible via Ctrl+O
- The options screen SHALL be accessible from typing screen (before timer starts) and results screen
- The options screen SHALL allow configuring which key advances to the next word:
  - **Space** (default): Press Space to advance to the next word
  - **Enter**: Press Enter to advance to the next word
  - **Either**: Press Space or Enter to advance to the next word
- Navigation in options screen:
  - ↑/↓ or Tab/Shift+Tab SHALL move the cursor between options
  - Enter or Space SHALL select the highlighted option
  - Number keys (1-3) SHALL quick-select the corresponding option
  - ESC SHALL return to the previous screen without changes
- The currently selected option SHALL be indicated with a green checkmark (✓)
- The cursor position SHALL be highlighted with a contrasting background colour
- Settings SHALL be persisted to `~/.config/baboon/settings.json`
- Settings SHALL be loaded on application startup
- Help text in typing screen SHALL dynamically reflect the configured advance key

### FR-037: Monthly Leaderboard
- The application SHALL maintain a leaderboard showing the top 10 scores per month
- Leaderboard entries SHALL be partitioned by month in "YYYY-MM" format
- Each user SHALL have at most one entry per month (best score replaces existing entry)
- The leaderboard screen SHALL display:
  - Rank (1-10) with gold/silver/bronze styling for top 3
  - Display name (up to 10 characters)
  - WPM with gradient colour bar
  - Accuracy percentage
  - Share button
- Users SHALL be able to browse previous months via dropdown
- Current month SHALL be shown by default
- The leaderboard SHALL be accessible from welcome screen and results screen
- ESC key SHALL return to previous screen

### FR-038: Arcade Name Entry
- When an authenticated user achieves a top 10 score, they SHALL be prompted to enter a name
- The name entry screen SHALL appear after round completion but before results
- Names SHALL be limited to 10 characters
- Allowed characters: A-Z, 0-9, space, underscore, hyphen
- Character slots SHALL display arcade-style visual feedback:
  - Empty slots: grey with underscore
  - Typed valid chars: green background and text
  - Invalid chars (profanity detected): red background and text
  - Current slot: Kartoza orange border with glow
- Real-time profanity validation SHALL highlight problematic characters
- ENTER key SHALL submit the name
- ESC key SHALL skip leaderboard entry and proceed to results
- Profanity filter SHALL detect:
  - Direct matches from blocklist
  - Leetspeak variants (e.g., "4" for "a", "3" for "e")

### FR-039: Leaderboard Badge Sharing
- Each leaderboard entry SHALL have a shareable SVG badge
- Badge SHALL display:
  - Simplified baboon mascot with keyboard
  - Player display name
  - WPM score
  - Accuracy percentage
  - Rank with star indicator (gold for #1, orange for #2-3)
  - Month/year
  - Site URL (baboon.kartoza.com)
  - Kartoza branding footer
- Badge dimensions: 400x200 pixels
- Share modal SHALL provide:
  - Badge preview
  - Copy Link button
  - Download SVG button
  - Share on X (Twitter) button
- Share URL format: `https://baboon.kartoza.com/leaderboard?highlight={id}`

## Technical Requirements

### TR-001: Cross-Platform Compatibility
- The application SHALL run on Linux, macOS, and Windows
- The application SHALL be buildable using Nix flakes for reproducible builds
- The terminal SHALL support 256-colour mode for proper gradient display

### TR-002: Terminal Interface
- The application SHALL use the Bubble Tea framework (github.com/charmbracelet/bubbletea)
- The application SHALL use lipgloss for styling (github.com/charmbracelet/lipgloss)
- The application SHALL use github.com/timlinux/blockfont for block letter rendering
- The application SHALL use tea.WithAltScreen() for fullscreen mode
- The application SHALL handle tea.WindowSizeMsg for responsive centering

### TR-004: Screen Layout Structure
- All screens (typing, results, options) SHALL use a fixed three-section layout:
  1. **Header**: Fixed at the top of the terminal (line 0), displaying "🐒 BABOON - Typing Practice" centered
  2. **Content**: Main content vertically centered in the available space between header and footer
  3. **Footer**: Fixed at the bottom of the terminal (last line), displaying context-sensitive help text
- The header and footer SHALL remain at the terminal boundaries regardless of window size
- Main content SHALL be horizontally centered within the terminal width
- Header text SHALL be cyan (colour 14) and bold
- Footer text SHALL be gray (colour 8) for help hints

### TR-003: Update Loop
- The application SHALL use tea.Tick with 100ms interval for WPM bar updates
- Tick messages SHALL continue throughout the typing session
- Window resize messages SHALL update width/height for centering calculations

## Business Rules

### BR-001: Word Selection
- Words are selected randomly with replacement (same word may appear multiple times)
- All words in dictionary have equal probability of selection
- Empty words or whitespace-only words SHALL be skipped
- All words SHALL be converted to lowercase before use

### BR-002: WPM Calculation
- Formula: WPM = (correct_characters / 5) / minutes_elapsed
- Standard word length is defined as 5 characters
- Only correctly typed characters contribute to WPM
- Time measured from first correct keystroke to round completion

### BR-003: Accuracy Calculation
- Formula: Accuracy = (correct_characters / total_characters) × 100
- Every keystroke counts toward total_characters (including errors)
- Backspace removes the last character from consideration
- Extra characters beyond word length count as incorrect

### BR-004: Best Time Logic
- Best time is the LOWEST (fastest) completion time
- On first session, current time becomes best time
- Subsequent sessions only update best if time < current best

### BR-005: New Best Detection
- WPM: New best if current > historical best (strictly greater)
- Accuracy: New best if current > historical best (strictly greater)
- Time: New best if current <= historical best (lower is better)
- First session (when historical best is 0) triggers celebration

### BR-007: Celebration Trigger
- A celebration is triggered when the user achieves a new personal best WPM OR accuracy
- The celebration occurs BEFORE the results screen is displayed
- Users can skip the celebration at any time by pressing any key
- Only the best WPM value is displayed in the celebration message (even if both WPM and accuracy are new bests)

### BR-006: Relative Colour Scaling
- Per-item statistics (fingers, rows, letters) use relative colour scaling
- This ensures meaningful visual differentiation when absolute values are clustered
- Algorithm:
  1. Calculate min and max values across all items with data
  2. For each value, compute: normalized = (value - min) / (max - min)
  3. Map normalized (0.0-1.0) to gradient colour index (0-11)
  4. Index 0 = red (worst in range), Index 11 = green (best in range)
- Items with no data are shown in gray (colour 240)
- For inverted metrics (seek time: lower is better):
  1. Compute inverted_value = max - value + min
  2. Apply same normalization to inverted_value
- This prevents all items appearing the same colour when accuracy is e.g. 92-97%

## Word Dictionary

### British English Spellings
The dictionary SHALL use British English spellings:
- colour (not color)
- behaviour (not behavior)
- centre (not center)
- defence (not defense)
- favour (not favor)
- realise (not realize)
- organisation (not organization)
- recognise (not recognize)
- programme (not program)
- labour, honour, neighbour
- travelling
- theatre, metre, litre, fibre
- Words ending in -ise (apologise, capitalise, emphasise, etc.)

## Architecture

The application follows a clean backend/frontend separation with a well-defined API:

### Backend Package (`backend/`)
The backend handles all game logic, statistics tracking, and state management. The frontend communicates exclusively through the `GameAPI` interface.

**Key Components:**
- `api.go` - Defines the `GameAPI` interface and data types
- `engine.go` - Implements the game engine

**GameAPI Interface:**
```go
type GameAPI interface {
    // Game Lifecycle
    StartRound()

    // Input Handling
    ProcessKeystroke(char string) KeystrokeResult
    ProcessBackspace() bool
    ClearInput() bool
    ProcessSpace() SpaceResult

    // State Queries
    GetGameState() GameState
    GetSessionStats() *stats.Stats
    GetHistoricalStats() *stats.HistoricalStats

    // Persistence
    SaveStats() error
}
```

### Frontend Package (`frontend/`)
The TUI frontend handles all rendering, user input, and visual presentation. It communicates with the backend exclusively through the `GameAPI` interface.

**Key Components:**
- `model.go` - Bubble Tea model (Init, Update, View)
- `views.go` - All rendering functions (typing screen, results screen)
- `styles.go` - Lipgloss styles and colour definitions
- `animations.go` - Spring-based animation logic

### Web Frontend (`web/`)
The web frontend is a React application that provides the same functionality as the TUI but with a beautiful, modern web interface.

**Technology Stack:**
- **React 18**: Modern React with hooks for state management
- **Chakra UI 2.x**: Component library with dark theme support
- **Framer Motion**: Physics-based animations with spring dynamics
- **Custom Theme**: Dark theme with Kartoza brand colours (orange #D4922A and blue #4A90A4)

**Key Components:**
- `App.js` - Main application component with state management and screen routing
- `api.js` - REST API client for backend communication
- `theme.js` - Custom Chakra UI theme with dark mode and chunky button styles
- `components/WelcomeScreen.js` - Landing page with animated logo and game options
- `components/TypingScreen.js` - Main typing interface with BlockLetter physics and live WPM bar
- `components/ResultsScreen.js` - Statistics display with all typing theory metrics

**Features:**
- Large chunky block letters with spring-based physics animations
- Real-time colour feedback (green for correct, red for incorrect)
- Live WPM bar with gradient colouring during typing
- Animated transitions between screens
- Letter statistics grid with accuracy and speed indicators
- Finger accuracy display and hand balance statistics
- Common error pattern display
- Responsive design for various screen sizes

**Physics Animations:**
- Letters bounce and scale on correct/incorrect keystrokes using Framer Motion springs
- Stat cards slide in with staggered spring animations
- Progress bars animate smoothly with spring dynamics
- UI elements have hover/tap scaling effects

### Data Flow

**Embedded TUI Mode (default):**
1. User input → Frontend Model → Game Engine (direct function calls)
2. Game state changes → Frontend queries engine → Render updated view
3. Statistics persist through Engine → Stats package → JSON file

**Web/Client Mode:**
1. User input → Frontend Model → REST Client → HTTP → REST Server → Game Engine
2. Game state changes → REST Client queries API → HTTP Response → Render updated view
3. Statistics persist through REST Server → Stats package → JSON file

### REST API

The backend exposes a RESTful API that the frontend communicates with via HTTP. This enables:
- Clean separation between frontend and backend processes
- Multiple concurrent frontend sessions connecting to a single backend
- Potential for alternative clients (web UI, mobile app, etc.)
- Network-based or remote gameplay

**Base URL:** `http://127.0.0.1:8787` (configurable via `-port` flag)

#### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/{provider}/login` | Initiate OAuth login (google, github, apple, microsoft) |
| GET | `/api/auth/{provider}/callback` | OAuth callback from provider |
| POST | `/api/auth/logout` | Logout (revoke tokens) |
| POST | `/api/auth/refresh` | Refresh access token |
| GET | `/api/auth/me` | Get current user info |
| DELETE | `/api/auth/account` | Delete user account |

#### User Stats Endpoints (Authenticated)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/user/stats` | Get authenticated user's stats |
| POST | `/api/user/stats/sync` | Merge local stats with server |
| DELETE | `/api/user/stats` | Delete all user stats |
| GET | `/api/user/export` | Export all user data as JSON |

#### Session Management

The REST API uses session-based routing. Each frontend client creates a session on startup and receives a unique session ID. All game operations are then scoped to that session.

**Session Lifecycle:**
1. Client calls `POST /api/sessions` to create a new session
2. Server responds with a unique `session_id` (32-character hex string)
3. Client includes session ID in all subsequent requests: `/api/sessions/{id}/...`
4. Client calls `DELETE /api/sessions/{id}` on exit to clean up

#### Endpoints

**Session Management:**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/sessions` | Create a new session |
| DELETE | `/api/sessions/{id}` | Delete a session |
| GET | `/api/sessions` | List all active sessions |
| GET | `/api/health` | Health check (includes active session count) |
| GET | `/api/config` | Get server configuration (AdSense key, etc.) |

**Game Operations (session-scoped):**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/sessions/{id}/round` | Start a new round |
| POST | `/api/sessions/{id}/keystroke` | Process a keystroke (with timing) |
| POST | `/api/sessions/{id}/backspace` | Process backspace |
| POST | `/api/sessions/{id}/clearinput` | Clear all typed characters for current word |
| POST | `/api/sessions/{id}/space` | Process space key (with timing) |
| POST | `/api/sessions/{id}/timing` | Submit final round timing data |
| GET | `/api/sessions/{id}/state` | Get current game state |
| GET | `/api/sessions/{id}/stats/session` | Get session statistics |
| GET | `/api/sessions/{id}/stats/historical` | Get historical statistics |
| POST | `/api/sessions/{id}/save` | Save statistics to disk |

#### Frontend Timing

All timing-critical measurements are performed on the frontend to avoid network latency affecting accuracy:

1. **Timer tracking**: The frontend tracks when the timer starts (first correct keystroke) and ends (round complete)
2. **Seek time measurement**: Time between keystrokes is measured locally and sent with each request
3. **Live WPM calculation**: Computed on the frontend using local timing data
4. **Duration calculation**: Total round duration is calculated on the frontend and submitted at round end

**How it works:**
- Each keystroke/space request includes a `seek_time_ms` field with the frontend-measured time since the previous keystroke
- When a round completes, the frontend calls `POST /api/sessions/{id}/timing` with:
  - `start_time_unix_ms`: Unix milliseconds when timer started
  - `end_time_unix_ms`: Unix milliseconds when round ended
  - `duration_ms`: Total duration in milliseconds
- The backend uses this timing data for WPM and accuracy calculations, ensuring network latency doesn't affect statistics

#### Request/Response Examples

**POST /api/sessions**
```json
// Request
{"punctuation_mode": false}

// Response (201 Created)
{"session_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"}
```

**GET /api/sessions**
```json
// Response
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

**GET /api/health**
```json
// Response
{"status": "healthy", "active_sessions": 3}
```

**POST /api/sessions/{id}/keystroke**
```json
// Request
{"char": "a"}

// Response
{"is_correct": true, "timer_started": true, "char_index": 0}
```

**POST /api/sessions/{id}/space**
```json
// Response
{"advanced": true, "round_complete": false, "treated_as_error": false}
```

**GET /api/sessions/{id}/state**
```json
// Response
{
  "words": ["hello", "world", ...],
  "current_word_idx": 0,
  "current_input": "hel",
  "timer_started": true,
  "punctuation_mode": false,
  "word_number": 1,
  "total_words": 30,
  "live_wpm": 45.2,
  "current_word": "hello",
  "previous_word": "",
  "next_word": "world",
  "next_words": ["world", "typing", "test"]
}
```

## File Structure

```
baboon/
├── flake.nix           # Nix flake for cross-platform builds
├── flake.lock          # Nix flake lock file
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── main.go             # Entry point - supports server, client, and combined modes
├── auth/                  # Authentication package
│   ├── auth.go            # Auth service and types
│   ├── errors.go          # Error definitions
│   ├── handlers.go        # HTTP handlers for auth endpoints
│   ├── jwt.go             # JWT generation and validation
│   ├── middleware.go      # Auth middleware
│   └── oauth.go           # OAuth provider configurations
├── database/              # Database package
│   ├── database.go        # Connection management and migrations
│   ├── stats.go           # User stats repository
│   ├── tokens.go          # Refresh token repository
│   └── users.go           # User repository
├── backend/
│   ├── api.go          # GameAPI interface and types (includes auth config)
│   ├── engine.go       # Game engine implementation
│   └── server.go       # REST API server with session management and auth
├── frontend/
│   ├── model.go        # Bubble Tea model with local timing
│   ├── views.go        # Rendering functions
│   ├── styles.go       # Lipgloss styles
│   ├── animations.go   # Spring animation logic
│   ├── celebration.go  # Personal best celebration particle system
│   ├── celebration_renderer.go  # Celebration screen rendering
│   └── client.go       # REST API client (implements GameAPI)
├── words/
│   └── words.go        # Dictionary of common words (British English)
├── stats/
│   ├── stats.go        # Statistics types, persistence, validation
│   └── keyboard.go     # Keyboard layout mappings (finger, hand, row)
├── settings/
│   └── settings.go     # User preferences (advance key setting)
├── scripts/
│   ├── start-backend.sh   # Start backend server in background
│   ├── stop-backend.sh    # Stop backend server
│   ├── status-backend.sh  # Check backend status and health
│   ├── launch-frontend.sh # Launch frontend client
│   └── release.sh         # Interactive version bump and release
├── web/                   # React web frontend
│   ├── package.json       # NPM dependencies
│   ├── package-lock.json  # NPM lockfile
│   ├── public/
│   │   └── index.html     # HTML template with fonts
│   └── src/
│       ├── index.js       # React entry point
│       ├── App.js         # Main application component
│       ├── api.js         # REST API client
│       ├── theme.js       # Chakra UI custom theme
│       └── components/
│           ├── WelcomeScreen.jsx  # Landing screen with local stats
│           ├── TypingScreen.jsx   # Typing practice screen
│           ├── ResultsScreen.jsx  # Statistics display
│           ├── AdSense.jsx        # Google AdSense component
│           ├── LoginButton.jsx    # OAuth login buttons
│           ├── UserMenu.jsx       # User profile dropdown
│           └── SyncDialog.jsx     # Stats merge confirmation dialog
│       └── contexts/
│           └── AuthContext.jsx    # Authentication state management
├── Makefile            # Build and run targets
├── SPECIFICATION.md    # This file
├── README.md           # User documentation
├── LICENSE             # MIT license
└── .gitignore          # Git ignore patterns
```

## Running Modes

The application supports multiple running modes:

### Embedded Mode (Default)
```bash
baboon              # TUI with embedded engine (single binary, no server)
baboon -p           # With punctuation mode
```
The TUI runs with the game engine embedded directly - no client-server architecture, no HTTP overhead. This is the simplest and most performant way to use Baboon. Stats persist to `~/.config/baboon/`.

### Server-Only Mode
```bash
baboon -server              # Run REST API server only (blocking)
baboon -server -port 9000   # On custom port
```
Runs the REST API server in the foreground. Required for the web frontend or for remote TUI clients. Writes PID to `$XDG_RUNTIME_DIR/baboon.pid`.

### Client-Only Mode
```bash
baboon -client              # Connect to existing backend
baboon -client -p           # With punctuation mode
baboon -client -port 9000   # Connect to custom port
```
Runs the TUI connecting to an already-running backend server via REST API. Useful when you want multiple clients to share a single backend, or when connecting remotely.

### Web Frontend Mode
```bash
make web-start              # Start backend + web frontend
make web-dev                # Start web dev server only (needs backend running)
make web-build              # Build for production
```
Starts the React web frontend on port 3000. The frontend proxies API requests to the backend on port 8787.

### Web Server Mode (Production)
```bash
baboon web                              # Serve web frontend + API on port 8787
baboon web -port 8080                   # On custom port
baboon web -adsense ca-pub-1234567890   # Enable Google AdSense ads
baboon web -dir ./web/dist              # Custom web directory
```
Serves the built web frontend with the REST API from a single binary. This is the recommended mode for production deployments. The `-adsense` flag enables Google AdSense advertising, which displays ads beneath the typing game.

## Management Scripts

Scripts are provided in the `scripts/` directory for managing the backend as a background service:

### start-backend.sh
Starts the backend server in the background.
```bash
./scripts/start-backend.sh           # Start on default port
./scripts/start-backend.sh -port 9000  # Custom port
./scripts/start-backend.sh -p        # With punctuation mode
```
- Checks if backend is already running
- Writes PID file for management
- Logs output to `$XDG_RUNTIME_DIR/baboon.log`

### stop-backend.sh
Stops the backend server gracefully.
```bash
./scripts/stop-backend.sh      # Graceful shutdown
./scripts/stop-backend.sh -f   # Force kill
```

### status-backend.sh
Checks the backend server status and health.
```bash
./scripts/status-backend.sh              # Check default port
./scripts/status-backend.sh -port 9000   # Check custom port
```
Shows: process status, health endpoint response, active session count.

### launch-frontend.sh
Launches a frontend client connected to the backend.
```bash
./scripts/launch-frontend.sh           # Connect to default port
./scripts/launch-frontend.sh -p        # With punctuation mode
./scripts/launch-frontend.sh -port 9000  # Connect to custom port
```
Checks that backend is running before launching.

## Release Management

### Version Bump and Release

The project provides an interactive release tool accessible via multiple methods:

```bash
# Via nix run (recommended)
nix run .#release

# Via make
make release

# Via direct script
./scripts/release.sh

# Check current version
make version
```

### Release Process

The release tool (`scripts/release.sh`) performs the following:

1. **Shows current state**:
   - Current version from `flake.nix`
   - Latest git tag
   - Uncommitted changes count

2. **Prompts for new version** with suggested options:
   - Patch bump (e.g., 1.4.0 → 1.4.1) for bugfixes
   - Minor bump (e.g., 1.4.0 → 1.5.0) for new features
   - Major bump (e.g., 1.4.0 → 2.0.0) for breaking changes
   - Custom version input

3. **Updates version**:
   - Modifies `version` in `flake.nix` (baboon package)
   - Creates commit: "Bump version to X.Y.Z"
   - Creates git tag: `vX.Y.Z`

4. **Optionally pushes** (with user confirmation):
   - Push to `origin main`
   - Push tag to trigger release workflow

### Automated Release Workflow

Pushing a tag matching `v*` triggers `.github/workflows/release.yml` which:

- Builds binaries for Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)
- Creates DEB packages for Debian/Ubuntu
- Creates RPM packages for Fedora/RHEL
- Creates Flatpak bundles
- Publishes GitHub release with all artifacts attached

### Version Location

The authoritative version is stored in `flake.nix` at the `buildGoModule` block:
```nix
packages = {
  default = pkgs.buildGoModule {
    pname = "baboon";
    version = "1.4.0";  # <-- Version defined here
    ...
  };
};
```

## Database Schema (Authentication)

When authentication is enabled via `BABOON_DATABASE_DSN`, the following tables are created:

### Users Table
```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    display_name TEXT,
    avatar_url TEXT,
    provider TEXT NOT NULL,      -- google, github, apple, microsoft
    provider_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    UNIQUE(provider, provider_id)
);
```

### User Stats Table
```sql
CREATE TABLE user_stats (
    user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    best_wpm REAL DEFAULT 0,
    best_accuracy REAL DEFAULT 0,
    best_time REAL DEFAULT 0,
    total_wpm REAL DEFAULT 0,
    total_accuracy REAL DEFAULT 0,
    total_time REAL DEFAULT 0,
    total_sessions INTEGER DEFAULT 0,
    last_session_date TIMESTAMP,
    -- Complex nested stats as JSON columns
    letter_accuracy TEXT,       -- JSON
    letter_seek_time TEXT,      -- JSON
    bigram_seek_time TEXT,      -- JSON
    finger_stats TEXT,          -- JSON
    hand_stats TEXT,            -- JSON
    row_stats TEXT,             -- JSON
    error_substitution TEXT,    -- JSON
    sfb_stats TEXT,             -- JSON
    hand_alternations INTEGER DEFAULT 0,
    same_hand_runs INTEGER DEFAULT 0,
    rhythm_stats TEXT,          -- JSON
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE
);
```

## Authentication Configuration

Authentication is configured via environment variables:

| Variable | Required | Description |
|----------|----------|-------------|
| `BABOON_DATABASE_DSN` | For auth | Database connection string (SQLite path or PostgreSQL URL) |
| `BABOON_JWT_SECRET` | For auth | Secret for signing JWTs (256-bit recommended) |
| `BABOON_BASE_URL` | For OAuth | Base URL for OAuth callbacks (e.g., `https://baboon.example.com`) |
| `BABOON_OAUTH_GOOGLE_CLIENT_ID` | Optional | Google OAuth client ID |
| `BABOON_OAUTH_GOOGLE_CLIENT_SECRET` | Optional | Google OAuth client secret |
| `BABOON_OAUTH_GITHUB_CLIENT_ID` | Optional | GitHub OAuth client ID |
| `BABOON_OAUTH_GITHUB_CLIENT_SECRET` | Optional | GitHub OAuth client secret |
| `BABOON_OAUTH_APPLE_CLIENT_ID` | Optional | Apple OAuth client ID |
| `BABOON_OAUTH_APPLE_CLIENT_SECRET` | Optional | Apple OAuth client secret |
| `BABOON_OAUTH_MICROSOFT_CLIENT_ID` | Optional | Microsoft OAuth client ID |
| `BABOON_OAUTH_MICROSOFT_CLIENT_SECRET` | Optional | Microsoft OAuth client secret |

### Graceful Degradation

- **No database configured**: Auth completely disabled, localStorage-only (current behaviour)
- **Database but no OAuth**: Database available but no login buttons shown
- **Partial OAuth**: Only configured providers shown (e.g., just Google)
- **Full OAuth**: All four providers available

### Database Drivers

- **SQLite**: Detected when DSN is a file path (e.g., `./baboon.db`)
- **PostgreSQL**: Detected when DSN starts with `postgres://` or `postgresql://`

## Stats File Format

Location: `~/.config/baboon/stats.json`

```json
{
  "best_wpm": 65.5,
  "best_accuracy": 98.2,
  "best_time": 45.3,
  "total_wpm": 850.5,
  "total_accuracy": 1420.8,
  "total_time": 725.0,
  "total_sessions": 15,
  "last_session_date": "2024-01-15T10:30:00Z",
  "letter_accuracy": {
    "a": {"presented": 100, "correct": 99},
    "b": {"presented": 45, "correct": 43}
  },
  "letter_seek_time": {
    "a": {"total_time_ms": 15000, "count": 100},
    "b": {"total_time_ms": 9000, "count": 45}
  },
  "bigram_seek_time": {
    "th": {"total_time_ms": 8500, "count": 50},
    "he": {"total_time_ms": 7200, "count": 48}
  },
  "finger_stats": {
    "0": {"presented": 200, "correct": 198, "total_time_ms": 30000, "count": 198},
    "1": {"presented": 180, "correct": 175, "total_time_ms": 27000, "count": 175}
  },
  "hand_stats": {
    "0": {"presented": 800, "correct": 790, "total_time_ms": 120000, "count": 790},
    "1": {"presented": 750, "correct": 740, "total_time_ms": 112500, "count": 740}
  },
  "row_stats": {
    "0": {"presented": 400, "correct": 395, "total_time_ms": 60000, "count": 395},
    "1": {"presented": 600, "correct": 595, "total_time_ms": 90000, "count": 595},
    "2": {"presented": 350, "correct": 340, "total_time_ms": 52500, "count": 340}
  },
  "error_substitution": {
    "a": {"s": 5, "q": 2},
    "e": {"r": 3, "w": 1}
  },
  "sfb_stats": {"count": 150, "total_time_ms": 45000},
  "hand_alternations": 1200,
  "same_hand_runs": 800,
  "rhythm_stats": {
    "total_seek_time_ms": 300000,
    "total_seek_time_sq": 75000000.0,
    "count": 2000,
    "last_variance": 0
  }
}
```

## Colour Palette Reference

### TUI Colour Codes

| Usage | Colour Code | Description |
|-------|-------------|-------------|
| Correct letter | 10 | Bright green |
| Incorrect letter | 9 | Bright red |
| Untyped letter | 8 | Gray |
| Title | 14 | Cyan |
| Labels | 7 | Light gray |
| Values | 15 | White |
| Session label | 6 | Cyan |
| New best star | 226 | Yellow |
| Help text | 8 | Gray |
| Empty bar | 236 | Dark gray |
| Gradient | 196→47 | Red through yellow to green |

### Web UI TUI-Style Design

The web frontend uses a TUI-consistent design that closely matches the terminal application.

**Design Principles:**
- Monospace fonts (Fira Code, JetBrains Mono) throughout
- Unicode block character rendering for word display
- TUI-style gradient progress bars using █ and ░ characters
- Terminal-inspired colour scheme with bright greens, reds, and cyan

**Block Font Rendering:**
- Letters rendered using Unicode block characters (█) matching the TUI
- Each letter is 6 lines tall with 1 character spacing
- Letters support A-Z uppercase for maximum readability
- BlockFont.jsx component provides all block font rendering

**Colour Scheme (TUI-Matching):**

| Colour | Hex Code | Usage |
|--------|----------|-------|
| Correct | #00ff00 | Correct keystrokes (ANSI 10) |
| Incorrect | #ff0000 | Incorrect keystrokes (ANSI 9) |
| Untyped | #808080 | Characters not yet typed (ANSI 8) |
| Current | #D4922A | Current letter to type (Kartoza orange) |
| Cyan | cyan.400 | Headers, titles, links |

**Gradient Progress Bar:**
The same red-to-green gradient used in the TUI (ANSI colours 196→47):
- #ff0000 (196), #ff5f00 (202), #ff8700 (208), #ffaf00 (214)
- #ffd700 (220), #ffff00 (226), #d7ff00 (190), #afff00 (154)
- #87ff00 (118), #5fff00 (82), #00ff00 (46), #00ff5f (47)

**Results Screen Layout:**
TUI-style label-value-bar columns with fixed character widths:
- Labels: 18 characters, right-aligned
- Values: 8 characters, right-aligned
- Bars: 30 character gradient bars

**Statistics Display:**
- Finger accuracy: LP LR LM LI RI RM RR RP with coloured dots
- Row accuracy: Top Home Bot with coloured dots
- Hand balance: L:XX% R:XX% with alternation percentage
- Letter matrix: A-Z header with accuracy/seek-time rows using coloured dots
- Common errors: expected→typed(count) format

**Brand Colour Palette (Orange)**
- 50: #fef6e9
- 100: #fce8c7
- 200: #f9d9a5
- 300: #f5c983
- 400: #e8a93d
- 500: #D4922A (primary)
- 600: #b87a22
- 700: #9c631a
- 800: #804c12
- 900: #64350a

**Brand Colour Palette (Blue)**
- 50: #e9f4f7
- 100: #c7e3ea
- 200: #a5d2dd
- 300: #83c1d0
- 400: #61b0c3
- 500: #4A90A4 (primary)
- 600: #3d7688
- 700: #305c6c
- 800: #234250
- 900: #162834

The Kartoza wallpaper (`web/public/kartoza-wallpaper.png`) is included in the project assets for reference.

## Version History

### v1.8.0
- TUI-style web UI redesign to match terminal application appearance
  - Block font rendering using Unicode characters (█) via BlockFont.jsx component
  - Same 6-line tall letters with 1 character spacing as TUI
  - Monospace font family (Fira Code, JetBrains Mono) throughout
  - TUI-matching colour scheme: #00ff00 (correct), #ff0000 (incorrect), #808080 (untyped)
- Results screen redesign matching TUI layout
  - Label-value-bar columns with fixed character widths (18-8-30)
  - Gradient progress bars using █ and ░ characters (ANSI 196→47 gradient)
  - Time bars inverted (lower is better)
  - Finger/row/hand accuracy displayed with coloured dots
  - Letter accuracy matrix (A-Z) with coloured dots
  - Common errors in expected→typed(count) format
- TypingScreen improvements
  - Block font word display with smooth carousel animations
  - TUI-style WPM bar with gradient and scale markers
  - Decorative separators matching terminal style
- WelcomeScreen redesign
  - Block font BABOON title
  - Terminal-inspired styling with cyan accents
  - Monospace font for all text

### v1.7.0
- User authentication system with OAuth SSO providers
  - Support for Google, GitHub, Apple, and Microsoft login
  - Self-hosted OAuth using `golang.org/x/oauth2` and JWT tokens
  - Database storage for user accounts and statistics (SQLite or PostgreSQL)
  - Cross-device stats synchronisation for authenticated users
- New backend packages:
  - `auth/`: Authentication service with JWT and OAuth support
  - `database/`: Database connectivity with SQLite and PostgreSQL drivers
- Frontend authentication components:
  - AuthContext for managing authentication state
  - LoginButton component for SSO provider buttons
  - UserMenu component for profile dropdown
  - SyncDialog for merging local stats with server on login
- Stats merge algorithm when user logs in with existing local stats:
  - Cumulative values (sessions, time): Added together
  - Best values (WPM, accuracy): Maximum kept
  - Map values (letter/finger stats): Merged additively
- Privacy features:
  - Data export endpoint for GDPR compliance
  - Account deletion with cascade delete
  - No keystroke logging - only aggregate statistics
- Graceful degradation:
  - Anonymous users continue using localStorage (no change)
  - Auth features hidden when not configured
  - Only configured OAuth providers shown
- New functional requirements: FR-032 through FR-036
- New user stories: US-007 (cross-device sync), US-008 (SSO login)

### v1.6.0
- Personal best celebration feature with fireworks animation
  - When user achieves a new best WPM or accuracy, a celebration screen is displayed
  - 8-second fireworks phase with particle explosions across the screen
  - 2-second message phase displaying "PERSONAL BEST" in block font with achieved WPM
  - 12 scheduled explosions with 30-50 particles each using celebratory colours
  - Physics-based particles with gravity, boundary bouncing, and text collision
  - Text destruction effect: particles destroy screen text creating debris particles
  - Skip celebration at any time by pressing any key
- New frontend files:
  - `celebration.go`: Particle system, physics simulation, explosion scheduling
  - `celebration_renderer.go`: Screen buffer rendering, particle-to-lipgloss colour mapping
- Added `StateCelebration` game state and `celebrationTickMsg` for 50 FPS updates
- Updated BR-005 (New Best Detection) with strict comparison for WPM/accuracy
- Added BR-007 (Celebration Trigger) business rule
- Added FR-030 (Personal Best Celebration) functional requirement

### v1.4.0
- TUI now uses embedded engine by default (no client-server architecture)
  - Running `baboon` directly embeds the game engine in the TUI binary
  - No HTTP server startup, no REST API overhead for TUI usage
  - Simpler, faster, single-binary experience
- Server mode (`-server`) still available for web frontend
- Client mode (`-client`) still available for connecting to remote servers
- Architecture documentation updated to reflect embedded mode

### v1.3.1
- Improved statistics colour display with relative scaling
  - Fixed issue where finger, row, and letter accuracy all appeared as identical green
  - When accuracy values are clustered (e.g., 92-97%), they now show meaningful colour differentiation
  - Lowest accuracy in the range shows red, highest shows green, with gradient between
  - Items with no data display in gray (colour 240)
- Added new colour function `GetRelativeColour()` for relative value visualization
- Updated FR-014 (Letter Statistics Display) with relative scaling reference
- Added FR-027 (Relative Colour Scaling for Statistics) as new functional requirement
- Added BR-006 (Relative Colour Scaling) with algorithm specification

### v1.3.0
- Options screen for configuring application settings
  - Accessible via 'o' key from typing screen (before timer starts) or results screen
  - Configurable advance key: Space (default), Enter, or Either
  - Navigation with ↑/↓ keys, Enter/Space to select, 1-3 for quick select
  - Settings persist to `~/.config/baboon/settings.json`
- New settings package (`settings/settings.go`) for user preferences
- Dynamic help text reflects currently configured advance key
- Added FR-026 (Options Screen) to functional requirements

### v1.2.0
- Show next 3 upcoming words in carousel display
  - Words displayed below the current word with decreasing opacity
  - First upcoming word has decorative arrows (▼), subsequent words shown plain
  - All upcoming words centered horizontally using lipgloss alignment
  - Provides better look-ahead for typing preparation
- Updated GameState API to include `NextWords` slice (array of up to 3 words)
- Backwards compatible: falls back to `NextWord` if `NextWords` is empty

### v1.1.0
- Beautiful carousel animation for word transitions
  - **Console (TUI)**: Smooth harmonica spring-based animations
    - Previous word fades in with animated greyscale opacity as it scrolls up
    - Current word slides up from below with spring physics
    - Next word fades in from below with staggered timing
    - Animation triggered on space key when advancing to next word
  - **Web**: Framer Motion spring animations
    - Previous word floats above at 50% scale with blur and fade
    - Current word displays large block letters with spring transitions
    - Next word floats below at 50% scale with blur
    - Decorative glow effect behind current word
- Fixed accuracy statistics exceeding 100%
  - Bug: When backspacing and retyping a character, "Correct" was counted multiple times while "Presented" was only counted once at round start
  - Fix: Track which character positions have been recorded as correct using a position map
  - Accuracy stats (letter, finger, hand, row) only recorded on first correct keystroke per position
  - Timing stats still recorded for all keystrokes (useful data regardless of retypes)
- Updated FR-001 to describe carousel word display layout

### v1.0.0
- First stable release
- Beautiful README with screenshots and badges
- GitHub Actions CI/CD workflows:
  - Test workflow: runs on push/PR, executes go test and go vet
  - Build workflow: cross-platform build verification (Linux, macOS, Windows)
  - Release workflow: automated builds on tag push with all package formats
- Pre-built binaries for multiple platforms:
  - Linux AMD64 and ARM64
  - macOS Intel and Apple Silicon
  - Windows AMD64
  - DEB package for Debian/Ubuntu
  - RPM package for Fedora/RHEL
  - Flatpak package
- macOS unsigned binary instructions in README
- Nix flake integration for system configurations

### v0.9.1
- Kartoza brand colour scheme applied to web frontend
  - Primary colour: Kartoza Orange (#D4922A)
  - Secondary colour: Kartoza Blue (#4A90A4)
  - Updated theme.js with full Kartoza colour palettes
  - Updated WelcomeScreen gradient title with brand colours
  - Updated TypingScreen with orange current letter highlight and blue/orange progress bars
  - Updated ResultsScreen with brand colours for stats, hand balance, and heatmaps
  - Kartoza wallpaper added to project assets (web/public/kartoza-wallpaper.png)

### v0.9.0
- React web frontend with 100% feature parity to TUI
  - Built with React 18, Chakra UI 2.x, and Framer Motion
  - Physics-based animations using Framer Motion spring dynamics
  - Large chunky block letters with bounce effects on keystrokes
  - Dark theme with Kartoza brand colours (orange and blue)
  - Custom theme with chunky button styles ("glow" and "chunky" variants)
- Web frontend components:
  - WelcomeScreen: Animated logo, connection status, game options
  - TypingScreen: Block letters with physics, live WPM bar, progress indicator
  - ResultsScreen: Full statistics display with animated stat cards
- Letter statistics grid with colour-coded accuracy and speed indicators
- Finger accuracy display and hand balance statistics
- Common error pattern tracking and display
- Responsive design for various screen sizes
- Makefile targets for web development:
  - `make web-install` - Install NPM dependencies
  - `make web-dev` - Start development server
  - `make web-build` - Build for production
  - `make web-start` - Start backend + web frontend together
- Proxy configuration for development (port 3000 → 8787)

### v0.8.0
- Management scripts for backend server lifecycle
  - `start-backend.sh` - Start backend in background with PID tracking
  - `stop-backend.sh` - Graceful or forced shutdown
  - `status-backend.sh` - Health check and session monitoring
  - `launch-frontend.sh` - Launch frontend against running backend
- Three running modes added:
  - Combined mode (default): Backend + frontend in same process
  - Server-only mode (`-server`): Run backend only, blocking
  - Client-only mode (`-client`): Connect to existing backend
- PID file written to `$XDG_RUNTIME_DIR/baboon.pid` in server mode
- Graceful shutdown handling with SIGINT/SIGTERM

### v0.7.0
- Frontend timing implementation to eliminate network latency effects
  - All timing-critical measurements now performed on the frontend
  - Seek times measured locally and sent with keystroke/space requests
  - Live WPM calculated on frontend using local timing data
  - Round duration submitted via dedicated timing endpoint
- New API methods:
  - `ProcessKeystrokeWithTiming(char, seekTimeMs)` - keystroke with timing
  - `ProcessSpaceWithTiming(seekTimeMs)` - space with timing
  - `SubmitTiming(startTime, endTime, durationMs)` - submit final timing
- New REST endpoint: `POST /api/sessions/{id}/timing`
- Backend no longer calls `time.Now()` for timing-critical operations
- Ensures accurate WPM and seek time statistics regardless of network conditions

### v0.6.0
- Multi-client session management for REST API
  - Each frontend client creates a unique session on startup
  - Session IDs are 32-character hex strings generated cryptographically
  - All game operations are scoped to sessions: `/api/sessions/{id}/...`
  - Sessions are automatically cleaned up when clients disconnect
- New session management endpoints:
  - `POST /api/sessions` - Create a new session
  - `DELETE /api/sessions/{id}` - Delete a session
  - `GET /api/sessions` - List all active sessions
- Health endpoint now reports active session count
- Thread-safe session storage with mutex protection
- Enables multiple concurrent players on a single backend server

### v0.5.0
- RESTful API implementation for frontend-backend communication
  - All frontend interactions routed through HTTP REST API
  - Backend server (`backend/server.go`) exposes REST endpoints
  - Frontend client (`frontend/client.go`) implements GameAPI via HTTP
- REST API endpoints:
  - `POST /api/round` - Start new round
  - `POST /api/keystroke` - Process keystroke
  - `POST /api/backspace` - Process backspace
  - `POST /api/space` - Process space
  - `GET /api/state` - Get game state
  - `GET /api/stats/session` - Get session statistics
  - `GET /api/stats/historical` - Get historical statistics
  - `POST /api/save` - Save statistics
  - `GET /api/health` - Health check
- Configurable port via `-port` flag (default: 8787)
- Thread-safe server with mutex protection

### v1.5.0
- Asciinema demo recording integration for terminal screencasts
  - `nix run .#demo-record` - Record terminal demo with asciinema
  - `nix run .#demo-play` - Play recorded demo locally
  - Automatic GIF conversion for README and documentation embedding
  - Demo GIF displayed prominently on README and documentation landing pages
  - Makefile targets: `make demo-record`, `make demo-play`
  - Neovim shortcuts: `<leader>par` (record), `<leader>pap` (play)
- Improved block letter font with smooth triangle corners
  - Uses filled triangles (◢ ◣ ◤ ◥) for smooth rounded corners
  - Added support for numbers 0-9 and additional punctuation (- ' ")
  - Eliminates blocky jagged appearance at letter edges

### v0.4.0
- Major architecture refactoring: clean backend/frontend separation
  - Backend package (`backend/`): Game engine with `GameAPI` interface
  - Frontend package (`frontend/`): TUI model, views, styles, animations
  - Main.go reduced to simple entry point wiring backend and frontend
- Clear API boundary: frontend communicates only through `GameAPI` interface
- Improved code organisation and maintainability

### v0.3.0
- Advanced typing theory statistics for effective touch typing practice:
  - Finger-specific accuracy and speed tracking (8 fingers mapped to QWERTY layout)
  - Keyboard row tracking (top, home, bottom row performance)
  - Hand balance and alternation rate tracking
  - Same-finger bigram (SFB) detection and timing
  - Rhythm consistency tracking (standard deviation of seek times)
  - Error substitution pattern tracking (which letters get confused)
- Results screen displays new typing theory metrics:
  - Finger accuracy row: LP LR LM LI | RI RM RR RP with colour-coded indicators
  - Row accuracy: Top, Home, Bot performance indicators
  - Hand balance: L:X% R:Y% distribution with alternation rate
  - Rhythm: Session StdDev vs historical average
  - Same-finger: SFB count and average timing
  - Common errors: Top 5 letter substitution patterns (e.g., a→s(12))
- Increased animated rows from 14 to 25 for new stats sections

### v0.2.0
- Per-letter accuracy tracking (tracks how often each letter a-z is presented and typed correctly)
- Letter statistics persist across sessions for cumulative tracking
- Results screen displays 26-letter accuracy row with red-to-green gradient
- Results screen displays 26-letter frequency row showing relative letter presentation counts
- Results screen displays 26-letter seek time row showing typing speed per letter
- Fixed 30 words with exactly 150 characters per round for consistent timing comparisons
- Smooth font with Unicode half-block characters (▀, ▄) for rounded edges
- Previous word displayed in top left, next word in top right during typing
- Adaptive word selection: weights by both letter frequency AND accuracy
- Words with low-accuracy letters are favoured to give more practice on weak letters
- Letter seek time tracking: measures time between keystrokes for each letter
- Improved seek time calculation:
  - Only records for correct keystrokes (not errors)
  - Records against expected letter (not what user typed)
  - Excludes first letter of each word (avoids word-reading time)
- Bigram (letter pair) seek time tracking: measures transition speed between letter pairs
- Letter statistics display redesigned: header row with letters, filled circles (●) for data
- Results screen animation: rows slide in sequentially using harmonica spring physics
- Punctuation mode (-p flag): words separated by random punctuation (, . ; : ! ?)

### v1.3.0
- MkDocs documentation site with Material theme
  - Comprehensive user documentation
  - Developer guides and API reference
  - Beautiful baboon mascot typing at keyboard
  - GitHub Pages deployment via GitHub Actions
  - Dark/light theme support
  - Full-text search
  - Responsive mobile design

### v0.1.0 (Initial Release)
- Basic typing practice with 30-word rounds
- Block letter word display with real-time colour feedback
- Live WPM bar during typing
- Results screen with paired comparison bars
- WPM, time, and accuracy tracking with averages
- Historical best comparison with star indicators
- British English word dictionary
- Cross-platform support via Nix flakes
- Statistics persistence with corruption detection
