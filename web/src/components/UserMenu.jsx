import React from 'react';
import {
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  MenuDivider,
  Button,
  Avatar,
  HStack,
  VStack,
  Text,
  Icon,
  useToast,
} from '@chakra-ui/react';
import { useAuth } from '../contexts/AuthContext.jsx';

// Simple icons
const LogoutIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4" />
    <polyline points="16,17 21,12 16,7" />
    <line x1="21" y1="12" x2="9" y2="12" />
  </svg>
);

const ExportIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" />
    <polyline points="7,10 12,15 17,10" />
    <line x1="12" y1="15" x2="12" y2="3" />
  </svg>
);

const SyncIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <polyline points="23,4 23,10 17,10" />
    <polyline points="1,20 1,14 7,14" />
    <path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" />
  </svg>
);

const DeleteIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <polyline points="3,6 5,6 21,6" />
    <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
    <line x1="10" y1="11" x2="10" y2="17" />
    <line x1="14" y1="11" x2="14" y2="17" />
  </svg>
);

function UserMenu() {
  const { user, logout, isAuthenticated, setShowSyncDialog, localStats } = useAuth();
  const toast = useToast();

  if (!isAuthenticated || !user) {
    return null;
  }

  const handleExport = async () => {
    try {
      const response = await fetch('/api/user/export', {
        credentials: 'include',
      });
      if (!response.ok) throw new Error('Export failed');

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'baboon-data-export.json';
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);

      toast({
        title: 'Data exported',
        description: 'Your data has been downloaded',
        status: 'success',
        duration: 3000,
      });
    } catch (e) {
      toast({
        title: 'Export failed',
        description: e.message,
        status: 'error',
        duration: 3000,
      });
    }
  };

  const handleSync = () => {
    if (localStats && localStats.total_sessions > 0) {
      setShowSyncDialog(true);
    } else {
      toast({
        title: 'Nothing to sync',
        description: 'No local stats found',
        status: 'info',
        duration: 3000,
      });
    }
  };

  const handleLogout = async () => {
    await logout();
    toast({
      title: 'Logged out',
      description: 'You have been logged out successfully',
      status: 'success',
      duration: 3000,
    });
  };

  return (
    <Menu>
      <MenuButton
        as={Button}
        variant="ghost"
        rounded="full"
        p={1}
        _hover={{ bg: 'whiteAlpha.100' }}
      >
        <HStack spacing={2}>
          <Avatar
            size="sm"
            name={user.display_name || user.email}
            src={user.avatar_url}
            bg="brand.500"
          />
          <Text color="gray.300" fontSize="sm" display={{ base: 'none', md: 'block' }}>
            {user.display_name || user.email.split('@')[0]}
          </Text>
        </HStack>
      </MenuButton>

      <MenuList bg="bg.card" borderColor="whiteAlpha.200" shadow="xl">
        <VStack px={4} py={3} spacing={1} align="start">
          <Text fontWeight="600" color="white">
            {user.display_name || 'Anonymous'}
          </Text>
          <Text fontSize="sm" color="gray.400">
            {user.email}
          </Text>
          <Text fontSize="xs" color="gray.500" textTransform="capitalize">
            via {user.provider}
          </Text>
        </VStack>

        <MenuDivider borderColor="whiteAlpha.200" />

        {localStats && localStats.total_sessions > 0 && (
          <MenuItem
            icon={<Icon as={SyncIcon} />}
            onClick={handleSync}
            _hover={{ bg: 'whiteAlpha.100' }}
          >
            Sync Local Stats ({localStats.total_sessions} sessions)
          </MenuItem>
        )}

        <MenuItem
          icon={<Icon as={ExportIcon} />}
          onClick={handleExport}
          _hover={{ bg: 'whiteAlpha.100' }}
        >
          Export My Data
        </MenuItem>

        <MenuDivider borderColor="whiteAlpha.200" />

        <MenuItem
          icon={<Icon as={LogoutIcon} />}
          onClick={handleLogout}
          _hover={{ bg: 'whiteAlpha.100' }}
          color="red.400"
        >
          Sign Out
        </MenuItem>
      </MenuList>
    </Menu>
  );
}

export default UserMenu;
