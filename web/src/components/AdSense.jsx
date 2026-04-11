import React, { useEffect, useRef, useState } from 'react';
import { Box, Text, Flex, Skeleton } from '@chakra-ui/react';

/**
 * Google AdSense component that displays an ad unit.
 * Shows a preview placeholder when no publisherId is provided.
 * Set showPreview={true} to always show the ad area (with placeholder if no publisherId).
 */
function AdSense({ publisherId, slot = 'auto', format = 'auto', responsive = true, showPreview = true }) {
  const adRef = useRef(null);
  const isLoaded = useRef(false);
  const [isLoading, setIsLoading] = useState(!!publisherId);

  useEffect(() => {
    // Only load if we have a publisher ID and haven't loaded yet
    if (!publisherId || isLoaded.current) return;

    setIsLoading(true);

    // Check if the AdSense script is already loaded
    const existingScript = document.querySelector(
      'script[src*="pagead2.googlesyndication.com"]'
    );

    if (!existingScript) {
      // Load the AdSense script
      const script = document.createElement('script');
      script.src = `https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=${publisherId}`;
      script.async = true;
      script.crossOrigin = 'anonymous';
      document.head.appendChild(script);

      script.onload = () => {
        pushAd();
        setIsLoading(false);
      };

      script.onerror = () => {
        setIsLoading(false);
      };
    } else {
      pushAd();
      setIsLoading(false);
    }

    isLoaded.current = true;
  }, [publisherId]);

  const pushAd = () => {
    try {
      if (window.adsbygoogle && adRef.current) {
        (window.adsbygoogle = window.adsbygoogle || []).push({});
      }
    } catch (e) {
      console.error('AdSense error:', e);
    }
  };

  // Don't render if no publisher ID and preview is disabled
  if (!publisherId && !showPreview) {
    return null;
  }

  return (
    <Box
      w="100%"
      maxW="728px"
      mx="auto"
      py={4}
      textAlign="center"
    >
      <Text fontSize="xs" color="gray.600" mb={2}>
        Advertisement
      </Text>
      <Box
        bg="bg.tertiary"
        borderRadius="lg"
        overflow="hidden"
        minH="90px"
        border="1px dashed"
        borderColor={publisherId ? 'transparent' : 'gray.600'}
      >
        {publisherId ? (
          isLoading ? (
            <Skeleton height="90px" borderRadius="lg" />
          ) : (
            <ins
              ref={adRef}
              className="adsbygoogle"
              style={{
                display: 'block',
                minWidth: '300px',
                minHeight: '90px',
              }}
              data-ad-client={publisherId}
              data-ad-slot={slot}
              data-ad-format={format}
              data-full-width-responsive={responsive ? 'true' : 'false'}
            />
          )
        ) : (
          <Flex
            align="center"
            justify="center"
            minH="90px"
            color="gray.500"
            fontSize="sm"
            flexDirection="column"
            gap={1}
          >
            <Text>Ad Preview Placeholder</Text>
            <Text fontSize="xs" color="gray.600">
              728x90 Leaderboard
            </Text>
          </Flex>
        )}
      </Box>
    </Box>
  );
}

export default AdSense;
