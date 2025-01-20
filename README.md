# Simple API Proxy

This is a very basic API proxy that adds relaxed CORS headers to API calls. This allows easier use of direct API calls in web browser javascript.

It also hides HTTP `Authorization: Bearer` access tokens to keep them out of that javascript.

## Usage

Configuration is set via environment variables. For example:

```
PORT=8787 API_HOSTNAME=docs.getgrist.com ACCESS_TOKEN=16274320ba0fd00dee589af6ebd21d5c664b0e3e go run main.go 
```

If you were trying to use this REST API endpoint:

```
https://docs.getgrist.com/api/docs/t00TDwkP4gQ3/tables/Locations/records
```

You could use this instead:

```
http://localhost:8787/api/docs/t00TDwkP4gQ3/tables/Locations/records
```

Because this proxy adds CORS headers, could then use something like the following in client-side javascript:

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
