# Dwarf

A performant URL shortener written in Go using a Redis as a store while exposing
an interface for others. Uses rGPC for communication.

### How it works

**GET** `/{short-hash}`

This will result in a 301 redirect or a 404 not found.

**Creating new shortened URL's via gRPC**

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

# Development

* [Protobuf/gRPC is required.](https://grpc.io/docs/quickstart/go.html)
* Go & dep


**Redis Store**
Spin up an instance of redis with:

```bash
docker run -p "6379:6379" --rm --name dwarf-redis redis:4-alpine
```

**Testing**

`go test github.com/LevInteractive/dwarf/ -v`
