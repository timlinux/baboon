import React, { useEffect, useState } from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Button,
  Flex,
  Link,
  useBreakpointValue,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';
import { GRADIENT_COLORS } from './BlockFont.jsx';
import AdSense from './AdSense.jsx';
import BrandingFooter from './BrandingFooter.jsx';

const MotionBox = motion(Box);
const MotionFlex = motion(Flex);

// Get accuracy color based on value
function getAccuracyColor(accuracy) {
  if (accuracy >= 95) return GRADIENT_COLORS[10];
  if (accuracy >= 90) return GRADIENT_COLORS[9];
  if (accuracy >= 85) return GRADIENT_COLORS[8];
  if (accuracy >= 80) return GRADIENT_COLORS[7];
  if (accuracy >= 75) return GRADIENT_COLORS[6];
  if (accuracy >= 70) return GRADIENT_COLORS[5];
  if (accuracy >= 65) return GRADIENT_COLORS[4];
  if (accuracy >= 60) return GRADIENT_COLORS[3];
  if (accuracy >= 50) return GRADIENT_COLORS[2];
  if (accuracy >= 40) return GRADIENT_COLORS[1];
  return GRADIENT_COLORS[0];
}

// Get relative color for values within a range
function getRelativeColor(value, minVal, maxVal) {
  if (maxVal <= minVal) return GRADIENT_COLORS[10];
  const normalized = (value - minVal) / (maxVal - minVal);
  const idx = Math.min(Math.floor(normalized * 11), 11);
  return GRADIENT_COLORS[Math.max(0, idx)];
}

// Get inverted relative color (for seek time - lower is better)
function getInvertedRelativeColor(value, minVal, maxVal) {
  if (maxVal <= minVal) return GRADIENT_COLORS[10];
  // Invert: lower values get higher (greener) colors
  const invertedValue = maxVal - value + minVal;
  return getRelativeColor(invertedValue, minVal, maxVal);
}

// Get frequency color (high frequency = red = needs more practice)
function getFrequencyColor(frequency) {
  // Invert: higher frequency = red (more exposure = lower index)
  const idx = Math.min(Math.floor((1 - frequency) * 11), 11);
  return GRADIENT_COLORS[Math.max(0, idx)];
}

// Get gradient color at position
function getGradientColor(position) {
  const idx = Math.min(Math.floor(position * (GRADIENT_COLORS.length - 1)), GRADIENT_COLORS.length - 1);
  return GRADIENT_COLORS[Math.max(0, idx)];
}

// Responsive gradient bar
function ResponsiveGradientBar({ value, maxValue, barWidth, showStar = false, inverted = false }) {
  const fillPercent = inverted
    ? Math.min(Math.max(1 - (value / maxValue), 0), 1)
    : Math.min(Math.max(value / maxValue, 0), 1);
  const filledWidth = Math.floor(barWidth * fillPercent);
  const emptyWidth = barWidth - filledWidth;

  return (
    <HStack spacing={0} display="inline-flex">
      {Array.from({ length: filledWidth }).map((_, i) => (
        <Text key={`fill-${i}`} as="span" color={getGradientColor(i / barWidth)}>█</Text>
      ))}
      {Array.from({ length: emptyWidth }).map((_, i) => (
        <Text key={`empty-${i}`} as="span" color="gray.700">░</Text>
      ))}
      {showStar && (
        <Text as="span" color="yellow.400" fontWeight="bold" ml={1}>*</Text>
      )}
    </HStack>
  );
}

// Animated stat row wrapper
function AnimatedStatRow({ children, delay = 0 }) {
  return (
    <MotionBox
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay, type: 'spring', stiffness: 200, damping: 20 }}
      w="100%"
    >
      {children}
    </MotionBox>
  );
}

