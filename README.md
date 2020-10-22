## Getting started

Build the image:

```
docker build -t pack-calc-go .
```

Run the dev container with hot reload testing:

```
docker run --name pack-calc-go-dev -v $(pwd):/app pack-calc-go
```

Reference the container name with `docker exec` to run arbitary commands, e.g. format all files:

```
docker exec -it pack-calc-go-dev gofmt -s -w .
```

## Microservice deployment

The [Serverless framework](https://serverless.com/) is used for a microservice deployment:

```
npm install -g serverless
serverless config credentials --provider aws --key <key> --secret <secret>
```

To build and deploy the function:

```
docker exec -it pack-calc-go-dev make build
make deploy
```

The API Gateway endpoint and API key will then be printed to console. An example request would be:

```
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: example_key" \
  -d '{"quantity": 12001, "packSizes": [250,500,1000,2000,5000]}' \
  https://example.execute-api.eu-west-1.amazonaws.com/dev/calculate
```

The service responds with a JSON object describing the number of required pack sizes:

```
{"250":1,"2000":1,"5000":2}
```