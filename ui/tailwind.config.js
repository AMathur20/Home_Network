/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'noc-bg': '#0a0a0c',
        'noc-card': '#16161a',
        'noc-cyan': '#00f5ff',
        'noc-emerald': '#10b981',
        'noc-forest': '#065f46',
        'noc-yellow': '#facc15',
        'noc-orange': '#f97316',
      },
      animation: {
        'pulse-glow': 'pulse-glow 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        'pulse-glow': {
          '0%, 100%': { opacity: 1, filter: 'drop-shadow(0 0 5px currentColor)' },
          '50%': { opacity: 0.7, filter: 'drop-shadow(0 0 15px currentColor)' },
        }
      }
    },
  },
  plugins: [],
}
