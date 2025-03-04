FROM golang:1.20-alpine AS builder

RUN apk add --no-cache make

WORKDIR /app

COPY . .

RUN make build

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/server /bin/server

EXPOSE 8080

CMD ["/bin/server"]

