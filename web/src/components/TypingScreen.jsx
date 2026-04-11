import React, { useEffect, useState, useRef } from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Flex,
  Progress,
  Link,
  useBreakpointValue,
} from '@chakra-ui/react';
import { motion, AnimatePresence } from 'framer-motion';
import AdSense from './AdSense.jsx';
import { AnimatedBlockFontWord, WpmBar } from './BlockFont.jsx';

const MotionBox = motion(Box);

// Calculate optimal font size based on word length and viewport
function useBlockFontSize(wordLength) {
  // Each letter is ~7 chars wide (6 + 1 spacing), estimate total width
  const charsPerLetter = 7;
  const estimatedWidth = wordLength * charsPerLetter;

  // Use viewport-based sizing
  const baseSizes = useBreakpointValue({
    base: { short: 'xs', medium: '2xs', long: '3xs' },
    sm: { short: 'sm', medium: 'xs', long: '2xs' },
    md: { short: 'md', medium: 'sm', long: 'xs' },
    lg: { short: 'lg', medium: 'md', long: 'sm' },
    xl: { short: 'xl', medium: 'lg', long: 'md' },
    '2xl': { short: '2xl', medium: 'xl', long: 'lg' },
  }) || { short: 'sm', medium: 'xs', long: '2xs' };

  if (wordLength <= 4) return baseSizes.short;
  if (wordLength <= 7) return baseSizes.medium;
  return baseSizes.long;
}

// Preview word (smaller, faded) for carousel
function PreviewWord({ word, position, fontSize }) {
  const isPrevious = position === 'previous';

  const previewSize = useBreakpointValue({
    base: 'lg',
    md: '2xl',
    lg: '3xl',
    xl: '4xl',
  });

  return (
    <MotionBox
      position="absolute"
      top={isPrevious ? '8%' : undefined}
      bottom={!isPrevious ? '8%' : undefined}
      initial={{ opacity: 0, y: isPrevious ? 30 : -30 }}
      animate={{ opacity: 0.5, y: 0 }}
      exit={{ opacity: 0, y: isPrevious ? -30 : 30 }}
      transition={{
        type: 'spring',
        stiffness: 200,
        damping: 25,
      }}
    >
      <Text
        fontFamily="'Fira Code', 'JetBrains Mono', monospace"
        fontSize={previewSize}
        color="gray.500"
        textTransform="lowercase"
        letterSpacing="0.2em"
      >
        {isPrevious ? `· · · ${word} · · ·` : `▼  ${word}  ▼`}
      </Text>
    </MotionBox>
  );
}

// Extra typed characters (errors beyond word length)
function ExtraChars({ chars }) {
  const fontSize = useBreakpointValue({ base: 'lg', md: 'xl', lg: '2xl' });

  return (
    <MotionBox
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -10 }}
    >
      <HStack spacing={1} mt={4}>
        {chars.split('').map((char, index) => (
          <MotionBox
            key={`extra-${index}`}
            initial={{ scale: 0, rotate: -10 }}
            animate={{ scale: 1, rotate: [0, -5, 5, 0] }}
            transition={{ type: 'spring', stiffness: 500, damping: 20 }}
          >
            <Box
              px={3}
              py={1}
              bg="rgba(255, 0, 0, 0.2)"
              borderRadius="lg"
              border="2px solid"
              borderColor="#ff0000"
            >
              <Text
                color="#ff0000"
                fontFamily="'Fira Code', monospace"
                fontSize={fontSize}
              >
                {char}
              </Text>
            </Box>
          </MotionBox>
        ))}
      </HStack>
    </MotionBox>
  );
}

