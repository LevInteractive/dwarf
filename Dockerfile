FROM golang:1.11.4-alpine

WORKDIR /src/github.com/LevInterctive/dwarf
ADD . /src/github.com/LevInterctive/dwarf

RUN apk add bash ca-certificates git gcc g++ libc-dev

RUN go build -o bin/dwarf
CMD ["./bin/dwarf"]
