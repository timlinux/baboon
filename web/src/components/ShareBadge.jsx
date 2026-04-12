import React, { useState, useCallback } from 'react';
import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalCloseButton,
  VStack,
  HStack,
  Text,
  Box,
  Image,
  useToast,
} from '@chakra-ui/react';
import { motion } from 'framer-motion';
import api from '../api.js';

const MotionBox = motion(Box);

function ShareBadge({ entry, onClose }) {
  const toast = useToast();
  const [isCopying, setIsCopying] = useState(false);

  const badgeUrl = api.getBadgeUrl(entry.id);
  const shareUrl = `${window.location.origin}/leaderboard?highlight=${entry.id}`;
  const shareText = `I ranked #${entry.rank} on the Baboon typing leaderboard with ${entry.wpm.toFixed(1)} WPM! Can you beat my score?`;
  const linkedInUrl = `https://www.linkedin.com/sharing/share-offsite/?url=${encodeURIComponent(shareUrl)}`;
  const mastodonText = `${shareText}\n\n${shareUrl}`;
  const redditUrl = `https://www.reddit.com/submit?url=${encodeURIComponent(shareUrl)}&title=${encodeURIComponent(shareText)}`;
  const emailUrl = `mailto:?subject=${encodeURIComponent('Check out my Baboon typing score!')}&body=${encodeURIComponent(`${shareText}\n\n${shareUrl}`)}`;

  const handleCopyLink = useCallback(async () => {
    setIsCopying(true);
    try {
      await navigator.clipboard.writeText(shareUrl);
      toast({
        title: 'Link copied!',
        status: 'success',
        duration: 2000,
        isClosable: true,
      });
    } catch (e) {
      toast({
        title: 'Failed to copy',
        description: 'Please copy the link manually',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
    setIsCopying(false);
  }, [shareUrl, toast]);

  const handleDownload = useCallback(() => {
    // Create a link element and trigger download
    const link = document.createElement('a');
    link.href = badgeUrl;
    link.download = `baboon-badge-${entry.display_name.toLowerCase().replace(/\s/g, '-')}.svg`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    toast({
      title: 'Badge downloaded!',
      status: 'success',
      duration: 2000,
      isClosable: true,
    });
  }, [badgeUrl, entry.display_name, toast]);

  const handleShareLinkedIn = useCallback(() => {
    window.open(linkedInUrl, '_blank', 'width=550,height=420');
  }, [linkedInUrl]);

  const handleShareMastodon = useCallback(() => {
    // Mastodon share - user can paste into their instance
    navigator.clipboard.writeText(mastodonText).then(() => {
      toast({
        title: 'Copied for Mastodon!',
        description: 'Paste into your Mastodon instance',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    });
  }, [mastodonText, toast]);

  const handleShareReddit = useCallback(() => {
    window.open(redditUrl, '_blank', 'width=550,height=420');
  }, [redditUrl]);

  const handleShareEmail = useCallback(() => {
    window.location.href = emailUrl;
  }, [emailUrl]);

  return (
    <Modal isOpen onClose={onClose} isCentered size="md">
      <ModalOverlay bg="blackAlpha.800" />
      <ModalContent
        bg="#1a2833"
        border="2px solid"
        borderColor="#D4922A"
        borderRadius="lg"
      >
        <ModalHeader
          fontFamily="'Fira Code', monospace"
          color="cyan.400"
          textAlign="center"
        >
          Share Your Score
        </ModalHeader>
        <ModalCloseButton color="gray.400" />

        <ModalBody pb={6}>
          <VStack spacing={6}>
            {/* Badge preview */}
            <MotionBox
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ type: 'spring' }}
              w="100%"
              maxW="400px"
              borderRadius="lg"
              overflow="hidden"
              border="1px solid"
              borderColor="gray.700"
            >
              <Image
                src={badgeUrl}
                alt={`${entry.display_name}'s Baboon Badge`}
                w="100%"
                fallbackSrc="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 400 200'%3E%3Crect fill='%231a2833' width='400' height='200'/%3E%3Ctext fill='%23808080' x='200' y='100' text-anchor='middle' font-family='monospace'%3ELoading badge...%3C/text%3E%3C/svg%3E"
              />
            </MotionBox>

            {/* Player stats */}
            <HStack spacing={4} fontFamily="'Fira Code', monospace">
              <Text color="#00ff00" fontWeight="bold">
                {entry.wpm.toFixed(1)} WPM
              </Text>
              <Text color="cyan.400">
                {entry.accuracy.toFixed(1)}%
              </Text>
              <Text color="#D4922A">
                Rank #{entry.rank}
              </Text>
            </HStack>

            {/* Action buttons */}
            <HStack spacing={4} w="100%">
              <MotionBox
                flex={1}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <Box
                  as="button"
                  onClick={handleCopyLink}
                  w="100%"
                  py={2}
                  borderRadius="md"
                  bg="transparent"
                  border="2px solid"
                  borderColor="cyan.500"
                  color="cyan.400"
                  fontFamily="'Fira Code', monospace"
                  fontSize="sm"
                  _hover={{ bg: 'rgba(74, 144, 164, 0.1)' }}
                  disabled={isCopying}
                >
                  {isCopying ? 'Copying...' : 'Copy Link'}
                </Box>
              </MotionBox>

              <MotionBox
                flex={1}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <Box
                  as="button"
                  onClick={handleDownload}
                  w="100%"
                  py={2}
                  borderRadius="md"
                  bg="transparent"
                  border="2px solid"
                  borderColor="#00ff00"
                  color="#00ff00"
                  fontFamily="'Fira Code', monospace"
                  fontSize="sm"
                  _hover={{ bg: 'rgba(0, 255, 0, 0.1)' }}
                >
                  Download
                </Box>
              </MotionBox>
            </HStack>

            {/* Social share buttons */}
            <HStack spacing={2} w="100%">
              <MotionBox
                flex={1}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <Box
                  as="button"
                  onClick={handleShareLinkedIn}
                  w="100%"
                  py={2}
                  borderRadius="md"
                  bg="transparent"
                  border="2px solid"
                  borderColor="#0077B5"
                  color="#0077B5"
                  fontFamily="'Fira Code', monospace"
                  fontSize="xs"
                  _hover={{ bg: 'rgba(0, 119, 181, 0.1)' }}
                >
                  LinkedIn
                </Box>
              </MotionBox>

              <MotionBox
                flex={1}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <Box
                  as="button"
                  onClick={handleShareMastodon}
                  w="100%"
                  py={2}
                  borderRadius="md"
                  bg="transparent"
                  border="2px solid"
                  borderColor="#6364FF"
                  color="#6364FF"
                  fontFamily="'Fira Code', monospace"
                  fontSize="xs"
                  _hover={{ bg: 'rgba(99, 100, 255, 0.1)' }}
                >
                  Mastodon
                </Box>
              </MotionBox>

              <MotionBox
                flex={1}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <Box
                  as="button"
                  onClick={handleShareReddit}
                  w="100%"
                  py={2}
                  borderRadius="md"
                  bg="transparent"
                  border="2px solid"
                  borderColor="#FF4500"
                  color="#FF4500"
                  fontFamily="'Fira Code', monospace"
                  fontSize="xs"
                  _hover={{ bg: 'rgba(255, 69, 0, 0.1)' }}
                >
                  Reddit
                </Box>
              </MotionBox>

              <MotionBox
                flex={1}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <Box
                  as="button"
                  onClick={handleShareEmail}
                  w="100%"
                  py={2}
                  borderRadius="md"
                  bg="transparent"
                  border="2px solid"
                  borderColor="gray.500"
                  color="gray.300"
                  fontFamily="'Fira Code', monospace"
                  fontSize="xs"
                  _hover={{ bg: 'rgba(255, 255, 255, 0.05)' }}
                >
                  Email
                </Box>
              </MotionBox>
            </HStack>

            {/* Share URL display */}
            <Box
              w="100%"
              p={4}
              bg="gray.800"
              borderRadius="md"
              border="2px solid"
              borderColor="cyan.500"
            >
              <VStack spacing={2}>
                <Text
                  fontFamily="'Fira Code', monospace"
                  fontSize="sm"
                  color="cyan.400"
                  fontWeight="bold"
                >
                  YOUR BADGE LINK
                </Text>
                <Text
                  fontFamily="'Fira Code', monospace"
                  fontSize="md"
                  color="white"
                  wordBreak="break-all"
                  textAlign="center"
                >
                  {shareUrl}
                </Text>
                <Text
                  fontFamily="'Fira Code', monospace"
                  fontSize="xs"
                  color="gray.400"
                >
                  Click "Copy Link" above to share!
                </Text>
              </VStack>
            </Box>
          </VStack>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
}

export default ShareBadge;
