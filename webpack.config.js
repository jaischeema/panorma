module.exports = {
  context: __dirname + "/client",
  entry: {
    javascript: "./app.js",
    html: "./index.html",
  },

  output: {
    filename: "app.js",
    path: __dirname + "/public",
  },

  module: {
    loaders: [{
      test: /\.js$/,
      exclude: /node_modules/,
      loaders: ["react-hot", "babel-loader"],
    }, {
      test: /\.html$/,
      loader: "file?name=[name].[ext]",
    }],
  },
}
