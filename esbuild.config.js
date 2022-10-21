import { build } from 'esbuild'
import * as dotenv from 'dotenv'
import { dirname } from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

dotenv.config({
  path: __dirname + '/.env',
})
const define = {}

for (const k in process.env) {
  define[`process.env.${k}`] = JSON.stringify(process.env[k])
}

build({
  entryPoints: ['./web/src/index.js'],
  outfile: './web/build/static/app.js',
  minify: true,
  bundle: true,
  loader: {
    '.js': 'jsx',
  },
  plugins: [],
  define,
}).catch(() => process.exit(1))
