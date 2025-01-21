# Simple API Proxy

This is a very basic API proxy that adds relaxed CORS headers to API calls. This allows easier use of direct API calls in web browser javascript.

It also hides HTTP `Authorization: Bearer` access tokens to keep them out of that javascript.

## Usage

Configuration is set via environment variables. For example:

```
LISTEN=localhost \
PORT=8787 \
API_HOSTNAME=docs.getgrist.com \
ACCESS_TOKEN=16274320ba0fd00dee589af6ebd21d5c664b0e3e \
./bin/apiproxy
```

The following environment variables are available:

`LISTEN` - The network address to listen on. Defaults to `0.0.0.0` if not set. Set this to `localhost` for testing or local development.

`PORT` - The port to listen on. Defaults to `8787`.

`API_HOSTNAME` - The hostname of the upstream API server.

`ACCESS_TOKEN` - The `Authorization: Bearer` token for the API server. This is currently the only supported authorization method.

### Container image

A [container image for apiproxy](https://github.com/freethoughtdesign/apiproxy/pkgs/container/apiproxy) is available at `ghcr.io/freethoughtdesign/apiproxy`, so this example will also work:

```
docker run --name apiproxy --rm -d \
  -p 8787:8787 \
  -e API_HOSTNAME=docs.getgrist.com \
  -e ACCESS_TOKEN=16274320ba0fd00dee589af6ebd21d5c664b0e3e \
  ghcr.io/freethoughtdesign/apiproxy:latest
```

Or in a `compose.yml` if using Docker Compose:

```
services:
  apiproxy:
    image: ghcr.io/freethoughtdesign/apiproxy:latest
    ports:
      - 8787:8787
    environment:
      API_HOSTNAME: docs.getgrist.com
      ACCESS_TOKEN: 16274320ba0fd00dee589af6ebd21d5c664b0e3e
```

## Building

```
go build -o ./bin/apiproxy main.go
```

## Development

You can quickly test during development with a command like this:

```
LISTEN=localhost PORT=8787 API_HOSTNAME=docs.getgrist.com ACCESS_TOKEN=16274320ba0fd00dee589af6ebd21d5c664b0e3e go run main.go 
```

You can also use Docker Compose to build and run the container from the `compose.yml` file in this repo. Copy `.env.example` to `.env` and modify it. Then run:

```
docker compose up --build -d
```


## Example

Let's say you were trying to use this REST API endpoint from a client-side, single-page javascript application:

```
https://docs.getgrist.com/api/docs/t00TDwkP4gQ3/tables/Locations/records
```

But the API service doesn't support direct use from browser javascript, and even if it did, you would have to include the `Authorization: Bearer` token in javascript code which would reveal it to anyone who looked at the source code.

You could instead run the proxy like so:

```
API_HOSTNAME=docs.getgrist.com ACCESS_TOKEN=16274320ba0fd00dee589af6ebd21d5c664b0e3e ./bin/apiproxy 
```

And access the API via:

```
http://localhost:8787/api/docs/t00TDwkP4gQ3/tables/Locations/records
```

Because this proxy adds CORS headers, you could then call the API in client-side javascript:

```javascript
async function getData() {
  const url = "http://localhost:8787/api/docs/t00TDwkP4gQ3/tables/Locations/records";

  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }

    const json = await response.json();
    console.log(json);
  } catch (error) {
    console.error(error.message);
  }
}
```

## Acknowledgements

- This [article on hashnode.dev](https://tobiojuolape.hashnode.dev/implementing-a-reverse-proxy-api-in-go) got me pointed the right direction.
