import React, { useState } from 'react';
import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  VStack,
  HStack,
  Text,
  Box,
  Stat,
  StatLabel,
  StatNumber,
  SimpleGrid,
  useToast,
  Icon,
} from '@chakra-ui/react';
import { useAuth } from '../contexts/AuthContext.jsx';

// Icons
const MergeIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <circle cx="18" cy="18" r="3" />
    <circle cx="6" cy="6" r="3" />
    <path d="M6 21V9a9 9 0 009 9" />
  </svg>
);

const TrashIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <polyline points="3,6 5,6 21,6" />
    <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
  </svg>
);

function StatsPreview({ stats, label, color }) {
  if (!stats) return null;

  return (
    <Box
      bg="whiteAlpha.50"
      borderRadius="xl"
      p={4}
      border="1px solid"
      borderColor="whiteAlpha.100"
    >
      <Text fontSize="sm" fontWeight="600" color={color} mb={3}>
        {label}
      </Text>
      <SimpleGrid columns={2} spacing={3}>
        <Stat size="sm">
          <StatNumber fontSize="lg" color="white">
            {stats.best_wpm?.toFixed(0) || 0}
          </StatNumber>
          <StatLabel fontSize="xs" color="gray.500">Best WPM</StatLabel>
        </Stat>
        <Stat size="sm">
          <StatNumber fontSize="lg" color="white">
            {stats.best_accuracy?.toFixed(0) || 0}%
          </StatNumber>
          <StatLabel fontSize="xs" color="gray.500">Best Accuracy</StatLabel>
        </Stat>
        <Stat size="sm">
          <StatNumber fontSize="lg" color="white">
            {stats.total_sessions || 0}
          </StatNumber>
          <StatLabel fontSize="xs" color="gray.500">Sessions</StatLabel>
        </Stat>
        <Stat size="sm">
          <StatNumber fontSize="lg" color="white">
            {((stats.total_time || 0) / 60).toFixed(0)}m
          </StatNumber>
          <StatLabel fontSize="xs" color="gray.500">Total Time</StatLabel>
        </Stat>
      </SimpleGrid>
    </Box>
  );
}

function SyncDialog() {
  const { showSyncDialog, setShowSyncDialog, pendingSyncData, syncStats, skipSync } = useAuth();
  const [isSyncing, setIsSyncing] = useState(false);
  const toast = useToast();

  const handleMerge = async () => {
    setIsSyncing(true);
    try {
      await syncStats(pendingSyncData);
      toast({
        title: 'Stats synced',
        description: 'Your local stats have been merged with your account',
        status: 'success',
        duration: 3000,
      });
    } catch (e) {
      toast({
        title: 'Sync failed',
        description: e.message,
        status: 'error',
        duration: 3000,
      });
    }
    setIsSyncing(false);
  };

  const handleSkip = () => {
    skipSync();
    toast({
      title: 'Sync skipped',
      description: 'Local stats have been discarded',
      status: 'info',
      duration: 3000,
    });
  };

  return (
    <Modal
      isOpen={showSyncDialog}
      onClose={() => setShowSyncDialog(false)}
      isCentered
      size="lg"
    >
      <ModalOverlay bg="blackAlpha.800" />
      <ModalContent
        bg="bg.card"
        borderRadius="2xl"
        border="1px solid"
        borderColor="whiteAlpha.200"
        shadow="2xl"
      >
        <ModalHeader>
          <VStack spacing={2} align="start">
            <Text fontSize="xl" fontWeight="bold" color="white">
              Sync Your Stats?
            </Text>
            <Text fontSize="sm" color="gray.400" fontWeight="normal">
              You have local stats from before signing in. Would you like to merge them with your account?
            </Text>
          </VStack>
        </ModalHeader>

        <ModalBody>
          <VStack spacing={4}>
            <StatsPreview
              stats={pendingSyncData}
              label="Local Stats (Browser)"
              color="brand.500"
            />

            <HStack spacing={2} color="gray.500">
              <Text fontSize="xs">+</Text>
              <Text fontSize="xs">Merge into your account</Text>
            </HStack>

            <Box
              bg="whiteAlpha.50"
              borderRadius="xl"
              p={4}
              border="1px dashed"
              borderColor="green.500"
              w="100%"
            >
              <Text fontSize="sm" color="green.400" fontWeight="600" mb={2}>
                Merge Strategy
              </Text>
              <VStack align="start" spacing={1}>
                <Text fontSize="xs" color="gray.400">
                  - Best scores: Highest value kept
                </Text>
                <Text fontSize="xs" color="gray.400">
                  - Total sessions: Added together
                </Text>
                <Text fontSize="xs" color="gray.400">
                  - Letter/finger stats: Combined
                </Text>
              </VStack>
            </Box>
          </VStack>
        </ModalBody>

        <ModalFooter>
          <HStack spacing={3} w="100%">
            <Button
              variant="ghost"
              leftIcon={<Icon as={TrashIcon} />}
              onClick={handleSkip}
              flex={1}
              color="gray.400"
              _hover={{ bg: 'whiteAlpha.100', color: 'red.400' }}
            >
              Discard Local
            </Button>
            <Button
              variant="glow"
              leftIcon={<Icon as={MergeIcon} />}
              onClick={handleMerge}
              isLoading={isSyncing}
              loadingText="Syncing..."
              flex={1}
            >
              Merge Stats
            </Button>
          </HStack>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}

export default SyncDialog;
