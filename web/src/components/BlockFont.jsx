import React from 'react';
import { Box, Text, HStack, VStack } from '@chakra-ui/react';
import { motion } from 'framer-motion';

const MotionBox = motion(Box);

// Block letters matching the TUI blockfont library
// Each letter is 6 lines tall, using Unicode block characters
const BLOCK_LETTERS = {
  'A': [
    "  ██  ",
    " ████ ",
    "██  ██",
    "██████",
    "██  ██",
    "      ",
  ],
  'B': [
    "█████ ",
    "██  ██",
    "█████ ",
    "██  ██",
    "█████ ",
    "      ",
  ],
  'C': [
    " ████ ",
    "██    ",
    "██    ",
    "██    ",
    " ████ ",
    "      ",
  ],
  'D': [
    "████  ",
    "██ ██ ",
    "██  ██",
    "██ ██ ",
    "████  ",
    "      ",
  ],
  'E': [
    "██████",
    "██    ",
    "████  ",
    "██    ",
    "██████",
    "      ",
  ],
  'F': [
    "██████",
    "██    ",
    "████  ",
    "██    ",
    "██    ",
    "      ",
  ],
  'G': [
    " ████ ",
    "██    ",
    "██ ███",
    "██  ██",
    " ████ ",
    "      ",
  ],
  'H': [
    "██  ██",
    "██  ██",
    "██████",
    "██  ██",
    "██  ██",
    "      ",
  ],
  'I': [
    "██████",
    "  ██  ",
    "  ██  ",
    "  ██  ",
    "██████",
    "      ",
  ],
  'J': [
    "  ████",
    "    ██",
    "    ██",
    "██  ██",
    " ████ ",
    "      ",
  ],
  'K': [
    "██  ██",
    "██ ██ ",
    "████  ",
    "██ ██ ",
    "██  ██",
    "      ",
  ],
  'L': [
    "██    ",
    "██    ",
    "██    ",
    "██    ",
    "██████",
    "      ",
  ],
  'M': [
    "██  ██",
    "██████",
    "██████",
    "██  ██",
    "██  ██",
    "      ",
  ],
  'N': [
    "██  ██",
    "███ ██",
    "██████",
    "██ ███",
    "██  ██",
    "      ",
  ],
  'O': [
    " ████ ",
    "██  ██",
    "██  ██",
    "██  ██",
    " ████ ",
    "      ",
  ],
  'P': [
    "█████ ",
    "██  ██",
    "█████ ",
    "██    ",
    "██    ",
    "      ",
  ],
  'Q': [
    " ████ ",
    "██  ██",
    "██  ██",
    "██ ██ ",
    " ██ ██",
    "      ",
  ],
  'R': [
    "█████ ",
    "██  ██",
    "█████ ",
    "██ ██ ",
    "██  ██",
    "      ",
  ],
  'S': [
    " ████ ",
    "██    ",
    " ████ ",
    "    ██",
    " ████ ",
    "      ",
  ],
  'T': [
    "██████",
    "  ██  ",
    "  ██  ",
    "  ██  ",
    "  ██  ",
    "      ",
  ],
  'U': [
    "██  ██",
    "██  ██",
    "██  ██",
    "██  ██",
    " ████ ",
    "      ",
  ],
  'V': [
    "██  ██",
    "██  ██",
    "██  ██",
    " ████ ",
    "  ██  ",
    "      ",
  ],
  'W': [
    "██  ██",
    "██  ██",
    "██████",
    "██████",
    "██  ██",
    "      ",
  ],
  'X': [
    "██  ██",
    " ████ ",
    "  ██  ",
    " ████ ",
    "██  ██",
    "      ",
  ],
  'Y': [
    "██  ██",
    " ████ ",
    "  ██  ",
    "  ██  ",
    "  ██  ",
    "      ",
  ],
  'Z': [
    "██████",
    "   ██ ",
    "  ██  ",
    " ██   ",
    "██████",
    "      ",
  ],
  ' ': [
    "    ",
    "    ",
    "    ",
    "    ",
    "    ",
    "    ",
  ],
};

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

// Render a single block letter character with color
function BlockChar({ char, color, fontSize = 'xs' }) {
  // Replace spaces with non-breaking spaces for proper rendering
  const displayChar = char === ' ' ? '\u00A0' : char;

  return (
    <Text
      as="span"
      color={color}
      fontFamily="'Fira Code', 'JetBrains Mono', 'Cascadia Code', monospace"
      fontSize={fontSize}
      lineHeight="1"
      whiteSpace="pre"
      letterSpacing="0"
      display="inline"
    >
      {displayChar}
    </Text>
  );
}

// Render a word in block font style
export function BlockFontWord({ word, input = '', showCurrent = true, fontSize = 'xs' }) {
  const upperWord = word.toUpperCase();
  const lines = [];

  // Build 6 lines of block characters
  for (let lineIdx = 0; lineIdx < 6; lineIdx++) {
    const lineChars = [];

    for (let charIdx = 0; charIdx < upperWord.length; charIdx++) {
      const char = upperWord[charIdx];
      const letter = BLOCK_LETTERS[char] || BLOCK_LETTERS[' '];
      const line = letter[lineIdx] || '';

      // Determine color based on typing status
      let color;
      if (charIdx < input.length) {
        // Character has been typed
        const typedCorrect = input[charIdx]?.toLowerCase() === word[charIdx]?.toLowerCase();
        color = typedCorrect ? COLORS.correct : COLORS.incorrect;
      } else if (charIdx === input.length && showCurrent) {
        // Current character to type
        color = COLORS.current;
      } else {
        // Not yet typed
        color = COLORS.untyped;
      }

      // Add each character of the block letter line
      for (const blockChar of line) {
        lineChars.push(
          <BlockChar
            key={`${charIdx}-${lineChars.length}`}
            char={blockChar}
            color={color}
            fontSize={fontSize}
          />
        );
      }

      // Add spacing between letters
      if (charIdx < upperWord.length - 1) {
        lineChars.push(
          <BlockChar
            key={`space-${charIdx}`}
            char=" "
            color={COLORS.untyped}
            fontSize={fontSize}
          />
        );
      }
    }

    lines.push(
      <Box key={lineIdx} whiteSpace="pre" lineHeight="1">
        {lineChars}
      </Box>
    );
  }

  return (
    <VStack spacing={0} align="center">
      {lines}
    </VStack>
  );
}

// Animated block font word with spring physics
export function AnimatedBlockFontWord({ word, input = '', wordKey, showCurrent = true, fontSize = 'xs' }) {
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
        fontFamily="'Fira Code', monospace"
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
        fontFamily="'Fira Code', monospace"
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
        fontFamily="'Fira Code', monospace"
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
        fontFamily="'Fira Code', monospace"
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
        fontFamily="'Fira Code', monospace"
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
        fontFamily="'Fira Code', monospace"
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
    <HStack spacing={2} fontFamily="'Fira Code', monospace" fontSize="sm">
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
      <HStack spacing={0} fontFamily="'Fira Code', monospace" fontSize="md">
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
      <Text color="gray.600" fontSize="xs" fontFamily="'Fira Code', monospace">
        0                        60                       120
      </Text>
    </VStack>
  );
}

export default BlockFontWord;
