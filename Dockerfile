
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/MihaPecnik/order-maching-system
COPY . .

RUN go get -d -v
RUN go build -o /go/bin/order-maching-system


FROM alpine

COPY --from=builder /go/bin/order-maching-system /go/bin/order-maching-system

CMD ["/go/bin/order-maching-system", "-migrate", "-postgres_url", "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable", "-populate" ]