import react from "@vitejs/plugin-react-swc";
import { defineConfig } from "vite";
import EnvironmentPlugin from "vite-plugin-environment";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react(), EnvironmentPlugin("all")],
    server: {
        port: 3000,
        host: "0.0.0.0",
        base: "/",
        proxy: {
            "/api/saul": {
                changeOrigin: true,
                target: "http://localhost:8080",
                secure: false,
            },
        },
    },
    resolve: {
        alias: {
            "@root": "/src",
            "@assets": "/src/assets",
            "@components": "/src/components",
            "@hooks": "/src/hooks",
            "@pages": "/src/pages",
            "@services": "/src/services",
        },
    },
    css: {
        postcss: {
            plugins: [require("tailwindcss"), require("autoprefixer")],
        },
    },
});
