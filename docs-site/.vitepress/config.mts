import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "pronouns.cc documentation",
  description: "pronouns.cc documentation",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Back to pronouns.cc", link: "https://pronouns.cc/" },
    ],

    sidebar: [
      {
        text: "API",
        items: [
          { text: "API reference", link: "/api/" },
          { text: "Rate limits", link: "/api/ratelimits" },
          { text: "Error messages", link: "/api/errors" },
        ],
      },
      {
        text: "Endpoints",
        items: [{ text: "Users", link: "/api/endpoints/users" }],
      },
    ],
  },
});
