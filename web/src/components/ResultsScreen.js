import React from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Button,
  Container,
  Flex,
  SimpleGrid,
  Badge,
  Tooltip,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';

const MotionBox = motion(Box);
const MotionFlex = motion(Flex);

// Animated stat card with physics
function StatCard({ label, value, unit, comparison, isBest, delay = 0, color = 'cyan' }) {
  const colors = {
    cyan: { bg: 'rgba(0, 204, 255, 0.1)', border: 'accent.cyan', text: 'accent.cyan' },
    green: { bg: 'rgba(0, 255, 136, 0.1)', border: 'accent.green', text: 'accent.green' },
    purple: { bg: 'rgba(170, 102, 255, 0.1)', border: 'accent.purple', text: 'accent.purple' },
    yellow: { bg: 'rgba(255, 204, 0, 0.1)', border: 'accent.yellow', text: 'accent.yellow' },
    orange: { bg: 'rgba(255, 136, 68, 0.1)', border: 'accent.orange', text: 'accent.orange' },
  };

  const c = colors[color] || colors.cyan;

  return (
    <MotionBox
      initial={{ opacity: 0, y: 30, scale: 0.9 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ delay, type: 'spring', bounce: 0.4 }}
      whileHover={{ scale: 1.02, y: -5 }}
    >
      <Box
        bg={c.bg}
        borderRadius="2xl"
        p={6}
        border="2px solid"
        borderColor={isBest ? 'accent.yellow' : c.border}
        position="relative"
        overflow="hidden"
        boxShadow={isBest ? '0 0 30px rgba(255, 204, 0, 0.3)' : 'none'}
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
        <VStack spacing={1} align="start">
          <Text color="gray.400" fontSize="sm" fontWeight="500">
            {label}
          </Text>
          <HStack align="baseline" spacing={1}>
            <Text fontSize="4xl" fontWeight="800" color={c.text}>
              {typeof value === 'number' ? value.toFixed(1) : value}
            </Text>
            {unit && (
              <Text color="gray.500" fontSize="lg">
                {unit}
              </Text>
            )}
          </HStack>
          {comparison && (
            <Text color="gray.500" fontSize="xs">
              {comparison}
            </Text>
          )}
        </VStack>
      </Box>
    </MotionBox>
  );
}

// Progress bar for comparing values
function ComparisonBar({ current, best, average, label, inverted = false, delay = 0 }) {
  const maxVal = Math.max(current, best, average, 1);
  const scale = inverted ? (v) => ((maxVal - v) / maxVal) * 100 : (v) => (v / maxVal) * 100;

  return (
    <MotionBox
      initial={{ opacity: 0, x: -30 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay, type: 'spring' }}
    >
      <VStack align="stretch" spacing={2}>
        <Text color="gray.400" fontSize="sm">{label}</Text>

        <HStack spacing={4}>
          <Text color="gray.500" fontSize="xs" w="50px">This run</Text>
          <Box flex={1} h="8px" bg="bg.tertiary" borderRadius="full" overflow="hidden">
            <MotionBox
              h="100%"
              bg="accent.cyan"
              borderRadius="full"
              initial={{ width: 0 }}
              animate={{ width: `${scale(current)}%` }}
              transition={{ delay: delay + 0.2, type: 'spring' }}
            />
          </Box>
          <Text color="white" fontSize="sm" fontWeight="bold" w="60px" textAlign="right">
            {current.toFixed(1)}
          </Text>
        </HStack>

        <HStack spacing={4}>
          <Text color="gray.500" fontSize="xs" w="50px">Best</Text>
          <Box flex={1} h="8px" bg="bg.tertiary" borderRadius="full" overflow="hidden">
            <MotionBox
              h="100%"
              bg="accent.green"
              borderRadius="full"
              initial={{ width: 0 }}
              animate={{ width: `${scale(best)}%` }}
              transition={{ delay: delay + 0.3, type: 'spring' }}
            />
          </Box>
          <Text color="accent.green" fontSize="sm" fontWeight="bold" w="60px" textAlign="right">
            {best.toFixed(1)}
          </Text>
        </HStack>

        <HStack spacing={4}>
          <Text color="gray.500" fontSize="xs" w="50px">Average</Text>
          <Box flex={1} h="8px" bg="bg.tertiary" borderRadius="full" overflow="hidden">
            <MotionBox
              h="100%"
              bg="accent.purple"
              borderRadius="full"
              initial={{ width: 0 }}
              animate={{ width: `${scale(average)}%` }}
              transition={{ delay: delay + 0.4, type: 'spring' }}
            />
          </Box>
          <Text color="accent.purple" fontSize="sm" fontWeight="bold" w="60px" textAlign="right">
            {average.toFixed(1)}
          </Text>
        </HStack>
      </VStack>
    </MotionBox>
  );
}

