import React from 'react';
import { Box, Text, HStack, VStack } from '@chakra-ui/react';
import { motion } from 'framer-motion';

const MotionBox = motion(Box);

// TUI-style colors matching the terminal app
const COLORS = {
  correct: '#00ff00',   // Bright green (ANSI 10)
  incorrect: '#ff0000', // Bright red (ANSI 9)
  untyped: '#808080',   // Gray (ANSI 8)
  current: '#D4922A',   // Kartoza orange for current letter
};

// Gradient colors from red to green (matches TUI gradient)
export const GRADIENT_COLORS = [
  '#ff0000', // 196 - red
  '#ff5f00', // 202
  '#ff8700', // 208
  '#ffaf00', // 214
  '#ffd700', // 220
  '#ffff00', // 226 - yellow
  '#d7ff00', // 190
  '#afff00', // 154
  '#87ff00', // 118
  '#5fff00', // 82
  '#00ff00', // 46
  '#00ff5f', // 47 - bright green
];

// Get gradient color based on position (0-1)
export function getGradientColor(position) {
  const idx = Math.min(Math.floor(position * (GRADIENT_COLORS.length - 1)), GRADIENT_COLORS.length - 1);
  return GRADIENT_COLORS[Math.max(0, idx)];
}

// Font family for monospace elements
const BLOCK_FONT_FAMILY = "'Fira Code', 'JetBrains Mono', 'Cascadia Code', monospace";

// Render a word as large bold uppercase letters with per-character coloring
export function BlockFontWord({ word, input = '', showCurrent = true, fontSize = '6xl' }) {
  const upperWord = word.toUpperCase();

  return (
    <HStack spacing={0} justify="center" userSelect="none">
      {[...upperWord].map((char, charIdx) => {
        let color;
        if (charIdx < input.length) {
          const typedCorrect = input[charIdx]?.toLowerCase() === word[charIdx]?.toLowerCase();
          color = typedCorrect ? COLORS.correct : COLORS.incorrect;
        } else if (charIdx === input.length && showCurrent) {
          color = COLORS.current;
        } else {
          color = COLORS.untyped;
        }

        return (
          <Text
            key={charIdx}
            as="span"
            color={color}
            fontFamily={BLOCK_FONT_FAMILY}
            fontSize={fontSize}
            fontWeight="900"
            lineHeight="1.1"
            letterSpacing="0.05em"
          >
            {char}
          </Text>
        );
      })}
    </HStack>
  );
}

// Animated word with spring physics
export function AnimatedBlockFontWord({ word, input = '', wordKey, showCurrent = true, fontSize = '6xl' }) {
  return (
    <MotionBox
      key={wordKey}
      initial={{ y: 60, opacity: 0, scale: 0.8 }}
      animate={{ y: 0, opacity: 1, scale: 1 }}
      exit={{ y: -80, opacity: 0, scale: 0.6 }}
      transition={{
        type: 'spring',
        stiffness: 300,
        damping: 25,
        mass: 0.8,
      }}
    >
      <BlockFontWord word={word} input={input} showCurrent={showCurrent} fontSize={fontSize} />
    </MotionBox>
  );
}

// TUI-style gradient progress bar
export function GradientBar({ value, maxValue, width = 30, showStar = false }) {
  const fillPercent = Math.min(Math.max(value / maxValue, 0), 1);
  const filledWidth = Math.floor(width * fillPercent);
  const emptyWidth = width - filledWidth;

  const chars = [];

  // Filled portion with gradient
  for (let i = 0; i < filledWidth; i++) {
    const position = i / width;
    const color = getGradientColor(position);
    chars.push(
      <Text
        key={`fill-${i}`}
        as="span"
        color={color}
        fontFamily={BLOCK_FONT_FAMILY}
        fontSize="md"
      >
        █
      </Text>
    );
  }

  // Empty portion
  for (let i = 0; i < emptyWidth; i++) {
    chars.push(
      <Text
        key={`empty-${i}`}
        as="span"
        color="gray.700"
        fontFamily={BLOCK_FONT_FAMILY}
        fontSize="md"
      >
        ░
      </Text>
    );
  }

  // Star for new best
  if (showStar) {
    chars.push(
      <Text
        key="star"
        as="span"
        color="yellow.400"
        fontFamily={BLOCK_FONT_FAMILY}
        fontSize="md"
        fontWeight="bold"
        ml={1}
      >
        *
      </Text>
    );
  }

  return (
    <Box display="inline-flex" alignItems="center">
      {chars}
    </Box>
  );
}

