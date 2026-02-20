import React, { useEffect } from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Flex,
  Progress,
} from '@chakra-ui/react';
import { motion, useSpring, useTransform } from 'framer-motion';

const MotionBox = motion(Box);

// Large block letter component with physics
function BlockLetter({ char, status, index }) {
  const colors = {
    correct: '#00ff88',
    incorrect: '#ff4466',
    pending: '#4a4a6a',
    current: '#ffffff',
  };

  const bgColors = {
    correct: 'rgba(0, 255, 136, 0.15)',
    incorrect: 'rgba(255, 68, 102, 0.15)',
    pending: 'transparent',
    current: 'rgba(255, 255, 255, 0.1)',
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
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.03 }}
    >
      <Flex
        w={{ base: '50px', md: '70px', lg: '80px' }}
        h={{ base: '70px', md: '90px', lg: '100px' }}
        align="center"
        justify="center"
        bg={bgColors[status]}
        borderRadius="2xl"
        border="3px solid"
        borderColor={status === 'current' ? 'accent.cyan' : 'transparent'}
        boxShadow={status === 'current' ? '0 0 30px rgba(0, 204, 255, 0.4)' : 'none'}
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

// WPM Progress Bar with gradient
function WpmBar({ wpm, maxWpm = 120 }) {
  const percentage = Math.min((wpm / maxWpm) * 100, 100);

  const getColor = (wpm) => {
    if (wpm < 30) return 'red.400';
    if (wpm < 50) return 'orange.400';
    if (wpm < 70) return 'yellow.400';
    if (wpm < 90) return 'green.400';
    return 'cyan.400';
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
          bgGradient="linear(to-r, red.500, orange.400, yellow.400, green.400, cyan.400)"
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
    <Flex minH="100vh" direction="column" p={4}>
      {/* Header */}
      <HStack justify="space-between" px={4} py={2}>
        <MotionBox
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
        >
          <Text color="gray.600" fontSize="lg" fontFamily="mono">
            {previousWord}
          </Text>
        </MotionBox>

        <VStack spacing={0}>
          <Text color="gray.400" fontSize="sm">
            Word
          </Text>
          <Text color="white" fontSize="2xl" fontWeight="bold">
            {wordNumber} / {totalWords}
          </Text>
        </VStack>

        <MotionBox
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
        >
          <Text color="gray.600" fontSize="lg" fontFamily="mono">
            {nextWord}
          </Text>
        </MotionBox>
      </HStack>

      {/* Progress bar */}
      <Box px={8} py={2}>
        <Progress
          value={(wordNumber / totalWords) * 100}
          size="sm"
          borderRadius="full"
          bg="bg.tertiary"
          sx={{
            '& > div': {
              bgGradient: 'linear(to-r, accent.purple, accent.cyan)',
            },
          }}
        />
      </Box>

      {/* Main word display */}
      <Flex flex={1} align="center" justify="center">
        <VStack spacing={8}>
          {/* Current word letters */}
          <HStack spacing={{ base: 1, md: 2 }} flexWrap="wrap" justify="center">
            {currentWord.split('').map((char, index) => (
              <BlockLetter
                key={`${char}-${index}`}
                char={char}
                status={getLetterStatus(index)}
                index={index}
              />
            ))}
          </HStack>

          {/* Extra typed characters (errors) */}
          {currentInput.length > currentWord.length && (
            <HStack spacing={1}>
              {currentInput.slice(currentWord.length).split('').map((char, index) => (
                <MotionBox
                  key={`extra-${index}`}
                  initial={{ scale: 0 }}
                  animate={{ scale: 1, rotate: [0, -5, 5, 0] }}
                  transition={{ type: 'spring' }}
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
          )}
        </VStack>
      </Flex>

      {/* WPM Bar */}
      <Box pb={8}>
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
