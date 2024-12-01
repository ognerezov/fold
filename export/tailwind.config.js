/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [],
  theme: {
    extend: {
      colors: {
        background: "var(--background)",
        foreground: "var(--foreground)",
        transparent: 'transparent',
        current: 'currentColor',
        paper: "#FCFBF7",
        deep: "#0C2D50",
        fire: "#CB3A27",
        volcano: "#CD5C58",
        wave: "#61A0A9",
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms')({
      strategy: 'class'
    }),
  ],
}

