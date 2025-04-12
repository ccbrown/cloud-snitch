import boto3
import mimetypes
import os

bucket = os.getenv('BUCKET')
prefix = os.getenv('PREFIX')
s3 = boto3.client('s3')

def handler(event, context):
    for path, subdirs, files in os.walk('/staging'):
        for file in files:
            file_path = os.path.join(path, file)
            object_key = os.path.join(prefix, file_path[len('/staging/'):])
            extra_args = {
                'CacheControl': 'max-age=31536000',
            }
            content_type, _ = mimetypes.guess_type(file_path)
            if content_type is not None:
                extra_args['ContentType'] = content_type
            print(f'uploading {file_path} to s3://{bucket}/{object_key}')
            s3.upload_file(file_path, bucket, object_key, ExtraArgs=extra_args)