function ResultsScreen({
  sessionStats,
  historicalStats,
  onNewRound,
  onBackToMenu,
  onShowLeaderboard,
  isLoading,
  adsenseEnabled,
  adsenseKey,
  practiceMode,
  versionInfo,
}) {
  const wpm = sessionStats?.wpm || 0;
  const accuracy = sessionStats?.accuracy || 0;
  const duration = sessionStats?.duration ? sessionStats.duration / 1e9 : 0;

  // Use n-gram specific stats when in ngram mode, fall back to regular stats
  const bestWpm = practiceMode === 'ngrams'
    ? (historicalStats?.ngram_best_wpm || historicalStats?.best_wpm || 0)
    : (historicalStats?.best_wpm || 0);
  const bestAccuracy = practiceMode === 'ngrams'
    ? (historicalStats?.ngram_best_accuracy || historicalStats?.best_accuracy || 0)
    : (historicalStats?.best_accuracy || 0);
  const bestTime = practiceMode === 'ngrams'
    ? (historicalStats?.ngram_best_time || historicalStats?.best_time || 0)
    : (historicalStats?.best_time || 0);

  // Use n-gram specific totals for averages when in ngram mode
  const totalSessions = practiceMode === 'ngrams'
    ? (historicalStats?.ngram_total_sessions || historicalStats?.total_sessions || 0)
    : (historicalStats?.total_sessions || 0);
  const avgWpm = totalSessions > 0
    ? (practiceMode === 'ngrams'
      ? (historicalStats?.ngram_total_wpm || historicalStats?.total_wpm || 0) / totalSessions
      : (historicalStats?.total_wpm || 0) / totalSessions)
    : 0;
  const avgAccuracy = totalSessions > 0
    ? (practiceMode === 'ngrams'
      ? (historicalStats?.ngram_total_accuracy || historicalStats?.total_accuracy || 0) / totalSessions
      : (historicalStats?.total_accuracy || 0) / totalSessions)
    : 0;
  const avgTime = totalSessions > 0
    ? (practiceMode === 'ngrams'
      ? (historicalStats?.ngram_total_time || historicalStats?.total_time || 0) / totalSessions
      : (historicalStats?.total_time || 0) / totalSessions)
    : 0;

  const isNewBestWpm = wpm >= bestWpm && totalSessions > 0;
  const isNewBestAccuracy = accuracy >= bestAccuracy && totalSessions > 0;
  const isNewBestTime = (bestTime === 0 || duration <= bestTime) && totalSessions > 0;

  // Responsive sizing
  const fontSize = useBreakpointValue({ base: 'md', md: 'lg', lg: 'xl', xl: '2xl' });
  const smallFontSize = useBreakpointValue({ base: 'sm', md: 'md', lg: 'lg', xl: 'xl' });
  const titleSize = useBreakpointValue({ base: '2xl', md: '4xl', lg: '5xl', xl: '6xl' });
  const barWidth = useBreakpointValue({ base: 20, md: 30, lg: 40, xl: 50 });
  const labelWidth = useBreakpointValue({ base: '16ch', md: '18ch', lg: '20ch', xl: '22ch' });
  const valueWidth = useBreakpointValue({ base: '8ch', md: '10ch', lg: '12ch' });

  // Handle keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === 'Enter' || e.key === 'Tab') {
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

  const maxWpmDisplay = 120;
  const maxTimeDisplay = 180;
  const maxAccuracyDisplay = 100;

  let animIdx = 0;

  // Stat row component
  const StatRow = ({ label, value, bar }) => (
    <HStack spacing={2} fontFamily="'Fira Code', monospace" fontSize={fontSize} justify="center" w="100%">
      <Text color="gray.400" w={labelWidth} textAlign="right" whiteSpace="nowrap">
        {label}
      </Text>
      <Text color="white" w={valueWidth} textAlign="right" whiteSpace="nowrap">
        {value}
      </Text>
      {bar}
    </HStack>
  );

  // Finger stats
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

  let minFingerAcc = 100, maxFingerAcc = 0;
  const fingerAccuracies = fingers.map(f => {
    const stats = historicalStats?.finger_stats?.[f.id];
    if (!stats || stats.presented === 0) return -1;
    const acc = (stats.correct / stats.presented) * 100;
    if (acc < minFingerAcc) minFingerAcc = acc;
    if (acc > maxFingerAcc) maxFingerAcc = acc;
    return acc;
  });

  // Hand stats
  const leftStats = historicalStats?.hand_stats?.[0] || { correct: 0 };
  const rightStats = historicalStats?.hand_stats?.[1] || { correct: 0 };
  const totalHand = leftStats.correct + rightStats.correct;
  const leftPct = totalHand > 0 ? (leftStats.correct / totalHand) * 100 : 50;
  const rightPct = 100 - leftPct;
  const totalTransitions = (historicalStats?.hand_alternations || 0) + (historicalStats?.same_hand_runs || 0);
  const altRate = totalTransitions > 0 ? ((historicalStats?.hand_alternations || 0) / totalTransitions) * 100 : 0;

  // Common errors
  const errors = [];
  if (historicalStats?.error_substitution) {
    Object.entries(historicalStats.error_substitution).forEach(([expected, typed]) => {
      Object.entries(typed).forEach(([typedChar, count]) => {
        errors.push({ expected, typed: typedChar, count });
      });
    });
  }
  errors.sort((a, b) => b.count - a.count);
  const topErrors = errors.slice(0, 5);

  // Per-letter stats
  const letters = 'abcdefghijklmnopqrstuvwxyz'.split('');

  // Calculate letter accuracies with min/max for relative coloring
  const letterAccuracies = letters.map(letter => {
    const stats = historicalStats?.letter_accuracy?.[letter];
    if (!stats || stats.presented === 0) return null;
    return (stats.correct / stats.presented) * 100;
  });
  const validAccuracies = letterAccuracies.filter(a => a !== null);
  const minLetterAcc = validAccuracies.length > 0 ? Math.min(...validAccuracies) : 0;
  const maxLetterAcc = validAccuracies.length > 0 ? Math.max(...validAccuracies) : 100;

  // Calculate letter frequencies with max for relative coloring
  const letterFrequencies = letters.map(letter => {
    const stats = historicalStats?.letter_accuracy?.[letter];
    return stats?.presented || 0;
  });
  const maxFrequency = Math.max(...letterFrequencies, 1);

  // Calculate letter seek times with min/max for relative coloring
  const letterSeekTimes = letters.map(letter => {
    const stats = historicalStats?.letter_seek_time?.[letter];
    if (!stats || stats.count === 0) return null;
    return stats.total_time_ms / stats.count;
  });
  const validSeekTimes = letterSeekTimes.filter(t => t !== null);
  const minSeekTime = validSeekTimes.length > 0 ? Math.min(...validSeekTimes) : 0;
  const maxSeekTime = validSeekTimes.length > 0 ? Math.max(...validSeekTimes) : 1000;

  return (
    <Flex minH="100vh" direction="column" p={{ base: 2, md: 4 }}>
      {/* Header */}
      <Flex justify="center" py={2}>
        <Text
          color="cyan.400"
          fontSize={smallFontSize}
          fontWeight="bold"
          fontFamily="'Fira Code', monospace"
        >
          {practiceMode === 'ngrams' ? 'BABOON - N-gram Training' : 'BABOON - Typing Practice'}
        </Text>
      </Flex>

      {/* Title */}
      <MotionBox
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ type: 'spring', bounce: 0.4 }}
        textAlign="center"
        py={{ base: 2, md: 4 }}
      >
        <Text
          fontSize={titleSize}
          fontWeight="bold"
          color="cyan.400"
          fontFamily="'Fira Code', monospace"
        >
          {practiceMode === 'ngrams' ? 'N-gram Training Complete!' : 'Round Complete!'}
        </Text>
      </MotionBox>

      {/* Stats */}
      <VStack spacing={{ base: 1, md: 2 }} align="center" flex={1} overflow="auto" pb={4}>
        {/* WPM Section */}
        <Box h={{ base: 1, md: 2 }} />
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <StatRow
            label="WPM this run:"
            value={wpm.toFixed(1)}
            bar={<ResponsiveGradientBar value={wpm} maxValue={maxWpmDisplay} barWidth={barWidth} showStar={isNewBestWpm} />}
          />
        </AnimatedStatRow>
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <StatRow
            label="WPM best:"
            value={bestWpm.toFixed(1)}
            bar={<ResponsiveGradientBar value={bestWpm} maxValue={maxWpmDisplay} barWidth={barWidth} />}
          />
        </AnimatedStatRow>
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <StatRow
            label="WPM average:"
            value={avgWpm.toFixed(1)}
            bar={<ResponsiveGradientBar value={avgWpm} maxValue={maxWpmDisplay} barWidth={barWidth} />}
          />
        </AnimatedStatRow>

        {/* Time Section */}
        <Box h={{ base: 1, md: 2 }} />
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <StatRow
            label="Time this run:"
            value={`${duration.toFixed(1)}s`}
            bar={<ResponsiveGradientBar value={duration} maxValue={maxTimeDisplay} barWidth={barWidth} showStar={isNewBestTime} inverted />}
          />
        </AnimatedStatRow>
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <StatRow
            label="Time best:"
            value={`${bestTime.toFixed(1)}s`}
            bar={<ResponsiveGradientBar value={bestTime} maxValue={maxTimeDisplay} barWidth={barWidth} inverted />}
          />
        </AnimatedStatRow>
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <StatRow
            label="Time average:"
            value={`${avgTime.toFixed(1)}s`}
            bar={<ResponsiveGradientBar value={avgTime} maxValue={maxTimeDisplay} barWidth={barWidth} inverted />}
          />
        </AnimatedStatRow>

        {/* Accuracy Section */}
        <Box h={{ base: 1, md: 2 }} />
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <StatRow
            label="Accuracy this run:"
            value={`${accuracy.toFixed(1)}%`}
            bar={<ResponsiveGradientBar value={accuracy} maxValue={maxAccuracyDisplay} barWidth={barWidth} showStar={isNewBestAccuracy} />}
          />
        </AnimatedStatRow>
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <StatRow
            label="Accuracy best:"
            value={`${bestAccuracy.toFixed(1)}%`}
            bar={<ResponsiveGradientBar value={bestAccuracy} maxValue={maxAccuracyDisplay} barWidth={barWidth} />}
          />
        </AnimatedStatRow>
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <StatRow
            label="Accuracy average:"
            value={`${avgAccuracy.toFixed(1)}%`}
            bar={<ResponsiveGradientBar value={avgAccuracy} maxValue={maxAccuracyDisplay} barWidth={barWidth} />}
          />
        </AnimatedStatRow>

        {/* Sessions */}
        <Box h={{ base: 1, md: 2 }} />
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <HStack spacing={2} fontFamily="'Fira Code', monospace" fontSize={fontSize} justify="center">
            <Text color="cyan.400" w={labelWidth} textAlign="right">
              {practiceMode === 'ngrams' ? 'N-gram sessions:' : 'Total sessions:'}
            </Text>
            <Text color="white">{totalSessions || 1}</Text>
          </HStack>
        </AnimatedStatRow>

        {/* Per-Letter Stats Matrix */}
        <Box h={{ base: 2, md: 4 }} />
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <Text
            color="cyan.400"
            fontFamily="'Fira Code', monospace"
            fontSize={smallFontSize}
            fontWeight="bold"
            textAlign="center"
          >
            Per-Key Statistics
          </Text>
        </AnimatedStatRow>

        {/* Letter Stats Table - using HTML table for proper alignment */}
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <Box
            as="table"
            fontFamily="'Fira Code', monospace"
            fontSize={{ base: 'xs', md: 'sm', lg: 'md' }}
            mx="auto"
            sx={{ borderCollapse: 'collapse' }}
          >
            <Box as="tbody">
              {/* Header row */}
              <Box as="tr">
                <Box as="td" w={{ base: '40px', md: '50px' }} />
                {letters.map(letter => (
                  <Box
                    as="td"
                    key={`hdr-${letter}`}
                    color="white"
                    fontWeight="bold"
                    textAlign="center"
                    w={{ base: '14px', md: '18px', lg: '22px' }}
                    px="1px"
                  >
                    {letter.toUpperCase()}
                  </Box>
                ))}
              </Box>

              {/* Accuracy row */}
              <Box as="tr">
                <Box as="td" color="gray.500" textAlign="right" pr={1}>Acc</Box>
                {letters.map((letter, i) => {
                  const acc = letterAccuracies[i];
                  const color = acc === null ? 'gray.700' : getRelativeColor(acc, minLetterAcc, maxLetterAcc);
                  return (
                    <Box as="td" key={`acc-${letter}`} color={color} textAlign="center">●</Box>
                  );
                })}
              </Box>

              {/* Frequency row */}
              <Box as="tr">
                <Box as="td" color="gray.500" textAlign="right" pr={1}>Freq</Box>
                {letters.map((letter, i) => {
                  const freq = letterFrequencies[i];
                  const normalizedFreq = freq / maxFrequency;
                  const color = freq === 0 ? 'gray.700' : getFrequencyColor(normalizedFreq);
                  return (
                    <Box as="td" key={`freq-${letter}`} color={color} textAlign="center">●</Box>
                  );
                })}
              </Box>

              {/* Seek Time row */}
              <Box as="tr">
                <Box as="td" color="gray.500" textAlign="right" pr={1}>Time</Box>
                {letters.map((letter, i) => {
                  const seekTime = letterSeekTimes[i];
                  const color = seekTime === null ? 'gray.700' : getInvertedRelativeColor(seekTime, minSeekTime, maxSeekTime);
                  return (
                    <Box as="td" key={`time-${letter}`} color={color} textAlign="center">●</Box>
                  );
                })}
              </Box>
            </Box>
          </Box>
        </AnimatedStatRow>

        {/* Legend for letter stats */}
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <HStack spacing={4} fontFamily="'Fira Code', monospace" fontSize={{ base: '2xs', md: 'xs' }} justify="center" color="gray.600">
            <Text>Acc: accuracy (green=good)</Text>
            <Text>Freq: exposure (red=high)</Text>
            <Text>Time: speed (green=fast)</Text>
          </HStack>
        </AnimatedStatRow>

        {/* Finger Stats */}
        <Box h={{ base: 1, md: 2 }} />
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <HStack spacing={2} fontFamily="'Fira Code', monospace" fontSize={smallFontSize} justify="center" flexWrap="wrap">
            <Text color="gray.400" w={labelWidth} textAlign="right">Finger accuracy:</Text>
            <HStack spacing={{ base: 1, md: 2 }}>
              {fingers.map((f, i) => {
                const acc = fingerAccuracies[i];
                const color = acc < 0 ? 'gray.600' : getRelativeColor(acc, minFingerAcc, maxFingerAcc);
                return (
                  <HStack key={f.id} spacing={0}>
                    <Text color="gray.500" fontSize={smallFontSize}>{f.label}</Text>
                    <Text color={color} fontSize={fontSize}>●</Text>
                  </HStack>
                );
              })}
            </HStack>
          </HStack>
        </AnimatedStatRow>

        {/* Hand Stats */}
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <HStack spacing={2} fontFamily="'Fira Code', monospace" fontSize={fontSize} justify="center">
            <Text color="gray.400" w={labelWidth} textAlign="right">Hand balance:</Text>
            <Text color="white">L:{leftPct.toFixed(0)}% R:{rightPct.toFixed(0)}%</Text>
            <Text color="gray.400">Alt:</Text>
            <Text color={getAccuracyColor(altRate)}>{altRate.toFixed(0)}%</Text>
          </HStack>
        </AnimatedStatRow>

        {/* Common Errors */}
        <AnimatedStatRow delay={animIdx++ * 0.03}>
          <HStack spacing={2} fontFamily="'Fira Code', monospace" fontSize={smallFontSize} justify="center" flexWrap="wrap">
            <Text color="gray.400" w={labelWidth} textAlign="right">Common errors:</Text>
            {topErrors.length === 0 ? (
              <Text color="gray.600">none</Text>
            ) : (
              topErrors.map((e) => (
                <HStack key={`${e.expected}-${e.typed}`} spacing={0}>
                  <Text color="red.400">{e.expected}→{e.typed}</Text>
                  <Text color="gray.600">({e.count})</Text>
                </HStack>
              ))
            )}
          </HStack>
        </AnimatedStatRow>

        {/* New best legend */}
        {(isNewBestWpm || isNewBestTime || isNewBestAccuracy) && (
          <AnimatedStatRow delay={animIdx++ * 0.03}>
            <Text color="yellow.400" fontFamily="'Fira Code', monospace" fontSize={fontSize} fontWeight="bold" mt={2} textAlign="center">
              * = New personal best!
            </Text>
          </AnimatedStatRow>
        )}
      </VStack>

      {/* Action Buttons */}
      <MotionFlex
        gap={4}
        justify="center"
        py={4}
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4, type: 'spring' }}
      >
        <MotionBox whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
          <Button
            size={{ base: 'md', md: 'lg' }}
            variant="outline"
            onClick={onNewRound}
            isLoading={isLoading}
            px={{ base: 6, md: 8 }}
            fontFamily="'Fira Code', monospace"
            fontSize={fontSize}
            borderColor="#00ff00"
            color="#00ff00"
            _hover={{
              bg: 'rgba(0, 255, 0, 0.1)',
              boxShadow: '0 0 20px rgba(0, 255, 0, 0.3)',
            }}
          >
            NEW ROUND
          </Button>
        </MotionBox>

        {onShowLeaderboard && (
          <MotionBox whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              size={{ base: 'md', md: 'lg' }}
              variant="outline"
              onClick={onShowLeaderboard}
              px={{ base: 6, md: 8 }}
              fontFamily="'Fira Code', monospace"
              fontSize={smallFontSize}
              borderColor="yellow.500"
              color="yellow.400"
              _hover={{
                bg: 'rgba(255, 215, 0, 0.1)',
                boxShadow: '0 0 20px rgba(255, 215, 0, 0.3)',
              }}
            >
              LEADERBOARD
            </Button>
          </MotionBox>
        )}

        <MotionBox whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
          <Button
            size={{ base: 'md', md: 'lg' }}
            variant="outline"
            onClick={onBackToMenu}
            px={{ base: 6, md: 8 }}
            fontFamily="'Fira Code', monospace"
            fontSize={smallFontSize}
            borderColor="gray.500"
            color="gray.400"
            _hover={{
              bg: 'rgba(128, 128, 128, 0.1)',
              boxShadow: '0 0 20px rgba(128, 128, 128, 0.3)',
              color: 'white',
            }}
          >
            MENU
          </Button>
        </MotionBox>
      </MotionFlex>

      {/* Keyboard hint */}
      <Text color="gray.600" fontSize={smallFontSize} textAlign="center" pb={2} fontFamily="'Fira Code', monospace">
        TAB or ENTER = new round | ESC = quit
      </Text>

      {/* Kartoza branding */}
      <Flex justify="center" pb={4}>
        <BrandingFooter versionInfo={versionInfo} fontSize={smallFontSize} />
      </Flex>

      {/* AdSense ad */}
      <Box pb={4}>
        <AdSense publisherId={adsenseKey} showPreview={true} />
      </Box>
    </Flex>
  );
}

export default ResultsScreen;
