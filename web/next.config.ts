import type { NextConfig } from 'next'
import type { RuleSetRule } from 'webpack'

const nextConfig: NextConfig = {
  experimental: {
    serverActions: {
      bodySizeLimit: '5mb',
    },
  },
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**', // TODO TEMPORARY: restrict this later
      },
    ],
  },
  typedRoutes: true,
  webpack(config) {
    const fileLoaderRule = config.module.rules.find(
      (rule: RuleSetRule) => rule.test instanceof RegExp && rule.test.test('.svg'),
    )

    config.module.rules.push(
      {
        ...fileLoaderRule,
        test: /\.svg$/i,
        resourceQuery: /url/, // import icon from './icon.svg?url'
      },

      {
        test: /\.svg$/i,
        issuer: fileLoaderRule.issuer,
        resourceQuery: { not: [...(fileLoaderRule.resourceQuery?.not || []), /url/] },
        use: [
          {
            loader: '@svgr/webpack',
            options: {
              icon: true,
            },
          },
        ],
      },
    )

    fileLoaderRule.exclude = /\.svg$/i

    return config
  },
}

export default nextConfig
