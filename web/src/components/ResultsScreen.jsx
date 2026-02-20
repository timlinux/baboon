import React, { useEffect } from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Button,
  Flex,
  Grid,
  GridItem,
  Badge,
  Tooltip,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';

const MotionBox = motion(Box);
const MotionFlex = motion(Flex);

// Hero stat - big number with comparison inline
function HeroStat({ label, value, unit, best, avg, isBest, delay = 0, color = 'orange' }) {
  // Kartoza color scheme
  const colors = {
    orange: { bg: 'rgba(212, 146, 42, 0.1)', border: 'brand.500', text: 'brand.500' },
    blue: { bg: 'rgba(74, 144, 164, 0.1)', border: 'kartoza.blue.500', text: 'kartoza.blue.500' },
    green: { bg: 'rgba(76, 175, 80, 0.1)', border: 'accent.green', text: 'accent.green' },
  };
  const c = colors[color] || colors.orange;

  return (
    <MotionBox
      initial={{ opacity: 0, y: 20, scale: 0.9 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ delay, type: 'spring', bounce: 0.4 }}
      whileHover={{ scale: 1.02 }}
      flex={1}
    >
      <Box
        bg={c.bg}
        borderRadius="2xl"
        p={4}
        border="2px solid"
        borderColor={isBest ? 'accent.yellow' : c.border}
        boxShadow={isBest ? '0 0 20px rgba(255, 204, 0, 0.3)' : 'none'}
        position="relative"
        textAlign="center"
      >
        {isBest && (
          <Badge
            position="absolute"
            top={2}
            right={2}
            colorScheme="yellow"
            fontSize="xs"
          >
            NEW BEST
          </Badge>
        )}
        <Text color="gray.400" fontSize="xs" fontWeight="500" mb={1}>
          {label}
        </Text>
        <HStack justify="center" align="baseline" spacing={1}>
          <Text fontSize={{ base: '3xl', md: '4xl' }} fontWeight="800" color={c.text}>
            {typeof value === 'number' ? value.toFixed(1) : value}
          </Text>
          <Text color="gray.500" fontSize="md">{unit}</Text>
        </HStack>
        <HStack justify="center" spacing={3} mt={1}>
          <Text color="gray.500" fontSize="xs">
            Best: <Text as="span" color="accent.green">{best.toFixed(1)}</Text>
          </Text>
          <Text color="gray.500" fontSize="xs">
            Avg: <Text as="span" color="kartoza.blue.500">{avg.toFixed(1)}</Text>
          </Text>
        </HStack>
      </Box>
    </MotionBox>
  );
}

