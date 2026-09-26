FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /simple-commerce .

FROM alpine:3.22
RUN addgroup -S app && adduser -S -G app app
COPY --from=builder /simple-commerce /usr/local/bin/simple-commerce
USER app
EXPOSE 8080
ENTRYPOINT ["simple-commerce"]
