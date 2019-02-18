const path    = require('path')
const webpack = require('webpack');

module.exports = {
  entry: './js-lib/bundle.js',
  output: {
    filename: 'webdaemon.js',
    path:     path.resolve(__dirname, 'js-dist')
  },

  target: "web",

  plugins: [
    new webpack.DefinePlugin({
      __DEV__: false,
      __PROD__: true,
      'process.env':{
        'NODE_ENV': JSON.stringify('production')
      },
    })
  ]
}