// Letter statistics grid
function LetterStatsGrid({ letterAccuracy, letterSeekTime, delay = 0 }) {
  const letters = 'abcdefghijklmnopqrstuvwxyz'.split('');

  const getAccuracyColor = (letter) => {
    const stats = letterAccuracy?.[letter];
    if (!stats || stats.presented === 0) return 'gray.600';
    const acc = (stats.correct / stats.presented) * 100;
    if (acc >= 95) return '#00ff88';
    if (acc >= 85) return '#88ff00';
    if (acc >= 75) return '#ffcc00';
    if (acc >= 60) return '#ff8844';
    return '#ff4466';
  };

  const getSeekColor = (letter) => {
    const stats = letterSeekTime?.[letter];
    if (!stats || stats.count === 0) return 'gray.600';
    const avg = stats.total_time_ms / stats.count;
    if (avg <= 150) return '#00ff88';
    if (avg <= 200) return '#88ff00';
    if (avg <= 250) return '#ffcc00';
    if (avg <= 350) return '#ff8844';
    return '#ff4466';
  };

  return (
    <MotionBox
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay }}
    >
      <VStack spacing={4} align="stretch">
        <Text color="gray.400" fontSize="sm" fontWeight="500">Letter Statistics</Text>

        {/* Letter headers */}
        <Flex justify="center" gap={1} flexWrap="wrap">
          {letters.map((letter) => (
            <Text
              key={letter}
              fontSize="xs"
              fontWeight="bold"
              color="gray.500"
              w="20px"
              textAlign="center"
              textTransform="uppercase"
            >
              {letter}
            </Text>
          ))}
        </Flex>

        {/* Accuracy row */}
        <VStack spacing={1}>
          <Text color="gray.500" fontSize="xs">Accuracy</Text>
          <Flex justify="center" gap={1} flexWrap="wrap">
            {letters.map((letter, i) => (
              <MotionBox
                key={letter}
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: delay + i * 0.02 }}
              >
                <Tooltip
                  label={`${letter.toUpperCase()}: ${
                    letterAccuracy?.[letter]
                      ? `${((letterAccuracy[letter].correct / letterAccuracy[letter].presented) * 100).toFixed(0)}%`
                      : 'N/A'
                  }`}
                >
                  <Box
                    w="20px"
                    h="20px"
                    borderRadius="full"
                    bg={getAccuracyColor(letter)}
                    opacity={letterAccuracy?.[letter]?.presented > 0 ? 1 : 0.3}
                  />
                </Tooltip>
              </MotionBox>
            ))}
          </Flex>
        </VStack>

        {/* Seek time row */}
        <VStack spacing={1}>
          <Text color="gray.500" fontSize="xs">Speed</Text>
          <Flex justify="center" gap={1} flexWrap="wrap">
            {letters.map((letter, i) => (
              <MotionBox
                key={letter}
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: delay + 0.5 + i * 0.02 }}
              >
                <Tooltip
                  label={`${letter.toUpperCase()}: ${
                    letterSeekTime?.[letter]
                      ? `${(letterSeekTime[letter].total_time_ms / letterSeekTime[letter].count).toFixed(0)}ms`
                      : 'N/A'
                  }`}
                >
                  <Box
                    w="20px"
                    h="20px"
                    borderRadius="md"
                    bg={getSeekColor(letter)}
                    opacity={letterSeekTime?.[letter]?.count > 0 ? 1 : 0.3}
                  />
                </Tooltip>
              </MotionBox>
            ))}
          </Flex>
        </VStack>
      </VStack>
    </MotionBox>
  );
}

// Finger statistics display
function FingerStats({ fingerStats, delay = 0 }) {
  const fingerLabels = ['LP', 'LR', 'LM', 'LI', '', '', 'RI', 'RM', 'RR', 'RP'];
  const fingerIndices = [0, 1, 2, 3, 6, 7, 8, 9];

  const getAccuracyColor = (finger) => {
    const stats = fingerStats?.[finger];
    if (!stats || stats.presented === 0) return 'gray.600';
    const acc = (stats.correct / stats.presented) * 100;
    if (acc >= 95) return 'accent.green';
    if (acc >= 85) return 'yellow.400';
    if (acc >= 75) return 'orange.400';
    return 'red.400';
  };

  return (
    <MotionBox
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay }}
    >
      <VStack spacing={3}>
        <Text color="gray.400" fontSize="sm" fontWeight="500">Finger Accuracy</Text>
        <HStack spacing={2}>
          {fingerIndices.map((finger, i) => {
            const stats = fingerStats?.[finger];
            const acc = stats && stats.presented > 0
              ? ((stats.correct / stats.presented) * 100).toFixed(0)
              : '-';
            return (
              <MotionBox
                key={finger}
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: delay + i * 0.05, type: 'spring' }}
              >
                <VStack spacing={1}>
                  <Text fontSize="xs" color="gray.500">{fingerLabels[finger]}</Text>
                  <Box
                    w="40px"
                    h="40px"
                    borderRadius="xl"
                    bg="bg.tertiary"
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    border="2px solid"
                    borderColor={getAccuracyColor(finger)}
                  >
                    <Text fontSize="sm" fontWeight="bold" color={getAccuracyColor(finger)}>
                      {acc}
                    </Text>
                  </Box>
                </VStack>
              </MotionBox>
            );
          })}
        </HStack>
      </VStack>
    </MotionBox>
  );
}

