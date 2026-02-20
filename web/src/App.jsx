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

const MotionBox = motion(Box);

function App() {
  const [screen, setScreen] = useState('welcome'); // welcome, typing, results
  const [punctuationMode, setPunctuationMode] = useState(false);
  const [gameState, setGameState] = useState(null);
  const [sessionStats, setSessionStats] = useState(null);
  const [historicalStats, setHistoricalStats] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Timing state (tracked on frontend to avoid latency)
  const [timerStarted, setTimerStarted] = useState(false);
  const [startTime, setStartTime] = useState(null);
  const [lastKeyTime, setLastKeyTime] = useState(null);
  const [correctChars, setCorrectChars] = useState(0);
  const [liveWpm, setLiveWpm] = useState(0);

  const toast = useToast();
  const wpmIntervalRef = useRef(null);

  // Check backend health and create session
  useEffect(() => {
    const init = async () => {
      try {
        await api.checkHealth();
        setIsConnected(true);
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
      await api.createSession(punctuationMode);
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

      // Refresh game state
      const state = await api.getState();
      setGameState(state);
    } catch (e) {
      console.error('Keystroke error:', e);
    }
  }, [lastKeyTime, timerStarted]);

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
        // Submit timing data
        const durationMs = startTime ? now - startTime : 0;
        await api.submitTiming(startTime, now, durationMs);
        await api.saveStats();

        // Get final stats
        const [session, historical] = await Promise.all([
          api.getSessionStats(),
          api.getHistoricalStats(),
        ]);
        setSessionStats(session);
        setHistoricalStats(historical);
        setScreen('results');
      } else {
        const state = await api.getState();
        setGameState(state);
      }
    } catch (e) {
      console.error('Space error:', e);
    }
  }, [lastKeyTime, startTime]);

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
  };

  if (isLoading && screen === 'welcome') {
    return (
      <Flex h="100vh" align="center" justify="center" bg="bg.primary">
        <VStack spacing={4}>
          <Spinner size="xl" color="accent.cyan" thickness="4px" />
          <Text color="gray.400">Connecting to backend...</Text>
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
              punctuationMode={punctuationMode}
              setPunctuationMode={setPunctuationMode}
              onStart={startGame}
              isLoading={isLoading}
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
              isLoading={isLoading}
            />
          </MotionBox>
        )}
      </AnimatePresence>
    </Box>
  );
}

export default App;
