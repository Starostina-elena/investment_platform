import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    output: "standalone",
    images: {
        unoptimized: true
    },
    trailingSlash: true,
    async rewrites() {
        return [
            {
                source: '/api/payment/:path*',
                destination: 'http://payment:8106/:path*'
            },
            {
                source: '/api/:path*',
                destination: 'http://nginx:80/api/:path*'
            },
            {
                source: '/media/:path*',
                destination: 'http://nginx:80/media/:path*'
            }
        ]
    }
};

export default nextConfig;
