/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{js,jsx,tsx,ts}'],
  presets: [require('@ones-design/tailwind-preset')],
  theme: {
    extend: {},
  },
  plugins: [],
  corePlugins: {
    preflight: false, // <== disable this!
  },
};
