 # Typesense Demo

> [Typesense](https://typesense.org/) companies search demo application, written using [Golang](https://go.dev/), [Vue.js](https://vuejs.org/) and [Axios](https://axios-http.com/).

## Usage

- First, start the provided [docker compose stack](docker-compose.yaml)

    ```shell
    docker compose up -d
    ```

- Then, populate the companies collection with sample data by calling [http://localhost:8080/populate](http://localhost:8080/populate)


- Finally, you can access the application on [http://localhost:8080](http://localhost:8080) and try to search companies.

## Configuration

You can configure the application in the [docker-compose.yaml](docker-compose.yaml) file, on the environment variables section:
- `TYPESENSE_HOST`: Typesense host
- `TYPESENSE_PORT`: Typesense port
- `TYPESENSE_API_KEY`: Typesense API key
