import { build } from 'esbuild'
import inlineImage from 'esbuild-plugin-inline-image'
import * as dotenv from 'dotenv'

dotenv.config()
const define = {}

for (const k in process.env) {
  define[`process.env.${k}`] = JSON.stringify(process.env[k])
}

build({
    entryPoints: ['./src/index.js'],
    outfile: './build/static/app.js',
    minify: true,
    bundle: true,
    loader: {
      '.js': 'jsx',
    },
    plugins: [inlineImage()],
    define,
  })
  .catch(() => process.exit(1))
