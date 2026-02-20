// API client for Baboon backend

const API_BASE = '/api';

class BaboonAPI {
  constructor() {
    this.sessionId = null;
    this.baseUrl = API_BASE;
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
}

export const api = new BaboonAPI();
export default api;
