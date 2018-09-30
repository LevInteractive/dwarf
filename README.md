# dwarf

A high-throughput URL shortener microservice built with Go.

* gRPC for communication
* Redis store out of the box
* Fast & simple

### Usage

See [start-dev.sh](start-dev.sh) for a complete list of available environmental variables.

#### GET `/{short-hash}` -> 301 redirection

Dwarf will deliver a 301 redirection to the destination URL or redirect to the fallback
URL specified with [`NOTFOUND_REDIRECT_URL`](start-dev.sh).

#### Creating short links

You must communicate with dwarf via gRPC in order to generate new shortened URLs.

```proto
service Dwarf {
	rpc Create(CreateRequest) returns (CreateResponse) {}
}

message CreateRequest {
	repeated string urls = 1;
}

message CreateResponse {
	repeated string urls = 2;
}
```

Your response will return a set of shortened urls in the same order:

```json
// -> Request
{ "urls": ["http://long-url.com/1", "http://long-url.com/2"] }

// -> Response
{ "urls": ["http://sh.ort/Mp", "http://sh.ort/uJ"] }
```

#### A dwarf gRPC client written with node.js

To generate short urls, use a gRPC client such as this [node client](https://github.com/LevInteractive/dwarf-client-javascript).


# Development

* [Protobuf/gRPC is required.](https://grpc.io/docs/quickstart/go.html)
* Go & dep


### Redis Store

Spin up an instance of redis with:

```bash
docker run -p "6379:6379" --rm --name dwarf-redis redis:4-alpine
```

### Testing

`go test github.com/LevInteractive/dwarf/ -v`

Note that the tests rely on a running redis instance.
