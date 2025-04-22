FROM node:23-alpine AS builder

RUN apk add --no-cache libc6-compat openjdk17-jre-headless

WORKDIR /opt/frontend

COPY package.json package-lock.json .
RUN npm ci

ARG OPENAPI_YAML
RUN mkdir -p ../backend/api/apispec && echo "$OPENAPI_YAML" > ../backend/api/apispec/openapi.yaml

COPY . .

ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_CDN_URL
ARG NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY
ARG NEXT_PUBLIC_AWS_ACCOUNT_ID
ARG NEXT_PUBLIC_PUBLIC_S3_BUCKET_NAME
ARG NEXT_PUBLIC_NO_INDEX

RUN npm run generate
RUN npm run build

FROM node:23-alpine AS runner

WORKDIR /opt/frontend
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.9.0 /lambda-adapter /opt/extensions/lambda-adapter

COPY --from=builder /opt/frontend/public ./public

RUN mkdir .next
RUN chown nextjs:nodejs .next
COPY --from=builder --chown=nextjs:nodejs /opt/frontend/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /opt/frontend/.next/static ./.next/static

USER nextjs
EXPOSE 8080
ENV PORT=8080
ENV HOSTNAME=0.0.0.0
CMD ["node", "server.js"]
