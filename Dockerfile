FROM golang:1.22-jammy AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app

FROM gcr.io/distroless/base-debian12

WORKDIR /app

USER nonroot:nonroot

COPY --from=build /app/app /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]