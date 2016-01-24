module.exports = {
  entry: './websrc/app/boot.ts',
  output: {
    filename: './app-engine/webapp/app.js'
  },
  devtool: 'source-map',
  resolve: {
    extensions: ['', '.webpack.js', '.web.js', '.ts', '.js']
  },
  module: {
    loaders: [
      { test: /\.ts$/, loader: 'ts-loader' }
    ]
  }
};