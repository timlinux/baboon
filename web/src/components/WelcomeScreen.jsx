import React from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Button,
  Switch,
  FormControl,
  FormLabel,
  Container,
  Flex,
  Badge,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  SimpleGrid,
  Link,
  Icon,
  Divider,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';
import { LoginButtons } from './LoginButton.jsx';
import UserMenu from './UserMenu.jsx';
import { useAuth } from '../contexts/AuthContext.jsx';

const MotionBox = motion(Box);
const MotionText = motion(Text);

function WelcomeScreen({ isConnected, punctuationMode, setPunctuationMode, onStart, isLoading, localStats }) {
  const { user, isAuthenticated, authConfig } = useAuth();

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
                bgGradient="linear(to-r, brand.500, kartoza.blue.500, brand.400)"
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

          {/* Connection Status and User Info */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.5 }}
          >
            <VStack spacing={4}>
              <HStack spacing={4}>
                <Badge
                  colorScheme={isConnected ? 'green' : 'red'}
                  fontSize="md"
                  px={4}
                  py={2}
                  borderRadius="full"
                >
                  {isConnected ? '● Connected to Backend' : '○ Backend Disconnected'}
                </Badge>

                {/* Show user menu if authenticated */}
                {isAuthenticated && user && <UserMenu />}
              </HStack>

              {/* Show welcome message for authenticated users */}
              {isAuthenticated && user && (
                <Text color="gray.400" fontSize="sm">
                  Welcome back, {user.display_name || user.email.split('@')[0]}!
                </Text>
              )}
            </VStack>
          </MotionBox>

          {/* Login Section (show only if auth is enabled and user is not logged in) */}
          {authConfig.auth_enabled && !isAuthenticated && (
            <MotionBox
              initial={{ opacity: 0, y: 30 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.52, type: 'spring' }}
              w="100%"
              maxW="400px"
            >
              <Box
                bg="bg.card"
                borderRadius="3xl"
                p={6}
                border="1px solid"
                borderColor="whiteAlpha.100"
                boxShadow="0 20px 60px rgba(0, 0, 0, 0.3)"
              >
                <VStack spacing={4}>
                  <Text fontSize="lg" fontWeight="600" color="white">
                    Sign In
                  </Text>
                  <LoginButtons showTitle={true} />
                </VStack>
              </Box>
            </MotionBox>
          )}

          {/* Personal Stats Card (if available) */}
          {localStats && localStats.total_sessions > 0 && (
            <MotionBox
              initial={{ opacity: 0, y: 30 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.55, type: 'spring' }}
              w="100%"
              maxW="500px"
            >
              <Box
                bg="bg.card"
                borderRadius="3xl"
                p={6}
                border="1px solid"
                borderColor="whiteAlpha.100"
                boxShadow="0 20px 60px rgba(0, 0, 0, 0.3)"
              >
                <VStack spacing={4}>
                  <Text fontSize="lg" fontWeight="600" color="gray.400">
                    Your Progress
                  </Text>
                  <SimpleGrid columns={3} spacing={4} w="100%">
                    <Stat textAlign="center">
                      <StatNumber fontSize="2xl" color="brand.500">
                        {localStats.best_wpm?.toFixed(0) || 0}
                      </StatNumber>
                      <StatLabel color="gray.500" fontSize="xs">Best WPM</StatLabel>
                    </Stat>
                    <Stat textAlign="center">
                      <StatNumber fontSize="2xl" color="green.400">
                        {localStats.best_accuracy?.toFixed(0) || 0}%
                      </StatNumber>
                      <StatLabel color="gray.500" fontSize="xs">Best Accuracy</StatLabel>
                    </Stat>
                    <Stat textAlign="center">
                      <StatNumber fontSize="2xl" color="kartoza.blue.500">
                        {localStats.total_sessions || 0}
                      </StatNumber>
                      <StatLabel color="gray.500" fontSize="xs">Sessions</StatLabel>
                    </Stat>
                  </SimpleGrid>
                </VStack>
              </Box>
            </MotionBox>
          )}

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
                    colorScheme="orange"
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
              aria-label={!isConnected ? "Start Typing (waiting for backend connection)" : "Start Typing"}
              _disabled={{
                opacity: 0.5,
                cursor: 'not-allowed',
                boxShadow: 'none',
              }}
              _focus={{
                outline: '2px solid',
                outlineColor: 'brand.500',
                outlineOffset: '2px',
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

          {/* Navigation Links */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1.1 }}
          >
            <HStack spacing={6} color="gray.500" fontSize="sm">
              <Link
                href="/docs/"
                color="kartoza.blue.400"
                _hover={{ color: 'kartoza.blue.300', textDecoration: 'underline' }}
              >
                Documentation
              </Link>
              <Link
                href="https://github.com/timlinux/baboon"
                isExternal
                color="gray.400"
                _hover={{ color: 'gray.300' }}
              >
                GitHub
              </Link>
            </HStack>
          </MotionBox>

          {/* Footer */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1.2 }}
          >
            <HStack spacing={2} color="gray.600" fontSize="xs">
              <Text>Made with</Text>
              <Text color="red.400">❤️</Text>
              <Text>by</Text>
              <Link
                href="https://kartoza.com"
                isExternal
                color="brand.500"
                _hover={{ color: 'brand.400' }}
              >
                Kartoza
              </Link>
              <Text>|</Text>
              <Link
                href="https://github.com/sponsors/timlinux"
                isExternal
                color="kartoza.blue.500"
                _hover={{ color: 'kartoza.blue.400' }}
              >
                Donate!
              </Link>
            </HStack>
          </MotionBox>
        </VStack>
      </Container>
    </Flex>
  );
}

export default WelcomeScreen;
