import React, { useEffect, useState, useRef } from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Flex,
  Progress,
} from '@chakra-ui/react';
import { motion, useSpring, useTransform, AnimatePresence } from 'framer-motion';

const MotionBox = motion(Box);
const MotionText = motion(Text);
const MotionFlex = motion(Flex);

// Large block letter component with physics
function BlockLetter({ char, status, index }) {
  // Kartoza color scheme
  const colors = {
    correct: '#4CAF50', // Green for correct
    incorrect: '#E53935', // Red for incorrect
    pending: '#6A6A6A', // Kartoza gray
    current: '#D4922A', // Kartoza orange for current
  };

  const bgColors = {
    correct: 'rgba(76, 175, 80, 0.15)',
    incorrect: 'rgba(229, 57, 53, 0.15)',
    pending: 'transparent',
    current: 'rgba(212, 146, 42, 0.15)',
  };

  const spring = useSpring(0, { stiffness: 500, damping: 30 });

  useEffect(() => {
    if (status === 'correct' || status === 'incorrect') {
      spring.set(1);
      setTimeout(() => spring.set(0), 150);
    }
  }, [status, spring]);

  const scale = useTransform(spring, [0, 1], [1, status === 'correct' ? 1.1 : 0.95]);
  const rotate = useTransform(spring, [0, 1], [0, status === 'incorrect' ? 5 : 0]);

  return (
    <MotionBox
      style={{ scale, rotate }}
      initial={{ opacity: 0, y: 20, scale: 0.8 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{
        delay: index * 0.02,
        type: 'spring',
        stiffness: 400,
        damping: 25
      }}
    >
      <Flex
        w={{ base: '50px', md: '70px', lg: '80px' }}
        h={{ base: '70px', md: '90px', lg: '100px' }}
        align="center"
        justify="center"
        bg={bgColors[status]}
        borderRadius="2xl"
        border="3px solid"
        borderColor={status === 'current' ? 'brand.500' : 'transparent'}
        boxShadow={status === 'current' ? '0 0 30px rgba(212, 146, 42, 0.4)' : 'none'}
        transition="all 0.2s ease"
      >
        <Text
          fontSize={{ base: '3xl', md: '4xl', lg: '5xl' }}
          fontWeight="800"
          fontFamily="mono"
          color={colors[status]}
          textTransform="lowercase"
          transition="color 0.15s ease"
        >
          {char}
        </Text>
      </Flex>
    </MotionBox>
  );
}

// Animated word display for carousel effect
function CarouselWord({ word, position, isCompleted }) {
  // Position: 'previous', 'current', 'next'
  const variants = {
    previous: {
      y: -120,
      scale: 0.4,
      opacity: 0.3,
      filter: 'blur(2px)',
      transition: {
        type: 'spring',
        stiffness: 200,
        damping: 25,
        mass: 1,
      }
    },
    current: {
      y: 0,
      scale: 1,
      opacity: 1,
      filter: 'blur(0px)',
      transition: {
        type: 'spring',
        stiffness: 300,
        damping: 30,
        mass: 0.8,
      }
    },
    next: {
      y: 100,
      scale: 0.5,
      opacity: 0.4,
      filter: 'blur(1px)',
      transition: {
        type: 'spring',
        stiffness: 200,
        damping: 25,
        mass: 1,
      }
    },
    exitUp: {
      y: -200,
      scale: 0.3,
      opacity: 0,
      filter: 'blur(4px)',
      transition: {
        type: 'spring',
        stiffness: 200,
        damping: 30,
      }
    },
    enterBelow: {
      y: 200,
      scale: 0.3,
      opacity: 0,
      filter: 'blur(4px)',
    }
  };

  return (
    <MotionBox
      position="absolute"
      variants={variants}
      initial={position === 'next' ? 'enterBelow' : false}
      animate={position}
      exit="exitUp"
      style={{
        originX: 0.5,
        originY: 0.5,
        zIndex: position === 'current' ? 10 : 1,
      }}
    >
      <Text
        fontSize={position === 'current' ? { base: '4xl', md: '6xl', lg: '7xl' } : { base: '2xl', md: '3xl', lg: '4xl' }}
        fontWeight="700"
        fontFamily="mono"
        color={position === 'current' ? 'white' : 'gray.500'}
        textTransform="lowercase"
        letterSpacing="0.15em"
        textShadow={position === 'current' ? '0 0 40px rgba(212, 146, 42, 0.3)' : 'none'}
        whiteSpace="nowrap"
      >
        {word}
      </Text>
    </MotionBox>
  );
}

// Word Carousel component
function WordCarousel({ previousWord, currentWord, nextWord, wordKey }) {
  return (
    <Box
      position="relative"
      h={{ base: '200px', md: '280px', lg: '320px' }}
      w="100%"
      overflow="hidden"
    >
      {/* Gradient overlay top */}
      <Box
        position="absolute"
        top={0}
        left={0}
        right={0}
        h="60px"
        bgGradient="linear(to-b, bg.primary, transparent)"
        zIndex={20}
        pointerEvents="none"
      />

      {/* Gradient overlay bottom */}
      <Box
        position="absolute"
        bottom={0}
        left={0}
        right={0}
        h="60px"
        bgGradient="linear(to-t, bg.primary, transparent)"
        zIndex={20}
        pointerEvents="none"
      />

      <Flex
        position="relative"
        h="100%"
        align="center"
        justify="center"
      >
        <AnimatePresence mode="popLayout">
          {/* Previous word - above */}
          {previousWord && (
            <CarouselWord
              key={`prev-${wordKey}`}
              word={previousWord}
              position="previous"
            />
          )}

          {/* Current word - center (just text, letters shown separately) */}
          <MotionBox
            key={`current-${wordKey}`}
            position="absolute"
            initial={{ y: 100, scale: 0.5, opacity: 0 }}
            animate={{ y: 0, scale: 1, opacity: 1 }}
            exit={{ y: -200, scale: 0.3, opacity: 0 }}
            transition={{
              type: 'spring',
              stiffness: 300,
              damping: 30,
            }}
            style={{ zIndex: 10 }}
          >
            {/* This is a placeholder - actual letters rendered below */}
          </MotionBox>

          {/* Next word - below */}
          {nextWord && (
            <CarouselWord
              key={`next-${wordKey}`}
              word={nextWord}
              position="next"
            />
          )}
        </AnimatePresence>
      </Flex>
    </Box>
  );
}

// WPM Progress Bar with gradient
function WpmBar({ wpm, maxWpm = 120 }) {
  const percentage = Math.min((wpm / maxWpm) * 100, 100);

  const getColor = (wpm) => {
    if (wpm < 30) return 'red.400';
    if (wpm < 50) return 'brand.500'; // Kartoza orange
    if (wpm < 70) return 'yellow.400';
    if (wpm < 90) return 'green.400';
    return 'kartoza.blue.500'; // Kartoza blue for best
  };

  return (
    <Box w="100%" maxW="600px">
      <HStack justify="space-between" mb={2}>
        <Text color="gray.500" fontSize="sm">WPM</Text>
        <HStack spacing={4}>
          <Text color="gray.600" fontSize="xs">0</Text>
          <Text color="gray.600" fontSize="xs">60</Text>
          <Text color="gray.600" fontSize="xs">120</Text>
        </HStack>
      </HStack>
      <Box
        h="20px"
        bg="bg.tertiary"
        borderRadius="full"
        overflow="hidden"
        position="relative"
      >
        <MotionBox
          h="100%"
          bgGradient="linear(to-r, red.500, brand.500, yellow.400, green.400, kartoza.blue.500)"
          borderRadius="full"
          initial={{ width: 0 }}
          animate={{ width: `${percentage}%` }}
          transition={{ type: 'spring', stiffness: 100, damping: 20 }}
        />
      </Box>
      <Flex justify="center" mt={3}>
        <MotionBox
          initial={{ scale: 0.8 }}
          animate={{ scale: 1 }}
          key={Math.floor(wpm)}
        >
          <Text
            fontSize="4xl"
            fontWeight="800"
            color={getColor(wpm)}
            fontFamily="mono"
          >
            {Math.round(wpm)}
          </Text>
        </MotionBox>
      </Flex>
    </Box>
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

  // Determine letter statuses
  const getLetterStatus = (index) => {
    if (index < currentInput.length) {
      return currentInput[index] === currentWord[index] ? 'correct' : 'incorrect';
    }
    if (index === currentInput.length) {
      return 'current';
    }
    return 'pending';
  };

  return (
    <Flex minH="100vh" direction="column" p={4} overflow="hidden">
      {/* Header with word counter */}
      <Flex justify="center" py={4}>
        <VStack spacing={0}>
          <Text color="gray.400" fontSize="sm">
            Word
          </Text>
          <Text color="white" fontSize="2xl" fontWeight="bold">
            {wordNumber} / {totalWords}
          </Text>
        </VStack>
      </Flex>

      {/* Progress bar */}
      <Box px={8} py={2}>
        <Progress
          value={(wordNumber / totalWords) * 100}
          size="sm"
          borderRadius="full"
          bg="bg.tertiary"
          sx={{
            '& > div': {
              bgGradient: 'linear(to-r, brand.500, kartoza.blue.500)',
              transition: 'all 0.5s ease',
            },
          }}
        />
      </Box>

      {/* Main word carousel area */}
      <Flex flex={1} direction="column" align="center" justify="center" position="relative">
        {/* Previous word - floating above with fade */}
        <AnimatePresence mode="wait">
          <MotionBox
            key={`prev-display-${wordKey}`}
            position="absolute"
            top={{ base: '5%', md: '10%' }}
            initial={{ y: 0, opacity: 0, scale: 0.6 }}
            animate={{ y: 0, opacity: 0.35, scale: 0.5 }}
            exit={{ y: -50, opacity: 0, scale: 0.3 }}
            transition={{
              type: 'spring',
              stiffness: 200,
              damping: 25,
            }}
          >
            <Text
              fontSize={{ base: '2xl', md: '4xl', lg: '5xl' }}
              fontWeight="600"
              fontFamily="mono"
              color="gray.500"
              textTransform="lowercase"
              letterSpacing="0.1em"
              filter="blur(1px)"
            >
              {previousWord}
            </Text>
          </MotionBox>
        </AnimatePresence>

        {/* Current word letters - main display */}
        <VStack spacing={6}>
          <AnimatePresence mode="wait">
            <MotionBox
              key={`letters-${wordKey}`}
              initial={{ y: 80, opacity: 0, scale: 0.7 }}
              animate={{ y: 0, opacity: 1, scale: 1 }}
              exit={{ y: -100, opacity: 0, scale: 0.5 }}
              transition={{
                type: 'spring',
                stiffness: 300,
                damping: 28,
                mass: 0.8,
              }}
            >
              <HStack spacing={{ base: 1, md: 2, lg: 3 }} flexWrap="wrap" justify="center">
                {currentWord.split('').map((char, index) => (
                  <BlockLetter
                    key={`${wordKey}-${char}-${index}`}
                    char={char}
                    status={getLetterStatus(index)}
                    index={index}
                  />
                ))}
              </HStack>
            </MotionBox>
          </AnimatePresence>

          {/* Extra typed characters (errors) */}
          <AnimatePresence>
            {currentInput.length > currentWord.length && (
              <MotionBox
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
              >
                <HStack spacing={1}>
                  {currentInput.slice(currentWord.length).split('').map((char, index) => (
                    <MotionBox
                      key={`extra-${index}`}
                      initial={{ scale: 0, rotate: -10 }}
                      animate={{ scale: 1, rotate: [0, -5, 5, 0] }}
                      transition={{ type: 'spring', stiffness: 500, damping: 20 }}
                    >
                      <Box
                        px={3}
                        py={1}
                        bg="rgba(255, 68, 102, 0.2)"
                        borderRadius="lg"
                        border="2px solid"
                        borderColor="accent.red"
                      >
                        <Text color="accent.red" fontFamily="mono" fontSize="xl">
                          {char}
                        </Text>
                      </Box>
                    </MotionBox>
                  ))}
                </HStack>
              </MotionBox>
            )}
          </AnimatePresence>
        </VStack>

        {/* Next word - floating below with fade */}
        <AnimatePresence mode="wait">
          <MotionBox
            key={`next-display-${wordKey}`}
            position="absolute"
            bottom={{ base: '5%', md: '10%' }}
            initial={{ y: 50, opacity: 0, scale: 0.3 }}
            animate={{ y: 0, opacity: 0.4, scale: 0.5 }}
            exit={{ y: 0, opacity: 1, scale: 0.7 }}
            transition={{
              type: 'spring',
              stiffness: 200,
              damping: 25,
            }}
          >
            <Text
              fontSize={{ base: '2xl', md: '4xl', lg: '5xl' }}
              fontWeight="600"
              fontFamily="mono"
              color="gray.500"
              textTransform="lowercase"
              letterSpacing="0.1em"
              filter="blur(1px)"
            >
              {nextWord}
            </Text>
          </MotionBox>
        </AnimatePresence>

        {/* Decorative glow effect behind current word */}
        <Box
          position="absolute"
          w="400px"
          h="200px"
          bg="radial-gradient(ellipse, rgba(212, 146, 42, 0.15) 0%, transparent 70%)"
          pointerEvents="none"
          zIndex={0}
        />
      </Flex>

      {/* WPM Bar */}
      <Box pb={6}>
        <Flex justify="center">
          <WpmBar wpm={timerStarted ? liveWpm : 0} />
        </Flex>
      </Box>

      {/* Footer hint */}
      <Flex justify="center" pb={4}>
        <Text color="gray.600" fontSize="sm">
          Press SPACE when done with word • ESC to exit
        </Text>
      </Flex>
    </Flex>
  );
}

export default TypingScreen;
