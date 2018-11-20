const webpack = require('webpack');

module.exports = {
  context: __dirname + "/src",
  entry: './index.js',
  output: {
    // path: __dirname + '/dist',
    path: __dirname + '/../server/static/js',
    publicPath: "assets/",
    filename: 'fish-cake.min.js',
    libraryTarget: 'var',
    library: 'svv',
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: ['babel-loader'],
      }
    ],
  },
  resolve: {
    extensions: ['*', '.js'],
  },
  externals: {
    "THREE": "THREE",
    "react": "React",
    "react-dom": "ReactDOM",
  },
  mode: "production",
};
