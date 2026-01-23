FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./cmd/app

FROM gcr.io/distroless/static-debian12

WORKDIR /app

USER nonroot:nonroot

COPY --from=build /app/app /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]