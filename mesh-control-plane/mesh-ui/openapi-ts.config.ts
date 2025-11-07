import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
    input: 'http://localhost:8080/swagger/v1/openapiv2.json',
    output: {
        path: './src/api/openapi',
        format: 'prettier',
        lint: 'eslint',
    },
    plugins: ['@hey-api/typescript'],
})