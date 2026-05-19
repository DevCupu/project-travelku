/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#E31B23',
          hover: '#B8141A',
          active: '#B8141A',
          disabled: '#F7C1C3',
        },
        accent: {
          blue: '#0056B3',
          soft: '#E6F0FA',
        }
      },
      borderRadius: {
        'lg': '24px',
        'md': '16px',
        'sm': '8px',
      }
    },
  },
  plugins: [],
}
