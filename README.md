 # Typesense Demo

> [Typesense](https://typesense.org/) companies search demo application, written using [Golang](https://go.dev/), [Vue.js](https://vuejs.org/) and [Axios](https://axios-http.com/).

![Screenshot](/home/jonathan/Dev/ankorstore/sandbox/typesense-go-demo/doc/screenshot.png)

## Usage

First, start the provided [docker compose stack](docker-compose.yaml)

```shell
docker compose up -d
```

Then head to, you can access the application on [http://localhost](http://localhost) to:
- first populate the companies collection
- and then try to search for them

## Configuration

You can configure the application in the [docker-compose.yaml](docker-compose.yaml) file, on the environment variables section:
- `TYPESENSE_HOST`: Typesense host
- `TYPESENSE_PORT`: Typesense port
- `TYPESENSE_API_KEY`: Typesense API key
