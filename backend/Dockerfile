FROM golang:alpine as builder
WORKDIR /app

COPY . .

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /build ./cmd

FROM debian AS runner
WORKDIR /

COPY --from=builder /build /build
COPY .env .env

EXPOSE 8082

ENTRYPOINT ["/build"]