// Hand balance display
function HandBalance({ handStats, handAlternations, sameHandRuns, delay = 0 }) {
  const leftStats = handStats?.[0] || { presented: 0, correct: 0 };
  const rightStats = handStats?.[1] || { presented: 0, correct: 0 };
  const total = leftStats.presented + rightStats.presented;
  const leftPct = total > 0 ? (leftStats.presented / total) * 100 : 50;
  const rightPct = 100 - leftPct;

  const totalTransitions = (handAlternations || 0) + (sameHandRuns || 0);
  const alternationRate = totalTransitions > 0
    ? ((handAlternations || 0) / totalTransitions) * 100
    : 0;

  return (
    <MotionBox
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay }}
    >
      <VStack spacing={3}>
        <Text color="gray.400" fontSize="sm" fontWeight="500">Hand Balance</Text>

        <HStack spacing={4} w="100%">
          <Text color="accent.purple" fontWeight="bold">L {leftPct.toFixed(0)}%</Text>
          <Box flex={1} h="12px" bg="bg.tertiary" borderRadius="full" overflow="hidden">
            <Flex h="100%">
              <MotionBox
                h="100%"
                bg="accent.purple"
                initial={{ width: 0 }}
                animate={{ width: `${leftPct}%` }}
                transition={{ delay: delay + 0.2, type: 'spring' }}
              />
              <MotionBox
                h="100%"
                bg="accent.cyan"
                initial={{ width: 0 }}
                animate={{ width: `${rightPct}%` }}
                transition={{ delay: delay + 0.3, type: 'spring' }}
              />
            </Flex>
          </Box>
          <Text color="accent.cyan" fontWeight="bold">R {rightPct.toFixed(0)}%</Text>
        </HStack>

        <HStack spacing={6}>
          <VStack spacing={0}>
            <Text color="gray.500" fontSize="xs">Alternation Rate</Text>
            <Text color="accent.green" fontSize="lg" fontWeight="bold">
              {alternationRate.toFixed(0)}%
            </Text>
          </VStack>
        </HStack>
      </VStack>
    </MotionBox>
  );
}