// TUI-style time bar (inverted - lower is better)
export function TimeBar({ value, maxValue, width = 30, showStar = false }) {
  // For time, lower is better, so we invert the fill
  const fillPercent = Math.min(Math.max(1 - (value / maxValue), 0), 1);
  const filledWidth = Math.floor(width * fillPercent);
  const emptyWidth = width - filledWidth;

  const chars = [];

  // Filled portion with gradient
  for (let i = 0; i < filledWidth; i++) {
    const position = i / width;
    const color = getGradientColor(position);
    chars.push(
      <Text
        key={`fill-${i}`}
        as="span"
        color={color}
        fontFamily={BLOCK_FONT_FAMILY}
        fontSize="md"
      >
        █
      </Text>
    );
  }

  // Empty portion
  for (let i = 0; i < emptyWidth; i++) {
    chars.push(
      <Text
        key={`empty-${i}`}
        as="span"
        color="gray.700"
        fontFamily={BLOCK_FONT_FAMILY}
        fontSize="md"
      >
        ░
      </Text>
    );
  }

  // Star for new best
  if (showStar) {
    chars.push(
      <Text
        key="star"
        as="span"
        color="yellow.400"
        fontFamily={BLOCK_FONT_FAMILY}
        fontSize="md"
        fontWeight="bold"
        ml={1}
      >
        *
      </Text>
    );
  }

  return (
    <Box display="inline-flex" alignItems="center">
      {chars}
    </Box>
  );
}

// TUI-style stat row with label, value, and gradient bar
export function StatRow({ label, value, bar, labelWidth = 18, valueWidth = 8 }) {
  return (
    <HStack spacing={2} fontFamily={BLOCK_FONT_FAMILY} fontSize="sm">
      <Text
        color="gray.400"
        w={`${labelWidth}ch`}
        textAlign="right"
        whiteSpace="nowrap"
      >
        {label}
      </Text>
      <Text
        color="white"
        w={`${valueWidth}ch`}
        textAlign="right"
        whiteSpace="nowrap"
      >
        {value}
      </Text>
      {bar}
    </HStack>
  );
}

// TUI-style WPM bar for typing screen
export function WpmBar({ wpm, maxWpm = 120 }) {
  const percentage = Math.min(Math.max(wpm / maxWpm, 0), 1);
  const barWidth = 50;
  const filledWidth = Math.floor(barWidth * percentage);
  const emptyWidth = barWidth - filledWidth;

  // Determine WPM color
  let wpmColor;
  if (wpm >= 60) {
    wpmColor = '#00ff00'; // Bright green
  } else if (wpm >= 40) {
    wpmColor = '#ffff00'; // Yellow
  } else {
    wpmColor = '#ff0000'; // Red
  }

  return (
    <VStack spacing={1}>
      <HStack spacing={0} fontFamily={BLOCK_FONT_FAMILY} fontSize="md">
        {/* Filled portion with gradient */}
        {Array.from({ length: filledWidth }).map((_, i) => {
          const position = i / barWidth;
          const color = getGradientColor(position);
          return (
            <Text key={`fill-${i}`} as="span" color={color}>█</Text>
          );
        })}
        {/* Empty portion */}
        {Array.from({ length: emptyWidth }).map((_, i) => (
          <Text key={`empty-${i}`} as="span" color="gray.700">░</Text>
        ))}
        {/* WPM label */}
        <Text
          as="span"
          color={wpmColor}
          fontWeight="bold"
          ml={2}
        >
          {Math.round(wpm)} WPM
        </Text>
      </HStack>
      {/* Scale */}
      <Text color="gray.600" fontSize="xs" fontFamily={BLOCK_FONT_FAMILY}>
        0                        60                       120
      </Text>
    </VStack>
  );
}

export default BlockFontWord;
