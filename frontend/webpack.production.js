/* eslint-disable */
const { merge } = require("webpack-merge");
const webpack = require('webpack');
const dotenv = require('dotenv');
const WebpackObfuscator = require("webpack-obfuscator");
const common = require("./webpack.common.js");

module.exports = merge(common, {
    mode: "production",
    plugins: [
        new webpack.DefinePlugin({
            'process.env.BACKEND_GATEWAY': JSON.stringify(process.env.BACKEND_GATEWAY),
            'process.env.BACKEND_GATEWAY_WS': JSON.stringify(process.env.BACKEND_GATEWAY_WS),
            'process.env.DOC_SERVER': JSON.stringify(process.env.DOC_SERVER),
            'process.env.WORD_FILE': JSON.stringify(process.env.WORD_FILE),
            'process.env.SLIDE_FILE': JSON.stringify(process.env.SLIDE_FILE),
            'process.env.SPREADSHEET_FILE': JSON.stringify(process.env.SPREADSHEET_FILE),
            'process.env.FILE_STALE_TIME': JSON.stringify(process.env.FILE_STALE_TIME),
            'process.env.FILE_CACHE_TIME': JSON.stringify(process.env.FILE_CACHE_TIME),
        }),
        // new WebpackObfuscator({
        //     rotateStringArray: true,
        // }),
    ],
});
