/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  reactStrictMode: true,
  swcMinify: true,
  experimental: {
    optimizePackageImports: ['lucide-react', '@radix-ui/react-dialog'],
  },
  // Explicitly disable Turbopack to use webpack
  turbopack: {},
  webpack: (config, { isServer }) => {
    if (!isServer) {
      // Exclude Node.js modules from client bundles
      config.resolve.fallback = {
        ...config.resolve.fallback,
        fs: false,
        path: false,
        os: false,
        crypto: false,
        stream: false,
        util: false,
        buffer: false,
        process: false,
      };
    }
    return config;
  },
}

module.exports = nextConfig

