import React from 'react';
import {
  Box,
  VStack,
  Text,
  Button,
  Switch,
  FormControl,
  FormLabel,
  Container,
  Flex,
  Badge,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';

const MotionBox = motion(Box);
const MotionText = motion(Text);

function WelcomeScreen({ isConnected, punctuationMode, setPunctuationMode, onStart, isLoading }) {
  return (
    <Flex minH="100vh" align="center" justify="center" p={8}>
      <Container maxW="container.md">
        <VStack spacing={12}>
          {/* Logo/Title */}
          <VStack spacing={4}>
            <MotionBox
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: 'spring', bounce: 0.5, duration: 0.8 }}
            >
              <Text
                fontSize={{ base: '6xl', md: '8xl' }}
                fontWeight="800"
                bgGradient="linear(to-r, accent.cyan, accent.purple, accent.green)"
                bgClip="text"
                letterSpacing="tight"
              >
                BABOON
              </Text>
            </MotionBox>
            <MotionText
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
              fontSize="xl"
              color="gray.400"
              textAlign="center"
            >
              Master your typing with beautiful practice
            </MotionText>
          </VStack>

          {/* Connection Status */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.5 }}
          >
            <Badge
              colorScheme={isConnected ? 'green' : 'red'}
              fontSize="md"
              px={4}
              py={2}
              borderRadius="full"
            >
              {isConnected ? '● Connected to Backend' : '○ Backend Disconnected'}
            </Badge>
          </MotionBox>

          {/* Options Card */}
          <MotionBox
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6, type: 'spring' }}
            w="100%"
            maxW="400px"
          >
            <Box
              bg="bg.card"
              borderRadius="3xl"
              p={8}
              border="1px solid"
              borderColor="whiteAlpha.100"
              boxShadow="0 20px 60px rgba(0, 0, 0, 0.4)"
            >
              <VStack spacing={6}>
                <Text fontSize="2xl" fontWeight="bold" color="white">
                  Game Options
                </Text>

                <FormControl display="flex" alignItems="center" justifyContent="space-between">
                  <FormLabel mb="0" fontSize="lg" color="gray.300">
                    Punctuation Mode
                  </FormLabel>
                  <Switch
                    size="lg"
                    colorScheme="cyan"
                    isChecked={punctuationMode}
                    onChange={(e) => setPunctuationMode(e.target.checked)}
                  />
                </FormControl>

                <Text fontSize="sm" color="gray.500" textAlign="center">
                  {punctuationMode
                    ? 'Words will be separated by punctuation'
                    : 'Standard word-by-word practice'}
                </Text>
              </VStack>
            </Box>
          </MotionBox>

          {/* Start Button */}
          <MotionBox
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.8, type: 'spring', bounce: 0.4 }}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <Button
              size="huge"
              variant="glow"
              onClick={onStart}
              isDisabled={!isConnected}
              isLoading={isLoading}
              loadingText="Starting..."
              _disabled={{
                opacity: 0.5,
                cursor: 'not-allowed',
                boxShadow: 'none',
              }}
            >
              Start Typing
            </Button>
          </MotionBox>

          {/* Instructions */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1 }}
          >
            <VStack spacing={2} color="gray.500" fontSize="sm">
              <Text>Type the words as they appear</Text>
              <Text>Press SPACE to advance to the next word</Text>
              <Text>Press ESC to exit at any time</Text>
            </VStack>
          </MotionBox>
        </VStack>
      </Container>
    </Flex>
  );
}

export default WelcomeScreen;
