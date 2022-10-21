require('dotenv').config()

const define = {}
for (const k in process.env) {
  define[`process.env.${k}`] = JSON.stringify(process.env[k])
}

require('esbuild').build({
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