// Compact letter heatmap - keyboard-style layout
function LetterHeatmap({ letterAccuracy, letterSeekTime, delay = 0 }) {
  const rows = [
    ['q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p'],
    ['a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l'],
    ['z', 'x', 'c', 'v', 'b', 'n', 'm'],
  ];

  const getColor = (letter, type) => {
    // Kartoza color scheme for heatmap
    if (type === 'accuracy') {
      const stats = letterAccuracy?.[letter];
      if (!stats || stats.presented === 0) return 'gray.700';
      const acc = (stats.correct / stats.presented) * 100;
      if (acc >= 95) return '#4CAF50'; // Green
      if (acc >= 85) return '#8BC34A'; // Light green
      if (acc >= 75) return '#D4922A'; // Kartoza orange
      if (acc >= 60) return '#E65100'; // Dark orange
      return '#E53935'; // Red
    } else {
      const stats = letterSeekTime?.[letter];
      if (!stats || stats.count === 0) return 'gray.700';
      const avg = stats.total_time_ms / stats.count;
      if (avg <= 150) return '#4CAF50'; // Green
      if (avg <= 200) return '#8BC34A'; // Light green
      if (avg <= 250) return '#D4922A'; // Kartoza orange
      if (avg <= 350) return '#E65100'; // Dark orange
      return '#E53935'; // Red
    }
  };

  const getTooltip = (letter, type) => {
    if (type === 'accuracy') {
      const stats = letterAccuracy?.[letter];
      if (!stats || stats.presented === 0) return `${letter.toUpperCase()}: N/A`;
      return `${letter.toUpperCase()}: ${((stats.correct / stats.presented) * 100).toFixed(0)}%`;
    } else {
      const stats = letterSeekTime?.[letter];
      if (!stats || stats.count === 0) return `${letter.toUpperCase()}: N/A`;
      return `${letter.toUpperCase()}: ${(stats.total_time_ms / stats.count).toFixed(0)}ms`;
    }
  };

  const KeyboardRow = ({ letters, type, rowDelay }) => (
    <HStack spacing={1} justify="center">
      {letters.map((letter, i) => (
        <MotionBox
          key={letter}
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ delay: rowDelay + i * 0.02 }}
        >
          <Tooltip label={getTooltip(letter, type)} placement="top" hasArrow>
            <Flex
              w={{ base: '22px', md: '26px' }}
              h={{ base: '22px', md: '26px' }}
              align="center"
              justify="center"
              bg={getColor(letter, type)}
              borderRadius="md"
              opacity={
                (type === 'accuracy' ? letterAccuracy?.[letter]?.presented : letterSeekTime?.[letter]?.count) > 0
                  ? 1
                  : 0.3
              }
            >
              <Text fontSize="xs" fontWeight="bold" color="gray.900">
                {letter.toUpperCase()}
              </Text>
            </Flex>
          </Tooltip>
        </MotionBox>
      ))}
    </HStack>
  );

  return (
    <MotionBox
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay }}
    >
      <HStack spacing={6} align="start">
        {/* Accuracy keyboard */}
        <VStack spacing={1}>
          <Text color="gray.400" fontSize="xs" fontWeight="500" mb={1}>Accuracy</Text>
          {rows.map((row, idx) => (
            <Box key={idx} pl={idx === 1 ? 2 : idx === 2 ? 4 : 0}>
              <KeyboardRow letters={row} type="accuracy" rowDelay={delay + idx * 0.1} />
            </Box>
          ))}
        </VStack>

        {/* Speed keyboard */}
        <VStack spacing={1}>
          <Text color="gray.400" fontSize="xs" fontWeight="500" mb={1}>Speed</Text>
          {rows.map((row, idx) => (
            <Box key={idx} pl={idx === 1 ? 2 : idx === 2 ? 4 : 0}>
              <KeyboardRow letters={row} type="speed" rowDelay={delay + 0.3 + idx * 0.1} />
            </Box>
          ))}
        </VStack>
      </HStack>
    </MotionBox>
  );
}

// Compact finger stats - inline display
function CompactFingerStats({ fingerStats, delay = 0 }) {
  const fingers = [
    { id: 0, label: 'LP' },
    { id: 1, label: 'LR' },
    { id: 2, label: 'LM' },
    { id: 3, label: 'LI' },
    { id: 6, label: 'RI' },
    { id: 7, label: 'RM' },
    { id: 8, label: 'RR' },
    { id: 9, label: 'RP' },
  ];

  const getColor = (finger) => {
    const stats = fingerStats?.[finger];
    if (!stats || stats.presented === 0) return 'gray.600';
    const acc = (stats.correct / stats.presented) * 100;
    if (acc >= 95) return 'accent.green';
    if (acc >= 85) return 'kartoza.blue.500';
    if (acc >= 75) return 'brand.500'; // Kartoza orange
    return 'red.400';
  };

  return (
    <MotionBox
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay }}
    >
      <VStack align="stretch" spacing={2}>
        <Text color="gray.400" fontSize="xs" fontWeight="500">Finger Accuracy</Text>
        <HStack spacing={1} justify="center">
          {fingers.map((f, i) => {
            const stats = fingerStats?.[f.id];
            const acc = stats && stats.presented > 0
              ? ((stats.correct / stats.presented) * 100).toFixed(0)
              : '-';
            return (
              <MotionBox
                key={f.id}
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: delay + i * 0.03, type: 'spring' }}
              >
                <VStack spacing={0}>
                  <Text fontSize="2xs" color="gray.500">{f.label}</Text>
                  <Flex
                    w="28px"
                    h="28px"
                    align="center"
                    justify="center"
                    bg="bg.tertiary"
                    borderRadius="lg"
                    border="2px solid"
                    borderColor={getColor(f.id)}
                  >
                    <Text fontSize="xs" fontWeight="bold" color={getColor(f.id)}>
                      {acc}
                    </Text>
                  </Flex>
                </VStack>
              </MotionBox>
            );
          })}
        </HStack>
      </VStack>
    </MotionBox>
  );
}

