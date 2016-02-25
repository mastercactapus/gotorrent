"use strict";

var webpack = require("webpack");

module.exports = {
    entry: "./src/main",
    output: {
      path: __dirname + "/dist",
            filename: 'bundle.js'
    },
    plugins: [
      // new webpack.HotModuleReplacementPlugin(),
      // new webpack.NoErrorsPlugin(),
    ],
    module: {
      loaders: [
        { test: /\.js$/, exclude: /node_modules/, loader: "babel" }
      ]
    }
}