// Common errors display
function CommonErrors({ errorSubstitution, delay = 0 }) {
  // Flatten and sort errors
  const errors = [];
  if (errorSubstitution) {
    Object.entries(errorSubstitution).forEach(([expected, typed]) => {
      Object.entries(typed).forEach(([typedChar, count]) => {
        errors.push({ expected, typed: typedChar, count });
      });
    });
  }
  errors.sort((a, b) => b.count - a.count);
  const topErrors = errors.slice(0, 5);

  if (topErrors.length === 0) {
    return (
      <MotionBox
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay }}
      >
        <VStack spacing={2}>
          <Text color="gray.400" fontSize="sm" fontWeight="500">Common Errors</Text>
          <Text color="gray.600" fontSize="sm">No errors recorded</Text>
        </VStack>
      </MotionBox>
    );
  }

  return (
    <MotionBox
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay }}
    >
      <VStack spacing={3}>
        <Text color="gray.400" fontSize="sm" fontWeight="500">Common Errors</Text>
        <HStack spacing={2} flexWrap="wrap" justify="center">
          {topErrors.map((error, i) => (
            <MotionBox
              key={`${error.expected}-${error.typed}`}
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: delay + i * 0.1, type: 'spring' }}
            >
              <Box
                px={3}
                py={2}
                bg="rgba(255, 68, 102, 0.1)"
                borderRadius="lg"
                border="1px solid"
                borderColor="accent.red"
              >
                <Text fontSize="sm" color="accent.red">
                  {error.expected}→{error.typed}
                  <Text as="span" color="gray.500" ml={1}>
                    ({error.count})
                  </Text>
                </Text>
              </Box>
            </MotionBox>
          ))}
        </HStack>
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
  const duration = sessionStats?.duration ? sessionStats.duration / 1e9 : 0; // nanoseconds to seconds

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

  return (
    <Box minH="100vh" py={8} px={4} overflowY="auto">
      <Container maxW="container.xl">
        <VStack spacing={8}>
          {/* Header */}
          <MotionBox
            initial={{ opacity: 0, y: -30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ type: 'spring', bounce: 0.4 }}
          >
            <Text
              fontSize={{ base: '4xl', md: '6xl' }}
              fontWeight="800"
              bgGradient="linear(to-r, accent.cyan, accent.green)"
              bgClip="text"
              textAlign="center"
            >
              Round Complete!
            </Text>
          </MotionBox>

          {/* Main stats cards */}
          <SimpleGrid columns={{ base: 1, md: 3 }} spacing={6} w="100%">
            <StatCard
              label="Words Per Minute"
              value={wpm}
              unit="WPM"
              comparison={`Best: ${bestWpm.toFixed(1)} | Avg: ${avgWpm.toFixed(1)}`}
              isBest={isNewBestWpm}
              color="cyan"
              delay={0.1}
            />
            <StatCard
              label="Accuracy"
              value={accuracy}
              unit="%"
              comparison={`Best: ${bestAccuracy.toFixed(1)}% | Avg: ${avgAccuracy.toFixed(1)}%`}
              isBest={isNewBestAccuracy}
              color="green"
              delay={0.2}
            />
            <StatCard
              label="Time"
              value={duration}
              unit="s"
              comparison={`Best: ${bestTime.toFixed(1)}s | Avg: ${avgTime.toFixed(1)}s`}
              isBest={isNewBestTime}
              color="purple"
              delay={0.3}
            />
          </SimpleGrid>

          {/* Comparison bars */}
          <Box
            w="100%"
            bg="bg.card"
            borderRadius="3xl"
            p={8}
            border="1px solid"
            borderColor="whiteAlpha.100"
          >
            <SimpleGrid columns={{ base: 1, md: 3 }} spacing={8}>
              <ComparisonBar
                label="WPM Comparison"
                current={wpm}
                best={bestWpm}
                average={avgWpm}
                delay={0.4}
              />
              <ComparisonBar
                label="Accuracy Comparison"
                current={accuracy}
                best={bestAccuracy}
                average={avgAccuracy}
                delay={0.5}
              />
              <ComparisonBar
                label="Time Comparison"
                current={duration}
                best={bestTime}
                average={avgTime}
                inverted
                delay={0.6}
              />
            </SimpleGrid>
          </Box>

          {/* Letter statistics */}
          <Box
            w="100%"
            bg="bg.card"
            borderRadius="3xl"
            p={8}
            border="1px solid"
            borderColor="whiteAlpha.100"
          >
            <LetterStatsGrid
              letterAccuracy={historicalStats?.letter_accuracy}
              letterSeekTime={historicalStats?.letter_seek_time}
              delay={0.7}
            />
          </Box>

          {/* Typing theory stats */}
          <SimpleGrid columns={{ base: 1, md: 3 }} spacing={6} w="100%">
            <Box
              bg="bg.card"
              borderRadius="3xl"
              p={6}
              border="1px solid"
              borderColor="whiteAlpha.100"
            >
              <FingerStats fingerStats={historicalStats?.finger_stats} delay={0.8} />
            </Box>

            <Box
              bg="bg.card"
              borderRadius="3xl"
              p={6}
              border="1px solid"
              borderColor="whiteAlpha.100"
            >
              <HandBalance
                handStats={historicalStats?.hand_stats}
                handAlternations={historicalStats?.hand_alternations}
                sameHandRuns={historicalStats?.same_hand_runs}
                delay={0.9}
              />
            </Box>

            <Box
              bg="bg.card"
              borderRadius="3xl"
              p={6}
              border="1px solid"
              borderColor="whiteAlpha.100"
            >
              <CommonErrors
                errorSubstitution={historicalStats?.error_substitution}
                delay={1.0}
              />
            </Box>
          </SimpleGrid>

          {/* Session count */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1.1 }}
          >
            <Text color="gray.500" fontSize="lg">
              Total Sessions: <Text as="span" color="accent.cyan" fontWeight="bold">
                {historicalStats?.total_sessions || 0}
              </Text>
            </Text>
          </MotionBox>

          {/* Action buttons */}
          <MotionFlex
            gap={4}
            flexWrap="wrap"
            justify="center"
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 1.2, type: 'spring' }}
          >
            <MotionBox whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
              <Button
                size="chunky"
                variant="glow"
                onClick={onNewRound}
                isLoading={isLoading}
              >
                New Round
              </Button>
            </MotionBox>

            <MotionBox whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
              <Button
                size="chunky"
                variant="ghost"
                onClick={onBackToMenu}
              >
                Back to Menu
              </Button>
            </MotionBox>
          </MotionFlex>
        </VStack>
      </Container>
    </Box>
  );
}

export default ResultsScreen;
