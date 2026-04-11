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
    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || 'Request failed');
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
}

export const api = new BaboonAPI();
export default api;
