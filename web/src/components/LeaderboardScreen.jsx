import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Flex,
  Select,
  Link,
  Spinner,
  useBreakpointValue,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';
import api from '../api.js';
import { BlockFontWord, GRADIENT_COLORS } from './BlockFont.jsx';
import { useAuth } from '../contexts/AuthContext.jsx';
import ShareBadge from './ShareBadge.jsx';
import AdSense from './AdSense.jsx';

const MotionBox = motion(Box);
const MotionFlex = motion(Flex);

// Get gradient color at position
function getGradientColor(position) {
  const idx = Math.min(Math.floor(position * (GRADIENT_COLORS.length - 1)), GRADIENT_COLORS.length - 1);
  return GRADIENT_COLORS[Math.max(0, idx)];
}

// Responsive gradient bar for WPM display
function WpmBar({ wpm, maxWpm = 120, barWidth = 20 }) {
  const fillPercent = Math.min(Math.max(wpm / maxWpm, 0), 1);
  const filledWidth = Math.floor(barWidth * fillPercent);
  const emptyWidth = barWidth - filledWidth;

  return (
    <HStack spacing={0} display="inline-flex">
      {Array.from({ length: filledWidth }).map((_, i) => (
        <Text key={`fill-${i}`} as="span" color={getGradientColor(i / barWidth)}>█</Text>
      ))}
      {Array.from({ length: emptyWidth }).map((_, i) => (
        <Text key={`empty-${i}`} as="span" color="gray.700">░</Text>
      ))}
    </HStack>
  );
}

// Rank badge with color
function RankBadge({ rank }) {
  const fontSize = useBreakpointValue({ base: 'lg', sm: 'xl', md: '2xl', lg: '3xl' });
  let color = 'gray.500';
  let symbol = '●';

  if (rank === 1) {
    color = '#FFD700'; // Gold
    symbol = '★';
  } else if (rank === 2) {
    color = '#C0C0C0'; // Silver
    symbol = '★';
  } else if (rank === 3) {
    color = '#CD7F32'; // Bronze
    symbol = '★';
  }

  return (
    <HStack spacing={1}>
      <Text color={color} fontSize={fontSize}>{symbol}</Text>
      <Text color={color} fontWeight="bold" fontSize={fontSize}>#{rank}</Text>
    </HStack>
  );
}

// Leaderboard entry row
function LeaderboardRow({ entry, isCurrentUser, delay, onShare }) {
  // Use larger fonts to maximize readability
  const fontSize = useBreakpointValue({ base: 'sm', sm: 'md', md: 'lg', lg: 'xl', xl: '2xl' });
  const barWidth = useBreakpointValue({ base: 8, sm: 12, md: 16, lg: 20, xl: 24 });

  // Calculate score if not provided
  const score = entry.score || (entry.wpm * entry.accuracy / 100);

  return (
    <MotionBox
      initial={{ opacity: 0, x: -30 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay, type: 'spring', stiffness: 200, damping: 20 }}
      w="100%"
    >
      <HStack
        w="100%"
        py={{ base: 2, md: 3, lg: 4 }}
        px={4}
        bg={isCurrentUser ? 'rgba(212, 146, 42, 0.2)' : 'transparent'}
        borderRadius="md"
        border={isCurrentUser ? '2px solid' : 'none'}
        borderColor="#D4922A"
        fontFamily="'Fira Code', monospace"
        fontSize={fontSize}
        justify="space-between"
        _hover={{ bg: 'rgba(255, 255, 255, 0.05)' }}
      >
        {/* Rank */}
        <Box w={{ base: '50px', sm: '60px', md: '70px', lg: '80px' }}>
          <RankBadge rank={entry.rank} />
        </Box>

        {/* Username */}
        <Text
          w={{ base: '70px', sm: '90px', md: '110px', lg: '130px' }}
          color={isCurrentUser ? '#D4922A' : 'cyan.400'}
          fontWeight="bold"
          textAlign="left"
          isTruncated
        >
          {entry.username || 'ANON'}
        </Text>

        {/* Message */}
        <Text
          flex={1}
          color={isCurrentUser ? '#D4922A' : 'gray.300'}
          textAlign="left"
          isTruncated
          display={{ base: 'none', md: 'block' }}
        >
          {entry.display_name}
        </Text>

        {/* Score with bar */}
        <HStack spacing={2} w={{ base: '100px', sm: '140px', md: '180px', lg: '220px' }}>
          <Text color="yellow.400" w={{ base: '50px', md: '70px', lg: '80px' }} textAlign="right" fontWeight="bold">
            {score.toFixed(1)}
          </Text>
          <Box display={{ base: 'none', md: 'block' }}>
            <WpmBar wpm={score} maxWpm={100} barWidth={barWidth} />
          </Box>
        </HStack>

        {/* WPM */}
        <Text color="#00ff00" w={{ base: '50px', sm: '60px', md: '70px', lg: '80px' }} textAlign="right">
          {entry.wpm.toFixed(1)}
        </Text>

        {/* Accuracy */}
        <Text color="cyan.400" w={{ base: '50px', sm: '60px', md: '70px', lg: '80px' }} textAlign="right" display={{ base: 'none', sm: 'block' }}>
          {entry.accuracy.toFixed(0)}%
        </Text>

        {/* Share button */}
        <Box w={{ base: '60px', md: '70px' }} textAlign="center">
          <Box
            as="button"
            onClick={() => onShare(entry)}
            px={{ base: 2, md: 3 }}
            py={1}
            borderRadius="md"
            border="1px solid"
            borderColor="cyan.500"
            color="cyan.400"
            fontSize={{ base: 'xs', md: 'sm' }}
            fontWeight="bold"
            _hover={{ bg: 'rgba(74, 144, 164, 0.2)', borderColor: 'cyan.400' }}
            cursor="pointer"
          >
            SHARE
          </Box>
        </Box>
      </HStack>
    </MotionBox>
  );
}

