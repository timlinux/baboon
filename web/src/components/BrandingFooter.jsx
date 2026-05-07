import React from 'react';
import { Flex, HStack, Text, Link } from '@chakra-ui/react';

function BrandingFooter({ versionInfo, fontSize }) {
  const size = fontSize || { base: 'xs', md: 'sm' };
  return (
    <Flex direction="column" align="center" gap={1}>
      <HStack spacing={2} color="gray.600" fontSize={size} fontFamily="'Fira Code', monospace">
        <Text>Made with</Text>
        <Text color="red.400">&#9829;</Text>
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
      {versionInfo?.version && (
        <HStack spacing={2} color="gray.500" fontSize={size} fontFamily="'Fira Code', monospace">
          <Link
            href={`https://github.com/timlinux/baboon/releases/tag/v${versionInfo.version}`}
            isExternal
            color="gray.500"
            _hover={{ color: 'cyan.400' }}
          >
            v{versionInfo.version}
          </Link>
          {versionInfo?.git_commit && versionInfo.git_commit !== 'unknown' && (
            <>
              <Text>&#8226;</Text>
              <Link
                href={`https://github.com/timlinux/baboon/commit/${versionInfo.git_commit}`}
                isExternal
                color="gray.500"
                _hover={{ color: 'cyan.400' }}
              >
                {versionInfo.git_commit}
              </Link>
            </>
          )}
        </HStack>
      )}
    </Flex>
  );
}

export default BrandingFooter;
