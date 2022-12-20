/* eslint-disable */
const { merge } = require("webpack-merge");
const WebpackObfuscator = require("webpack-obfuscator");
const common = require("./webpack.common.js");

module.exports = merge(common, {
    mode: "production",
    plugins: [
        new WebpackObfuscator({
            rotateStringArray: true,
        }),
    ],
});
