import React, { useEffect } from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Button,
  Container,
  Flex,
  Stat,
  StatLabel,
  StatNumber,
  SimpleGrid,
  Link,
  Switch,
  FormControl,
  FormLabel,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';
import { LoginButtons } from './LoginButton.jsx';
import UserMenu from './UserMenu.jsx';
import { useAuth } from '../contexts/AuthContext.jsx';
import { BlockFontWord } from './BlockFont.jsx';
import AdSense from './AdSense.jsx';
import BrandingFooter from './BrandingFooter.jsx';

const MotionBox = motion(Box);
const MotionText = motion(Text);

function WelcomeScreen({ isConnected, onStart, onShowLeaderboard, onShowAbout, isLoading, localStats, adsenseEnabled, adsenseKey, versionInfo, perfectMode, onTogglePerfectMode, settings, setSettings }) {
  const { user, isAuthenticated, authConfig } = useAuth();

  // Handle keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === 'Enter' && isConnected && !isLoading) {
        e.preventDefault();
        onStart();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isConnected, isLoading, onStart]);

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
              <BlockFontWord word="BABOON" input="" showCurrent={false} fontSize="5xl" />
            </MotionBox>
            <MotionText
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
              fontSize={{ base: 'xl', md: '2xl', lg: '3xl' }}
              color="gray.500"
              textAlign="center"
              fontFamily="'Fira Code', monospace"
            >
              Typing Practice Made Fun!
            </MotionText>
          </VStack>

          {/* User Info */}
          {isAuthenticated && user && (
            <MotionBox
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.5 }}
            >
              <VStack spacing={4}>
                <UserMenu />
                <Text color="gray.400" fontSize={{ base: 'lg', md: 'xl', lg: '2xl' }} fontFamily="'Fira Code', monospace">
                  Welcome back, {user.display_name || user.email.split('@')[0]}!
                </Text>
              </VStack>
            </MotionBox>
          )}

          {/* Login Section (show only if auth is enabled with providers and user is not logged in) */}
          {authConfig.auth_enabled && authConfig.providers?.length > 0 && !isAuthenticated && (
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
                  <Text fontSize={{ base: 'xl', md: '2xl', lg: '3xl' }} fontWeight="600" color="cyan.400" fontFamily="'Fira Code', monospace">
                    SIGN IN
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
                  <Text fontSize={{ base: 'xl', md: '2xl', lg: '3xl' }} fontWeight="600" color="cyan.400" fontFamily="'Fira Code', monospace">
                    Your Progress
                  </Text>
                  <SimpleGrid columns={3} spacing={4} w="100%">
                    <Stat textAlign="center">
                      <StatNumber fontSize={{ base: '3xl', md: '4xl', lg: '5xl' }} color="#00ff00" fontFamily="'Fira Code', monospace">
                        {localStats.best_wpm?.toFixed(0) || 0}
                      </StatNumber>
                      <StatLabel color="gray.500" fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace">Best WPM</StatLabel>
                    </Stat>
                    <Stat textAlign="center">
                      <StatNumber fontSize={{ base: '3xl', md: '4xl', lg: '5xl' }} color="#00ff00" fontFamily="'Fira Code', monospace">
                        {localStats.best_accuracy?.toFixed(0) || 0}%
                      </StatNumber>
                      <StatLabel color="gray.500" fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace">Best Acc</StatLabel>
                    </Stat>
                    <Stat textAlign="center">
                      <StatNumber fontSize={{ base: '3xl', md: '4xl', lg: '5xl' }} color="cyan.400" fontFamily="'Fira Code', monospace">
                        {localStats.total_sessions || 0}
                      </StatNumber>
                      <StatLabel color="gray.500" fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace">Sessions</StatLabel>
                    </Stat>
                  </SimpleGrid>
                </VStack>
              </Box>
            </MotionBox>
          )}

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
              px={{ base: 12, md: 16, lg: 20 }}
              py={{ base: 6, md: 8, lg: 10 }}
              fontSize={{ base: '2xl', md: '3xl', lg: '4xl' }}
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
              START
            </Button>
          </MotionBox>

          {/* Perfect Mode Toggle */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.9 }}
          >
            <FormControl display="flex" alignItems="center" justifyContent="center">
              <FormLabel
                htmlFor="perfect-mode"
                mb="0"
                color={perfectMode ? 'red.400' : 'gray.500'}
                fontSize={{ base: 'md', md: 'lg', lg: 'xl' }}
                fontFamily="'Fira Code', monospace"
                fontWeight={perfectMode ? 'bold' : 'normal'}
                cursor="pointer"
              >
                PERFECT MODE
              </FormLabel>
              <Switch
                id="perfect-mode"
                isChecked={perfectMode}
                onChange={onTogglePerfectMode}
                colorScheme="red"
                size="lg"
                sx={{
                  '& .chakra-switch__track': {
                    bg: perfectMode ? 'red.500' : 'gray.600',
                  },
                }}
              />
            </FormControl>
            {perfectMode && (
              <Text
                color="red.400"
                fontSize={{ base: 'xs', md: 'sm' }}
                fontFamily="'Fira Code', monospace"
                textAlign="center"
                mt={1}
              >
                Any mistake restarts the round!
              </Text>
            )}
          </MotionBox>

          {/* Practice Mode Toggle */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.95 }}
          >
            <Box textAlign="center">
              <Text fontSize={{ base: 'sm', md: 'md' }} color="gray.400" mb={2} fontFamily="'Fira Code', monospace">
                Practice Mode
              </Text>
              <HStack spacing={4} justify="center">
                <Button
                  size="sm"
                  variant={settings?.practiceMode === 'words' ? 'solid' : 'outline'}
                  colorScheme="orange"
                  fontFamily="'Fira Code', monospace"
                  onClick={() => setSettings(prev => ({ ...prev, practiceMode: 'words' }))}
                >
                  Words
                </Button>
                <Button
                  size="sm"
                  variant={settings?.practiceMode === 'ngrams' ? 'solid' : 'outline'}
                  colorScheme="orange"
                  fontFamily="'Fira Code', monospace"
                  onClick={() => setSettings(prev => ({ ...prev, practiceMode: 'ngrams' }))}
                >
                  N-grams
                </Button>
              </HStack>
              {settings?.practiceMode === 'ngrams' && (
                <Text
                  color="orange.400"
                  fontSize={{ base: 'xs', md: 'sm' }}
                  fontFamily="'Fira Code', monospace"
                  textAlign="center"
                  mt={1}
                >
                  Practice common letter combinations
                </Text>
              )}
            </Box>
          </MotionBox>

          {/* Instructions */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1 }}
          >
            <Text color="gray.600" fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace">
              Enter to start, Esc to return here, Tab restart
            </Text>
          </MotionBox>

          {/* Navigation Links */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1.1 }}
          >
            <HStack spacing={{ base: 4, md: 8, lg: 10 }} color="gray.500" fontSize={{ base: 'lg', md: 'xl', lg: '2xl' }} fontFamily="'Fira Code', monospace" flexWrap="wrap" justify="center">
              {onShowLeaderboard && (
                <Link
                  as="button"
                  onClick={onShowLeaderboard}
                  color="yellow.400"
                  fontWeight="bold"
                  _hover={{ color: 'yellow.300', textDecoration: 'underline' }}
                >
                  LEADERBOARD
                </Link>
              )}
              {onShowAbout && (
                <Link
                  as="button"
                  onClick={onShowAbout}
                  color="#D4922A"
                  fontWeight="bold"
                  _hover={{ color: '#E5A33B', textDecoration: 'underline' }}
                >
                  ABOUT
                </Link>
              )}
              <Link
                href="https://timlinux.github.io/baboon/"
                isExternal
                color="cyan.500"
                fontWeight="bold"
                _hover={{ color: 'cyan.400', textDecoration: 'underline' }}
              >
                DOCS
              </Link>
              <Link
                href="https://github.com/timlinux/baboon"
                isExternal
                color="gray.500"
                fontWeight="bold"
                _hover={{ color: 'gray.400', textDecoration: 'underline' }}
              >
                GITHUB
              </Link>
            </HStack>
          </MotionBox>

          {/* Footer */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1.2 }}
          >
            <BrandingFooter versionInfo={versionInfo} fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} />
          </MotionBox>

          {/* AdSense ad */}
          <Box pb={4}>
            <AdSense publisherId={adsenseKey} showPreview={true} />
          </Box>
        </VStack>
      </Container>
    </Flex>
  );
}

export default WelcomeScreen;
