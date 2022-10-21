require('dotenv').config()
const appENVs = [
  'REACT_APP_ENV',
  'REACT_APP_API_URL',
  'REACT_APP_TELLER_ENV',
  'REACT_APP_TELLER_APPLICATION_ID',
]

const define = {}
for (const k in process.env) {
  if (appENVs.includes(k)) {
    define[`process.env.${k}`] = JSON.stringify(process.env[k])
  }
}

require('esbuild')
  .build({
    entryPoints: ['./web/src/index.js'],
    outfile: './web/build/static/app.js',
    minify: true,
    bundle: true,
    loader: {
      '.js': 'jsx',
    },
    plugins: [],
    define,
  })
  .catch(() => process.exit(1))
