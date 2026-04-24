import React, { useState, useEffect, useCallback, useRef } from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Flex,
  Link,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';
import api from '../api.js';
import { BlockFontWord } from './BlockFont.jsx';
import AdSense from './AdSense.jsx';

const MotionBox = motion(Box);
const MotionFlex = motion(Flex);

const MAX_NAME_LENGTH = 10;

// Character slot states
const SLOT_EMPTY = 'empty';
const SLOT_VALID = 'valid';
const SLOT_INVALID = 'invalid';
const SLOT_CURRENT = 'current';

// Slot colors matching the arcade theme
const SLOT_COLORS = {
  [SLOT_EMPTY]: { bg: 'gray.800', border: 'gray.600', text: 'gray.500' },
  [SLOT_VALID]: { bg: 'green.900', border: '#00ff00', text: '#00ff00' },
  [SLOT_INVALID]: { bg: 'red.900', border: 'red.500', text: 'red.400' },
  [SLOT_CURRENT]: { bg: 'gray.800', border: '#D4922A', text: '#D4922A' },
};

function CharacterSlot({ char, state, index }) {
  const colors = SLOT_COLORS[state];
  const isCurrent = state === SLOT_CURRENT;

  return (
    <MotionBox
      w={{ base: '36px', sm: '44px', md: '56px', lg: '68px', xl: '80px' }}
      h={{ base: '52px', sm: '64px', md: '80px', lg: '96px', xl: '112px' }}
      bg={colors.bg}
      border={{ base: '3px solid', md: '4px solid' }}
      borderColor={colors.border}
      borderRadius="md"
      display="flex"
      alignItems="center"
      justifyContent="center"
      initial={{ scale: 0, opacity: 0 }}
      animate={{
        scale: 1,
        opacity: 1,
        boxShadow: isCurrent ? `0 0 25px ${colors.border}` : 'none',
      }}
      transition={{ delay: index * 0.03, type: 'spring', stiffness: 300 }}
    >
      <Text
        fontSize={{ base: '2xl', sm: '3xl', md: '4xl', lg: '5xl', xl: '6xl' }}
        fontFamily="'Fira Code', monospace"
        fontWeight="bold"
        color={colors.text}
      >
        {char || '_'}
      </Text>
    </MotionBox>
  );
}

