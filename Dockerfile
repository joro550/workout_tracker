FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN GOOS=linux go build -o /out/workout

FROM debian:bookworm-slim

COPY --from=builder /out/ /out/

EXPOSE 8080

CMD ["/out/workout"]
