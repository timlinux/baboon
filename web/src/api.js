// API client for Baboon backend

const API_BASE = '/api';

class BaboonAPI {
  constructor() {
    this.sessionId = null;
    this.baseUrl = API_BASE;
  }

  // Helper to make authenticated requests
  async _fetch(url, options = {}) {
    const response = await fetch(url, {
      ...options,
      credentials: 'include', // Include cookies for auth
    });
    return response;
  }

  async _fetchJSON(url, options = {}) {
    const response = await this._fetch(url, options);

    // Check content-type to ensure we're getting JSON
    const contentType = response.headers.get('content-type');
    const isJson = contentType && contentType.includes('application/json');

    if (!response.ok) {
      if (isJson) {
        const error = await response.json().catch(() => ({ error: response.statusText }));
        throw new Error(error.error || 'Request failed');
      } else {
        // Non-JSON error response (e.g., HTML 404 page)
        throw new Error(`Request failed: ${response.status} ${response.statusText}`);
      }
    }

    // Ensure response is JSON before parsing
    if (!isJson) {
      throw new Error('Server returned non-JSON response. The API may not be available.');
    }

    return response.json();
  }

  async createSession(punctuationMode = false) {
    const response = await fetch(`${this.baseUrl}/sessions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ punctuation_mode: punctuationMode }),
    });
    const data = await response.json();
    this.sessionId = data.session_id;
    return data;
  }

  async deleteSession() {
    if (!this.sessionId) return;
    await fetch(`${this.baseUrl}/sessions/${this.sessionId}`, {
      method: 'DELETE',
    });
    this.sessionId = null;
  }

  async startRound() {
    const response = await fetch(`${this.baseUrl}/sessions/${this.sessionId}/round`, {
      method: 'POST',
    });
    return response.json();
  }

  async processKeystroke(char, seekTimeMs = 0) {
    const response = await fetch(`${this.baseUrl}/sessions/${this.sessionId}/keystroke`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ char, seek_time_ms: seekTimeMs }),
    });
    return response.json();
  }

  async processBackspace() {
    const response = await fetch(`${this.baseUrl}/sessions/${this.sessionId}/backspace`, {
      method: 'POST',
    });
    return response.json();
  }

  async processSpace(seekTimeMs = 0) {
    const response = await fetch(`${this.baseUrl}/sessions/${this.sessionId}/space`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ seek_time_ms: seekTimeMs }),
    });
    return response.json();
  }

  async submitTiming(startTimeMs, endTimeMs, durationMs) {
    const response = await fetch(`${this.baseUrl}/sessions/${this.sessionId}/timing`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        start_time_unix_ms: startTimeMs,
        end_time_unix_ms: endTimeMs,
        duration_ms: durationMs,
      }),
    });
    return response.json();
  }

  async getState() {
    const response = await fetch(`${this.baseUrl}/sessions/${this.sessionId}/state`);
    return response.json();
  }

  async getSessionStats() {
    const response = await fetch(`${this.baseUrl}/sessions/${this.sessionId}/stats/session`);
    return response.json();
  }

  async getHistoricalStats() {
    const response = await fetch(`${this.baseUrl}/sessions/${this.sessionId}/stats/historical`);
    return response.json();
  }

  async saveStats() {
    const response = await fetch(`${this.baseUrl}/sessions/${this.sessionId}/save`, {
      method: 'POST',
    });
    return response.json();
  }

  async checkHealth() {
    const response = await fetch(`${this.baseUrl}/health`);
    return response.json();
  }

  async getConfig() {
    const response = await fetch(`${this.baseUrl}/config`);
    return response.json();
  }

  async getVersion() {
    const response = await fetch(`${this.baseUrl}/version`);
    return response.json();
  }

  // Authentication methods

  // Get current authenticated user
  async getCurrentUser() {
    return this._fetchJSON(`${this.baseUrl}/auth/me`);
  }

  // Logout (revoke tokens)
  async logout() {
    return this._fetchJSON(`${this.baseUrl}/auth/logout`, {
      method: 'POST',
    });
  }

  // Refresh access token
  async refreshToken() {
    return this._fetchJSON(`${this.baseUrl}/auth/refresh`, {
      method: 'POST',
    });
  }

  // Get user stats (authenticated)
  async getUserStats() {
    return this._fetchJSON(`${this.baseUrl}/user/stats`);
  }

  // Sync local stats with server (merge)
  async syncStats(localStats) {
    return this._fetchJSON(`${this.baseUrl}/user/stats/sync`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ local_stats: localStats }),
    });
  }

  // Delete user stats
  async deleteUserStats() {
    return this._fetchJSON(`${this.baseUrl}/user/stats`, {
      method: 'DELETE',
    });
  }

  // Export user data
  async exportUserData() {
    const response = await this._fetch(`${this.baseUrl}/user/export`);
    if (!response.ok) {
      throw new Error('Export failed');
    }
    return response.blob();
  }

  // Delete user account
  async deleteAccount() {
    return this._fetchJSON(`${this.baseUrl}/auth/account`, {
      method: 'DELETE',
    });
  }

  // Leaderboard methods

  // Get leaderboard entries for a month (default: current month)
  async getLeaderboard(monthYear = null) {
    const params = monthYear ? `?month=${monthYear}` : '';
    return this._fetchJSON(`${this.baseUrl}/v1/leaderboard${params}`);
  }

  // Get available months with leaderboard entries
  async getLeaderboardMonths() {
    return this._fetchJSON(`${this.baseUrl}/v1/leaderboard/months`);
  }

  // Check if a WPM and accuracy would qualify for top 10
  // Score is calculated as WPM * (accuracy/100)
  async checkLeaderboardQualification(wpm, accuracy = 100) {
    return this._fetchJSON(`${this.baseUrl}/v1/leaderboard/check?wpm=${wpm}&accuracy=${accuracy}`);
  }

  // Submit a leaderboard entry (requires authentication)
  async submitLeaderboardEntry(displayName, wpm, accuracy, durationSeconds) {
    return this._fetchJSON(`${this.baseUrl}/v1/leaderboard/submit`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        display_name: displayName,
        wpm,
        accuracy,
        duration_seconds: durationSeconds,
      }),
    });
  }

  // Validate a display name for profanity
  async validateDisplayName(name) {
    return this._fetchJSON(`${this.baseUrl}/v1/leaderboard/validate?name=${encodeURIComponent(name)}`);
  }

  // Get badge URL for a leaderboard entry
  getBadgeUrl(entryId) {
    return `${this.baseUrl}/v1/leaderboard/badge/${entryId}.svg`;
  }
}

export const api = new BaboonAPI();
export default api;