function LeaderboardScreen({ onBack, adsenseEnabled, adsenseKey }) {
  const { user, isAuthenticated } = useAuth();
  const [entries, setEntries] = useState([]);
  const [months, setMonths] = useState([]);
  const [selectedMonth, setSelectedMonth] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [shareEntry, setShareEntry] = useState(null);

  // Use larger title size to maximize screen real estate
  const titleSize = useBreakpointValue({ base: 'sm', sm: 'md', md: 'lg', lg: 'xl', xl: '2xl' });

  // Fetch available months
  useEffect(() => {
    const fetchMonths = async () => {
      try {
        const result = await api.getLeaderboardMonths();
        setMonths(result.months || []);
        if (result.months && result.months.length > 0) {
          setSelectedMonth(result.months[0]); // Current/most recent month
        }
      } catch (e) {
        console.error('Failed to fetch months:', e);
        // Default to current month
        const currentMonth = new Date().toISOString().slice(0, 7);
        setMonths([currentMonth]);
        setSelectedMonth(currentMonth);
      }
    };
    fetchMonths();
  }, []);

  // Fetch leaderboard when month changes
  useEffect(() => {
    if (!selectedMonth) return;

    const fetchLeaderboard = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const result = await api.getLeaderboard(selectedMonth);
        setEntries(result.entries || []);
      } catch (e) {
        setError(e.message || 'Failed to load leaderboard');
        setEntries([]);
      }
      setIsLoading(false);
    };
    fetchLeaderboard();
  }, [selectedMonth]);

  // Handle keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === 'Escape') {
        if (shareEntry) {
          setShareEntry(null);
        } else {
          onBack();
        }
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [onBack, shareEntry]);

  // Format month for display
  const formatMonth = (monthYear) => {
    if (!monthYear) return '';
    const [year, month] = monthYear.split('-');
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    const monthIdx = parseInt(month, 10) - 1;
    if (monthIdx >= 0 && monthIdx < 12) {
      return `${months[monthIdx]} ${year}`;
    }
    return monthYear;
  };

  const handleShare = useCallback((entry) => {
    setShareEntry(entry);
  }, []);

  return (
    <Flex minH="100vh" direction="column" p={{ base: 2, md: 4 }}>
      {/* Header */}
      <Flex justify="center" py={4}>
        <MotionBox
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ type: 'spring', bounce: 0.5 }}
        >
          <BlockFontWord word="TOP 10" input="" showCurrent={false} fontSize={titleSize} />
        </MotionBox>
      </Flex>

      {/* Month selector */}
      <MotionFlex
        justify="center"
        mb={4}
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
      >
        <Select
          value={selectedMonth}
          onChange={(e) => setSelectedMonth(e.target.value)}
          maxW={{ base: '180px', md: '220px', lg: '260px' }}
          bg="gray.900"
          borderColor="gray.700"
          color="white"
          fontFamily="'Fira Code', monospace"
          fontSize={{ base: 'md', md: 'lg', lg: 'xl' }}
          h={{ base: '40px', md: '48px', lg: '56px' }}
          _hover={{ borderColor: 'cyan.400' }}
          _focus={{ borderColor: 'cyan.400' }}
        >
          {months.map((month) => (
            <option key={month} value={month} style={{ background: '#1a2833' }}>
              {formatMonth(month)}
            </option>
          ))}
        </Select>
      </MotionFlex>

      {/* Column headers */}
      <MotionBox
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.3 }}
        mb={2}
      >
        <HStack
          w="100%"
          px={4}
          fontFamily="'Fira Code', monospace"
          fontSize={{ base: 'sm', sm: 'md', md: 'lg', lg: 'xl' }}
          color="gray.500"
          justify="space-between"
        >
          <Box w={{ base: '50px', sm: '60px', md: '70px', lg: '80px' }} textAlign="left">#</Box>
          <Box w={{ base: '70px', sm: '90px', md: '110px', lg: '130px' }} textAlign="left">USER</Box>
          <Text flex={1} display={{ base: 'none', md: 'block' }} textAlign="left">MESSAGE</Text>
          <Box w={{ base: '100px', sm: '140px', md: '180px', lg: '220px' }} textAlign="right">SCORE</Box>
          <Box w={{ base: '50px', sm: '60px', md: '70px', lg: '80px' }} textAlign="right">WPM</Box>
          <Box w={{ base: '50px', sm: '60px', md: '70px', lg: '80px' }} textAlign="right" display={{ base: 'none', sm: 'block' }}>ACC</Box>
          <Box w={{ base: '60px', md: '70px' }} textAlign="center" color="cyan.500">SHARE</Box>
        </HStack>
      </MotionBox>

      {/* Divider */}
      <Box h="1px" bg="gray.700" mx={4} mb={2} />

      {/* Leaderboard entries */}
      <VStack spacing={1} flex={1} overflow="auto" align="stretch">
        {isLoading && (
          <Flex justify="center" py={8}>
            <Spinner color="cyan.400" size="lg" />
          </Flex>
        )}

        {error && (
          <Text color="red.400" textAlign="center" py={4} fontFamily="'Fira Code', monospace">
            {error}
          </Text>
        )}

        {!isLoading && !error && entries.length === 0 && (
          <VStack py={8} spacing={4}>
            <Text color="gray.500" textAlign="center" fontFamily="'Fira Code', monospace">
              No entries yet for {formatMonth(selectedMonth)}
            </Text>
            <Text color="gray.600" textAlign="center" fontFamily="'Fira Code', monospace" fontSize="sm">
              Be the first to claim a spot!
            </Text>
          </VStack>
        )}

        {!isLoading && entries.map((entry, idx) => (
          <LeaderboardRow
            key={entry.id}
            entry={entry}
            isCurrentUser={isAuthenticated && user && entry.user_id === user.id}
            delay={0.4 + idx * 0.05}
            onShare={handleShare}
          />
        ))}
      </VStack>

      {/* Back button */}
      <MotionFlex
        justify="center"
        py={4}
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.8 }}
      >
        <Box
          as="button"
          onClick={onBack}
          px={{ base: 6, md: 8, lg: 10 }}
          py={{ base: 2, md: 3 }}
          borderRadius="md"
          border="2px solid"
          borderColor="gray.600"
          color="gray.400"
          fontFamily="'Fira Code', monospace"
          fontSize={{ base: 'md', md: 'lg', lg: 'xl' }}
          _hover={{ borderColor: 'cyan.400', color: 'cyan.400' }}
        >
          BACK
        </Box>
      </MotionFlex>

      {/* Keyboard hint */}
      <Text color="gray.600" fontSize={{ base: 'sm', md: 'md' }} textAlign="center" pb={2} fontFamily="'Fira Code', monospace">
        ESC = back
      </Text>

      {/* Kartoza branding */}
      <Flex justify="center" pb={4}>
        <HStack spacing={2} color="gray.600" fontSize={{ base: 'sm', md: 'md' }} fontFamily="'Fira Code', monospace">
          <Text>Made with</Text>
          <Text color="red.400">♥</Text>
          <Text>by</Text>
          <Link href="https://kartoza.com" isExternal color="cyan.500" _hover={{ color: 'cyan.400' }}>
            Kartoza
          </Link>
          <Text>|</Text>
          <Link href="https://github.com/sponsors/timlinux" isExternal color="cyan.500" _hover={{ color: 'cyan.400' }}>
            Donate!
          </Link>
          <Text>|</Text>
          <Link href="https://github.com/timlinux/baboon" isExternal color="gray.500" _hover={{ color: 'gray.400' }}>
            GitHub
          </Link>
        </HStack>
      </Flex>

      {/* AdSense ad */}
      <Box pb={4}>
        <AdSense publisherId={adsenseKey} showPreview={true} />
      </Box>

      {/* Share modal */}
      {shareEntry && (
        <ShareBadge
          entry={shareEntry}
          onClose={() => setShareEntry(null)}
        />
      )}
    </Flex>
  );
}

export default LeaderboardScreen;
