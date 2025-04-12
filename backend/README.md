# backend

## Dependencies

The only dependency of the backend is Go. You can install the latest version from [here](https://go.dev/).

## Testing

To run the tests, you'll need to run a local DynamoDB server. You can do this using Docker like so:

```bash
docker run -p 8000:8000 --rm -it amazon/dynamodb-local
```

Then, you can run the tests using `go test -v ./...`.
