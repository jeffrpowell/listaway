const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
    mode: 'development',
    entry: {
      itemCreate: './app/pages/itemCreate.js',
      listCreate: './app/pages/listCreate.js',
      listEdit: './app/pages/listEdit.js',
      listItems: './app/pages/listItems.js',
      lists: './app/pages/lists.js',
      registerAdmin: './app/pages/registerAdmin.js',
      login: './app/pages/login.js',
      sharedList: './app/pages/sharedList.js',
      sharedList404: './app/pages/sharedList404.js',
      userAdmin: './app/pages/userAdmin.js',
      userCreate: './app/pages/userCreate.js',
      allUsers: './app/pages/allUsers.js',
      resetForm: './app/pages/resetForm.js',
      collectionCreate: './app/pages/collectionCreate.js',
      collectionDetail: './app/pages/collectionDetail.js',
      collectionEdit: './app/pages/collectionEdit.js',
      sharedCollection: './app/pages/sharedCollection.js',
      sharedCollection404: './app/pages/sharedCollection404.js',
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
        },
        {
          test: /\.css$/,
          include: /node_modules/,
          use: [
              MiniCssExtractPlugin.loader,
              'css-loader',
          ],
        },
        {
          test: /\.(woff|woff2|eot|ttf|otf)$/,
          type: 'asset/resource',
          generator: {
              filename: 'fonts/[name][ext]',
          },
        },
      ]
    },
    plugins: [
      new MiniCssExtractPlugin({
        filename: '[name].css',
      }),
      new CopyWebpackPlugin({
        patterns: [
          { from: 'app/**/*.html', to: '[name][ext]' },
          { from: 'app/**/*.png', to: '[name][ext]' },
        ],
      }),
    ],
    optimization: {
      minimize: false,
      splitChunks: {
        chunks: 'all',
        cacheGroups: {
          default: false,
          vendors: {
            test: /[\\/]node_modules[\\/]/,
            name: 'vendors',
            chunks: 'all',
            enforce: true,
          },
        },
      },
    },
  }