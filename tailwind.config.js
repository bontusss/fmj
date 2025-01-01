/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        '**/*.{html,templ}',
        
        './node_modules/preline/dist/*.js',
        
    ],
    theme: {
        extend: {},
    },
    // darkMode: "media",
    plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/typography'),
        
        require('preline/plugin'),
        
    ],
}
