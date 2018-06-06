# dwarf

A URL shortener application optimized for performance with basic security.

## Install

```bash
git clone https://github.com/LevInteractive/dwarf.git
cd dwarf
npm install
```

## Configure

Create `.env` file by copying `.env.sample`.

```bash
cp .env.sample .env
```

So setup as needed:

```bash
# Mongo Connection
MONGO_CONNECTION_STRING=mongodb://127.0.0.1:27017
MONGO_DATABASE=dwarf
# Port to run
PORT=3001
# Base (public) url
# Base (public) url - that's your real domain in nginx or something like that
# BASE_URL=https://shortu.rl
BASE_URL=http://localhost:3001
# Api key to access
API_KEY=RANDOM_LARGE_HASH
```

And your good to run:

```bash
node app.js
```

### No `.env` file

You can start directly using env vars:

```bash
MONGO_CONNECTION_STRING=mongodb://127.0.0.1:27017  MONGO_DATABASE=dwarf \
  PORT=3001 BASE_URL=http://localhost:3001 API_KEY=RANDOM_LARGE_HASH \
  node app.js
```

But remember that `.env` file has precendence and will be used instead if available.

## Consuming

### Shortening a single url

You're free to send a request url twice - server checks if the url already
exists and just return the `shortUrl` for you.

Send a json post request to `BASE_URL/create` with json data
`{"longUrl": "http://example.com", "apiKey": "YOUR_API_KEY" }`:

```bash
curl -X POST \
  http://localhost:3001/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -H 'postman-token: c11a9eea-f7bf-f1e9-5ca3-1ecc5197168d' \
  -d '{"longUrl": "https://lev-interactive.com", "apiKey": "YOUR_API_KEY"}'
```

Returns:

```javascript
{
  "longUrl": "https://lev-interactive.com",
  "shortUrl": "http://localhost:3001/3YA"
}
```

Or 404 with errors

```javascript
// Bad Json
{
    "error": true,
    "message": "Bad JSON"
}
// Invalid apiKey
{
  "error": true,
  "message": "You need to send a valid apiKey"
}
// Not a string
{
  "error": true,
  "message": "longUrl is not a string"
}
// No longUrl
{
  error: true,
  "message": "You need to send a longUrl to be shorten"
}
// Invalid URL
{
  "error": true,
  "message": "Invalid URL format. Input URL must comply to the following: http(s)://(www.)domain.ext(/)(path)"
}
```

You can also force a `code`:

```bash
curl -X POST \
  http://localhost:3001/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -H 'postman-token: c11a9eea-f7bf-f1e9-5ca3-1ecc5197168d' \
  -d '{"longUrl": "https://lev-interactive.com", "apiKey": "YOUR_API_KEY", "code": "lev"}'
```

Returns:

```javascript
{
  "longUrl": "https://lev-interactive.com",
  "shortUrl": "http://localhost:3001/lev"
}
```

### Batch sortening urls

Send an array of urls on `longUrl`:

```bash
curl -X POST \
  http://localhost:3001/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -H 'postman-token: c11a9eea-f7bf-f1e9-5ca3-1ecc5197168d' \
  -d '{"longUrl": ["https://lev-interactive.com", "https://github.com", "htttp://duckduckgo.com"], "apiKey": "YOUR_API_KEY"}'
```

Returns (see that errors are returned on each url block so you can parse):

```javascript
[
  {
    longUrl: "https://lev-interactive.com",
    shortUrl: "http://localhost:3001/3YA"
  },
  {
    longUrl: "https://github.com",
    shortUrl: "http://localhost:3001/3YA"
  },
  {
    longUrl: "htttp://duckduckgo.com",
    error: true,
    message:
      "Invalid URL format. Input URL must comply to the following: http(s)://(www.)domain.ext(/)(path)"
  }
];
```

So far, using `code` to batch create custom urls are not available.
