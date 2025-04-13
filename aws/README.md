# aws

This directory contains a production-ready AWS deployment for Cloud Snitch. It features...

- Zero-downtime deployments - Deployments happen without disrupting user activity. Users get the latest version of the app whenever they refresh.
- An active-active multi-region architecture - Critical functionality is uninterrupted by full region outages.
- Scaling to (near) zero – An idle deployment costs just a few dollars a month.
- One-line deployments – Deployments are as simple as running a command like `npx cdk deploy '*-dev'`. The CDK does all the heavy lifting.

## Deploying for the First Time

### Preparation

Make sure you have Node.js and npm installed.

Generate a public/private key pair for URL signing:

```bash
openssl genrsa -out private_key.pem 2048
openssl rsa -pubout -in private_key.pem -out public_key.pem
```

Generate a password encryption key secret:

```bash
openssl rand -base64 32
```

Create a Stripe account with two products: one for the individual plan and one for the team plan. The products should be associated with features granting the "individual-features" and "team-features" entitlements respectively. If you don't intend to charge for the service, you can use Stripe's test mode or assign $0.00 prices to the products. If your price should be based on the number of active AWS accounts the team is using, add a "use_account_quantity" metadata key with a value of "true" to the price.

Lastly, you'll need an AWS account to deploy to. **It's strongly recommended that you use a new AWS account with no other resources in it.**

### Configuring the Environment

Modify bin/aws.ts to configure an environment with all the required parameters.

### Deploying the Base Stack

Deploy the base stack only, with a command like...

```bash
AWS_PROFILE=cloud-snitch-dev npx cdk deploy global-base-dev
```

After the deployment, there will be three secrets in AWS which must be set before continuing:

- The password encryption key secret.
- The private signing key secret.
- The Stripe secret key secret.

Set these to the values generated previously.

### Deploying the Remaining Stacks

You can deploy the remaining stacks with:

```bash
AWS_PROFILE=cloud-snitch-dev npx cdk deploy '*-dev'
```

### Email

Lastly, to send emails, you'll need to request that the account be moved out of the SES sandbox in all desired regions.

At this point, your deployment should be fully functional.

## Updating an Environment

Typically an update is done by deploying all the stacks like so:

```bash
AWS_PROFILE=cloud-snitch-dev npx cdk deploy '*-dev'
```
