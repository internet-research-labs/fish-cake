const webpack = require('webpack');
const UglifyJsPlugin = require('uglifyjs-webpack-plugin');

module.exports = {
  context: __dirname + "/src",
  entry: './index.js',
  output: {
    // path: __dirname + '/dist',
    path: __dirname + '/../server/static/js',
    filename: 'fish-cake.min.js',
    publicPath: "assets/",
    libraryTarget: 'var',
    library: 'svv',
  },
  // devtool: 'source-map',
  plugins: new webpack.optimize.UglifyJsPlugin(),
  module: {
    loaders: [
      {
        test: /\.(frag|vert)?$/,
        exclude: /node_modules/,
        loader: 'webpack-glsl',
      },
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loader: 'babel',
        query: {
          presets: ['es2015', 'react'],
        },
      },
    ],
  },
  externals: {
    "THREE": "THREE",
    "react": "React",
    "react-dom": "ReactDOM",
  },
};
