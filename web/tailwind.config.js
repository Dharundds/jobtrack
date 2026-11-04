/** @type {import('tailwindcss').Config} */
export default {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Source Sans 3"', "system-ui", "sans-serif"],
        mono: ['"JetBrains Mono"', "ui-monospace", "monospace"],
      },
      boxShadow: {
        board: "0 1px 0 rgb(var(--border)), 0 4px 16px rgb(0 0 0 / 0.12)",
        "board-lg": "0 2px 0 rgb(var(--border)), 0 12px 32px rgb(0 0 0 / 0.18)",
      },
      colors: {
        surface: {
          DEFAULT: "rgb(var(--surface) / <alpha-value>)",
          muted: "rgb(var(--surface-muted) / <alpha-value>)",
          elevated: "rgb(var(--surface-elevated) / <alpha-value>)",
        },
        foreground: {
          DEFAULT: "rgb(var(--foreground) / <alpha-value>)",
          muted: "rgb(var(--foreground-muted) / <alpha-value>)",
        },
        border: {
          DEFAULT: "rgb(var(--border) / <alpha-value>)",
        },
        accent: {
          DEFAULT: "rgb(var(--accent) / <alpha-value>)",
          foreground: "rgb(var(--accent-foreground) / <alpha-value>)",
        },
        rail: {
          DEFAULT: "rgb(var(--rail) / <alpha-value>)",
          foreground: "rgb(var(--rail-foreground) / <alpha-value>)",
          muted: "rgb(var(--rail-muted) / <alpha-value>)",
        },
        status: {
          applied: "rgb(var(--status-applied) / <alpha-value>)",
          replied: "rgb(var(--status-replied) / <alpha-value>)",
          interview: "rgb(var(--status-interview) / <alpha-value>)",
          offer: "rgb(var(--status-offer) / <alpha-value>)",
          rejected: "rgb(var(--status-rejected) / <alpha-value>)",
        },
      },
      borderRadius: {
        md: "0.375rem",
        lg: "0.5rem",
        xl: "0.625rem",
      },
    },
  },
  plugins: [],
};
