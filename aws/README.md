# aws

## Updating an Environment

Typically an update is done by deploying all the stacks like so:

```bash
AWS_PROFILE=cloud-snitch-dev npx cdk deploy '*-dev'
```

## Initial Deploy

When deploying a new environment for the first time...

You'll need to generate a public/private key pair:

```bash
openssl genrsa -out private_key.pem 2048
openssl rsa -pubout -in private_key.pem -out public_key.pem
```

After deploying the base, but before deploying the regional stacks:

- Set the password encryption key secret. You can generate a new value for it via `openssl rand -base64 32`.
- Set the private signing key secret.
- Set the Stripe secret key secret.

Finally:

- Request that the account be moved out of the SES sandbox in all desired regions.
