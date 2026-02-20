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
    // Kartoza brand colors
    brand: {
      50: '#fef6e9',
      100: '#fce8c7',
      200: '#f9d9a5',
      300: '#f5c983',
      400: '#e8a93d',
      500: '#D4922A', // Kartoza Orange - primary
      600: '#b87a22',
      700: '#9c631a',
      800: '#804c12',
      900: '#64350a',
    },
    // Kartoza Blue palette
    kartoza: {
      blue: {
        50: '#e9f4f7',
        100: '#c7e3ea',
        200: '#a5d2dd',
        300: '#83c1d0',
        400: '#61b0c3',
        500: '#4A90A4', // Kartoza Blue - secondary
        600: '#3d7688',
        700: '#305c6c',
        800: '#234250',
        900: '#162834',
      },
      orange: '#D4922A',
      gray: {
        light: '#C4C4C4',
        medium: '#9A9A9A',
        dark: '#6A6A6A',
      },
    },
    bg: {
      primary: '#1a2833',
      secondary: '#243442',
      tertiary: '#2e4050',
      card: '#1f3040',
    },
    accent: {
      green: '#4CAF50',
      red: '#E53935',
      yellow: '#FFC107',
      cyan: '#4A90A4', // Kartoza blue as cyan
      purple: '#7E57C2',
      orange: '#D4922A', // Kartoza orange
    },
    gray: {
      750: '#3a4a5a',
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
          bg: 'brand.500', // Kartoza orange
          color: 'white',
          boxShadow: '0 0 30px rgba(212, 146, 42, 0.4)',
          _hover: {
            bg: 'kartoza.blue.500', // Kartoza blue on hover
            boxShadow: '0 0 50px rgba(74, 144, 164, 0.6)',
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
