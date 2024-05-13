FROM golang:latest as build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -race -o /app .

FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates
COPY --from=build /app /app

ENTRYPOINT ["/app"]