// Compact hand balance
function CompactHandBalance({ handStats, handAlternations, sameHandRuns, delay = 0 }) {
  const leftStats = handStats?.[0] || { presented: 0 };
  const rightStats = handStats?.[1] || { presented: 0 };
  const total = leftStats.presented + rightStats.presented;
  const leftPct = total > 0 ? (leftStats.presented / total) * 100 : 50;
  const rightPct = 100 - leftPct;

  const totalTransitions = (handAlternations || 0) + (sameHandRuns || 0);
  const alternationRate = totalTransitions > 0
    ? ((handAlternations || 0) / totalTransitions) * 100
    : 0;

  return (
    <MotionBox
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay }}
    >
      <VStack align="stretch" spacing={2}>
        <Text color="gray.400" fontSize="xs" fontWeight="500">Hand Balance</Text>
        <HStack spacing={2}>
          <Text color="brand.500" fontSize="sm" fontWeight="bold" w="40px">
            L {leftPct.toFixed(0)}%
          </Text>
          <Box flex={1} h="10px" bg="bg.tertiary" borderRadius="full" overflow="hidden">
            <Flex h="100%">
              <MotionBox
                h="100%"
                bg="brand.500"
                initial={{ width: 0 }}
                animate={{ width: `${leftPct}%` }}
                transition={{ delay: delay + 0.1, type: 'spring' }}
              />
              <MotionBox
                h="100%"
                bg="kartoza.blue.500"
                initial={{ width: 0 }}
                animate={{ width: `${rightPct}%` }}
                transition={{ delay: delay + 0.2, type: 'spring' }}
              />
            </Flex>
          </Box>
          <Text color="kartoza.blue.500" fontSize="sm" fontWeight="bold" w="40px" textAlign="right">
            R {rightPct.toFixed(0)}%
          </Text>
        </HStack>
        <HStack justify="center">
          <Text color="gray.500" fontSize="xs">
            Alternation: <Text as="span" color="accent.green" fontWeight="bold">{alternationRate.toFixed(0)}%</Text>
          </Text>
        </HStack>
      </VStack>
    </MotionBox>
  );
}

// Compact common errors
function CompactErrors({ errorSubstitution, delay = 0 }) {
  const errors = [];
  if (errorSubstitution) {
    Object.entries(errorSubstitution).forEach(([expected, typed]) => {
      Object.entries(typed).forEach(([typedChar, count]) => {
        errors.push({ expected, typed: typedChar, count });
      });
    });
  }
  errors.sort((a, b) => b.count - a.count);
  const topErrors = errors.slice(0, 4);

  return (
    <MotionBox
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay }}
    >
      <VStack align="stretch" spacing={2}>
        <Text color="gray.400" fontSize="xs" fontWeight="500">Common Errors</Text>
        {topErrors.length === 0 ? (
          <Text color="gray.600" fontSize="xs" textAlign="center">No errors</Text>
        ) : (
          <HStack spacing={1} flexWrap="wrap" justify="center">
            {topErrors.map((error, i) => (
              <MotionBox
                key={`${error.expected}-${error.typed}`}
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: delay + i * 0.05, type: 'spring' }}
              >
                <Box
                  px={2}
                  py={1}
                  bg="rgba(255, 68, 102, 0.15)"
                  borderRadius="md"
                  border="1px solid"
                  borderColor="accent.red"
                >
                  <Text fontSize="xs" color="accent.red" fontWeight="500">
                    {error.expected}→{error.typed}
                    <Text as="span" color="gray.500" ml={1}>({error.count})</Text>
                  </Text>
                </Box>
              </MotionBox>
            ))}
          </HStack>
        )}
      </VStack>
    </MotionBox>
  );
}

