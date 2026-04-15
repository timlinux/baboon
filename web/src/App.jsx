import React, { useState, useEffect, useCallback, useRef } from 'react';
import {
  Box,
  VStack,
  Text,
  Flex,
  useToast,
  Spinner,
} from '@chakra-ui/react';
import { motion, AnimatePresence } from 'framer-motion';
import api from './api.js';
import TypingScreen from './components/TypingScreen.jsx';
import ResultsScreen from './components/ResultsScreen.jsx';
import WelcomeScreen from './components/WelcomeScreen.jsx';
import LeaderboardScreen from './components/LeaderboardScreen.jsx';
import NameEntryScreen from './components/NameEntryScreen.jsx';
import AboutScreen from './components/AboutScreen.jsx';
import SyncDialog from './components/SyncDialog.jsx';
import { AuthProvider, useAuth } from './contexts/AuthContext.jsx';

const MotionBox = motion(Box);

// Local storage keys
const STORAGE_KEYS = {
  HISTORICAL_STATS: 'baboon_historical_stats',
};

// Load data from local storage
const loadFromStorage = (key, defaultValue = null) => {
  try {
    const item = localStorage.getItem(key);
    return item ? JSON.parse(item) : defaultValue;
  } catch (e) {
    console.error('Error loading from storage:', e);
    return defaultValue;
  }
};

// Save data to local storage
const saveToStorage = (key, value) => {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch (e) {
    console.error('Error saving to storage:', e);
  }
};

