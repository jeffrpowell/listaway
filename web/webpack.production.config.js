const path = require('path');
const TerserPlugin = require('terser-webpack-plugin');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');
const { optimization } = require('./webpack.development.config');

module.exports = {
  extends: path.resolve(__dirname, './webpack.development.config.js'),
  mode: 'production',
  optimization: {
    minimize: true,
    minimizer: [
      new TerserPlugin({
        terserOptions: {
          compress: {
            drop_console: true, // Removes console logs
          },
        },
      }),
      new CssMinimizerPlugin(),
    ],
  }
}