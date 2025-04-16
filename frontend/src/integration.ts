const REVISION = 2;

export const INTEGRATION_TEMPLATE_URL = `${process.env.NEXT_PUBLIC_CDN_URL || ''}/integration-v${REVISION}.cfn.yaml`;
export const INTEGRATION_TEMPLATE_S3_URL = `https://s3.amazonaws.com/${process.env.NEXT_PUBLIC_PUBLIC_S3_BUCKET_NAME || ''}/integration-v${REVISION}.cfn.yaml`;
