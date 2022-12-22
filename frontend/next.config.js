const { withSentryConfig } = require("@sentry/nextjs");

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:8080/:path*", // proxy to backend
      },
      {
        source: "/media/:path*",
        destination: "http://localhost:9000/pronouns.cc/:path*", // proxy to media server
      }
    ];
  },
  sentry: {
    hideSourceMaps: true,
  },
};

const sentryWebpackPluginOptions = {
  silent: true, // Suppresses all logs
};

module.exports = withSentryConfig(nextConfig, sentryWebpackPluginOptions);
