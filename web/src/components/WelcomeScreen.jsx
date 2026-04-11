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
  SimpleGrid,
  Link,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';
import { LoginButtons } from './LoginButton.jsx';
import UserMenu from './UserMenu.jsx';
import { useAuth } from '../contexts/AuthContext.jsx';
import { BlockFontWord } from './BlockFont.jsx';

const MotionBox = motion(Box);
const MotionText = motion(Text);

function WelcomeScreen({ isConnected, punctuationMode, setPunctuationMode, onStart, isLoading, localStats }) {
  const { user, isAuthenticated, authConfig } = useAuth();

  return (
    <Flex minH="100vh" align="center" justify="center" p={8}>
      <Container maxW="container.lg">
        <VStack spacing={10}>
          {/* Logo/Title - TUI style block font */}
          <VStack spacing={4}>
            <MotionBox
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: 'spring', bounce: 0.5, duration: 0.8 }}
            >
              <BlockFontWord word="BABOON" input="" showCurrent={false} fontSize="2xs" />
            </MotionBox>
            <MotionText
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
              fontSize="lg"
              color="gray.500"
              textAlign="center"
              fontFamily="'Fira Code', monospace"
            >
              Terminal-style typing practice
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
                  fontFamily="'Fira Code', monospace"
                >
                  {isConnected ? '● Connected' : '○ Disconnected'}
                </Badge>

                {/* Show user menu if authenticated */}
                {isAuthenticated && user && <UserMenu />}
              </HStack>

              {/* Show welcome message for authenticated users */}
              {isAuthenticated && user && (
                <Text color="gray.400" fontSize="sm" fontFamily="'Fira Code', monospace">
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
                bg="gray.900"
                borderRadius="xl"
                p={6}
                border="1px solid"
                borderColor="gray.700"
              >
                <VStack spacing={4}>
                  <Text fontSize="lg" fontWeight="600" color="cyan.400" fontFamily="'Fira Code', monospace">
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
                bg="gray.900"
                borderRadius="xl"
                p={6}
                border="1px solid"
                borderColor="gray.700"
              >
                <VStack spacing={4}>
                  <Text fontSize="lg" fontWeight="600" color="cyan.400" fontFamily="'Fira Code', monospace">
                    Your Progress
                  </Text>
                  <SimpleGrid columns={3} spacing={4} w="100%">
                    <Stat textAlign="center">
                      <StatNumber fontSize="2xl" color="#00ff00" fontFamily="'Fira Code', monospace">
                        {localStats.best_wpm?.toFixed(0) || 0}
                      </StatNumber>
                      <StatLabel color="gray.500" fontSize="xs" fontFamily="'Fira Code', monospace">Best WPM</StatLabel>
                    </Stat>
                    <Stat textAlign="center">
                      <StatNumber fontSize="2xl" color="#00ff00" fontFamily="'Fira Code', monospace">
                        {localStats.best_accuracy?.toFixed(0) || 0}%
                      </StatNumber>
                      <StatLabel color="gray.500" fontSize="xs" fontFamily="'Fira Code', monospace">Best Acc</StatLabel>
                    </Stat>
                    <Stat textAlign="center">
                      <StatNumber fontSize="2xl" color="cyan.400" fontFamily="'Fira Code', monospace">
                        {localStats.total_sessions || 0}
                      </StatNumber>
                      <StatLabel color="gray.500" fontSize="xs" fontFamily="'Fira Code', monospace">Sessions</StatLabel>
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
              bg="gray.900"
              borderRadius="xl"
              p={6}
              border="1px solid"
              borderColor="gray.700"
            >
              <VStack spacing={6}>
                <Text fontSize="lg" fontWeight="bold" color="cyan.400" fontFamily="'Fira Code', monospace">
                  Options
                </Text>

                <FormControl display="flex" alignItems="center" justifyContent="space-between">
                  <FormLabel mb="0" fontSize="md" color="gray.300" fontFamily="'Fira Code', monospace">
                    Punctuation
                  </FormLabel>
                  <Switch
                    size="lg"
                    colorScheme="green"
                    isChecked={punctuationMode}
                    onChange={(e) => setPunctuationMode(e.target.checked)}
                  />
                </FormControl>

                <Text fontSize="xs" color="gray.600" textAlign="center" fontFamily="'Fira Code', monospace">
                  {punctuationMode
                    ? 'Words include punctuation marks'
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
              size="lg"
              variant="outline"
              onClick={onStart}
              isDisabled={!isConnected}
              isLoading={isLoading}
              loadingText="Starting..."
              px={12}
              py={6}
              fontSize="xl"
              fontFamily="'Fira Code', monospace"
              fontWeight="bold"
              borderColor="#00ff00"
              color="#00ff00"
              bg="transparent"
              _hover={{
                bg: 'rgba(0, 255, 0, 0.1)',
                borderColor: '#00ff00',
                boxShadow: '0 0 20px rgba(0, 255, 0, 0.3)',
              }}
              _disabled={{
                opacity: 0.5,
                cursor: 'not-allowed',
                boxShadow: 'none',
              }}
              _focus={{
                outline: '2px solid',
                outlineColor: '#00ff00',
                outlineOffset: '2px',
              }}
              aria-label={!isConnected ? "Start Typing (waiting for backend connection)" : "Start Typing"}
            >
              [ START ]
            </Button>
          </MotionBox>

          {/* Instructions */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1 }}
          >
            <VStack spacing={1} color="gray.600" fontSize="xs" fontFamily="'Fira Code', monospace">
              <Text>Type the words as they appear</Text>
              <Text>Press SPACE to advance | Last word auto-completes</Text>
              <Text>Press ESC to exit</Text>
            </VStack>
          </MotionBox>

          {/* Navigation Links */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1.1 }}
          >
            <HStack spacing={6} color="gray.500" fontSize="sm" fontFamily="'Fira Code', monospace">
              <Link
                href="/docs/"
                color="cyan.500"
                _hover={{ color: 'cyan.400', textDecoration: 'underline' }}
              >
                Documentation
              </Link>
              <Link
                href="https://github.com/timlinux/baboon"
                isExternal
                color="gray.500"
                _hover={{ color: 'gray.400' }}
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
            <HStack spacing={2} color="gray.600" fontSize="xs" fontFamily="'Fira Code', monospace">
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
            </HStack>
          </MotionBox>
        </VStack>
      </Container>
    </Flex>
  );
}

export default WelcomeScreen;
