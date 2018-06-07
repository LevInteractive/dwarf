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
# Allowed hosts
# * => Allow all hosts
# http://example.com => allows example.com domain (exactly, http and https are different)
# /\.example\.com/ => allows any subdmain at example.com (REGEX)
# 200.200.100.100 => allows a host defined by IP (for server side requests)
# http://example.com,/.example\.com/,200.200.100.100 => any of above, comma separated list
WHITELIST=*
```

**IMPORTANT CORS WARNING**: You need to check carefully your WHITELIST configuration
to avoid CORS errors on your application.

And your good to run:

```bash
node app.js
```

### Overriding `.env` file or no `.env` file at all

You can start directly using env vars (no `.env` file):

```bash
MONGO_CONNECTION_STRING=mongodb://127.0.0.1:27017  MONGO_DATABASE=dwarf \
  PORT=3001 BASE_URL=http://localhost:3001 API_KEY=RANDOM_LARGE_HASH \
  WHITELIST="*"
  node app.js
```

or just override one or more settings - the rest will use `.env` file settings:

```bash
  BASE_URL=myshrt.ul WHITELIST="example.com,/\.localhost:3555/" node app.js
```

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

## Example frontend module

The following code uses [request-promise](https://www.npmjs.com/package/request-promise)
module for async requests, so you need to install on your code:

```
npm install request-promise
```

Considering the module:

```javascript
/**
 * dwarf-shortener wrapper
 * @file dwarf-shortener.js
 */

import rp from "request-promise";

exports.shorten = async function(apiUrl, apiKey, longUrl, code) {
  const options = {
    method: "POST",
    uri: `${apiUrl}/create`,
    json: { longUrl, code, apiKey },
    headers: {
      Origin: window.location.origin
    }
  };

  return rp(options)
    .then(function(data) {
      return data.shortUrl;
    })
    .catch(function(err) {
      console.error("[DWARF SHORTENER ERROR]", err.message);
      return longUrl;
    });
};

exports.batchShorten = async function(apiUrl, apiKey, longUrls) {
  if (process.env.NODE_ENV === "test") {
    return await Promise.all(
      longUrls.map(async longUrl => {
        return { longUrl, shortUrl: longUrl };
      })
    );
  }

  const options = {
    method: "POST",
    uri: `${apiUrl}/create`,
    json: { longUrl: longUrls, apiKey },
    headers: {
      Origin: window.location.origin
    }
  };

  return rp(options)
    .then(function(data) {
      return data;
    })
    .catch(async function(err) {
      error("[DWARF SHORTENER BATCH ERROR]", err);
      return await Promise.all(
        longUrls.map(async longUrl => {
          return { longUrl, err: err.message };
        })
      );
    });
};
```

You can consume using:

```javascript
import { shorten, batchShorten } from "PATH_TO/dwarf-shortner";
const shortUrl = shorten(
  "http://myshrt.url",
  "MY_API_KEY",
  "http://longurl.example.com"
);
console.log(shortUrl);]
// # sample output
// http://myshrt.url/3Yc

const urls = batchShorten("http://myshrt.url", "MY_API_KEY", [
  "http://longurl1.example.com",
  "http://longurl2.example.com",
  "http://longurl3.example.com",
]);
console.log(urls);
// # sample output
// [
//   { "longUrl": "http://longurl1.example.com", "shortUrl": "http://myshrt.url/3Yc"},
//   { "longUrl": "http://longurl2.example.com", "shortUrl": "http://myshrt.url/3Yd"},
//   { "longUrl": "http://longurl3.example.com", "shortUrl": "http://myshrt.url/3Ye"}
// ]
```
