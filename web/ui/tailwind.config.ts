import { config } from "@qttdev/ui/tailwind";
import type { Config } from "tailwindcss";

export default {
  presets:[config],
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
      },
      borderRadius: {
        "4xl": "2rem",
      }
    },
  },
  plugins: [],
} satisfies Config;
