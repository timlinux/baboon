import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import api from '../api.js';

// Auth context for managing user authentication state
const AuthContext = createContext(null);

// Local storage key for local stats (used before login and for anonymous users)
const LOCAL_STATS_KEY = 'baboon_historical_stats';

// Load local stats from localStorage
const loadLocalStats = () => {
  try {
    const item = localStorage.getItem(LOCAL_STATS_KEY);
    return item ? JSON.parse(item) : null;
  } catch (e) {
    console.error('Error loading local stats:', e);
    return null;
  }
};

// Save local stats to localStorage
const saveLocalStats = (stats) => {
  try {
    localStorage.setItem(LOCAL_STATS_KEY, JSON.stringify(stats));
  } catch (e) {
    console.error('Error saving local stats:', e);
  }
};

// Clear local stats from localStorage
const clearLocalStats = () => {
  try {
    localStorage.removeItem(LOCAL_STATS_KEY);
  } catch (e) {
    console.error('Error clearing local stats:', e);
  }
};

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [authConfig, setAuthConfig] = useState({
    auth_enabled: false,
    providers: [],
  });
  const [localStats, setLocalStats] = useState(() => loadLocalStats());
  const [showSyncDialog, setShowSyncDialog] = useState(false);
  const [pendingSyncData, setPendingSyncData] = useState(null);

  // Check if user is authenticated on mount and handle OAuth callback
  useEffect(() => {
    const init = async () => {
      setIsLoading(true);

      // Check URL for OAuth callback params
      const params = new URLSearchParams(window.location.search);
      const authSuccess = params.get('auth_success');
      const authError = params.get('auth_error');

      if (authError) {
        console.error('OAuth error:', authError);
        // Clear URL params
        window.history.replaceState({}, '', window.location.pathname);
      }

      if (authSuccess) {
        // Clear URL params
        window.history.replaceState({}, '', window.location.pathname);
      }

      try {
        // Get auth config
        const config = await api.getConfig();
        setAuthConfig({
          auth_enabled: config.auth_enabled || false,
          providers: config.auth_providers || [],
        });

        // Try to get current user
        if (config.auth_enabled) {
          try {
            const userData = await api.getCurrentUser();
            setUser(userData);
            setIsAuthenticated(true);

            // Check if we have local stats that need syncing
            const local = loadLocalStats();
            if (local && local.total_sessions > 0) {
              setPendingSyncData(local);
              setShowSyncDialog(true);
            }
          } catch (e) {
            // Not authenticated, that's fine
            setUser(null);
            setIsAuthenticated(false);
          }
        }
      } catch (e) {
        console.error('Failed to initialize auth:', e);
      }

      setIsLoading(false);
    };

    init();
  }, []);

  // Login with OAuth provider
  const login = useCallback((provider) => {
    // Redirect to OAuth login endpoint
    window.location.href = `/api/auth/${provider}/login`;
  }, []);

  // Logout
  const logout = useCallback(async () => {
    try {
      await api.logout();
    } catch (e) {
      console.error('Logout error:', e);
    }
    setUser(null);
    setIsAuthenticated(false);
  }, []);

  // Sync local stats with server
  const syncStats = useCallback(async (localStatsToSync) => {
    if (!isAuthenticated) return null;

    try {
      const result = await api.syncStats(localStatsToSync);
      // Clear local stats after successful sync
      clearLocalStats();
      setLocalStats(null);
      setPendingSyncData(null);
      setShowSyncDialog(false);
      return result.stats;
    } catch (e) {
      console.error('Sync error:', e);
      throw e;
    }
  }, [isAuthenticated]);

  // Skip syncing (keep server stats, discard local)
  const skipSync = useCallback(() => {
    clearLocalStats();
    setLocalStats(null);
    setPendingSyncData(null);
    setShowSyncDialog(false);
  }, []);

  // Update local stats (for anonymous users or before sync)
  const updateLocalStats = useCallback((stats) => {
    setLocalStats(stats);
    saveLocalStats(stats);
  }, []);

  // Refresh user data
  const refreshUser = useCallback(async () => {
    if (!authConfig.auth_enabled) return;

    try {
      const userData = await api.getCurrentUser();
      setUser(userData);
      setIsAuthenticated(true);
    } catch (e) {
      setUser(null);
      setIsAuthenticated(false);
    }
  }, [authConfig.auth_enabled]);

  const value = {
    user,
    isLoading,
    isAuthenticated,
    authConfig,
    localStats,
    showSyncDialog,
    pendingSyncData,
    login,
    logout,
    syncStats,
    skipSync,
    updateLocalStats,
    refreshUser,
    setShowSyncDialog,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

// Hook to use auth context
export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

export default AuthContext;
