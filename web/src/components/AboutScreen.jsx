import React, { useEffect } from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Flex,
  Container,
  Link,
  Divider,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';
import { BlockFontWord } from './BlockFont.jsx';
import AdSense from './AdSense.jsx';

const MotionBox = motion(Box);
const MotionText = motion(Text);

function AboutScreen({ onBack, adsenseEnabled, adsenseKey }) {
  // Handle keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === 'Escape') {
        onBack();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [onBack]);

  return (
    <Flex minH="100vh" direction="column" p={{ base: 4, md: 8 }} overflow="auto">
      <Container maxW="container.md">
        <VStack spacing={8} align="stretch">
          {/* Header */}
          <MotionBox
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: 'spring', bounce: 0.5 }}
            textAlign="center"
          >
            <BlockFontWord word="WHY" input="" showCurrent={false} fontSize="6xl" />
          </MotionBox>

          {/* AI Era Section */}
          <MotionBox
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
          >
            <VStack
              bg="gray.900"
              borderRadius="xl"
              p={6}
              border="1px solid"
              borderColor="gray.700"
              spacing={4}
              align="stretch"
            >
              <Text
                fontSize={{ base: 'xl', md: '2xl', lg: '3xl' }}
                fontFamily="'Fira Code', monospace"
                fontWeight="bold"
                color="cyan.400"
              >
                In the Age of AI, Typing Matters More Than Ever
              </Text>
              <Text
                fontSize={{ base: 'md', md: 'lg', lg: 'xl' }}
                fontFamily="'Fira Code', monospace"
                color="gray.300"
                lineHeight="1.8"
              >
                We live in an AI-centric world where tools like ChatGPT, Claude, and countless
                other AI assistants are transforming how we work. But here's the thing: these
                tools are only as fast as you can communicate with them.
              </Text>
              <Text
                fontSize={{ base: 'md', md: 'lg', lg: 'xl' }}
                fontFamily="'Fira Code', monospace"
                color="gray.300"
                lineHeight="1.8"
              >
                The faster and more accurately you type, the more efficiently you can iterate
                with AI, refine your prompts, and turn ideas into reality. Your typing speed
                is now a bottleneck for your productivity in ways we never anticipated.
              </Text>
            </VStack>
          </MotionBox>

          {/* Inspiration Section */}
          <MotionBox
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.4 }}
          >
            <VStack
              bg="gray.900"
              borderRadius="xl"
              p={6}
              border="1px solid"
              borderColor="gray.700"
              spacing={4}
              align="stretch"
            >
              <Text
                fontSize={{ base: 'xl', md: '2xl', lg: '3xl' }}
                fontFamily="'Fira Code', monospace"
                fontWeight="bold"
                color="#D4922A"
              >
                Standing on the Shoulders of Giants
              </Text>
              <Text
                fontSize={{ base: 'md', md: 'lg', lg: 'xl' }}
                fontFamily="'Fira Code', monospace"
                color="gray.300"
                lineHeight="1.8"
              >
                We were inspired by amazing typing practice sites like{' '}
                <Link href="https://monkeytype.com" isExternal color="cyan.400" _hover={{ color: 'cyan.300' }}>
                  Monkeytype
                </Link>{' '}
                and{' '}
                <Link href="https://www.keybr.com" isExternal color="cyan.400" _hover={{ color: 'cyan.300' }}>
                  Keybr
                </Link>
                . These tools have helped millions of people improve their typing skills.
              </Text>
              <Text
                fontSize={{ base: 'md', md: 'lg', lg: 'xl' }}
                fontFamily="'Fira Code', monospace"
                color="gray.300"
                lineHeight="1.8"
              >
                But we wanted to add something new to the typing scene: a fun, competitive,
                retro-styled speed typing game that you can use throughout your day to
                continuously hone your skills.
              </Text>
            </VStack>
          </MotionBox>

          {/* What Makes Baboon Different */}
          <MotionBox
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6 }}
          >
            <VStack
              bg="gray.900"
              borderRadius="xl"
              p={6}
              border="1px solid"
              borderColor="gray.700"
              spacing={4}
              align="stretch"
            >
              <Text
                fontSize={{ base: 'xl', md: '2xl', lg: '3xl' }}
                fontFamily="'Fira Code', monospace"
                fontWeight="bold"
                color="#00ff00"
              >
                What Makes Baboon Different
              </Text>
              <VStack spacing={3} align="stretch" pl={4}>
                <HStack align="flex-start">
                  <Text color="yellow.400" fontFamily="'Fira Code', monospace" fontSize={{ base: 'lg', md: 'xl', lg: '2xl' }}>*</Text>
                  <Text fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace" color="gray.300">
                    <Text as="span" color="cyan.400" fontWeight="bold">Monthly Leaderboards</Text> - Compete with others and see your name in lights
                  </Text>
                </HStack>
                <HStack align="flex-start">
                  <Text color="yellow.400" fontFamily="'Fira Code', monospace" fontSize={{ base: 'lg', md: 'xl', lg: '2xl' }}>*</Text>
                  <Text fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace" color="gray.300">
                    <Text as="span" color="cyan.400" fontWeight="bold">Arcade-Style Name Entry</Text> - Just like the classic games of the 80s
                  </Text>
                </HStack>
                <HStack align="flex-start">
                  <Text color="yellow.400" fontFamily="'Fira Code', monospace" fontSize={{ base: 'lg', md: 'xl', lg: '2xl' }}>*</Text>
                  <Text fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace" color="gray.300">
                    <Text as="span" color="cyan.400" fontWeight="bold">Shareable Badges</Text> - Brag about your achievements on social media
                  </Text>
                </HStack>
                <HStack align="flex-start">
                  <Text color="yellow.400" fontFamily="'Fira Code', monospace" fontSize={{ base: 'lg', md: 'xl', lg: '2xl' }}>*</Text>
                  <Text fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace" color="gray.300">
                    <Text as="span" color="cyan.400" fontWeight="bold">Retro Terminal Aesthetic</Text> - Beautiful TUI and web interfaces
                  </Text>
                </HStack>
                <HStack align="flex-start">
                  <Text color="yellow.400" fontFamily="'Fira Code', monospace" fontSize={{ base: 'lg', md: 'xl', lg: '2xl' }}>*</Text>
                  <Text fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace" color="gray.300">
                    <Text as="span" color="cyan.400" fontWeight="bold">Quick Sessions</Text> - 30 words per round, perfect for breaks
                  </Text>
                </HStack>
                <HStack align="flex-start">
                  <Text color="yellow.400" fontFamily="'Fira Code', monospace" fontSize={{ base: 'lg', md: 'xl', lg: '2xl' }}>*</Text>
                  <Text fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace" color="gray.300">
                    <Text as="span" color="cyan.400" fontWeight="bold">Deep Statistics</Text> - Per-key accuracy, finger stats, common errors
                  </Text>
                </HStack>
              </VStack>
            </VStack>
          </MotionBox>

          {/* Call to Action */}
          <MotionBox
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.8 }}
            textAlign="center"
          >
            <VStack spacing={4}>
              <Text
                fontSize={{ base: 'lg', md: 'xl', lg: '2xl' }}
                fontFamily="'Fira Code', monospace"
                color="gray.400"
              >
                Ready to up your typing game?
              </Text>
              <MotionBox
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <Box
                  as="button"
                  onClick={onBack}
                  px={{ base: 8, md: 12, lg: 16 }}
                  py={{ base: 3, md: 4, lg: 5 }}
                  borderRadius="md"
                  border="2px solid"
                  borderColor="#00ff00"
                  color="#00ff00"
                  fontFamily="'Fira Code', monospace"
                  fontWeight="bold"
                  fontSize={{ base: 'xl', md: '2xl', lg: '3xl' }}
                  _hover={{ bg: 'rgba(0, 255, 0, 0.1)', boxShadow: '0 0 20px rgba(0, 255, 0, 0.3)' }}
                >
                  START TYPING
                </Box>
              </MotionBox>
            </VStack>
          </MotionBox>

          <Divider borderColor="gray.700" />

          {/* Kartoza branding */}
          <MotionBox
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1 }}
          >
            <VStack spacing={2}>
              <HStack spacing={2} color="gray.600" fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace">
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
                <Text>|</Text>
                <Link href="https://github.com/timlinux/baboon" isExternal color="gray.500" _hover={{ color: 'gray.400' }}>
                  GitHub
                </Link>
              </HStack>
              <Text color="gray.700" fontSize={{ base: 'md', md: 'lg', lg: 'xl' }} fontFamily="'Fira Code', monospace">
                ESC = back
              </Text>
            </VStack>
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

export default AboutScreen;
