# frontend

This is the Cloud Snitch frontend.

## Getting Started

- Install Node
- Run `npm install` to install dependencies
- Run `npm run generate` to generate the API code
- Run `ENV=prod.env npm run dev` to run the dev frontend against the prod API

## Docker build

To perform the docker build, the OpenAPI spec must be provided via build arg:

```bash
docker build --build-arg "OPENAPI_YAML=$(cat ../backend/api/apispec/openapi.yaml)" -t cloud-snitch-frontend .
```

## Handy Commands

- `npm run check` - Type check the code
- `npm run lint` - Lint the code
- `npm run lint -- --fix` - Lint the code and automatically fix issues
