FROM golang:1.15.6-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build ./cmd/server/server.go

FROM alpine as production

WORKDIR /app
COPY --from=builder /app/server .

ENTRYPOINT [ "./server" ]