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
import { BlockFontWord, WpmBar } from './BlockFont.jsx';

const MotionBox = motion(Box);

// Calculate optimal font size based on word length and viewport
function useBlockFontSize(wordLength) {
  const baseSizes = useBreakpointValue({
    base: { short: '4xl', medium: '3xl', long: '2xl' },
    sm: { short: '5xl', medium: '4xl', long: '3xl' },
    md: { short: '6xl', medium: '5xl', long: '4xl' },
    lg: { short: '7xl', medium: '6xl', long: '5xl' },
    xl: { short: '8xl', medium: '7xl', long: '6xl' },
    '2xl': { short: '9xl', medium: '8xl', long: '7xl' },
  }) || { short: '5xl', medium: '4xl', long: '3xl' };

  if (wordLength <= 4) return baseSizes.short;
  if (wordLength <= 7) return baseSizes.medium;
  return baseSizes.long;
}

// Slot machine word reel - all words slide up together on transition
function WordReel({ previousWord, currentWord, currentInput, nextWords, wordKey, blockFontSize }) {
  const previewSize = useBreakpointValue({
    base: 'lg',
    md: '2xl',
    lg: '3xl',
    xl: '4xl',
  });
  const nextSize = useBreakpointValue({
    base: 'md',
    md: 'lg',
    lg: 'xl',
    xl: '2xl',
  });

  return (
    <AnimatePresence mode="wait" initial={false}>
      <MotionBox
        key={wordKey}
        initial={{ y: 60 }}
        animate={{ y: 0 }}
        exit={{ y: -60 }}
        transition={{
          duration: 0.15,
          ease: 'easeInOut',
        }}
        display="flex"
        flexDirection="column"
        alignItems="center"
        gap={6}
      >
        {/* Previous word - faded, small, above */}
        <Box minH="2em" display="flex" alignItems="center" justifyContent="center">
          {previousWord && (
            <Text
              fontFamily="'Fira Code', 'JetBrains Mono', monospace"
              fontSize={previewSize}
              color="gray.600"
              textTransform="lowercase"
              letterSpacing="0.2em"
              opacity={0.4}
            >
              {previousWord}
            </Text>
          )}
        </Box>

        {/* Current word - large, bold, colored */}
        <BlockFontWord
          word={currentWord}
          input={currentInput}
          showCurrent={true}
          fontSize={blockFontSize}
        />

        {/* Next words - fading out below */}
        <VStack spacing={1} minH="4em">
          {nextWords.map((word, index) => (
            <Text
              key={`next-${index}-${word}`}
              fontFamily="'Fira Code', 'JetBrains Mono', monospace"
              fontSize={nextSize}
              color={index === 0 ? 'gray.400' : index === 1 ? 'gray.500' : 'gray.600'}
              textTransform="lowercase"
              letterSpacing="0.15em"
              opacity={1 - index * 0.25}
            >
              {index === 0 ? `▼  ${word}  ▼` : word}
            </Text>
          ))}
        </VStack>
      </MotionBox>
    </AnimatePresence>
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
  onClearInput,
  onSpace,
  onExit,
  onRestart,
  adsenseEnabled,
  adsenseKey,
  perfectMode,
  practiceMode,
}) {
  const currentWord = gameState?.current_word || '';
  const currentInput = gameState?.current_input || '';
  const previousWord = gameState?.previous_word || '';
  const nextWords = gameState?.next_words || [];
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
      if (e.key === 'Escape') {
        onExit();
        return;
      }

      if (e.key === 'Tab') {
        e.preventDefault();
        onRestart();
        return;
      }

      // Ctrl+Backspace clears entire word input
      if (e.key === 'Backspace' && (e.ctrlKey || e.metaKey)) {
        e.preventDefault();
        onClearInput();
        return;
      }

      // Ignore other modifier key combos (except shift)
      if (e.ctrlKey || e.metaKey || e.altKey) return;

      if (e.key === 'Backspace') {
        e.preventDefault();
        onBackspace();
        return;
      }

      if (e.key === ' ' || e.key === 'Enter') {
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
  }, [onKeystroke, onBackspace, onClearInput, onSpace, onExit, onRestart]);

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
        {practiceMode === 'ngrams' ? 'N-gram' : 'Word'} {wordNumber} of {totalWords}: {currentWord}
      </Box>

      {/* Header */}
      <Flex justify="center" py={2}>
        <Text
          color="cyan.400"
          fontSize={headerSize}
          fontWeight="bold"
          fontFamily="'Fira Code', monospace"
        >
          {practiceMode === 'ngrams' ? 'BABOON - N-gram Training' : 'BABOON - Typing Practice'}
        </Text>
      </Flex>

      {/* Word counter */}
      <Flex justify="center" py={2} role="status" aria-label={`Progress: ${practiceMode === 'ngrams' ? 'n-gram' : 'word'} ${wordNumber} of ${totalWords}`}>
        <Text
          color="cyan.400"
          fontSize={counterSize}
          fontWeight="bold"
          fontFamily="'Fira Code', monospace"
        >
          {practiceMode === 'ngrams' ? `N-gram ${wordNumber}/${totalWords}` : `Word ${wordNumber}/${totalWords}`}
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

      {/* Main word reel area - slot machine style */}
      <Flex flex={1} direction="column" align="center" justify="center" position="relative" overflow="hidden">
        {/* Decorative separator */}
        <Box position="absolute" top="20%">
          <Text color="gray.700" fontFamily="monospace" fontSize={headerSize}>
            ─────────────────────────────────
          </Text>
        </Box>

        {/* Slot machine word reel */}
        <WordReel
          previousWord={previousWord}
          currentWord={currentWord}
          currentInput={currentInput}
          nextWords={nextWords}
          wordKey={wordKey}
          blockFontSize={blockFontSize}
        />

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
        <HStack spacing={2}>
          <Text
            color="gray.600"
            fontSize={{ base: 'xs', md: 'sm' }}
            fontFamily="'Fira Code', monospace"
            textAlign="center"
          >
            Type to start | SPACE/ENTER to continue | TAB to restart | ESC to quit
          </Text>
          {perfectMode && (
            <Text
              color="red.400"
              fontSize={{ base: 'xs', md: 'sm' }}
              fontFamily="'Fira Code', monospace"
              fontWeight="bold"
            >
              [PERFECT MODE]
            </Text>
          )}
        </HStack>
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
            href="https://github.com/sponsors/kartoza"
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