function ResultsScreen({
  sessionStats,
  historicalStats,
  onNewRound,
  onBackToMenu,
  isLoading,
}) {
  const wpm = sessionStats?.wpm || 0;
  const accuracy = sessionStats?.accuracy || 0;
  const duration = sessionStats?.duration ? sessionStats.duration / 1e9 : 0;

  const bestWpm = historicalStats?.best_wpm || 0;
  const bestAccuracy = historicalStats?.best_accuracy || 0;
  const bestTime = historicalStats?.best_time || 0;

  const avgWpm = historicalStats?.total_sessions > 0
    ? historicalStats.total_wpm / historicalStats.total_sessions
    : 0;
  const avgAccuracy = historicalStats?.total_sessions > 0
    ? historicalStats.total_accuracy / historicalStats.total_sessions
    : 0;
  const avgTime = historicalStats?.total_sessions > 0
    ? historicalStats.total_time / historicalStats.total_sessions
    : 0;

  const isNewBestWpm = wpm >= bestWpm && historicalStats?.total_sessions > 0;
  const isNewBestAccuracy = accuracy >= bestAccuracy && historicalStats?.total_sessions > 0;
  const isNewBestTime = (bestTime === 0 || duration <= bestTime) && historicalStats?.total_sessions > 0;

  // Handle keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === 'Tab') {
        e.preventDefault();
        if (!isLoading) {
          onNewRound();
        }
      } else if (e.key === 'Escape') {
        onBackToMenu();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [onNewRound, onBackToMenu, isLoading]);

  return (
    <Flex minH="100vh" direction="column" p={4}>
      {/* Header */}
      <MotionBox
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ type: 'spring', bounce: 0.4 }}
        textAlign="center"
        py={2}
      >
        <Text
          fontSize={{ base: '3xl', md: '5xl' }}
          fontWeight="800"
          bgGradient="linear(to-r, brand.500, kartoza.blue.500)"
          bgClip="text"
        >
          Round Complete!
        </Text>
        <Text color="gray.500" fontSize="sm">
          Session #{historicalStats?.total_sessions || 1}
        </Text>
      </MotionBox>

      {/* Hero Stats Row */}
      <HStack spacing={4} py={4} px={{ base: 0, md: 8 }}>
        <HeroStat
          label="Words Per Minute"
          value={wpm}
          unit="WPM"
          best={bestWpm}
          avg={avgWpm}
          isBest={isNewBestWpm}
          color="orange"
          delay={0.1}
        />
        <HeroStat
          label="Accuracy"
          value={accuracy}
          unit="%"
          best={bestAccuracy}
          avg={avgAccuracy}
          isBest={isNewBestAccuracy}
          color="green"
          delay={0.15}
        />
        <HeroStat
          label="Time"
          value={duration}
          unit="s"
          best={bestTime}
          avg={avgTime}
          isBest={isNewBestTime}
          color="blue"
          delay={0.2}
        />
      </HStack>

      {/* Main Content - Two Column Grid */}
      <Grid
        templateColumns={{ base: '1fr', lg: '1fr 1fr' }}
        gap={4}
        flex={1}
        px={{ base: 0, md: 8 }}
      >
        {/* Left Column - Letter Heatmaps */}
        <GridItem>
          <Box
            bg="bg.card"
            borderRadius="2xl"
            p={4}
            border="1px solid"
            borderColor="whiteAlpha.100"
            h="100%"
          >
            <Text color="gray.300" fontSize="sm" fontWeight="600" mb={3}>
              Letter Performance
            </Text>
            <LetterHeatmap
              letterAccuracy={historicalStats?.letter_accuracy}
              letterSeekTime={historicalStats?.letter_seek_time}
              delay={0.3}
            />
          </Box>
        </GridItem>

        {/* Right Column - Typing Analysis */}
        <GridItem>
          <Box
            bg="bg.card"
            borderRadius="2xl"
            p={4}
            border="1px solid"
            borderColor="whiteAlpha.100"
            h="100%"
          >
            <Text color="gray.300" fontSize="sm" fontWeight="600" mb={3}>
              Typing Analysis
            </Text>
            <VStack spacing={4} align="stretch">
              <CompactFingerStats
                fingerStats={historicalStats?.finger_stats}
                delay={0.4}
              />
              <CompactHandBalance
                handStats={historicalStats?.hand_stats}
                handAlternations={historicalStats?.hand_alternations}
                sameHandRuns={historicalStats?.same_hand_runs}
                delay={0.5}
              />
              <CompactErrors
                errorSubstitution={historicalStats?.error_substitution}
                delay={0.6}
              />
            </VStack>
          </Box>
        </GridItem>
      </Grid>

      {/* Action Buttons - Fixed at bottom */}
      <MotionFlex
        gap={4}
        justify="center"
        py={4}
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.7, type: 'spring' }}
      >
        <MotionBox whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
          <Button
            size="lg"
            variant="glow"
            onClick={onNewRound}
            isLoading={isLoading}
            px={8}
          >
            New Round
          </Button>
        </MotionBox>

        <MotionBox whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
          <Button
            size="lg"
            variant="ghost"
            onClick={onBackToMenu}
            px={8}
          >
            Back to Menu
          </Button>
        </MotionBox>
      </MotionFlex>

      {/* Keyboard hint */}
      <Text color="gray.600" fontSize="xs" textAlign="center" pb={2}>
        Press TAB for new round • ESC to exit
      </Text>
    </Flex>
  );
}

export default ResultsScreen;
