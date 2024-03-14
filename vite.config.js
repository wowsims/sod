import path from "path";
import glob from "glob";
import { defineConfig } from 'vite'

export default defineConfig(({ command, mode }) => ({
	base: "/sod/",
	root: path.join(__dirname, "ui"),
	build: {
		outDir: path.join(__dirname, "dist", "sod"),
		minify: mode === "development" ? false : "terser",
		sourcemap: command === "serve" ? "inline" : "false",
		target: ["es2020"],
		rollupOptions: {
			input: {
				...glob.sync(path.resolve(__dirname, "ui", "**/index.html").replace(/\\/g, "/")).reduce((acc, cur) => {
					const name = path.relative(__dirname, cur);
					acc[name] = cur;
					return acc;
				}, {}),
				// Add shared.scss as a separate entry if needed or handle it separately
			},
			output: {
				assetFileNames: () => "bundle/[name]-[hash].style.css",
				entryFileNames: () => "bundle/[name]-[hash].entry.js",
				chunkFileNames: () => "bundle/[name]-[hash].chunk.js",
			},
		},
		server: {
			origin: 'http://localhost:3000',
		},
	}
}));