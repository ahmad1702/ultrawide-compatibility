/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./views/**/*.html', './router/**/*.go'],
  theme: {
    extend: {}
  },
  plugins: [
    require('daisyui'),
    require("tailwindcss-animate")
  ],
  safelist: [
    'table'
  ],
  daisyui: {
  }
}
