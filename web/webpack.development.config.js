const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
    mode: 'development',
    entry: {
        login: './app/pages/login.js',
        registerAdmin: './app/pages/registerAdmin.js',
        lists: './app/pages/lists.js',
        listCreate: './app/pages/listCreate.js',
        listEdit: './app/pages/listEdit.js',
        listItems: './app/pages/listItems.js',
        itemCreate: './app/pages/itemCreate.js',
    },
    output: {
        filename: '[name].js',
        path: __dirname + '/dist',
        clean: true
    },
    module: {
      rules: [
        {
          test: /\.css$/,
          exclude: /node_modules/,
          use: [
            MiniCssExtractPlugin.loader,
            'css-loader',
            'postcss-loader',
          ],
        }
      ]
    },
    plugins: [
      new MiniCssExtractPlugin({
        filename: 'styles.css',
      }),
      new CopyWebpackPlugin({
        patterns: [
          { from: 'app/**/*.html', to: '[name][ext]' },
        ],
      }),
    ],
    optimization: {
      minimize: false,
      splitChunks: {
        chunks: 'all',
      },
    },
  }