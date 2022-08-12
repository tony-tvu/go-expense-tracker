const path = require("path")
const HtmlWebpackPlugin = require("html-webpack-plugin")

module.exports = {
  entry: "./web/src/index.js",
  output: {
    path: path.join(__dirname, "/build/static"),
    filename: "bundle.[contenthash].js",
    clean: true,
  },
  devtool: "source-map",
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: "babel-loader",
        },
      },
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      },
    ],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: "web/index.html",
      filename: "index.html",
      // favicon: "public/favicon.ico",
    }),
  ],
}