function App() {
  const { isAuthenticated, user, updateLocalStats, localStats: authLocalStats } = useAuth();

  const [screen, setScreen] = useState('welcome'); // welcome, typing, results, leaderboard, nameEntry
  const [leaderboardRank, setLeaderboardRank] = useState(null);
  const [pendingSubmission, setPendingSubmission] = useState(null); // Store session data for name entry
  const [gameState, setGameState] = useState(null);
  const [sessionStats, setSessionStats] = useState(null);
  const [historicalStats, setHistoricalStats] = useState(null);
  const [localHistoricalStats, setLocalHistoricalStats] = useState(() =>
    loadFromStorage(STORAGE_KEYS.HISTORICAL_STATS, null)
  );
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [config, setConfig] = useState({ adsense_enabled: false, adsense_key: '' });
  const [versionInfo, setVersionInfo] = useState({ version: '', git_commit: '' });

  // Timing state (tracked on frontend to avoid latency)
  const [timerStarted, setTimerStarted] = useState(false);
  const [startTime, setStartTime] = useState(null);
  const [lastKeyTime, setLastKeyTime] = useState(null);
  const [correctChars, setCorrectChars] = useState(0);
  const [liveWpm, setLiveWpm] = useState(0);

  const toast = useToast();
  const wpmIntervalRef = useRef(null);

  // Check backend health and fetch config (show menu first, don't auto-start game)
  useEffect(() => {
    const init = async () => {
      try {
        await api.checkHealth();
        setIsConnected(true);

        // Fetch server config (includes AdSense key if set)
        try {
          const serverConfig = await api.getConfig();
          setConfig(serverConfig);
        } catch (e) {
          console.log('Config endpoint not available, using defaults');
        }

        // Fetch version info
        try {
          const version = await api.getVersion();
          setVersionInfo(version);
        } catch (e) {
          console.log('Version endpoint not available');
        }
      } catch (e) {
        toast({
          title: 'Backend not available',
          description: 'Make sure the Baboon backend is running on port 8787',
          status: 'error',
          duration: null,
          isClosable: true,
        });
      }
      setIsLoading(false);
    };
    init();
  }, [toast]);

  // Live WPM calculation
  useEffect(() => {
    if (timerStarted && screen === 'typing') {
      wpmIntervalRef.current = setInterval(() => {
        if (startTime && correctChars > 0) {
          const elapsed = (Date.now() - startTime) / 1000 / 60; // minutes
          if (elapsed > 0) {
            setLiveWpm((correctChars / 5) / elapsed);
          }
        }
      }, 100);
    }
    return () => {
      if (wpmIntervalRef.current) {
        clearInterval(wpmIntervalRef.current);
      }
    };
  }, [timerStarted, screen, startTime, correctChars]);

  const startGame = async () => {
    try {
      setIsLoading(true);
      await api.createSession(false);
      await api.startRound();
      const state = await api.getState();
      setGameState(state);
      setTimerStarted(false);
      setStartTime(null);
      setLastKeyTime(null);
      setCorrectChars(0);
      setLiveWpm(0);
      setScreen('typing');
    } catch (e) {
      toast({
        title: 'Error starting game',
        description: e.message,
        status: 'error',
      });
    }
    setIsLoading(false);
  };

  // Common function to handle round completion
  const handleRoundComplete = useCallback(async (now) => {
    // Submit timing data
    const durationMs = startTime ? now - startTime : 0;
    await api.submitTiming(startTime || now, now, durationMs);
    await api.saveStats();

    // Get final stats
    const [session, historical] = await Promise.all([
      api.getSessionStats(),
      api.getHistoricalStats(),
    ]);
    setSessionStats(session);
    setHistoricalStats(historical);

    // Save stats to local storage for persistence
    setLocalHistoricalStats(historical);
    saveToStorage(STORAGE_KEYS.HISTORICAL_STATS, historical);

    // Also update auth context's local stats (for sync feature)
    if (updateLocalStats) {
      updateLocalStats(historical);
    }

    // If authenticated, sync to server
    if (isAuthenticated) {
      try {
        await api.syncStats(historical);
      } catch (e) {
        console.log('Failed to sync stats to server:', e);
        // Non-fatal - stats are still saved locally
      }
    }

    // Check if this score qualifies for the leaderboard (works for everyone)
    // Score is calculated as WPM * (accuracy/100)
    try {
      const wpm = session.wpm;
      const accuracy = session.accuracy;
      const qualification = await api.checkLeaderboardQualification(wpm, accuracy);
      if (qualification.qualifies) {
        setLeaderboardRank(qualification.rank);
        setPendingSubmission({
          wpm: session.wpm,
          accuracy: session.accuracy,
          duration: session.duration ? session.duration / 1e9 : 0,
        });
        setScreen('nameEntry');
        return;
      }
    } catch (e) {
      console.log('Failed to check leaderboard qualification:', e);
      // Non-fatal - continue to results (leaderboard may not be configured)
    }

    setScreen('results');
  }, [startTime, isAuthenticated, updateLocalStats]);

  const handleKeystroke = useCallback(async (char) => {
    const now = Date.now();
    const seekTimeMs = lastKeyTime ? now - lastKeyTime : 0;

    try {
      const result = await api.processKeystroke(char, seekTimeMs);

      // Start timer on first correct keystroke
      if (result.timer_started && !timerStarted) {
        setTimerStarted(true);
        setStartTime(now);
      }

      if (result.is_correct) {
        setCorrectChars(c => c + 1);
      }

      setLastKeyTime(now);

      // Check if round completed (last correct char of last word)
      if (result.round_complete) {
        await handleRoundComplete(now);
      } else {
        // Refresh game state
        const state = await api.getState();
        setGameState(state);
      }
    } catch (e) {
      console.error('Keystroke error:', e);
    }
  }, [lastKeyTime, timerStarted, handleRoundComplete]);

  const handleBackspace = useCallback(async () => {
    try {
      await api.processBackspace();
      const state = await api.getState();
      setGameState(state);
    } catch (e) {
      console.error('Backspace error:', e);
    }
  }, []);

  const handleSpace = useCallback(async () => {
    const now = Date.now();
    const seekTimeMs = lastKeyTime ? now - lastKeyTime : 0;

    try {
      const result = await api.processSpace(seekTimeMs);
      setLastKeyTime(now);

      if (result.round_complete) {
        await handleRoundComplete(now);
      } else {
        const state = await api.getState();
        setGameState(state);
      }
    } catch (e) {
      console.error('Space error:', e);
    }
  }, [lastKeyTime, handleRoundComplete]);

  const handleNewRound = async () => {
    try {
      setIsLoading(true);
      await api.startRound();
      const state = await api.getState();
      setGameState(state);
      setTimerStarted(false);
      setStartTime(null);
      setLastKeyTime(null);
      setCorrectChars(0);
      setLiveWpm(0);
      setScreen('typing');
    } catch (e) {
      toast({
        title: 'Error starting round',
        description: e.message,
        status: 'error',
      });
    }
    setIsLoading(false);
  };

  const handleBackToMenu = async () => {
    await api.deleteSession();
    setScreen('welcome');
    setGameState(null);
    setSessionStats(null);
    setHistoricalStats(null);
    setLeaderboardRank(null);
    setPendingSubmission(null);
  };

  const handleLeaderboardSubmit = async (displayName) => {
    if (!pendingSubmission) return;

    await api.submitLeaderboardEntry(
      displayName,
      pendingSubmission.wpm,
      pendingSubmission.accuracy,
      pendingSubmission.duration
    );
    setPendingSubmission(null);
    setScreen('results');
  };

  const handleSkipLeaderboard = () => {
    setPendingSubmission(null);
    setScreen('results');
  };

  const handleShowLeaderboard = () => {
    setScreen('leaderboard');
  };

  const handleShowAbout = () => {
    setScreen('about');
  };

  if (isLoading || (screen === 'typing' && !gameState)) {
    return (
      <Flex h="100vh" align="center" justify="center" bg="bg.primary">
        <VStack spacing={4}>
          <Spinner size="xl" color="accent.cyan" thickness="4px" />
          <Text color="gray.400" fontFamily="'Fira Code', monospace" fontSize={{ base: 'lg', md: 'xl' }}>Connecting to backend...</Text>
        </VStack>
      </Flex>
    );
  }

  return (
    <Box
      minH="100vh"
      bg="bg.primary"
      bgGradient="radial(ellipse at top, bg.secondary, bg.primary)"
      overflow="hidden"
    >
      <AnimatePresence mode="wait">
        {screen === 'welcome' && (
          <MotionBox
            key="welcome"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.5 }}
          >
            <WelcomeScreen
              isConnected={isConnected}
              onStart={startGame}
              onShowLeaderboard={handleShowLeaderboard}
              onShowAbout={handleShowAbout}
              isLoading={isLoading}
              localStats={localHistoricalStats}
              adsenseEnabled={config.adsense_enabled}
              adsenseKey={config.adsense_key}
              versionInfo={versionInfo}
            />
          </MotionBox>
        )}

        {screen === 'typing' && gameState && (
          <MotionBox
            key="typing"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 1.05 }}
            transition={{ duration: 0.4 }}
          >
            <TypingScreen
              gameState={gameState}
              liveWpm={liveWpm}
              timerStarted={timerStarted}
              onKeystroke={handleKeystroke}
              onBackspace={handleBackspace}
              onSpace={handleSpace}
              onExit={handleBackToMenu}
              onRestart={handleNewRound}
              adsenseEnabled={config.adsense_enabled}
              adsenseKey={config.adsense_key}
            />
          </MotionBox>
        )}

        {screen === 'results' && sessionStats && historicalStats && (
          <MotionBox
            key="results"
            initial={{ opacity: 0, y: 50 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -50 }}
            transition={{ duration: 0.5, type: 'spring', bounce: 0.3 }}
          >
            <ResultsScreen
              sessionStats={sessionStats}
              historicalStats={historicalStats}
              onNewRound={handleNewRound}
              onBackToMenu={handleBackToMenu}
              onShowLeaderboard={handleShowLeaderboard}
              isLoading={isLoading}
              adsenseEnabled={config.adsense_enabled}
              adsenseKey={config.adsense_key}
            />
          </MotionBox>
        )}

        {screen === 'nameEntry' && pendingSubmission && (
          <MotionBox
            key="nameEntry"
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 1.1 }}
            transition={{ duration: 0.4, type: 'spring', bounce: 0.3 }}
          >
            <NameEntryScreen
              wpm={pendingSubmission.wpm}
              accuracy={pendingSubmission.accuracy}
              rank={leaderboardRank}
              onSubmit={handleLeaderboardSubmit}
              onSkip={handleSkipLeaderboard}
              adsenseEnabled={config.adsense_enabled}
              adsenseKey={config.adsense_key}
            />
          </MotionBox>
        )}

        {screen === 'leaderboard' && (
          <MotionBox
            key="leaderboard"
            initial={{ opacity: 0, x: 50 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -50 }}
            transition={{ duration: 0.4 }}
          >
            <LeaderboardScreen
              onBack={() => setScreen(sessionStats ? 'results' : 'welcome')}
              adsenseEnabled={config.adsense_enabled}
              adsenseKey={config.adsense_key}
            />
          </MotionBox>
        )}

        {screen === 'about' && (
          <MotionBox
            key="about"
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -30 }}
            transition={{ duration: 0.4 }}
          >
            <AboutScreen
              onBack={() => setScreen('welcome')}
              adsenseEnabled={config.adsense_enabled}
              adsenseKey={config.adsense_key}
            />
          </MotionBox>
        )}
      </AnimatePresence>

      {/* Sync Dialog for merging local stats with server */}
      <SyncDialog />
    </Box>
  );
}

// Wrap App with AuthProvider
function AppWithAuth() {
  return (
    <AuthProvider>
      <App />
    </AuthProvider>
  );
}

export default AppWithAuth;
