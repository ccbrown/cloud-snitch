# backend

## Dependencies

The only dependency of the backend is Go. You can install the latest version from [here](https://go.dev/).

## Testing

To run the tests, you'll need to run a local DynamoDB server. You can do this using Docker like so:

```bash
docker run -p 8000:8000 --rm -it amazon/dynamodb-local
```

Then, you can run the tests using `go test -v ./...`.

## Code Layout

- [api](api): The API which the frontend uses. This is a thin layer on top of the business logic.
- [app](app): The business logic.
- [cmd](cmd): The CLI entrypoints for the application.
- [geoip](geoip): A package for looking up geolocation information for IP addresses.
- [model](model): The data models used internally by the application.
- [report](report): Pulls data from CloudTrail logs and generates reports based on it.
- [store](store): The data store for the application. This is a thin layer on top of DynamoDB.
