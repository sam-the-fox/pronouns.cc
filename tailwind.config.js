module.exports = {
  darkMode: "class",
  content: ["./frontend/index.html", "./frontend/src/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("@tailwindcss/forms")],
};
