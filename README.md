# Dwarf

Crazy fast URL shortener written in Go using a Redis as a store while exposing
an interface for others.

Meant to be used as a microservice.

### API

**GET** `/{short-hash}`

This will result in a 301 redirect or a 404 not found.

**POST** `/create`

Send a JSON payload with the following shape:

```json
{
  "urls": ["http://my-url.com", "http://my-other-url.io"]
}
```

Your response will return a set of shortened urls in the same order:

```json
{
  "urls": ["http://sh.ort/Mp", "http://sh.ort/uJ"]
}
```


Better docs coming.
