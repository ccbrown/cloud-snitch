#!/usr/bin/env python3
import boto3

ssm = boto3.client('ssm')

def get_region_ids():
    ret = []
    params = {
        'Path': '/aws/service/global-infrastructure/regions',
    }
    while True:
        resp = ssm.get_parameters_by_path(**params)
        for p in resp['Parameters']:
            ret.append(p['Value'])
        if 'NextToken' in resp:
            params['NextToken'] = resp['NextToken']
        else:
            break
    return ret


def get_region_info(region_id):
    ret = {}
    params = {
        'Path': f'/aws/service/global-infrastructure/regions/{region_id}',
    }
    while True:
        resp = ssm.get_parameters_by_path(**params)
        for p in resp['Parameters']:
            name = p['Name'].split('/')[-1]
            ret[name] = p['Value']
        if 'NextToken' in resp:
            params['NextToken'] = resp['NextToken']
        else:
            break
    return ret


"""
Locations are approximates, based on sources like...

https://github.com/turnkeylinux/aws-datacenters/blob/master/input/datacenters
https://gist.github.com/tobilg/ba6a5e1635478d13efdea5c1cd8227de

Some locations are just estimated based on the city name.
"""
locations = {
	'us-east-1':      [38.13, -78.45],
	'us-east-2':      [39.96, -83],
	'us-west-1':      [37.35, -121.96],
	'us-west-2':      [46.15, -123.88],
	'af-south-1':     [-33.93, 18.42],
	'ap-east-1':      [22.27, 114.16],
	'ap-south-2':     [17.4065, 78.4772],
	'ap-southeast-3': [-6.125, 106.655],
	'ap-southeast-5': [4.2105, 101.9758],
	'ap-southeast-4': [-37.8136, 144.9631],
	'ap-south-1':     [19.08, 72.88],
	'ap-northeast-2': [37.56, 126.98],
	'ap-northeast-3': [34.69, 135.49],
	'ap-southeast-1': [1.37, 103.8],
	'ap-southeast-2': [-33.86, 151.2],
	'ap-southeast-7': [15.87, 100.9925],
	'ap-northeast-1': [35.41, 139.42],
	'ca-central-1':   [45.5, -73.6],
	'ca-west-1':      [51.0447, -114.0719],
	'eu-central-1':   [50, 8],
	'eu-west-1':      [53, -8],
	'eu-west-2':      [51, -0.1],
	'eu-south-1':     [45.43, 9.29],
	'eu-west-3':      [48.86, 2.35],
	'eu-south-2':     [40.4637, -3.7492],
	'eu-north-1':     [59.25, 17.81],
	'eu-central-2':   [47.3769, 8.5417],
	'il-central-1':   [32.0853, 34.7818],
	'mx-central-1':   [19.4326, -99.1332],
	'me-south-1':     [26.10, 50.46],
	'me-central-1':   [23.4241, 53.8478],
	'sa-east-1':      [-23.34, -46.38],
	'us-gov-east-1':  [38.944, -77.455],
	'us-gov-west-1':  [37.618, -122.375],
	'cn-north-1':     [40.080, 116.584],
	'cn-northwest-1': [38.321667, 106.3925],
}

for id in get_region_ids():
    info = get_region_info(id)
    loc = locations[id]
    print(f'"{id}": {{"{info["longName"]}", "{info["geolocationCountry"]}", "{info["geolocationRegion"]}", "{info["partition"]}", {loc[0]}, {loc[1]}}},')
