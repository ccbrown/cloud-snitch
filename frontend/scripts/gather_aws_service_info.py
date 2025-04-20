#!/usr/bin/env python3
import boto3
import json
import time

orgs = boto3.client('organizations')
iam = boto3.client('iam')

roots = orgs.list_roots()
root = roots['Roots'][0]
entity_path = root['Arn'].split('/', 1)[1]

job = iam.generate_organizations_access_report(EntityPath=entity_path)
job_id = job['JobId']

params = {
    'JobId': job_id,
}

services = []

while True:
    report = iam.get_organizations_access_report(**params)
    if report['JobStatus'] == 'IN_PROGRESS':
        time.sleep(1)
    elif report['JobStatus'] == 'COMPLETED':
        for service in report['AccessDetails']:
            services.append({
                'name': service['ServiceName'],
                'namespace': service['ServiceNamespace'],
            })

        if report['IsTruncated']:
            params['Marker'] = report['Marker']
        else:
            break
    else:
        raise Exception(f"Unexpected job status: {report['JobStatus']}")

print(json.dumps(services, indent=2))