function NameEntryScreen({ wpm, accuracy, rank, onSubmit, onSkip, adsenseEnabled, adsenseKey }) {
  const [name, setName] = useState('');
  const [slotStates, setSlotStates] = useState(Array(MAX_NAME_LENGTH).fill(SLOT_EMPTY));
  const [validationError, setValidationError] = useState('');
  const [isValidating, setIsValidating] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const validationTimeoutRef = useRef(null);

  // Update slot states based on name and validation
  const updateSlotStates = useCallback((currentName, invalidChars = []) => {
    const states = Array(MAX_NAME_LENGTH).fill(SLOT_EMPTY);
    for (let i = 0; i < currentName.length; i++) {
      if (invalidChars.includes(i)) {
        states[i] = SLOT_INVALID;
      } else {
        states[i] = SLOT_VALID;
      }
    }
    // Mark current slot (first empty after typed chars)
    if (currentName.length < MAX_NAME_LENGTH) {
      states[currentName.length] = SLOT_CURRENT;
    }
    setSlotStates(states);
  }, []);

  // Validate name with debounce
  const validateName = useCallback(async (newName) => {
    if (!newName.trim()) {
      setValidationError('');
      updateSlotStates(newName, []);
      return;
    }

    setIsValidating(true);
    try {
      const result = await api.validateDisplayName(newName);
      if (!result.valid) {
        setValidationError(result.reason);
        updateSlotStates(newName, result.invalid_chars || []);
      } else {
        setValidationError('');
        updateSlotStates(newName, []);
      }
    } catch (e) {
      console.error('Validation error:', e);
      // On error, allow the name (server-side will validate on submit)
      setValidationError('');
      updateSlotStates(newName, []);
    }
    setIsValidating(false);
  }, [updateSlotStates]);

  // Handle keyboard input
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (isSubmitting) return;

      // Handle Enter to submit
      if (e.key === 'Enter') {
        e.preventDefault();
        if (name.trim() && !validationError) {
          handleSubmit();
        }
        return;
      }

      // Handle Escape to skip
      if (e.key === 'Escape') {
        e.preventDefault();
        onSkip();
        return;
      }

      // Handle Backspace
      if (e.key === 'Backspace') {
        e.preventDefault();
        if (name.length > 0) {
          const newName = name.slice(0, -1);
          setName(newName);
          // Clear validation error immediately when user corrects input
          setValidationError('');
          updateSlotStates(newName, []);
          // Debounce validation
          if (validationTimeoutRef.current) {
            clearTimeout(validationTimeoutRef.current);
          }
          validationTimeoutRef.current = setTimeout(() => validateName(newName), 200);
        }
        return;
      }

      // Handle character input (A-Z, 0-9, space, underscore, hyphen)
      if (name.length >= MAX_NAME_LENGTH) return;

      const char = e.key;
      const isValidChar = /^[A-Za-z0-9 _-]$/.test(char);
      if (!isValidChar) return;

      e.preventDefault();
      const newName = name + char.toUpperCase();
      setName(newName);
      // Clear validation error immediately when user corrects input
      setValidationError('');
      updateSlotStates(newName, []);

      // Debounce validation
      if (validationTimeoutRef.current) {
        clearTimeout(validationTimeoutRef.current);
      }
      validationTimeoutRef.current = setTimeout(() => validateName(newName), 200);
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
      if (validationTimeoutRef.current) {
        clearTimeout(validationTimeoutRef.current);
      }
    };
  }, [name, validationError, isSubmitting, onSkip, validateName, updateSlotStates]);

  // Initial slot state
  useEffect(() => {
    updateSlotStates('', []);
  }, [updateSlotStates]);

  const handleSubmit = async () => {
    if (!name.trim() || validationError || isSubmitting) return;

    setIsSubmitting(true);
    try {
      await onSubmit(name.trim());
    } catch (e) {
      setValidationError(e.message || 'Failed to submit');
      setIsSubmitting(false);
    }
  };

  return (
    <Flex minH="100vh" direction="column" align="center" justify="center" p={4}>
      {/* Header */}
      <MotionBox
        initial={{ y: -50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ type: 'spring', bounce: 0.4 }}
        mb={{ base: 4, md: 6, lg: 8 }}
      >
        <BlockFontWord word="TOP 10" input="" showCurrent={false} fontSize="5xl" />
      </MotionBox>

      {/* New High Score announcement */}
      <MotionBox
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ delay: 0.3, type: 'spring', bounce: 0.5 }}
        mb={{ base: 4, md: 6 }}
      >
        <VStack spacing={{ base: 2, md: 4 }}>
          <Text
            fontSize={{ base: '3xl', sm: '4xl', md: '5xl', lg: '6xl' }}
            fontFamily="'Fira Code', monospace"
            fontWeight="bold"
            color="yellow.400"
            textAlign="center"
          >
            NEW HIGH SCORE!
          </Text>
          <HStack spacing={{ base: 4, md: 6, lg: 8 }}>
            <Text
              fontSize={{ base: '2xl', sm: '3xl', md: '4xl', lg: '5xl' }}
              fontFamily="'Fira Code', monospace"
              color="#00ff00"
            >
              {wpm.toFixed(1)} WPM
            </Text>
            <Text
              fontSize={{ base: '2xl', sm: '3xl', md: '4xl', lg: '5xl' }}
              fontFamily="'Fira Code', monospace"
              color="cyan.400"
            >
              {accuracy.toFixed(1)}%
            </Text>
            <Text
              fontSize={{ base: '2xl', sm: '3xl', md: '4xl', lg: '5xl' }}
              fontFamily="'Fira Code', monospace"
              color="#D4922A"
            >
              RANK #{rank}
            </Text>
          </HStack>
        </VStack>
      </MotionBox>

      {/* Enter your message prompt */}
      <MotionBox
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.5 }}
        mb={{ base: 4, md: 6 }}
      >
        <Text
          fontSize={{ base: 'xl', sm: '2xl', md: '3xl', lg: '4xl' }}
          fontFamily="'Fira Code', monospace"
          color="gray.400"
          textAlign="center"
        >
          ENTER YOUR MESSAGE
        </Text>
      </MotionBox>

      {/* Character slots */}
      <MotionFlex
        gap={{ base: 1, md: 2 }}
        mb={6}
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.6 }}
        flexWrap="wrap"
        justify="center"
      >
        {slotStates.map((state, i) => (
          <CharacterSlot
            key={i}
            char={name[i] || ''}
            state={state}
            index={i}
          />
        ))}
      </MotionFlex>

      {/* Validation error */}
      {validationError && (
        <MotionBox
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          mb={4}
        >
          <Text
            fontSize="md"
            fontFamily="'Fira Code', monospace"
            color="red.400"
            textAlign="center"
          >
            {validationError}
          </Text>
        </MotionBox>
      )}

      {/* Instructions */}
      <MotionBox
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.8 }}
        mt={{ base: 4, md: 6 }}
      >
        <VStack spacing={{ base: 1, md: 2 }} color="gray.600" fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace">
          <Text>Type A-Z, 0-9, SPACE, UNDERSCORE, HYPHEN</Text>
          <Text>ENTER to submit | ESC to skip</Text>
          {isValidating && <Text color="cyan.400">Validating...</Text>}
          {isSubmitting && <Text color="cyan.400">Submitting...</Text>}
        </VStack>
      </MotionBox>

      {/* Kartoza branding */}
      <MotionBox
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1 }}
        position="absolute"
        bottom={{ base: 16, md: 20 }}
      >
        <HStack spacing={2} color="gray.600" fontSize={{ base: 'sm', md: 'md' }} fontFamily="'Fira Code', monospace">
          <Text>Made with</Text>
          <Text color="red.400">♥</Text>
          <Text>by</Text>
          <Link href="https://kartoza.com" isExternal color="cyan.500" _hover={{ color: 'cyan.400' }}>
            Kartoza
          </Link>
          <Text>|</Text>
          <Link href="https://github.com/sponsors/kartoza" isExternal color="cyan.500" _hover={{ color: 'cyan.400' }}>
            Donate!
          </Link>
        </HStack>
      </MotionBox>

      {/* AdSense ad */}
      <Box position="absolute" bottom={4} w="100%" px={4}>
        <AdSense publisherId={adsenseKey} showPreview={true} />
      </Box>
    </Flex>
  );
}

export default NameEntryScreen;
