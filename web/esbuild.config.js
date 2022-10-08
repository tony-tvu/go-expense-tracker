import { build } from 'esbuild'
import inlineImage from 'esbuild-plugin-inline-image'
import * as dotenv from 'dotenv'
import { dirname } from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)
const parts = __dirname.split('/')

let envPath = ''
for (let i = 1; i < parts.length - 1; i++) {
  envPath += `/${parts[i]}`
}

dotenv.config({
  path: `${envPath}/.env`,
})
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
}).catch(() => process.exit(1))