function TypingScreen({
  gameState,
  liveWpm,
  timerStarted,
  onKeystroke,
  onBackspace,
  onSpace,
  onExit,
  adsenseEnabled,
  adsenseKey,
}) {
  const currentWord = gameState?.current_word || '';
  const currentInput = gameState?.current_input || '';
  const previousWord = gameState?.previous_word || '';
  const nextWord = gameState?.next_word || '';
  const wordNumber = gameState?.word_number || 1;
  const totalWords = gameState?.total_words || 30;

  // Track word changes for animation keys
  const [wordKey, setWordKey] = useState(0);
  const prevWordRef = useRef(currentWord);

  useEffect(() => {
    if (currentWord !== prevWordRef.current) {
      setWordKey(prev => prev + 1);
      prevWordRef.current = currentWord;
    }
  }, [currentWord]);

  // Calculate responsive font size based on word length
  const blockFontSize = useBlockFontSize(currentWord.length);

  // Responsive sizes for other elements
  const headerSize = useBreakpointValue({ base: 'md', md: 'lg', lg: 'xl' });
  const counterSize = useBreakpointValue({ base: 'lg', md: 'xl', lg: '2xl' });

  // Handle keyboard input
  useEffect(() => {
    const handleKeyDown = (e) => {
      // Ignore if modifier keys are pressed (except shift)
      if (e.ctrlKey || e.metaKey || e.altKey) return;

      if (e.key === 'Escape') {
        onExit();
        return;
      }

      if (e.key === 'Backspace') {
        e.preventDefault();
        onBackspace();
        return;
      }

      if (e.key === ' ') {
        e.preventDefault();
        onSpace();
        return;
      }

      // Only accept printable characters
      if (e.key.length === 1) {
        e.preventDefault();
        onKeystroke(e.key);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [onKeystroke, onBackspace, onSpace, onExit]);

  const extraChars = currentInput.length > currentWord.length
    ? currentInput.slice(currentWord.length)
    : '';

  return (
    <Flex minH="100vh" direction="column" p={{ base: 2, md: 4 }} overflow="hidden" role="main" aria-label="Typing practice">
      {/* Screen reader announcements */}
      <Box
        position="absolute"
        width="1px"
        height="1px"
        padding="0"
        margin="-1px"
        overflow="hidden"
        clip="rect(0, 0, 0, 0)"
        whiteSpace="nowrap"
        border="0"
        aria-live="polite"
        aria-atomic="true"
      >
        Word {wordNumber} of {totalWords}: {currentWord}
      </Box>

      {/* Header */}
      <Flex justify="center" py={2}>
        <Text
          color="cyan.400"
          fontSize={headerSize}
          fontWeight="bold"
          fontFamily="'Fira Code', monospace"
        >
          BABOON - Typing Practice
        </Text>
      </Flex>

      {/* Word counter */}
      <Flex justify="center" py={2} role="status" aria-label={`Progress: word ${wordNumber} of ${totalWords}`}>
        <Text
          color="cyan.400"
          fontSize={counterSize}
          fontWeight="bold"
          fontFamily="'Fira Code', monospace"
        >
          Word {wordNumber}/{totalWords}
        </Text>
      </Flex>

      {/* Progress bar */}
      <Box px={{ base: 4, md: 8 }} py={2}>
        <Progress
          value={(wordNumber / totalWords) * 100}
          size="sm"
          borderRadius="full"
          bg="gray.800"
          aria-label={`Progress: ${Math.round((wordNumber / totalWords) * 100)}% complete`}
          sx={{
            '& > div': {
              bgGradient: 'linear(to-r, red.500, orange.400, yellow.400, green.400)',
              transition: 'all 0.5s ease',
            },
          }}
        />
      </Box>

      {/* Main word carousel area */}
      <Flex flex={1} direction="column" align="center" justify="center" position="relative">
        {/* Decorative separator */}
        <Box position="absolute" top="20%">
          <Text color="gray.700" fontFamily="monospace" fontSize={headerSize}>
            ─────────────────────────────────
          </Text>
        </Box>

        {/* Previous word */}
        <AnimatePresence mode="wait">
          {previousWord && (
            <PreviewWord key={`prev-${wordKey}`} word={previousWord} position="previous" />
          )}
        </AnimatePresence>

        {/* Current word - Block font */}
        <VStack spacing={4}>
          <AnimatePresence mode="wait">
            <AnimatedBlockFontWord
              key={`current-${wordKey}`}
              word={currentWord}
              input={currentInput}
              wordKey={wordKey}
              fontSize={blockFontSize}
            />
          </AnimatePresence>

          {/* Extra typed characters (errors) */}
          <AnimatePresence>
            {extraChars && <ExtraChars chars={extraChars} />}
          </AnimatePresence>
        </VStack>

        {/* Next word */}
        <AnimatePresence mode="wait">
          {nextWord && (
            <PreviewWord key={`next-${wordKey}`} word={nextWord} position="next" />
          )}
        </AnimatePresence>

        {/* Decorative separator */}
        <Box position="absolute" bottom="20%">
          <Text color="gray.700" fontFamily="monospace" fontSize={headerSize}>
            ─────────────────────────────────
          </Text>
        </Box>

        {/* Decorative glow effect */}
        <Box
          position="absolute"
          w={{ base: '300px', md: '500px', lg: '600px' }}
          h={{ base: '150px', md: '250px', lg: '300px' }}
          bg="radial-gradient(ellipse, rgba(212, 146, 42, 0.1) 0%, transparent 70%)"
          pointerEvents="none"
          zIndex={0}
        />
      </Flex>

      {/* WPM Bar */}
      <Box pb={{ base: 4, md: 6 }}>
        <Flex justify="center">
          <WpmBar wpm={timerStarted ? liveWpm : 0} />
        </Flex>
      </Box>

      {/* Footer hint */}
      <Flex justify="center" pb={2}>
        <Text
          color="gray.600"
          fontSize={{ base: 'xs', md: 'sm' }}
          fontFamily="'Fira Code', monospace"
          textAlign="center"
        >
          Type to start | SPACE to continue | Last word auto-completes | ESC to quit
        </Text>
      </Flex>

      {/* Kartoza branding */}
      <Flex justify="center" pb={4}>
        <HStack spacing={2} color="gray.600" fontSize={{ base: 'xs', md: 'sm' }} fontFamily="'Fira Code', monospace">
          <Text>Made with</Text>
          <Text color="red.400">♥</Text>
          <Text>by</Text>
          <Link
            href="https://kartoza.com"
            isExternal
            color="cyan.500"
            _hover={{ color: 'cyan.400' }}
          >
            Kartoza
          </Link>
          <Text>|</Text>
          <Link
            href="https://github.com/sponsors/timlinux"
            isExternal
            color="cyan.500"
            _hover={{ color: 'cyan.400' }}
          >
            Donate!
          </Link>
          <Text>|</Text>
          <Link
            href="https://github.com/timlinux/baboon"
            isExternal
            color="gray.500"
            _hover={{ color: 'gray.400' }}
          >
            GitHub
          </Link>
        </HStack>
      </Flex>

      {/* AdSense ad */}
      <Box pb={4}>
        <AdSense publisherId={adsenseKey} showPreview={true} />
      </Box>
    </Flex>
  );
}

export default TypingScreen;
