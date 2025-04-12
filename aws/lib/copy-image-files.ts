import {
    aws_ecr_assets as ecr_assets,
    aws_lambda as lambda,
    aws_s3 as s3,
    Duration,
    triggers,
    Stack,
} from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as path from 'path';

interface Props {
    image: ecr_assets.DockerImageAsset;
    directory: string;
    bucket: s3.IBucket;
    prefix: string;
}

export class CopyImageFiles extends Construct {
    private f: triggers.TriggerFunction;

    constructor(scope: Construct, id: string, props: Props) {
        super(scope, id);

        const stack = Stack.of(this);

        const dockerImageAsset = new ecr_assets.DockerImageAsset(this, 'DockerImage', {
            directory: path.join(__dirname, './copy-image-files'),
            buildArgs: {
                IMAGE: `${stack.account}.dkr.ecr.${stack.region}.amazonaws.com/${props.image.repository.repositoryName}:${props.image.imageTag}`,
                SRC: props.directory,
            },
            extraHash: stack.region,
            platform: ecr_assets.Platform.LINUX_ARM64,
        });
        dockerImageAsset.node.addDependency(props.image);

        this.f = new triggers.TriggerFunction(this, 'Trigger', {
            architecture: lambda.Architecture.ARM_64,
            handler: lambda.Handler.FROM_IMAGE,
            runtime: lambda.Runtime.FROM_IMAGE,
            environment: {
                BUCKET: props.bucket.bucketName,
                PREFIX: props.prefix,
            },
            code: new lambda.EcrImageCode(dockerImageAsset.repository, {
                tagOrDigest: dockerImageAsset.imageTag,
            }),
            memorySize: 2048,
            timeout: Duration.minutes(15),
        });
        props.bucket.grantReadWrite(this.f);
    }

    executeBefore(...scopes: Construct[]): void {
        this.f.executeBefore(...scopes);
    }
}
