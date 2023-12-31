# 📦 Shipment Generator
Sample REST API for handling shipments for given quantity and sizes.
The API aims to optimize the number of items and packages shipped (e.g. minimum of excess items and minimum number of packages), prioritizing items over packages.

## Test Live
Navigate to [https://shipment-generator-api-zfbs2ytkgq-uc.a.run.app/](https://shipment-generator-api-zfbs2ytkgq-uc.a.run.app/)

## Run Locally

On your Terminal, run 
```
docker compose up
```

### Unit tests

Run tests with
```
docker exec -it shipment-generator go test
```

> Note: Docker image must be up and running to run the tests.

## Test on Browser

While the Docker container is running, you can visit [http://localhost:8080](http://localhost:8080) to test the app. You will be presented with an interface like below:

<img src="/docs/img/ui-1.png" width="320">

You can test the app by changing the input. Once you click on **Submit**, you will be redirected to the order page where a summary of your shipment is present:

<img src="/docs/img/ui-2.png" width="320">

## Test with cURL

You can call the `/api/order` endpoint with a GET request from your Terminal, i.e.:
```
curl -X GET "http://localhost:8080/api/order?quantity=1235&sizes=250,500,1000,2000,5000"
```

Sample response:

<img src="/docs/img/terminal-1.png" width="480">


## Note on Pack Sizes
The pack sizes are decided in the following hierarchical order:
1. If the pack sizes are specified as a comma-separated list in the request
2. If the pack sizes are specified in `DEFAULT_PACK_SIZES` environment variable (you can change it in `docker-compose.yaml` to play around)
3. If both are not present, a default array of sizes `250,500,1000,2000,5000` are used.
