import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
    assetPrefix: process.env.NEXT_PUBLIC_CDN_URL,
    output: process.env.NODE_ENV === 'development' ? undefined : 'standalone',
    images: process.env.NEXT_PUBLIC_CDN_URL
        ? {
              loader: 'custom',
              loaderFile: './src/image-loader.ts',
          }
        : undefined,
    experimental: {
        turbo: {
            rules: {
                '.svg': {
                    loaders: ['@svgr/webpack'],
                    as: '.js',
                },
            },
        },
    },
};

export default nextConfig;
