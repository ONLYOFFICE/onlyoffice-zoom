/* eslint-disable */
const { merge } = require("webpack-merge");
const common = require("./webpack.common.js");
const Dotenv = require("dotenv-webpack");

module.exports = merge(common, {
    mode: "development",
    devtool: "inline-source-map",
    devServer: {
        historyApiFallback: true,
        client: {
            logging: "info",
            overlay: true,
        },
        compress: true,
        open: true,
        static: "./build",
        port: 3000,
        allowedHosts: "all",
    },
    stats: {
        errorDetails: true,
    },
    plugins: [
        new Dotenv(),
    ],
});
