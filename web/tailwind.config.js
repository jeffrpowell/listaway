/** @type {import('tailwindcss').Config} */
module.exports = {
  relative: true,
  content: ["./app/**/*.{html,js}"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

