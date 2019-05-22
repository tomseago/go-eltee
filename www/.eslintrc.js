module.exports = {
    root: true,
    extends: "airbnb",

    // Plugings
    plugins: [
        "react",
        "babel",
    ],

    parser: "babel-eslint",
    
    rules: {
        "indent": ["warn", 4],
        "quotes": ["warn", "double"],
        "no-trailing-spaces": 0,

        // For node.js we want to use the console
        "no-console": 0,

        "no-multiple-empty-lines": 0,
        "no-underscore-dangle": 0,

        // Destroys proper structured code
        "no-else-return": 0,
        "class-methods-use-this": 0,

        //"max-len": [2, 120],
        "max-len": 0,

        // OMG these are annoying
        "react/jsx-indent": 0,
        "react/jsx-indent-props": 0,
    },

    env: {
        node: true
    },
};