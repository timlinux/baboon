import { extendTheme } from '@chakra-ui/react';

const theme = extendTheme({
  config: {
    initialColorMode: 'dark',
    useSystemColorMode: false,
  },
  fonts: {
    heading: '"Inter", sans-serif',
    body: '"Inter", sans-serif',
    mono: '"JetBrains Mono", monospace',
  },
  colors: {
    brand: {
      50: '#e6f7ff',
      100: '#b3e0ff',
      200: '#80caff',
      300: '#4db3ff',
      400: '#1a9dff',
      500: '#0080e6',
      600: '#0066b3',
      700: '#004d80',
      800: '#00334d',
      900: '#001a26',
    },
    bg: {
      primary: '#0d0d1a',
      secondary: '#1a1a2e',
      tertiary: '#252542',
      card: '#16162a',
    },
    accent: {
      green: '#00ff88',
      red: '#ff4466',
      yellow: '#ffcc00',
      cyan: '#00ccff',
      purple: '#aa66ff',
      orange: '#ff8844',
    },
    gray: {
      750: '#2d2d44',
    },
  },
  styles: {
    global: {
      body: {
        bg: 'bg.primary',
        color: 'white',
      },
    },
  },
  components: {
    Button: {
      baseStyle: {
        fontWeight: 'bold',
        borderRadius: '2xl',
        transition: 'all 0.3s ease',
      },
      sizes: {
        chunky: {
          h: '80px',
          minW: '200px',
          fontSize: '2xl',
          px: '12',
        },
        huge: {
          h: '100px',
          minW: '300px',
          fontSize: '3xl',
          px: '16',
        },
      },
      variants: {
        glow: {
          bg: 'accent.cyan',
          color: 'bg.primary',
          boxShadow: '0 0 30px rgba(0, 204, 255, 0.4)',
          _hover: {
            bg: 'accent.green',
            boxShadow: '0 0 50px rgba(0, 255, 136, 0.6)',
            transform: 'scale(1.05)',
          },
          _active: {
            transform: 'scale(0.98)',
          },
        },
        ghost: {
          color: 'gray.400',
          _hover: {
            bg: 'whiteAlpha.100',
            color: 'white',
          },
        },
      },
    },
    Card: {
      baseStyle: {
        container: {
          bg: 'bg.card',
          borderRadius: '3xl',
          border: '1px solid',
          borderColor: 'whiteAlpha.100',
        },
      },
    },
  },
});

export default theme;
