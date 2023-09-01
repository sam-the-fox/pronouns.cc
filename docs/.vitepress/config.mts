import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "pronouns.cc documentation",
  description: "pronouns.cc documentation",
  markdown: {
    anchor: { level: [2, 3] },
  },
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    siteTitle: "pronouns.cc",
    logo: "/logo.svg",
    nav: [
      { text: "Home", link: "/" },
      { text: "Back to pronouns.cc", link: "https://pronouns.cc/" },
    ],
    outline: {
      level: [2, 3],
    },
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
        items: [
          { text: "Object reference", link: "/api/endpoints/" },
          { text: "Users", link: "/api/endpoints/users" },
          { text: "Members", link: "/api/endpoints/members" },
        ],
      },
    ],
  },
});
